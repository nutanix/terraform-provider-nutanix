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

variable "network_function_name" {
  type    = string
  default = "tf-network-function"
}

variable "ingress_nic_reference" {
  type = string
}

variable "egress_nic_reference" {
  type    = string
  default = ""
}

