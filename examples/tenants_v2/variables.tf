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

variable "tenant_name_1" {
  type        = string
  description = "The name of the first tenant to create for demonstration."
}

variable "tenant_name_2" {
  type        = string
  description = "The name of the second tenant to create for demonstration."
}

variable "tenant_admin_password" {
  type        = string
  description = "The password for the tenant local admin user."
  sensitive   = true
}