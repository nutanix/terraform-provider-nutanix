variable "nutanix_username" {
  type = string
}

variable "nutanix_password" {
  type = string
}

variable "nutanix_endpoint" {
  type = string
}

variable "nutanix_port" {
  type    = string
  default = "9440"
}

variable "role_ext_id" {
  type        = string
  description = "External identifier of the role to assign."
}

variable "identity_ext_id" {
  type        = string
  description = "External identifier of the identity (user/group) to assign the role to."
}

variable "idp_ext_id" {
  type        = string
  description = "External identifier of the identity provider."
  default     = ""
}

variable "scope_template_name" {
  type        = string
  description = "Name of the scope template."
  default     = ""
}

variable "project_ext_id" {
  type        = string
  description = "External identifier of the project."
  default     = ""
}
