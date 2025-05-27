package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccEraProfileDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "versions.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "system_profile", "true"),
				),
			},
		},
	})
}

func TestAccEraProfileDataSource_ById(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileDataSourceConfigByID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "versions.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "system_profile", "true"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "engine_type", "postgres_database"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "type", "Software"),
				),
			},
		},
	})
}

func TestAccEraProfileDataSource_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "versions.#", "1"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "system_profile", "true"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "engine_type", "postgres_database"),
					resource.TestCheckResourceAttr("data.nutanix_ndb_profile.test", "type", "Database_Parameter"),
				),
			},
		},
	})
}

func testAccEraProfileDataSourceConfig() string {
	return `
		data "nutanix_ndb_profiles" "test1" {}

		data "nutanix_ndb_profile" "test" {
			profile_id = data.nutanix_ndb_profiles.test1.profiles.0.id
		}
	`
}

func testAccEraProfileDataSourceConfigByID() string {
	return `
		data "nutanix_ndb_profiles" "test1" {
			engine = "postgres_database"
			profile_type = "Software"
		}

		data "nutanix_ndb_profile" "test" {
			profile_id = data.nutanix_ndb_profiles.test1.profiles.0.id
		}
	`
}

func testAccEraProfileDataSourceConfigByName() string {
	return `
		data "nutanix_ndb_profiles" "test1" {
			engine = "postgres_database"
			profile_type = "Database_Parameter"
		}

		data "nutanix_ndb_profile" "test" {
			profile_name = data.nutanix_ndb_profiles.test1.profiles.0.name
		}
	`
}
