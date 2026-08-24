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

# Prerequisite tenant + domain allocation that own the external network connection.
resource "nutanix_tenant_v2" "example" {
  name        = var.tenant_name
  description = "Tenant for external network connection example"

  admin_info {
    username   = "tenant-admin"
    first_name = "First"
    last_name  = "Last"
    email      = "tenant-admin@example.com"
    password   = var.tenant_admin_password
  }
}

resource "nutanix_tenant_domain_allocation_v2" "example" {
  tenant_ext_id = nutanix_tenant_v2.example.ext_id
  domain_ext_id = var.domain_ext_id
}

# Configure the external network connection in NAT mode. Only one of nat / no_nat
# should be set per configuration.
resource "nutanix_tenant_external_network_connection_v2" "example" {
  tenant_ext_id            = nutanix_tenant_v2.example.ext_id
  domain_allocation_ext_id = nutanix_tenant_domain_allocation_v2.example.ext_id

  config {
    nat {
      transit_vpc_ext_id    = var.transit_vpc_ext_id
      overlay_subnet_ext_id = var.overlay_subnet_ext_id
      floating_ip_count     = 2
    }
  }
}

# Fetch the external network connection settings for the domain allocation.
data "nutanix_tenant_external_network_connection_v2" "get-external-network-connection" {
  tenant_ext_id            = nutanix_tenant_v2.example.ext_id
  domain_allocation_ext_id = nutanix_tenant_domain_allocation_v2.example.ext_id
  depends_on               = [nutanix_tenant_external_network_connection_v2.example]
}
