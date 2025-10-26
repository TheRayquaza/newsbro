import logging
import os
import time
from datetime import datetime
from typing import Dict, List, Set

import numpy as np
import pydantic
from qdrant_client import QdrantClient

from abstract.mlflow_model import MlflowModel
from abstract.producer import InferenceProducer, InferenceProducerConfig


class RecommendationMessage(pydantic.BaseModel):
    """Message format for recommendations sent to Kafka."""
    model: str
    user_id: int
    article_id: int
    score: float
    created_at: datetime


def poll_recommendations(
    producer: InferenceProducer,
    qdrant: QdrantClient,
    model: MlflowModel,
    config: dict,
    processed_article_ids: Set[int]
) -> int:
    """
    Poll Qdrant for new articles and generate recommendations.
    
    Returns:
        Number of recommendation messages produced.
    """
    produced_count = 0

    try:
        # Load all articles from Qdrant
        articles_map = _load_all_articles(qdrant, config["ARTICLES_COLLECTION"], producer.logger)
        if not articles_map:
            producer.logger.warning("No articles found in Qdrant")
            return 0

        # Identify new articles
        current_article_ids = set(articles_map.keys())
        new_article_ids = current_article_ids - processed_article_ids

        if not new_article_ids:
            producer.logger.debug("No new articles to process")
            return 0

        producer.logger.info(f"Found {len(new_article_ids)} new articles to process")

        # Load feedbacks
        feedbacks = _load_all_feedbacks(qdrant, config["FEEDBACKS_COLLECTION"], producer.logger)

        # Build user histories
        user_liked_articles = _build_user_liked_articles(feedbacks, articles_map, config["USER_HISTORY_LIMIT"], producer.logger)
        
        if not user_liked_articles:
            producer.logger.warning("No user likes found, skipping recommendations")
            processed_article_ids.update(new_article_ids)
            return 0

        # Compute mean TF-IDF vectors for each user
        user_mean_vectors = _compute_user_mean_vectors(user_liked_articles, model, producer.logger)

        if not user_mean_vectors:
            producer.logger.warning("Could not compute user vectors")
            processed_article_ids.update(new_article_ids)
            return 0

        # Process each new article
        new_articles = [articles_map[aid] for aid in new_article_ids]
        
        # Extract texts and compute TF-IDF vectors for new articles
        article_texts = [
            article.get("abstract") or article.get("content") or article.get("title") or ""
            for article in new_articles
        ]
        article_vectors = model.transform(article_texts).toarray()

        # Generate recommendations
        recommendations: List[RecommendationMessage] = []

        for article_id, article_vector in zip(new_article_ids, article_vectors):
            for user_id, user_mean_vector in user_mean_vectors.items():
                # Compute similarity
                similarity = _compute_cosine_similarity(article_vector, user_mean_vector)

                # Check threshold
                if similarity >= config["SIMILARITY_THRESHOLD"]:
                    recommendation = RecommendationMessage(
                        model=config["MODEL_NAME"],
                        user_id=int(user_id),
                        article_id=int(article_id),
                        score=float(similarity),
                        created_at=datetime.now(datetime.timezone.utc)
                    )
                    recommendations.append(recommendation)

        # Produce recommendations
        if recommendations:
            producer.produce(recommendations)
            produced_count = len(recommendations)
            producer.logger.info(f"Produced {produced_count} recommendations for {len(new_article_ids)} new articles")

        # Mark articles as processed
        processed_article_ids.update(new_article_ids)

    except Exception as e:
        producer.logger.error(f"Error in poll_recommendations: {e}", exc_info=True)

    return produced_count


def _load_all_articles(qdrant: QdrantClient, collection_name: str, logger: logging.Logger) -> Dict[int, dict]:
    """Load all articles from Qdrant and return as a dict mapping id -> payload."""
    try:
        points, _ = qdrant.scroll(
            collection_name=collection_name,
            limit=10000,
            with_payload=True,
            with_vectors=False,
        )
        articles = {}
        for point in points:
            try:
                articles[int(point.id)] = point.payload
            except Exception as e:
                logger.warning(f"Skipping malformed article point: {e}")
                continue
        return articles
    except Exception as e:
        logger.error(f"Error loading articles: {e}", exc_info=True)
        return {}


def _load_all_feedbacks(qdrant: QdrantClient, collection_name: str, logger: logging.Logger) -> List[dict]:
    """Load all feedback records from Qdrant."""
    try:
        points, _ = qdrant.scroll(
            collection_name=collection_name,
            limit=10000,
            with_payload=True,
            with_vectors=False,
        )
        feedbacks = []
        for point in points:
            try:
                feedbacks.append(point.payload)
            except Exception as e:
                logger.warning(f"Skipping malformed feedback point: {e}")
                continue
        return feedbacks
    except Exception as e:
        logger.error(f"Error loading feedbacks: {e}", exc_info=True)
        return []


