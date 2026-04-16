#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type = string
}
variable "nutanix_password" {
  type = string
}
variable "nutanix_pc_endpoint" {
  type = string
}
variable "nutanix_pe_endpoint" {
  type = string
}
variable "nutanix_port" {
  type = string
}

variable "cluster_ext_id" {
  type = string
}
variable "bucket_name" {
  type = string
}
variable "region" {
  type = string
}
variable "access_key_id" {
  type = string
}
variable "secret_access_key" {
  type = string
}