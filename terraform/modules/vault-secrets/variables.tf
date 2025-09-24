# modules/vault-secret/variables.tf
variable "vault_address" {
  description = "Vault server address"
  type        = string
}

variable "vault_token" {
  description = "Vault root or management token"
  type        = string
  sensitive   = true
}

variable "mount_path" {
  description = "KV mount path (default: secret)"
  type        = string
  default     = "secret"
}

variable "secret_path" {
  description = "Secret path under the mount"
  type        = string
}

variable "secrets" {
  description = "Map of key/value secrets to store"
  type        = map(string)
}
