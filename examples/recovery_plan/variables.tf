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
variable "source_cluster_uuid" {
  type = string
}
variable "target_cluster_uuid" {
  type = string
}
variable "source_az_url" {
  type = string
}
variable "target_az_url" {
  type = string
}
