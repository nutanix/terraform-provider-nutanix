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
variable "app_name" {
  type = string
}
variable "blueprint_name" {
  type = string
}
variable "app_description" {
  type = string
}
variable "patch_name" {
  type = string
}
variable "config_name" {
  type = string
}
variable "memory_size_mib" {
  type = number
}
variable "num_sockets" {
  type = number
}
variable "num_vcpus_per_socket" {
  type = number
}
variable "category_value" {
  type = string
}
variable "add_operation" {
  type = string
}
variable "delete_operation" {
  type = string
}
variable "disk_size_mib" {
  type = number
}
variable "index" {
  type = string
}
variable "subnet_uuid" {
  type = string
}