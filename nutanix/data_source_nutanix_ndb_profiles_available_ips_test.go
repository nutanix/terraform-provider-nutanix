package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraProfilesAvailableIPsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesAvailableIPsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profile_available_ips.ips", "available_ips.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profile_available_ips.ips", "available_ips.0.id"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profile_available_ips.ips", "available_ips.0.name"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profile_available_ips.ips", "available_ips.0.type"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profile_available_ips.ips", "available_ips.0.managed"),
				),
			},
		},
	})
}

func testAccEraProfilesAvailableIPsDataSourceConfig() string {
	return `
		data "nutanix_ndb_profiles" "test" {
			profile_type = "Network"
		}

		data "nutanix_ndb_profile_available_ips" "ips"{
			profile_id =  data.nutanix_ndb_profiles.test.profiles.0.id
		}
	`
}
