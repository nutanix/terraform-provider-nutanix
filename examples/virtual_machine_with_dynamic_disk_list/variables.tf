
#variable definations
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
variable "cluster_name" {
  type = string
}
variable "subnet_name" {
  type = string
}
variable "disk_sizes" {
  type = list(string)
  default = [1024,2048]
}