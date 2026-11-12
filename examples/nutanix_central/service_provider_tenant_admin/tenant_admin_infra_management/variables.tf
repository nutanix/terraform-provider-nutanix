# ==============================================================================
# Input Variables
# ==============================================================================


# ------------------------------------------------------------------------------
# NUTANIX CENTRAL CONNECTION
# ------------------------------------------------------------------------------

variable "nutanix_central_endpoint" {
  type        = string
  description = "Service Provider FQDN only (the part after '.tenants.' in the Tenant Admin URL). Example: nconprem-example.nxncm.com"
}

variable "tenant_username" {
  type        = string
  description = "Tenant Admin username."
}

variable "tenant_password" {
  type        = string
  sensitive   = true
  description = "Tenant Admin password."
}

variable "nutanix_port" {
  type        = number
  default     = 443
  description = "Port for Nutanix Central API."
}


# ------------------------------------------------------------------------------
# TENANT CONFIGURATION
# ------------------------------------------------------------------------------

variable "tenant_name" {
  type        = string
  description = "Tenant name from the Tenant Admin URL (the segment before '.tenants.')."
}

variable "domain_name" {
  type        = string
  description = "Name of the domain."
}



# ------------------------------------------------------------------------------
# VIRTUAL MACHINE CONFIGURATION (Optional)
# ------------------------------------------------------------------------------
variable "image_url" {
  type        = string
  description = "URL of the image to use for the VM."
}

variable "vm_name" {
  type        = string
  default     = "tenant-workload-vm"
  description = "Name for the tenant VM."
}

variable "vm_num_sockets" {
  type        = number
  default     = 2
  description = "Number of CPU sockets for the VM."
}

variable "vm_num_cores_per_socket" {
  type        = number
  default     = 2
  description = "Number of cores per CPU socket."
}

variable "vm_memory_gb" {
  type        = number
  default     = 8
  description = "Memory size in GB for the VM."
}

variable "vm_disk_gb" {
  type        = number
  default     = 100
  description = "Boot disk size in GB for the VM."
}
