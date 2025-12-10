use crate::config::Config;
use crate::error::{DriftError, Result};
use crate::reporting::discord::AlertSeverity;
use crate::reporting::discord::DiscordReporter;
use crate::storage::postgres::DriftSnapshot;
use crate::storage::postgres::PostgresStorage;
use chrono::Utc;
use std::sync::Arc;
use tokio::time::interval;
use tracing::instrument;
use tracing::{error, info};

pub struct DriftCalculator;

impl DriftCalculator {
    #[instrument(skip(embeddings))]
    pub fn calculate_centroid(embeddings: &[Vec<f32>]) -> Result<Vec<f32>> {
        if embeddings.is_empty() {
            return Err(DriftError::InsufficientData(
                "No embeddings provided".to_string(),
            ));
        }

        let dim = embeddings[0].len();
        let mut centroid = vec![0.0f32; dim];

        for embedding in embeddings {
            if embedding.len() != dim {
                // print embedding
                error!("Invalid embedding dimension: {} {:?}", dim, embedding);
                return Err(DriftError::InvalidEmbeddingDimension {
                    expected: dim,
                    actual: embedding.len(),
                });
            }

            for (i, &val) in embedding.iter().enumerate() {
                centroid[i] += val;
            }
        }

        let count = embeddings.len() as f32;
        for val in &mut centroid {
            *val /= count;
        }

        Ok(centroid)
    }

    #[instrument(skip(a, b))]
    pub fn cosine_distance(a: &[f32], b: &[f32]) -> Result<f64> {
        if a.len() != b.len() {
            return Err(DriftError::InvalidEmbeddingDimension {
                expected: a.len(),
                actual: b.len(),
            });
        }

        let dot_product: f32 = a.iter().zip(b.iter()).map(|(x, y)| x * y).sum();
        let norm_a: f32 = a.iter().map(|x| x * x).sum::<f32>().sqrt();
        let norm_b: f32 = b.iter().map(|x| x * x).sum::<f32>().sqrt();

        if norm_a == 0.0 || norm_b == 0.0 {
            return Ok(1.0); // Maximum distance
        }

        let cosine_similarity = dot_product / (norm_a * norm_b);
        Ok((1.0 - cosine_similarity as f64).clamp(0.0, 2.0))
    }

    #[instrument(skip(current, baseline))]
    pub fn calculate_psi(current: &[Vec<f32>], baseline: &[Vec<f32>]) -> Result<f64> {
        if current.is_empty() || baseline.is_empty() {
            return Err(DriftError::InsufficientData(
                "Insufficient data for PSI calculation".to_string(),
            ));
        }

        // Simplified PSI calculation on first dimension (for demonstration)
        // In production, you'd want more sophisticated binning strategies
        let bins = 10;
        let current_vals: Vec<f32> = current.iter().map(|v| v[0]).collect();
        let baseline_vals: Vec<f32> = baseline.iter().map(|v| v[0]).collect();

        let min_val = current_vals
            .iter()
            .chain(baseline_vals.iter())
            .copied()
            .min_by(|a, b| a.partial_cmp(b).unwrap())
            .unwrap();
        let max_val = current_vals
            .iter()
            .chain(baseline_vals.iter())
            .copied()
            .max_by(|a, b| a.partial_cmp(b).unwrap())
            .unwrap();

        let bin_width = (max_val - min_val) / bins as f32;
        if bin_width == 0.0 {
            return Ok(0.0);
        }

        let mut current_hist = vec![0; bins];
        let mut baseline_hist = vec![0; bins];

        for &val in &current_vals {
            let bin = ((val - min_val) / bin_width).floor() as usize;
            let bin = bin.min(bins - 1);
            current_hist[bin] += 1;
        }

        for &val in &baseline_vals {
            let bin = ((val - min_val) / bin_width).floor() as usize;
            let bin = bin.min(bins - 1);
            baseline_hist[bin] += 1;
        }

        let mut psi = 0.0;
        for i in 0..bins {
            let current_pct = (current_hist[i] as f64 + 0.0001) / current_vals.len() as f64;
            let baseline_pct = (baseline_hist[i] as f64 + 0.0001) / baseline_vals.len() as f64;
            psi += (current_pct - baseline_pct) * (current_pct / baseline_pct).ln();
        }

        Ok(psi.abs())
    }

