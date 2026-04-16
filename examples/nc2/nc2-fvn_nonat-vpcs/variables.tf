variable "NUTANIX_USERNAME" {
  description = "Username for Prism Central"
  type        = string
  sensitive   = true
}

variable "NUTANIX_PASSWORD" {
  description = "Password for Prism Central"
  type        = string
  sensitive   = true
}

variable "NUTANIX_ENDPOINT" {
  description = "Prism Central endpoint"
  type        = string
}

variable "NUTANIX_INSECURE" {
  description = "Whether to skip SSL verification"
  type        = bool
  default     = true
}

variable "NUTANIX_PORT" {
  description = "Prism Central port"
  type        = number
  default     = 9440
}

variable "NUTANIX_WAIT_TIMEOUT" {}


variable "VPC" {
  type = list(string)
}

variable "SUBNET_A" {
  type = list(string)
}

variable "SUBNET_B" {
  type = list(string)
}
