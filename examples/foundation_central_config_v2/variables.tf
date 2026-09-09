variable "nutanix_username" {
  type        = string
  description = "Username for the Nutanix Prism Central / Foundation Central endpoint."
}

variable "nutanix_password" {
  type        = string
  description = "Password for the Nutanix Prism Central / Foundation Central endpoint."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Prism Central endpoint (IP or FQDN)."
}

variable "foundation_endpoint" {
  type        = string
  description = "Foundation Central (FCVM) endpoint (IP or FQDN) that serves the Foundation Central config APIs."
}

variable "foundation_port" {
  type        = string
  description = "Port for the Foundation Central (FCVM) endpoint."
  default     = "9440"
}

variable "ahv_installation_timeout_minutes" {
  type        = number
  description = "Timeout in minutes for AHV installation."
  default     = 120
}

variable "aos_download_timeout_minutes" {
  type        = number
  description = "Timeout in minutes for AOS download."
  default     = 60
}