    #[instrument(skip(current, baseline))]
    pub fn calculate_kl_divergence(current: &[Vec<f32>], baseline: &[Vec<f32>]) -> Result<f64> {
        if current.is_empty() || baseline.is_empty() {
            return Err(DriftError::InsufficientData(
                "Insufficient data for KL divergence".to_string(),
            ));
        }

        // Simplified KL divergence on first dimension
        let bins = 10;
        let current_vals: Vec<f32> = current.iter().map(|v| v[0]).collect();
        let baseline_vals: Vec<f32> = baseline.iter().map(|v| v[0]).collect();

        let min_val = current_vals
            .iter()
            .chain(baseline_vals.iter())
            .copied()
            .min_by(|a, b| a.partial_cmp(b).unwrap())
            .unwrap();
        let max_val = current_vals
            .iter()
            .chain(baseline_vals.iter())
            .copied()
            .max_by(|a, b| a.partial_cmp(b).unwrap())
            .unwrap();

        let bin_width = (max_val - min_val) / bins as f32;
        if bin_width == 0.0 {
            return Ok(0.0);
        }

        let mut current_hist = vec![0; bins];
        let mut baseline_hist = vec![0; bins];

        for &val in &current_vals {
            let bin = ((val - min_val) / bin_width).floor() as usize;
            let bin = bin.min(bins - 1);
            current_hist[bin] += 1;
        }

        for &val in &baseline_vals {
            let bin = ((val - min_val) / bin_width).floor() as usize;
            let bin = bin.min(bins - 1);
            baseline_hist[bin] += 1;
        }

        let mut kl = 0.0;
        for i in 0..bins {
            let p = (current_hist[i] as f64 + 0.0001) / current_vals.len() as f64;
            let q = (baseline_hist[i] as f64 + 0.0001) / baseline_vals.len() as f64;
            kl += p * (p / q).ln();
        }

        Ok(kl.max(0.0))
    }

    pub fn determine_severity(
        cosine_dist: f64,
        psi: f64,
        kl: f64,
        thresholds: &crate::config::AlertThresholds,
    ) -> String {
        let max_metric = cosine_dist.max(psi).max(kl);

        if max_metric >= thresholds.drift_critical {
            "high".to_string()
        } else if max_metric >= thresholds.drift_warning {
            "medium".to_string()
        } else if max_metric > 0.05 {
            "low".to_string()
        } else {
            "none".to_string()
        }
    }
}

pub async fn drift_calculation_loop(
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

        let embedding_json =
            serde_json::to_value(current_centroid).unwrap_or(serde_json::Value::Null);

        let snapshot = DriftSnapshot {
            snapshot_time: Utc::now(),
            embedding_centroid: embedding_json,
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_centroid_calculation() {
        let embeddings = vec![
            vec![1.0, 2.0, 3.0],
            vec![2.0, 3.0, 4.0],
            vec![3.0, 4.0, 5.0],
        ];

        let centroid = DriftCalculator::calculate_centroid(&embeddings).unwrap();
        assert_eq!(centroid, vec![2.0, 3.0, 4.0]);
    }

    #[test]
    fn test_cosine_distance() {
        let a = vec![1.0, 0.0, 0.0];
        let b = vec![1.0, 0.0, 0.0];
        let dist = DriftCalculator::cosine_distance(&a, &b).unwrap();
        assert!(dist < 0.001); // Should be very close to 0

        let c = vec![1.0, 0.0, 0.0];
        let d = vec![-1.0, 0.0, 0.0];
        let dist2 = DriftCalculator::cosine_distance(&c, &d).unwrap();
        assert!(dist2 > 1.99); // Should be close to 2
    }
}
