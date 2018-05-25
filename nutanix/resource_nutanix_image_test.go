package nutanix

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixImage_basic(t *testing.T) {
	r := rand.Int31()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
			{
				Config: testAccNutanixImageConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
		},
	})
}

// Skipping this test for now, since it is difficult to implement on the CI
func TestAccNutanixImage_basic_uploadLocal(t *testing.T) {

	t.Skip()

	r := rand.Int31()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageLocalConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
			{
				Config: testAccNutanixImageLocalConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
		},
	})
}

func testAccCheckNutanixImageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixImageDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_image" {
			continue
		}
		for {
			_, err := conn.API.V3.GetImage(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}

	}

	return nil
}

func testAccNutanixImageConfig(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu"
  description = "Ubuntu"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}
`)
}

func testAccNutanixImageConfigUpdate(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu Updated"
  description = "Ubuntu Updated"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}
`)
}

func testAccNutanixImageLocalConfig(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu"
  description = "Ubuntu"
  source_path  = "/Users/thetonymaster/development/src/github.com/terraform-providers/terraform-provider-nutanix/mini.iso"
}
`)
}

func testAccNutanixImageLocalConfigUpdate(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu Updated"
  description = "Ubuntu Updated"
  source_path  = "/Users/thetonymaster/development/src/github.com/terraform-providers/terraform-provider-nutanix/alp.iso"
}
`)
}
