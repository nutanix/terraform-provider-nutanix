variable "nutanix_username" {
  type = string
}

variable "nutanix_password" {
  type      = string
  sensitive = true
}

variable "nutanix_endpoint" {
  type = string
}

variable "directory_service_ext_id" {
  type        = string
  description = "The ExtID of the directory service to use for mapping."
}
