package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixRecoveryPlansDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
