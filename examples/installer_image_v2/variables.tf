#define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username used to authenticate with Prism Central."
  type        = string
}

variable "nutanix_password" {
  description = "Password used to authenticate with Prism Central."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)."
  type        = string
}

variable "image_name" {
  description = "Name of the installer image to register."
  type        = string
}

variable "image_type" {
  description = "Type of the installer image. One of AOS, AHV or ESX."
  type        = string
}

variable "image_url" {
  description = "URL from where the installer image can be downloaded."
  type        = string
}

variable "image_version" {
  description = "Version of the installer image."
  type        = string
}
