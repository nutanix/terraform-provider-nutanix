variable "nutanix_central_username" {
  type        = string
  description = "Username for Nutanix Central (NC Admin)."
}

variable "nutanix_central_password" {
  type        = string
  description = "Password for Nutanix Central."
  sensitive   = true
}

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint base FQDN."
}

variable "nutanix_pc_username" {
  type        = string
  description = "Prism Central username used for local-domain registration."
}

variable "nutanix_pc_password" {
  type        = string
  description = "Prism Central password used for local-domain registration."
  sensitive   = true
}

variable "nutanix_pc_endpoint" {
  type        = string
  description = "Prism Central endpoint (IP or FQDN) used for local-domain registration."
}

variable "registered_domain_name" {
  type        = string
  description = "Display name of the registered domain."
}

variable "registered_domain_fqdn" {
  type        = string
  description = "Fully qualified domain name of the registered Prism Central."
}

variable "domain_ext_id" {
  type        = string
  description = "Unique identifier for the Prism Central domain."
}

variable "location_ext_id" {
  type        = string
  description = "Location external identifier associated with the registered domain."
}

variable "nc_target_url" {
  type        = string
  description = "URL Prism Central should use to communicate with Nutanix Central."
}

variable "nc_tenant_ext_id" {
  type        = string
  description = "Unique identifier of the tenant associated with the registration."
}

variable "registered_domain_api_key" {
  type        = string
  description = "API key used by Prism Central to authenticate to Nutanix Central."
  sensitive   = true
}

variable "registered_domain_key_id" {
  type        = string
  description = "Key id paired with the registered domain API key."
  sensitive   = true
}

variable "tunnel_endpoint" {
  type        = string
  description = "Tunnel endpoint used for NC-PC connectivity."
  default     = "https://tunnel.example.com"
}

variable "client_signing_request" {
  type        = string
  description = "Client-signed request value used to establish trust with Nutanix Central."
  sensitive   = true
}

variable "nc_api_key" {
  type        = string
  description = "API key issued by Nutanix Central for local-domain registration on PC."
  sensitive   = true
}

variable "nc_key_id" {
  type        = string
  description = "Key id used by Prism Central for the Nutanix Central auth token."
  sensitive   = true
}

variable "nc_identity_provider_ext_id" {
  type        = string
  description = "Identity provider unique identifier (directory-service or SAML based)."
}
