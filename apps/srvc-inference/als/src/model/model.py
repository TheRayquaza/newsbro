import mlflow
import mlflow.pyfunc
import numpy as np
from typing import Optional, Dict
import pandas as pd


class ALSModel:
    """
    Alternating Least Squares (ALS) model for implicit feedback recommendation.
    
    This implements a matrix factorization approach where:
    - User-item interactions are represented as a sparse matrix R
    - User factors U and item factors I are learned such that R ≈ U @ I.T
    - ALS alternates between fixing U and solving for I, then fixing I and solving for U
    """

    def __init__(
        self,
        n_factors: int = 50,
        regularization: float = 0.01,
        iterations: int = 15,
        random_state: int = 42,
    ):
        """
        Initialize ALS model.
        
        Args:
            n_factors: Number of latent factors
            regularization: Regularization parameter (lambda)
            iterations: Number of ALS iterations
            random_state: Random seed for reproducibility
        """
        self.n_factors = n_factors
        self.regularization = regularization
        self.iterations = iterations
        self.random_state = random_state
        
        self.user_factors: Optional[np.ndarray] = None
        self.item_factors: Optional[np.ndarray] = None
        self.user_id_map: Optional[Dict[int, int]] = None
        self.item_id_map: Optional[Dict[int, int]] = None
        self.reverse_user_map: Optional[Dict[int, int]] = None
        self.reverse_item_map: Optional[Dict[int, int]] = None

    def fit(self, interactions: pd.DataFrame) -> 'ALSModel':
        """
        Fit the ALS model on user-item interactions.
        
        Args:
            interactions: DataFrame with columns ['user_id', 'item_id', 'rating']
                         For implicit feedback, rating is typically 1 for positive interactions
        
        Returns:
            self
        """
        np.random.seed(self.random_state)
        
        # Create mappings from original IDs to matrix indices
        unique_users = interactions['user_id'].unique()
        unique_items = interactions['item_id'].unique()
        
        self.user_id_map = {uid: idx for idx, uid in enumerate(unique_users)}
        self.item_id_map = {iid: idx for idx, iid in enumerate(unique_items)}
        self.reverse_user_map = {idx: uid for uid, idx in self.user_id_map.items()}
        self.reverse_item_map = {idx: iid for iid, idx in self.item_id_map.items()}
        
        n_users = len(unique_users)
        n_items = len(unique_items)
        
        # Build interaction matrix
        interaction_matrix = np.zeros((n_users, n_items), dtype=np.float32)
        for _, row in interactions.iterrows():
            user_idx = self.user_id_map[row['user_id']]
            item_idx = self.item_id_map[row['item_id']]
            interaction_matrix[user_idx, item_idx] = row['rating']
        
        # Initialize factors with small random values
        self.user_factors = np.random.normal(0, 0.1, (n_users, self.n_factors)).astype(np.float32)
        self.item_factors = np.random.normal(0, 0.1, (n_items, self.n_factors)).astype(np.float32)
        
        # ALS iterations
        reg_eye = self.regularization * np.eye(self.n_factors, dtype=np.float32)
        
        for iteration in range(self.iterations):
            # Fix item factors, solve for user factors
            for u in range(n_users):
                # Get items this user interacted with
                item_indices = np.where(interaction_matrix[u, :] > 0)[0]
                if len(item_indices) == 0:
                    continue
                
                item_factors_u = self.item_factors[item_indices]
                ratings_u = interaction_matrix[u, item_indices]
                
                # Solve: user_factor = (I^T @ I + λΙ)^(-1) @ I^T @ r
                gram = item_factors_u.T @ item_factors_u + reg_eye
                self.user_factors[u] = np.linalg.solve(gram, item_factors_u.T @ ratings_u)
            
            # Fix user factors, solve for item factors
            for i in range(n_items):
                # Get users who interacted with this item
                user_indices = np.where(interaction_matrix[:, i] > 0)[0]
                if len(user_indices) == 0:
                    continue
                
                user_factors_i = self.user_factors[user_indices]
                ratings_i = interaction_matrix[user_indices, i]
                
                # Solve: item_factor = (U^T @ U + λΙ)^(-1) @ U^T @ r
                gram = user_factors_i.T @ user_factors_i + reg_eye
                self.item_factors[i] = np.linalg.solve(gram, user_factors_i.T @ ratings_i)
            
            # Calculate loss (optional, for monitoring)
            if iteration % 5 == 0:
                predictions = self.user_factors @ self.item_factors.T
                mask = interaction_matrix > 0
                mse = np.mean((interaction_matrix[mask] - predictions[mask]) ** 2)
                reg_loss = self.regularization * (
                    np.sum(self.user_factors ** 2) + np.sum(self.item_factors ** 2)
                )
                print(f"Iteration {iteration}: MSE={mse:.4f}, Reg Loss={reg_loss:.4f}")
        
        return self

    def predict(self, user_id: int, item_id: int) -> float:
        """
        Predict rating for a user-item pair.
        
        Args:
            user_id: Original user ID
            item_id: Original item ID
        
        Returns:
            Predicted rating (dot product of user and item factors)
        """
        if user_id not in self.user_id_map or item_id not in self.item_id_map:
            return 0.0
        
        user_idx = self.user_id_map[user_id]
        item_idx = self.item_id_map[item_id]
        
        return float(self.user_factors[user_idx] @ self.item_factors[item_idx])

    def get_user_factor(self, user_id: int) -> Optional[np.ndarray]:
        """Get latent factor for a user."""
        if user_id not in self.user_id_map:
            return None
        user_idx = self.user_id_map[user_id]
        return self.user_factors[user_idx]

    def get_item_factor(self, item_id: int) -> Optional[np.ndarray]:
        """Get latent factor for an item."""
        if item_id not in self.item_id_map:
            return None
        item_idx = self.item_id_map[item_id]
        return self.item_factors[item_idx]

    def recommend(self, user_id: int, top_k: int = 10, exclude_seen: bool = True) -> list:
        """
        Recommend top K items for a user.
        
        Args:
            user_id: Original user ID
            top_k: Number of recommendations
            exclude_seen: Whether to exclude items user has already seen
        
        Returns:
            List of (item_id, score) tuples
        """
        if user_id not in self.user_id_map:
            return []
        
        user_idx = self.user_id_map[user_id]
        scores = self.user_factors[user_idx] @ self.item_factors.T
        
        # Get top K indices
        top_indices = np.argsort(scores)[-top_k:][::-1]
        
        recommendations = [
            (self.reverse_item_map[idx], float(scores[idx]))
            for idx in top_indices
        ]
        
        return recommendations

    def save_to_mlflow(
        self,
        experiment_name: str = "als_experiment",
        run_name: str = "ALS_Training",
        tracking_uri: str = "http://mlflow.localhost:8084",
    ):
        """Log the ALS model to MLflow."""
        mlflow.set_tracking_uri(tracking_uri)
        mlflow.set_experiment(experiment_name)
        
        with mlflow.start_run(run_name=run_name):
            # Log parameters
            mlflow.log_param("n_factors", self.n_factors)
            mlflow.log_param("regularization", self.regularization)
            mlflow.log_param("iterations", self.iterations)
            mlflow.log_param("random_state", self.random_state)
            
            # Log model artifacts
            mlflow.pyfunc.log_model(
                artifact_path="als_model",
                python_model=ALSMLflowWrapper(self),
            )
            
            print(f"✅ ALS model logged to MLflow experiment: {experiment_name}")


class ALSMLflowWrapper(mlflow.pyfunc.PythonModel):
    """Wrapper to make ALS model compatible with MLflow."""
    
    def __init__(self, als_model: ALSModel):
        self.als_model = als_model
    
    def predict(self, context, model_input):
        """
        Predict method for MLflow compatibility.
        Expects DataFrame with columns: user_id, item_id
        """
        predictions = []
        for _, row in model_input.iterrows():
            pred = self.als_model.predict(row['user_id'], row['item_id'])
            predictions.append(pred)
        return np.array(predictions)
    
    def get_item_factor(self, item_id: int) -> Optional[np.ndarray]:
        """Proxy method to get item factor."""
        return self.als_model.get_item_factor(item_id)
    
    def get_user_factor(self, user_id: int) -> Optional[np.ndarray]:
        """Proxy method to get user factor."""
        return self.als_model.get_user_factor(user_id)
