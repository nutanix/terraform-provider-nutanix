package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameNetworkFunctionV2 = "data.nutanix_network_function_v2.test"

func TestAccV2NutanixNetworkFunctionDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()

	subnet_name := fmt.Sprintf("tf-test-subnet-%d", r)
	vmm_1_name := fmt.Sprintf("tf-test-vm-1-%d", r)
	vmm_2_name := fmt.Sprintf("tf-test-vm-2-%d", r)
	name := fmt.Sprintf("tf-test-network-function-%d", r)

	networkFunctionConfig := testAccNetworkFunctionV2ConfigPrerequisites(subnet_name, vmm_1_name, vmm_2_name) + testAccNetworkFunctionV2EgressIngressConfig(name)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Prerequisites
			{
				Config: networkFunctionConfig,
				Check: resource.ComposeTestCheckFunc(
					waitForNetworkFunctionHealth(resourceNameNetworkFunctionV2_1, "data_plane_health_status", "HEALTHY"),
				),
			},
			{
				Config: networkFunctionConfig + testAccNetworkFunctionV2DataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNetworkFunctionV2, "id"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionV2, "name", name),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionV2, "description", "First Network function managed by Terraform"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionV2, "high_availability_mode", "ACTIVE_PASSIVE"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionV2, "nic_pairs.#", "2"),
					testAccCheckNetworkFunctionDataSourcePair(
						datasourceNameNetworkFunctionV2,
						"nutanix_virtual_machine_v2.vm-1",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					testAccCheckNetworkFunctionDataSourcePair(
						datasourceNameNetworkFunctionV2,
						"nutanix_virtual_machine_v2.vm-2",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionV2, "nic_pairs.0.data_plane_health_status", "HEALTHY"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionV2, "nic_pairs.1.data_plane_health_status", "HEALTHY"),
					testAccCheckNetworkFunctionNICPairHAStateCounts(
						datasourceNameNetworkFunctionV2,
						map[string]int{
							"ACTIVE":  1,
							"PASSIVE": 1,
						},
					),
				),
			},
		},
	})
}

func testAccNetworkFunctionV2DataSourceConfig() string {
	return `


data "nutanix_network_function_v2" "test" {
  ext_id = nutanix_network_function_v2.ntf-1.id
}
`
}
