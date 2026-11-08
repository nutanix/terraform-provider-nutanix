# ==============================================================================
# Module Provider Requirements
# ==============================================================================
#
# This module expects three aliased nutanix provider configurations to be
# supplied by the calling module:
#   • nutanix.pc     - Prism Central admin  (PC discovery, AD, users, storage, VPCs, local domain)
#   • nutanix.nc_iam - Nutanix Central IAM  (NC directory service, LDAP user, authz policies)
#   • nutanix.nc     - Nutanix Central LDAP admin (domain/tenant operations: registered domain)
#
# ==============================================================================

terraform {
  required_providers {
    nutanix = {
      source                = "nutanix/nutanix"
      configuration_aliases = [nutanix.pc, nutanix.nc, nutanix.nc_iam]
    }
  }
}
