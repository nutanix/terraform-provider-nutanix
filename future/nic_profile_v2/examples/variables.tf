variable "nutanix_username" {
  type        = string
  description = "Username for the Nutanix Prism Central account."
}

variable "nutanix_password" {
  type        = string
  description = "Password for the Nutanix Prism Central account."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Prism Central IP address or FQDN."
}

variable "nutanix_port" {
  type        = number
  description = "Prism Central port."
  default     = 9440
}

variable "nic_family" {
  type        = string
  description = "Host NIC device family in the format vendor_id:device_id (e.g. 15b3:101d for an NVIDIA/Mellanox ConnectX device)."
  default     = "15b3:101d"
}

variable "host_nic_ext_ids" {
  type        = list(string)
  description = "List of Host NIC UUIDs to associate with the SR-IOV NIC profile. Leave empty to create the profile without associating Host NICs."
  default     = []
}
