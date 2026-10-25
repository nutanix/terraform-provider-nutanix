#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "Username to connect to Prism Central."
}

variable "nutanix_password" {
  type        = string
  description = "Password to connect to Prism Central."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Prism Central endpoint (IP or FQDN)."
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
