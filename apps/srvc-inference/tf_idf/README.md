# TF-IDF

## Dev

* linting:

```
cd ..
nix develop .
mypy .
ruff check . --fix
bandit -r .
black .
```

* run:

```
cd ..
docker compose up -d
nix develop .
python3 -m tf_idf.main tf_idf/.env # ~/Downloads/model.pkl (this is optional)
```

* train:

```
python3 -m tf_idf.src.model.train articles.csv
```

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
| `QDRANT_URL`                    | `http://localhost:6333` | URL for Qdrant vector database                          |
| `MODEL_URI`                     | *(None)*                | URI or path to the ML model                             |
| `TOP_N_RECOMMENDATIONS`         | `10`                    | Number of top recommendations to return                 |
| `MLFLOW_S3_ENDPOINT_URL`        | `http://minio-svc.localhost:8080` | S3 Endpoint                                   |
| `AWS_ACCESS_KEY_ID`             | `minio`                 | S3 access key                                           |
| `AWS_SECRET_ACCESS_KEY`         | `minio123`              | S3 secret key                                           |
| `MLFLOW_TRACKING_URI`           | `http://mlflow.localhost:8080` | Mlflow tracking URI                              |
