package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraProfilesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profiles.test", "profiles.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profiles.test", "profiles.0.id"),
				),
			},
		},
	})
}

func TestAccEraProfilesDataSource_ByEngine(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesDataSourceConfigByEngine(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profiles.test", "profiles.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profiles.test", "profiles.0.id"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.engine_type", "postgres_database"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.system_profile", "true"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.topology", "ALL"),
				),
			},
		},
	})
}

func TestAccEraProfilesDataSource_ByProfileType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesDataSourceConfigByProfileType(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profiles.test", "profiles.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_profiles.test", "profiles.0.id"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.engine_type", "postgres_database"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.type", "Network"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.system_profile", "false"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profiles.test", "profiles.0.topology", "ALL"),
				),
			},
		},
	})
}

func testAccEraProfilesDataSourceConfig() string {
	return `
		data "nutanix_ndb_profiles" "test" { }
	`
}

func testAccEraProfilesDataSourceConfigByEngine() string {
	return `
		data "nutanix_ndb_profiles" "test" {
			engine = "postgres_database"
		}
	`
}

func testAccEraProfilesDataSourceConfigByProfileType() string {
	return `
		data "nutanix_ndb_profiles" "test" {
			engine = "postgres_database"
			profile_type = "Network"
		}
	`
}
