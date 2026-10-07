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
  description = "UUID of the cluster the SNMP trap belongs to."
  type        = string
}
variable "snmp_trap_v2_ipv4" {
  description = "IPv4 address of the V2 SNMP trap receiver."
  type        = string
}
variable "snmp_trap_v2_port" {
  description = "Port for the V2 SNMP trap."
  type        = number
  default     = 162
}
variable "snmp_trap_v3_ipv4" {
  description = "IPv4 address of the V3 SNMP trap receiver."
  type        = string
}
variable "snmp_trap_v3_port" {
  description = "Port for the V3 SNMP trap."
  type        = number
  default     = 163
}
variable "snmp_v3_username" {
  description = "Username for the SNMP V3 user that the V3 trap will reference."
  type        = string
  default     = "tf-snmp-v3-user"
}
