variable "nutanix_username" {
  type        = string
  description = "Nutanix Prism Central username"
}

variable "nutanix_password" {
  type        = string
  description = "Nutanix Prism Central password"
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Nutanix Prism Central endpoint IP or FQDN"
}

variable "cluster_ext_id" {
  type        = string
  description = "The UUID of the target Prism Element cluster"
}

variable "host_ext_id" {
  type        = string
  description = "The UUID of the target AHV Host"
}

variable "host_nic" {
  type        = string
  description = "The physical NIC name on the host (e.g., eth0)"
  default     = "eth0"
}

variable "existing_bridge_name" {
  type        = string
  description = "The name of a pre-existing OVS bridge to migrate (e.g., br1)"
  default     = "br1"
}

variable "project_ext_id" {
  type        = string
  description = "The UUID of the project to share the virtual switch with"
}
