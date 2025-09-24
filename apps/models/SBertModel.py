from .BaseRecommendationModel import BaseRecommendationModel
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
from dotenv import load_dotenv
import numpy as np
import mlflow
import pandas as pd

load_dotenv()

mlflow.autolog()
mlflow.set_tracking_uri("http://mlflow.localhost:8080")
mlflow.set_experiment("sbert_experiment")

class SBertModel(BaseRecommendationModel):
    
    def __init__(self, name="S-BERT Model", description="Semantic similarity using sentence transformers"):
        super().__init__(name=name, description=description)
        self.model = None
        self.news_df = None
        self.news_embeddings = None
        self.is_trained = False

    def fit(self, news_df: pd.DataFrame, **kwargs):
        """Train the recommendation model with MLflow tracking"""
        print("Training S-BERT model...")

        with mlflow.start_run(run_name="SBERT Training"):
            self.model = SentenceTransformer('all-MiniLM-L6-v2')
            mlflow.log_param("embedding_model", "all-MiniLM-L6-v2")

            self.news_df = news_df.copy()
            mlflow.log_param("num_documents", len(self.news_df))
            
            self.news_embeddings = self.model.encode(
                self.news_df['content'].tolist(), 
                show_progress_bar=True
            )

            sim_matrix = cosine_similarity(self.news_embeddings)
            avg_sim = (sim_matrix.sum() - len(sim_matrix)) / (len(sim_matrix)**2 - len(sim_matrix))
            mlflow.log_metric("average_similarity", avg_sim)

            mlflow.sentence_transformers.log_model(self.model, "sbert_model")

        self.is_trained = True
        print("S-BERT model training completed.")

    def recommend(self, title: str, top_n: int = 6) -> pd.DataFrame:
        """Get recommendations for a given article by title"""
        if not self.is_trained:
            raise ValueError("Model not trained yet. Call fit() first.")
        
        matching_articles = self.news_df[self.news_df['title'] == title]
        if matching_articles.empty:
            raise ValueError(f"Title '{title}' not found in news dataset.")
        
        idx = matching_articles.index[0]

        query_embedding = self.news_embeddings[idx].reshape(1, -1)
        sim_scores = cosine_similarity(query_embedding, self.news_embeddings)[0]

        sim_indices = np.argsort(sim_scores)[::-1]
        sim_indices = sim_indices[sim_indices != idx][:top_n]
        top_scores = sim_scores[sim_indices]

        recommendations = self.news_df.iloc[sim_indices].copy()
        recommendations['similarity_score'] = top_scores

        return recommendations[['title', 'abstract', 'similarity_score']]
