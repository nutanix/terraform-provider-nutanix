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
  description = "External ID of the node to image."
  type        = string
}

variable "patched_image_ext_id" {
  description = "External ID of a patched hypervisor or host OS image."
  type        = string
  default     = ""
}

variable "aos_image_ext_id" {
  description = "External ID of the AOS image managed by Foundation Central."
  type        = string
  default     = ""
}
