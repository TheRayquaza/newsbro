import pickle  # nosec B403
import sys

from sklearn.metrics.pairwise import cosine_similarity


def load_model(model_path: str):
    """Load a saved TF-IDF + SVD pipeline model from a pickle file."""
    with open(model_path, "rb") as f:
        model = pickle.load(
            f
        )  # nosec B301 - safe here, model file is locally generated and trusted
    print(f"✅ Model loaded from {model_path}")
    return model


def compute_similarity(model, text1: str, text2: str) -> float:
    """Compute cosine similarity between two texts using the trained model."""
    vectors = model.transform([text1, text2])
    similarity = cosine_similarity(vectors[0:1], vectors[1:2])[0][0]
    return float(similarity)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python test_tfidf_svd_similarity.py <path_to_model.pkl>")
        sys.exit(1)

    model_path = sys.argv[1]
    model = load_model(model_path)

    print("\n📝 Example texts for similarity check:")
    text_a = input("Enter first text: ").strip()
    text_b = input("Enter second text: ").strip()

    sim = compute_similarity(model, text_a, text_b)
    print(f"\n🔍 Cosine similarity: {sim:.4f}")
