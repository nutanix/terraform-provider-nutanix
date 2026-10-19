terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.5.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Create a tenant membership by referencing an existing principal (identity) via
# its external identifier. Use this variant when the principal has already been
# resolved to an identity ext_id in NC IAM.
resource "nutanix_tenant_membership_v2" "by_identity_ext_id" {
  domain_ext_id   = var.domain_ext_id
  role_ext_id     = var.role_ext_id
  identity_ext_id = var.identity_ext_id
}

# Create a tenant membership by resolving an external (AD/LDAP/SAML) user via the
# identity provider. Use this variant when you only know the principal's lookup
# value (UPN / username / email) and its identity provider.
resource "nutanix_tenant_membership_v2" "by_identity_details" {
  domain_ext_id = var.domain_ext_id
  role_ext_id   = var.role_ext_id

  identity_details {
    identity_provider_ext_id = var.identity_provider_ext_id
    principal_value          = var.principal_value
    type                     = var.identity_type
  }
}

# Read a single tenant membership by its ext_id.
data "nutanix_tenant_membership_v2" "get_membership" {
  ext_id = nutanix_tenant_membership_v2.by_identity_ext_id.ext_id
}

# List all tenant memberships for the tenant across all domains.
data "nutanix_tenant_memberships_v2" "list_memberships" {
  depends_on = [nutanix_tenant_membership_v2.by_identity_ext_id]
}

# List with filter/limit — filter memberships by their domain.
data "nutanix_tenant_memberships_v2" "filtered_memberships" {
  limit  = 10
  filter = "domainExtId eq '${var.domain_ext_id}'"
}
