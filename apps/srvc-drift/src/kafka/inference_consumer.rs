use crate::config::KafkaConfig;
use crate::error::{DriftError, Result};
use crate::kafka::models::InferenceCommand;
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
