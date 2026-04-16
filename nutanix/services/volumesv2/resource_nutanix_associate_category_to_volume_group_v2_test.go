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

	dataSourceAssociatedCategory := "data.nutanix_category_v2.associated_vg"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAssociateCategoryToVolumeGroupResourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "description", desc),
					resource.TestCheckResourceAttrSet(resourceAssociateCategoryToVolumeGroup, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceAssociateCategoryToVolumeGroup, "categories.0.ext_id"),
					resource.TestCheckResourceAttr(resourceAssociateCategoryToVolumeGroup, "categories.0.entity_type", "CATEGORY"),
					resource.TestCheckResourceAttr(dataSourceAssociatedCategory, "associations.0.resource_type", "VOLUMEGROUP"),
				),
			},
		},
	})
}

func testAccAssociateCategoryToVolumeGroupResourceConfig(name, desc string) string {
	return fmt.Sprintf(`
#pull cluster data
data "nutanix_clusters_v2" "clusters" {}

#pull desired cluster data from setup
locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

// Create a volume group
resource "nutanix_volume_group_v2" "test" {
  name                               = "%[1]s"
  description                        = "%[2]s"
  should_load_balance_vm_attachments = false
  created_by                         = "example"
  cluster_reference                  = local.cluster_ext_id
}



#creating category
resource "nutanix_category_v2" "vg-category" {
  key         = "%[1]s_vg_category"
  value       = "category_example_value"
  description = "category example to associate with volume group"
}


# Associate categories to volume group
resource "nutanix_associate_category_to_volume_group_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  categories {
    ext_id = nutanix_category_v2.vg-category.id
  }
}

# pull associated category data
data "nutanix_category_v2" "associated_vg" {
  ext_id     = nutanix_category_v2.vg-category.id
  expand     = "associations"
  depends_on = [nutanix_associate_category_to_volume_group_v2.test]
}


`, name, desc)
}
