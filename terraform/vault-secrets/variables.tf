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

variable "capacitor_auth" {
  type    = string
  default = "noauth"
}

variable "capacitor_impersonate_sa_rules" {
  type    = string
  default = "noauth=flux-system:capacitor-next-builtin-editor"
}

variable "capacitor_session_hash_key" {
  type      = string
  sensitive = true
}

variable "capacitor_session_block_key" {
  type      = string
  sensitive = true
}

variable "capacitor_registry_yaml" {
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

variable "srvc_inference_sbert_access_key_id" {
  type      = string
  sensitive = true
}

variable "srvc_inference_sbert_secret_access_key" {
  type      = string
  sensitive = true
}

variable "srvc_inference_sbert_redis_password" {
  type      = string
  sensitive = true
}

variable "srvc_inference_sbert_redis_db" {
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

# Datadog variables
variable "datadog_api_key" {
  type      = string
  sensitive = true
}

variable "datadog_app_key" {
  type      = string
  sensitive = true
}
