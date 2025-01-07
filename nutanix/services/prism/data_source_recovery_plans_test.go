package prism_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixRecoveryPlansDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRecoveryPlansDataSourceConfig(),
			},
		},
	})
}

func testAccRecoveryPlansDataSourceConfig() string {
	return `
		data "nutanix_recovery_plans" "test" {}
	`
}
