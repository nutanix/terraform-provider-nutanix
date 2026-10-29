variable "nutanix_central_username" {
  type        = string
  description = "Username for Nutanix Central."
}

variable "nutanix_central_password" {
  type        = string
  description = "Password for Nutanix Central."
  sensitive   = true
}

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint."
}

variable "domain_ext_id" {
  type        = string
  description = "External identifier of the existing domain to allow for tenancy."
}
