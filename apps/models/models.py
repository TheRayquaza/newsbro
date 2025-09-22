import numpy as np
import pandas as pd
from abc import ABC, abstractmethod
from typing import List, Dict, Any, Optional, Tuple
import logging
import mlflow
import mlflow.sklearn
import mlflow.pyfunc
from mlflow.models.signature import infer_signature
from mlflow.types.schema import Schema, ColSpec
import tempfile
import pickle
import os
from datetime import datetime

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
    def fit(self, papers_df: pd.DataFrame, **kwargs) -> None:
        """Train the recommendation model"""
        pass
    
    @abstractmethod
    def recommend(self, paper_index: int, top_n: int = 6, **kwargs) -> pd.DataFrame:
        """Get recommendations for a given paper"""
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

class MLflowRecommendationWrapper(mlflow.pyfunc.PythonModel):
    """MLflow wrapper for recommendation models"""
    
    def __init__(self, model, papers_df):
        self.model = model
        self.papers_df = papers_df
    
    def load_context(self, context):
        """Load model context"""
        pass
    
    def predict(self, context, model_input):
        """Make predictions using the wrapped model"""
        # model_input should be a DataFrame with columns: paper_index, top_n
        recommendations = []
        
        for _, row in model_input.iterrows():
            paper_index = int(row['paper_index'])
            top_n = int(row.get('top_n', 6))
            
            rec = self.model.recommend(paper_index, top_n)
            rec['query_paper_index'] = paper_index
            recommendations.append(rec)
        
        return pd.concat(recommendations, ignore_index=True) if recommendations else pd.DataFrame()

class SimilarityMatrixModel(BaseRecommendationModel):
    """Cosine similarity-based recommendation model using precomputed similarity matrix"""
    
    def __init__(self):
        super().__init__(
            name="Similarity Matrix",
            description="Uses precomputed cosine similarity matrix for fast recommendations"
        )
        self.similarity_matrix = None
        self.papers_df = None
    
    def fit(self, papers_df: pd.DataFrame, similarity_matrix: np.ndarray = None, 
            log_mlflow: bool = True, experiment_name: str = "recommendation_experiments", **kwargs) -> None:
        """Load the precomputed similarity matrix and papers data"""
        try:
            # Set up MLflow experiment
            if log_mlflow:
                mlflow.set_experiment(experiment_name)
                
            with mlflow.start_run(run_name=f"{self.name}_{datetime.now().strftime('%Y%m%d_%H%M%S')}") as run:
                self.mlflow_run_id = run.info.run_id
                
                if similarity_matrix is not None:
                    self.similarity_matrix = similarity_matrix
                else:
                    # Load from file if not provided
                    self.similarity_matrix = np.load('../data/similarity_matrix.npy', allow_pickle=True)
                
                self.papers_df = papers_df.copy()
                
                # Validate dimensions
                if len(self.papers_df) != self.similarity_matrix.shape[0]:
                    raise ValueError(f"Papers dataframe length ({len(self.papers_df)}) doesn't match similarity matrix shape ({self.similarity_matrix.shape[0]})")
                
                self.is_trained = True
                
                if log_mlflow:
                    # Log parameters
                    mlflow.log_param("model_name", self.name)
                    mlflow.log_param("num_papers", len(self.papers_df))
                    mlflow.log_param("similarity_matrix_shape", self.similarity_matrix.shape)
                    
                    # Log metrics
                    mlflow.log_metric("matrix_density", np.count_nonzero(self.similarity_matrix) / self.similarity_matrix.size)
                    mlflow.log_metric("mean_similarity", np.mean(self.similarity_matrix))
                    mlflow.log_metric("std_similarity", np.std(self.similarity_matrix))
                    
                    # Save similarity matrix as artifact
                    with tempfile.NamedTemporaryFile(suffix='.npy', delete=False) as tmp_file:
                        np.save(tmp_file.name, self.similarity_matrix)
                        mlflow.log_artifact(tmp_file.name, "similarity_matrix")
                        os.unlink(tmp_file.name)
                    
                    # Log model
                    wrapper = MLflowRecommendationWrapper(self, self.papers_df)
                    
                    # Create sample input for signature inference
                    sample_input = pd.DataFrame({
                        'paper_index': [0, 1],
                        'top_n': [5, 10]
                    })
                    
                    # Infer signature
                    sample_output = wrapper.predict(None, sample_input)
                    signature = infer_signature(sample_input, sample_output)
                    
                    mlflow.pyfunc.log_model(
                        artifact_path="model",
                        python_model=wrapper,
                        signature=signature,
                        registered_model_name=f"RecommendationModel_{self.name.replace(' ', '_')}"
                    )
                
                logger.info(f"SimilarityMatrixModel loaded successfully with {len(self.papers_df)} papers")
                
        except Exception as e:
            logger.error(f"Error loading similarity matrix model: {e}")
            raise
    
    def recommend(self, paper_index: int, top_n: int = 6, **kwargs) -> pd.DataFrame:
        """Get similar papers using cosine similarity"""
        if not self.is_trained:
            raise ValueError("Model not trained. Call fit() first.")
        
        if paper_index >= len(self.papers_df) or paper_index < 0:
            raise ValueError(f"Invalid paper index: {paper_index}")
        
        try:
            # Get similarity scores for the given paper
            similarity_scores = self.similarity_matrix[paper_index]
            
            # Get indices of most similar papers (excluding the paper itself)
            similar_indices = similarity_scores.argsort()[::-1][1:top_n+1]
            
            # Return the similar papers
            recommendations = self.papers_df.iloc[similar_indices].copy()
            
            # Add similarity scores to the results
            recommendations['similarity_score'] = similarity_scores[similar_indices]
            
            logger.info(f"Generated {len(recommendations)} recommendations for paper index {paper_index}")
            return recommendations
            
        except Exception as e:
            logger.error(f"Error generating recommendations: {e}")
            raise    

