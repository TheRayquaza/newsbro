import os
from typing import Optional

import pydantic


class Config(pydantic.BaseModel):
    # General
    environment: str = os.getenv("ENVIRONMENT", "development")
    log_level: str = os.getenv("LOG_LEVEL", "INFO")
    thread_pool_size: int = int(os.getenv("THREAD_POOL_SIZE", "4"))

    # Kafka
    kafka_bootstrap_servers: str = os.getenv(
        "KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"
    )
    kafka_batch_size: int = int(os.getenv("KAFKA_BATCH_SIZE", "64"))
    kafka_batch_interval: int = int(os.getenv("KAFKA_BATCH_INTERVAL", "2"))

    ## Producer
    kafka_producer_topic: str = os.getenv("KAFKA_PRODUCER_TOPIC", "inference-topic")

    ## Article Consumer
    kafka_article_consumer_topic: str = os.getenv(
        "KAFKA_ARTICLE_CONSUMER_TOPIC", "article-aggregate"
    )
    kafka_article_consumer_group: str = os.getenv(
        "KAFKA_ARTICLE_CONSUMER_GROUP", "article-sbert-group"
    )

    ## Feedback Consumer
    kafka_feedback_consumer_topic: str = os.getenv(
        "KAFKA_FEEDBACK_CONSUMER_TOPIC", "feedback-aggregate"
    )
    kafka_feedback_consumer_group: str = os.getenv(
        "KAFKA_FEEDBACK_CONSUMER_GROUP", "feedback-sbert-group"
    )

    # Qdrant
    qdrant_url: str = os.getenv("QDRANT_URL", "http://localhost:6333")
    qdrant_api_key: Optional[str] = os.getenv("QDRANT_API_KEY", None)
    articles_collection: str = os.getenv("QDRANT_ARTICLES_COLLECTION", "articles_sbert")

    # Redis
    feedback_retention_days: int = int(os.getenv("FEEDBACK_RETENTION_DAYS", "30"))
    redis_sentinels: str = os.getenv(
        "REDIS_SENTINELS", "localhost:26379,localhost:26380,localhost:26381"
    )
    redis_master_name: str = os.getenv("REDIS_MASTER_NAME", "mymaster")
    redis_password: Optional[str] = os.getenv("REDIS_PASSWORD", None)
    redis_db: int = int(os.getenv("REDIS_DB", "0"))
    redis_user_profile_prefix: str = os.getenv("REDIS_USER_PROFILE_KEY", "user:sbert")
    top_k_articles: int = int(os.getenv("TOP_K_ARTICLES", "10"))
    redis_scan_batch_size: int = int(os.getenv("REDIS_SCAN_BATCH_SIZE", "100"))

    # Model
    model_uri: str = os.getenv("MODEL_URI", "")
    tracking_uri: str = os.getenv(
        "MLFLOW_TRACKING_URI", "https://mlflow.internal.newsbro.cc"
    )
    model_name: str = os.getenv("MODEL_NAME", "sbert")
    similarity_threshold: float = float(os.getenv("MODEL_SIMILARITY_THRESHOLD", "0.7"))
    article_vector_features: int = int(os.getenv("ARTICLE_VECTOR_FEATURES", "384"))

    # API
    api_port: int = int(os.getenv("API_PORT", "8080"))
    api_host: str = os.getenv("API_HOST", "")
