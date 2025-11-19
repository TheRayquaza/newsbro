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