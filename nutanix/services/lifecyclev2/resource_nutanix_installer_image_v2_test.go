package lifecyclev2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Test plan for nutanix_installer_image_v2 (lifecycle / InstallerImages)
// ----------------------------------------------------------------------
// The resource wraps the asynchronous v4 InstallerImages API
// (CreateImage / GetImageById / UpdateImageById / DeleteImageById) plus the
// singular (GetImageById) and plural (ListImages) datasources. Coverage:
//
//  1. TestAccV2NutanixInstallerImageResource_Basic
//     - Step 1 (Create): register a REMOTE_URL image with every writable field
//       populated (name, type, source, url, version). Assert each literal value
//       on the resource, and assert the server-computed attributes (ext_id,
//       file_status) are present.
//     - Step 2 (Update): change name, version and url and assert the new literal
//       values are reflected in state.
//     - Step 3 (Datasource singular): read the image back by ext_id and assert
//       every attribute matches the resource via TestCheckResourceAttrPair.
//     - Step 4 (Datasource plural): read the image list, assert count > 0 and
//       spot-check nested attributes.
//     - Step 5 (Import): import the resource by ext_id and verify state.
//
//  2. TestAccV2NutanixInstallerImageResource_InvalidType (negative)
//     - Supply an invalid `type` enum and assert the schema ValidateFunc rejects it.
//
//  3. TestAccV2NutanixInstallerImagesDataSource_WithInvalidFilter (negative)
//     - Supply a malformed filter and assert the list datasource errors.
//
// CheckDestroy calls GetImageById for the primary resource and asserts it errors
// (the image no longer exists) — see testAccCheckNutanixInstallerImageDestroy.

const (
	resourceNameInstallerImage        = "nutanix_installer_image_v2.test"
	datasourceNameInstallerImage      = "data.nutanix_installer_image_v2.test"
	datasourceNameInstallerImagesList = "data.nutanix_installer_images_v2.list"
)

func TestAccV2NutanixInstallerImageResource_Basic(t *testing.T) {
	// The API enforces a maximum name length of 40 characters, so keep names short.
	r := acctest.RandIntRange(0, 100000)
	name := fmt.Sprintf("tf-img-%d", r)
	nameUpdated := fmt.Sprintf("tf-img-upd-%d", r)
	imageType := testVars.Lifecycle.InstallerImage.Type
	url := testVars.Lifecycle.InstallerImage.URL
	urlUpdated := testVars.Lifecycle.InstallerImage.UpdatedURL
	version := testVars.Lifecycle.InstallerImage.Version
	versionUpdated := version + "-b"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixInstallerImageDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with all writable fields populated.
			{
				Config: testInstallerImageConfig(name, imageType, url, version),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "type", imageType),
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "source", "REMOTE_URL"),
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "url", url),
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "version", version),
					resource.TestCheckResourceAttrSet(resourceNameInstallerImage, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameInstallerImage, "file_status"),
				),
			},
			// Step 2: Update name, version and url; assert new literal values.
			{
				Config: testInstallerImageConfig(nameUpdated, imageType, urlUpdated, versionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "url", urlUpdated),
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "version", versionUpdated),
					// unchanged fields remain stable
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "type", imageType),
					resource.TestCheckResourceAttr(resourceNameInstallerImage, "source", "REMOTE_URL"),
				),
			},
			// Step 3: Read back via the singular datasource and match every attribute.
			{
				Config: testInstallerImageWithSingularDatasourceConfig(nameUpdated, imageType, urlUpdated, versionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "name", resourceNameInstallerImage, "name"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "type", resourceNameInstallerImage, "type"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "source", resourceNameInstallerImage, "source"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "url", resourceNameInstallerImage, "url"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "version", resourceNameInstallerImage, "version"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "ext_id", resourceNameInstallerImage, "ext_id"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "file_status", resourceNameInstallerImage, "file_status"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "metadata_status", resourceNameInstallerImage, "metadata_status"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "metadata_download_url", resourceNameInstallerImage, "metadata_download_url"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "certificate_chain", resourceNameInstallerImage, "certificate_chain"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "created_time", resourceNameInstallerImage, "created_time"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "tenant_id", resourceNameInstallerImage, "tenant_id"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "checksum.#", resourceNameInstallerImage, "checksum.#"),
					resource.TestCheckResourceAttrPair(datasourceNameInstallerImage, "links.#", resourceNameInstallerImage, "links.#"),
				),
			},
			// Step 4: Read the list datasource and spot-check.
			{
				Config: testInstallerImageWithPluralDatasourceConfig(nameUpdated, imageType, urlUpdated, versionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameInstallerImagesList, "images.#"),
					resource.TestCheckResourceAttrSet(datasourceNameInstallerImagesList, "images.0.ext_id"),
				),
			},
			// Step 5: Import the resource by ext_id and verify state.
			{
				ResourceName:      resourceNameInstallerImage,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccV2NutanixInstallerImageResource_InvalidType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testInstallerImageConfig("tf-invalid-type", "INVALID_TYPE", "https://example.com/img.tar.gz", "1.0"),
				ExpectError: regexp.MustCompile(`expected type to be one of`),
			},
		},
	})
}

func TestAccV2NutanixInstallerImagesDataSource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testInstallerImagesInvalidFilterConfig(),
				ExpectError: regexp.MustCompile(`error while fetching images`),
			},
		},
	})
}

// testAccCheckNutanixInstallerImageDestroy verifies that after the test the
// installer image no longer exists — GetImageById must return an error.
func testAccCheckNutanixInstallerImageDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).LifecycleAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_installer_image_v2" {
			continue
		}
		if _, err := conn.InstallerImagesAPIInstance.GetImageById(utils.StringPtr(rs.Primary.ID)); err == nil {
			return fmt.Errorf("installer image (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testInstallerImageConfig(name, imageType, url, version string) string {
	return fmt.Sprintf(`
resource "nutanix_installer_image_v2" "test" {
  name    = "%[1]s"
  type    = "%[2]s"
  source  = "REMOTE_URL"
  url     = "%[3]s"
  version = "%[4]s"
}
`, name, imageType, url, version)
}

func testInstallerImageWithSingularDatasourceConfig(name, imageType, url, version string) string {
	return testInstallerImageConfig(name, imageType, url, version) + `
data "nutanix_installer_image_v2" "test" {
  ext_id = nutanix_installer_image_v2.test.ext_id
}
`
}

func testInstallerImageWithPluralDatasourceConfig(name, imageType, url, version string) string {
	return testInstallerImageConfig(name, imageType, url, version) + `
data "nutanix_installer_images_v2" "list" {
  depends_on = [nutanix_installer_image_v2.test]
}
`
}

func testInstallerImagesInvalidFilterConfig() string {
	return `
data "nutanix_installer_images_v2" "invalid" {
  filter = "this is not a valid filter"
}
`
}
