use crate::config::QdrantConfig;
use crate::error::Result;
use qdrant_client::Qdrant;
use qdrant_client::qdrant::GetPointsBuilder;
use qdrant_client::qdrant::PointId;
use qdrant_client::qdrant::point_id::PointIdOptions;
use qdrant_client::qdrant::vectors_output::VectorsOptions;
use tracing::{info, instrument};

#[derive(Clone)]
pub struct QdrantStorage {
    client: Qdrant,
    config: QdrantConfig,
}

impl QdrantStorage {
    pub async fn new(config: &QdrantConfig) -> Result<Self> {
        let client = Qdrant::from_url(config.url.as_str())
            .api_key(config.api_key.clone())
            .build()
            .expect("Failed to create Qdrant client");

        info!("Qdrant client established");

        Ok(Self {
            client,
            config: config.clone(),
        })
    }

    #[instrument(skip(self))]
    pub async fn retrieve_embedding_by_id(&self, id: u32) -> Result<Option<Vec<f32>>> {
        let search_result = self
            .client
            .get_points(
                GetPointsBuilder::new(
                    &self.config.collection_name,
                    vec![PointId {
                        point_id_options: Some(PointIdOptions::Num(id as u64)),
                    }],
                )
                .with_vectors(true),
            )
            .await
            .unwrap();

        if let Some(point) = search_result.result.first()
            && let Some(vectors) = &point.vectors
        {
            match vectors.vectors_options.as_ref().unwrap() {
                VectorsOptions::Vector(vector) => {
                    return Ok(Some(vector.data.clone()));
                }
                _ => {
                    return Ok(None);
                }
            }
        }
        Ok(None)
    }
}
