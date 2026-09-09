package lifecyclev2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test plan for nutanix_create_connection_by_hardware_provider_id_v2 (action resource)
// ----------------------------------------------------------------------------------
// This is the standalone create-action surface for the connection create API. The
// test creates a connection through the action resource and asserts the created
// connection's ext_id and name.
//
//  1. TestAccV2NutanixCreateConnectionByHardwareProviderIDResource_Basic

func TestAccV2NutanixCreateConnectionByHardwareProviderIDResource_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-create-conn-%d", acctest.RandInt())
	createResourceName := "nutanix_create_connection_by_hardware_provider_id_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCreateConnectionByHardwareProviderIDConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(createResourceName, "name", name),
					resource.TestCheckResourceAttrSet(createResourceName, "ext_id"),
				),
			},
		},
	})
}

func testCreateConnectionByHardwareProviderIDConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_create_connection_by_hardware_provider_id_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  name                     = "%[2]s"
  region                   = "us-west"
  access_details {
    auth {
      basic_auth {
        username = "create-user"
        password = "create-pass"
      }
    }
    endpoint {
      url_endpoint {
        url = "https://hardware-provider.example.com"
      }
    }
  }
}
`, testVars.Lifecycle.HardwareProviders.HardwareProviderExtID, name)
}
