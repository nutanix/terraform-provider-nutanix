variable "nutanix_username" {
  type = string
}

variable "nutanix_password" {
  type      = string
  sensitive = true
}

variable "nutanix_endpoint" {
  type = string
}

variable "nutanix_port" {
  type    = number
  default = 9440
}

variable "cluster_name" {
  type    = string
  default = ""
}

variable "image_name" {
  type = string
}

variable "network_function_name" {
  type    = string
  default = "tf-network-function-inline"
}

variable "network_function_description" {
  type    = string
  default = "Inline network function managed by Terraform"
}

variable "management_subnet_name" {
  type    = string
  default = "tf-network-function-mgmt"
}

variable "management_subnet_vlan_id" {
  type = number
}

variable "management_subnet_network" {
  type = string
}

variable "management_subnet_prefix_length" {
  type    = number
  default = 24
}

variable "management_subnet_gateway" {
  type = string
}

variable "management_subnet_pool_start" {
  type = string
}

variable "management_subnet_pool_end" {
  type = string
}

variable "nf_vm_admin_password" {
  type      = string
  sensitive = true
}

variable "vm_num_sockets" {
  type    = number
  default = 2
}

variable "vm_num_cores_per_socket" {
  type    = number
  default = 2
}

variable "vm_memory_size_bytes" {
  type    = number
  default = 4294967296
}

variable "vm_disk_size_bytes" {
  type    = number
  default = 21474836480
}
