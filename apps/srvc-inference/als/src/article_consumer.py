import json
import logging
import os
from datetime import datetime, timezone
from typing import Dict, List, Optional

import numpy as np
import pydantic
from pydantic import parse_obj_as
from qdrant_client import QdrantClient
from qdrant_client.http.models import Distance, PointStruct, VectorParams
from redis.exceptions import RedisError
from redis.sentinel import Sentinel

from abstract.consumer import InferenceConsumer, InferenceConsumerConfig
from abstract.message import ArticleAggregate, InferenceCommand
from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer


class ALSArticleConsumerConfig(pydantic.BaseModel):
    article_vector_features: int = int(os.getenv("ARTICLE_VECTOR_FEATURES", "50"))
    articles_collection: str = os.getenv("QDRANT_ARTICLES_COLLECTION", "articles")
    qdrant_url: str = os.getenv("QDRANT_URL", "localhost:6333")
    qdrant_api_key: Optional[str] = os.getenv("QDRANT_API_KEY", None)
    similarity_threshold: float = float(os.getenv("MODEL_SIMILARITY_THRESHOLD", "0.3"))
    model_name: str = os.getenv("MODEL_NAME", "als")
    redis_sentinels: str = os.getenv(
        "REDIS_SENTINELS", "localhost:26379,localhost:26380,localhost:26381"
    )
    redis_master_name: str = os.getenv("REDIS_MASTER_NAME", "mymaster")
    redis_password: Optional[str] = os.getenv("REDIS_PASSWORD", None)
    redis_db: int = int(os.getenv("REDIS_DB", "0"))
    redis_scan_batch_size: int = int(os.getenv("REDIS_SCAN_BATCH_SIZE", "100"))
    redis_user_profile_prefix: str = os.getenv("REDIS_USER_PROFILE_KEY", "user_profile")


