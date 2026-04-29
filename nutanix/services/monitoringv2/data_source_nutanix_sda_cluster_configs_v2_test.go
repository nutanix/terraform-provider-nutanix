package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameClusterConfigs = "data.nutanix_sda_cluster_configs_v2.test"

func TestAccV2NutanixSdaClusterConfigsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSdaClusterConfigsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameClusterConfigs, "cluster_configs.#"),
					resource.TestCheckResourceAttrSet(datasourceNameClusterConfigs, "system_defined_policy_ext_id"),
				),
			},
		},
	})
}

func testAccSdaClusterConfigsDatasourceConfig() string {
	return `
		data "nutanix_sda_cluster_configs_v2" "test" {
			system_defined_policy_ext_id = "` + testVars.SystemDefinedPolicyExtID + `"
		}
	`
}
