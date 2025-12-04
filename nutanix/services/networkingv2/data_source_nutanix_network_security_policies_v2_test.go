package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameNSPS = "data.nutanix_network_security_policies_v2.test"

func TestAccV2NutanixNSPsDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nsp-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSPsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.#"),
					checkAttributeLength(datasourceNameNSPS, "network_policies", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixNSPsDataSource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nsp-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSPsDataSourceConfigWithFilter(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.#"),
					resource.TestCheckResourceAttr(datasourceNameNSPS, "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.0.state"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.0.rules.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.0.is_system_defined"),
				),
			},
		},
	})
}

func TestAccV2NutanixNSPsDataSource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSPsDataSourceWithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNSPS, "network_policies.#"),
					resource.TestCheckResourceAttr(datasourceNameNSPS, "network_policies.#", "0"),
				),
			},
		},
	})
}

func testAccNSPsDataSourceConfig(name string) string {
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

	data "nutanix_network_security_policies_v2" "test" {
		depends_on = [nutanix_network_security_policy_v2.test]
}
	`, name)
}

func testAccNSPsDataSourceConfigWithFilter(name string) string {
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


	data "nutanix_network_security_policies_v2" "test" {
		filter = "name eq '${nutanix_network_security_policy_v2.test.name}'"
	}
	`, name)
}

func testAccNSPsDataSourceWithInvalidFilterConfig() string {
	return `
	data "nutanix_network_security_policies_v2" "test" {
		filter = "name eq 'invalid_name'"
	}
	`
}
