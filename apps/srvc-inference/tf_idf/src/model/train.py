import sys

import pandas as pd
from tqdm import tqdm

from tf_idf.src.model.model import TfidfSVDPipelineModel

tqdm.pandas()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 -m tf_idf.src.model.train <path_to_csv>")
        sys.exit(1)

    csv_file = sys.argv[1]
    print(f"📘 Loading dataset from {csv_file}")
    df = pd.read_csv(csv_file)

    df["content"] = (
        df["title"].fillna("") + " " + df["abstract"].fillna("")
    ).str.strip()
    df = df[df["content"].str.len() > 0].reset_index(drop=True)
    print(f"📊 Loaded {len(df)} documents")

    model = TfidfSVDPipelineModel(svd_components=100)
    model.fit(df["content"].progress_apply(str))

    model.save_to_mlflow(
        experiment_name="tfidf_experiment_4",
        tracking_uri="http://mlflow.localhost:8084",
    )

    print("✅ TF-IDF + SVD pipeline trained and saved to MLflow successfully.")
