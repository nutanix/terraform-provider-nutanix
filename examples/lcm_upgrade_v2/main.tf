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

# List Prism Central
data "nutanix_clusters_v2" "pc" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  pcExtID      = data.nutanix_clusters_v2.pc.cluster_entities[0].ext_id
}

data "nutanix_lcm_entities_v2" "lcm-entities" {
  filter = "entityModel eq 'Calm Policy Engine'"
}

data "nutanix_lcm_entity_v2" "entity-before-upgrade" {
  ext_id = data.nutanix_lcm_entities_v2.lcm-entities.entities[0].ext_id
}

# perform inventory
resource "nutanix_lcm_perform_inventory_v2" "inventory" {
  x_cluster_id = local.pcExtID
  depends_on   = [data.nutanix_lcm_entity_v2.entity-before-upgrade]
}

resource "nutanix_lcm_prechecks_v2" "pre-checks" {
  x_cluster_id = local.pcExtID
  entity_update_specs {
    entity_uuid = data.nutanix_lcm_entity_v2.entity-before-upgrade.ext_id
    to_version  = "4.0.0"
  }
  depends_on = [nutanix_lcm_perform_inventory_v2.inventory]
}

# check if there is any operation in progress before starting the upgrade
data "nutanix_lcm_status_v2" "status-before-upgrade" {
  x_cluster_id = local.pcExtID
  lifecycle {
    postcondition {
      condition     = self.in_progress_operation[0].operation_type == "" && self.in_progress_operation[0].operation_id == ""
      error_message = "operation is in progress: ${self.in_progress_operation[0].operation_type}"
    }
  }
  depends_on = [nutanix_lcm_prechecks_v2.pre-checks]
}

# upgrade the entity
resource "nutanix_lcm_upgrade_v2" "upgrade" {
  entity_update_specs {
    entity_uuid = data.nutanix_lcm_entity_v2.entity-before-upgrade.ext_id
    to_version  = "4.0.0"
  }
  depends_on = [data.nutanix_lcm_status_v2.status-before-upgrade]
}

# check if there is any operation in progress after upgrade
data "nutanix_lcm_status_v2" "status-after-upgrade" {
  x_cluster_id = local.pcExtID
  lifecycle {
    postcondition {
      condition     = self.in_progress_operation[0].operation_type == "" && self.in_progress_operation[0].operation_id == ""
      error_message = "operation is in progress: ${self.in_progress_operation[0].operation_type}"
    }
  }
  depends_on = [nutanix_lcm_upgrade_v2.upgrade]
}

# fetch the entity after upgrade
data "nutanix_lcm_entity_v2" "entity-after-upgrade" {
  ext_id = data.nutanix_lcm_entities_v2.lcm-entities.entities[0].ext_id
  lifecycle {
    postcondition {
      condition     = self.ext_id == data.nutanix_lcm_entity_v2.entity-before-upgrade.ext_id
      error_message = "entity ext id changed"
    }
    postcondition {
      condition     = self.entity_version == "4.0.0"
      error_message = "entity version is not upgraded, current version: ${self.entity_version}"
    }
    postcondition {
      condition     = self.entity_model == "Calm Policy Engine"
      error_message = "entity model is changed, current model: ${self.entity_model}"
    }
  }
  depends_on = [nutanix_lcm_upgrade_v2.upgrade, data.nutanix_lcm_status_v2.status-after-upgrade]
}


