import json
import logging
from datetime import datetime, timezone
from typing import List, Optional

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


class SBERTFeedbackConsumerConfig(pydantic.BaseModel):
    articles_collection: str
    qdrant_url: str
    qdrant_api_key: Optional[str]
    feedback_retention_days: int
    redis_sentinels: str
    redis_master_name: str
    redis_password: Optional[str]
    redis_db: int
    redis_user_profile_prefix: str
    top_k_articles: int


class SBERTFeedbackConsumer(InferenceConsumer):
    def __init__(
        self,
        model: MlflowModel,
        producer: InferenceProducer,
        logger: logging.Logger,
        consumer_config: InferenceConsumerConfig,
        config: SBERTFeedbackConsumerConfig,
    ):
        super().__init__(
            logger,
            consumer_config,
            dict_to_msg=lambda d: parse_obj_as(FeedbackAggregate, d),
        )
        self.model = model
        self.producer = producer
        self.qdrant = QdrantClient(
            url=config.qdrant_url,
            api_key=config.qdrant_api_key,
            prefer_grpc=False,
        )
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
        """Process a batch of feedback records and store them in Qdrant."""
        try:
            for feedback in batch:
                if feedback.value == 1:
                    try:
                        redis_key = f"{self.config.redis_user_profile_prefix}:{feedback.user_id}"

                        existing_profile = self.redis_client.get(redis_key)

                        if existing_profile:
                            profile_data = json.loads(existing_profile)
                            article_ids = profile_data.get("article_ids", [])
                            article_ids.insert(0, feedback.news_id)
                            article_ids = article_ids[: self.config.top_k_articles]
                        else:
                            article_ids = [feedback.news_id]

                        self.logger.info(
                            f"Updating profile for user {feedback.user_id} with {len(article_ids)} articles: {article_ids}"
                        )

                        article_vectors = self._fetch_article_vectors(article_ids)
                        if article_vectors is None or len(article_vectors) == 0:
                            self.logger.warning(
                                f"No article vectors found for user {feedback.user_id}"
                            )
                            continue

                        mean_vector = np.mean(article_vectors, axis=0)

                        redis_value = {
                            "vector": mean_vector.tolist(),
                            "article_ids": article_ids,
                            "updated_at": datetime.now(timezone.utc).isoformat(),
                            "num_articles": len(article_ids),
                        }
                        self.redis_client.set(
                            redis_key,
                            json.dumps(redis_value),
                            ex=self.config.feedback_retention_days
                            * 86400,  # TTL in seconds
                        )
                        self.logger.info(
                            f"Updated mean vector for user {feedback.user_id} "
                            f"based on {len(article_ids)} articles"
                        )
                    except RedisError as e:
                        self.logger.error(
                            f"Redis error updating user {feedback.user_id} profile: {e}"
                        )
                    except Exception as e:
                        self.logger.error(
                            f"Error updating mean vector for user {feedback.user_id}: {e}",
                            exc_info=True,
                        )
                else:
                    self.logger.info(
                        f"Skipping negative feedback "
                        f"(user: {feedback.user_id}, article: {feedback.news_id}, value: {feedback.value})"
                    )

            self.logger.info(f"Processed batch of {len(batch)} feedbacks")

        except Exception as e:
            self.logger.error(f"Error processing batch: {e}", exc_info=True)

    def _fetch_article_vectors(self, article_ids: List[int]) -> Optional[np.ndarray]:
        try:
            points = self.qdrant.retrieve(
                collection_name=self.config.articles_collection,
                ids=article_ids,
                with_vectors=True,
            )

            return np.array([point.vector for point in points])

        except Exception as e:
            self.logger.error(f"Error fetching article vectors: {e}")
            return None

    def get_user_profile(self, user_id: int) -> Optional[dict]:
        try:
            redis_key = f"{self.config.redis_user_profile_prefix}:{user_id}"
            value = self.redis_client.get(redis_key)

            if value:
                return json.loads(value)
            return None
        except Exception as e:
            self.logger.error(f"Error retrieving user profile for {user_id}: {e}")
            return None
