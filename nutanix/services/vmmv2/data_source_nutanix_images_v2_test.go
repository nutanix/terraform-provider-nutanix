package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameImages = "data.nutanix_images_v2.test"

func TestAccV2NutanixImagesDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagePreConfigV2(name, desc) + testAccImagesDataSourceConfigV2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesDatasource_WithFilters(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagePreConfigV2(name, desc) + testAccImagesDataSourceConfigV2WithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.type"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.description"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.last_update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.owner_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.placement_policy_status.#"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.cluster_location_ext_ids.#"),
					resource.TestCheckResourceAttrSet(datasourceNameImages, "images.0.source.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesDataSourceConfigV2WithInvalidFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameImages, "images.#", "0"),
				),
			},
		},
	})
}

func testImagePreConfigV2(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_images_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			type = "ISO_IMAGE"
			source{
				url_source{
					url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
				}
			}
		}
	`, name, desc)
}

func testAccImagesDataSourceConfigV2() string {
	return `
		data "nutanix_images_v2" "test"{
			depends_on = [
				resource.nutanix_images_v2.test
			]
		}
	`
}

func testAccImagesDataSourceConfigV2WithFilters() string {
	return `

		data "nutanix_images_v2" "test"{
			page=0
			limit=10
			filter="name eq '${nutanix_images_v2.test.name}'"
		}
`
}

func testAccImagesDataSourceConfigV2WithInvalidFilters() string {
	return `

		data "nutanix_images_v2" "test"{
			filter="name eq 'invalid-name'"
		}
`
}
