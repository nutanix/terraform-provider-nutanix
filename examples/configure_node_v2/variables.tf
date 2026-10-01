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
  description = "External ID of the node to configure."
  type        = string
}

variable "group_id" {
  description = "Identifier of the group to which the node belongs (required for Cisco nodes)."
  type        = string
  default     = ""
}

variable "server_identity_pool_ext_id" {
  description = "UUID of the server identity pool used for server configuration."
  type        = string
  default     = ""
}

variable "cluster_ext_id" {
  description = "External ID of the cluster for post-cluster configuration."
  type        = string
  default     = ""
}
