import logging
import os
from datetime import datetime, timedelta
from typing import List

from qdrant_client import QdrantClient
from qdrant_client.http.models import (
    Distance,
    FieldCondition,
    Filter,
    Range,
    VectorParams,
)

from abstract.consumer import InferenceConsumer, InferenceConsumerConfig
from abstract.message import FeedbackAggregate
from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer


def health_hook(consumer: InferenceConsumer) -> bool:
    try:
        consumer.state["qdrant"].get_collections()
        return True
    except Exception as e:
        consumer.logger.error(f"Health check failed: {e}")
        return False


def bootstrap_hook(consumer: InferenceConsumer, qdrant_collection: str) -> None:
    try:
        collections = [
            c.name for c in consumer.state["qdrant"].get_collections().collections
        ]
        if qdrant_collection not in collections:
            consumer.state["qdrant"].create_collection(
                collection_name=qdrant_collection,
                vectors_config=VectorParams(size=1, distance=Distance.COSINE),
            )
            consumer.state["qdrant"].create_payload_index(
                collection_name=qdrant_collection,
                field_name="timestamp",
                field_schema="float",
            )
            consumer.state["qdrant"].create_payload_index(
                collection_name=qdrant_collection,
                field_name="user_id",
                field_schema="integer",
            )
            consumer.state["qdrant"].create_payload_index(
                collection_name=qdrant_collection,
                field_name="type",
                field_schema="keyword",
            )
            consumer.logger.info(
                f"Collection '{qdrant_collection}' created with indexes."
            )
        else:
            consumer.logger.info(f"Collection '{qdrant_collection}' already exists.")
    except Exception as e:
        consumer.logger.error(f"Error creating collection: {e}")
        raise


def process_hook(consumer: InferenceConsumer, batch: List[FeedbackAggregate]) -> None:
    pass
    # try:
    #     for feedback in batch:
    #         try:
    # user_id = feedback.user_id
    # article_id = feedback.news_id
    # feedback_type = feedback.value
    # feedback_date = feedback.date

    # try:
    #     feedback_datetime = datetime.fromisoformat(feedback_date)
    #     timestamp = feedback_datetime.timestamp()
    # except Exception:
    #     timestamp = datetime.now().timestamp()

    # point = PointStruct(
    #     id=feedback_id,
    #     vector=[0.0],
    #     payload={
    #         "user_id": user_id,
    #         "article_id": article_id,
    #         "type": feedback_type,
    #         "date": feedback_date,
    #         "timestamp": timestamp,
    #         "text": feedback.get("text", ""),
    #         "rating": feedback.get("rating")
    #     }
    # )
    # consumer.state["qdrant"].upsert(collection_name=consumer.config["QDRANT_COLLECTION"], points=[point])
    # consumer.logger.info(f"✅ Feedback {feedback_id} (user: {user_id}, article: {article_id}, type: {feedback_type}) saved.")

    #         except Exception as e:
    #             consumer.logger.error(f"Error processing message: {e}", exc_info=True)
    # except Exception as e:
    #     consumer.logger.error(f"Error processing batch: {e}", exc_info=True)


def create_feedback_consumer(
    model: MlflowModel,
    producer: InferenceProducer,
    logger: logging.Logger,
    config: InferenceConsumerConfig,
) -> InferenceConsumer:
    consumer = InferenceConsumer(
        logger,
        config,
        {
            "FEEDBACK_RETENTION_DAYS": int(os.getenv("FEEDBACK_RETENTION_DAYS", "30")),
            "QDRANT_COLLECTION": "feedbacks",
        },
        health_hook=health_hook,
        process_hook=process_hook,
        bootstrap_hook=lambda c: bootstrap_hook(
            c, os.getenv("QDRANT_COLLECTION", "feedbacks")
        ),
    )
    consumer.state["qdrant"] = QdrantClient(
        url=os.getenv("QDRANT_URL", "http://localhost:6333")
    )
    consumer.state["producer"] = producer
    consumer.state["model"] = model
    bootstrap_hook(consumer, consumer.config["QDRANT_COLLECTION"])
    return consumer


def cleanup_old_feedbacks(consumer: InferenceConsumer) -> None:
    try:
        cutoff_timestamp = (
            datetime.now() - timedelta(days=consumer.config["FEEDBACK_RETENTION_DAYS"])
        ).timestamp()

        # Scroll through old points
        old_points, _ = consumer.state["qdrant"].scroll(
            collection_name=consumer.config["QDRANT_COLLECTION"],
            scroll_filter=Filter(
                must=[FieldCondition(key="timestamp", range=Range(lt=cutoff_timestamp))]
            ),
            limit=1000,
            with_payload=False,
            with_vectors=False,
        )

        if old_points:
            old_ids = [p.id for p in old_points]
            consumer.state["qdrant"].delete(
                collection_name=consumer.config["QDRANT_COLLECTION"],
                points_selector=old_ids,
            )
            consumer.logger.info(
                f"🗑️  {len(old_ids)} old feedbacks deleted (older than {consumer.config['FEEDBACK_RETENTION_DAYS']} days.)"
            )
        else:
            consumer.logger.info("No old feedbacks to delete.")
    except Exception as e:
        consumer.logger.error(f"Error during feedback cleanup: {e}", exc_info=True)
