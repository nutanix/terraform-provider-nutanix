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

variable "volume_group_name" {
  type = string
}

variable "volume_group_ext_id" {
  type = string
}

variable "volume_group_disk_ext_id" {
  type = string
}

variable "volume_group_sharing_status" {
  type = string
}

variable "volume_group_target_secret" {
  type = string
}

variable "disk_data_source_reference_ext_id" {
  type = string
}

variable "vg_iscsi_ext_id" {
  type = string
}

variable "vg_iscsi_initiator_name" {
  type = string
}

variable "vg_vm_ext_id" {
  type = string
}

variable "volume_iscsi_client_ext_id" {
  type = string
}