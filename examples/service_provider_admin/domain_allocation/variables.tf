variable "nutanix_central_username" {
  type        = string
  description = "Username for Nutanix Central (Service Provider Admin)."
}

variable "nutanix_central_password" {
  type        = string
  description = "Password for Nutanix Central."
  sensitive   = true
}

variable "nutanix_central_endpoint" {
  type        = string
  description = "Nutanix Central endpoint base FQDN."
}

variable "tenant_name" {
  type        = string
  description = "The name of the prerequisite tenant."
}

variable "tenant_admin_password" {
  type        = string
  description = "The password for the tenant local admin user."
  sensitive   = true
}

variable "domain_ext_id" {
  type        = string
  description = "External identifier of the domain to allocate to the tenant."
}

variable "transit_vpc_ext_id" {
  type        = string
  description = "External identifier of the transit VPC for the allocation's NAT configuration."
}

variable "overlay_subnet_ext_id" {
  type        = string
  description = "External identifier of the overlay subnet for the allocation's NAT configuration."
}

variable "cluster_uuid" {
  type        = string
  description = "The UUID of the cluster used as a placement target."
}

variable "storage_container_uuid" {
  type        = string
  description = "The UUID of the storage container used as a placement target."
}
