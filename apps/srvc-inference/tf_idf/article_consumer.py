import os
import json
import time
import logging
import pandas as pd
from pathlib import Path
from kafka import KafkaConsumer
from sklearn.feature_extraction.text import TfidfVectorizer
from qdrant_client import QdrantClient
from qdrant_client.http.models import PointStruct, Distance, VectorParams
from threading import Lock, Thread
from dotenv import load_dotenv

load_dotenv()

# --- Configuration from ENV ---
QDRANT_URL = os.getenv("QDRANT_URL", "http://localhost:6333")
KAFKA_BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
BATCH_SIZE = int(os.getenv("BATCH_SIZE", "50"))
BATCH_INTERVAL = int(os.getenv("BATCH_INTERVAL", "2"))
TFIDF_MAX_FEATURES = int(os.getenv("TFIDF_MAX_FEATURES", "12"))
CSV_DATA_PATH = os.getenv("CSV_DATA_PATH", "../data/news_cleaned.csv")

# --- Logging ---
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

# --- Qdrant Client ---
client = QdrantClient(url=QDRANT_URL)
collection_name = "articles"

# --- Initialize TF-IDF from CSV ---
def load_and_train_vectorizer():
    """Load CSV data and train TF-IDF vectorizer"""
    csv_path = Path(CSV_DATA_PATH)
    
    if not csv_path.exists():
        logger.warning(f"CSV file not found at {CSV_DATA_PATH}. Using minimal corpus.")
        initial_corpus = [
            "First sample article to initialize the vectorizer",
            "Second sample article with different content"
        ]
    else:
        logger.info(f"Loading training data from {CSV_DATA_PATH}")
        try:
            df = pd.read_csv(csv_path)
            
            # Validate required columns
            if 'abstract' not in df.columns:
                logger.error("CSV must contain 'abstract' column")
                raise ValueError("Missing 'abstract' column in CSV")
    
            initial_corpus = df['abstract'].head(1000).tolist()  # Use first 1000 for training
            logger.info(f"Loaded {len(initial_corpus)} documents for TF-IDF training")
            
        except Exception as e:
            logger.error(f"Error loading CSV: {e}. Using minimal corpus.")
            initial_corpus = [
                "First sample article to initialize the vectorizer",
                "Second sample article with different abstract"
            ]
    
    # Train vectorizer
    vectorizer = TfidfVectorizer(
        max_features=TFIDF_MAX_FEATURES,
        stop_words='english',
        min_df=2,
        max_df=0.8
    )
    vectorizer.fit(initial_corpus)
    logger.info(f"TF-IDF vectorizer trained with {TFIDF_MAX_FEATURES} features.")
    
    return vectorizer

# --- Create Qdrant Collection ---
def create_collection_if_not_exists():
    """Create Qdrant collection with proper vector size"""
    try:
        collections = [c.name for c in client.get_collections().collections]
        if collection_name not in collections:
            client.create_collection(
                collection_name=collection_name,
                vectors_config=VectorParams(
                    size=TFIDF_MAX_FEATURES,
                    distance=Distance.COSINE
                )
            )
            logger.info(f"Collection '{collection_name}' created with vector size {TFIDF_MAX_FEATURES}.")
        else:
            # Verify vector size matches
            collection_info = client.get_collection(collection_name)
            vector_size = collection_info.config.params.vectors.size
            if vector_size != TFIDF_MAX_FEATURES:
                logger.error(
                    f"Collection exists but vector size mismatch: "
                    f"expected {TFIDF_MAX_FEATURES}, got {vector_size}"
                )
                raise ValueError("Vector size mismatch with existing collection")
            logger.info(f"Collection '{collection_name}' already exists.")
    except Exception as e:
        logger.error(f"Error managing collection: {e}")
        raise

# Initialize
vectorizer = load_and_train_vectorizer()
create_collection_if_not_exists()

# --- Kafka Consumer ---
consumer = KafkaConsumer(
    'article-aggregate',
    bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
    auto_offset_reset='earliest',
    enable_auto_commit=True,
    group_id='tfidf-batch-group',
    value_deserializer=lambda x: json.loads(x.decode('utf-8'))
)
logger.info("Kafka consumer started on topic 'article-aggregate'.")

# --- Batch Processing ---
batch_lock = Lock()
batch_buffer = []

def process_batch_data(batch):
    """Process a batch of articles: upsert or delete"""
    try:
        to_upsert = [a for a in batch if a.get("action", "add") in ["add", "update"]]
        to_delete = [a for a in batch if a.get("action") == "delete"]

        if to_upsert:
            texts = [a.get("abstract", "") for a in to_upsert]
            vectors = vectorizer.transform(texts).toarray().tolist()
            
            # Validate vector dimensions
            for v in vectors:
                if len(v) != TFIDF_MAX_FEATURES:
                    logger.error(f"Vector dimension mismatch: expected {TFIDF_MAX_FEATURES}, got {len(v)}")
                    return
            
            points = [
                PointStruct(
                    id=int(a["id"]),
                    vector=v,
                    payload={
                        "title": a.get("title", ""),
                        "abstract": a.get("abstract", ""),
                        "category": a.get("category", ""),
                        "created_at": a.get("created_at", "")
                    }
                )
                for a, v in zip(to_upsert, vectors)
            ]
            client.upsert(collection_name=collection_name, points=points)
            logger.info(f"[+] Upserted batch of {len(points)} articles")

        if to_delete:
            ids_to_delete = [int(a["id"]) for a in to_delete]
            client.delete(collection_name=collection_name, points_selector=ids_to_delete)
            logger.info(f"[-] Deleted batch of {len(ids_to_delete)} articles")

    except Exception as e:
        logger.error(f"Error processing batch: {e}", exc_info=True)

def periodic_batch_processor():
    """Process batch periodically"""
    global batch_buffer
    while True:
        time.sleep(BATCH_INTERVAL)
        with batch_lock:
            if not batch_buffer:
                continue
            batch = batch_buffer[:]
            batch_buffer = []
        
        process_batch_data(batch)

# Start periodic batch processor
Thread(target=periodic_batch_processor, daemon=True).start()
logger.info("Periodic batch processor thread started.")

# --- Main loop: read from Kafka ---
try:
    for message in consumer:
        article = message.value
        
        with batch_lock:
            batch_buffer.append(article)
            # Process immediately if batch is full
            if len(batch_buffer) >= BATCH_SIZE:
                batch = batch_buffer[:]
                batch_buffer = []
                # Process in separate thread to avoid blocking
                Thread(target=process_batch_data, args=(batch,), daemon=True).start()
                
except KeyboardInterrupt:
    logger.info("Consumer shutdown requested.")
finally:
    consumer.close()
    logger.info("Kafka consumer closed.")
