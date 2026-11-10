terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
    }
  }
}

# Service Provider Tenant Admin — manage memberships for users/groups in an
# allocated domain. IAM calls are routed to NC IAM via use_nc_iam.
#
# Note (KB-22890): Calls to Nutanix Central may return 404 when the SP FQDN Host
# header includes an explicit :443 port. This can affect Terraform and other API
# clients. Until a fix ships, add both <host> and <host>:443 to
# ALLOWED_CALLBACK_HOSTS on the Dion service and restart it.
# Details: https://portal.nutanix.com/page/documents/kbs/details?targetId=kA0VO000000EVaD0AW
provider "nutanix" {
  username                 = var.nutanix_username
  password                 = var.nutanix_password
  nutanix_central_endpoint = var.nutanix_central_endpoint
  domain_name              = var.domain_name
  tenant_name              = var.tenant_name
  use_nc_iam               = true
  port                     = 443
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

output "memberships" {
  value = data.nutanix_tenant_memberships_v2.list_memberships.memberships
}