def _build_user_liked_articles(
    feedbacks: List[dict],
    articles_map: Dict[int, dict],
    history_limit: int,
    logger: logging.Logger
) -> Dict[int, List[dict]]:
    """
    Build a mapping of user_id -> list of liked article payloads (last N sorted by timestamp).
    """
    user_likes: Dict[int, List[dict]] = {}

    for feedback in feedbacks:
        # Check if this is a like (value == 1)
        value = feedback.get("value")
        if value != 1:
            continue

        try:
            user_id = int(feedback.get("user_id"))
            # Try different possible field names for article ID
            article_id = int(
                feedback.get("news_id")
                or feedback.get("article_id")
                or feedback.get("article")
                or 0
            )
            
            if article_id <= 0 or article_id not in articles_map:
                continue

            # Get timestamp for sorting
            timestamp = feedback.get("timestamp") or feedback.get("date") or feedback.get("created_at")
            try:
                ts = float(timestamp) if timestamp else datetime.now(datetime.timezone.utc).timestamp()
            except Exception:
                ts = datetime.now(datetime.timezone.utc).timestamp()

            user_likes.setdefault(user_id, []).append({
                "article_id": article_id,
                "article": articles_map[article_id],
                "timestamp": ts
            })
        except Exception as e:
            logger.warning(f"Error processing feedback: {e}")
            continue

    # Sort by timestamp and keep only last N articles per user
    user_history: Dict[int, List[dict]] = {}
    for user_id, likes in user_likes.items():
        sorted_likes = sorted(likes, key=lambda x: x["timestamp"])
        last_n = sorted_likes[-history_limit:]
        user_history[user_id] = [item["article"] for item in last_n]

    return user_history


def _compute_user_mean_vectors(
    user_liked_articles: Dict[int, List[dict]],
    model: MlflowModel,
    logger: logging.Logger
) -> Dict[int, np.ndarray]:
    """
    Compute mean TF-IDF vector for each user based on their liked articles.
    """
    user_mean_vectors = {}

    for user_id, articles in user_liked_articles.items():
        texts = []
        for article in articles:
            # Get text content (try abstract first, then content, then title)
            text = (
                article.get("abstract")
                or article.get("content")
                or article.get("title")
                or ""
            )
            texts.append(text)

        if not texts:
            continue

        try:
            # Transform texts to TF-IDF vectors
            vectors = model.transform(texts).toarray()
            # Compute mean vector
            mean_vector = vectors.mean(axis=0)
            user_mean_vectors[user_id] = mean_vector
        except Exception as e:
            logger.error(f"Error computing vectors for user {user_id}: {e}")
            continue

    return user_mean_vectors


def _compute_cosine_similarity(vec1: np.ndarray, vec2: np.ndarray) -> float:
    """Compute cosine similarity between two vectors."""
    norm1 = np.linalg.norm(vec1)
    norm2 = np.linalg.norm(vec2)
    
    if norm1 == 0 or norm2 == 0:
        return 0.0
    
    return float(np.dot(vec1, vec2) / (norm1 * norm2))


def run_producer(
    producer: InferenceProducer,
    poll_interval_seconds: int = 10
) -> None:
    """
    Run the producer in a continuous polling loop.
    
    Args:
        producer: InferenceProducer instance with state initialized
        poll_interval_seconds: Time to wait between polling cycles
    """
    config = producer.state["config"]
    qdrant = producer.state["qdrant"]
    model = producer.state["model"]
    processed_article_ids = producer.state["processed_article_ids"]
    
    producer.logger.info(
        f"Starting TF-IDF recommendation producer with poll interval: {poll_interval_seconds}s, "
        f"threshold: {config['SIMILARITY_THRESHOLD']}, history limit: {config['USER_HISTORY_LIMIT']}"
    )
    
    try:
        while True:
            count = poll_recommendations(producer, qdrant, model, config, processed_article_ids)
            if count > 0:
                producer.logger.info(f"Poll cycle complete: {count} recommendations produced")
            time.sleep(poll_interval_seconds)
    except KeyboardInterrupt:
        producer.logger.info("Stopping TF-IDF recommendation producer")
    except Exception as e:
        producer.logger.error(f"Fatal error in producer run loop: {e}", exc_info=True)
        raise


def create_recommendation_producer(
    model: MlflowModel,
    logger: logging.Logger,
    producer_config: InferenceProducerConfig,
) -> InferenceProducer:
    """
    Create and initialize a TF-IDF recommendation producer.
    
    Args:
        model: MLflow TF-IDF model for transforming text
        logger: Logger instance
        producer_config: Configuration for the Kafka producer
        
    Returns:
        Configured InferenceProducer instance
    """
    config = {
        "MODEL_NAME": os.getenv("MODEL_NAME", "tfidf"),
        "ARTICLES_COLLECTION": os.getenv("QDRANT_COLLECTION", "articles"),
        "FEEDBACKS_COLLECTION": os.getenv("FEEDBACKS_COLLECTION", "feedbacks"),
        "SIMILARITY_THRESHOLD": float(os.getenv("SIMILARITY_THRESHOLD", "0.7")),
        "USER_HISTORY_LIMIT": int(os.getenv("USER_HISTORY_LIMIT", "10")),
        "QDRANT_URL": os.getenv("QDRANT_URL", "http://localhost:6333"),
    }
    
    producer = InferenceProducer(logger, producer_config)
    
    producer.state["config"] = config
    producer.state["qdrant"] = QdrantClient(url=config["QDRANT_URL"])
    producer.state["model"] = model
    producer.state["processed_article_ids"] = set()
    
    return producer
