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
}

variable "vault_address" {
  description = "Vault server address"
  type        = string
  default     = "http://localhost:8200"
}

variable "vault_root_token" {
  description = "Vault root token"
  type        = string
  sensitive   = true
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

resource "vault_kv_secret_v2" "postgres-minio" {
  mount = vault_mount.kv.path
  name  = "postgres/minio"

  data_json = jsonencode({
    user = var.postgres_minio_username
    password = var.postgres_minio_password
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

# ==========================================
# Variables for all secrets
# ==========================================

# PostgreSQL variables
variable "mlflow_postgres_username" {
  type      = string
  sensitive = true
}

variable "mlflow_postgres_password" {
  type      = string
  sensitive = true
}

variable "repo_account_postgres_username" {
  type      = string
  sensitive = true
}

variable "repo_account_postgres_password" {
  type      = string
  sensitive = true
}

variable "repo_article_postgres_username" {
  type      = string
  sensitive = true
}

variable "repo_article_postgres_password" {
  type      = string
  sensitive = true
}

variable "postgres_admin_username" {
  type      = string
  sensitive = true
}

variable "postgres_admin_password" {
  type      = string
  sensitive = true
}

variable "srvc_scrapping_postgres_username" {
  type      = string
  sensitive = true
}

variable "srvc_scrapping_postgres_password" {
  type      = string
  sensitive = true
}

# Cloudflare variables
variable "cloudflare_api_token" {
  type      = string
  sensitive = true
}

variable "cloudflare_email" {
  type      = string
  sensitive = true
}

# MinIO variables
variable "minio_root_user" {
  type      = string
  sensitive = true
}

variable "minio_root_password" {
  type      = string
  sensitive = true
}

variable "minio_console_username" {
  type      = string
  sensitive = true
}

variable "minio_console_password" {
  type      = string
  sensitive = true
}

variable "minio_console_disabled" {
  type    = string
  default = "false"
}

variable "minio_console_policies" {
  type    = string
  default = "consoleAdmin"
}

variable "minio_console_set_policies" {
  type    = string
  default = "false"
}

variable "minio_minio_username" {
  type      = string
  sensitive = true
}

variable "minio_minio_password" {
  type      = string
  sensitive = true
}

variable "minio_minio_disabled" {
  type    = string
  default = "false"
}

variable "minio_minio_policies" {
  type    = string
  default = "readwrite"
}

variable "minio_minio_set_policies" {
  type    = string
  default = "false"
}

# Capacitor variables
variable "capacitor_token" {
  type      = string
  sensitive = true
}

# OIDC variables
variable "oidc_jwt_secret" {
  type      = string
  sensitive = true
}

variable "oidc_client_id" {
  type      = string
  sensitive = true
}

variable "oidc_client_secret" {
  type      = string
  sensitive = true
}

variable "oidc_issuer_url" {
  type      = string
  sensitive = true
}

# Service Inference variables
variable "srvc_inference_tfidf_access_key_id" {
  type      = string
  sensitive = true
}

variable "srvc_inference_tfidf_secret_access_key" {
  type      = string
  sensitive = true
}

variable "srvc_inference_tfidf_redis_password" {
  type      = string
  sensitive = true
}

variable "srvc_inference_tfidf_redis_db" {
  type    = string
  default = "0"
}

# Docker variables
variable "docker_ghcr_auth" {
  type        = string
  sensitive   = true
  description = "Base64 encoded username:token for GitHub Container Registry"
}

# Repo Feed variables
variable "repo_feed_redis_password" {
  type      = string
  sensitive = true
}

variable "repo_feed_redis_db" {
  type    = string
  default = "0"
}

# Service Scrapping variables
variable "srvc_scrapping_webhook_url" {
  type      = string
  sensitive = true
}

# Qdrant variables
variable "qdrant_api_key" {
  type      = string
  sensitive = true
}

# ==========================================
# Outputs
# ==========================================

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
      vault_kv_secret_v2.srvc_scrapping_discord.name
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