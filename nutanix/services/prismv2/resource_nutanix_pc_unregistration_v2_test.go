package prismv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameUnregisterCluster = "nutanix_pc_unregistration_v2.test"

func TestAccV2NutanixUnregisterClusterResource_Unregister_PC_PC(t *testing.T) {
	if testVars.Prism.Unregister.PcExtID == "" {
		t.Skip("Skipping test as it requires PC to unregister not provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnregisterClusterResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUnregisterCluster, "pc_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameUnregisterCluster, "ext_id"),
				),
			},
		},
	})
}

func testAccUnregisterClusterResourceConfig() string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "cls" {
	filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}
locals {
  pcExtID = data.nutanix_clusters_v2.cls.cluster_entities.0.ext_id
  config   = (jsondecode(file("%[1]s")))
  unregister = local.config.prism.unregister
}

resource "nutanix_pc_unregistration_v2" "test"{
  pc_ext_id = local.pcExtID # local pc ext id
  ext_id = local.unregister.pc_ext_id # remote pc ext id
}
`, filepath)
}
