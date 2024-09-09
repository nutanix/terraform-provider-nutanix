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


