package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixVolumeGroupDataSource_basic(t *testing.T) {
	t.Skip()
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_volume_group.test", "name", "Ubuntu"),
					resource.TestCheckResourceAttr(
						"data.nutanix_volume_group.test", "description", "VG Test Description"),
				),
			},
		},
	})
}

func testAccVolumeGroupDataSourceConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_volume_group" "test" {
  name        = "VG Test"
  description = "VG Test Description"
  
}

data "nutanix_volume_group" "test" {
	volume_group_id = "${nutanix_volume_group.test.id}"
}
`)
}
