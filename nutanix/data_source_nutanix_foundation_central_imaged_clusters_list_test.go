package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFCClusterListDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFCClusterListDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_foundation_central_imaged_clusters_list.cls", "imaged_clusters.#"),
				),
			},
		},
	})
}

func testAccFCClusterListDataSourceConfig() string {
	return `
	data "nutanix_foundation_central_imaged_clusters_list" "cls" {}
	`
}
