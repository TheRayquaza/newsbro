import re

from sklearn.base import BaseEstimator, TransformerMixin


class TextCleaner(BaseEstimator, TransformerMixin):
    """
    Scikit-learn transformer for text cleaning and English-language filtering.
    Removes punctuation, numbers, extra spaces, and optionally filters non-English text.
    """

    def __init__(self):
        pass

    def _clean_text(self, text: str) -> str:
        text = re.sub(r"[^A-Za-z\s]", " ", str(text))
        return re.sub(r"\s+", " ", text).strip().lower()

    def fit(self, X, y=None):
        return self

    def transform(self, X):
        return [self._clean_text(doc) for doc in X]
