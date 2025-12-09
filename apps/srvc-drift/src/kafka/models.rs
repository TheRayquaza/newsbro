use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InferenceCommand {
    pub user_id: u32,
    pub model: String,
    pub score: f64,
    pub date: DateTime<Utc>,
    pub article: Option<serde_json::Value>, // Placeholder for aggregate.ArticleAggregate
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FeedbackAggregate {
    pub user_id: u32,
    pub news_id: u32,
    pub value: i32,      // 0 = dislike, 1 = like
    pub is_active: bool, // Remove feedback if false
    pub date: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(rename_all = "lowercase")]
pub enum FeedbackType {
    Like,
    Dislike,
}

impl InferenceCommand {
    pub fn validate_embedding_dimension(
        &self,
        expected: usize,
    ) -> Result<(), crate::error::DriftError> {
        if self.embedding.len() != expected {
            return Err(crate::error::DriftError::InvalidEmbeddingDimension {
                expected,
                actual: self.embedding.len(),
            });
        }
        Ok(())
    }
}
