## Vault

variable "vault_token" {
  type = string
  sensitive = true
}

variable "vault_address" {
  type = string
}

## Capacitor

variable "capacitor_license_key" {
  type = string
    sensitive = true

}

variable "capacitor_auth" {
  type = string
    sensitive = true

}

variable "capacitor_impersonate_sa_rules" {
  type = string
}

variable "capacitor_session_hash_key" {
  type = string
}

variable "capacitor_session_block_key" {
  type = string
  sensitive = true
}

variable "capacitor_registry_yaml" {
  type = string
}

## Cloudflare

variable "cloudflare_token" {
  type = string
  sensitive = true
}
