package lcmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameLcmPreChecks = "nutanix_lcm_prechecks_v2.pre-checks"

func TestAccV2NutanixLcmPrechecks_Basic(t *testing.T) {
	datasourceNameLcmEntities := "data.nutanix_lcm_entities_v2.lcm-entities"
	datasourceNameLcmEntityBeforeUpgrade := "data.nutanix_lcm_entity_v2.entity-before-upgrade"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testLcmPrechecksConfig(),
				Check: resource.ComposeTestCheckFunc(
					// check if the entity model is present
					resource.TestCheckResourceAttr(datasourceNameLcmEntities, "entities.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameLcmEntities, "entities.0.entity_model", testVars.Lcm.EntityModel),
					// check if the entity is present
					resource.TestCheckResourceAttrSet(datasourceNameLcmEntityBeforeUpgrade, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameLcmEntityBeforeUpgrade, "entity_model", testVars.Lcm.EntityModel),
					resource.TestCheckResourceAttrSet(datasourceNameLcmEntityBeforeUpgrade, "entity_version"),
					// perform inventory checks
					resource.TestCheckResourceAttrSet(resourceNameLcmPerformInventory, "x_cluster_id"),
					// pre-checks checks
					resource.TestCheckResourceAttrSet(resourceNameLcmPreChecks, "x_cluster_id"),
					resource.TestCheckResourceAttr(resourceNameLcmPreChecks, "entity_update_specs.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameLcmPreChecks, "entity_update_specs.0.entity_uuid", datasourceNameLcmEntityBeforeUpgrade, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameLcmPreChecks, "entity_update_specs.0.to_version", testVars.Lcm.EntityModelVersion),
				),
			},
		},
	})
}

func testLcmPrechecksConfig() string {
	return fmt.Sprintf(`
# list Clusters
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

# List Prism Central
data "nutanix_clusters_v2" "pc" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
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



`, filepath)
}
