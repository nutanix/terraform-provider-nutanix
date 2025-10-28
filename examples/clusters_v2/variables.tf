#define the type of variables to be used in terraform file
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
  type = string
}
variable "pe_username" {
  type = string
}
variable "pe_password" {
  type = string
}

variable "node_ip" {
  type = string
}

variable "username" {
  type = string
}

variable "password" {
  type = string
}

variable "nodes_ip" {
  type = list(string)
}
