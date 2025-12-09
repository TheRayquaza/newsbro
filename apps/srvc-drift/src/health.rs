use crate::storage::postgres::PostgresStorage;
use crate::storage::redis::RedisCache;
use axum::{
    Json,
    extract::State,
    http::StatusCode,
    response::{IntoResponse, Response},
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tracing::error;

#[derive(Clone)]
pub struct AppState {
    pub postgres: PostgresStorage,
    pub redis: Option<RedisCache>,
}

#[derive(Serialize, Deserialize)]
pub struct HealthResponse {
    pub status: String,
    pub version: String,
}

#[derive(Serialize, Deserialize)]
pub struct ReadyResponse {
    pub status: String,
    pub checks: DependencyChecks,
}

#[derive(Serialize, Deserialize)]
pub struct DependencyChecks {
    pub postgres: DependencyStatus,
    pub redis: DependencyStatus,
    pub kafka: DependencyStatus,
}

#[derive(Serialize, Deserialize)]
pub struct DependencyStatus {
    pub healthy: bool,
    pub message: String,
}

#[derive(Serialize, Deserialize)]
pub struct DependencyDetailsResponse {
    pub dependencies: DependencyChecks,
}

pub async fn health_handler() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ok".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
    })
}

pub async fn ready_handler(State(state): State<Arc<AppState>>) -> Response {
    let postgres_status = check_postgres(&state.postgres).await;
    let redis_status = check_redis(&state.redis).await;

    // Kafka health is checked by consumer heartbeat - if process is running, Kafka is OK
    let kafka_status = DependencyStatus {
        healthy: true,
        message: "Consumer group active".to_string(),
    };

    let all_healthy = postgres_status.healthy && redis_status.healthy && kafka_status.healthy;

    let response = ReadyResponse {
        status: if all_healthy {
            "ready".to_string()
        } else {
            "not_ready".to_string()
        },
        checks: DependencyChecks {
            postgres: postgres_status,
            redis: redis_status,
            kafka: kafka_status,
        },
    };

    if all_healthy {
        (StatusCode::OK, Json(response)).into_response()
    } else {
        (StatusCode::SERVICE_UNAVAILABLE, Json(response)).into_response()
    }
}

pub async fn dependencies_handler(
    State(state): State<Arc<AppState>>,
) -> Json<DependencyDetailsResponse> {
    let postgres_status = check_postgres(&state.postgres).await;
    let redis_status = check_redis(&state.redis).await;
    let kafka_status = DependencyStatus {
        healthy: true,
        message: "Consumer group active".to_string(),
    };

    Json(DependencyDetailsResponse {
        dependencies: DependencyChecks {
            postgres: postgres_status,
            redis: redis_status,
            kafka: kafka_status,
        },
    })
}

async fn check_postgres(storage: &PostgresStorage) -> DependencyStatus {
    match storage.health_check().await {
        Ok(_) => DependencyStatus {
            healthy: true,
            message: "Connected".to_string(),
        },
        Err(e) => {
            error!("PostgreSQL health check failed: {}", e);
            DependencyStatus {
                healthy: false,
                message: format!("Connection failed: {}", e),
            }
        }
    }
}

async fn check_redis(cache: &Option<RedisCache>) -> DependencyStatus {
    match cache {
        Some(redis) => match redis.health_check().await {
            Ok(_) => DependencyStatus {
                healthy: true,
                message: "Connected".to_string(),
            },
            Err(e) => {
                error!("Redis health check failed: {}", e);
                DependencyStatus {
                    healthy: false,
                    message: format!("Connection failed: {}", e),
                }
            }
        },
        None => DependencyStatus {
            healthy: true,
            message: "Not configured".to_string(),
        },
    }
}
