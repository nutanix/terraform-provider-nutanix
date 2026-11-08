# ==============================================================================
# Module Outputs
# ==============================================================================

# ------------------------------------------------------------------------------
# PRISM CENTRAL IDENTIFIERS
# ------------------------------------------------------------------------------

output "pc_ext_id" {
  description = "External ID of the Prism Central instance."
  value       = local.pc_ext_id
}

output "pc_local_domain_id" {
  description = "ID of the completed PC local-domain registration."
  value       = nutanix_local_domain_v2.local_domain.id
}

# ------------------------------------------------------------------------------
# REGISTERED DOMAIN (Nutanix Central)
# ------------------------------------------------------------------------------

output "registered_domain_ext_id" {
  description = "External ID of the registered domain in Nutanix Central."
  value       = nutanix_registered_domain_v2.registered_domain.ext_id
}

output "registered_domain_dns_name" {
  description = "DNS label assigned to the registered domain."
  value       = nutanix_registered_domain_v2.registered_domain.dns_name
}

output "registration_config" {
  description = "NC-issued registration credentials and connection metadata. Contains API keys - handle securely."
  value       = data.nutanix_registered_domain_v2.registration_creds.registration_config
  sensitive   = true
}

# ------------------------------------------------------------------------------
# INFRASTRUCTURE IDs (for Tenant Domain Allocation)
# ------------------------------------------------------------------------------

output "cluster_ext_id" {
  description = "UUID of the cluster available for tenant domain allocation."
  value       = local.cluster_ext_id
}

output "cpaas_storage_container_ext_id" {
  description = "External ID of the storage container for tenant workloads."
  value       = local.cpaas_storage_container_ext_id
}

output "storage_container_self_service_ext_id" {
  description = "External ID of the storage container for tenant workloads."
  value       = local.storage_container_self_service_ext_id
}

output "database_storage_container_ext_id" {
  description = "External ID of the storage container for tenant workloads."
  value       = local.database_storage_container_ext_id
}

# ------------------------------------------------------------------------------
# NETWORK INFRASTRUCTURE - NAT
# ------------------------------------------------------------------------------

output "nat_transit_vpc_ext_id" {
  description = "External ID of the NAT-enabled transit VPC."
  value       = local.nat_transit_vpc_ext_id
}

output "nat_overlay_subnet_ext_id" {
  description = "External ID of the NAT overlay external subnet."
  value       = local.nat_overlay_subnet_ext_id
}

# ------------------------------------------------------------------------------
# NETWORK INFRASTRUCTURE - NO-NAT
# ------------------------------------------------------------------------------

output "no_nat_transit_vpc_ext_id" {
  description = "External ID of the No-NAT transit VPC."
  value       = local.no_nat_transit_vpc_ext_id
}

output "no_nat_overlay_subnet_ext_id" {
  description = "External ID of the No-NAT overlay external subnet."
  value       = local.no_nat_overlay_subnet_ext_id
}
