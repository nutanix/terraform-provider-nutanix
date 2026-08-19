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
  description = "UUID of the cluster whose SNMP config is managed."
  type        = string
}
variable "snmp_is_enabled" {
  description = "Desired SNMP enabled flag on the cluster."
  type        = bool
  default     = true
}
variable "snmp_port" {
  description = "SNMP transport port."
  type        = number
  default     = 161
}
variable "snmp_protocol" {
  description = "SNMP transport protocol. One of TCP, TCP6, UDP, UDP6."
  type        = string
  default     = "UDP"
}
