import json
import logging
import os
from datetime import datetime, timezone
from typing import Dict, List, Optional, Tuple

import numpy as np
import pydantic
from pydantic import parse_obj_as
from qdrant_client import QdrantClient
from redis.exceptions import RedisError
from redis.sentinel import Sentinel

from abstract.consumer import InferenceConsumer, InferenceConsumerConfig
from abstract.message import FeedbackAggregate
from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer


class ALSFeedbackConsumerConfig(pydantic.BaseModel):
    articles_collection: str = os.getenv("QDRANT_ARTICLES_COLLECTION", "articles")
    qdrant_url: str = os.getenv("QDRANT_URL", "http://localhost:6333")
    qdrant_api_key: Optional[str] = os.getenv("QDRANT_API_KEY", None)
    feedback_retention_days: int = int(os.getenv("FEEDBACK_RETENTION_DAYS", "30"))
    redis_sentinels: str = os.getenv(
        "REDIS_SENTINELS", "localhost:26379,localhost:26380,localhost:26381"
    )
    redis_master_name: str = os.getenv("REDIS_MASTER_NAME", "mymaster")
    redis_password: Optional[str] = os.getenv("REDIS_PASSWORD", None)
    redis_db: int = int(os.getenv("REDIS_DB", "0"))
    redis_user_profile_prefix: str = os.getenv("REDIS_USER_PROFILE_KEY", "user_profile")
    redis_interaction_prefix: str = os.getenv("REDIS_INTERACTION_KEY", "user_interactions")
    als_factors: int = int(os.getenv("ARTICLE_VECTOR_FEATURES", "50"))
    als_regularization: float = float(os.getenv("ALS_REGULARIZATION", "0.01"))
    als_iterations: int = int(os.getenv("ALS_ITERATIONS", "5"))


