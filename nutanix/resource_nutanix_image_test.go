package nutanix

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceName = "nutanix_image.acctest-test"

func TestAccNutanixImage_basic(t *testing.T) {
	rInt := acctest.RandInt()

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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixImage_Update(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "image_type", "ISO_IMAGE"),
				),
			},
			{
				Config: testAccNutanixImageConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("Ubuntu-%d-updated", rInt)),
					resource.TestCheckResourceAttr(resourceName, "image_type", "DISK_IMAGE"),
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

func TestAccNutanixImage_WithCategoriesAndCluster(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "categories.#", "2"),
					//resource.TestCheckResourceAttr(resourceName, "categories.os_type", "ubuntu"),
					//resource.TestCheckResourceAttr(resourceName, "categories.os_version", "current"),
				),
			},
			{
				Config: testAccNutanixImageConfigWithCategoriesUpdated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "2"),
					//resource.TestCheckResourceAttr(resourceName, "categories.os_type", "ubuntu"),
					//resource.TestCheckResourceAttr(resourceName, "categories.os_version", "18.04"),
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

func TestAccNutanixImage_WithLargeImageURL(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageConfigWithLargeImageURL(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("Ubuntu-%d-server", rInt)),
					testAccCheckNutanixImageExists(resourceName),
				),
			},
		},
	})
}

func TestAccNutanixImage_uploadLocal(t *testing.T) {
	//Skipping Because in GCP still failing
	if isGCPEnvironment() {
		t.Skip()
	}

	// Get the Working directory
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("TestAccNutanixImage_uploadLocal failed to get working directory %s", err)
	}

	filepath := dir + "/alpine.iso"

	defer os.Remove(filepath)
	//Small Alpine image
	image := "http://dl-cdn.alpinelinux.org/alpine/v3.8/releases/x86_64/alpine-virt-3.8.1-x86_64.iso"

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if err := downloadFile(filepath, image); err != nil {
				t.Errorf("TestAccNutanixImage_uploadLocal failed to download image %s", err)
			}
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageLocalConfig(rInt, filepath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.acctest-testLocal"),
				),
			},
			{
				Config: testAccNutanixImageLocalConfigUpdate(rInt, filepath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.acctest-testLocal"),
					resource.TestCheckResourceAttr("nutanix_image.acctest-testLocal", "description", "new description"),
				),
			},
		},
	})
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
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
  image_type = "ISO_IMAGE"
}
`, r)
}

func testAccNutanixImageConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-test" {
  name        = "Ubuntu-%d-updated"
  description = "Ubuntu Updated"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
  image_type = "DISK_IMAGE"
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

func testAccNutanixImageLocalConfigUpdate(rNumb int, rFile string) string {
	return fmt.Sprintf(`
resource "nutanix_image" "acctest-testLocal" {
  name        = "random-local-image-%d"
  description = "new description"
  source_path  = "%s"
}
`, rNumb, rFile)
}

func testAccNutanixImageConfigWithCategories(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "os_version"{
	name = "os_version"
	description = "testacc-os-version"
}

resource "nutanix_category_value" "os_version_value"{
	name = nutanix_category_key.os_version.id
	description = "testacc-os-current"
	value = "os_current"
}

resource "nutanix_category_key" "os_type"{
	name = "os_type"
	description = "testacc-os-type"
}

resource "nutanix_category_value" "ubuntu"{
	name = nutanix_category_key.os_type.id
	description = "testacc-ubuntu"
	value = "ubuntu"
}

data "nutanix_clusters" "clusters"{}

locals {
	cluster1 = [
	  for cluster in data.nutanix_clusters.clusters.entities :
	  cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_image" "acctest-test-categories" {
  name        = "Ubuntu-%d"
  description = "Ubuntu"

 categories {
	name  = nutanix_category_key.os_type.id
	value =	nutanix_category_value.ubuntu.id
 }

 categories {
	name  = nutanix_category_key.os_version.id
	value =	nutanix_category_value.os_version_value.id
 }

 cluster_references{
	 uuid = local.cluster1
 }

  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"

}
`, r)
}

func testAccNutanixImageConfigWithCategoriesUpdated(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "os_version"{
	name = "os_version"
	description = "testacc-os-version"
}

resource "nutanix_category_value" "os_version_value"{
	name = nutanix_category_key.os_version.id
	description = "testacc-os-current"
	value = "os_current"
}

resource "nutanix_category_value" "os_version_value_updated"{
	name = nutanix_category_key.os_version.id
	description = "testacc-ubuntu18"
	value = "18.08"
}

resource "nutanix_category_key" "os_type"{
	name = "os_type"
	description = "testacc-os-type"
}

resource "nutanix_category_value" "ubuntu"{
	name = nutanix_category_key.os_type.id
	description = "testacc-ubuntu"
	value = "ubuntu"
}

data "nutanix_clusters" "clusters"{}

locals {
	cluster1 = [
	  for cluster in data.nutanix_clusters.clusters.entities :
	  cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_image" "acctest-test-categories" {
  	name        = "Ubuntu-%d"
  	description = "Ubuntu"

	categories {
	   name  = nutanix_category_key.os_type.id
	   value =	nutanix_category_value.ubuntu.id
	}

	categories {
	   name  = nutanix_category_key.os_version.id
	   value = nutanix_category_value.os_version_value_updated.id
	}

	cluster_references{
		uuid = local.cluster1
	}

  	source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"

}
`, r)
}

func testAccNutanixImageConfigWithLargeImageURL(r int) string {
	return fmt.Sprintf(`
		resource "nutanix_image" "acctest-test" {
			name        = "Ubuntu-%d-server"
			description = "Ubuntu Server"
			source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
		}
	`, r)
}
