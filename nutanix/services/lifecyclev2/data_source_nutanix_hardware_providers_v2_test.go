package lifecyclev2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_hardware_providers_v2 / nutanix_hardware_provider_v2
// -------------------------------------------------------------------------
// Hardware providers are read-only from Terraform (no Create). These datasource
// tests validate the list datasource returns providers on the cluster and that
// the singular datasource can read one by ext_id.
//
//  1. TestAccV2NutanixHardwareProvidersDatasource_Basic: list all hardware providers
//     and assert the collection count is present.
//  2. TestAccV2NutanixHardwareProviderDatasource_Basic: read a single hardware
//     provider by ext_id (read from test config) and assert its attributes.

const (
	datasourceNameHardwareProviders = "data.nutanix_hardware_providers_v2.test"
	datasourceNameHardwareProvider  = "data.nutanix_hardware_provider_v2.test"
)

func TestAccV2NutanixHardwareProvidersDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testHardwareProvidersDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameHardwareProviders, "hardware_providers.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixHardwareProviderDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testHardwareProviderDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameHardwareProvider, "name"),
					resource.TestCheckResourceAttrSet(datasourceNameHardwareProvider, "ext_id"),
				),
			},
		},
	})
}

func testHardwareProvidersDatasourceConfig() string {
	return `
data "nutanix_hardware_providers_v2" "test" {}
`
}

func testHardwareProviderDatasourceConfig() string {
	return `
data "nutanix_hardware_providers_v2" "all" {}

data "nutanix_hardware_provider_v2" "test" {
  ext_id = data.nutanix_hardware_providers_v2.all.hardware_providers.0.ext_id
}
`
}
