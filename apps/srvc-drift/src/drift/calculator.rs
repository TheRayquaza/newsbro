use crate::error::{DriftError, Result};
use tracing::instrument;

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
        Ok((1.0 - cosine_similarity as f64).max(0.0).min(2.0))
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
