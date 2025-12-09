use crate::config::RedisConfig;
use crate::error::Result;
use crate::kafka::models::FeedbackType;
use redis::aio::ConnectionManager;
use redis::{AsyncCommands, Client};
use tracing::{info, instrument};

#[derive(Debug, Clone)]
pub struct RedisCache {
    manager: ConnectionManager,
}

impl RedisCache {
    pub async fn new(config: &RedisConfig) -> Result<Self> {
        let client = Client::open(config.url.as_str())?;
        let manager = ConnectionManager::new(client).await?;

        info!("Redis connection manager established");

        Ok(Self { manager })
    }

    #[instrument(skip(self))]
    pub async fn increment_feedback(&self, feedback_type: &FeedbackType) -> Result<()> {
        let key = match feedback_type {
            FeedbackType::Like => "feedback:likes",
            FeedbackType::Dislike => "feedback:dislikes",
        };

        let mut conn = self.manager.clone();
        conn.incr::<_, _, ()>(key, 1).await?;

        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn get_feedback_counts(&self) -> Result<(i64, i64)> {
        let mut conn = self.manager.clone();

        let likes: i64 = conn.get("feedback:likes").await.unwrap_or(0);
        let dislikes: i64 = conn.get("feedback:dislikes").await.unwrap_or(0);

        Ok((likes, dislikes))
    }

    #[instrument(skip(self))]
    pub async fn reset_feedback_counts(&self) -> Result<()> {
        let mut conn = self.manager.clone();

        conn.set::<_, _, ()>("feedback:likes", 0).await?;
        conn.set::<_, _, ()>("feedback:dislikes", 0).await?;

        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn increment_inference_count(&self) -> Result<()> {
        let mut conn = self.manager.clone();
        conn.incr::<_, _, ()>("inference:count", 1).await?;
        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn get_and_reset_inference_count(&self) -> Result<i64> {
        let mut conn = self.manager.clone();
        let count: i64 = conn.get_del("inference:count").await.unwrap_or(0);
        Ok(count)
    }

    #[instrument(skip(self, response_time))]
    pub async fn add_response_time(&self, response_time: i32) -> Result<()> {
        let mut conn = self.manager.clone();
        conn.rpush::<_, _, ()>("response_times", response_time)
            .await?;
        conn.ltrim::<_, _, ()>("response_times", -1000, -1).await?; // Keep last 1000
        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn get_response_times(&self) -> Result<Vec<i32>> {
        let mut conn = self.manager.clone();
        let times: Vec<i32> = conn
            .lrange("response_times", 0, -1)
            .await
            .unwrap_or_default();
        Ok(times)
    }

    pub async fn health_check(&self) -> Result<()> {
        let mut conn = self.manager.clone();
        redis::cmd("PING")
            .query_async::<_, String>(&mut conn)
            .await?;
        Ok(())
    }
}