class RecommendationEngine:
    """Main recommendation engine that manages different models"""
    
    def __init__(self, experiment_name: str = "recommendation_experiments"):
        self.models = {}
        self.active_model = None
        self.experiment_name = experiment_name
        mlflow.set_experiment(experiment_name)
    
    def register_model(self, model: BaseRecommendationModel) -> None:
        """Register a new recommendation model"""
        self.models[model.name] = model
        logger.info(f"Registered model: {model.name}")
    
    def get_available_models(self) -> List[str]:
        """Get list of available model names"""
        return list(self.models.keys())
    
    def get_model_info(self, model_name: str) -> Dict[str, str]:
        """Get information about a specific model"""
        if model_name not in self.models:
            raise ValueError(f"Model '{model_name}' not found")
        return self.models[model_name].get_info()
    
    def set_active_model(self, model_name: str) -> None:
        """Set the active recommendation model"""
        if model_name not in self.models:
            raise ValueError(f"Model '{model_name}' not found")
        
        self.active_model = self.models[model_name]
        logger.info(f"Active model set to: {model_name}")
    
    def fit(self, papers_df: pd.DataFrame, model_name: Optional[str] = None, log_mlflow: bool = True, **kwargs) -> None:
        """Train a specific model or the active model"""
        if model_name:
            if model_name not in self.models:
                raise ValueError(f"Model '{model_name}' not found")
            self.models[model_name].fit(papers_df, log_mlflow=log_mlflow, experiment_name=self.experiment_name, **kwargs)
        elif self.active_model:
            self.active_model.fit(papers_df, log_mlflow=log_mlflow, experiment_name=self.experiment_name, **kwargs)
        else:
            raise ValueError("No model specified and no active model set")
    
    def recommend(self, paper_index: int, top_n: int = 6, model_name: Optional[str] = None, 
                 log_prediction: bool = False, **kwargs) -> pd.DataFrame:
        """Get recommendations using specified model or active model"""
        if model_name:
            if model_name not in self.models:
                raise ValueError(f"Model '{model_name}' not found")
            model = self.models[model_name]
        elif self.active_model:
            model = self.active_model
        else:
            raise ValueError("No model specified and no active model set")
        
        if not model.is_trained:
            raise ValueError(f"Model '{model.name}' is not trained")
        
        recommendations = model.recommend(paper_index, top_n, **kwargs)
        
        # Optionally log prediction to MLflow
        if log_prediction and model.mlflow_run_id:
            with mlflow.start_run(run_id=model.mlflow_run_id, nested=True) as run:
                mlflow.log_param("query_paper_index", paper_index)
                mlflow.log_param("top_n", top_n)
                mlflow.log_metric("num_recommendations", len(recommendations))
                if 'similarity_score' in recommendations.columns:
                    mlflow.log_metric("avg_similarity_score", recommendations['similarity_score'].mean())
                    mlflow.log_metric("min_similarity_score", recommendations['similarity_score'].min())
                    mlflow.log_metric("max_similarity_score", recommendations['similarity_score'].max())
        
        return recommendations
    
    def evaluate_model(self, papers_df: pd.DataFrame, test_queries: List[int], 
                      model_name: Optional[str] = None, top_n: int = 6) -> Dict[str, float]:
        """Evaluate model performance on test queries"""
        if model_name:
            if model_name not in self.models:
                raise ValueError(f"Model '{model_name}' not found")
            model = self.models[model_name]
        elif self.active_model:
            model = self.active_model
        else:
            raise ValueError("No model specified and no active model set")
        
        if not model.is_trained:
            raise ValueError(f"Model '{model.name}' is not trained")
        
        with mlflow.start_run(run_name=f"evaluation_{model.name}_{datetime.now().strftime('%Y%m%d_%H%M%S')}") as run:
            total_recommendations = 0
            avg_similarity_scores = []
            
            for query_idx in test_queries:
                try:
                    recommendations = model.recommend(query_idx, top_n)
                    total_recommendations += len(recommendations)
                    
                    if 'similarity_score' in recommendations.columns:
                        avg_similarity_scores.extend(recommendations['similarity_score'].tolist())
                        
                except Exception as e:
                    logger.warning(f"Error evaluating query {query_idx}: {e}")
            
            # Calculate evaluation metrics
            metrics = {
                "num_test_queries": len(test_queries),
                "total_recommendations": total_recommendations,
                "avg_recommendations_per_query": total_recommendations / len(test_queries) if test_queries else 0,
                "overall_avg_similarity": np.mean(avg_similarity_scores) if avg_similarity_scores else 0,
                "overall_std_similarity": np.std(avg_similarity_scores) if avg_similarity_scores else 0
            }
            
            # Log evaluation metrics
            mlflow.log_param("model_evaluated", model.name)
            mlflow.log_param("top_n", top_n)
            for metric_name, metric_value in metrics.items():
                mlflow.log_metric(metric_name, metric_value)
            
            logger.info(f"Model evaluation completed: {metrics}")
            return metrics
    
    def load_model_from_mlflow(self, model_name: str, model_version: Optional[str] = None) -> None:
        """Load a model from MLflow model registry"""
        try:
            if model_version:
                model_uri = f"models:/{model_name}/{model_version}"
            else:
                model_uri = f"models:/{model_name}/latest"
            
            loaded_model = mlflow.pyfunc.load_model(model_uri)
            
            # Note: This is a simplified example. In practice, you'd need to reconstruct
            # the full model object with its similarity matrix and papers_df
            logger.info(f"Model loaded from MLflow: {model_uri}")
            return loaded_model
            
        except Exception as e:
            logger.error(f"Error loading model from MLflow: {e}")
            raise

