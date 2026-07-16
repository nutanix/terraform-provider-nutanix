#define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central account."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central account."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "IP address or FQDN of the Nutanix Prism Central."
  type        = string
}

variable "welcome_banner_content" {
  description = "Content of the welcome banner displayed on login."
  type        = string
}

variable "welcome_banner_is_enabled" {
  description = "Flag to denote whether the welcome banner is enabled or not."
  type        = bool
}
