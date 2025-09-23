from .BaseRecommendationModel import BaseRecommendationModel
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np
import pandas as pd

class SBertModel(BaseRecommendationModel):
    
    def __init__(self, name="S-BERT Model", description="Semantic similarity using sentence transformers"):
        super().__init__(name=name, description=description)
        self.model = None
        self.news_df = None
        self.news_embeddings = None
        self.similarity_matrix = None

    def fit(self, news_df: pd.DataFrame, **kwargs):
        """Train the recommendation model"""
        print("Training S-BERT model...")
        self.model = SentenceTransformer('all-MiniLM-L6-v2')
        self.news_df = news_df.copy()
        
        # Ensure we have the content column
        if 'content' not in self.news_df.columns:
            self.news_df['content'] = self.news_df['title'].astype(str) + " " + self.news_df['abstract'].astype(str)
        
        self.news_embeddings = self.model.encode(self.news_df['content'].tolist(), show_progress_bar=True)
        self.similarity_matrix = cosine_similarity(self.news_embeddings, self.news_embeddings)
        self.is_trained = True
        print("S-BERT model training completed.")

    def recommend(self, title: str, top_n: int = 6) -> pd.DataFrame:
        """Get recommendations for a given article by title"""
        if not self.is_trained:
            raise ValueError("Model not trained yet. Call fit() first.")
        
        # Find the article by title
        matching_articles = self.news_df[self.news_df['title'] == title]
        if matching_articles.empty:
            raise ValueError(f"Title '{title}' not found in news dataset.")
        
        idx = matching_articles.index[0]
        
        # Get similarity scores
        sim_scores = list(enumerate(self.similarity_matrix[idx]))
        sim_scores = sorted(sim_scores, key=lambda x: x[1], reverse=True)
        
        # Get top articles (excluding the original)
        top_articles_indices = [i[0] for i in sim_scores[1:top_n+1]]
        recommendations = self.news_df.iloc[top_articles_indices].copy()
        
        # Add similarity scores
        recommendations['similarity_score'] = [sim_scores[i+1][1] for i in range(len(top_articles_indices))]
        
        return recommendations[['title', 'abstract', 'similarity_score']]
