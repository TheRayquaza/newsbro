variable "dns_records" {
    type = list(object({
        name    = string
        type    = string
        content = string
        ttl     = number
        proxied = bool
    }))
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