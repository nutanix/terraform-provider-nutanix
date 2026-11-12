# ==============================================================================
# Input Variables
# ==============================================================================
variable "pc_endpoint" {
  type        = string
  description = "Prism Central endpoint (e.g., https://pc.nutanix.com)."
}

variable "pc_username" {
  type        = string
  description = "Prism Central username."
}

variable "pc_password" {
  type        = string
  sensitive   = true
  description = "Prism Central password."
}

variable "pc_name" {
  type        = string
  description = "Name of the Prism Central."
}

variable "registered_domain_name" {
  type        = string
  description = "Name of the registered domain."
}

# ------------------------------------------------------------------------------
# NUTANIX CENTRAL CONNECTION
# ------------------------------------------------------------------------------

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint (e.g., nconprem-xxx.nxncm.com)."
}

variable "nutanix_central_username" {
  type        = string
  description = "Nutanix Central LDAP username."
}

variable "nutanix_central_password" {
  type        = string
  sensitive   = true
  description = "Nutanix Central admin password."
}
