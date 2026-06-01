terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">=2.0"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  insecure = true
}

# Create a Directory Server Config
resource "nutanix_directory_server_config_v2" "example" {
  directory_service_reference = var.directory_service_ext_id

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "ALL"
  }

  is_default_category_enabled           = true
  should_keep_default_category_on_login = false
}

# Read a single Directory Server Config
data "nutanix_directory_server_config_v2" "example" {
  ext_id = nutanix_directory_server_config_v2.example.ext_id
}

# List all Directory Server Configs
data "nutanix_directory_server_configs_v2" "all" {
  depends_on = [nutanix_directory_server_config_v2.example]
}

# List all Category Mappings
data "nutanix_category_mappings_v2" "all" {}
