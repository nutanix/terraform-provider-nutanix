variable "ndb_endpoint" {
  type = string
}

variable "ndb_username" {
  type = string
}

variable "ndb_password" {
  type      = string
  sensitive = true
}

variable "pe_name" {
  type = string
}

variable "pe_cluster_ip" {
  type = string
}

variable "pe_username" {
  type = string
}

variable "pe_password" {
  type      = string
  sensitive = true
}

variable "pc_ip" {
  type    = string
  default = ""
}

variable "pc_name" {
  type    = string
  default = ""
}

variable "pc_username" {
  type    = string
  default = ""
}

variable "pc_password" {
  type      = string
  sensitive = true
  default   = ""
}

variable "storage_container" {
  type    = string
  default = ""
}

variable "dns_servers" {
  type = list(string)
}

variable "ntp_servers" {
  type = list(string)
}
