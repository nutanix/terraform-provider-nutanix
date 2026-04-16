package networkingv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRoute = "data.nutanix_route_v2.test"

func TestAccV2NutanixRouteDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-route-%d", r)
	desc := "test terraform route description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteDataSourceConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameRoute, "name", name),
					resource.TestCheckResourceAttr(datasourceNameRoute, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameRoute, "vpc_reference"),
					resource.TestCheckResourceAttrSet(datasourceNameRoute, "route_table_ext_id"),
					resource.TestCheckResourceAttr(datasourceNameRoute, "destination.0.ipv4.0.ip.0.value", "10.0.0.2"),
					resource.TestCheckResourceAttr(datasourceNameRoute, "destination.0.ipv4.0.prefix_length", "32"),
					resource.TestCheckResourceAttr(datasourceNameRoute, "next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(datasourceNameRoute, "next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttrSet(datasourceNameRoute, "metadata.0.owner_reference_id"),
					resource.TestCheckResourceAttr(datasourceNameRoute, "metadata.0.project_reference_id", testVars.Networking.Subnets.ProjectID),
					resource.TestCheckResourceAttr(datasourceNameRoute, "route_type", "STATIC"),
				),
			},
		},
	})
}

func testAccRouteDataSourceConfig(name, desc string, r int) string {
	return testRoute1Config(name, desc, r) + `
		data "nutanix_route_v2" "test"{
			ext_id             = nutanix_routes_v2.test-1.id
  			route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id		   
		}`
}
