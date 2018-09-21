package nutanix

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixImage_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "nutanix_image.acctest-test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists(resourceName),
				),
			},
			{
				Config: testAccNutanixImageConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("Ubuntu-%d-updated", rInt)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixImageWithCategories(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "nutanix_image.acctest-test-categories"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageConfigWithCategories(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "categories.os_type", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "categories.os_version", "current"),
				),
			},
			{
				Config: testAccNutanixImageConfigWithCategoriesUpdated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "categories.os_type", "ubuntu"),
					resource.TestCheckResourceAttr(resourceName, "categories.os_version", "18.04"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixImage_basic_uploadLocal(t *testing.T) {
	// Skipping as this test needs functional work
	t.Skip()
	// function guts inspired by resource_aws_s3_bucket_object_test.go
	tmpFile, err := ioutil.TempFile("", "tf-acc-image-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	rInt := acctest.RandInt()
	// first write some data to the tempfile just so it's not 0 bytes.
	err = ioutil.WriteFile(tmpFile.Name(), []byte("{anything will do }"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageLocalConfig(rInt, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.acctest-testLocal"),
				),
			},
			{
				Config: testAccNutanixImageLocalConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.acctest-testLocal"),
					resource.TestCheckResourceAttr("nutanix_image.acctest-testLocal", "name", "image-updateName"),
				),
			},
		},
	})
}

func testAccCheckNutanixImageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
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

func testAccNutanixImageConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-test" {
  name        = "Ubuntu-%d"
  description = "Ubuntu"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}
`, r)
}

func testAccNutanixImageConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-test" {
  name        = "Ubuntu-%d-updated"
  description = "Ubuntu Updated"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}
`, r)
}

func testAccNutanixImageLocalConfig(rNumb int, rFile string) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-testLocal" {
  name        = "random-local-image-%d"
  description = "some description"
  source_path  = "%s"
}
`, rNumb, rFile)
}

func testAccNutanixImageLocalConfigUpdate(rNumb int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-testLocal" {
  name        = "random-local-image-%d"
  description = "update my description!"
}
`, rNumb)
}

func testAccNutanixImageConfigWithCategories(r int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-test-categories" {
  name        = "Ubuntu-%d"
  description = "Ubuntu"

	categories = {
		os_type = "ubuntu"
		os_version = "current"
	}

  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"

}
`, r)
}

func testAccNutanixImageConfigWithCategoriesUpdated(r int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-test-categories" {
  name        = "Ubuntu-%d"
  description = "Ubuntu"

	categories = {
		os_type = "ubuntu"
		os_version = "18.04"
	}

  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"

}
`, r)
}
