variable "nutanix_username" {
  type        = string
  description = "Username for Nutanix Central."
}

variable "nutanix_password" {
  type        = string
  description = "Password for Nutanix Central."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Nutanix Central endpoint."
}