class ALSFeedbackConsumer(InferenceConsumer):
    def __init__(
        self,
        model: MlflowModel,
        producer: InferenceProducer,
        logger: logging.Logger,
        consumer_config: InferenceConsumerConfig,
        config: ALSFeedbackConsumerConfig,
    ):
        super().__init__(
            logger,
            consumer_config,
            dict_to_msg=lambda d: parse_obj_as(FeedbackAggregate, d),
        )
        self.model = model
        self.producer = producer
        self.qdrant = QdrantClient(url=config.qdrant_url, api_key=config.qdrant_api_key)
        self.config = config
        self._init_redis()

    def _init_redis(self) -> None:
        try:
            sentinel_list: list[tuple[str, int]] = [
                (host, int(port))
                for s in self.config.redis_sentinels.split(",")
                for host, port in [s.strip().split(":")]
            ]

            self.sentinel = Sentinel(
                sentinel_list,
                socket_timeout=5.0,
                password=self.config.redis_password,
            )

            self.redis_client = self.sentinel.master_for(
                self.config.redis_master_name,
                socket_timeout=5.0,
                password=self.config.redis_password,
                db=self.config.redis_db,
            )

            self.logger.info("Redis Sentinel connection established")
        except Exception as e:
            self.logger.error(f"Failed to initialize Redis Sentinel: {e}")
            raise

    def health(self) -> bool:
        """Check health of consumer, Qdrant, and Redis connections."""
        if not super().health():
            return False
        try:
            self.qdrant.get_collections()
            self.redis_client.ping()
            return True
        except Exception as e:
            self.logger.error(f"Health check failed: {e}")
            return False

    def process(self, batch: List[FeedbackAggregate]) -> None:
        """
        Process a batch of feedback records.
        Update user-item interactions and recompute user latent factors using ALS.
        """
        try:
            # Group feedback by user
            user_feedbacks: Dict[int, List[FeedbackAggregate]] = {}
            for feedback in batch:
                if feedback.user_id not in user_feedbacks:
                    user_feedbacks[feedback.user_id] = []
                user_feedbacks[feedback.user_id].append(feedback)

            # Update each user's profile
            for user_id, feedbacks in user_feedbacks.items():
                try:
                    self._update_user_profile(user_id, feedbacks)
                except Exception as e:
                    self.logger.error(
                        f"Error updating user {user_id} profile: {e}",
                        exc_info=True,
                    )

            self.logger.info(f"Processed batch of {len(batch)} feedbacks")

        except Exception as e:
            self.logger.error(f"Error processing batch: {e}", exc_info=True)

    def _update_user_profile(
        self, user_id: int, feedbacks: List[FeedbackAggregate]
    ) -> None:
        """Update user profile with new feedback using ALS update."""
        try:
            # Load existing interactions
            interactions = self._load_user_interactions(user_id)

            # Add new interactions
            for feedback in feedbacks:
                # Store interaction with implicit feedback value (1 for positive)
                if feedback.value == 1:
                    interactions[feedback.news_id] = {
                        "rating": 1.0,
                        "timestamp": datetime.now(timezone.utc).isoformat(),
                    }
                    self.logger.info(
                        f"Added interaction for user {user_id}, article {feedback.news_id}"
                    )
                else:
                    # For negative feedback, we can either ignore or use implicit negative
                    # Here we'll just log it but not add to interactions
                    self.logger.info(
                        f"Skipping negative feedback "
                        f"(user: {user_id}, article: {feedback.news_id})"
                    )

            if not interactions:
                self.logger.warning(f"No positive interactions for user {user_id}")
                return

            # Save interactions to Redis
            self._save_user_interactions(user_id, interactions)

            # Compute new user latent factor using ALS
            user_factor = self._compute_user_factor(user_id, interactions)

            if user_factor is not None:
                # Save user profile to Redis
                redis_key = f"{self.config.redis_user_profile_prefix}:{user_id}"
                redis_value = {
                    "vector": user_factor.tolist(),
                    "num_interactions": len(interactions),
                    "updated_at": datetime.now(timezone.utc).isoformat(),
                }
                self.redis_client.set(
                    redis_key,
                    json.dumps(redis_value),
                    ex=self.config.feedback_retention_days * 86400,
                )
                self.logger.info(
                    f"Updated user factor for user {user_id} "
                    f"based on {len(interactions)} interactions"
                )

        except Exception as e:
            self.logger.error(f"Error updating user {user_id} profile: {e}", exc_info=True)

    def _load_user_interactions(self, user_id: int) -> Dict[int, dict]:
        """Load user's interaction history from Redis."""
        try:
            redis_key = f"{self.config.redis_interaction_prefix}:{user_id}"
            value = self.redis_client.get(redis_key)

            if value:
                return json.loads(value)
            return {}
        except Exception as e:
            self.logger.error(f"Error loading interactions for user {user_id}: {e}")
            return {}

    def _save_user_interactions(self, user_id: int, interactions: Dict[int, dict]) -> None:
        """Save user's interaction history to Redis."""
        try:
            redis_key = f"{self.config.redis_interaction_prefix}:{user_id}"
            self.redis_client.set(
                redis_key,
                json.dumps(interactions),
                ex=self.config.feedback_retention_days * 86400,
            )
        except Exception as e:
            self.logger.error(f"Error saving interactions for user {user_id}: {e}")

    def _compute_user_factor(
        self, user_id: int, interactions: Dict[int, dict]
    ) -> Optional[np.ndarray]:
        """
        Compute user latent factor using ALS formula:
        user_factor = (Item^T * Item + λI)^(-1) * Item^T * ratings
        
        Where Item is the matrix of item factors for items the user interacted with.
        """
        try:
            # Get article IDs
            article_ids = list(interactions.keys())
            
            # Fetch item factors from Qdrant
            item_factors = self._fetch_article_vectors(article_ids)
            
            if item_factors is None or len(item_factors) == 0:
                self.logger.warning(f"No item factors found for user {user_id}")
                return None

            # Build ratings vector (implicit feedback: all 1s for positive interactions)
            ratings = np.array([interactions[aid]["rating"] for aid in article_ids], dtype=np.float32)

            # ALS closed-form solution
            # user_factor = (Item^T @ Item + λI)^(-1) @ Item^T @ ratings
            regularization = self.config.als_regularization * np.eye(self.config.als_factors)
            
            gram_matrix = item_factors.T @ item_factors + regularization
            user_factor = np.linalg.solve(gram_matrix, item_factors.T @ ratings)

            return user_factor.astype(np.float32)

        except Exception as e:
            self.logger.error(f"Error computing user factor for user {user_id}: {e}", exc_info=True)
            return None

    def _fetch_article_vectors(self, article_ids: List[int]) -> Optional[np.ndarray]:
        """Fetch item latent factors from Qdrant."""
        try:
            points = self.qdrant.retrieve(
                collection_name=self.config.articles_collection,
                ids=article_ids,
                with_vectors=True,
            )

            if not points:
                return None

            return np.array([point.vector for point in points], dtype=np.float32)

        except Exception as e:
            self.logger.error(f"Error fetching article vectors: {e}")
            return None

    def get_user_profile(self, user_id: int) -> Optional[dict]:
        """Retrieve a specific user profile from Redis."""
        try:
            redis_key = f"{self.config.redis_user_profile_prefix}:{user_id}"
            value = self.redis_client.get(redis_key)

            if value:
                return json.loads(value)
            return None
        except Exception as e:
            self.logger.error(f"Error retrieving user profile for {user_id}: {e}")
            return None
