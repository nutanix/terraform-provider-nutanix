# ==============================================================================
# Input Variables
# ==============================================================================


# ------------------------------------------------------------------------------
# NUTANIX CENTRAL CONNECTION
# ------------------------------------------------------------------------------

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint (e.g., nconprem-xxx.nxncm.com)."
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
  description = "Name of the tenant."
}

variable "domain_name" {
  type        = string
  description = "Name of the domain."
}

# ------------------------------------------------------------------------------
# AD USER/GROUP SEARCH
# ------------------------------------------------------------------------------

variable "ad_user_query" {
  type        = string
  default     = "user_001"
  description = "Query string to search for user in Active Directory."
}

variable "ad_group_query" {
  type        = string
  default     = "user_group_1"
  description = "Query string to search for group in Active Directory."
}



# ------------------------------------------------------------------------------
# Directory Service Configuration (Optional)
# ------------------------------------------------------------------------------

variable "directory_service_name" {
  type        = string
  description = "User Group CRUD role external ID."
}

variable "directory_service_url" {
  type        = string
  description = "Directory service URL."
}

variable "directory_service_username" {
  type        = string
  description = "Directory service username."
}

variable "directory_service_password" {
  type        = string
  sensitive   = true
  description = "Directory service password."
}

variable "directory_service_domain_name" {
  type        = string
  description = "Directory service domain name."
}

variable "whitelisted_groups" {
  type        = list(string)
  default     = []
  description = "List of AD groups to whitelist for the directory service."
}