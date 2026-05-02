package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameSdaClusterConfig = "nutanix_sda_cluster_config_v2.test"

func TestAccV2NutanixSdaClusterConfigResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSdaClusterConfigResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSdaClusterConfig, "id"),
					resource.TestCheckResourceAttrSet(resourceNameSdaClusterConfig, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSdaClusterConfig, "system_defined_policy_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSdaClusterConfig, "alert_config.#"),
				),
			},
		},
	})
}

func testSdaClusterConfigResourceConfig() string {
	return fmt.Sprintf(`
data "nutanix_sda_policies_v2" "list" {}

data "nutanix_sda_cluster_configs_v2" "list_configs" {
  system_defined_policy_ext_id = data.nutanix_sda_policies_v2.list.sda_policies.0.ext_id
}

resource "nutanix_sda_cluster_config_v2" "test" {
  system_defined_policy_ext_id = data.nutanix_sda_policies_v2.list.sda_policies.0.ext_id
  ext_id                       = data.nutanix_sda_cluster_configs_v2.list_configs.cluster_configs.0.ext_id
}
`)
}
