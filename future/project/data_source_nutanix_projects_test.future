package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixProjectsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceConfig(),
			},
		},
	})
}

func testAccProjectsDataSourceConfig() string {
	return `
		data "nutanix_projects" "test" {}
	`
}
