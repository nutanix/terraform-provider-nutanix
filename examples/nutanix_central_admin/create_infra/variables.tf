variable "nutanix_username" {
  type        = string
  description = "Username for the Nutanix Central Admin principal."
}

variable "nutanix_password" {
  type        = string
  description = "Password for the Nutanix Central Admin principal."
  sensitive   = true
}

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint base FQDN."
}

variable "domain_name" {
  type        = string
  description = "Domain name of the registered Prism Central to target."
}

variable "subnet_name" {
  type        = string
  description = "Name of an existing subnet on the registered domain."
}

variable "image_name" {
  type        = string
  description = "Name for the disk image created by the NC admin."
  default     = "tf-nc-admin-image"
}

variable "image_url" {
  type        = string
  description = "Source URL for the disk image."
}

variable "vm_name" {
  type        = string
  description = "Name for the virtual machine created on the registered domain."
  default     = "tf-nc-admin-vm"
}
