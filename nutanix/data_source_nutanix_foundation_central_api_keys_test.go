package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixFCAPIKeysDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeysDataSourceConfig(),
			},
		},
	})
}

func TestAccNutanixFCAPIKeysDataSource_KeyUUID(t *testing.T) {
	apiKeyName := acctest.RandomWithPrefix("test-key")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeysDataSourceConfigWithKeyUUID(apiKeyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_foundation_central_api_keys.k1", "alias", apiKeyName),
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_api_keys.k1", "alias"),
				),
			},
		},
	})
}

func testAccAPIKeysDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_list_api_keys" "test"{}
	`
}

func testAccAPIKeysDataSourceConfigWithKeyUUID(apiKeyName string) string {
	return fmt.Sprintf(`
		resource "nutanix_foundation_central_api_keys" "apk"{
				alias = "%s"
			}
	
		data "nutanix_foundation_central_api_keys" "k1"{
		    key_uuid = "${nutanix_foundation_central_api_keys.apk.key_uuid}"
		}
		
	 `, apiKeyName)
}
