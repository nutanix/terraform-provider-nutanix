variable "nutanix_central_endpoint" {
  description = "Nutanix Central Endpoint"
  type        = string
}

variable "nutanix_password" {
  description = "Nutanix Password"
  type        = string
}

variable "nutanix_username" {
  description = "Nutanix Username"
  type        = string
}

variable "tenant_name" {
  description = "Tenant Name"
  type        = string
}

variable "tenant_admin_username" {
  description = "Tenant Admin Username"
  type        = string
}

variable "tenant_admin_password" {
  description = "Tenant Admin Password"
  type        = string
}

variable "domain_ext_id" {
  type = string
}

variable "cluster_ext_id" {
  type = string
}

variable "cpaas_storage_container_ext_id" {
  type = string
}

variable "database_storage_container_ext_id" {
  type = string
}

variable "self_service_storage_container_ext_id" {
  type = string
}

variable "nat_transit_vpc_ext_id" {
  type = string
}

variable "nat_overlay_subnet_ext_id" {
  type = string
}

variable "no_nat_transit_vpc_ext_id" {
  type = string
}

variable "no_nat_overlay_subnet_ext_id" {
  type = string
}

variable "externally_routable_prefixes" {
  type = list(string)
} 
