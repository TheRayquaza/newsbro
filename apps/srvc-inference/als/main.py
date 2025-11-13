import logging
import os
import sys

import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException

from abstract.consumer import InferenceConsumerConfig
from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer, InferenceProducerConfig

if __name__ == "__main__":
    if os.getenv("ENVIRONMENT") != "production":
        print("Loading .env file for development environment")
        env_file = sys.argv[1] if len(sys.argv) > 1 else ".env"
        if os.path.exists(env_file):
            load_dotenv(env_file)
        elif len(sys.argv) > 1:
            raise FileNotFoundError(f"Environment file {env_file} not found")

    from als.src.article_consumer import (
        ALSArticleConsumer,
        ALSArticleConsumerConfig,
    )
    from als.src.config import Config  # because Config evaluation is done on runtime
    from als.src.feedback_consumer import (
        ALSFeedbackConsumer,
        ALSFeedbackConsumerConfig,
    )

    config = Config()

    logging.basicConfig(
        level=(logging.INFO if config.log_level.upper() == "INFO" else logging.DEBUG),
        format="%(asctime)s [%(levelname)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
    logger = logging.getLogger(__name__)

    app = FastAPI(
        description="ALS (Alternating Least Squares) based article recommendation API using Qdrant",
        version="1.0.0",
    )

    if os.getenv("ENVIRONMENT") == "production":
        model = MlflowModel(
            model_uri=config.model_uri,
            tracking_uri=config.tracking_uri,
        )
    else:
        model = MlflowModel(
            model_uri=sys.argv[2] if len(sys.argv) > 2 else config.model_uri,
            tracking_uri=config.tracking_uri,
            from_pickle=True,
        )

    producer_config = InferenceProducerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_producer_topic=config.kafka_producer_topic,
    )
    producer = InferenceProducer(logger, producer_config)

    # Feedback Consumer (processes user interactions and updates user factors)
    c1_consumer_config = InferenceConsumerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_consumer_topic=config.kafka_feedback_consumer_topic,
        kafka_consumer_group=config.kafka_feedback_consumer_group,
        batch_size=config.kafka_batch_size,
    )
    c1_config = ALSFeedbackConsumerConfig()
    c1 = ALSFeedbackConsumer(model, producer, logger, c1_consumer_config, c1_config)

    # Article Consumer (processes new articles and generates recommendations)
    c2_consumer_config = InferenceConsumerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_consumer_topic=config.kafka_article_consumer_topic,
        kafka_consumer_group=config.kafka_article_consumer_group,
        batch_size=config.kafka_batch_size,
    )
    c2_config = ALSArticleConsumerConfig()
    c2 = ALSArticleConsumer(model, producer, logger, c2_consumer_config, c2_config)
    
    # Bootstrap collections
    c2.bootstrap()
    
    # Start consumers
    c1.run()
    c2.run()

    @app.get("/health")
    async def health_check():
        """Health check endpoint."""
        if not c1.health() or not c2.health():
            raise HTTPException(status_code=503, detail="Service Unhealthy")
        return {"status": "healthy"}
    
    @app.get("/user/{user_id}/profile")
    async def get_user_profile(user_id: int):
        """Get user profile including latent factor."""
        profile = c1.get_user_profile(user_id)
        if not profile:
            raise HTTPException(status_code=404, detail="User profile not found")
        return profile

    uvicorn.run(app, host=config.api_host, port=config.api_port)
