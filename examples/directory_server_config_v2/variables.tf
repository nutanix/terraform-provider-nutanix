#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "Prism Central username."
}

variable "nutanix_password" {
  type        = string
  description = "Prism Central password."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "Prism Central IP address or FQDN."
}

variable "directory_service_reference" {
  type        = string
  description = "The ExtID of the directory service that will be used for mapping."
}
