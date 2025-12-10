use crate::schema::{drift_snapshots, feedback_metrics, feedback_records, inference_log};
use chrono::{DateTime, Utc};
use diesel::prelude::*;

// Insertable models
#[derive(Debug, Clone, Queryable, Insertable, Selectable)]
#[diesel(table_name = inference_log)]
pub struct InferenceLog {
    pub user_id: i32,
    pub embedding: serde_json::Value,
    pub model: String,
    pub score: f32,
    pub date: DateTime<Utc>,
}

#[derive(Debug, Clone, Insertable, Queryable, Selectable)]
#[diesel(table_name = feedback_metrics)]
pub struct FeedbackMetrics {
    pub time_bucket: chrono::NaiveDateTime,
    pub bucket_size_seconds: i32,
    pub total_inferences: i32,
    pub likes: i32,
    pub dislikes: i32,
    pub like_ratio: f64,
}

#[derive(Debug, Clone, Insertable, Queryable, Selectable)]
#[diesel(table_name = drift_snapshots)]
pub struct DriftSnapshot {
    pub snapshot_time: DateTime<Utc>,
    pub embedding_centroid: serde_json::Value,
    pub psi_score: Option<f64>,
    pub kl_divergence: Option<f64>,
    pub cosine_distance_from_baseline: Option<f64>,
    pub sample_count: i32,
    pub drift_severity: String,
}

#[derive(Debug, Clone, Insertable, Queryable, Selectable)]
#[diesel(table_name = feedback_records)]
pub struct FeedbackRecord {
    pub user_id: i32,
    pub feedback_type: String,
    pub timestamp: DateTime<Utc>,
}
