use crate::config::KafkaConfig;
use crate::error::{DriftError, Result};
use crate::kafka::models::FeedbackAggregate;
use crate::storage::postgres::PostgresStorage;
use crate::storage::redis::RedisCache;
use rdkafka::config::ClientConfig;
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::message::Message;
use tokio::sync::mpsc;
use tracing::{error, info, warn};

pub struct FeedbackConsumer {
    consumer: StreamConsumer,
}

impl FeedbackConsumer {
    pub fn new(config: &KafkaConfig) -> Result<Self> {
        let consumer: StreamConsumer = ClientConfig::new()
            .set("bootstrap.servers", &config.bootstrap_servers)
            .set("group.id", format!("{}-feedback", config.consumer_group))
            .set("enable.auto.commit", "true")
            .set("auto.offset.reset", "latest")
            .set("session.timeout.ms", "30000")
            .create()
            .map_err(DriftError::Kafka)?;

        consumer
            .subscribe(&[&config.feedback_topic])
            .map_err(DriftError::Kafka)?;

        info!(
            "Feedback consumer initialized for topic: {}",
            config.feedback_topic
        );

        Ok(Self { consumer })
    }

    pub async fn start(self, tx: mpsc::Sender<FeedbackAggregate>) {
        info!("Starting feedback consumer");

        loop {
            match self.consumer.recv().await {
                Ok(msg) => {
                    if let Some(payload) = msg.payload() {
                        match serde_json::from_slice::<FeedbackAggregate>(payload) {
                            Ok(feedback) => {
                                if tx.send(feedback).await.is_err() {
                                    error!("Failed to send feedback to channel");
                                    break;
                                }
                            }
                            Err(e) => {
                                warn!("Failed to deserialize feedback: {}", e);
                            }
                        }
                    }
                }
                Err(e) => {
                    error!("Kafka error receiving feedback message: {}", e);
                    tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
                }
            }
        }
    }
}

pub async fn process_feedback(
    mut rx: mpsc::Receiver<FeedbackAggregate>,
    _postgres: PostgresStorage,
    redis: Option<RedisCache>,
) {
    info!("Feedback processor started");

    while let Some(feedback) = rx.recv().await {
        // Update Redis counters
        if let Some(cache) = &redis
            && let Err(e) = cache.increment_feedback(&feedback.value).await
        {
            error!("Failed to increment feedback: {}", e);
        }
    }

    warn!("Feedback processor stopped");
}
