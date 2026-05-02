package monitoringv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSdaClusterConfigs = "data.nutanix_sda_cluster_configs_v2.test"

func TestAccV2NutanixSdaClusterConfigsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSdaClusterConfigsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSdaClusterConfigs, "cluster_configs.#"),
					resource.TestCheckResourceAttrSet(datasourceNameSdaClusterConfigs, "system_defined_policy_ext_id"),
				),
			},
		},
	})
}

func testSdaClusterConfigsDatasourceConfig() string {
	return fmt.Sprintf(`
data "nutanix_sda_policies_v2" "list" {}

data "nutanix_sda_cluster_configs_v2" "test" {
  system_defined_policy_ext_id = data.nutanix_sda_policies_v2.list.sda_policies.0.ext_id
}
`)
}
