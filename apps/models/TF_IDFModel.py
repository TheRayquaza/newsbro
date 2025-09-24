from .BaseRecommendationModel import BaseRecommendationModel
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import linear_kernel
from dotenv import load_dotenv
import pandas as pd
import mlflow
import mlflow.sklearn

load_dotenv()

mlflow.autolog()
mlflow.set_tracking_uri("http://mlflow.localhost:8080")
mlflow.set_experiment("tfidf_experiment")

class TfidfModel(BaseRecommendationModel):
    
    def __init__(self, name="TF-IDF Model", description="TF-IDF based content similarity"):
        super().__init__(name=name, description=description)
        self.vectorizer = None
        self.tfidf_matrix = None
        self.news_df = None
        self.is_trained = False

    def fit(self, news_df: pd.DataFrame, **kwargs):
        """Train the TF-IDF recommendation model"""
        self.news_df = news_df.copy()
        
        with mlflow.start_run(run_name="TF-IDF Model Training"):
            self.vectorizer = TfidfVectorizer(stop_words='english', **kwargs)
            self.tfidf_matrix = self.vectorizer.fit_transform(self.news_df['content'])

            mlflow.log_params(kwargs)
            mlflow.log_param("num_documents", len(self.news_df))

            mlflow.sklearn.log_model(self.vectorizer, "tfidf_vectorizer")

        self.is_trained = True

    def recommend(self, title: str, top_n: int = 5) -> pd.DataFrame:
        """Get top-N recommended articles for a given title"""
        if not self.is_trained:
            raise ValueError("Model not trained yet. Call fit() first.")

        if title not in self.news_df['title'].values:
            raise ValueError(f"Title '{title}' not found in news dataset.")

        idx = self.news_df[self.news_df['title'] == title].index[0]

        sim_scores = linear_kernel(self.tfidf_matrix[idx:idx+1], self.tfidf_matrix).flatten()

        sim_indices = sim_scores.argsort()[::-1][1:top_n+1]

        recommendations = self.news_df.iloc[sim_indices].copy()
        recommendations['similarity_score'] = sim_scores[sim_indices]

        return recommendations[['title', 'abstract', 'similarity_score']]
