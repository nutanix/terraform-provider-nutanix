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

variable "entity_group_name" {
  type        = string
  description = "Name for the entity group (matches test pattern)"
  default     = "entity_group_basic"
}

variable "entity_group_description" {
  type        = string
  description = "Description for the entity group"
  default     = "Entity group example with categories"
}
