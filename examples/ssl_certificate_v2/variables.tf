# Variable definitions for SSL certificate example

variable "nutanix_username" {
  type        = string
  description = "Nutanix Prism Central username"
}

variable "nutanix_password" {
  type        = string
  description = "Nutanix Prism Central password"
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Nutanix Prism Central endpoint (IP or FQDN)"
}

variable "nutanix_port" {
  type        = string
  description = "Nutanix Prism Central port"
  default     = "9440"
}

variable "passphrase" {
  type        = string
  description = "Passphrase for the SSL certificate (optional)"
  sensitive   = true
  default     = ""
}

variable "private_key" {
  type        = string
  description = "Private key for the SSL certificate in PEM format"
  sensitive   = true
  default     = ""
}

variable "public_certificate" {
  type        = string
  description = "Public certificate for the SSL certificate in PEM format"
  default     = ""
}

variable "ca_chain" {
  type        = string
  description = "Certificate authority (CA) chain in PEM format"
  default     = ""
}

