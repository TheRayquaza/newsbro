import os, numpy as np, pandas as pd, mlflow
from kserve import Model, ModelServer
from typing import Dict, Any
from dotenv import load_dotenv

def parse_v2_inputs(payload: Dict[str, Any]):
    if "inputs" in payload:
        x = payload["inputs"][0]["data"]
        return np.array(x)
    if "instances" in payload:  # V1 compatibility
        inst = payload["instances"]
        # accept list-of-rows OR dict-of-columns
        if isinstance(inst, dict):
            return pd.DataFrame(inst)
        return np.array(inst)
    raise ValueError("Missing 'inputs' or 'instances'")

def to_v2_output(pred):
    arr = np.array(pred)
    return {
        "model_name": os.environ.get("MODEL_NAME", "model"),
        "outputs": [{
            "name": "output-0",
            "shape": list(arr.shape),
            "datatype": "FP32" if arr.dtype.kind == "f" else "INT64",
            "data": arr.tolist()
        }]
    }

class MlflowPyfunc(Model):
    def __init__(self, name, uri):
        super().__init__(name)
        self.uri = uri
        self.model = None

    async def load(self):
        try:
            print(f"Loading model from: {self.uri}")
            # Set MLflow tracking URI and S3 config for model loading
            mlflow.set_tracking_uri(os.getenv("MLFLOW_TRACKING_URI", "http://localhost:8080"))
            self.model = mlflow.pyfunc.load_model(self.uri)
            self.ready = True
            print(f"Model {self.name} loaded successfully")
        except Exception as e:
            print(f"Failed to load model {self.name}: {e}")
            self.ready = False
            raise e

    def predict(self, payload: Dict[str, Any], headers: Dict[str, str] = None):
        if not self.ready:
            raise Exception(f"Model {self.name} is not ready")
        X = parse_v2_inputs(payload)
        y = self.model.predict(X)
        return to_v2_output(y)

if __name__ == "__main__":
    load_dotenv()  # take environment variables from .env.
    
    # Set AWS/S3 credentials for MLflow
    os.environ["AWS_ACCESS_KEY_ID"] = os.getenv("AWS_ACCESS_KEY_ID", "minio")
    os.environ["AWS_SECRET_ACCESS_KEY"] = os.getenv("AWS_SECRET_ACCESS_KEY", "minio123")
    os.environ["MLFLOW_S3_ENDPOINT_URL"] = os.getenv("MLFLOW_S3_ENDPOINT_URL", "http://minio-svc.localhost:8080/")
    
    m = MlflowPyfunc(os.getenv("MODEL_NAME", "model"), os.getenv("MODEL_URI"))
    ModelServer(workers=1).start(models=[m])