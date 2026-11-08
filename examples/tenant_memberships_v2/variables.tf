variable "nutanix_username" {
  type        = string
  description = "Username for Nutanix Central."
}

variable "nutanix_password" {
  type        = string
  description = "Password for Nutanix Central."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Nutanix Central endpoint."
}

variable "tenant_name" {
  type        = string
  description = "The name of the tenant to create for demonstration."
}

variable "tenant_admin_password" {
  type        = string
  description = "The password for the tenant local admin user."
  sensitive   = true
}

variable "domain_name" {
  type        = string
  description = "The name of the domain to create for demonstration."
}