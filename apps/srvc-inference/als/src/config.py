import os

import pydantic


class Config(pydantic.BaseModel):
    # General
    environment: str = os.getenv("ENVIRONMENT", "development")
    log_level: str = os.getenv("LOG_LEVEL", "INFO")

    # Kafka
    kafka_bootstrap_servers: str = os.getenv("KAFKA_BROKERS", "localhost:9092")
    kafka_batch_size: int = int(os.getenv("KAFKA_BATCH_SIZE", "50"))
    ## Producer
    kafka_producer_topic: str = os.getenv(
        "KAFKA_INFERENCE_COMMAND_TOPIC", "inference-command"
    )
    ## Article Consumer
    kafka_article_consumer_topic: str = os.getenv(
        "KAFKA_ARTICLE_AGGREGATE_TOPIC", "articles-aggregate"
    )
    kafka_article_consumer_group: str = os.getenv(
        "KAFKA_ARTICLE_CONSUMER_GROUP", "article-als-group"
    )
    ## Feedback Consumer
    kafka_feedback_consumer_topic: str = os.getenv(
        "KAFKA_FEEDBACK_AGGREGATE_TOPIC", "feedback-aggregate"
    )
    kafka_feedback_consumer_group: str = os.getenv(
        "KAFKA_FEEDBACK_CONSUMER_GROUP", "feedback-als-group"
    )

    # Qdrant
    qdrant_url: str = os.getenv("QDRANT_URL", "http://localhost:6333")

    # Model
    model_uri: str = os.getenv("MODEL_URI", "")
    model_name: str = os.getenv("MODEL_NAME", "als")
    tracking_uri: str = os.getenv("MLFLOW_TRACKING_URI", "http://localhost:5000")

    # API
    api_port: int = int(os.getenv("API_PORT", "8000"))
    api_host: str = os.getenv("API_HOST", "0.0.0.0")
