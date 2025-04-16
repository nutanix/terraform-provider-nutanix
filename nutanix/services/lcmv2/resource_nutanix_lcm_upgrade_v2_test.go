package lcmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameLcmUpgrade = "nutanix_lcm_upgrade_v2.upgrade"

func TestAccV2NutanixLcmUpgrade_Basic(t *testing.T) {
	datasourceNameLcmEntities := "data.nutanix_lcm_entities_v2.lcm-entities"
	datasourceNameLcmEntityBeforeUpgrade := "data.nutanix_lcm_entity_v2.entity-before-upgrade"
	datasourceNameLcmStatusBeforeUpgrade := "data.nutanix_lcm_status_v2.status-before-upgrade"
	resourceNameLcmPerformInventory := "nutanix_lcm_perform_inventory_v2.inventory"
	resourceNameLcmPreChecks := "nutanix_lcm_prechecks_v2.pre-checks"
	datasourceNameLcmStatusAfterUpgrade := "data.nutanix_lcm_status_v2.status-after-upgrade"
	datasourceNameLcmEntityAfterUpgrade := "data.nutanix_lcm_entity_v2.entity-after-upgrade"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		// CheckDestroy: testCheckDestroyProtectedResource,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testLcmUpgradeConfig(),
				Check: resource.ComposeTestCheckFunc(
					// check if the entity model is present
					resource.TestCheckResourceAttr(datasourceNameLcmEntities, "entities.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameLcmEntities, "entities.0.entity_model", testVars.Lcm.EntityModel),
					// check if the entity is present
					resource.TestCheckResourceAttrSet(datasourceNameLcmEntityBeforeUpgrade, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameLcmEntityBeforeUpgrade, "entity_model", testVars.Lcm.EntityModel),
					resource.TestCheckResourceAttrSet(datasourceNameLcmEntityBeforeUpgrade, "entity_version"),
					// check if there is any operation in progress before starting the upgrade
					resource.TestCheckResourceAttr(datasourceNameLcmStatusBeforeUpgrade, "in_progress_operation.0.operation_type", ""),
					resource.TestCheckResourceAttr(datasourceNameLcmStatusBeforeUpgrade, "in_progress_operation.0.operation_id", ""),
					// perform inventory checks
					resource.TestCheckResourceAttrSet(resourceNameLcmPerformInventory, "x_cluster_id"),
					// pre-checks checks
					resource.TestCheckResourceAttrSet(resourceNameLcmPreChecks, "x_cluster_id"),
					resource.TestCheckResourceAttr(resourceNameLcmPreChecks, "entity_update_specs.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameLcmPreChecks, "entity_update_specs.0.entity_uuid", datasourceNameLcmEntityBeforeUpgrade, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameLcmPreChecks, "entity_update_specs.0.to_version", testVars.Lcm.EntityModelVersion),
					// upgrade checks
					resource.TestCheckResourceAttr(resourceNameLcmUpgrade, "entity_update_specs.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameLcmUpgrade, "entity_update_specs.0.entity_uuid", datasourceNameLcmEntityBeforeUpgrade, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameLcmUpgrade, "entity_update_specs.0.to_version", testVars.Lcm.EntityModelVersion),
					// lcm status after upgrade
					resource.TestCheckResourceAttr(datasourceNameLcmStatusAfterUpgrade, "in_progress_operation.0.operation_type", ""),
					resource.TestCheckResourceAttr(datasourceNameLcmStatusAfterUpgrade, "in_progress_operation.0.operation_id", ""),
					// entity after upgrade
					resource.TestCheckResourceAttrSet(datasourceNameLcmEntityAfterUpgrade, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameLcmEntityAfterUpgrade, "entity_version", testVars.Lcm.EntityModelVersion),
					resource.TestCheckResourceAttr(datasourceNameLcmEntityAfterUpgrade, "entity_model", testVars.Lcm.EntityModel),
				),
			},
		},
	})
}

func testLcmUpgradeConfig() string {
	return fmt.Sprintf(`

# List Prism Central
data "nutanix_clusters_v2" "pc" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  pcExtID      = data.nutanix_clusters_v2.pc.cluster_entities[0].ext_id
  config = jsondecode(file("%[1]s"))
  lcm          = local.config.lcm
}

data "nutanix_lcm_entities_v2" "lcm-entities" {
  filter = "entityModel eq '${local.lcm.entity_model}'"
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
    to_version  = local.lcm.entity_model_version
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
    to_version  = local.lcm.entity_model_version
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
      condition     = self.entity_version == local.lcm.entity_model_version
      error_message = "entity version is not upgraded, current version: ${self.entity_version}"
    }
    postcondition {
      condition     = self.entity_model == local.lcm.entity_model
      error_message = "entity model is changed, current model: ${self.entity_model}"
    }
  }
  depends_on = [nutanix_lcm_upgrade_v2.upgrade, data.nutanix_lcm_status_v2.status-after-upgrade]
}

`, filepath)
}
