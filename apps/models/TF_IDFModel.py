from models import BaseRecommendationModel
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import pandas as pd

class TfidfModel(BaseRecommendationModel):
    
    def fit(self, news_df: pd.DataFrame, **kwargs):
        """Train the TF-IDF recommendation model"""
        self.news_df = news_df.copy()
        # Combine title and abstract as content
        self.news_df['content'] = self.news_df['title'] + " " + self.news_df['abstract']
        
        # Fit TF-IDF vectorizer
        self.vectorizer = TfidfVectorizer(stop_words='english')
        self.tfidf_matrix = self.vectorizer.fit_transform(self.news_df['content'])
        
        # Compute similarity matrix
        self.similarity_matrix = cosine_similarity(self.tfidf_matrix, self.tfidf_matrix)
        
        self.is_trained = True

    def recommend(self, title: str, top_n: int = 5) -> pd.DataFrame:
        """Get recommendations for a given article by title"""
        if not self.is_trained:
            raise ValueError("Model not trained yet. Call fit() first.")
        
        if title not in self.news_df['title'].values:
            raise ValueError(f"Title '{title}' not found in news dataset.")
        
        idx = self.news_df[self.news_df['title'] == title].index[0]
        sim_scores = list(enumerate(self.similarity_matrix[idx]))
        sim_scores = sorted(sim_scores, key=lambda x: x[1], reverse=True)
        top_articles = [i[0] for i in sim_scores[1:top_n+1]]
        
        return self.news_df.iloc[top_articles][['title', 'abstract']]
