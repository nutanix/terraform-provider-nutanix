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
						"data.nutanix_image.nutanix_image", "name", "CentOS7-ISO"),
					resource.TestCheckResourceAttr(
						"data.nutanix_image.nutanix_image", "source_uri", "http://10.7.1.7/data1/ISOs/CentOS-7-x86_64-Minimal-1503-01.iso"),
				),
			},
		},
	})
}

func testAccImageDataSourceConfig(r int) string {
	return fmt.Sprintf(`
provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port = 9440
}

resource "nutanix_image" "test" {
  name        = "CentOS-LAMP-APP.qcow2"
  description = "CentOS LAMP - App"
  image_type  = "DISK_IMAGE"
  source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-APP.qcow2"

  metadata = {
    kind = "image"
  }
}

data "nutanix_image" "test" {
	metadata = {
		kind = "image"
	}
	image_id = "${nutanix_image.test.id}"
}
`, r)
}
