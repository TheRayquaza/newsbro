import pickle  # nosec B403
import sys

from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity


def load_model(model_path: str):
    """Load a saved SBERT model from a pickle file or load pre-trained model."""
    try:
        # Try loading from pickle first
        with open(model_path, "rb") as f:
            model = pickle.load(
                f
            )  # nosec B301 - safe here, model file is locally generated and trusted
        print(f"✅ Model loaded from {model_path}")
    except Exception:
        # If pickle fails, assume it's a model name and load from sentence-transformers
        print(f"Loading pre-trained model: {model_path}")
        model = SentenceTransformer(model_path)
        print("✅ Pre-trained model loaded")

    return model


def compute_similarity(model, text1: str, text2: str) -> float:
    """Compute cosine similarity between two texts using the SBERT model."""
    embeddings = model.encode([text1, text2], convert_to_numpy=True)
    similarity = cosine_similarity(embeddings[0:1], embeddings[1:2])[0][0]
    return float(similarity)


def test_batch_encoding(model, texts: list) -> None:
    """Test batch encoding with multiple texts."""
    print(f"\n📦 Testing batch encoding with {len(texts)} texts...")
    embeddings = model.encode(texts, show_progress_bar=True, convert_to_numpy=True)
    print(f"✅ Generated {len(embeddings)} embeddings")
    print(f"   Embedding shape: {embeddings[0].shape}")
    print(f"   Embedding dimension: {embeddings.shape[1]}")

    # Compute pairwise similarities
    print("\n🔍 Pairwise similarities:")
    for i in range(len(texts)):
        for j in range(i + 1, len(texts)):
            sim = cosine_similarity(embeddings[i : i + 1], embeddings[j : j + 1])[0][0]
            print(f"   Text {i+1} ↔ Text {j+1}: {sim:.4f}")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python test.py <path_to_model.pkl or model_name>")
        print("\nExamples:")
        print("  python test.py all-MiniLM-L6-v2")
        print("  python test.py model.pkl")
        sys.exit(1)

    model_path = sys.argv[1]
    model = load_model(model_path)

    print("\n📊 Model information:")
    if hasattr(model, "get_sentence_embedding_dimension"):
        print(f"   Embedding dimension: {model.get_sentence_embedding_dimension()}")
    if hasattr(model, "max_seq_length"):
        print(f"   Max sequence length: {model.max_seq_length}")

    # Interactive similarity test
    print("\n📝 Example texts for similarity check:")
    text_a = input("Enter first text: ").strip()
    text_b = input("Enter second text: ").strip()

    if text_a and text_b:
        sim = compute_similarity(model, text_a, text_b)
        print(f"\n🔍 Cosine similarity: {sim:.4f}")

    # Batch test
    print("\n" + "=" * 50)
    print("Would you like to test batch encoding? (y/n): ", end="")
    if input().strip().lower() == "y":
        print("\nEnter texts (one per line, empty line to finish):")
        texts: list[str] = []
        while True:
            text = input(f"Text {len(texts)+1}: ").strip()
            if not text:
                break
            texts.append(text)

        if len(texts) >= 2:
            test_batch_encoding(model, texts)
        else:
            print("Need at least 2 texts for batch testing")
