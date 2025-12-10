use thiserror::Error;

#[derive(Error, Debug)]
pub enum DriftError {
    #[error("Kafka error: {0}")]
    Kafka(#[from] rdkafka::error::KafkaError),

    #[error("Database error: {0}")]
    Database(#[from] diesel::result::Error),

    #[error("Redis error: {0}")]
    Redis(#[from] redis::RedisError),

    #[error("Serialization error: {0}")]
    Serialization(#[from] serde_json::Error),

    #[error("HTTP error: {0}")]
    Http(#[from] reqwest::Error),

    #[error("Configuration error: {0}")]
    Config(String),

    #[error("Insufficient data: {0}")]
    InsufficientData(String),

    #[error("Invalid embedding dimension: expected {expected}, got {actual}")]
    InvalidEmbeddingDimension { expected: usize, actual: usize },

    #[error("Service unavailable: {0}")]
    ServiceUnavailable(String),

    #[error("Error: {0}")]
    InvalidInput(String),

    #[error("Embedding not found: {0}")]
    EmbeddingNotFound(String),
}

pub type Result<T> = std::result::Result<T, DriftError>;
