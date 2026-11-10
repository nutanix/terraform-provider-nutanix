#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "Username used to authenticate against Prism Central / Nutanix Central."
}

variable "nutanix_password" {
  type        = string
  description = "Password used to authenticate against Prism Central / Nutanix Central."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "IP address or FQDN of the Prism Central / Nutanix Central endpoint."
}
