package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixImageDataSource_basic(t *testing.T) {

	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_image.test", "description", "Ubuntu mini ISO"),
					resource.TestCheckResourceAttr(
						"data.nutanix_image.test",
						"source_uri",
						"http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"),
				),
			},
		},
	})
}

func testAccImageDataSourceConfig(rNumber int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu-%d"
  description = "Ubuntu mini ISO"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}


data "nutanix_image" "test" {
	image_id = "${nutanix_image.test.id}"
}
`, rNumber)
}
