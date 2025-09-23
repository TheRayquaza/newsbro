from models import BaseRecommendationModel
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np
import pandas as pd

class SimModel(BaseRecommendationModel):

    def fit(self, news_df: pd.DataFrame, **kwargs):
        """Train the recommendation model"""
        self.model = SentenceTransformer('all-MiniLM-L6-v2')
        self.news_df = news_df
        self.news_embeddings = self.model.encode(news_df['content'], show_progress_bar=True)
        self.similarity_matrix = cosine_similarity(self.news_embeddings, self.news_embeddings)
        self.is_trained = True

    def recommend(self, title: str, top_n: int = 6) -> pd.DataFrame:
        """Get recommendations for a given paper by title"""
        if not self.is_trained:
            raise ValueError("Model not trained yet. Call fit() first.")
        
        idx = self.news_df[self.news_df['title'] == title].index[0]
        sim_scores = list(enumerate(self.similarity_matrix[idx]))
        sim_scores = sorted(sim_scores, key=lambda x: x[1], reverse=True)
        top_articles = [i[0] for i in sim_scores[1:top_n+1]]
        return self.news_df.iloc[top_articles][['title', 'abstract']]
