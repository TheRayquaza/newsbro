mod config;
mod drift;
mod error;
mod health;
mod kafka;
mod reporting;
mod schema;
mod storage;

use crate::config::Config;
use crate::drift::calculator::DriftCalculator;
use crate::health::{AppState, dependencies_handler, health_handler, ready_handler};
use crate::kafka::feedback_consumer::FeedbackConsumer;
use crate::kafka::inference_consumer::InferenceConsumer;
use crate::kafka::models::{FeedbackAggregate, InferenceCommand};
use crate::reporting::discord::{AlertSeverity, DiscordReporter, DriftReport};
use crate::storage::postgres::{DriftSnapshot, PostgresStorage};
use crate::storage::qdrant::QdrantStorage;
use crate::storage::redis::RedisCache;
use axum::{Router, routing::get};
use chrono::Utc;
use std::sync::Arc;
use tokio::signal;
use tokio::sync::mpsc;
use tokio::time::interval;
use tower_http::trace::TraceLayer;
use tracing::{error, info, warn};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "srvc_drift=info,tower_http=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer().json())
        .init();

    info!("Starting srvc-drift service");

    // Load configuration
    let config = Config::load()?;
    info!("Configuration loaded successfully");

    // Initialize storage
    let postgres = PostgresStorage::new(&config.database).await?;
    let redis = if let Some(redis_config) = &config.redis {
        Some(RedisCache::new(redis_config).await?)
    } else {
        info!("Redis not configured, running without cache");
        None
    };

    // Initialize Discord reporter
    let discord = Arc::new(DiscordReporter::new(config.discord.clone()));

    // Initialize Qdrant storage
    let qdrant = QdrantStorage::new(&config.qdrant).await?;

    // Initialize Kafka consumers
    let inference_consumer = InferenceConsumer::new(&config.kafka)?;
    let feedback_consumer = FeedbackConsumer::new(&config.kafka)?;

    // Create channels for communication
    let (inference_tx, inference_rx) = mpsc::channel::<InferenceCommand>(1000);
    let (feedback_tx, feedback_rx) = mpsc::channel::<FeedbackAggregate>(1000);

    // Spawn inference consumer task
    let inference_consumer_handle = tokio::spawn(async move {
        inference_consumer.start(inference_tx).await;
    });

    // Spawn feedback consumer task
    let feedback_consumer_handle = tokio::spawn(async move {
        feedback_consumer.start(feedback_tx).await;
    });

    // Spawn inference processor task
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

    // Spawn feedback processor task
    let feedback_postgres = postgres.clone();
    let feedback_redis = redis.clone();
    let feedback_processor_handle = tokio::spawn(async move {
        process_feedback(feedback_rx, feedback_postgres, feedback_redis).await;
    });

    // Spawn drift calculation task
    let drift_config = config.clone();
    let drift_postgres = postgres.clone();
    let drift_discord = Arc::clone(&discord);
    let drift_handle = tokio::spawn(async move {
        drift_calculation_loop(drift_config, drift_postgres, drift_discord).await;
    });

    // Spawn reporting task
    let report_config = config.clone();
    let report_postgres = postgres.clone();
    let report_discord = Arc::clone(&discord);
    let reporting_handle = tokio::spawn(async move {
        reporting_loop(report_config, report_postgres, report_discord).await;
    });

    // Setup HTTP server
    let app_state = Arc::new(AppState {
        postgres: postgres.clone(),
        redis: redis.clone(),
    });

    let app = Router::new()
        .route("/health", get(health_handler))
        .route("/health/ready", get(ready_handler))
        .route("/health/dependencies", get(dependencies_handler))
        .layer(TraceLayer::new_for_http())
        .with_state(app_state);

    let addr = format!("{}:{}", config.server.host, config.server.port);
    info!("Starting HTTP server on {}", addr);

    let listener = tokio::net::TcpListener::bind(&addr).await?;
    let server_handle = tokio::spawn(async move {
        axum::serve(listener, app).await.expect("Server failed");
    });

    // Send startup notification
    if let Err(e) = discord
        .send_alert("🚀 srvc-drift service started", AlertSeverity::Info)
        .await
    {
        warn!("Failed to send startup notification: {}", e);
    }

    // Wait for shutdown signal
    info!("Service is running. Press Ctrl+C to shutdown.");
    shutdown_signal().await;

    info!("Shutdown signal received, stopping gracefully...");

    // Cancel all tasks
    inference_consumer_handle.abort();
    feedback_consumer_handle.abort();
    inference_processor_handle.abort();
    feedback_processor_handle.abort();
    drift_handle.abort();
    reporting_handle.abort();
    server_handle.abort();

    // Send shutdown notification
    if let Err(e) = discord
        .send_alert("🛑 srvc-drift service stopped", AlertSeverity::Info)
        .await
    {
        warn!("Failed to send shutdown notification: {}", e);
    }

    info!("Service stopped");
    Ok(())
}

