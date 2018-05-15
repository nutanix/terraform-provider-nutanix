package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleClustersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClustersDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_clusters.basic_web", "entities.#", "3"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccClustersDataSourceConfig = `
provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}
data "nutanix_clusters" "basic_web" {
	metadata = {
		length = 3
	}
}`
