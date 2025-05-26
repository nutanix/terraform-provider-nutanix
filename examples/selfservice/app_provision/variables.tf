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
variable "file_name" {
  type = string
}
variable "substrate_value" {
  type = string
  default = <<EOT
  <jq-formatted-value>
EOT 
}
variable "substrate_name" {
  type = string
}
variable "system_action_name" {
  type = string
}