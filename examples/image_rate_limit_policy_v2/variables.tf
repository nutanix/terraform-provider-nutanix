variable "nutanix_username" {
  type        = string
  description = "The username for Nutanix Prism Central"
}

variable "nutanix_password" {
  type        = string
  sensitive   = true
  description = "The password for Nutanix Prism Central"
}

variable "nutanix_endpoint" {
  type        = string
  description = "The endpoint (IP address) of Nutanix Prism Central"
}
