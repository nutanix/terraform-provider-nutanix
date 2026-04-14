#define the type of variables to be used in terraform file

# Provider Credentials (nutanix_username / nutanix_password): 
# These are used by the Terraform provider to authenticate 
# with your Prism Central (PC) to execute the plan.
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
# PE Credentials (pe_username / pe_password):
# These are the credentials of the PE to be used for the ssh to the node.
variable "pe_username" {
  type = string
}
variable "pe_password" {
  type = string
}

variable "node_ip" {
  type = string
}

# Registration Credentials (username / password): 
# These are the new credentials you want to set for the PE. 
# The local-exec command sets these on the PE, and the nutanix_pc_registration_v2 resource then uses them to complete the registration. 
variable "username" {
  type = string
}

variable "password" {
  type = string
}

variable "nodes_ip" {
  type = list(string)
}
