#define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central account."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central account."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)."
  type        = string
}
