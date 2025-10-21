import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from pathlib import Path
import logging

class TFIDFModel:
    def __init__(self):
        self.logger = logging.getLogger(__name__)

    def fit(self, path: str, params: dict = {}) -> TfidfVectorizer:
        """Load CSV data and train TF-IDF vectorizer"""
        csv_path = Path(path)
        max_features = params.get("TFIDF_MAX_FEATURES", 12)

        if not csv_path.exists():
            self.logger.warning(f"CSV file not found at {path}. Using minimal corpus.")
            return None
        else:
            self.logger.info(f"Loading training data from {path}")
            try:
                df = pd.read_csv(csv_path)
                
                if 'abstract' not in df.columns:
                    self.logger.error("CSV must contain 'abstract' column")
                    raise ValueError("Missing 'abstract' column in CSV")
        
                initial_corpus = df['abstract'].head(1000).tolist()
                self.logger.info(f"Loaded {len(initial_corpus)} documents for TF-IDF training")
                
            except Exception as e:
                self.logger.error(f"Error loading CSV: {e}. Using minimal corpus.")
                initial_corpus = [
                    "First sample article to initialize the vectorizer",
                    "Second sample article with different abstract"
                ]

        vectorizer = TfidfVectorizer(
            max_features=max_features,
            stop_words='english',
            min_df=2,
            max_df=0.8
        )
        vectorizer.fit(initial_corpus)
        self.logger.info(f"TF-IDF vectorizer trained with {max_features} features.")
        
        return vectorizer
