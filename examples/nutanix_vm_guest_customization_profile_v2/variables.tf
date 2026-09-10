# =============================================================================
# Provider credentials
# =============================================================================

variable "nutanix_username" {
  type        = string
  description = "Prism Central username"
}

variable "nutanix_password" {
  type        = string
  sensitive   = true
  description = "Prism Central password"
}

variable "nutanix_endpoint" {
  type        = string
  description = "Prism Central IP or hostname"
}

# =============================================================================
# Infrastructure
# =============================================================================

variable "subnet_name" {
  type        = string
  description = "Name of the subnet to look up via data source"
  default     = "VLAN 225"
}

variable "source_vm_uuid" {
  type        = string
  description = "Ext ID of a pre-created Windows VM to use as the template source and clone source"
}

# =============================================================================
# GC Profile — Administrator password
# =============================================================================

variable "admin_password" {
  type        = string
  sensitive   = true
  description = "Administrator password set in GC profiles"
}

# =============================================================================
# Domain settings (used in GC profiles)
# =============================================================================

variable "domain_name" {
  type        = string
  description = "Active Directory domain FQDN (e.g., qa.nutanix.com)"
}

variable "domain_username" {
  type        = string
  description = "Domain join username (e.g., Administrator@qa.nutanix.com)"
}

variable "domain_password" {
  type        = string
  sensitive   = true
  description = "Domain join password"
}

variable "domain_dns_server" {
  type        = string
  description = "Preferred DNS server IP for the domain"
}

# =============================================================================
# Override values (used in template deploy and clone overrides)
# =============================================================================

variable "override_admin_password" {
  type        = string
  sensitive   = true
  description = "Administrator password set in override (different from profile)"
}

variable "override_domain_name" {
  type        = string
  description = "Domain FQDN used in override (can be different from profile domain)"
  default     = ""
}

variable "override_domain_username" {
  type        = string
  description = "Domain join username for override"
  default     = ""
}

variable "override_domain_password" {
  type        = string
  sensitive   = true
  description = "Domain join password for override"
  default     = ""
}

variable "override_dns_server" {
  type        = string
  description = "DNS server IP used in override"
  default     = ""
}

# =============================================================================
# Networking — Static IPs for deploy/clone
# =============================================================================

variable "deploy_ip_address" {
  type        = string
  description = "Static IP address for the deployed VM"
}

variable "clone_ip_address" {
  type        = string
  description = "Static IP address for the cloned VM"
}

variable "subnet_prefix_length" {
  type        = number
  description = "Subnet prefix length (e.g., 27)"
  default     = 27
}

variable "subnet_gateway" {
  type        = string
  description = "Default gateway IP for the subnet"
}

# =============================================================================
# Windows product key (optional)
# =============================================================================

variable "windows_product_key" {
  type        = string
  sensitive   = true
  description = "Windows product key for activation"
  default     = ""
}
