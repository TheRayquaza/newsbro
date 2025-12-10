use crate::config::DiscordConfig;
use crate::error::Result;
use crate::storage::postgres::FeedbackMetrics;
use reqwest::Client;
use serde_json::json;
use tracing::{error, info, instrument};

pub struct DiscordReporter {
    client: Client,
    config: DiscordConfig,
}

impl DiscordReporter {
    pub fn new(config: DiscordConfig) -> Self {
        Self {
            client: Client::new(),
            config,
        }
    }

    #[instrument(skip(self, metrics, drift_info))]
    pub async fn send_report(
        &self,
        metrics: &[FeedbackMetrics],
        drift_info: Option<DriftReport>,
    ) -> Result<()> {
        if metrics.is_empty() {
            info!("No metrics to report, skipping Discord notification");
            return Ok(());
        }

        let total_inferences: i32 = metrics.iter().map(|m| m.total_inferences).sum();
        let total_likes: i32 = metrics.iter().map(|m| m.likes).sum();
        let total_dislikes: i32 = metrics.iter().map(|m| m.dislikes).sum();

        let like_ratio = if total_inferences > 0 {
            total_likes as f64 / total_inferences as f64
        } else {
            0.0
        };

        let color = self.determine_color(like_ratio, drift_info.as_ref());

        let mut fields = vec![
            json!({
                "name": "📊 Total Inferences",
                "value": format!("{}", total_inferences),
                "inline": true
            }),
            json!({
                "name": "👍 Likes",
                "value": format!("{}", total_likes),
                "inline": true
            }),
            json!({
                "name": "👎 Dislikes",
                "value": format!("{}", total_dislikes),
                "inline": true
            }),
            json!({
                "name": "📈 Like Ratio",
                "value": format!("{:.2}%", like_ratio * 100.0),
                "inline": true
            }),
        ];

        if let Some(drift) = drift_info {
            fields.push(json!({
                "name": "🔄 Drift Status",
                "value": format!(
                    "Severity: **{}**\nCosine: {:.4}\nPSI: {:.4}\nKL: {:.4}",
                    drift.severity.to_uppercase(),
                    drift.cosine_distance.unwrap_or(0.0),
                    drift.psi_score.unwrap_or(0.0),
                    drift.kl_divergence.unwrap_or(0.0)
                ),
                "inline": false
            }));
        }

        let embed = json!({
            "title": "🤖 ML Model Performance Report",
            "color": color,
            "fields": fields,
            "timestamp": chrono::Utc::now().to_rfc3339(),
            "footer": {
                "text": "srvc-drift monitoring"
            }
        });

        let payload = json!({
            "username": self.config.username,
            "embeds": [embed]
        });

        let response = self
            .client
            .post(&self.config.webhook_url)
            .json(&payload)
            .send()
            .await?;

        if !response.status().is_success() {
            error!("Discord webhook failed with status: {}", response.status());
            let body = response.text().await.unwrap_or_default();
            error!("Response body: {}", body);
        } else {
            info!("Discord report sent successfully");
        }

        Ok(())
    }

    fn determine_color(&self, like_ratio: f64, drift: Option<&DriftReport>) -> u32 {
        if like_ratio < self.config.alert_thresholds.dislike_ratio_warning {
            return 0xE74C3C;
        }

        if let Some(drift_info) = drift {
            match drift_info.severity.as_str() {
                "high" => return 0xE74C3C,   // Red
                "medium" => return 0xF39C12, // Yellow
                "low" => return 0x3498DB,    // Blue
                _ => {}
            }
        }

        0x2ECC71 // Green
    }

    #[instrument(skip(self))]
    pub async fn send_alert(&self, message: &str, severity: AlertSeverity) -> Result<()> {
        let color = match severity {
            AlertSeverity::Critical => 0xE74C3C, // Red
            AlertSeverity::Warning => 0xF39C12,  // Yellow
            AlertSeverity::Info => 0x3498DB,     // Blue
        };

        let emoji = match severity {
            AlertSeverity::Critical => "🚨",
            AlertSeverity::Warning => "⚠️",
            AlertSeverity::Info => "ℹ️",
        };

        let embed = json!({
            "title": format!("{} Alert", emoji),
            "description": message,
            "color": color,
            "timestamp": chrono::Utc::now().to_rfc3339(),
        });

        let payload = json!({
            "username": self.config.username,
            "embeds": [embed]
        });

        self.client
            .post(&self.config.webhook_url)
            .json(&payload)
            .send()
            .await?;

        info!("Alert sent to Discord: {}", message);

        Ok(())
    }
}

#[derive(Debug, Clone)]
pub struct DriftReport {
    pub severity: String,
    pub cosine_distance: Option<f64>,
    pub psi_score: Option<f64>,
    pub kl_divergence: Option<f64>,
    pub sample_count: i32,
}

#[derive(Debug, Clone, Copy)]
pub enum AlertSeverity {
    Critical,
    Warning,
    Info,
}
