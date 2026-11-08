# ==============================================================================
# Module Input Variables
# ==============================================================================

# ------------------------------------------------------------------------------
# DOMAIN / PRISM CENTRAL IDENTITY
# ------------------------------------------------------------------------------

variable "pc_name" {
  description = "Name of the Prism Central instance to discover (matched case-insensitively against nutanix_pcs_v2)."
  type        = string
}

variable "registered_domain_name" {
  description = "Name to assign to the domain registered in Nutanix Central."
  type        = string
}

# ------------------------------------------------------------------------------
# ACTIVE DIRECTORY CONFIGURATION
# ------------------------------------------------------------------------------

variable "pc_directory_service" {
  description = "Active Directory configuration for the Prism Central directory service."
  type = object({
    name                     = string
    url                      = string
    domain_name              = string
    service_account_username = string
    service_account_password = string
  })
}

variable "nc_directory_service" {
  description = "Active Directory configuration for the Nutanix Central IAM directory service."
  type = object({
    name                     = string
    url                      = string
    domain_name              = string
    service_account_username = string
    service_account_password = string
  })
}

# ------------------------------------------------------------------------------
# LDAP USERS
# ------------------------------------------------------------------------------

variable "pc_ldap_username" {
  description = "LDAP username to create as a Prism Admin on Prism Central."
  type        = string
}

variable "nc_ldap_username" {
  description = "LDAP username to create as an admin on Nutanix Central IAM."
  type        = string
}

# ------------------------------------------------------------------------------
# ROLE NAMES (display names resolved against the roles_v2 data sources)
# ------------------------------------------------------------------------------

variable "pc_admin_role_name" {
  description = "Display name of the Prism Central role granted to the PC LDAP user."
  type        = string
  default     = "Super Admin"
}

variable "nc_admin_role_name" {
  description = "Display name of the Nutanix Central role granted to the NC LDAP user."
  type        = string
  default     = "Nutanix Central Admin"
}

variable "service_provider_admin_role_name" {
  description = "Display name of the Service Provider Admin role granted to the NC LDAP user."
  type        = string
  default     = "Service Provider Admin"
}
