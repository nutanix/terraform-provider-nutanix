# Define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "Username for Nutanix Prism Central"
}

variable "nutanix_password" {
  type        = string
  description = "Password for Nutanix Prism Central"
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Endpoint for Nutanix Prism Central (IP address or FQDN)"
}

variable "nutanix_port" {
  type        = number
  description = "Port for Nutanix Prism Central"
  default     = 9440
}
