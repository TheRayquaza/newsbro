use crate::config::KafkaConfig;
use crate::error::{DriftError, Result};
use crate::kafka::models::InferenceCommand;
use crate::storage::postgres::PostgresStorage;
use crate::storage::qdrant::QdrantStorage;
use crate::storage::redis::RedisCache;
use rdkafka::config::ClientConfig;
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::message::Message;
use tokio::sync::mpsc;
use tracing::{error, info, warn};

pub struct InferenceConsumer {
    consumer: StreamConsumer,
}

impl InferenceConsumer {
    pub fn new(config: &KafkaConfig) -> Result<Self> {
        let consumer: StreamConsumer = ClientConfig::new()
            .set("bootstrap.servers", &config.bootstrap_servers)
            .set("group.id", &config.consumer_group)
            .set("enable.auto.commit", "true")
            .set("auto.offset.reset", "latest")
            .set("session.timeout.ms", "30000")
            .create()
            .map_err(DriftError::Kafka)?;

        consumer
            .subscribe(&[&config.inference_topic])
            .map_err(DriftError::Kafka)?;

        info!(
            "Inference consumer initialized for topic: {}",
            config.inference_topic
        );

        Ok(Self { consumer })
    }

    pub async fn start(self, tx: mpsc::Sender<InferenceCommand>) {
        info!("Starting inference consumer");

        loop {
            match self.consumer.recv().await {
                Ok(msg) => {
                    if let Some(payload) = msg.payload() {
                        match serde_json::from_slice::<InferenceCommand>(payload) {
                            Ok(command) => {
                                if tx.send(command).await.is_err() {
                                    error!("Failed to send inference command to channel");
                                    break;
                                }
                            }
                            Err(e) => {
                                warn!("Failed to deserialize inference command: {}", e);
                            }
                        }
                    }
                }
                Err(e) => {
                    error!("Kafka error receiving inference message: {}", e);
                    tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
                }
            }
        }
    }
}

pub async fn process_inferences(
    mut rx: mpsc::Receiver<InferenceCommand>,
    postgres: PostgresStorage,
    qdrant: QdrantStorage,
    redis: Option<RedisCache>,
) {
    info!("Inference processor started");

    while let Some(inference) = rx.recv().await {
        let article = inference.article.as_ref().map(|article| {
            info!(
                "Processing inference for user_id={}, article_id={}",
                inference.user_id, article.id
            );
            article.clone()
        });
        if article.is_none() {
            warn!(
                "Inference missing article data for user_id={}, skipping",
                inference.user_id
            );
            continue;
        }
        let article = article.unwrap();
        let collection_name = format!("{}_{}", qdrant.config.collection_prefix, inference.model);
        // Retrieve from qdrant
        let embedding = match qdrant.retrieve_embedding_by_id(collection_name, article.id).await {
            Ok(emb) => emb,
            Err(e) => {
                error!("Failed to retrieve embedding from Qdrant: {}", e);
                continue;
            }
        }
        .unwrap_or_else(|| {
            warn!("No embedding found in Qdrant for article_id={}", article.id);
            vec![]
        });
        if let Err(e) = postgres.insert_inference(&inference, &embedding).await {
            error!("Failed to insert inference: {}", e);
        }

        // Update Redis counters
        if let Some(cache) = &redis
            && let Err(e) = cache.increment_inference_count().await
        {
            error!("Failed to increment inference count: {}", e);
        }
    }

    warn!("Inference processor stopped");
}
