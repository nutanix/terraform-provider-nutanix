package iamv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameOperations = "data.nutanix_operations_v2.test"

func TestAccV2NutanixOperationsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOperationsV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameOperations, "operations.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixOperationsDatasource_WithLimit(t *testing.T) {
	limit := 3
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOperationsV2DatasourceWithLimitConfig(limit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameOperations, "operations.#"),
					resource.TestCheckResourceAttr(datasourceNameOperations, "operations.#", strconv.Itoa(limit)),
				),
			},
		},
	})
}

func testOperationsV2DatasourceConfig() string {
	return `
		data "nutanix_operations_v2" "test" {}
	`
}

func testOperationsV2DatasourceWithLimitConfig(limit int) string {
	return fmt.Sprintf(`

		data "nutanix_operations_v2" "test" {
		  limit = %d
		}
	`, limit)
}
