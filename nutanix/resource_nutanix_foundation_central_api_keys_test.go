package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixFCAPIKey(t *testing.T) {
	name := acctest.RandomWithPrefix("test-key")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixFCAPIKeyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_foundation_central_api_keys.test", "alias", name),
					resource.TestCheckResourceAttr(resourcesAddressGroup, "ip_address_block_list.0.prefix_length", "24"),
				),
			},
		},
	})
}

func testAccNutanixFCAPIKeyConfig(name string) string {
	return fmt.Sprintf(`
	resource "nutanix_foundation_central_api_keys" "test"{
		alias = "%s"
	}
`, name)
}
