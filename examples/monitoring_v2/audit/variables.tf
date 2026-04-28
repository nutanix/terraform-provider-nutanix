variable "nutanix_username" {
  description = "Nutanix username"
  type        = string
  default     = ""
}

variable "nutanix_password" {
  description = "Nutanix password"
  type        = string
  sensitive   = true
  default     = ""
}

variable "nutanix_endpoint" {
  description = "Nutanix Prism Central endpoint"
  type        = string
  default     = ""
}

variable "nutanix_port" {
  description = "Nutanix Prism Central port"
  type        = string
  default     = "9440"
}

variable "nutanix_insecure" {
  description = "Allow insecure SSL connection"
  type        = bool
  default     = true
}
