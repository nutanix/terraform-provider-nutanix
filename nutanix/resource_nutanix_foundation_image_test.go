package nutanix

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func GetResourceState(s *terraform.State, key string) map[string]string {
	moduleState := s.RootModule()
	return moduleState.Resources[key].Primary.Attributes
}

func testAccCheckNosImageExists(filename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*Client).FoundationClientAPI
		ctx := context.TODO()
		resp, err := conn.FileManagement.ListNOSPackages(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch nos packages from FVM")
		}

		for _, v := range *resp {
			if v == filename {
				return nil
			}
		}
		return fmt.Errorf("upload for nos package %s failed. Image not found in FVM", filename)
	}
}

func testAccCheckNosImageDestroy(filename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*Client).FoundationClientAPI
		ctx := context.TODO()
		resp, err := conn.FileManagement.ListNOSPackages(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch nos packages from FVM")
		}

		for _, v := range *resp {
			if v == filename {
				return fmt.Errorf("teraform destroy for nos package %s failed. It still exists in FVM", filename)
			}
		}
		return nil
	}
}

func TestAccFoundationImageResource_NOSUpload(t *testing.T) {
	nameForUpload := "nos_upload"
	resourcePathForUpload := "nutanix_foundation_image." + nameForUpload
	r := acctest.RandIntRange(0, 500)
	filename := fmt.Sprintf("test_nos_image-%d.tar.gz", r)
	nosFile := "test_nos_image.tar.gz"

	// get file file path to config having nodes info
	path, _ := os.Getwd()
	filepath := path + "/../test_files/" + nosFile

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccFoundationPreCheck(t) },
		CheckDestroy: testAccCheckNosImageDestroy(filename),
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImageResourceUpload(nameForUpload, filename, "nos", filepath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNosImageExists(filename),
					resource.TestCheckResourceAttrSet(resourcePathForUpload, "md5sum"),
					resource.TestCheckResourceAttrSet(resourcePathForUpload, "in_whitelist"),
					resource.TestCheckResourceAttrSet(resourcePathForUpload, "name"),
				),
			},
		},
	})
}

// Check negative scenario incase the resource errors out for incorrect installer type
func TestAccFoundationImageResource_Error(t *testing.T) {
	nameForUpload := "iso_upload"
	r := acctest.RandIntRange(0, 500)
	filename := fmt.Sprintf("test_alpine-%d.iso", r)
	file := "alpine.iso"

	// get file file path to config having nodes info
	path, _ := os.Getwd()
	filepath := path + "/../test_files/" + file

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccFoundationPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testImageResourceUpload(nameForUpload, filename, "check", filepath),
				ExpectError: regexp.MustCompile("installer_type check should be one of "),
			},
		},
	})
}

func testImageResourceUpload(name, filename, instType, filepath string) string {
	return fmt.Sprintf(`
	resource "nutanix_foundation_image" "%s"{
		filename = "%s"
		installer_type = "%s"
		source = "%s"
	}
	`, name, filename, instType, filepath)
}
