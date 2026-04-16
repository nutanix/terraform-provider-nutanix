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

variable "user" {
  type = string
}
variable "password" {
  type = string
}
variable "endpoint" {
  type = string
}
variable "insecure" {
  type = bool
}
variable "port" {
  type = number
}

variable "operations" {
  type = list(string)
  default = [
    "operation_1_ext_id",
    "operation_2_ext_id",
    "operation_3_ext_id",
  ]
}


