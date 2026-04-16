terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.1"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}
# Categories used by the entity group (matches test pattern: VM selected by CATEGORY_EXT_ID)
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_example_${count.index}_key"
  value       = "tf_entity_group_example_${count.index}_value"
  description = "Category for entity group example ${count.index}"
}

# Entity group with allowed_config using categories (matches TestAccNutanixEntityGroupV2Resource_Basic)
resource "nutanix_entity_group_v2" "basic" {
  name        = var.entity_group_name
  description = var.entity_group_description

  allowed_config {
    entities {
      type             = "VM"
      selected_by      = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        nutanix_category_v2.categories[0].id,
        nutanix_category_v2.categories[1].id
      ]
    }
  }
}

# Entity group with allowed_config (addresses and ip_ranges) - alternative style
# API requires type + selected_by per entity; (IP_VALUES, ADDRESS_GROUP). Only one entity per (selected_by, type) - combine addresses and ip_ranges in one block.
resource "nutanix_entity_group_v2" "with_allowed" {
  name        = "entity_group_with_allowed"
  description = "Entity group with allowed entities (addresses and ip_ranges)"

  allowed_config {
    entities {
      type        = "ADDRESS_GROUP"
      selected_by = "IP_VALUES"
      addresses {
        ipv4_addresses {
          value         = "10.0.0.0"
          prefix_length = 24
        }
      }
      ip_ranges {
        ipv4_ranges {
          start_ip = "192.168.1.1"
          end_ip   = "192.168.1.10"
        }
      }
    }
  }
}

# Get entity group by ext_id (matches TestAccNutanixEntityGroupV2Datasource_Basic)
data "nutanix_entity_group_v2" "by_id" {
  ext_id = nutanix_entity_group_v2.basic.id
}

# List entity groups (matches TestAccNutanixEntityGroupsV2Datasource_Basic)
data "nutanix_entity_groups_v2" "list" {
  depends_on = [nutanix_entity_group_v2.basic, nutanix_entity_group_v2.with_allowed]
}

# List entity groups with filter (matches TestAccNutanixEntityGroupsV2Datasource_WithFilter)
data "nutanix_entity_groups_v2" "filtered" {
  filter     = "name eq '${nutanix_entity_group_v2.basic.name}'"
  depends_on = [nutanix_entity_group_v2.basic]
}

# List entity groups with limit (matches TestAccNutanixEntityGroupsV2Datasource_WithLimit)
data "nutanix_entity_groups_v2" "with_limit" {
  limit      = 1
  depends_on = [nutanix_entity_group_v2.basic]
}
