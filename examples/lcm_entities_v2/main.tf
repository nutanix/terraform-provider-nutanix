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

# get all entities
data "nutanix_lcm_entities_v2" "lcm-entities" {}

# get filtered entities
data "nutanix_lcm_entities_v2" "lcm-entities-filtered" {
  filter = "entityModel eq 'Calm Policy Engine'"
}

# get limited entities
data "nutanix_lcm_entities_v2" "lcm-entities-limited" {
  limit = 5
}

# get filtered and limited entities
data "nutanix_lcm_entities_v2" "lcm-entities-filter-and-limit" {
  filter = "startswith(entityModel,'Calm')"
  limit  = 5
}

# get specific entity
data "nutanix_lcm_entity_v2" "entity" {
  ext_id = data.nutanix_lcm_entities_v2.lcm-entities.entities[0].ext_id
}