class ALSArticleConsumer(InferenceConsumer):
    def __init__(
        self,
        model: MlflowModel,
        producer: InferenceProducer,
        logger: logging.Logger,
        consumer_config: InferenceConsumerConfig,
        config: ALSArticleConsumerConfig,
    ):
        super().__init__(
            logger,
            consumer_config,
            dict_to_msg=lambda d: parse_obj_as(ArticleAggregate, d),
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

    def bootstrap(self) -> None:
        """Initialize Qdrant collection if it doesn't exist."""
        collection_name = self.config.articles_collection

        try:
            existing_collections = {
                c.name for c in self.qdrant.get_collections().collections
            }

            if collection_name not in existing_collections:
                self.qdrant.create_collection(
                    collection_name=collection_name,
                    vectors_config=VectorParams(
                        size=self.config.article_vector_features,
                        distance=Distance.DOT,  # For ALS, dot product is more appropriate
                    ),
                )
                self.logger.info(
                    f"Collection '{collection_name}' created with vector size "
                    f"{self.config.article_vector_features}."
                )
        except Exception as e:
            self.logger.error(f"Error creating collection '{collection_name}': {e}")
            raise

    def process(self, batch: List[ArticleAggregate]) -> None:
        """Process a batch of articles - upsert active ones and delete inactive ones."""
        try:
            to_upsert = [article for article in batch if article.is_active]
            to_delete = [article for article in batch if not article.is_active]

            if to_upsert:
                self._upsert_articles(to_upsert)

            if to_delete:
                self._delete_articles(to_delete)

        except Exception as e:
            self.logger.error(f"Error processing batch: {e}", exc_info=True)

    def _upsert_articles(self, articles: List[ArticleAggregate]) -> None:
        """
        For ALS, we compute item latent factors from the model.
        If articles are new, we initialize them with the average item factor or random.
        """
        article_vectors = []
        
        for article in articles:
            # Try to get existing item factor from model
            item_vector = self.model.get_item_factor(article.id)
            
            if item_vector is None:
                # New article - initialize with average or random
                item_vector = self._initialize_item_vector()
            
            article_vectors.append(item_vector)

        points = [
            PointStruct(id=article.id, vector=vector.tolist(), payload=article.dict())
            for article, vector in zip(articles, article_vectors, strict=True)
        ]

        self.qdrant.upsert(
            collection_name=self.config.articles_collection,
            points=points,
        )
        self.logger.info(f"Upserted batch of {len(points)} articles")

        self.producer.produce(self.recommend(articles, article_vectors))

    def _initialize_item_vector(self) -> np.ndarray:
        """Initialize a new item vector with small random values."""
        return np.random.normal(0, 0.1, self.config.article_vector_features).astype(np.float32)

    def _delete_articles(self, articles: List[ArticleAggregate]) -> None:
        """Delete articles from Qdrant."""
        ids_to_delete = [int(article.id) for article in articles]
        self.qdrant.delete(
            collection_name=self.config.articles_collection,
            points_selector=ids_to_delete,
        )
        self.logger.warning(f"Deleted batch of {len(ids_to_delete)} articles")

    def recommend(
        self, new_articles: List[ArticleAggregate], article_vectors: List[np.ndarray]
    ) -> List[InferenceCommand]:
        """
        Generate recommendations by computing dot product between item and user latent factors.
        """
        if not new_articles:
            self.logger.debug("No articles to process for recommendations")
            return []

        try:
            user_profiles = self._load_user_profiles_from_redis()

            if not user_profiles:
                self.logger.warning(
                    "No user profiles found in Redis, skipping recommendations"
                )
                return []

            self.logger.info(f"Loaded {len(user_profiles)} user profiles from Redis")

            return self._generate_recommendations(
                new_articles, article_vectors, user_profiles
            )

        except Exception as e:
            self.logger.error(f"Error in recommend: {e}", exc_info=True)
            return []

    def _load_user_profiles_from_redis(self) -> Dict[int, np.ndarray]:
        """
        Load all user latent factors from Redis using SCAN for efficiency.
        Returns dict mapping user_id -> user_factor.
        """
        user_profiles = {}
        cursor = 0
        pattern = f"{self.config.redis_user_profile_prefix}:*"

        try:
            while True:
                cursor, keys = self.redis_client.scan(
                    cursor=cursor,
                    match=pattern,
                    count=self.config.redis_scan_batch_size,
                )

                for key in keys:
                    try:
                        user_id = int(key.decode("utf-8").split(":")[1])

                        value = self.redis_client.get(key)
                        if value:
                            profile_data = json.loads(value)
                            vector = np.array(profile_data["vector"], dtype=np.float32)
                            user_profiles[user_id] = vector

                    except Exception as e:
                        self.logger.warning(
                            f"Error loading profile from key {key}: {e}"
                        )
                        continue

                if cursor == 0:
                    break

        except RedisError as e:
            self.logger.error(f"Redis error loading user profiles: {e}")
        except Exception as e:
            self.logger.error(f"Error loading user profiles: {e}", exc_info=True)

        return user_profiles

    def _generate_recommendations(
        self,
        new_articles: List[ArticleAggregate],
        article_vectors: List[np.ndarray],
        user_profiles: Dict[int, np.ndarray],
    ) -> List[InferenceCommand]:
        """
        Generate recommendations by computing dot product between item factors
        and user factors (standard ALS prediction: user_factor @ item_factor.T).
        """
        if not new_articles or not article_vectors or not user_profiles:
            return []

        # Convert to matrix: shape (n_articles, n_factors)
        article_matrix = np.vstack(article_vectors).astype(np.float32)

        user_ids = list(user_profiles.keys())
        # shape (n_users, n_factors)
        user_matrix = np.vstack([user_profiles[uid] for uid in user_ids]).astype(np.float32)

        # Compute predictions: shape (n_articles, n_users)
        # prediction[i, j] = item_i @ user_j
        prediction_matrix = article_matrix @ user_matrix.T

        threshold = self.config.similarity_threshold

        # Find predictions above threshold
        above_threshold = prediction_matrix >= threshold
        article_indices, user_indices = np.where(above_threshold)

        current_time = datetime.now(timezone.utc)

        self.logger.info(
            f"Processing {len(article_indices)} recommendations "
            f"from {len(new_articles)} articles and {len(user_ids)} users"
        )

        recommendations = [
            InferenceCommand(
                model=self.config.model_name,
                user_id=int(user_ids[user_idx]),
                article=new_articles[art_idx],
                score=float(prediction_matrix[art_idx, user_idx]),
                date=current_time,
            )
            for art_idx, user_idx in zip(article_indices, user_indices, strict=True)
        ]

        if len(recommendations) != 0:
            self.logger.info(
                f"Generated {len(recommendations)} recommendations "
                f"(avg score: {sum(r.score for r in recommendations) / len(recommendations):.4f})"
            )
        else:
            self.logger.info(
                f"No recommendations generated above threshold {threshold:.4f}"
            )

        return recommendations

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
