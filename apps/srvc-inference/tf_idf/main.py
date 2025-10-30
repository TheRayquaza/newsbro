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
            print(os.getenv("TFIDF_MAX_FEATURES"))
        elif len(sys.argv) > 1:
            raise FileNotFoundError(f"Environment file {env_file} not found")

    from tf_idf.src.article_consumer import (
        TFIDFArticleConsumer,
        TFIDFArticleConsumerConfig,
    )
    from tf_idf.src.config import Config  # because Config evaluation is done on runtime
    from tf_idf.src.feedback_consumer import (
        TFIDFFeedbackConsumer,
        TFIDFFeedbackConsumerConfig,
    )

    config = Config()

    logging.basicConfig(
        level=(logging.INFO if config.log_level.upper() == "INFO" else logging.DEBUG),
        format="%(asctime)s [%(levelname)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
    logger = logging.getLogger(__name__)

    app = FastAPI(
        description="TF-IDF based article recommendation API using Qdrant",
        version="0.0.0",
    )

    model = MlflowModel(
        model_uri=config.model_uri,
        tracking_uri=config.tracking_uri,
    )

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
    )
    c1_config = TFIDFFeedbackConsumerConfig()
    logger.info(c1_config)
    c1 = TFIDFFeedbackConsumer(model, producer, logger, c1_consumer_config, c1_config)

    c2_consumer_config = InferenceConsumerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_consumer_topic=config.kafka_article_consumer_topic,
        kafka_consumer_group=config.kafka_article_consumer_group,
        batch_size=config.kafka_batch_size,
    )
    c2_config = TFIDFArticleConsumerConfig()
    logger.info(c2_config)
    c2 = TFIDFArticleConsumer(model, producer, logger, c2_consumer_config, c2_config)
    # c1.bootstrap() no more bootstrap
    c2.bootstrap()
    c1.run()
    c2.run()

    @app.get("/health")
    async def health_check():
        if not c1.health() or not c2.health():
            raise HTTPException(status_code=503, detail="Service Unhealthy")
        return {"status": "healthy"}

    uvicorn.run(app, host=config.api_host, port=config.api_port)
