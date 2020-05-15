package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixHostsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHostsDataSourceConfig(),
			},
		},
	})
}

func testAccHostsDataSourceConfig() string {
	return `
		data "nutanix_hosts" "test" {}
	`
}
