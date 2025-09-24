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

    print("Loading news and behaviors Dataframes")
    col_names = ["impression_id", "user_id", "time", "history", "impressions"]
    behaviors = pd.read_csv("../data/behaviors.tsv", sep="\t", names=col_names, quoting=3)

    col_names = ["id", "category", "subcategory", "title", "abstract", "link", "title_entities", "abstract_entities"]
    news = pd.read_csv("../data/news.tsv", sep="\t", names=col_names, quoting=3)
    
    print("Loading entity and relation vectors")
    """
    def load_vec_file(path):
        embeddings = {}
        with open(path, "r", encoding="utf-8") as f:
            for line in f:
                parts = line.strip().split()
                if len(parts) <= 2:
                    continue
                key = parts[0]
                vec = np.array(list(map(float, parts[1:])))
                embeddings[key] = vec
        return embeddings
    
    #entity_embeddings = load_vec_file("../data/entity_embedding.vec")
    #relation_embeddings = load_vec_file("../data/relation_embedding.vec")
    """
    print("Reducing the dataset")
    
    # --- Step 1: Stratified sampling by category/subcategory ---
    sample_size = 5000
    news = news.groupby(['category', 'subcategory'], group_keys=False).apply(
        lambda x: x.sample(frac=sample_size/len(news), random_state=42)
    )
    news = news.reset_index(drop=True)
    valid_news_ids = set(news['id'])

    # --- Step 2: Filter histories and impressions ---
    def filter_history(history_str, valid_ids):
        if pd.isna(history_str) or history_str == "":
            return ""
        return " ".join([nid for nid in history_str.split() if nid in valid_ids])

    def filter_impressions(impressions_str, valid_ids):
        if pd.isna(impressions_str) or impressions_str == "":
            return ""
        return " ".join([pair for pair in impressions_str.split() if pair.split("-")[0] in valid_ids])

    behaviors['history'] = behaviors['history'].apply(lambda x: filter_history(x, valid_news_ids))
    behaviors['impressions'] = behaviors['impressions'].apply(lambda x: filter_impressions(x, valid_news_ids))

    # --- Step 3: Remove behaviors with no impressions left ---
    behaviors = behaviors[behaviors['impressions'] != ""]
    
    print("Dropping unused columns and null values")
    news = news.dropna(subset=['abstract'])

    columns_to_drop = ['category', 'subcategory', 'title_entities', 'abstract_entities']
    news = news.drop(columns=columns_to_drop, errors='ignore')
    
    behaviors = behaviors.dropna(subset=['history'])
    
    news['content'] = news['title'] + " " + news['abstract']
    news = news.reset_index(drop=True)
    
    news.to_csv('../data/news_cleaned.csv')
    behaviors.to_csv('../data/behaviors_cleaned.csv')
    
    print("FINISHING PREPROCESSING")

