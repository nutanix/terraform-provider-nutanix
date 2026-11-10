variable "nutanix_central_username" {
  type        = string
  description = "Username for Nutanix Central (Service Provider Admin)."
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

variable "tenant_name" {
  type        = string
  description = "The name of the tenant to create."
}

variable "tenant_description" {
  type        = string
  description = "A brief description of the tenant."
  default     = "Tenant created via Terraform"
}

variable "admin_username" {
  type        = string
  description = "The username of the tenant admin user."
}

variable "admin_first_name" {
  type        = string
  description = "The first name of the tenant admin user."
}

variable "admin_last_name" {
  type        = string
  description = "The last name of the tenant admin user."
}

variable "admin_email" {
  type        = string
  description = "The email address of the tenant admin user."
}

variable "admin_password" {
  type        = string
  description = "The password of the tenant admin user."
  sensitive   = true
}

variable "quota_memory" {
  type        = number
  description = "The amount of memory allocated as quota, in bytes."
  default     = 8589934592 # 8 GB
}

variable "quota_storage" {
  type        = number
  description = "The amount of storage allocated as quota, in bytes."
  default     = 53687091200 # 50 GB
}

variable "quota_vcpu" {
  type        = number
  description = "The number of vCPUs allocated as quota."
  default     = 4
}
