use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InferenceCommand {
    pub user_id: u32,
    pub model: String,
    pub score: f64,
    pub date: DateTime<Utc>,
    pub article: Option<ArticleAggregate>, // Placeholder for aggregate.ArticleAggregate
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ArticleAggregate {
    pub id: u32,
    pub category: String,
    pub subcategory: String,
    pub title: String,
    #[serde(rename = "abstract")]
    pub abstract_text: String,
    pub link: String,
    pub rss_link: String,
    pub published_at: DateTime<Utc>,
    pub is_active: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FeedbackAggregate {
    pub user_id: u32,
    pub news_id: u32,
    pub value: i32,      // 0 = dislike, 1 = like
    pub is_active: bool, // Remove feedback if false
    pub date: DateTime<Utc>,
}
