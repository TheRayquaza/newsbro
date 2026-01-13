import sys
import pandas as pd
from tqdm import tqdm

from als.src.model.model import ALSModel

tqdm.pandas()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 -m als.src.model.train <path_to_interactions_csv>")
        print("\nExpected CSV format:")
        print("  - Columns: user_id, item_id (or news_id), rating (optional, defaults to 1)")
        print("  - Example: user_id,item_id,rating")
        print("            1,101,1")
        print("            1,102,1")
        print("            2,101,1")
        sys.exit(1)

    csv_file = sys.argv[1]
    print(f"📘 Loading interaction data from {csv_file}")
    
    df = pd.read_csv(csv_file)
    
    # Validate required columns
    if 'user_id' not in df.columns:
        print("❌ Error: CSV must contain 'user_id' column")
        sys.exit(1)
    
    # Handle item_id column (could be 'item_id' or 'news_id')
    if 'item_id' not in df.columns and 'news_id' in df.columns:
        df['item_id'] = df['news_id']
    elif 'item_id' not in df.columns:
        print("❌ Error: CSV must contain 'item_id' or 'news_id' column")
        sys.exit(1)
    
    # Add rating column if not present (implicit feedback = 1)
    if 'rating' not in df.columns:
        df['rating'] = 1.0
        print("ℹ️  No 'rating' column found, using implicit feedback (rating=1)")
    
    # Filter to only positive interactions
    df = df[df['rating'] > 0].reset_index(drop=True)
    
    print(f"📊 Loaded {len(df)} positive interactions")
    print(f"   - Unique users: {df['user_id'].nunique()}")
    print(f"   - Unique items: {df['item_id'].nunique()}")
    print(f"   - Sparsity: {100 * (1 - len(df) / (df['user_id'].nunique() * df['item_id'].nunique())):.2f}%")
    
    # Initialize and train ALS model
    print("\n🔧 Training ALS model...")
    model = ALSModel(
        n_factors=50,           # Number of latent factors
        regularization=0.01,    # L2 regularization
        iterations=15,          # Number of ALS iterations
        random_state=42
    )
    
    model.fit(df[['user_id', 'item_id', 'rating']])
    
    # Log to MLflow
    print("\n💾 Saving model to MLflow...")
    model.save_to_mlflow(
        experiment_name="als_experiment_1",
        tracking_uri="http://mlflow.localhost:8084",
    )

    print("\n✅ ALS model trained and saved to MLflow successfully!")
    
    # Show some example recommendations
    print("\n🎯 Example recommendations:")
    sample_user = df['user_id'].iloc[0]
    recommendations = model.recommend(sample_user, top_k=5)
    print(f"   Top 5 items for user {sample_user}:")
    for item_id, score in recommendations:
        print(f"      - Item {item_id}: score={score:.4f}")
