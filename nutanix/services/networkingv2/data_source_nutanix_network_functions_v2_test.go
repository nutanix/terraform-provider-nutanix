package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
)

const datasourceNameNetworkFunctionsV2 = "data.nutanix_network_functions_v2.test"

func TestAccV2NutanixNetworkFunctionsDataSource_Basic(t *testing.T) {
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
				Config: networkFunctionConfig + testAccNetworkFunctionsV2DataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					common.CheckAttributeLength(datasourceNameNetworkFunctionsV2, "network_functions", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixNetworkFunctionsDataSource_FilterAndLimit(t *testing.T) {
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
			{
				Config: networkFunctionConfig,
				Check: resource.ComposeTestCheckFunc(
					waitForNetworkFunctionHealth(resourceNameNetworkFunctionV2_1, "data_plane_health_status", "HEALTHY"),
				),
			},
			{
				Config: networkFunctionConfig + testAccNetworkFunctionsV2DataSourceConfigFilterAndLimit(),
				Check: resource.ComposeTestCheckFunc(
					common.CheckAttributeLengthEqual(datasourceNameNetworkFunctionsV2, "network_functions", 1),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "filter", "name eq '"+name+"'"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "limit", "1"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.name", name),
					resource.TestCheckResourceAttrSet(datasourceNameNetworkFunctionsV2, "network_functions.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.description", "First Network function managed by Terraform"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.failure_handling", "FAIL_CLOSE"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.high_availability_mode", "ACTIVE_PASSIVE"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.traffic_forwarding_mode", "INLINE"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.nic_pairs.#", "2"),
					testAccCheckNetworkFunctionDataSourcePairWithPrefix(
						datasourceNameNetworkFunctionsV2,
						"network_functions.0.",
						"nutanix_virtual_machine_v2.vm-1",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					testAccCheckNetworkFunctionDataSourcePairWithPrefix(
						datasourceNameNetworkFunctionsV2,
						"network_functions.0.",
						"nutanix_virtual_machine_v2.vm-2",
						"nics.2.ext_id",
						"nics.1.ext_id",
					),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.nic_pairs.0.data_plane_health_status", "HEALTHY"),
					resource.TestCheckResourceAttr(datasourceNameNetworkFunctionsV2, "network_functions.0.nic_pairs.1.data_plane_health_status", "HEALTHY"),
					testAccCheckNetworkFunctionNICPairHAStateCountsWithPrefix(
						datasourceNameNetworkFunctionsV2,
						"network_functions.0.",
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

func testAccNetworkFunctionsV2DataSourceConfig() string {
	return `


data "nutanix_network_functions_v2" "test" {
  depends_on = [nutanix_network_function_v2.ntf-1]
}
  `
}

func testAccNetworkFunctionsV2DataSourceConfigFilterAndLimit() string {
	return `

data "nutanix_network_functions_v2" "test" {
  filter = "name eq '${nutanix_network_function_v2.ntf-1.name}'"
  limit = 1
}
  `
}
