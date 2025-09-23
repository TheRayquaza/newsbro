import pandas as pd
from abc import ABC, abstractmethod
from typing import Dict
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
#mlflow.set_tracking_uri("http://mlflow.localhost:8080")

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
    def recommend(self, news_index: int, top_n: int = 6, **kwargs) -> pd.DataFrame:
        """Get recommendations for a given news"""
        pass
    
    def get_info(self) -> Dict[str, str]:
        """Get model information"""
        return {
            "name": self.name,
            "description": self.description,
            "is_trained": self.is_trained,
            "mlflow_run_id": self.mlflow_run_id,
            "model_version": self.model_version
        }
