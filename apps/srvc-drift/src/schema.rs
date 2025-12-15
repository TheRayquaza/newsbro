// @generated automatically by Diesel CLI.

diesel::table! {
    feedback_metrics (id) {
        id -> Int8,
        time_bucket -> Timestamptz,
        bucket_size_seconds -> Int4,
        total_inferences -> Int4,
        likes -> Int4,
        dislikes -> Int4,
        like_ratio -> Float8,
        created_at -> Nullable<Timestamptz>,
    }
}

diesel::table! {
    drift_snapshots (id) {
        id -> Int8,
        snapshot_time -> Timestamptz,
        embedding_centroid -> Jsonb,
        psi_score -> Nullable<Float8>,
        kl_divergence -> Nullable<Float8>,
        cosine_distance_from_baseline -> Nullable<Float8>,
        sample_count -> Int4,
        drift_severity -> Text,
        created_at -> Nullable<Timestamptz>,
    }
}

diesel::table! {
    inference_log (id) {
        id -> Int8,
        user_id -> Int4,
        embedding -> Jsonb,
        model -> Text,
        score -> Float4,
        date -> Timestamptz,
        created_at -> Nullable<Timestamptz>,
    }
}

diesel::table! {
    feedback_records (id) {
        id -> Int8,
        user_id -> Int4,
        feedback_type -> Text,
        timestamp -> Timestamptz,
        metadata -> Nullable<Jsonb>,
        created_at -> Nullable<Timestamptz>,
    }
}

diesel::allow_tables_to_appear_in_same_query!(
    feedback_metrics,
    drift_snapshots,
    inference_log,
    feedback_records,
);
