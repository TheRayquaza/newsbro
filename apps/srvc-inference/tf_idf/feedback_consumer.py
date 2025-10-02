import os
import json
import logging
from kafka import KafkaConsumer
from qdrant_client import QdrantClient
from qdrant_client.http.models import PointStruct, Distance, VectorParams, Filter, FieldCondition, Range
from datetime import datetime, timedelta
from dotenv import load_dotenv

load_dotenv()

# --- Configuration from ENV ---
QDRANT_URL = os.getenv("QDRANT_URL", "http://localhost:6333")
KAFKA_BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
FEEDBACK_RETENTION_DAYS = int(os.getenv("FEEDBACK_RETENTION_DAYS", "30"))

# --- Logging ---
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

# --- Qdrant Client ---
qdrant = QdrantClient(url=QDRANT_URL)
QDRANT_COLLECTION = "feedbacks"

# Create collection if not exists
def create_feedback_collection():
    """Create feedback collection with proper configuration"""
    try:
        collections = [c.name for c in qdrant.get_collections().collections]
        if QDRANT_COLLECTION not in collections:
            qdrant.create_collection(
                collection_name=QDRANT_COLLECTION,
                vectors_config=VectorParams(size=1, distance=Distance.COSINE)
            )
            # Create index on timestamp for faster cleanup queries
            qdrant.create_payload_index(
                collection_name=QDRANT_COLLECTION,
                field_name="timestamp",
                field_schema="float"
            )
            # Create indexes for user_id and article_id for faster filtering
            qdrant.create_payload_index(
                collection_name=QDRANT_COLLECTION,
                field_name="user_id",
                field_schema="integer"
            )
            qdrant.create_payload_index(
                collection_name=QDRANT_COLLECTION,
                field_name="type",
                field_schema="keyword"
            )
            logger.info(f"Collection '{QDRANT_COLLECTION}' created with indexes.")
        else:
            logger.info(f"Collection '{QDRANT_COLLECTION}' already exists.")
    except Exception as e:
        logger.error(f"Error creating collection: {e}")
        raise

create_feedback_collection()

# --- Kafka Consumer ---
consumer = KafkaConsumer(
    'feedback-aggregate',
    bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
    value_deserializer=lambda m: json.loads(m.decode('utf-8')),
    auto_offset_reset='earliest',
    enable_auto_commit=True,
    group_id='feedback-consumer-group'
)
logger.info(f"Kafka consumer started on topic 'feedback-aggregate'.")

def cleanup_old_feedbacks():
    """Delete feedbacks older than FEEDBACK_RETENTION_DAYS"""
    try:
        cutoff_timestamp = (datetime.now() - timedelta(days=FEEDBACK_RETENTION_DAYS)).timestamp()
        
        # Scroll through old points
        old_points, _ = qdrant.scroll(
            collection_name=QDRANT_COLLECTION,
            scroll_filter=Filter(
                must=[
                    FieldCondition(
                        key="timestamp",
                        range=Range(lt=cutoff_timestamp)
                    )
                ]
            ),
            limit=1000,
            with_payload=False,
            with_vectors=False
        )
        
        if old_points:
            old_ids = [p.id for p in old_points]
            qdrant.delete(
                collection_name=QDRANT_COLLECTION,
                points_selector=old_ids
            )
            logger.info(f"🗑️  {len(old_ids)} old feedbacks deleted (older than {FEEDBACK_RETENTION_DAYS} days).")
        else:
            logger.info(f"No old feedbacks to delete.")
    except Exception as e:
        logger.error(f"Error during feedback cleanup: {e}", exc_info=True)

# --- Main loop ---
message_count = 0
try:
    for message in consumer:
        try:
            feedback = message.value
            feedback_id = int(feedback.get("id"))
            user_id = int(feedback.get("user_id"))
            article_id = int(feedback.get("article_id"))
            feedback_type = feedback.get("type", "like")  # like, dislike, view, etc.
            feedback_date = feedback.get("date", datetime.now().isoformat())
            
            # Parse date to timestamp
            try:
                feedback_datetime = datetime.fromisoformat(feedback_date)
                timestamp = feedback_datetime.timestamp()
            except Exception:
                timestamp = datetime.now().timestamp()

            # Save feedback to Qdrant
            point = PointStruct(
                id=feedback_id,
                vector=[0.0],  # placeholder vector (not used for feedbacks)
                payload={
                    "user_id": user_id,
                    "article_id": article_id,
                    "type": feedback_type,
                    "date": feedback_date,
                    "timestamp": timestamp,
                    "text": feedback.get("text", ""),
                    "rating": feedback.get("rating")
                }
            )
            qdrant.upsert(collection_name=QDRANT_COLLECTION, points=[point])
            logger.info(f"✅ Feedback {feedback_id} (user: {user_id}, article: {article_id}, type: {feedback_type}) saved.")

            message_count += 1
            
            # Cleanup old feedbacks periodically (every 100 messages)
            if message_count % 100 == 0:
                cleanup_old_feedbacks()
                
        except Exception as e:
            logger.error(f"Error processing message: {e}", exc_info=True)
            
except KeyboardInterrupt:
    logger.info("Consumer shutdown requested.")
finally:
    consumer.close()
    logger.info("Kafka consumer closed.")
