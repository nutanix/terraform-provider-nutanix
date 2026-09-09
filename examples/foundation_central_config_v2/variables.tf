#define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central API."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central API."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)."
  type        = string
}

variable "ahv_installation_timeout_minutes" {
  description = "Timeout in minutes for AHV installation."
  type        = number
  default     = 60
}

variable "aos_download_timeout_minutes" {
  description = "Timeout in minutes for AOS download."
  type        = number
  default     = 60
}
