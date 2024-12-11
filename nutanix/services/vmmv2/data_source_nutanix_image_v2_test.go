package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameImage = "data.nutanix_image_v2.test"

func TestAccV2NutanixImageDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceConfigV2(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameImage, "name", name),
					resource.TestCheckResourceAttr(datasourceNameImage, "type", "ISO_IMAGE"),
					resource.TestCheckResourceAttr(datasourceNameImage, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "placement_policy_status.#"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "cluster_location_ext_ids.#"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "source.#"),
				),
			},
		},
	})
}

func testAccImageDataSourceConfigV2(name, desc string) string {
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

		data "nutanix_image_v2" "test"{
			ext_id = resource.nutanix_images_v2.test.id
		}
`, name, desc)
}
