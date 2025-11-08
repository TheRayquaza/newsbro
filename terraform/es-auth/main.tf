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

variable "kubernetes_host" {
  description = "Kubernetes API host"
  type        = string
  default     = "https://kubernetes.default.svc.cluster.local"
}

variable "kubernetes_ca_cert" {
  description = "Kubernetes CA certificate (base64 encoded)"
  type        = string
  sensitive   = true
}

variable "vault_auth_token" {
  description = "Token from vault-auth service account"
  type        = string
  sensitive   = true
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
# Create Policy for External Secrets
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

# List capability for the secret mount
path "secret/metadata/*" {
  capabilities = ["list"]
}
EOT
}

# ==========================================
# Create Kubernetes Auth Role for flux-system
# ==========================================

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
# Outputs
# ==========================================

output "kubernetes_auth_path" {
  value       = vault_auth_backend.kubernetes.path
  description = "Path where Kubernetes auth is mounted"
}

output "policy_name" {
  value       = vault_policy.external_secrets.name
  description = "Name of the External Secrets policy"
}

output "flux_system_role" {
  value       = vault_kubernetes_auth_backend_role.flux_system.role_name
  description = "Kubernetes auth role name for flux-system"
}

output "auth_configuration" {
  value = {
    kubernetes_host = var.kubernetes_host
    auth_path       = vault_auth_backend.kubernetes.path
    policy          = vault_policy.external_secrets.name
  }
  description = "Authentication configuration summary"
}
