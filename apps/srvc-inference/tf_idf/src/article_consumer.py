import logging
import os
from typing import List
import datetime

from qdrant_client import QdrantClient
from qdrant_client.http.models import Distance, VectorParams, PointStruct

from abstract.consumer import InferenceConsumer, InferenceConsumerConfig
from abstract.message import ArticleAggregate, FeedbackAggregate, InferenceCommand
from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer

from tf_idf.src.recommend import recommend

def health(consumer: InferenceConsumer) -> bool:
    try:
        consumer.state["qdrant"].get_collections()
        return True
    except Exception as e:
        consumer.logger.error(f"Health check failed: {e}")
        return False

def bootstrap(consumer: InferenceConsumer, qdrant_collection: str) -> None:
    TFIDF_MAX_FEATURES = consumer.config["TFIDF_MAX_FEATURES"]
    try:
        collections = [
            c.name for c in consumer.state["qdrant"].get_collections().collections
        ]
        if qdrant_collection not in collections:
            consumer.state["qdrant"].create_collection(
                collection_name=qdrant_collection,
                vectors_config=VectorParams(
                    size=TFIDF_MAX_FEATURES, distance=Distance.COSINE
                ),
            )
            consumer.logger.info(
                f"Collection '{qdrant_collection}' created with vector size {consumer.config['TFIDF_MAX_FEATURES']}."
            )
        else:
            collection_info = consumer.state["qdrant"].get_collection(qdrant_collection)
            vector_size = collection_info.config.params.vectors.size
            if vector_size != TFIDF_MAX_FEATURES:
                consumer.logger.error(
                    f"Collection exists but vector size mismatch: "
                    f"expected {consumer.config['TFIDF_MAX_FEATURES']}, got {vector_size}"
                )
                raise ValueError("Vector size mismatch with existing collection")
            consumer.logger.info(f"Collection '{qdrant_collection}' already exists.")

    except Exception as e:
        consumer.logger.error(f"Error creating collection: {e}")
        raise

def process(consumer: InferenceConsumer, batch: List[ArticleAggregate]) -> None:
    TFIDF_MAX_FEATURES = consumer.config["TFIDF_MAX_FEATURES"]
    QDRANT_COLLECTION = consumer.config["QDRANT_COLLECTION"]
    qdrant: QdrantClient = consumer.state["qdrant"]
    model: MlflowModel = consumer.state["model"]

    try:
        to_upsert: List[ArticleAggregate] = []
        to_delete: List[ArticleAggregate] = []

        for article in batch:
            try:
                if article["is_active"]:
                    to_upsert.append(article)
                else:
                    to_delete.append(article)

            except Exception as e:
                consumer.logger.error(f"Error processing message: {e}", exc_info=True)

        if to_upsert:
            texts = [a["abstract"] for a in to_upsert]
            vectors = model.transform(texts).toarray().tolist()

            for v in vectors:
                if len(v) != TFIDF_MAX_FEATURES:
                    consumer.logger.error(
                        f"Vector dimension mismatch: expected {consumer.config['TFIDF_MAX_FEATURES']}, got {len(v)}"
                    )
                    return

            points = [
                PointStruct(
                    id=int(a["id"]),
                    vector=v,
                    payload=a,
                )
                for a, v in zip(to_upsert, vectors)
            ]
            qdrant.upsert(collection_name=QDRANT_COLLECTION, points=points)
            consumer.logger.info(f"✅ Upserted batch of {len(points)} articles")
            
            recommendations = recommend(consumer, to_upsert)
            recommendation_commands: List[InferenceCommand] = []
            for article, user_scores in recommendations:
                for user in user_scores.keys():
                    recommendation_commands.append(InferenceCommand(
                        user,
                        consumer.config["MODEL_NAME"],
                        article=article,
                        date=datetime.utcnow()
                    ))
            producer: InferenceProducer = consumer.state["producer"]
            consumer.logger.info(f"Generating {len(recommendation_commands)} recommendations")
            producer.produce(recommendation_commands)

        if to_delete:
            ids_to_delete = [int(a["id"]) for a in to_delete]
            qdrant.delete(
                collection_name=QDRANT_COLLECTION,
                points_selector=ids_to_delete
            )
            consumer.logger.info(f"🗑️  Deleted batch of {len(ids_to_delete)} articles")

    except Exception as e:
        consumer.logger.error(f"Error processing batch: {e}", exc_info=True)

def create_article_consumer(
    model: MlflowModel,
    producer: InferenceProducer,
    logger: logging.Logger,
    consumer_config: InferenceConsumerConfig,
) -> InferenceConsumer:
    config = {
        "TFIDF_MAX_FEATURES": int(os.getenv("TFIDF_MAX_FEATURES", "12")),
        "QDRANT_COLLECTION": os.getenv("QDRANT_COLLECTION", "articles"),
        "QDRANT_URL": os.getenv("QDRANT_URL", "http://localhost:6333"),
    }
    consumer = InferenceConsumer(
        logger,
        consumer_config,
        config,
        health_hook=health,
        process_hook=lambda batch: process(consumer, batch)
    )

    consumer.state["qdrant"] = QdrantClient(
        url=config["QDRANT_URL"]
    )
    consumer.state["producer"] = producer
    consumer.state["model"] = model
    bootstrap(consumer, config["QDRANT_COLLECTION"])

    return consumer
