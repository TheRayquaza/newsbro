import os
from fastapi import FastAPI, HTTPException, Query
from qdrant_client import QdrantClient
from qdrant_client.http.models import Filter, FieldCondition, MatchValue
from pydantic import BaseModel
from typing import List, Optional
import logging
from dotenv import load_dotenv

load_dotenv()

QDRANT_URL = os.getenv("QDRANT_URL", "http://localhost:6333")
API_PORT = int(os.getenv("API_PORT", "8000"))
TOP_N_RECOMMENDATIONS = int(os.getenv("TOP_N_RECOMMENDATIONS", "10"))

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Article Recommendation API",
    description="TF-IDF based article recommendation API using Qdrant",
    version="1.0.0"
)

qdrant = QdrantClient(url=QDRANT_URL)
ARTICLES_COLLECTION = "articles"
FEEDBACKS_COLLECTION = "feedbacks"

# --- Response Models ---
class RecommendationItem(BaseModel):
    id: int
    title: str
    content: str
    abstract: Optional[str] = ""
    score: float

class RecommendationResponse(BaseModel):
    user_id: int
    recommendations: List[RecommendationItem]
    model: str
    based_on_articles: List[int]

# --- Health Check ---
@app.get("/health")
async def health_check():
    """Health check endpoint"""
    try:
        collections = qdrant.get_collections()
        return {
            "status": "healthy",
            "qdrant_url": QDRANT_URL,
            "collections": [c.name for c in collections.collections]
        }
    except Exception as e:
        logger.error(f"Health check failed: {e}")
        return {"status": "unhealthy", "error": str(e)}

# --- Main Recommendation Endpoint ---
@app.get("/api/v1/recommendation", response_model=RecommendationResponse)
async def get_recommendations(
    user_id: int = Query(..., description="User ID"),
    model: str = Query("tfidf", description="Recommender System (currently only 'tfidf' supported)"),
    top_n: int = Query(10, ge=1, le=100, description="Number of recommendations to return"),
    feedback_type: str = Query("like", description="Feedback type to filter (like, view, dislike, etc.)"),
    max_feedback_count: int = Query(5, ge=1, le=20, description="Maximum number of feedbacks to use")
):
    """
    Generates article recommendations based on user feedback history.
    """
    try:
        logger.info(f"Generating recommendations for user={user_id} with model={model}")
        
        # Step 1: Get the user's recent feedbacks
        feedbacks, _ = qdrant.scroll(
            collection_name=FEEDBACKS_COLLECTION,
            scroll_filter=Filter(
                must=[
                    FieldCondition(
                        key="user_id",
                        match=MatchValue(value=user_id)
                    ),
                    FieldCondition(
                        key="type",
                        match=MatchValue(value=feedback_type)
                    )
                ]
            ),
            limit=max_feedback_count,
            with_payload=True,
            with_vectors=False,
            order_by="timestamp"  # Get most recent first
        )
        
        if not feedbacks:
            raise HTTPException(
                status_code=404,
                detail=f"No feedback of type '{feedback_type}' found for user {user_id}"
            )
        
        # Extract article IDs from feedbacks
        feedback_article_ids = [
            int(fb.payload.get("article_id"))
            for fb in feedbacks
            if fb.payload.get("article_id") is not None
        ]
        logger.info(f"User {user_id} has {len(feedback_article_ids)} '{feedback_type}' feedbacks: {feedback_article_ids}")
        
        if not feedback_article_ids:
            raise HTTPException(
                status_code=404,
                detail=f"No valid article IDs found in user {user_id}'s feedback"
            )
        
        # Step 2: Get the TF-IDF vectors for the feedback articles
        try:
            feedback_articles = qdrant.retrieve(
                collection_name=ARTICLES_COLLECTION,
                ids=feedback_article_ids,
                with_vectors=True,
                with_payload=True
            )
        except Exception as e:
            logger.error(f"Error retrieving articles: {e}")
            raise HTTPException(
                status_code=500,
                detail=f"Error retrieving feedback articles: {str(e)}"
            )
        
        if not feedback_articles:
            raise HTTPException(
                status_code=404,
                detail="Could not retrieve feedback articles from vector database"
            )
        
        # Filter out articles without vectors
        valid_articles = [art for art in feedback_articles if art.vector is not None]
        if not valid_articles:
            raise HTTPException(
                status_code=404,
                detail="No valid vectors found for feedback articles"
            )
        
        # Step 3: Compute the average vector
        vectors = [article.vector for article in valid_articles]
        avg_vector = [sum(values) / len(vectors) for values in zip(*vectors)]
        
        logger.info(f"Computed average vector from {len(vectors)} articles")
        
        # Step 4: Search for similar articles
        search_results = qdrant.search(
            collection_name=ARTICLES_COLLECTION,
            query_vector=avg_vector,
            limit=top_n + len(feedback_article_ids),  # Request extra to account for filtering
            with_payload=True
        )
        
        # Step 5: Filter out already seen articles
        recommendations = []
        for result in search_results:
            if int(result.id) not in feedback_article_ids:
                recommendations.append(RecommendationItem(
                    id=int(result.id),
                    title=result.payload.get("title", ""),
                    content=result.payload.get("content", ""),
                    abstract=result.payload.get("abstract", ""),
                    score=float(result.score)
                ))
            if len(recommendations) >= top_n:
                break
        
        if not recommendations:
            raise HTTPException(
                status_code=404,
                detail="Could not generate recommendations. All similar articles have already been seen."
            )
        
        logger.info(f"✅ Generated {len(recommendations)} recommendations for user_id={user_id}")
        
        return RecommendationResponse(
            user_id=user_id,
            recommendations=recommendations,
            model=model,
            based_on_articles=feedback_article_ids
        )
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Unexpected error during recommendation generation: {e}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=API_PORT)
