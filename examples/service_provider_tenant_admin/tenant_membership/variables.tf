variable "nutanix_username" {
  type        = string
  description = "Username for the Tenant Admin principal."
}

variable "nutanix_password" {
  type        = string
  description = "Password for the Tenant Admin principal."
  sensitive   = true
}

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint base FQDN."
}

variable "domain_name" {
  type        = string
  description = "Domain name of the Prism Central allocated to the tenant."
}

variable "tenant_name" {
  type        = string
  description = "Tenant name used when operating as Tenant Admin."
}

variable "domain_ext_id" {
  type        = string
  description = "External identifier of the domain the membership is created in."
}

variable "role_ext_id" {
  type        = string
  description = "External identifier of the role assigned to the principal."
}

variable "identity_ext_id" {
  type        = string
  description = "External identifier of an existing principal entity (used by the by_identity_ext_id example)."
}

variable "identity_provider_ext_id" {
  type        = string
  description = "External identifier of the identity provider (used by the by_identity_details example)."
}

variable "principal_value" {
  type        = string
  description = "The unique lookup value for the user, such as a UPN, username, or email address."
}

variable "identity_type" {
  type        = string
  description = "The identity type: USER or GROUP."
  default     = "USER"
}
