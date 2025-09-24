module "capacitor" {
  source       = "./modules/vault-secrets"
  vault_address = var.vault_address
  vault_token   = var.vault_token
  mount_path    = "kv"
  secret_path   = "capacitor"

  secrets = {
    LICENSE_KEY = var.capacitor_license_key,
    AUTH = var.capacitor_auth,
    IMPERSONATE_SA_RULES = var.capacitor_impersonate_sa_rules,
    SESSION_HASH_KEY =var.capacitor_session_hash_key,
    SESSION_BLOCK_KEY = var.capacitor_session_block_key,
    "registry.yaml" = var.capacitor_registry_yaml
  }
}