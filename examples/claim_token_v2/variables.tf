# define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP address or FQDN)."
  type        = string
}

variable "claim_token_name" {
  description = "Name of the claim token."
  type        = string
}

variable "claim_token_expiry_time" {
  description = "Expiry time of the claim token in RFC3339 format."
  type        = string
}

variable "claim_token_max_usage_count" {
  description = "Maximum number of times the claim token can be used for registering nodes."
  type        = number
}

variable "node_ext_id" {
  description = "External ID of a node registered with the claim token (for the singular node datasource)."
  type        = string
}
