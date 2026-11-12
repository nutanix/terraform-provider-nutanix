variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint (e.g., nconprem-xxx.nxncm.com)."
}

variable "nutanix_central_username" {
  type        = string
  description = "Nutanix Central LDAP username."
}

variable "nutanix_central_password" {
  type        = string
  sensitive   = true
  description = "Nutanix Central admin password."
}

variable "nc_admin_role_name" {
  type        = string
  description = "Nutanix Central admin role name."
}

variable "service_provider_admin_role_name" {
  type        = string
  description = "Nutanix Central service provider admin role name."
}