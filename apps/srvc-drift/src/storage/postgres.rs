use crate::config::DatabaseConfig;
use crate::error::{DriftError, Result};
use crate::kafka::models::InferenceCommand;
use crate::schema::{drift_snapshots, feedback_metrics, inference_log};
pub use crate::storage::models::{
    FeedbackMetrics,
    InferenceLog,
    DriftSnapshot,
};
use chrono::{Duration, Utc};
use diesel::prelude::*;
use diesel_async::pooled_connection::AsyncDieselConnectionManager;
use diesel_async::pooled_connection::deadpool::Pool;
use diesel_async::{AsyncPgConnection, RunQueryDsl};
use tracing::{info, instrument};

pub type DbPool = Pool<AsyncPgConnection>;

#[derive(Clone)]
pub struct PostgresStorage {
    pool: DbPool,
}

impl PostgresStorage {
    pub async fn new(config: &DatabaseConfig) -> Result<Self> {
        let config_manager = AsyncDieselConnectionManager::<AsyncPgConnection>::new(&config.url);

        let pool = Pool::builder(config_manager)
            .max_size(config.max_connections as usize)
            .build()
            .map_err(|e| {
                DriftError::Database(diesel::result::Error::DatabaseError(
                    diesel::result::DatabaseErrorKind::UnableToSendCommand,
                    Box::new(e.to_string()),
                ))
            })?;

        info!("PostgreSQL connection pool established");

        Ok(Self { pool })
    }

    #[instrument(skip(self, inference))]
    pub async fn insert_inference(
        &self,
        inference: &InferenceCommand,
        embedding: &Vec<f32>,
    ) -> Result<()> {
        let embedding_json = serde_json::to_value(embedding)?;

        let new_inference = InferenceLog {
            user_id: inference.user_id as i32,
            embedding: embedding_json,
            model: inference.model.clone(),
            score: inference.score as f32,
            date: inference.date,
        };

        let mut conn = self.pool.get().await.map_err(|e| {
            crate::error::DriftError::Database(diesel::result::Error::DatabaseError(
                diesel::result::DatabaseErrorKind::UnableToSendCommand,
                Box::new(e.to_string()),
            ))
        })?;

        diesel::insert_into(inference_log::table)
            .values(&new_inference)
            .execute(&mut conn)
            .await?;

        Ok(())
    }

    #[instrument(skip(self, snapshot))]
    pub async fn insert_drift_snapshot(&self, snapshot: &DriftSnapshot) -> Result<()> {

        let new_snapshot = DriftSnapshot {
            snapshot_time: snapshot.snapshot_time,
            embedding_centroid: snapshot.embedding_centroid.clone(),
            psi_score: snapshot.psi_score,
            kl_divergence: snapshot.kl_divergence,
            cosine_distance_from_baseline: snapshot.cosine_distance_from_baseline,
            sample_count: snapshot.sample_count,
            drift_severity: snapshot.drift_severity.clone(),
        };

        let mut conn = self.pool.get().await.map_err(|e| {
            crate::error::DriftError::Database(diesel::result::Error::DatabaseError(
                diesel::result::DatabaseErrorKind::UnableToSendCommand,
                Box::new(e.to_string()),
            ))
        })?;

        diesel::insert_into(drift_snapshots::table)
            .values(&new_snapshot)
            .execute(&mut conn)
            .await?;

        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn get_recent_metrics(&self, hours: i64) -> Result<Vec<FeedbackMetrics>> {
        let cutoff = Utc::now() - Duration::hours(hours);

        let mut conn = self.pool.get().await.map_err(|e| {
            crate::error::DriftError::Database(diesel::result::Error::DatabaseError(
                diesel::result::DatabaseErrorKind::UnableToSendCommand,
                Box::new(e.to_string()),
            ))
        })?;

        let results: Vec<FeedbackMetrics> = feedback_metrics::table
            .filter(feedback_metrics::time_bucket.ge(cutoff))
            .order(feedback_metrics::time_bucket.desc())
            .select(FeedbackMetrics::as_select())
            .load::<FeedbackMetrics>(&mut conn)
            .await?;

        let metrics = results
            .into_iter()
            .map(|m| FeedbackMetrics {
                time_bucket: m.time_bucket,
                total_inferences: m.total_inferences,
                bucket_size_seconds: m.bucket_size_seconds,
                likes: m.likes,
                dislikes: m.dislikes,
                like_ratio: m.like_ratio,
            })
            .collect();

        Ok(metrics)
    }

    #[instrument(skip(self))]
    pub async fn get_embeddings_for_drift(&self, hours: i64, limit: i64) -> Result<Vec<Vec<f32>>> {
        let cutoff = Utc::now() - Duration::hours(hours);

        let mut conn = self.pool.get().await.map_err(|e| {
            crate::error::DriftError::Database(diesel::result::Error::DatabaseError(
                diesel::result::DatabaseErrorKind::UnableToSendCommand,
                Box::new(e.to_string()),
            ))
        })?;

        let results = inference_log::table
            .filter(inference_log::date.ge(cutoff))
            .order(inference_log::date.desc())
            .limit(limit)
            .select(inference_log::embedding)
            .load::<serde_json::Value>(&mut conn)
            .await?;

        let embeddings: Vec<Vec<f32>> = results
            .into_iter()
            .filter_map(|json| serde_json::from_value(json).ok())
            .collect();

        Ok(embeddings)
    }

    pub async fn health_check(&self) -> Result<()> {
        let mut conn = self.pool.get().await.map_err(|e| {
            crate::error::DriftError::Database(diesel::result::Error::DatabaseError(
                diesel::result::DatabaseErrorKind::UnableToSendCommand,
                Box::new(e.to_string()),
            ))
        })?;

        diesel::sql_query("SELECT 1").execute(&mut conn).await?;

        Ok(())
    }
}
