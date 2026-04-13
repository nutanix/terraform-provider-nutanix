package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameResourceGroupV2 = "data.nutanix_resource_group_v2.test"

func TestAccV2NutanixResourceGroupDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-rg-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testResourceGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupV2ResourceConfig(name) + testAccResourceGroupV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameResourceGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameResourceGroupV2, "name", name),
				),
			},
		},
	})
}

func testAccResourceGroupV2DatasourceConfig() string {
	return `
data "nutanix_resource_group_v2" "test" {
  ext_id = nutanix_resource_group_v2.test.id
  depends_on = [nutanix_resource_group_v2.test]
}
`
}
