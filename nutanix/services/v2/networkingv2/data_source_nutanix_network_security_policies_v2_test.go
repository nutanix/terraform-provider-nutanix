package networkingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameNsps = "data.nutanix_network_security_policies_v2.test"

func TestAccNutanixNSPsDataSourceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNspsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.state"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.rules.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.is_system_defined"),
				),
			},
		},
	})
}

func TestAccNutanixNSPsDataSourceV2_WithFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNspsDataSourceConfigWithFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.#"),
					resource.TestCheckResourceAttr(datasourceNameNsps, "network_policies.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.links.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.state"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.rules.#"),
					resource.TestCheckResourceAttrSet(datasourceNameNsps, "network_policies.0.is_system_defined"),
				),
			},
		},
	})
}

func testAccNspsDataSourceConfig() string {
	return `

	data "nutanix_network_security_policies_v2" "test" { }
	`
}

func testAccNspsDataSourceConfigWithFilter() string {
	return `

	data "nutanix_network_security_policies_v2" "dtest" { }

	locals {
		nsp_name = data.nutanix_network_security_policies_v2.dtest.network_policies.0.name
	}

	data "nutanix_network_security_policies_v2" "test" {
		filter = "name eq '${local.nsp_name}'"
	}
	`
}
