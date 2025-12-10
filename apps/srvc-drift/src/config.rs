use anyhow::{Context, Result};
use serde::Deserialize;
use std::time::Duration;

#[derive(Debug, Deserialize, Clone)]
pub struct Config {
    pub kafka: KafkaConfig,
    pub database: DatabaseConfig,
    pub redis: Option<RedisConfig>,
    pub discord: DiscordConfig,
    pub qdrant: QdrantConfig,
    pub drift: DriftConfig,
    pub server: ServerConfig,
    pub reporting: ReportingConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct KafkaConfig {
    pub bootstrap_servers: String,
    pub inference_topic: String,
    pub feedback_topic: String,
    pub consumer_group: String,
    #[serde(default = "default_batch_size")]
    pub max_batch_size: usize,
}

#[derive(Debug, Deserialize, Clone)]
pub struct DatabaseConfig {
    pub url: String,
    #[serde(default = "default_max_connections")]
    pub max_connections: u32,
    #[serde(default = "default_min_connections")]
    pub min_connections: u32,
    #[serde(default = "default_acquire_timeout")]
    pub acquire_timeout_seconds: u64,
}

#[derive(Debug, Deserialize, Clone)]
pub struct RedisConfig {
    pub url: String,
    #[serde(default = "default_redis_pool_size")]
    pub pool_size: u32,
}

#[derive(Debug, Deserialize, Clone)]
pub struct DiscordConfig {
    pub webhook_url: String,
    #[serde(default = "default_discord_username")]
    pub username: String,
    pub alert_thresholds: AlertThresholds,
}

#[derive(Debug, Deserialize, Clone)]
pub struct QdrantConfig {
    pub url: String,
    pub collection_name: String,
    pub api_key: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct AlertThresholds {
    #[serde(default = "default_drift_warning")]
    pub drift_warning: f64,
    #[serde(default = "default_drift_critical")]
    pub drift_critical: f64,
    #[serde(default = "default_dislike_ratio_warning")]
    pub dislike_ratio_warning: f64,
}

#[derive(Debug, Deserialize, Clone)]
pub struct DriftConfig {
    #[serde(default = "default_calculation_interval")]
    pub calculation_interval_seconds: u64,
    #[serde(default = "default_lookback_window")]
    pub lookback_window_hours: u64,
    #[serde(default = "default_baseline_window")]
    pub baseline_window_hours: u64,
    #[serde(default = "default_min_samples")]
    pub min_samples_for_drift: usize,
    #[serde(default = "default_embedding_dim")]
    pub embedding_dimension: usize,
}

#[derive(Debug, Deserialize, Clone)]
pub struct ServerConfig {
    #[serde(default = "default_host")]
    pub host: String,
    #[serde(default = "default_port")]
    pub port: u16,
}

#[derive(Debug, Deserialize, Clone)]
pub struct ReportingConfig {
    pub interval_seconds: u64,
    #[serde(default)]
    pub include_embeddings_sample: bool,
}

// Default value functions
fn default_batch_size() -> usize {
    1000
}
fn default_max_connections() -> u32 {
    20
}
fn default_min_connections() -> u32 {
    5
}
fn default_acquire_timeout() -> u64 {
    30
}
fn default_redis_pool_size() -> u32 {
    10
}
fn default_discord_username() -> String {
    "Drift Monitor".to_string()
}
fn default_drift_warning() -> f64 {
    0.1
}
fn default_drift_critical() -> f64 {
    0.3
}
fn default_dislike_ratio_warning() -> f64 {
    0.3
}
fn default_calculation_interval() -> u64 {
    300
}
fn default_lookback_window() -> u64 {
    24
}
fn default_baseline_window() -> u64 {
    168
}
fn default_min_samples() -> usize {
    100
}
fn default_embedding_dim() -> usize {
    1536
}
fn default_host() -> String {
    "0.0.0.0".to_string()
}
fn default_port() -> u16 {
    8080
}

impl Config {
    pub fn load() -> Result<Self> {
        dotenvy::dotenv().ok();

        let config = config::Config::builder()
            .add_source(config::File::with_name(".env").required(false))
            .add_source(config::Environment::with_prefix("DRIFT").separator("__"))
            .build()
            .context("Failed to build configuration")?;

        let cfg: Config = config
            .try_deserialize()
            .context("Failed to deserialize configuration")?;

        cfg.validate()?;

        Ok(cfg)
    }

    fn validate(&self) -> Result<()> {
        anyhow::ensure!(
            !self.kafka.bootstrap_servers.is_empty(),
            "Kafka bootstrap servers cannot be empty"
        );
        anyhow::ensure!(
            !self.discord.webhook_url.is_empty(),
            "Discord webhook URL cannot be empty"
        );
        anyhow::ensure!(
            self.reporting.interval_seconds > 0,
            "Reporting interval must be greater than 0"
        );
        Ok(())
    }

    pub fn reporting_interval(&self) -> Duration {
        Duration::from_secs(self.reporting.interval_seconds)
    }

    pub fn drift_calculation_interval(&self) -> Duration {
        Duration::from_secs(self.drift.calculation_interval_seconds)
    }
}
