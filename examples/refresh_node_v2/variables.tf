variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central endpoint."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central endpoint."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)."
  type        = string
}

variable "node_ext_id" {
  description = "External ID of the node to refresh."
  type        = string
}
