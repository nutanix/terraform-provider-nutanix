# define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central API"
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central API"
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)"
  type        = string
}

variable "hardware_provider_ext_id" {
  description = "External ID of an existing hardware provider to scope connections to"
  type        = string
}

variable "connection_username" {
  description = "Username used for basic authentication to the hardware provider connection"
  type        = string
}

variable "connection_password" {
  description = "Password used for basic authentication to the hardware provider connection"
  type        = string
  sensitive   = true
}
