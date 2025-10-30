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


class TFIDFArticleConsumerConfig(pydantic.BaseModel):
    tf_idf_max_features: int = int(os.getenv("TFIDF_MAX_FEATURES", "12"))
    articles_collection: str = os.getenv("ARTICLES_COLLECTION", "articles")
    qdrant_url: str = os.getenv("QDRANT_URL", "http://localhost:6333")
    similarity_threshold: float = float(os.getenv("SIMILARITY_THRESHOLD", "0.7"))
    model_name: str = os.getenv("MODEL_NAME", "tfidf")
    redis_sentinels: str = os.getenv(
        "REDIS_SENTINELS", "localhost:26379,localhost:26380,localhost:26381"
    )
    redis_master_name: str = os.getenv("REDIS_MASTER_NAME", "mymaster")
    redis_password: Optional[str] = os.getenv("REDIS_PASSWORD", None)
    redis_db: int = int(os.getenv("REDIS_DB", "0"))
    redis_scan_batch_size: int = int(os.getenv("REDIS_SCAN_BATCH_SIZE", "100"))


class TFIDFArticleConsumer(InferenceConsumer):
    def __init__(
        self,
        model: MlflowModel,
        producer: InferenceProducer,
        logger: logging.Logger,
        consumer_config: InferenceConsumerConfig,
        config: TFIDFArticleConsumerConfig,
    ):
        super().__init__(
            logger,
            consumer_config,
            dict_to_msg=lambda d: parse_obj_as(ArticleAggregate, d),
        )
        self.model = model
        self.producer = producer
        self.qdrant = QdrantClient(url=config.qdrant_url)
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

            self.logger.info("✅ Redis Sentinel connection established")
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
                        size=self.config.tf_idf_max_features,
                        distance=Distance.COSINE,
                    ),
                )
                self.logger.info(
                    f"Collection '{collection_name}' created with vector size "
                    f"{self.config.tf_idf_max_features}."
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
        texts = [article.abstract for article in articles]
        vectors = self.model.transform(texts).toarray().tolist()

        points = [
            PointStruct(id=article.id, vector=vector, payload=article.dict())
            for article, vector in zip(articles, vectors, strict=False)
        ]

        self.qdrant.upsert(
            collection_name=self.config.articles_collection,
            points=points,
        )
        self.logger.info(f"✅ Upserted batch of {len(points)} articles")

        recommendations = self.recommend(articles, vectors)
        if recommendations:
            self.logger.info(f"Generating {len(recommendations)} recommendations")
            self.producer.produce(recommendations)

    def _delete_articles(self, articles: List[ArticleAggregate]) -> None:
        """Delete articles from Qdrant."""
        ids_to_delete = [int(article.id) for article in articles]
        self.qdrant.delete(
            collection_name=self.config.articles_collection,
            points_selector=ids_to_delete,
        )
        self.logger.info(f"🗑️  Deleted batch of {len(ids_to_delete)} articles")

    def recommend(
        self, new_articles: List[ArticleAggregate], article_vectors: List[List[float]]
    ) -> List[InferenceCommand]:
        """
        Generate recommendations by comparing new articles to user preference vectors from Redis.
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

            recommendations = self._generate_recommendations(
                new_articles, article_vectors, user_profiles
            )

            self.logger.info(
                f"Generated {len(recommendations)} recommendations for "
                f"{len(new_articles)} articles across {len(user_profiles)} users"
            )
            return recommendations

        except Exception as e:
            self.logger.error(f"Error in recommend: {e}", exc_info=True)
            return []

    def _load_user_profiles_from_redis(self) -> Dict[int, np.ndarray]:
        """
        Load all user profile vectors from Redis using SCAN for efficiency.
        Returns dict mapping user_id -> mean_vector.
        """
        user_profiles = {}
        cursor = 0
        pattern = "user_profile:*"

        try:
            while True:
                cursor, keys = self.redis_client.scan(
                    cursor=cursor,
                    match=pattern,
                    count=self.config.redis_scan_batch_size,
                )

                for key in keys:
                    try:
                        # Extract user_id from key (format: "user_profile:123")
                        user_id = int(key.decode("utf-8").split(":")[1])

                        # Get profile data
                        value = self.redis_client.get(key)
                        if value:
                            profile_data = json.loads(value)
                            vector = np.array(profile_data["vector"])
                            user_profiles[user_id] = vector

                    except Exception as e:
                        self.logger.warning(
                            f"Error loading profile from key {key}: {e}"
                        )
                        continue

                # Break when cursor returns to 0 (full scan complete)
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
        article_vectors: List[List[float]],
        user_profiles: Dict[int, np.ndarray],
    ) -> List[InferenceCommand]:
        """
        Generate recommendations by computing similarity between article vectors
        and user profile vectors.
        """
        recommendations = []

        for article, article_vector_list in zip(
            new_articles, article_vectors, strict=False
        ):
            article_vector = np.array(article_vector_list)

            for user_id, user_vector in user_profiles.items():
                try:
                    similarity = self._compute_cosine_similarity(
                        article_vector, user_vector
                    )

                    if similarity >= self.config.similarity_threshold:
                        recommendation = InferenceCommand(
                            model=self.config.model_name,
                            user_id=int(user_id),
                            article=article,
                            score=float(similarity),
                            date=datetime.now(timezone.utc),
                        )
                        recommendations.append(recommendation)

                        self.logger.info(
                            f"✅ Recommendation: User {user_id} - Article {article.id} "
                            f"(similarity: {similarity:.4f})"
                        )

                except Exception as e:
                    self.logger.error(
                        f"Error computing similarity for user {user_id}, "
                        f"article {article.id}: {e}"
                    )
                    continue

        return recommendations

    @staticmethod
    def _compute_cosine_similarity(vec1: np.ndarray, vec2: np.ndarray) -> float:
        """Compute cosine similarity between two vectors."""
        norm1 = np.linalg.norm(vec1)
        norm2 = np.linalg.norm(vec2)

        if norm1 == 0 or norm2 == 0:
            return 0.0

        return float(np.dot(vec1, vec2) / (norm1 * norm2))

    def get_user_profile(self, user_id: int) -> Optional[dict]:
        """Retrieve a specific user profile from Redis."""
        try:
            redis_key = f"user_profile:{user_id}"
            value = self.redis_client.get(redis_key)

            if value:
                return json.loads(value)
            return None
        except Exception as e:
            self.logger.error(f"Error retrieving user profile for {user_id}: {e}")
            return None
