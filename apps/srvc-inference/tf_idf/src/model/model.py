import mlflow
import mlflow.sklearn
from sklearn.decomposition import TruncatedSVD
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.pipeline import Pipeline
from sklearn.preprocessing import Normalizer

from tf_idf.src.model.preprocessing import TextCleaner


class TfidfSVDPipelineModel:
    """
    Wrapper for a unified sklearn pipeline:
    TextCleaner -> TF-IDF -> TruncatedSVD -> Normalizer
    """

    def __init__(self, max_features=None, svd_components=100, stop_words="english"):
        self.max_features = max_features
        self.svd_components = svd_components
        self.stop_words = stop_words
        self.pipeline = None

    def build(self):
        """Builds the sklearn pipeline"""
        self.pipeline = Pipeline(
            [
                ("cleaner", TextCleaner()),
                (
                    "tfidf",
                    TfidfVectorizer(
                        stop_words=self.stop_words, max_features=self.max_features
                    ),
                ),
                (
                    "svd",
                    TruncatedSVD(n_components=self.svd_components, random_state=42),
                ),
                ("normalize", Normalizer(copy=False)),
            ]
        )
        return self.pipeline

    def fit(self, documents):
        """Fits the pipeline on training documents"""
        if self.pipeline is None:
            self.build()
        self.pipeline.fit(documents)
        return self

    def transform(self, documents):
        """Transforms raw documents into dense vectors"""
        if self.pipeline is None:
            raise RuntimeError("Pipeline is not trained or built yet.")
        return self.pipeline.transform(documents)

    def predict(self, documents):
        """MLflow requires predict() to log pyfunc flavor"""
        return self.transform(documents)
    
    def save_to_mlflow(
        self,
        experiment_name="tfidf_pipeline_experiment",
        run_name="TFIDF_Pipeline_Training",
        tracking_uri="https://mlflow.internal.newsbro.cc",
    ):
        """Logs the pipeline to MLflow"""
        mlflow.set_tracking_uri(tracking_uri)
        mlflow.set_experiment(experiment_name)
        with mlflow.start_run(run_name=run_name):
            mlflow.sklearn.log_model(self, "tfidf_full_pipeline")
