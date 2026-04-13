package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccV2NutanixResourceGroupResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-rg-%d", r)
	updateName := fmt.Sprintf("tf-test-rg-%d-update", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testResourceGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupV2ResourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameResourceGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "name", name),
				),
			},
			{
				Config: testAccResourceGroupV2ResourceConfig(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameResourceGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "name", updateName),
				),
			},
		},
	})
}

func testAccResourceGroupV2ResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_resource_group_v2" "test" {
  name = "%s"
}
`, name)
}
