package networkingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRouteTables = "data.nutanix_route_tables_v2.test"

func TestAccV2NutanixRouteTablesDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesDataSourceConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRouteTables, "route_tables.#"),
					checkAttributeLength(datasourceNameRouteTables, "route_tables", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixRouteTablesDataSource_WithFilter(t *testing.T) {
	r := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesDataSourceWithFilterConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRouteTables, "route_tables.#"),
					resource.TestCheckResourceAttr(datasourceNameRouteTables, "route_tables.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameRouteTables, "route_tables.0.ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixRouteTablesDataSource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTablesDataSourceWithInvalidFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRouteTables, "route_tables.#"),
					resource.TestCheckResourceAttr(datasourceNameRouteTables, "route_tables.#", "0"),
				),
			},
		},
	})
}

func testAccRouteTablesDataSourceConfig(r int) string {
	return testRouteTableInfoVpc1Config(r) + `data "nutanix_route_tables_v2" "test" {}`
}

func testAccRouteTablesDataSourceWithFilterConfig(r int) string {
	return testRouteTableInfoVpc1Config(r) + `
	data "nutanix_route_tables_v2" "test" {
		filter     = "vpcReference eq '${nutanix_vpc_v2.test-1.id}'"
		depends_on = [nutanix_vpc_v2.test-1]
	}`
}

func testAccRouteTablesDataSourceWithInvalidFilterConfig() string {
	return `
	data "nutanix_route_tables_v2" "test" {
		filter = "vpcReference eq 'invalid'"
	}`
}
