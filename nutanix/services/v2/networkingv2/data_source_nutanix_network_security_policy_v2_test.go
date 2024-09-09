package networkingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameNsp = "data.nutanix_network_security_policy_v2.test"

func TestAccNutanixNSPDataSourceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNspDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNsp, "name"),
					resource.TestCheckResourceAttrSet(datasourceNameNsp, "links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsp, "state"),
					resource.TestCheckResourceAttrSet(datasourceNameNsp, "rules.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsp, "rules.0.spec.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsp, "rules.0.type"),
					resource.TestCheckResourceAttrSet(datasourceNameNsp, "is_system_defined"),
				),
			},
		},
	})
}

func testAccNspDataSourceConfig() string {
	return `

	data "nutanix_network_security_policies_v2" "test" { }

	data "nutanix_network_security_policy_v2" "test" {
		ext_id = data.nutanix_network_security_policies_v2.test.network_policies.0.ext_id
		depends_on = [data.nutanix_network_security_policies_v2.test]
	}
	`
}
