import pickle  # nosec B403
import sys
import numpy as np


def load_model(model_path: str):
    """Load a saved ALS model from a pickle file."""
    with open(model_path, "rb") as f:
        model = pickle.load(
            f
        )  # nosec B301 - safe here, model file is locally generated and trusted
    print(f"✅ ALS model loaded from {model_path}")
    return model


def compute_similarity(model, user_id: int, item_id: int) -> float:
    """
    Compute predicted rating for a user-item pair using the trained ALS model.
    This is the dot product of user and item latent factors.
    """
    prediction = model.predict(user_id, item_id)
    return float(prediction)


def show_user_recommendations(model, user_id: int, top_k: int = 10):
    """Show top K recommendations for a user."""
    recommendations = model.recommend(user_id, top_k=top_k)
    
    if not recommendations:
        print(f"❌ No recommendations found for user {user_id}")
        return
    
    print(f"\n🎯 Top {top_k} recommendations for user {user_id}:")
    for idx, (item_id, score) in enumerate(recommendations, 1):
        print(f"   {idx}. Item {item_id}: score={score:.4f}")


def show_item_similar_users(model, item_id: int, top_k: int = 10):
    """Show top K users most likely to be interested in an item."""
    item_factor = model.get_item_factor(item_id)
    
    if item_factor is None:
        print(f"❌ Item {item_id} not found in model")
        return
    
    # Compute scores for all users
    scores = model.user_factors @ item_factor
    top_user_indices = np.argsort(scores)[-top_k:][::-1]
    
    print(f"\n👥 Top {top_k} users most interested in item {item_id}:")
    for idx, user_idx in enumerate(top_user_indices, 1):
        user_id = model.reverse_user_map[user_idx]
        score = scores[user_idx]
        print(f"   {idx}. User {user_id}: score={score:.4f}")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python test_als_similarity.py <path_to_model.pkl>")
        print("\nThis script tests the ALS model by:")
        print("  1. Computing predicted ratings for user-item pairs")
        print("  2. Showing top recommendations for users")
        print("  3. Finding users most interested in specific items")
        sys.exit(1)

    model_path = sys.argv[1]
    model = load_model(model_path)
    
    print(f"\n📊 Model Statistics:")
    print(f"   - Users: {len(model.user_id_map)}")
    print(f"   - Items: {len(model.item_id_map)}")
    print(f"   - Latent factors: {model.n_factors}")
    print(f"   - Regularization: {model.regularization}")
    print(f"   - Iterations: {model.iterations}")
    
    while True:
        print("\n" + "="*50)
        print("Choose an option:")
        print("  1. Predict rating for user-item pair")
        print("  2. Get recommendations for a user")
        print("  3. Find users interested in an item")
        print("  4. Exit")
        
        choice = input("\nEnter choice (1-4): ").strip()
        
        if choice == "1":
            try:
                user_id = int(input("Enter user ID: ").strip())
                item_id = int(input("Enter item ID: ").strip())
                
                score = compute_similarity(model, user_id, item_id)
                print(f"\n🔍 Predicted rating: {score:.4f}")
                
                if score > 0.5:
                    print("   → High interest (score > 0.5)")
                elif score > 0.2:
                    print("   → Moderate interest (0.2 < score < 0.5)")
                else:
                    print("   → Low interest (score < 0.2)")
                    
            except ValueError:
                print("❌ Invalid input. Please enter numeric IDs.")
        
        elif choice == "2":
            try:
                user_id = int(input("Enter user ID: ").strip())
                top_k = int(input("Enter number of recommendations (default 10): ").strip() or "10")
                show_user_recommendations(model, user_id, top_k)
            except ValueError:
                print("❌ Invalid input. Please enter numeric values.")
        
        elif choice == "3":
            try:
                item_id = int(input("Enter item ID: ").strip())
                top_k = int(input("Enter number of users (default 10): ").strip() or "10")
                show_item_similar_users(model, item_id, top_k)
            except ValueError:
                print("❌ Invalid input. Please enter numeric values.")
        
        elif choice == "4":
            print("\n👋 Goodbye!")
            break
        
        else:
            print("❌ Invalid choice. Please enter 1-4.")
