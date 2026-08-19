#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "Prism Central username."
}

variable "nutanix_password" {
  type        = string
  description = "Prism Central password."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Prism Central endpoint (IP or FQDN)."
}

variable "nc_api_key" {
  type        = string
  description = "API key issued by Nutanix Central, used by Prism Central for authentication."
  sensitive   = true
}

variable "nc_key_id" {
  type        = string
  description = "Key id used by Prism Central to form the authentication token for Nutanix Central."
  sensitive   = true
}

variable "nc_domain_ext_id" {
  type        = string
  description = "Unique identifier for the Prism Central (domain ext id)."
}

variable "nc_identity_provider_ext_id" {
  type        = string
  description = "Identity provider unique identifier (directory-service or SAML based)."
}

variable "nc_target_url" {
  type        = string
  description = "URL Prism Central should use to communicate with Nutanix Central."
}

variable "nc_tenant_ext_id" {
  type        = string
  description = "Unique identifier of the tenant."
}
