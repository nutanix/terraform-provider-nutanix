package iamv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameOperation = "data.nutanix_operation_v2.test"

func TestAccV2NutanixOperationDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOperationV2DatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameOperation, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameOperation, "display_name"),
					resource.TestCheckResourceAttrSet(datasourceNameOperation, "entity_type"),
				),
			},
		},
	})
}

func testOperationV2DatasourceConfig() string {
	return `
		data "nutanix_operations_v2" "test" {}

		data "nutanix_operation_v2" "test" {
			ext_id = data.nutanix_operations_v2.test.operations.0.ext_id		
}
	`
}
