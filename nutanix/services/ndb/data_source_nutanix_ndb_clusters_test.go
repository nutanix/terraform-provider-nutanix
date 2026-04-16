package ndb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNDBClustersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNDBClustersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_clusters.test", "clusters.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_ndb_clusters.test", "clusters.0.id"),
				),
			},
		},
	})
}

func testAccNDBClustersDataSourceConfig() string {
	return `
		data "nutanix_ndb_clusters" "test" { }
	`
}
