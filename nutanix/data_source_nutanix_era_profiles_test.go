package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraProfilesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_era_profiles.test", "profiles.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_era_profiles.test", "profiles.0.id"),
				),
			},
		},
	})
}

func TestAccEraProfilesDataSource_ByEngine(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesDataSourceConfigByEngine(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_era_profiles.test", "profiles.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_era_profiles.test", "profiles.0.id"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.engine_type", "postgres_database"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.system_profile", "true"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.topology", "ALL"),
				),
			},
		},
	})
}

func TestAccEraProfilesDataSource_ByProfileType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfilesDataSourceConfigByProfileType(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_era_profiles.test", "profiles.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_era_profiles.test", "profiles.0.id"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.status", "READY"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.engine_type", "postgres_database"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.type", "Network"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.system_profile", "false"),
					resource.TestCheckResourceAttr("data.nutanix_era_profiles.test", "profiles.0.topology", "ALL"),
				),
			},
		},
	})
}

func testAccEraProfilesDataSourceConfig() string {
	return `
		data "nutanix_era_profiles" "test" { }
	`
}

func testAccEraProfilesDataSourceConfigByEngine() string {
	return `
		data "nutanix_era_profiles" "test" {
			engine = "postgres_database"
		}
	`
}

func testAccEraProfilesDataSourceConfigByProfileType() string {
	return `
		data "nutanix_era_profiles" "test" {
			engine = "postgres_database"
			profile_type = "Network"
		}
	`
}
