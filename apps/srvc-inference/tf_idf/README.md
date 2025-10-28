# TF-IDF

## Dev

* linting:

```bash
mypy .
ruff check . --fix
bandit -r .
black .
```

* run:

```bash
cd apps/
docker compose up -d
cd srvc-inference
nix develop .

# Using local pickle
python3 -m tf_idf.main tf_idf/.env ~/Downloads/model.pkl
# Using artifact in mlflow
# can require ssh connection to mlflow
python3 -m tf_idf.main tf_idf/.env
```

* train: `python3 -m tf_idf.src.model.train articles.csv`

* test: `python3 -m tf_idf.src.model.test <path.pkl>`

## Deployment

| Environment Variable            | Default Value           | Description                                             |
| ------------------------------- | ----------------------- | ------------------------------------------------------- |
| `ENVIRONMENT`                   | `development`           | Application environment (development, production, etc.) |
| `API_PORT`                      | `8000`                  | Port for API server                                     |
| `API_HOST`                      | `0.0.0.0`               | Host for API server                                     |
| `LOG_LEVEL`                     | `INFO`                  | Logging level                                           |
| `KAFKA_BOOTSTRAP_SERVERS`       | `localhost:9092`        | Kafka bootstrap servers                                 |
| `KAFKA_PRODUCER_TOPIC`          | `inference-topic`       | Kafka topic for producer messages                       |
| `KAFKA_ARTICLE_CONSUMER_TOPIC`  | `article-aggregate`     | Kafka topic for article consumer                        |
| `KAFKA_ARTICLE_CONSUMER_GROUP`  | `article-tfidf-group`   | Kafka consumer group for articles                       |
| `KAFKA_FEEDBACK_CONSUMER_TOPIC` | `feedback-aggregate`    | Kafka topic for feedback consumer                       |
| `KAFKA_FEEDBACK_CONSUMER_GROUP` | `feedback-tfidf-group`  | Kafka consumer group for feedback                       |
| `KAFKA_BATCH_SIZE`              | `64`                    | Kafka batch size (before triggering ingestion)          |
| `KAFKA_BATCH_INTERVAL`          | `10`                    | Seconds before hold values gets release and processed   |
| `QDRANT_URL`                    | `http://localhost:6333` | URL for Qdrant vector database                          |
| `QDRANT_ARTICLES_COLLECTION`    | `articles` | articles collection name in Qdrant                          |
| `MODEL_URI`                     | *(None)*                | URI or path to the ML model                             |
| `MODEL_NAME`                     | `tfidf`                |  Model name used for recommendation                             |
| `FEEDBACK_RETENTION_DAYS`     | `30` | Number of days before redis drop user vector |
| `TOP_K_ARTICLES`         | `10`                    | Number of top recommendations to return                 |
| `MLFLOW_S3_ENDPOINT_URL`        | `http://minio-svc.localhost:8080` | S3 Endpoint                                   |
| `AWS_ACCESS_KEY_ID`             | `minio`                 | S3 access key                                           |
| `AWS_SECRET_ACCESS_KEY`         | `minio123`              | S3 secret key                                           |
| `MLFLOW_TRACKING_URI`           | `http://mlflow.localhost:8080` | Mlflow tracking URI                              |
| `REDIS_SENTINELS`               | `localhost:26379,localhost:26380,localhost:26381` | Redis Sentinel endpoints (comma-separated)    |
| `REDIS_MASTER_NAME`             | `mymaster`              | Redis Sentinel master name                    |
| `REDIS_PASSWORD`                | `None`                  | Redis password (optional)                     |
| `REDIS_DB`                      | `0`                     | Redis database number                         |
