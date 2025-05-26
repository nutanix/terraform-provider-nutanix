package networkingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRouteTable = "data.nutanix_route_table_v2.test"

func TestAccV2NutanixRouteTableDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()

	//goland:noinspection ALL
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableDataSourceWithFilterConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRouteTable, "ext_id"),
				),
			},
		},
	})
}

func testAccRouteTableDataSourceWithFilterConfig(r int) string {
	return testRouteTableInfoVpc1Config(r) + `
	data "nutanix_route_tables_v2" "test" {
		filter     = "vpcReference eq '${nutanix_vpc_v2.test-1.id}'"
		depends_on = [nutanix_vpc_v2.test-1]
	}
    
	data "nutanix_route_table_v2" "test" {
		ext_id             = data.nutanix_route_tables_v2.test.route_tables[0].ext_id
		depends_on         = [data.nutanix_route_tables_v2.test]
	}

`
}
