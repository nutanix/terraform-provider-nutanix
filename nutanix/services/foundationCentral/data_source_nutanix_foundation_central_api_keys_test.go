package fc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCAPIKeysDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeysDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_list_api_keys.test", "api_keys.#"),
				),
			},
		},
	})
}

func TestAccFCAPIKeysDataSource_KeyUUID(t *testing.T) {
	apiKeyName := acctest.RandomWithPrefix("test-key")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeysDataSourceConfigWithKeyUUID(apiKeyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_foundation_central_api_keys.k1", "alias", apiKeyName),
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_api_keys.k1", "alias"),
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_api_keys.k1", "created_timestamp"),
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_api_keys.k1", "current_time"),
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
