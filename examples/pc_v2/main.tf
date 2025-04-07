terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.1.0"
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

# list all PC (Domain Managers)
data "nutanix_pcs_v2" "list-pcs" {}

# list all PC (Domain Managers) with select query
data "nutanix_pcs_v2" "list-pcs-select" {
  select = "config"
}

# Fetch a single PC (Domain Manager) by its external ID
data "nutanix_pc_v2" "get-pc" {
  ext_id = "75dde184-3a0e-4f59-a185-03ca1efead17"
}
