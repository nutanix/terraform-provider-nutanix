package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameClusterConfig = "nutanix_sda_cluster_config_v2.test"

func TestAccV2NutanixSdaClusterConfigResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSdaClusterConfigResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameClusterConfig, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameClusterConfig, "system_defined_policy_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameClusterConfig, "is_enabled"),
					resource.TestCheckResourceAttrSet(resourceNameClusterConfig, "last_modified_by_user"),
				),
			},
		},
	})
}

func testAccSdaClusterConfigResourceConfig() string {
	return `
		resource "nutanix_sda_cluster_config_v2" "test" {
			system_defined_policy_ext_id = "` + testVars.SystemDefinedPolicyExtID + `"
			ext_id                       = "` + testVars.ClusterExtID + `"
		}
	`
}
