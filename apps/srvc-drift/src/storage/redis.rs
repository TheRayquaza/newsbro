use crate::error::DriftError;
use crate::{config::RedisConfig, error::Result};
use redis::aio::ConnectionManager;
use redis::{AsyncCommands, Client};
use tracing::{info, instrument};

#[derive(Clone)]
pub struct RedisCache {
    manager: ConnectionManager,
}

impl RedisCache {
    pub async fn new(config: &RedisConfig) -> Result<Self> {
        let client = Self::create_sentinel_client(&config.url)?;
        let manager = ConnectionManager::new(client).await?;

        info!("Redis Sentinel connection manager established");

        Ok(Self { manager })
    }

    fn create_sentinel_client(sentinel_url: &str) -> Result<Client> {
        let client = Client::open(sentinel_url)?;
        Ok(client)
    }

    #[instrument(skip(self))]
    pub async fn increment_feedback(&self, feedback_value: &i32) -> Result<()> {
        let key = match feedback_value {
            1 => "feedback:likes",
            0 => "feedback:dislikes",
            _ => {
                return Err(DriftError::InvalidInput(
                    "Feedback value must be 0 or 1".to_string(),
                ));
            }
        };

        let mut conn = self.manager.clone();
        conn.incr::<_, _, ()>(key, 1).await?;

        Ok(())
    }

    #[instrument(skip(self))]
    pub async fn increment_inference_count(&self) -> Result<()> {
        let mut conn = self.manager.clone();
        conn.incr::<_, _, ()>("inference:count", 1).await?;
        Ok(())
    }

    pub async fn health_check(&self) -> Result<()> {
        let mut conn = self.manager.clone();
        redis::cmd("PING")
            .query_async::<_, String>(&mut conn)
            .await?;
        Ok(())
    }
}
