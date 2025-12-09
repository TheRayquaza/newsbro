use crate::config::DatabaseConfig;
use crate::error::Result;
use crate::kafka::models::InferenceCommand;
use chrono::{DateTime, Duration, Utc};
use sqlx::Row;
use sqlx::postgres::{PgPool, PgPoolOptions};
use tracing::{info, instrument};

#[derive(Debug, Clone)]
pub struct PostgresStorage {
    pool: PgPool,
}

#[derive(Debug, Clone)]
pub struct FeedbackMetrics {
    pub time_bucket: DateTime<Utc>,
    pub total_inferences: i32,
    pub likes: i32,
    pub dislikes: i32,
    pub like_ratio: f64,
    pub avg_response_time_ms: Option<f64>,
    pub p95_response_time_ms: Option<f64>,
    pub p99_response_time_ms: Option<f64>,
}

#[derive(Debug, Clone)]
pub struct DriftSnapshot {
    pub snapshot_time: DateTime<Utc>,
    pub embedding_centroid: Vec<f32>,
    pub psi_score: Option<f64>,
    pub kl_divergence: Option<f64>,
    pub cosine_distance_from_baseline: Option<f64>,
    pub sample_count: i32,
    pub drift_severity: String,
}

impl PostgresStorage {
    pub async fn new(config: &DatabaseConfig) -> Result<Self> {
        let pool = PgPoolOptions::new()
            .max_connections(config.max_connections)
            .min_connections(config.min_connections)
            .acquire_timeout(std::time::Duration::from_secs(
                config.acquire_timeout_seconds,
            ))
            .connect(&config.url)
            .await?;

        info!("PostgreSQL connection pool established");

        Ok(Self { pool })
    }

    #[instrument(skip(self, inference))]
    pub async fn insert_inference(&self, inference: &InferenceCommand) -> Result<()> {
        let embedding_json = serde_json::to_value(&inference.embedding)?;

        sqlx::query(
            r#"
            INSERT INTO inference_log 
                (user_id, embedding, model, score, date, metadata)
            VALUES ($1, $2, $3, $4, $5, $6)
            "#,
        )
        .bind(inference.user_id)
        .bind(&embedding_json)
        .bind(&inference.model)
        .bind(inference.date)
        .bind(inference.score)
        .bind(&inference.metadata)
        .execute(&self.pool)
        .await?;

        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn insert_feedback_metrics(&self, metrics: &FeedbackMetrics) -> Result<()> {
        sqlx::query(
            r#"
            INSERT INTO feedback_metrics 
                (time_bucket, bucket_size_seconds, total_inferences, likes, dislikes, 
                 like_ratio, avg_response_time_ms, p95_response_time_ms, p99_response_time_ms)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            "#,
        )
        .bind(metrics.time_bucket)
        .bind(3600) // 1 hour bucket
        .bind(metrics.total_inferences)
        .bind(metrics.likes)
        .bind(metrics.dislikes)
        .bind(metrics.like_ratio)
        .bind(metrics.avg_response_time_ms)
        .bind(metrics.p95_response_time_ms)
        .bind(metrics.p99_response_time_ms)
        .execute(&self.pool)
        .await?;

        Ok(())
    }

    #[instrument(skip(self, snapshot))]
    pub async fn insert_drift_snapshot(&self, snapshot: &DriftSnapshot) -> Result<()> {
        let centroid_json = serde_json::to_value(&snapshot.embedding_centroid)?;

        sqlx::query(
            r#"
            INSERT INTO drift_snapshots 
                (snapshot_time, embedding_centroid, psi_score, kl_divergence, 
                 cosine_distance_from_baseline, sample_count, drift_severity)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            "#,
        )
        .bind(snapshot.snapshot_time)
        .bind(&centroid_json)
        .bind(snapshot.psi_score)
        .bind(snapshot.kl_divergence)
        .bind(snapshot.cosine_distance_from_baseline)
        .bind(snapshot.sample_count)
        .bind(&snapshot.drift_severity)
        .execute(&self.pool)
        .await?;

        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn get_recent_metrics(&self, hours: i64) -> Result<Vec<FeedbackMetrics>> {
        let cutoff = Utc::now() - Duration::hours(hours);

        let rows = sqlx::query(
            r#"
            SELECT time_bucket, total_inferences, likes, dislikes, like_ratio,
                   avg_response_time_ms, p95_response_time_ms, p99_response_time_ms
            FROM feedback_metrics
            WHERE time_bucket >= $1
            ORDER BY time_bucket DESC
            "#,
        )
        .bind(cutoff)
        .fetch_all(&self.pool)
        .await?;

        let metrics = rows
            .into_iter()
            .map(|row| FeedbackMetrics {
                time_bucket: row.get("time_bucket"),
                total_inferences: row.get("total_inferences"),
                likes: row.get("likes"),
                dislikes: row.get("dislikes"),
                like_ratio: row.get("like_ratio"),
                avg_response_time_ms: row.get("avg_response_time_ms"),
                p95_response_time_ms: row.get("p95_response_time_ms"),
                p99_response_time_ms: row.get("p99_response_time_ms"),
            })
            .collect();

        Ok(metrics)
    }

    #[instrument(skip(self))]
    pub async fn get_embeddings_for_drift(&self, hours: i64, limit: i64) -> Result<Vec<Vec<f32>>> {
        let cutoff = Utc::now() - Duration::hours(hours);

        let rows = sqlx::query(
            r#"
            SELECT embedding
            FROM inference_log
            WHERE timestamp >= $1
            ORDER BY timestamp DESC
            LIMIT $2
            "#,
        )
        .bind(cutoff)
        .bind(limit)
        .fetch_all(&self.pool)
        .await?;

        let embeddings = rows
            .into_iter()
            .filter_map(|row| {
                let json: serde_json::Value = row.get("embedding");
                serde_json::from_value(json).ok()
            })
            .collect();

        Ok(embeddings)
    }

    #[instrument(skip(self))]
    pub async fn get_baseline_centroid(&self, hours: i64) -> Result<Option<Vec<f32>>> {
        let cutoff = Utc::now() - Duration::hours(hours);

        let row = sqlx::query(
            r#"
            SELECT embedding_centroid
            FROM drift_snapshots
            WHERE snapshot_time >= $1
            ORDER BY snapshot_time ASC
            LIMIT 1
            "#,
        )
        .bind(cutoff)
        .fetch_optional(&self.pool)
        .await?;

        if let Some(row) = row {
            let json: serde_json::Value = row.get("embedding_centroid");
            Ok(serde_json::from_value(json).ok())
        } else {
            Ok(None)
        }
    }

    pub async fn health_check(&self) -> Result<()> {
        sqlx::query("SELECT 1").fetch_one(&self.pool).await?;
        Ok(())
    }
}
