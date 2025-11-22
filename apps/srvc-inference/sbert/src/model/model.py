import mlflow
from dotenv import load_dotenv
from sentence_transformers import SentenceTransformer

load_dotenv()

mlflow.autolog()
mlflow.set_tracking_uri("https://mlfow.internal.newsbro.cc")


class SBERTModel:
    """
    Wrapper for Sentence-BERT model with MLflow integration.
    Uses pre-trained transformer models to generate sentence embeddings.
    """

    def __init__(self, model_name="all-MiniLM-L6-v2", device=None):
        """
        Initialize SBERT model.

        Args:
            model_name: Name of the pre-trained model from sentence-transformers
                       Popular options:
                       - 'all-MiniLM-L6-v2': Fast, 384 dimensions (default)
                       - 'all-mpnet-base-v2': Better quality, 768 dimensions
                       - 'paraphrase-multilingual-MiniLM-L12-v2': Multilingual
            device: Device to run model on ('cuda', 'cpu', or None for auto)
        """
        self.model_name = model_name
        self.device = device
        self.model = None

    def load(self):
        """Load the pre-trained SBERT model"""
        self.model = SentenceTransformer(self.model_name, device=self.device)
        return self

    def encode(self, texts, batch_size=32, show_progress_bar=False):
        """
        Encode texts into embeddings.

        Args:
            texts: List of strings or single string to encode
            batch_size: Batch size for encoding
            show_progress_bar: Whether to show progress bar

        Returns:
            numpy array of embeddings
        """
        if self.model is None:
            raise RuntimeError("Model is not loaded yet. Call load() first.")

        return self.model.encode(
            texts,
            batch_size=batch_size,
            show_progress_bar=show_progress_bar,
            convert_to_numpy=True,
        )

    def save_to_mlflow(
        self,
        experiment_name="sbert_experiment",
        run_name="SBERT_Model_Training",
        tracking_uri="http://mlflow.localhost:8084",
    ):
        """
        Save the SBERT model to MLflow.
        Note: This saves the model name and configuration, not retrained weights.
        """
        mlflow.set_tracking_uri(tracking_uri)
        mlflow.set_experiment(experiment_name)

        with mlflow.start_run(run_name=run_name):
            # Log model parameters
            mlflow.log_param("model_name", self.model_name)
            mlflow.log_param(
                "embedding_dimension", self.model.get_sentence_embedding_dimension()
            )
            mlflow.log_param("max_seq_length", self.model.max_seq_length)

            # Log the model using sentence-transformers flavor
            mlflow.sentence_transformers.log_model(
                self.model,
                artifact_path="sbert_model",
                registered_model_name=f"sbert_{self.model_name.replace('/', '_')}",
            )

            print(f"✅ Model logged to MLflow experiment '{experiment_name}'")

    def get_embedding_dimension(self):
        """Get the dimension of the output embeddings"""
        if self.model is None:
            raise RuntimeError("Model is not loaded yet.")
        return self.model.get_sentence_embedding_dimension()
