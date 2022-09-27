package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEraClustersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraClustersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_era_clusters.test", "clusters.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_era_clusters.test", "clusters.0.id"),
				),
			},
		},
	})
}

func testAccEraClustersDataSourceConfig() string {
	return `
		data "nutanix_era_clusters" "test" { }
	`
}
