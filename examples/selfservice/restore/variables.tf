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
variable "app_uuid" {
  type = string
}

variable "restore_action_name" {
  type = string
}

variable "snapshot_name" {
  type = string
}