from collections import defaultdict
from typing import List, Dict
from qdrant_client import QdrantClient
from abstract.consumer import InferenceConsumer
from abstract.message import ArticleAggregate, FeedbackAggregate
from abstract.mlflow_model import MlflowModel

def recommend(
    consumer: InferenceConsumer,
    articles: List[ArticleAggregate],
    limit: int = 10
) -> List[tuple[ArticleAggregate, Dict[str, float]]]:
    """
    Generate recommendation scores for given articles based on each user's last 10 liked articles.

    Returns:
        List of tuples (article, {user_id: score})
    """
    model: MlflowModel = consumer.state["model"]
    qdrant: QdrantClient = consumer.state["qdrant"]
    qdrant_collection = consumer.config["QDRANT_COLLECTION"]

    all_points, _ = qdrant.scroll(
        collection_name=qdrant_collection,
        limit=10000,
        with_payload=True,
        with_vectors=False
    )

    user_liked_articles = defaultdict(list)
    for point in all_points:
        payload: FeedbackAggregate = point.payload
        if payload.value == 1:
            user_liked_articles[payload.user_id].append(payload)
    
    user_mean_vectors = {}
    for user_id, liked_articles in user_liked_articles.items():
        last_liked: List[FeedbackAggregate] = liked_articles[-limit:] if len(liked_articles) >= limit else liked_articles

        if last_liked:
            texts = [article['abstract'] for article in last_liked]
            vectors = model.transform(texts).toarray()
            mean_vector = vectors.mean(axis=0)
            user_mean_vectors[user_id] = mean_vector

    article_texts = [a.abstract for a in articles]
    article_vectors = model.transform(article_texts).toarray()

    recommendations = []

    for article, article_vector in zip(articles, article_vectors):
        user_scores = {}
        for user_id, user_mean_vector in user_mean_vectors.items():
            similarity = (article_vector @ user_mean_vector) / (
                (article_vector @ article_vector) ** 0.5 * 
                (user_mean_vector @ user_mean_vector) ** 0.5
            )
            user_scores[user_id] = float(similarity)

        recommendations.append((article, user_scores))

    return recommendations