# Factory function to create pre-configured recommendation engine
def create_default_engine(experiment_name: str = "recommendation_experiments") -> RecommendationEngine:
    """Create a recommendation engine with default models"""
    engine = RecommendationEngine(experiment_name=experiment_name)
    
    # Register available models
    engine.register_model(SimilarityMatrixModel())
    
    """
    try:
        engine.register_model(ContentBasedModel())
    except ImportError:
        logger.warning("scikit-learn not available, ContentBasedModel not registered")
    """
    
    # Set default active model
    engine.set_active_model("Similarity Matrix")
    return engine

# Utility functions for MLflow management
def list_experiments() -> pd.DataFrame:
    """List all MLflow experiments"""
    experiments = mlflow.search_experiments()
    return pd.DataFrame([{
        'experiment_id': exp.experiment_id,
        'name': exp.name,
        'lifecycle_stage': exp.lifecycle_stage,
        'creation_time': datetime.fromtimestamp(exp.creation_time / 1000)
    } for exp in experiments])

def list_runs(experiment_name: str) -> pd.DataFrame:
    """List runs for a specific experiment"""
    experiment = mlflow.get_experiment_by_name(experiment_name)
    if not experiment:
        raise ValueError(f"Experiment '{experiment_name}' not found")
    
    runs = mlflow.search_runs(experiment_ids=[experiment.experiment_id])
    return runs

def compare_models(experiment_name: str, metric_name: str = "mean_similarity") -> pd.DataFrame:
    """Compare models based on a specific metric"""
    runs_df = list_runs(experiment_name)
    
    comparison = runs_df[['run_id', 'run_name', 'status', f'metrics.{metric_name}', 'start_time']].copy()
    comparison = comparison.dropna(subset=[f'metrics.{metric_name}'])
    comparison = comparison.sort_values(f'metrics.{metric_name}', ascending=False)
    
    return comparison
