package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraProfilesAvailableIPsDataSource_basic(t *testing.T) {
	networkName := testVars.NDB.TestStaticNetwork
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesAvailableIPsDataSourceConfig(networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network_available_ips.ips", "available_ips.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network_available_ips.ips", "available_ips.0.id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network_available_ips.ips", "available_ips.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network_available_ips.ips", "available_ips.0.type"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network_available_ips.ips", "available_ips.0.managed"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_network_available_ips.ips", "available_ips.0.ip_addresses.#"),
				),
			},
		},
	})
}

func testAccEraProfilesAvailableIPsDataSourceConfig(networkName string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_profiles" "test" {
			profile_type = "Network"
		}

		locals {
			my_vlan_static = {
			  for obj in data.nutanix_ndb_profiles.test.profiles :
			  obj.name => obj
			  if obj.name == "%[1]s"
			}["%[1]s"]
		  }

		data "nutanix_ndb_network_available_ips" "ips"{
			profile_id =  local.my_vlan_static.id
		}
	`, networkName)
}
