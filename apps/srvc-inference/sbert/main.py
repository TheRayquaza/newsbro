import logging
import os
import signal
import sys

import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException

from abstract.consumer import InferenceConsumerConfig
from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer, InferenceProducerConfig
from sbert.src.article_consumer import (
    SBERTArticleConsumer,
    SBERTArticleConsumerConfig,
)
from sbert.src.feedback_consumer import (
    SBERTFeedbackConsumer,
    SBERTFeedbackConsumerConfig,
)

if __name__ == "__main__":
    # Load .env if not in production
    if os.getenv("ENVIRONMENT") != "production":
        print("Loading .env file for development environment")
        env_file = sys.argv[1] if len(sys.argv) > 1 else ".env"
        if os.path.exists(env_file):
            load_dotenv(env_file)
        elif len(sys.argv) > 1:
            raise FileNotFoundError(f"Environment file {env_file} not found")

    # Load configuration from env
    from sbert.src.config import Config

    config = Config()

    # Setup logging
    logging.basicConfig(
        level=(logging.INFO if config.log_level.upper() == "INFO" else logging.DEBUG),
        format="%(asctime)s [%(levelname)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
    logger = logging.getLogger(__name__)

    # FastAPI app
    app = FastAPI(
        title="Article Recommendation API - SBERT",
        description="SBERT based article recommendation API using Qdrant",
        version="0.0.0",
    )

    # Load model
    if os.getenv("ENVIRONMENT") == "production":
        model = MlflowModel(
            model_uri=config.model_uri,
            tracking_uri=config.tracking_uri,
        )
    else:
        model = MlflowModel(
            model_uri=sys.argv[2] if len(sys.argv) > 2 else config.model_uri,
            tracking_uri=config.tracking_uri,
            # from_pickle=True,
        )

    # Setup producer
    producer_config = InferenceProducerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_producer_topic=config.kafka_producer_topic,
    )
    producer = InferenceProducer(logger, producer_config)

    c1_consumer_config = InferenceConsumerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_consumer_topic=config.kafka_feedback_consumer_topic,
        kafka_consumer_group=config.kafka_feedback_consumer_group,
        batch_size=config.kafka_batch_size,
        batch_interval=config.kafka_batch_interval,
        thread_pool_size=config.thread_pool_size,
    )
    c1_config = SBERTFeedbackConsumerConfig(
        articles_collection=config.articles_collection,
        qdrant_url=config.qdrant_url,
        qdrant_api_key=config.qdrant_api_key,
        feedback_retention_days=config.feedback_retention_days,
        redis_sentinels=config.redis_sentinels,
        redis_master_name=config.redis_master_name,
        redis_password=config.redis_password,
        redis_db=config.redis_db,
        redis_user_profile_prefix=config.redis_user_profile_prefix,
        top_k_articles=config.top_k_articles,
    )
    c1 = SBERTFeedbackConsumer(model, producer, logger, c1_consumer_config, c1_config)

    c2_consumer_config = InferenceConsumerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_consumer_topic=config.kafka_article_consumer_topic,
        kafka_consumer_group=config.kafka_article_consumer_group,
        batch_size=config.kafka_batch_size,
        batch_interval=config.kafka_batch_interval,
        thread_pool_size=config.thread_pool_size,
    )
    c2_config = SBERTArticleConsumerConfig(
        model_name=config.model_name,
        articles_collection=config.articles_collection,
        article_vector_features=config.article_vector_features,
        qdrant_url=config.qdrant_url,
        qdrant_api_key=config.qdrant_api_key,
        similarity_threshold=config.similarity_threshold,
        redis_sentinels=config.redis_sentinels,
        redis_master_name=config.redis_master_name,
        redis_password=config.redis_password,
        redis_db=config.redis_db,
        redis_scan_batch_size=config.redis_scan_batch_size,
        redis_user_profile_prefix=config.redis_user_profile_prefix,
    )
    c2 = SBERTArticleConsumer(model, producer, logger, c2_consumer_config, c2_config)
    c2.bootstrap()

    # Run consumers
    c1.run()
    c2.run()

    def shutdown_handler(signum, frame):
        logger.info("Shutdown signal received. Shutting down consumers...")
        c1.shutdown()
        c2.shutdown()
        logger.info("Consumers shut down successfully.")
        sys.exit(0)

    signal.signal(signal.SIGINT, shutdown_handler)
    signal.signal(signal.SIGTERM, shutdown_handler)

    # Health check endpoint
    @app.get("/health")
    async def health_check():
        if not c1.health() or not c2.health():
            raise HTTPException(status_code=503, detail="Service Unhealthy")
        return {"status": "healthy"}

    uvicorn.run(app, host=config.api_host, port=config.api_port)
