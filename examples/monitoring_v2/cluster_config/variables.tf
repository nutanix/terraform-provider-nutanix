variable "nutanix_username" {
  type = string
}

variable "nutanix_password" {
  type      = string
  sensitive = true
}

variable "nutanix_endpoint" {
  type = string
}

variable "nutanix_port" {
  type    = string
  default = "9440"
}

variable "system_defined_policy_ext_id" {
  type        = string
  description = "Unique ID of the System-Defined Alert Policy."
}

variable "cluster_ext_id" {
  type        = string
  description = "Cluster UUID."
}
