import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from sentence_transformers import SentenceTransformer
from sklearn.preprocessing import OneHotEncoder
from sklearn.metrics.pairwise import cosine_similarity
import os
import json

def preprocessing():
    
    print("STARTING PREPROCESSING\n")
    papers = pd.read_csv('../data/arXiv-DataFrame.csv')

    print("Dropping columns : ['Unnamed: 0', 'id']")
    papers.drop(columns=['Unnamed: 0', 'id'], inplace=True, errors='ignore')

    papers.drop_duplicates(keep='first', inplace=True)
    print("Dropping duplicates")


    papers[papers.duplicated(keep=False)] \
        .sort_values(['Title', 'Author'])

    print("Reducing data")

    # -------------------------------
    # 1. Limit authors to max 20 papers
    # -------------------------------
    max_author_count = 20

    def truncate_group(df, group_col, max_count):
        truncated_df = df.groupby(group_col).head(max_count)
        print(f"After truncating {group_col}: {len(truncated_df)} rows")
        return truncated_df

    papers = truncate_group(papers, 'Author', max_author_count)

    # -------------------------------
    # 2. Limit categories to max 100 papers
    # -------------------------------
    max_category_count = 100
    papers = truncate_group(papers, 'Primary Category', max_category_count)

    print("Changing Categories' Names")
    with open("../data/map_category.json", "r") as f:
        category_map = json.load(f)

    papers["Category"] = papers["Primary Category"].map(category_map).fillna("Other")


    print("Creating Text Embeddings : Title + Summary")
    
    # Combine Title and Summary
    papers['text'] = papers['Title'].fillna('') + " " + papers['Summary'].fillna('')

    try:
        # Try to load existing embeddings
        text_embeddings = np.load('../data/text_embeddings.npy')
        if not isinstance(text_embeddings, np.ndarray):
            raise ValueError("Loaded embeddings are not a NumPy array. Recomputing embeddings.")
    except (FileNotFoundError, ValueError):
        # If file doesn't exist or is invalid, compute embeddings
        model = SentenceTransformer('all-MiniLM-L6-v2') 
        text_embeddings = model.encode(papers['text'].tolist(), show_progress_bar=True)
        np.save('../data/text_embeddings.npy', text_embeddings)
    
    print("Creating Categorical Embeddings : Author + Primary Category + Category")
    categorical_cols = ['Author', 'Primary Category', 'Category']

    try:
        # Try to load existing categorical embeddings
        cat_embeddings = np.load('../data/cat_embeddings.npy', allow_pickle=True)
        if not isinstance(cat_embeddings, np.ndarray):
            raise ValueError("Loaded embeddings are not a NumPy array. Recomputing embeddings.")
    except (FileNotFoundError, ValueError):
        # If file doesn't exist or is invalid, compute embeddings
        encoder = OneHotEncoder(handle_unknown='ignore', sparse_output=False)
        cat_embeddings = encoder.fit_transform(papers[categorical_cols].fillna(''))
        np.save('../data/cat_embeddings.npy', cat_embeddings)

    combined_embeddings = np.hstack([text_embeddings, cat_embeddings])
    
    print("Creating Similarity Matrix")
    try:
        # Try to load existing similarity matrix
        similarity_matrix = np.load('../data/similarity_matrix.npy', allow_pickle=True)
        if not isinstance(cat_embeddings, np.ndarray):
            raise ValueError("Loaded similarity matrix is not a NumPy array. Recomputing similarity matrix.")
    except (FileNotFoundError, ValueError):
        similarity_matrix = cosine_similarity(combined_embeddings)
        np.save('../data/similarity_matrix.npy', similarity_matrix)
    
    file_path = "../data/papers.csv"

    if not os.path.exists(file_path):
        papers.to_csv(file_path, index=False)
        print(f"✅ Saved DataFrame to {file_path}")
    else:
        print(f"⚡ File {file_path} already exists, not overwriting.")
    
    print("FINISHING PREPROCESSING")

