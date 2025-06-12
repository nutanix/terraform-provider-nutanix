package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRoutes = "data.nutanix_routes_v2.test"

func TestAccV2NutanixRoutesDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-route-%d", r)
	desc := "test terraform route description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutesDataSourceConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRoutes, "routes.#"),
					checkAttributeLength(datasourceNameRoutes, "routes", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixRoutesDataSource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-route-%d", r)
	desc := "test terraform route description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutesDataSourceWithFilterConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRoutes, "routes.#"),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.0.description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameRoutes, "routes.0.vpc_reference"),
					resource.TestCheckResourceAttrSet(datasourceNameRoutes, "route_table_ext_id"),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.0.destination.0.ipv4.0.ip.0.value", "10.0.0.2"),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.0.destination.0.ipv4.0.prefix_length", "32"),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.0.next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(datasourceNameRoutes, "routes.0.next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttrSet(datasourceNameRoutes, "routes.0.metadata.0.owner_reference_id"),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.0.metadata.0.project_reference_id", testVars.Networking.Subnets.ProjectID),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.0.route_type", "STATIC"),
				),
			},
		},
	})
}

func TestAccV2NutanixRoutesDataSource_WithInvalidFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-route-%d", r)
	desc := "test terraform route description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutesDataSourceWithInvalidFilterConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRoutes, "routes.#"),
					resource.TestCheckResourceAttr(datasourceNameRoutes, "routes.#", "0"),
				),
			},
		},
	})
}

func testAccRoutesDataSourceConfig(name, desc string, r int) string {
	return testRoute1Config(name, desc, r) + `
		data "nutanix_routes_v2" "test"{
  			route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
			depends_on = [nutanix_routes_v2.test-1]
		}`
}

func testAccRoutesDataSourceWithFilterConfig(name, desc string, r int) string {
	return testRoute1Config(name, desc, r) + `
		data "nutanix_routes_v2" "test"{
			filter             = "name eq '${nutanix_routes_v2.test-1.name}'"
  			route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
			depends_on = [nutanix_routes_v2.test-1]
		}`
}

func testAccRoutesDataSourceWithInvalidFilterConfig(name, desc string, r int) string {
	return testRoute1Config(name, desc, r) + `
		data "nutanix_routes_v2" "test"{
			filter             = "name eq 'invalid'"
  			route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
			depends_on = [nutanix_routes_v2.test-1]
		}`
}
