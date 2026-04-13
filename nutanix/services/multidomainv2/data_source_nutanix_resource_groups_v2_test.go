package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameResourceGroupsV2 = "data.nutanix_resource_groups_v2.test"

func TestAccV2NutanixResourceGroupsDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-rg-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testResourceGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupV2ResourceConfig(name) + testAccResourceGroupsV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameResourceGroupsV2, "resource_groups.#"),
					checkAttributeLength(dataSourceNameResourceGroupsV2, "resource_groups", 1),
				),
			},
		},
	})
}

func testAccResourceGroupsV2DatasourceConfig() string {
	return `
data "nutanix_resource_groups_v2" "test" {
  depends_on = [nutanix_resource_group_v2.test]
}
`
}
