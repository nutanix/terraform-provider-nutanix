#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "The username for the Nutanix Prism Central."
}

variable "nutanix_password" {
  type        = string
  description = "The password for the Nutanix Prism Central."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "The endpoint (IP/FQDN) for the Nutanix Prism Central."
}

variable "tenant_name" {
  type        = string
  description = "The name of the tenant to create."
}

variable "tenant_admin_password" {
  type        = string
  description = "The password for the tenant local admin user."
  sensitive   = true
}
