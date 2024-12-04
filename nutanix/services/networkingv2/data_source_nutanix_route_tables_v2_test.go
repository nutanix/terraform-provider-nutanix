package networkingv2_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"testing"
)

const datasourceNameRouteTables = "data.nutanix_route_tables_v2.test"

func TestAccNutanixRouteTablesV2DataSource_Basic(t *testing.T) {
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
func TestAccNutanixRouteTablesV2DataSource_WithFilter(t *testing.T) {
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
