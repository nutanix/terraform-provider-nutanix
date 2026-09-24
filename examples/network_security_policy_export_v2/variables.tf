#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "Username for Nutanix Prism Central."
}

variable "nutanix_password" {
  type        = string
  sensitive   = true
  description = "Password for Nutanix Prism Central."
}

variable "nutanix_endpoint" {
  type        = string
  description = "Endpoint (IP or FQDN) for Nutanix Prism Central."
}
