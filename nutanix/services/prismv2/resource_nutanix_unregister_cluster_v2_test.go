package prismv2_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameUnregisterCluster = "nutanix_unregister_cluster_v2.test"

func TestAccV2NutanixUnregisterClusterResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUnregisterClusterResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						aJson, _ := json.MarshalIndent(s.RootModule().Resources[resourceNameUnregisterCluster].Primary.Attributes, "", "  ")
						fmt.Println("############################################")
						fmt.Println(fmt.Sprintf("Resource Attributes: \n%v", string(aJson)))
						fmt.Println("############################################")

						return nil
					},
					resource.TestCheckResourceAttrSet(resourceNameUnregisterCluster, "pc_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameUnregisterCluster, "ext_id"),
				),
			},
		},
	})
}

func testAccUnregisterClusterResourceConfig() string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "clusters" {}

locals {
  pcExtID = [
	for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
	cluster.ext_id if cluster.config[0].cluster_function[0] == "PRISM_CENTRAL"
  ][0]
  config   = (jsondecode(file("%[1]s")))
  prism = local.config.prism
}

resource "nutanix_unregister_cluster_v2" "test"{
  pc_ext_id = local.pcExtID
  ext_id = local.prism.cluster_ext_id
}
`, filepath)
}
