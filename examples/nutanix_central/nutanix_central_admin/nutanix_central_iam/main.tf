# Nutanix Central IAM
# --------------------
# To leverage Nutanix Central IAM features (authorization policies, directory services, etc.),
# set `use_nc_iam = true` in the provider configuration.
# ==============================================================================
# PROVIDERS
# ==============================================================================

# Nutanix Central IAM
provider "nutanix" {
  nutanix_central_endpoint = var.nutanix_central_endpoint
  username                 = var.nutanix_central_username
  password                 = var.nutanix_central_password
  use_nc_iam               = true
  port                     = 443
}

# ==============================================================================
# in order to create auth polices or directory services ..etc, use the corresponding resources or data sources.