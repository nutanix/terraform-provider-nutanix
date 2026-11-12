# ==============================================================================
# Tenant Admin: Tenant Membership Workflow
# ==============================================================================
#
# What this file does:
#   Shows how a Tenant Admin grants users and groups access to an allocated
#   domain by creating tenant memberships (role assignments via NC IAM).
#
# Flow:
#   1. Authenticate as Tenant Admin via nutanix_central_endpoint + tenant_name
#      + domain_name (see provider comments below for how to map the Tenant URL).
#   2. Create (or reuse) an Active Directory directory service and whitelist groups.
#   3. Discover the allocated domain and the "Tenant Consumer" role.
#   4. Search AD for a user and a group (UPN / distinguishedName + objectGUID).
#   5. Create tenant memberships for that user and group on the domain.
#   6. Read back single and list membership data sources for verification.
#
# Prerequisites:
#   1. SP Admin has created the tenant and allocated a domain
#   2. Tenant Admin credentials are available
#   3. AD/LDAP reachable with a service account and searchable users/groups
#
# ==============================================================================

# ==============================================================================
# TERRAFORM CONFIGURATION
# ==============================================================================

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.5.1"
    }
  }
}


# ==============================================================================
# PROVIDER CONFIGURATION FOR TENANT ADMIN
# ==============================================================================
#
# How to set these values from the Tenant Admin URL you receive from the SP Admin:
#
#   https://<tenant_name>.tenants.<service_provider_fqdn>
#
# Example:
#   https://my-tenant.tenants.nconprem-example.nxncm.com
#
#   - tenant_name              = "my-tenant"
#   - nutanix_central_endpoint = "nconprem-example.nxncm.com"   # the service provider FQDN only
#   - domain_name              = allocated PC domain name (from the SP Admin domain allocation)
#
# Do not put the full tenant URL into nutanix_central_endpoint — use only the
# service provider FQDN (the part after ".tenants.").
provider "nutanix" {
  username                 = var.tenant_username
  password                 = var.tenant_password
  nutanix_central_endpoint = var.nutanix_central_endpoint # Tenant Endpoint, Shared to tenant Admin by Service Provider Admin
  domain_name              = var.domain_name              # Domain name for which tenant Admin need to create the infrastructure in
  tenant_name              = var.tenant_name              # Tenant Name. The Tenant name will be shared by Service Provider Admin
  port                     = 443
}


# ==============================================================================
# Add another or Fetch DIRECTORY SERVICE (Active Directory)
# ==============================================================================

resource "nutanix_directory_services_v2" "ad" {
  name           = var.directory_service_name
  url            = var.directory_service_url
  directory_type = "ACTIVE_DIRECTORY"
  domain_name    = var.directory_service_domain_name

  service_account {
    username = var.directory_service_username
    password = var.directory_service_password
  }

  white_listed_groups = var.whitelisted_groups

  lifecycle {
    ignore_changes = [service_account[0].password]
  }
}

data "nutanix_directory_services_v2" "ad" {
  filter = "name eq '${var.directory_service_name}'"
  depends_on = [ nutanix_directory_services_v2.ad ]
}

# ==============================================================================
# DATA SOURCES
# ==============================================================================

# Get registered domains
data "nutanix_registered_domains_v2" "domains" {}

# Get roles for tenant consumer
data "nutanix_roles_v2" "roles" {
  filter = "contains(displayName, 'Tenant Consumer')"
}

# Search for user in AD
data "nutanix_directory_service_users_search_v2" "user" {
  directory_service_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
  query                    = var.ad_user_query
  is_wildcard_search       = false
}

# Search for group in AD
data "nutanix_directory_service_users_search_v2" "group" {
  directory_service_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
  query                    = var.ad_group_query
  is_wildcard_search       = false
}

# ==============================================================================
# LOCAL VALUES
# ==============================================================================

locals {
  domain_ext_id = data.nutanix_registered_domains_v2.domains.registered_domains[0].domain_ext_id

  # Find Tenant Consumer role
  tenant_admin_role_ext_id = data.nutanix_roles_v2.roles.roles[0].ext_id

  # Extract user attributes from AD search
  search_results_user = data.nutanix_directory_service_users_search_v2.user.search_results
  user_principal_name = one(flatten([
    for result in local.search_results_user : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "userPrincipalName"
    ]
  ]))
  user_identity_ext_id = one(flatten([
    for result in local.search_results_user : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "objectGUID"
    ]
  ]))

  # Extract group attributes from AD search
  search_results_group = data.nutanix_directory_service_users_search_v2.group.search_results
  group_principal_name = one(flatten([
    for result in local.search_results_group : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "distinguishedName"
    ]
  ]))
  group_identity_ext_id = one(flatten([
    for result in local.search_results_group : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "objectGUID"
    ]
  ]))
}


# ==============================================================================
# TENANT MEMBERSHIPS
# ==============================================================================

# Add user as Tenant Consumer
resource "nutanix_tenant_membership_v2" "user" {
  domain_ext_id = local.domain_ext_id
  role_ext_id   = local.tenant_admin_role_ext_id 

  identity_details {
    identity_provider_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
    principal_value          = local.user_principal_name
    type                     = "USER"
  }
}

# Add group as Tenant Consumer
resource "nutanix_tenant_membership_v2" "group" {
  domain_ext_id = local.domain_ext_id
  role_ext_id   = local.tenant_admin_role_ext_id

  identity_details {
    identity_provider_ext_id = data.nutanix_directory_services_v2.ad.directory_services[0].ext_id
    principal_value          = local.group_principal_name
    type                     = "GROUP"
  }
}

# ==============================================================================
# DATA SOURCES - Membership Validation
# ==============================================================================

data "nutanix_tenant_membership_v2" "user" {
  ext_id = nutanix_tenant_membership_v2.user.ext_id
}

data "nutanix_tenant_membership_v2" "group" {
  ext_id = nutanix_tenant_membership_v2.group.ext_id
}

data "nutanix_tenant_memberships_v2" "all" {
  filter = "domainExtId eq '${local.domain_ext_id}'"
  depends_on = [
    nutanix_tenant_membership_v2.user,
    nutanix_tenant_membership_v2.group,
  ]
}