import logging
import os

import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI

from abstract.mlflow_model import MlflowModel
from sbert.src.config import Config

if __name__ == "__main__":
    if os.getenv("ENVIRONMENT") != "production":
        print("Loading .env file for development environment")
        load_dotenv()
    config = Config()

    print(config)

    logging.basicConfig(
        level=(logging.INFO if config.log_level.upper() == "INFO" else logging.DEBUG),
        format="%(asctime)s [%(levelname)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
    logger = logging.getLogger(__name__)

    app = FastAPI(
        title="Article Recommendation API - SBERT",
        description="SBERT based article recommendation API using Qdrant",
        version="0.0.0",
    )

    model = MlflowModel(
        model_uri=config.model_uri,
        tracking_uri=config.tracking_uri,
    )

    # producer_config = InferenceProducerConfig(
    #     kafka_bootstrap_servers=config.kafka_bootstrap_servers,
    #     kafka_producer_topic=config.kafka_producer_topic,
    # )
    # producer = InferenceProducer(logger, producer_config)

    # c1_config = InferenceConsumerConfig(
    #     kafka_bootstrap_servers=config.kafka_bootstrap_servers,
    #     kafka_consumer_topic=config.kafka_feedback_consumer_topic,
    #     kafka_consumer_group=config.kafka_feedback_consumer_group,
    # )
    # c1 = create_feedback_consumer(model, producer, logger, c1_config)

    # c2_config = InferenceConsumerConfig(
    #     kafka_bootstrap_servers=config.kafka_bootstrap_servers,
    #     kafka_consumer_topic=config.kafka_article_consumer_topic,
    #     kafka_consumer_group=config.kafka_article_consumer_group,
    # )
    # c2 = create_article_consumer(model, producer, logger, c2_config)

    # c1.run()
    # c2.run()

    @app.get("/health")
    async def health_check():
        # if not c1.health() or not c2.health():
        #     raise HTTPException(status_code=503, detail="Service Unhealthy")
        return {"status": "healthy"}

    uvicorn.run(app, host=config.api_host, port=config.api_port)
