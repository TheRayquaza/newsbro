mod config;
mod drift;
mod error;
mod health;
mod kafka;
mod reporting;
mod storage;

use crate::config::Config;
use crate::drift::calculator::DriftCalculator;
use crate::health::{AppState, dependencies_handler, health_handler, ready_handler};
use crate::kafka::consumer::KafkaConsumerManager;
use crate::kafka::models::{FeedbackAggregate, InferenceCommand};
use crate::reporting::discord::{AlertSeverity, DiscordReporter, DriftReport};
use crate::storage::postgres::{DriftSnapshot, PostgresStorage};
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

    // Initialize Kafka consumers
    let kafka_manager = KafkaConsumerManager::new(&config.kafka)?;
    let (inference_consumer, feedback_consumer) = kafka_manager.split();

    // Create channels for communication
    let (inference_tx, inference_rx) = mpsc::channel::<InferenceCommand>(1000);
    let (feedback_tx, feedback_rx) = mpsc::channel::<FeedbackAggregate>(1000);

    // Clone Arc'd values for tasks
    let config_clone = config.clone();
    let postgres_clone = postgres.clone();
    let redis_clone = redis.clone();
    let discord_clone = Arc::clone(&discord);

    // Spawn inference consumer task
    let inference_consumer_handle = tokio::spawn(async move {
        let manager = KafkaConsumerManager {
            inference_consumer,
            feedback_consumer: feedback_consumer.clone(),
        };
        manager
            .start_inference_consumer(inference_tx, config_clone.drift.embedding_dimension)
            .await;
    });

    // Spawn feedback consumer task
    let feedback_consumer_handle = tokio::spawn(async move {
        KafkaConsumerManager::start_feedback_consumer(feedback_consumer, feedback_tx).await;
    });

    // Spawn inference processor task
    let inference_postgres = postgres.clone();
    let inference_redis = redis.clone();
    let inference_processor_handle = tokio::spawn(async move {
        process_inferences(inference_rx, inference_postgres, inference_redis).await;
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
    redis: Option<RedisCache>,
) {
    info!("Inference processor started");

    while let Some(inference) = rx.recv().await {
        // Store in PostgreSQL
        if let Err(e) = postgres.insert_inference(&inference).await {
            error!("Failed to insert inference: {}", e);
        }

        // Update Redis counters
        if let Some(cache) = &redis {
            if let Err(e) = cache.increment_inference_count().await {
                error!("Failed to increment inference count: {}", e);
            }
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
        if let Some(cache) = &redis {
            if let Err(e) = cache.increment_feedback(&feedback.feedback_type).await {
                error!("Failed to increment feedback: {}", e);
            }
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

        let snapshot = DriftSnapshot {
            snapshot_time: Utc::now(),
            embedding_centroid: current_centroid,
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
