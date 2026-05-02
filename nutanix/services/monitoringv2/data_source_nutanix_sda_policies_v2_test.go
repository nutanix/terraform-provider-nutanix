package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameSdaPolicies = "data.nutanix_sda_policies_v2.test"

func TestAccV2NutanixSdaPoliciesDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSdaPoliciesDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameSdaPolicies, "sda_policies.#"),
				),
			},
		},
	})
}

func testSdaPoliciesDatasourceConfig() string {
	return `
data "nutanix_sda_policies_v2" "test" {}
`
}
