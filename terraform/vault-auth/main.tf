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

# ==========================================
# Enable Kubernetes Auth Method
# ==========================================

resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
  path = "kubernetes"
  
  description = "Kubernetes auth backend for External Secrets Operator"
}

resource "vault_kubernetes_auth_backend_config" "kubernetes" {
  backend            = vault_auth_backend.kubernetes.path
  kubernetes_host    = var.kubernetes_host
  kubernetes_ca_cert = var.kubernetes_ca_cert
  token_reviewer_jwt = var.vault_auth_token
}

# ==========================================
# External Secrets Operator Role
# ==========================================

resource "vault_policy" "external_secrets" {
  name = "external-secrets-policy"

  policy = <<EOT
# Read access to all secrets paths used by ExternalSecrets

# PostgreSQL secrets
path "kv/data/mlflow/postgres" {
  capabilities = ["read", "list"]
}

path "kv/data/repo-account/postgres" {
  capabilities = ["read", "list"]
}

path "kv/data/repo-article/postgres" {
  capabilities = ["read", "list"]
}

path "kv/data/postgres/admin" {
  capabilities = ["read", "list"]
}

path "kv/data/srvc-scrapping/postgres" {
  capabilities = ["read", "list"]
}

# Cloudflare secrets
path "kv/data/cloudflare" {
  capabilities = ["read", "list"]
}

# MinIO secrets
path "kv/data/minio/credentials" {
  capabilities = ["read", "list"]
}

path "kv/data/minio/users/console" {
  capabilities = ["read", "list"]
}

path "kv/data/minio/users/minio" {
  capabilities = ["read", "list"]
}

# Capacitor secrets
path "kv/data/capacitor" {
  capabilities = ["read", "list"]
}

# OIDC secrets
path "kv/data/oidc" {
  capabilities = ["read", "list"]
}

# Service Inference secrets
path "kv/data/srvc-inference/tfidf/s3" {
  capabilities = ["read", "list"]
}

path "kv/data/srvc-inference/tfidf/qdrant" {
  capabilities = ["read", "list"]
}

path "kv/data/srvc-inference/tfidf/redis" {
  capabilities = ["read", "list"]
}

path "kv/data/srvc-inference/sbert/s3" {
  capabilities = ["read", "list"]
}

path "kv/data/srvc-inference/sbert/qdrant" {
  capabilities = ["read", "list"]
}

path "kv/data/srvc-inference/sbert/redis" {
  capabilities = ["read", "list"]
}

# Docker registry secrets
path "kv/data/docker/ghcr" {
  capabilities = ["read", "list"]
}

# Repo Feed secrets
path "kv/data/repo-feed/redis" {
  capabilities = ["read", "list"]
}

# Service Scrapping secrets
path "kv/data/srvc-scrapping/discord" {
  capabilities = ["read", "list"]
}

# Qdrant secrets
path "kv/data/qdrant" {
  capabilities = ["read", "list"]
}

# Datadog secrets
path "kv/data/datadog" {
  capabilities = ["read", "list"]
}

# List capability for the secret mount
path "secret/metadata/*" {
  capabilities = ["list"]
}
EOT
}

resource "vault_kubernetes_auth_backend_role" "flux_system" {
  backend                          = vault_auth_backend.kubernetes.path
  role_name                        = "flux-system"
  bound_service_account_names      = ["vault-auth"]
  bound_service_account_namespaces = ["kube-system"]
  token_ttl                        = 3600
  token_max_ttl                    = 86400
  token_policies                   = [vault_policy.external_secrets.name]
  audience                         = null
}

# ==========================================
# Datadog Integration
# ==========================================

resource "vault_policy" "datadog_agent" {
  name = "datadog-agent-policy"
  
  policy = <<EOT
path "kv/data/datadog" {
  capabilities = ["read", "list"]
}

# Read access to Vault metrics for monitoring
path "sys/metrics" {
  capabilities = ["read"]
}

# Health check endpoint
path "sys/health" {
  capabilities = ["read"]
}

# Seal status for monitoring
path "sys/seal-status" {
  capabilities = ["read"]
}
EOT
}

resource "vault_kubernetes_auth_backend_role" "datadog_agent" {
  backend                          = vault_auth_backend.kubernetes.path
  role_name                        = "datadog-agent"
  bound_service_account_names      = ["datadog-agent"]
  bound_service_account_namespaces = ["datadog"]
  token_ttl                        = 3600
  token_max_ttl                    = 86400
  token_policies                   = [vault_policy.datadog_agent.name]
  audience                         = null
}
