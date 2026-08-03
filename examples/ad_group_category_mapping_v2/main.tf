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

# Fetch the objectGUID for the User group. This is used to map the Active Directory group to Category for ID Based Security.
data "nutanix_directory_service_users_search_v2" "infra_team_group" {
  directory_service_ext_id = var.directory_service_reference
  query = "developer_group_infra_team" # This is the name of the Active Directory group that will be mapped to the Category.
  is_wildcard_search = true
}

locals {
  search_results_infra_team_group = data.nutanix_directory_service_users_search_v2.infra_team_group.search_results
  object_infra_team_group_guid = one(flatten([
    for result in local.search_results_infra_team_group : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "objectGUID"
    ]
  ]))
}

# Map an Active Directory group to a Prism Central category.
# AD Group Category Mapping can be done only when Directory Server Config is configured.
resource "nutanix_ad_group_category_mapping_v2" "example" {
  name           = "developer_group_infra_team" # This is the name of the Category Mapping.
  # Recommendation is to use the name of the Active Directory group as the name of the Category Mapping.

  category_value = "infra_team" # This is the value of the category that will be mapped to the Active Directory group.
  
  category_name = "ADGroup"
  # By default it is ADGroup.
  # Recommendation is to use the same key for all the Active Directory groups Category mappings. This is the key of the category that will be mapped to the Active Directory group.
  ad_info {
    directory_service_reference = var.directory_service_reference
    object_identifier           = local.object_infra_team_group_guid
  }
  # object_identifier is the objectGUID for the Active Directory group. Cannot be updated once created.
}

# Read a single Category Mapping by ext_id.
data "nutanix_ad_group_category_mapping_v2" "get-mapping" {
  ext_id = nutanix_ad_group_category_mapping_v2.example.ext_id
}

# List all Category Mappings, optionally filtered.
data "nutanix_ad_group_category_mappings_v2" "mappings-list" {
  limit  = 10
  filter = "name eq '${nutanix_ad_group_category_mapping_v2.example.name}'"
}
