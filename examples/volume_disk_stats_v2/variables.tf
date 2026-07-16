#define the type of variables to be used in terraform file

variable "nutanix_username" {
  description = "Prism Central username."
  type        = string
}

variable "nutanix_password" {
  description = "Prism Central password."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)."
  type        = string
}

variable "volume_group_ext_id" {
  description = "The external identifier of the Volume Group that owns the disk."
  type        = string
}

variable "volume_disk_ext_id" {
  description = "The external identifier of the Volume Disk whose stats are queried."
  type        = string
}

variable "start_time" {
  description = "The start time (RFC-3339) from which the stats should be reported."
  type        = string
}

variable "end_time" {
  description = "The end time (RFC-3339) until which the stats should be reported."
  type        = string
}

variable "sampling_interval" {
  description = "The sampling interval in seconds at which the stats should be reported."
  type        = number
  default     = 30
}

variable "stat_type" {
  description = "The down-sampling operator. One of SUM, MIN, MAX, AVG, COUNT, LAST."
  type        = string
  default     = "AVG"
}
