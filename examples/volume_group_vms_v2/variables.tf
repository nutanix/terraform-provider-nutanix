#define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Prism Central username."
  type        = string
}

variable "nutanix_password" {
  description = "Prism Central password."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)."
  type        = string
}

variable "volume_group_name" {
  description = "Name of the volume group to create."
  type        = string
}

variable "vm_name" {
  description = "Name of the virtual machine to create and attach."
  type        = string
}
