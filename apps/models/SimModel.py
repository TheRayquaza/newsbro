

from models import BaseRecommendationModel
import numpy as np
import pandas as pd

class SimModel(BaseRecommendationModel):
    def fit(self, papers_df: pd.DataFrame, **kwargs) -> None:
        """Train the recommendation model"""
        pass
    
    def recommend(self, paper_index: int, top_n: int = 6, **kwargs) -> pd.DataFrame:
        """Get recommendations for a given paper"""
        pass
