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

# Configure identity categorization for Flow Network Security ID-based security.
resource "nutanix_directory_server_config_v2" "example" {
  directory_service_reference           = var.directory_service_reference
  is_default_category_enabled           = true
  # is_default_category_enabled = true - Allowed only for CONTAINS match type

  should_keep_default_category_on_login = false
  # should_keep_default_category_on_login = true
  # Allowed only for CONTAINS match type
  # Allowed only when is_default_category_enabled = true

  # Determine which entities are eligible for dynamic categorization.
  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "CONTAINS"
    criteria     = "DeveloperVM"
    # Criteria should be the name of the VMs on which the category should be applied.
    # Criteria is allowed only for CONTAINS match type
    # Criteria is not allowed for ALL match type

    # If you want to apply the category to all the VMs, match_type should be ALL.
  }

  # Domain controllers used to scrape Windows logon events.
  domain_controllers {
    fqdn {
      value = "dc01.example.com"
    }
  }
}

# Read a single Directory Server configuration by ext_id.
data "nutanix_directory_server_config_v2" "get-config" {
  ext_id = nutanix_directory_server_config_v2.example.ext_id
}

# List all Directory Server configurations.
data "nutanix_directory_server_configs_v2" "configs-list" {}
