package fc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCAPIKey_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-key")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		// CheckDestroy: testAccCheckNutanixAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFCAPIKeyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_foundation_central_api_keys.test", "alias", name),
				),
			},
		},
	})
}

func testAccFCAPIKeyConfig(name string) string {
	return fmt.Sprintf(`
	resource "nutanix_foundation_central_api_keys" "test"{
		alias = "%s"
	}
`, name)
}
