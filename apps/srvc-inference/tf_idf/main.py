import logging
import os
import sys

import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException

from abstract.consumer import InferenceConsumerConfig
from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer, InferenceProducerConfig
from tf_idf.src.article_consumer import (
    TFIDFArticleConsumer,
    TFIDFArticleConsumerConfig,
)
from tf_idf.src.feedback_consumer import (
    TFIDFFeedbackConsumer,
    TFIDFFeedbackConsumerConfig,
)

if __name__ == "__main__":
    if os.getenv("ENVIRONMENT") != "production":
        env_file = sys.argv[1] if len(sys.argv) > 1 else ".env"
        if os.path.exists(env_file):
            load_dotenv(env_file)
        elif len(sys.argv) > 1:
            raise FileNotFoundError(f"Environment file {env_file} not found")

    from tf_idf.src.config import Config  # because Config evaluation is done on runtime

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

    c1_consumer_config = InferenceConsumerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_consumer_topic=config.kafka_feedback_consumer_topic,
        kafka_consumer_group=config.kafka_feedback_consumer_group,
        batch_size=config.kafka_batch_size,
        batch_interval=config.kafka_batch_interval,
        thread_pool_size=config.thread_pool_size,
    )
    c1_config = TFIDFFeedbackConsumerConfig(
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
    c1 = TFIDFFeedbackConsumer(model, producer, logger, c1_consumer_config, c1_config)

    c2_consumer_config = InferenceConsumerConfig(
        kafka_bootstrap_servers=config.kafka_bootstrap_servers,
        kafka_consumer_topic=config.kafka_article_consumer_topic,
        kafka_consumer_group=config.kafka_article_consumer_group,
        batch_size=config.kafka_batch_size,
        batch_interval=config.kafka_batch_interval,
        thread_pool_size=config.thread_pool_size,
    )
    c2_config = TFIDFArticleConsumerConfig(
        articles_collection=config.articles_collection,
        qdrant_url=config.qdrant_url,
        qdrant_api_key=config.qdrant_api_key,
        redis_sentinels=config.redis_sentinels,
        redis_master_name=config.redis_master_name,
        redis_password=config.redis_password,
        redis_db=config.redis_db,
        redis_scan_batch_size=config.redis_scan_batch_size,
        redis_user_profile_prefix=config.redis_user_profile_prefix,
        article_vector_features=config.article_vector_features,
        similarity_threshold=config.similarity_threshold,
        model_name=config.model_name,
    )
    c2 = TFIDFArticleConsumer(model, producer, logger, c2_consumer_config, c2_config)
    c2.bootstrap()
    c1.run()
    c2.run()

    @app.get("/health")
    async def health_check():
        if not c1.health() or not c2.health():
            raise HTTPException(status_code=503, detail="Service Unhealthy")
        return {"status": "healthy"}

    uvicorn.run(app, host=config.api_host, port=config.api_port)
