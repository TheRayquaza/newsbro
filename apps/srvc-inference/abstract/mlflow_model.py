from typing import Any, Optional

import mlflow
import mlflow.pyfunc


class MlflowModel:
    def __init__(
        self, model_uri: str, tracking_uri: str, flavor: Optional[str] = None
    ) -> None:
        """
        Load MLflow model with automatic flavor detection.

        Args:
            model_uri: URI of the model (e.g., "models:/model-name/version")
            tracking_uri: MLflow tracking server URI
            flavor: Optional specific flavor to use (pyfunc, sklearn, pytorch, tensorflow, etc.)
        """
        print(f"Loading model from URI: {model_uri} with tracking URI: {tracking_uri}")
        mlflow.set_tracking_uri(tracking_uri)

        if flavor:
            self.model = self._load_with_flavor(model_uri, flavor)
        else:
            self.model = self._auto_load_model(model_uri)

    def _auto_load_model(self, model_uri: str) -> Any:
        """
        Automatically detect and load model with appropriate flavor.
        """
        try:
            return mlflow.pyfunc.load_model(model_uri)
        except Exception as e:
            print(f"pyfunc loading failed: {e}")

            try:
                client = mlflow.MlflowClient()
                if model_uri.startswith("models:/"):
                    parts = model_uri.replace("models:/", "").split("/")
                    model_name = parts[0]
                    version = parts[1] if len(parts) > 1 else None

                    if version:
                        model_version = client.get_model_version(model_name, version)
                        run_id = model_version.run_id
                    else:
                        versions = client.search_model_versions(f"name='{model_name}'")
                        if not versions:
                            raise ValueError(
                                f"No versions found for model {model_name}"
                            )
                        run_id = versions[0].run_id

                    run = client.get_run("" if not run_id else run_id)
                    flavors = list(run.data.params.keys())
                    print(f"Available flavors in model: {flavors}")

                elif model_uri.startswith("runs:/"):
                    run_id = model_uri.split("/")[1]
                    if not run_id:
                        raise ValueError(f"Invalid runs:/ URI: {model_uri}")
                    run = client.get_run(run_id)
                    flavors = list(run.data.params.keys())
                    print(f"Available flavors in model: {flavors}")

            except Exception as flavor_e:
                print(f"Could not detect flavors: {flavor_e}")

            for flavor in [
                "sklearn",
                "pytorch",
                "tensorflow",
                "keras",
                "xgboost",
                "lightgbm",
            ]:
                try:
                    print(f"Trying to load with {flavor} flavor...")
                    return self._load_with_flavor(model_uri, flavor)
                except Exception as flavor_err:
                    print(f"{flavor} loading failed: {flavor_err}")
                    continue
        raise ValueError(
            f"Could not load model from {model_uri} with any known flavor."
        )

    def _load_with_flavor(self, model_uri: str, flavor: str) -> Any:
        """Load model using a specific flavor."""
        flavor_loaders = {
            "pyfunc": lambda uri: mlflow.pyfunc.load_model(uri),
            "sklearn": lambda uri: mlflow.sklearn.load_model(uri),
            "pytorch": lambda uri: mlflow.pytorch.load_model(uri),
            "tensorflow": lambda uri: mlflow.tensorflow.load_model(uri),
            "keras": lambda uri: mlflow.keras.load_model(uri),
            "xgboost": lambda uri: mlflow.xgboost.load_model(uri),
            "lightgbm": lambda uri: mlflow.lightgbm.load_model(uri),
            "spark": lambda uri: mlflow.spark.load_model(uri),
        }

        if flavor not in flavor_loaders:
            raise ValueError(
                f"Unknown flavor: {flavor}. Available: {list(flavor_loaders.keys())}"
            )

        return flavor_loaders[flavor](model_uri)

    def predict(self, input_data: Any) -> Any:
        """Make predictions using the loaded model."""
        if hasattr(self.model, "predict"):
            return self.model.predict(input_data)
        elif callable(self.model):
            return self.model(input_data)
        else:
            raise AttributeError(
                "Model does not have a predict method or is not callable"
            )

    def __getattr__(self, name: str) -> Any:
        """
        Redirect any attribute access to the underlying model.
        This is called when an attribute is not found in MlflowModel.
        """
        return getattr(self.model, name)

    def __call__(self, *args: Any, **kwargs: Any) -> Any:
        """
        Make the wrapper callable, redirecting to the model's __call__ method if it exists.
        """
        if callable(self.model):
            return self.model(*args, **kwargs)
        return self.predict(*args, **kwargs)

    def __dir__(self):
        """
        Show both wrapper and model attributes when using dir() or autocomplete.
        """
        return list(set(super().__dir__() + dir(self.model)))
