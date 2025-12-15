mod config;
mod drift;
mod error;
mod health;
mod kafka;
mod reporting;
mod schema;
mod storage;

use crate::config::Config;
use crate::drift::calculator::drift_calculation_loop;
use crate::health::{AppState, health_handler};
use crate::kafka::feedback_consumer::{FeedbackConsumer, process_feedback};
use crate::kafka::inference_consumer::{InferenceConsumer, process_inferences};
use crate::kafka::models::{FeedbackAggregate, InferenceCommand};
use crate::reporting::discord::{DiscordReporter, reporting_loop};
use crate::storage::postgres::PostgresStorage;
use crate::storage::qdrant::QdrantStorage;
use crate::storage::redis::RedisCache;
use axum::{Router, routing::get};
use std::sync::Arc;
use tokio::signal;
use tokio::sync::mpsc;
use tower_http::trace::TraceLayer;
use tracing::info;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "srvc_drift=info,tower_http=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer().json())
        .init();

    info!("Starting srvc-drift service");

    let config = Config::load()?;
    info!("Configuration loaded successfully");

    let postgres = PostgresStorage::new(&config.database).await?;
    let redis = if let Some(redis_config) = &config.redis {
        Some(RedisCache::new(redis_config).await?)
    } else {
        info!("Redis not configured, running without cache");
        None
    };

    let discord = Arc::new(DiscordReporter::new(config.discord.clone()));

    let qdrant = QdrantStorage::new(&config.qdrant).await?;

    let inference_consumer = InferenceConsumer::new(&config.kafka)?;
    let feedback_consumer = FeedbackConsumer::new(&config.kafka)?;

    let (inference_tx, inference_rx) = mpsc::channel::<InferenceCommand>(1000);
    let (feedback_tx, feedback_rx) = mpsc::channel::<FeedbackAggregate>(1000);

    let inference_consumer_handle = tokio::spawn(async move {
        inference_consumer.start(inference_tx).await;
    });

    let feedback_consumer_handle = tokio::spawn(async move {
        feedback_consumer.start(feedback_tx).await;
    });

    let inference_postgres = postgres.clone();
    let inference_redis = redis.clone();
    let inference_qdrant = qdrant.clone();
    let inference_processor_handle = tokio::spawn(async move {
        process_inferences(
            inference_rx,
            inference_postgres,
            inference_qdrant,
            inference_redis,
        )
        .await;
    });

    let feedback_postgres = postgres.clone();
    let feedback_redis = redis.clone();
    let feedback_processor_handle = tokio::spawn(async move {
        process_feedback(feedback_rx, feedback_postgres, feedback_redis).await;
    });

    let drift_config = config.clone();
    let drift_postgres = postgres.clone();
    let drift_discord = Arc::clone(&discord);
    let drift_handle = tokio::spawn(async move {
        drift_calculation_loop(drift_config, drift_postgres, drift_discord).await;
    });

    let report_config = config.clone();
    let report_postgres = postgres.clone();
    let report_discord = Arc::clone(&discord);
    let reporting_handle = tokio::spawn(async move {
        reporting_loop(report_config, report_postgres, report_discord).await;
    });

    let app_state = Arc::new(AppState {
        postgres: postgres.clone(),
        redis: redis.clone(),
    });

    let app = Router::new()
        .route("/health", get(health_handler))
        .layer(TraceLayer::new_for_http())
        .with_state(app_state);

    let addr = format!("{}:{}", config.server.host, config.server.port);
    info!("Starting HTTP server on {}", addr);

    let listener = tokio::net::TcpListener::bind(&addr).await?;
    let server_handle = tokio::spawn(async move {
        axum::serve(listener, app).await.expect("Server failed");
    });

    info!("Service is running. Press Ctrl+C to shutdown.");
    shutdown_signal().await;

    info!("Shutdown signal received, stopping gracefully...");

    inference_consumer_handle.abort();
    feedback_consumer_handle.abort();
    inference_processor_handle.abort();
    feedback_processor_handle.abort();
    drift_handle.abort();
    reporting_handle.abort();
    server_handle.abort();

    info!("Service stopped");
    Ok(())
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }
}
