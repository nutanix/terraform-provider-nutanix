package nutanix

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func TestAccNutanixImageDataSource_name(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceConfigName(rInt),
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

func TestAccNutanixImageDataSource_conflicts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccImageDataSourceConfigConflicts(),
				ExpectError: regexp.MustCompile("conflicts with"),
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
	image_id = nutanix_image.test.id
}
`, rNumber)
}

func testAccImageDataSourceConfigName(rNumber int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu-%d"
  description = "Ubuntu mini ISO"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}


data "nutanix_image" "test" {
	image_name = nutanix_image.test.name
}
`, rNumber)
}

func testAccImageDataSourceConfigConflicts() string {
	return `
data "nutanix_image" "test" {
	image_name = "test-name"
	image_id   = "test-id"
}
`
}
