from abc import ABC, abstractmethod
from typing import Any, Dict, List

import pandas as pd


class BaseRecommendationModel(ABC):
    """Abstract base class for recommendation models"""

    def __init__(self, name: str, description: str):
        self.name = name
        self.description = description
        self.is_trained = False
        self.mlflow_run_id = None
        self.model_version = None

    @abstractmethod
    def fit(self, news_df: pd.DataFrame, **kwargs) -> None:
        """Train the recommendation model"""
        pass

    @abstractmethod
    def recommend(self, news: List[Dict[str, Any]], **kwargs) -> pd.DataFrame:
        """Get recommendations for a given news"""
        pass

    def get_info(self) -> Dict[str, str | bool | None]:
        """Get model information"""
        return {
            "name": self.name,
            "description": self.description,
            "is_trained": self.is_trained,
            "mlflow_run_id": self.mlflow_run_id,
            "model_version": self.model_version,
        }
