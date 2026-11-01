#define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username for the Nutanix Central / Prism Central API."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Central / Prism Central API."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Nutanix Central / Prism Central endpoint (IP or FQDN)."
  type        = string
}

variable "registered_domain_name" {
  description = "Display name of the registered domain."
  type        = string
}

variable "registered_domain_fqdn" {
  description = "Fully qualified domain name of the registered domain (Prism Central)."
  type        = string
}

variable "client_signing_request" {
  description = "Client-signed request value used to establish trust with Nutanix Central."
  type        = string
  sensitive   = true
}

variable "registered_domain_api_key" {
  description = "API key used by Prism Central to authenticate to Nutanix Central."
  type        = string
  sensitive   = true
}
