import sys

import pandas as pd
from tqdm import tqdm

from sbert.src.model.model import SBERTModel

tqdm.pandas()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 -m sbert.src.model.train <path_to_csv> [model_name]")
        print("\nAvailable model options:")
        print("  - all-MiniLM-L6-v2 (default): Fast, 384 dimensions")
        print("  - all-mpnet-base-v2: Better quality, 768 dimensions")
        print("  - paraphrase-multilingual-MiniLM-L12-v2: Multilingual support")
        sys.exit(1)

    csv_file = sys.argv[1]
    model_name = sys.argv[2] if len(sys.argv) > 2 else "all-MiniLM-L6-v2"
    
    print(f"📘 Loading dataset from {csv_file}")
    df = pd.read_csv(csv_file)

    df["content"] = (
        df["title"].fillna("") + " " + df["abstract"].fillna("")
    ).str.strip()
    df = df[df["content"].str.len() > 0].reset_index(drop=True)
    print(f"📊 Loaded {len(df)} documents")

    print(f"🤖 Loading SBERT model: {model_name}")
    model = SBERTModel(model_name=model_name)
    model.load()
    
    print(f"📐 Embedding dimension: {model.get_embedding_dimension()}")
    
    # Generate embeddings for a sample to verify the model works
    print("🔍 Testing model with sample data...")
    sample_size = min(10, len(df))
    sample_embeddings = model.encode(
        df["content"].head(sample_size).tolist(),
        show_progress_bar=True
    )
    print(f"✅ Generated {len(sample_embeddings)} sample embeddings of shape {sample_embeddings[0].shape}")

    # Save to MLflow
    print("💾 Saving model to MLflow...")
    model.save_to_mlflow(
        experiment_name="sbert_experiment",
        tracking_uri="https://mlflow.internal.newsbro.cc",
    )

    print("✅ SBERT model prepared and saved to MLflow successfully.")
    print(f"\nModel details:")
    print(f"  - Model name: {model_name}")
    print(f"  - Embedding dimension: {model.get_embedding_dimension()}")
    print(f"  - Max sequence length: {model.model.max_seq_length}")
