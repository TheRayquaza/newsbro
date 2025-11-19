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
