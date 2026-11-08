terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.5.1"
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

# Create a tenant to demonstrate the data source
resource "nutanix_tenant_v2" "example_tenant" {
  name        = var.tenant_name
  description = "Tenant for demonstrating tenant memberships data source"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = var.tenant_admin_password
  }
}

# Create a domain for the tenant
resource "nutanix_registered_domain_v2" "example_domain" {
  provider = nutanix.nc
  name     = var.domain_name
  # Note: In a real scenario, you would need to properly configure the domain registration
  # This is simplified for demonstration purposes
}

# Create a role for the tenant (assuming we have a role to use)
# For simplicity, we'll use a data source to get an existing role
data "nutanix_roles_v2" "example_role" {
  # This would typically filter for a specific role name
  # For demonstration, we'll just get the first role available
}

# Lists all tenant memberships.
data "nutanix_tenant_memberships_v2" "example" {
  # Note: This data source doesn't require specific parameters for basic listing
  # It will list all memberships across all domains
}

# Output the list of memberships for verification
output "memberships" {
  value = data.nutanix_tenant_memberships_v2.example.memberships
}