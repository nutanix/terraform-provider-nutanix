package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixFCAPIKeysListDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeysListDataSourceConfig(),
			},
		},
	})
}

func testAccAPIKeysListDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_list_api_keys" "test"{}
	`
}
