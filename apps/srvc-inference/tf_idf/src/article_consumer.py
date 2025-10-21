import logging
import os
from typing import List

from qdrant_client import QdrantClient
from qdrant_client.http.models import Distance, VectorParams

from abstract.consumer import InferenceConsumer, InferenceConsumerConfig
from abstract.message import ArticleAggregate
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
                vectors_config=VectorParams(
                    size=consumer.config["TFIDF_MAX_FEATURES"], distance=Distance.COSINE
                ),
            )
            consumer.logger.info(
                f"Collection '{qdrant_collection}' created with vector size {consumer.config['TFIDF_MAX_FEATURES']}."
            )
        else:
            collection_info = consumer.state["qdrant"].get_collection(qdrant_collection)
            vector_size = collection_info.config.params.vectors.size
            if vector_size != consumer.config["TFIDF_MAX_FEATURES"]:
                consumer.logger.error(
                    f"Collection exists but vector size mismatch: "
                    f"expected {consumer.config['TFIDF_MAX_FEATURES']}, got {vector_size}"
                )
                raise ValueError("Vector size mismatch with existing collection")
            consumer.logger.info(f"Collection '{qdrant_collection}' already exists.")

    except Exception as e:
        consumer.logger.error(f"Error creating collection: {e}")
        raise


def process_hook(consumer: InferenceConsumer, batch: List[ArticleAggregate]) -> None:
    # try:
    #     to_upsert = []
    #     to_delete = []

    #     for article in batch:
    #         try:
    #             action = article.get("action", "add")

    #             if action in ["add", "update"]:
    #                 to_upsert.append(article)
    #             elif action == "delete":
    #                 to_delete.append(article)

    #         except Exception as e:
    #             consumer.logger.error(f"Error processing message: {e}", exc_info=True)

    #     if to_upsert:
    #         texts = [a.get("abstract", "") for a in to_upsert]
    #         vectors = consumer.state["model"].transform(texts).toarray().tolist()

    #         for v in vectors:
    #             if len(v) != consumer.config["TFIDF_MAX_FEATURES"]:
    #                 consumer.logger.error(
    #                     f"Vector dimension mismatch: expected {consumer.config['TFIDF_MAX_FEATURES']}, got {len(v)}"
    #                 )
    #                 return

    #         points = [
    #             PointStruct(
    #                 id=int(a["id"]),
    #                 vector=v,
    #                 payload={
    #                     "title": a.get("title", ""),
    #                     "abstract": a.get("abstract", ""),
    #                     "category": a.get("category", ""),
    #                     "created_at": a.get("created_at", "")
    #                 }
    #             )
    #             for a, v in zip(to_upsert, vectors)
    #         ]
    #         consumer.state["qdrant"].upsert(collection_name=consumer.config["QDRANT_COLLECTION"], points=points)
    #         consumer.logger.info(f"✅ Upserted batch of {len(points)} articles")

    #     if to_delete:
    #         ids_to_delete = [int(a["id"]) for a in to_delete]
    #         consumer.state["qdrant"].delete(
    #             collection_name=consumer.config["QDRANT_COLLECTION"],
    #             points_selector=ids_to_delete
    #         )
    #         consumer.logger.info(f"🗑️  Deleted batch of {len(ids_to_delete)} articles")

    # except Exception as e:
    #     consumer.logger.error(f"Error processing batch: {e}", exc_info=True)
    pass


def create_article_consumer(
    model: MlflowModel,
    producer: InferenceProducer,
    logger: logging.Logger,
    config: InferenceConsumerConfig,
) -> InferenceConsumer:
    consumer = InferenceConsumer(
        logger,
        config,
        {
            "TFIDF_MAX_FEATURES": int(os.getenv("TFIDF_MAX_FEATURES", "12")),
            "CSV_DATA_PATH": os.getenv("CSV_DATA_PATH", "../data/news_cleaned.csv"),
            "QDRANT_COLLECTION": "articles",
        },
        health_hook=health_hook,
        process_hook=process_hook,
        bootstrap_hook=lambda c: bootstrap_hook(
            c, os.getenv("QDRANT_COLLECTION", "articles")
        ),
    )

    consumer.state["qdrant"] = QdrantClient(
        url=os.getenv("QDRANT_URL", "http://localhost:6333")
    )
    consumer.state["producer"] = producer
    consumer.state["model"] = model
    bootstrap_hook(consumer, consumer.config["QDRANT_COLLECTION"])

    return consumer
