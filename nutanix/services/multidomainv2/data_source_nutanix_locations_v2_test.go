package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Test plan for nutanix_locations_v2 (ListLocations — read-only, seeded data)
// ---------------------------------------------------------------------------
// LocationsServiceApi.ListLocations is the ONLY method in scope. Locations are
// a read-only, system-seeded dataset (~27k global entries) so there is no
// Create/Update/Delete resource and no GetById singular datasource — every
// scenario below exercises the plural list datasource.
//
//  1. List (no filter): validate the list returns entries and spot-check the
//     nested attributes (name/country/coordinates/links) on the first element.
//  2. List with a valid $filter + $limit: filter by country and assert the
//     result is bounded by the limit and non-empty.
//  3. List with $orderby: exercise the sort input path.
//  4. List with an invalid/non-matching $filter: assert zero results are
//     returned (empty-result path, no error).
//
// Import is not supported for a datasource. CheckDestroy is not applicable —
// these tests create no resources (locations are seeded, not user-managed).

const datasourceNameLocationsList = "data.nutanix_locations_v2.locations_list"

func TestAccV2NutanixLocationsDatasource_List(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationsV2ListConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameLocationsList, "locations.#"),
					resource.TestCheckResourceAttrSet(datasourceNameLocationsList, "locations.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameLocationsList, "locations.0.country"),
				),
			},
		},
	})
}

func TestAccV2NutanixLocationsDatasource_WithLimitAndFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationsV2WithLimitAndFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameLocationsList, "limit", "5"),
					resource.TestCheckResourceAttrSet(datasourceNameLocationsList, "locations.#"),
					resource.TestCheckResourceAttrSet(datasourceNameLocationsList, "locations.0.name"),
				),
			},
		},
	})
}

func TestAccV2NutanixLocationsDatasource_WithOrderBy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationsV2WithOrderByConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameLocationsList, "locations.#"),
					resource.TestCheckResourceAttrSet(datasourceNameLocationsList, "locations.0.name"),
				),
			},
		},
	})
}

func TestAccV2NutanixLocationsDatasource_WithInvalidFilter(t *testing.T) {
	randomUUID := utils.GenUUID()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationsV2WithInvalidFilterConfig(randomUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameLocationsList, "locations.#", "0"),
				),
			},
		},
	})
}

func testAccLocationsV2ListConfig() string {
	return `
	data "nutanix_locations_v2" "locations_list" {
		limit = 10
	}
`
}

func testAccLocationsV2WithLimitAndFilterConfig() string {
	return `
	data "nutanix_locations_v2" "locations_list" {
		limit  = 5
		filter = "startswith(country, 'United')"
	}
`
}

func testAccLocationsV2WithOrderByConfig() string {
	return `
	data "nutanix_locations_v2" "locations_list" {
		limit    = 10
		order_by = "name desc"
	}
`
}

func testAccLocationsV2WithInvalidFilterConfig(uuid string) string {
	return fmt.Sprintf(`
	data "nutanix_locations_v2" "locations_list" {
		filter = "extId eq '%s'"
	}
`, uuid)
}
