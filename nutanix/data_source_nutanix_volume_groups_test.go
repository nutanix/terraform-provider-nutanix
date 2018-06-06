package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixVolumeGroupsDataSource_basic(t *testing.T) {
	// skipping as this API is not yet GA (will GA in upcoming AOS release)
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_volume_groups.test", "entities.#", "2"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccVolumeGroupsDataSourceConfig = `
resource "nutanix_volume_group" "test" {
  name        = "VG Test"
  description = "VG Test Description"
  
}

resource "nutanix_volume_group" "test-1" {
  name        = "VG Test-1"
  description = "VG Test-1 Description"
  
}

data "nutanix_volume_groups" "test" {
	metadata = {
		length = 2
	}
}
`
