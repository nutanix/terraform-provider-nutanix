package networkingv2_test

import (
	"fmt"
)

const datasourceNameRouteTables = "data.nutanix_route_tables_v2.test"

//func TestAccNutanixRouteTablesDataSourceV2_basic(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:  func() { acc.TestAccPreCheck(t) },
//		Providers: acc.TestAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccRouteTablesDataSourceConfig(),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttrSet(datasourceNameRouteTables, "route_tables.#"),
//				),
//			},
//		},
//	})
//}

func testAccRouteTablesDataSourceConfig() string {
	return `data "nutanix_route_tables_v2" "test" {}`
}

func testAccRouteTablesDataSourceConfigWithFilter(name, desc string) string {
	return fmt.Sprintf(`


		data "nutanix_route_tables_v2" "test" {
			filter = "name eq '%[1]s'"
			depends_on = [
				resource.nutanix_address_groups_v2.test
			]
		}
	`, name, desc)
}
