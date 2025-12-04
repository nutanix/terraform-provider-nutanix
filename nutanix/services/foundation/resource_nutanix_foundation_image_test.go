//go:build !unit

package foundation_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func GetResourceState(s *terraform.State, key string) map[string]string {
	moduleState := s.RootModule()
	return moduleState.Resources[key].Primary.Attributes
}

func testAccCheckNosImageExists(filename string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client).FoundationClientAPI
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
		conn := acc.TestAccProvider.Meta().(*conns.Client).FoundationClientAPI
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

	// Get the Working directory
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("TestAccFoundationImageResource_NOSUpload failed to get working directory %s", err)
	}

	filepath := fmt.Sprintf("%s/%s", dir, nosFile)

	defer os.Remove(filepath)

	// get image url from env variables
	image := os.Getenv("NOS_IMAGE_TEST_URL")
	if image == "" {
		t.Fatal("NOS_IMAGE_TEST_URL is empty. Please set env variable NOS_IMAGE_TEST_URL")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if err := downloadFile(filepath, image); err != nil {
				t.Errorf("TestAccFoundationImageResource_NOSUpload failed to download image %s", err)
			}
			acc.TestAccFoundationPreCheck(t)
		},
		CheckDestroy: testAccCheckNosImageDestroy(filename),
		Providers:    acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImageResourceUpload(nameForUpload, filename, "nos", filepath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNosImageExists(filename),
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

	// Get the Working directory
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("TestAccFoundationImageResource_NOSUpload failed to get working directory %s", err)
	}

	filepath := fmt.Sprintf("%s/%s", dir, file)

	defer os.Remove(filepath)

	// get image url from env variables
	image := "http://dl-cdn.alpinelinux.org/alpine/v3.8/releases/x86_64/alpine-virt-3.8.1-x86_64.iso"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if err := downloadFile(filepath, image); err != nil {
				t.Errorf("TestAccFoundationImageResource_Error failed to download image %s", err)
			}
			acc.TestAccFoundationPreCheck(t)
		},
		Providers: acc.TestAccProviders,
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
