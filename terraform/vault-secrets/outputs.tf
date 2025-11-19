output "vault_mount_path" {
  value       = vault_mount.kv.path
  description = "The path where the KV secrets engine is mounted"
}

output "secrets_created" {
  value = {
    postgres = [
      vault_kv_secret_v2.mlflow_postgres.name,
      vault_kv_secret_v2.repo_account_postgres.name,
      vault_kv_secret_v2.repo_article_postgres.name,
      vault_kv_secret_v2.postgres_admin.name,
      vault_kv_secret_v2.srvc_scrapping_postgres.name
    ]
    minio = [
      vault_kv_secret_v2.minio_credentials.name,
      vault_kv_secret_v2.minio_user_console.name,
      vault_kv_secret_v2.minio_user_minio.name
    ]
    services = [
      vault_kv_secret_v2.srvc_inference_tfidf_s3.name,
      vault_kv_secret_v2.srvc_inference_tfidf_qdrant.name,
      vault_kv_secret_v2.srvc_inference_tfidf_redis.name,
      vault_kv_secret_v2.repo_feed_redis.name,
      vault_kv_secret_v2.srvc_scrapping_discord.name,
      vault_kv_secret_v2.datadog.name
    ]
    infrastructure = [
      vault_kv_secret_v2.cloudflare.name,
      vault_kv_secret_v2.capacitor.name,
      vault_kv_secret_v2.oidc.name,
      vault_kv_secret_v2.docker_ghcr.name,
      vault_kv_secret_v2.qdrant.name
    ]
  }
  description = "List of all secrets created in Vault"
}
