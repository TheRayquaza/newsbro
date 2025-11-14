terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 4.0"
    }
  }
}

provider "vault" {
  address = var.vault_address
  token   = var.vault_root_token
  skip_tls_verify= true
}

# Enable KV v2 secrets engine if not already enabled
resource "vault_mount" "kv" {
  path        = "kv"
  type        = "kv"
  options     = { version = "2" }
  description = "KV Version 2 secret engine mount"
}

# ==========================================
# PostgreSQL Database Secrets
# ==========================================

resource "vault_kv_secret_v2" "mlflow_postgres" {
  mount = vault_mount.kv.path
  name  = "mlflow/postgres"

  data_json = jsonencode({
    username = var.mlflow_postgres_username
    password = var.mlflow_postgres_password
  })
}

resource "vault_kv_secret_v2" "repo_account_postgres" {
  mount = vault_mount.kv.path
  name  = "repo-account/postgres"

  data_json = jsonencode({
    username = var.repo_account_postgres_username
    password = var.repo_account_postgres_password
  })
}

resource "vault_kv_secret_v2" "repo_article_postgres" {
  mount = vault_mount.kv.path
  name  = "repo-article/postgres"

  data_json = jsonencode({
    username = var.repo_article_postgres_username
    password = var.repo_article_postgres_password
  })
}

resource "vault_kv_secret_v2" "postgres_admin" {
  mount = vault_mount.kv.path
  name  = "postgres/admin"

  data_json = jsonencode({
    username = var.postgres_admin_username
    password = var.postgres_admin_password
  })
}

resource "vault_kv_secret_v2" "srvc_scrapping_postgres" {
  mount = vault_mount.kv.path
  name  = "srvc-scrapping/postgres"

  data_json = jsonencode({
    username = var.srvc_scrapping_postgres_username
    password = var.srvc_scrapping_postgres_password
  })
}

# ==========================================
# Cloudflare Secrets
# ==========================================

resource "vault_kv_secret_v2" "cloudflare" {
  mount = vault_mount.kv.path
  name  = "cloudflare"

  data_json = jsonencode({
    api_token = var.cloudflare_api_token
    email     = var.cloudflare_email
  })
}

# ==========================================
# MinIO Secrets
# ==========================================

resource "vault_kv_secret_v2" "minio_credentials" {
  mount = vault_mount.kv.path
  name  = "minio/credentials"

  data_json = jsonencode({
    rootUser     = var.minio_root_user
    rootPassword = var.minio_root_password
  })
}

resource "vault_kv_secret_v2" "minio_user_console" {
  mount = vault_mount.kv.path
  name  = "minio/users/console"

  data_json = jsonencode({
    username    = var.minio_console_username
    password    = var.minio_console_password
    disabled    = var.minio_console_disabled
    policies    = var.minio_console_policies
    setPolicies = var.minio_console_set_policies
  })
}

resource "vault_kv_secret_v2" "minio_user_minio" {
  mount = vault_mount.kv.path
  name  = "minio/users/minio"

  data_json = jsonencode({
    username    = var.minio_minio_username
    password    = var.minio_minio_password
    disabled    = var.minio_minio_disabled
    policies    = var.minio_minio_policies
    setPolicies = var.minio_minio_set_policies
  })
}

# ==========================================
# Capacitor Secrets
# ==========================================

resource "vault_kv_secret_v2" "capacitor" {
  mount = vault_mount.kv.path
  name  = "capacitor"

  data_json = jsonencode({
    token = var.capacitor_token
  })
}

# ==========================================
# Datadog Secrets
# ==========================================

resource "vault_kv_secret_v2" "datadog" {
  mount = vault_mount.kv.path
  name  = "datadog"

  data_json = jsonencode({
    api_key = var.datadog_api_key
  })
}

# ==========================================
# OIDC Secrets
# ==========================================

resource "vault_kv_secret_v2" "oidc" {
  mount = vault_mount.kv.path
  name  = "oidc"

  data_json = jsonencode({
    JWT_SECRET    = var.oidc_jwt_secret
    client_id     = var.oidc_client_id
    client_secret = var.oidc_client_secret
    issuer_url    = var.oidc_issuer_url
  })
}

# ==========================================
# Service Inference Secrets
# ==========================================

resource "vault_kv_secret_v2" "srvc_inference_tfidf_s3" {
  mount = vault_mount.kv.path
  name  = "srvc-inference/tfidf/s3"

  data_json = jsonencode({
    access_key_id     = var.srvc_inference_tfidf_access_key_id
    secret_access_key = var.srvc_inference_tfidf_secret_access_key
  })
}

resource "vault_kv_secret_v2" "srvc_inference_tfidf_qdrant" {
  mount = vault_mount.kv.path
  name  = "srvc-inference/tfidf/qdrant"

  data_json = jsonencode({
    qdrant_api_key = var.qdrant_api_key
  })
}

resource "vault_kv_secret_v2" "srvc_inference_tfidf_redis" {
  mount = vault_mount.kv.path
  name  = "srvc-inference/tfidf/redis"

  data_json = jsonencode({
    password = var.srvc_inference_tfidf_redis_password
    db       = var.srvc_inference_tfidf_redis_db
  })
}

# ==========================================
# Docker Registry Secrets
# ==========================================

resource "vault_kv_secret_v2" "docker_ghcr" {
  mount = vault_mount.kv.path
  name  = "docker/ghcr"

  data_json = jsonencode({
    auth = var.docker_ghcr_auth
  })
}

# ==========================================
# Repo Feed Secrets
# ==========================================

resource "vault_kv_secret_v2" "repo_feed_redis" {
  mount = vault_mount.kv.path
  name  = "repo-feed/redis"

  data_json = jsonencode({
    password = var.repo_feed_redis_password
    db       = var.repo_feed_redis_db
  })
}

# ==========================================
# Service Scrapping Secrets
# ==========================================

resource "vault_kv_secret_v2" "srvc_scrapping_discord" {
  mount = vault_mount.kv.path
  name  = "srvc-scrapping/discord"

  data_json = jsonencode({
    webhook_url = var.srvc_scrapping_webhook_url
  })
}

# ==========================================
# Qdrant Secrets
# ==========================================

resource "vault_kv_secret_v2" "qdrant" {
  mount = vault_mount.kv.path
  name  = "qdrant"

  data_json = jsonencode({
    api_key = var.qdrant_api_key
  })
}
