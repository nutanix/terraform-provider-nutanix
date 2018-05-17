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
						"data.nutanix_image.test", "name", "CentOS-LAMP-APP.qcow2"),
					resource.TestCheckResourceAttr(
						"data.nutanix_image.test", "source_uri", "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-APP.qcow2"),
				),
			},
		},
	})
}

func testAccImageDataSourceConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "CentOS-LAMP-APP.qcow2"
  description = "CentOS LAMP - App"
  image_type  = "DISK_IMAGE"
  source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-APP.qcow2"
}

data "nutanix_image" "test" {
	image_id = "${nutanix_image.test.id}"
}
`)
}
