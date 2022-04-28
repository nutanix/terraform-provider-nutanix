package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNutanixFCClusterListDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCClusterListDataSourceConfig(),
			},
		},
	})
}

func testAccFCClusterListDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_imaged_clusters_list" "cls" {}
	`
}
