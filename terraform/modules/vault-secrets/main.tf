# modules/vault-secret/main.tf
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
  token   = var.vault_token
}

resource "vault_kv_secret_v2" "this" {
  mount = var.mount_path
  name  = var.secret_path

  data_json = jsonencode(var.secrets)
}
