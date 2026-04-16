package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameNsp = "data.nutanix_network_security_policy_v2.test"

func TestAccV2NutanixNSPDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nsp-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNspDataSourceConfig(name),
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

func testAccNspDataSourceConfig(name string) string {
	return fmt.Sprintf(`

	data "nutanix_categories_v2" "test" {}

	resource "nutanix_network_security_policy_v2" "test" {
		name = "%[1]s"
		description = "test nsp description"
		state = "SAVE"
		type = "ISOLATION"
		rules{
		  type = "TWO_ENV_ISOLATION"
		  spec{
			two_env_isolation_rule_spec{
			  first_isolation_group = [
				data.nutanix_categories_v2.test.categories[0].ext_id,
			  ]
			  second_isolation_group =  [
				data.nutanix_categories_v2.test.categories[1].ext_id,
			  ]
			}
		  }
		}
		is_hitlog_enabled = true
		depends_on = [data.nutanix_categories_v2.test]
	  }

	data "nutanix_network_security_policy_v2" "test" {
		ext_id = nutanix_network_security_policy_v2.test.ext_id
		depends_on = [nutanix_network_security_policy_v2.test]
	}
	`, name)
}
