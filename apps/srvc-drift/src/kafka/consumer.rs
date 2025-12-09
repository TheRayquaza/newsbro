use crate::config::KafkaConfig;
use crate::error::{DriftError, Result};
use crate::kafka::models::{FeedbackAggregate, InferenceCommand};
use rdkafka::config::ClientConfig;
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::message::Message;
use tokio::sync::mpsc;
use tracing::{error, info, warn};

pub struct KafkaConsumerManager {
    inference_consumer: StreamConsumer,
    feedback_consumer: StreamConsumer,
}

impl KafkaConsumerManager {
    pub fn new(config: &KafkaConfig) -> Result<Self> {
        let inference_consumer: StreamConsumer = ClientConfig::new()
            .set("bootstrap.servers", &config.bootstrap_servers)
            .set("group.id", &config.consumer_group)
            .set("enable.auto.commit", "true")
            .set("auto.offset.reset", "latest")
            .set("session.timeout.ms", "30000")
            .create()
            .map_err(DriftError::Kafka)?;

        let feedback_consumer: StreamConsumer = ClientConfig::new()
            .set("bootstrap.servers", &config.bootstrap_servers)
            .set("group.id", &format!("{}-feedback", config.consumer_group))
            .set("enable.auto.commit", "true")
            .set("auto.offset.reset", "latest")
            .set("session.timeout.ms", "30000")
            .create()
            .map_err(DriftError::Kafka)?;

        inference_consumer
            .subscribe(&[&config.inference_topic])
            .map_err(DriftError::Kafka)?;

        feedback_consumer
            .subscribe(&[&config.feedback_topic])
            .map_err(DriftError::Kafka)?;

        info!(
            "Kafka consumers initialized for topics: {}, {}",
            config.inference_topic, config.feedback_topic
        );

        Ok(Self {
            inference_consumer,
            feedback_consumer,
        })
    }

    pub async fn start_inference_consumer(
        self,
        tx: mpsc::Sender<InferenceCommand>,
        expected_dim: usize,
    ) {
        info!("Starting inference consumer");

        loop {
            match self.inference_consumer.recv().await {
                Ok(msg) => {
                    if let Some(payload) = msg.payload() {
                        match serde_json::from_slice::<InferenceCommand>(payload) {
                            Ok(command) => {
                                if let Err(e) = command.validate_embedding_dimension(expected_dim) {
                                    warn!("Invalid embedding dimension: {}", e);
                                    continue;
                                }

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

    pub async fn start_feedback_consumer(
        feedback_consumer: StreamConsumer,
        tx: mpsc::Sender<FeedbackAggregate>,
    ) {
        info!("Starting feedback consumer");

        loop {
            match feedback_consumer.recv().await {
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

    pub fn split(self) -> (StreamConsumer, StreamConsumer) {
        (self.inference_consumer, self.feedback_consumer)
    }
}
