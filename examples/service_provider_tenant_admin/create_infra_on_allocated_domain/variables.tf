variable "nutanix_username" {
  type        = string
  description = "Username for the Tenant Admin principal."
}

variable "nutanix_password" {
  type        = string
  description = "Password for the Tenant Admin principal."
  sensitive   = true
}

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint base FQDN."
}

variable "domain_name" {
  type        = string
  description = "Domain name of the Prism Central allocated to the tenant."
}

variable "tenant_name" {
  type        = string
  description = "Tenant name used when operating as Tenant Admin."
}

variable "subnet_name" {
  type        = string
  description = "Name of an existing VLAN subnet available in the allocated domain."
}

variable "image_name" {
  type        = string
  description = "Name for the disk image created by the tenant admin."
  default     = "tf-tenant-admin-image"
}

variable "image_url" {
  type        = string
  description = "Source URL for the disk image."
}

variable "vm_name" {
  type        = string
  description = "Name for the virtual machine created on the allocated domain."
  default     = "tf-tenant-admin-vm"
}