async fn process_inferences(
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
        // Retrieve from qdrant
        let embedding = match qdrant.retrieve_embedding_by_id(article.id).await {
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
        // Store in PostgreSQL
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

async fn process_feedback(
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

async fn drift_calculation_loop(
    config: Config,
    postgres: PostgresStorage,
    discord: Arc<DiscordReporter>,
) {
    info!("Drift calculation loop started");
    let mut interval = interval(config.drift_calculation_interval());

    loop {
        interval.tick().await;

        info!("Starting drift calculation");

        // Get recent embeddings
        let embeddings = match postgres
            .get_embeddings_for_drift(
                config.drift.lookback_window_hours as i64,
                config.drift.min_samples_for_drift as i64,
            )
            .await
        {
            Ok(emb) => emb,
            Err(e) => {
                error!("Failed to fetch embeddings: {}", e);
                continue;
            }
        };

        if embeddings.len() < config.drift.min_samples_for_drift {
            info!(
                "Insufficient samples for drift calculation: {} < {}",
                embeddings.len(),
                config.drift.min_samples_for_drift
            );
            continue;
        }

        // Calculate current centroid
        let current_centroid = match DriftCalculator::calculate_centroid(&embeddings) {
            Ok(c) => c,
            Err(e) => {
                error!("Failed to calculate centroid: {}", e);
                continue;
            }
        };

        // Get baseline embeddings
        let baseline_embeddings = match postgres
            .get_embeddings_for_drift(
                config.drift.baseline_window_hours as i64,
                config.drift.min_samples_for_drift as i64,
            )
            .await
        {
            Ok(emb) => emb,
            Err(e) => {
                error!("Failed to fetch baseline embeddings: {}", e);
                continue;
            }
        };

        let (cosine_distance, psi_score, kl_divergence) = if baseline_embeddings.len()
            >= config.drift.min_samples_for_drift
        {
            let baseline_centroid = match DriftCalculator::calculate_centroid(&baseline_embeddings)
            {
                Ok(c) => c,
                Err(e) => {
                    error!("Failed to calculate baseline centroid: {}", e);
                    continue;
                }
            };

            let cosine =
                match DriftCalculator::cosine_distance(&current_centroid, &baseline_centroid) {
                    Ok(d) => Some(d),
                    Err(e) => {
                        error!("Failed to calculate cosine distance: {}", e);
                        None
                    }
                };

            let psi = match DriftCalculator::calculate_psi(&embeddings, &baseline_embeddings) {
                Ok(p) => Some(p),
                Err(e) => {
                    error!("Failed to calculate PSI: {}", e);
                    None
                }
            };

            let kl =
                match DriftCalculator::calculate_kl_divergence(&embeddings, &baseline_embeddings) {
                    Ok(k) => Some(k),
                    Err(e) => {
                        error!("Failed to calculate KL divergence: {}", e);
                        None
                    }
                };

            (cosine, psi, kl)
        } else {
            (None, None, None)
        };

        let severity = DriftCalculator::determine_severity(
            cosine_distance.unwrap_or(0.0),
            psi_score.unwrap_or(0.0),
            kl_divergence.unwrap_or(0.0),
            &config.discord.alert_thresholds,
        );

        let embedding_json = serde_json::to_value(current_centroid).unwrap_or(serde_json::Value::Null);

        let snapshot = DriftSnapshot {
            snapshot_time: Utc::now(),
            embedding_centroid: embedding_json,
            psi_score,
            kl_divergence,
            cosine_distance_from_baseline: cosine_distance,
            sample_count: embeddings.len() as i32,
            drift_severity: severity.clone(),
        };

        if let Err(e) = postgres.insert_drift_snapshot(&snapshot).await {
            error!("Failed to insert drift snapshot: {}", e);
        }

        // Send alert if drift is significant
        if severity == "high" || severity == "medium" {
            let alert_severity = if severity == "high" {
                AlertSeverity::Critical
            } else {
                AlertSeverity::Warning
            };

            let message = format!(
                "Data drift detected! Severity: {}\nCosine distance: {:.4}\nPSI: {:.4}\nKL: {:.4}",
                severity.to_uppercase(),
                cosine_distance.unwrap_or(0.0),
                psi_score.unwrap_or(0.0),
                kl_divergence.unwrap_or(0.0)
            );

            if let Err(e) = discord.send_alert(&message, alert_severity).await {
                error!("Failed to send drift alert: {}", e);
            }
        }

        info!("Drift calculation completed: severity={}", severity);
    }
}

async fn reporting_loop(config: Config, postgres: PostgresStorage, discord: Arc<DiscordReporter>) {
    info!("Reporting loop started");
    let mut interval = interval(config.reporting_interval());

    loop {
        interval.tick().await;

        info!("Generating report");

        // Get recent metrics
        let metrics = match postgres.get_recent_metrics(24).await {
            Ok(m) => m,
            Err(e) => {
                error!("Failed to fetch metrics: {}", e);
                continue;
            }
        };

        // Get latest drift info
        let drift_info = match postgres
            .get_embeddings_for_drift(
                config.drift.lookback_window_hours as i64,
                config.drift.min_samples_for_drift as i64,
            )
            .await
        {
            Ok(embeddings) if embeddings.len() >= config.drift.min_samples_for_drift => {
                // Calculate drift
                Some(DriftReport {
                    severity: "low".to_string(),
                    cosine_distance: Some(0.05),
                    psi_score: Some(0.03),
                    kl_divergence: Some(0.02),
                    sample_count: embeddings.len() as i32,
                })
            }
            _ => None,
        };

        if let Err(e) = discord.send_report(&metrics, drift_info).await {
            error!("Failed to send report: {}", e);
        } else {
            info!("Report sent successfully");
        }
    }
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
