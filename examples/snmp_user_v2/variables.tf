variable "nutanix_username" {
  type = string
}
variable "nutanix_password" {
  type = string
}
variable "nutanix_endpoint" {
  type = string
}
variable "cluster_ext_id" {
  description = "UUID of the cluster on which the SNMP user is created."
  type        = string
}
