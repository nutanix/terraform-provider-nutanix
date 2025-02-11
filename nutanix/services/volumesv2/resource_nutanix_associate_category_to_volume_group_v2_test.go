package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceAssociateCategoryToVolumeGroup = "nutanix_associate_category_to_volume_group_v2.test"

func TestAccV2NutanixAssociateCategoryToVolumeGroupResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group to associate category"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
			// Create a volume group
			{
				Config: testAccVolumeGroupResourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "should_load_balance_vm_attachments", "false"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "sharing_status", "SHARED"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "created_by", "admin"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "iscsi_features.0.enabled_authentications", "CHAP"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "storage_features.0.flash_mode.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "is_hidden", "false"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "usage_type", "USER"),
				),
			},
			// Associate Category to Volume Group
			{
				Config: testAccVolumeGroupResourceConfig(name, desc) +
					testAccAssociateCategoryToVolumeGroupResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceAssociateCategoryToVolumeGroup, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceAssociateCategoryToVolumeGroup, "categories.0.ext_id"),
					resource.TestCheckResourceAttr(resourceAssociateCategoryToVolumeGroup, "categories.0.entity_type", "CATEGORY"),
				),
			},
		},
	})
}

func testAccAssociateCategoryToVolumeGroupResourceConfig() string {
	return `
# List categories
data "nutanix_categories_v2" "categories" {}

resource "nutanix_associate_category_to_volume_group_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.0.ext_id
  }
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.1.ext_id
  }
  categories{
    ext_id = data.nutanix_categories_v2.categories.categories.2.ext_id
  }
}

`
}
