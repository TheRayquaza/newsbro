# Database Configuration
export DRIFT__DATABASE__URL=postgresql://username:password@localhost:5432/srvc_drift
export DRIFT__DATABASE__MAX_CONNECTIONS=20
export DRIFT__DATABASE__MIN_CONNECTIONS=5
export DRIFT__DATABASE__ACQUIRE_TIMEOUT_SECONDS=30

# Kafka Configuration
export DRIFT__KAFKA__BOOTSTRAP_SERVERS=localhost:9092
export DRIFT__KAFKA__INFERENCE_TOPIC=inference-command
export DRIFT__KAFKA__FEEDBACK_TOPIC=feedback-aggregate
export DRIFT__KAFKA__CONSUMER_GROUP=srvc-drift-group
export DRIFT__KAFKA__MAX_BATCH_SIZE=1000

# Redis Configuration (optional)
export DRIFT__REDIS__URL=redis://localhost:6379
export DRIFT__REDIS__POOL_SIZE=10

# Discord Configuration
export DRIFT__DISCORD__WEBHOOK_URL="https://discord.com/api/webhooks/1448348111398244584/gEzO5jFz5bYoHnDWgiLt4uc889JcV6ub1YKVZbda69aRwJquZHmohVcGvk4ZT9u_PdW9"
export DRIFT__DISCORD__USERNAME="Drift Monitor"
export DRIFT__DISCORD__ALERT_THRESHOLDS__DRIFT_WARNING=0.1
export DRIFT__DISCORD__ALERT_THRESHOLDS__DRIFT_CRITICAL=0.3
export DRIFT__DISCORD__ALERT_THRESHOLDS__DISLIKE_RATIO_WARNING=0.3

# Qdrant Configuration
export DRIFT__QDRANT__URL="http://localhost:6334"
export DRIFT__QDRANT__COLLECTION_NAME=articles
export DRIFT__QDRANT__API_KEY=""

# Drift Detection Configuration
export DRIFT__DRIFT__CALCULATION_INTERVAL_SECONDS=300
export DRIFT__DRIFT__LOOKBACK_WINDOW_HOURS=24
export DRIFT__DRIFT__BASELINE_WINDOW_HOURS=168
export DRIFT__DRIFT__MIN_SAMPLES_FOR_DRIFT=100
export DRIFT__DRIFT__EMBEDDING_DIMENSION=1536

# Server Configuration
export DRIFT__SERVER__HOST=0.0.0.0
export DRIFT__SERVER__PORT=8080

# Reporting Configuration
export DRIFT__REPORTING__INTERVAL_SECONDS=3600
export DRIFT__REPORTING__INCLUDE_EMBEDDINGS_SAMPLE=false



# Logging Configuration
export RUST_LOG=info,srvc_drift=debug