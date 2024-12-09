package volumesv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVolumeGroup = "nutanix_volume_group_v2.test"

func TestAccNutanixVolumeGroupV2Resource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		// CheckDestroy: testAccCheckNutanixVolumeGroupV4Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupResourceConfig(filepath, name, desc),
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
		},
	})
}

func TestAccNutanixVolumeGroupV2Resource_RequiredAttr(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupV2RequiredAttributes(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", name),
					testAndCheckComputedValues(resourceNameVolumeGroup),
				),
			},
		},
	})
}

func TestAccNutanixVolumeGroupV2Resource_WithNoName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVolumeGroupV2ConfigWithNoName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixVolumeGroupV2Resource_WithNoClusterReference(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-volume-group-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVolumeGroupV2ConfigWithNoClusterReference(name),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// VG just required attributes
func testAccVolumeGroupV2RequiredAttributes(name string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters" "clusters" {}

	locals{
		cluster1 = [
			for cluster in data.nutanix_clusters.clusters.entities :
			cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_volume_group_v2" "test" {
		name                               = "%s"
		cluster_reference                  = local.cluster1
	  }

`, name)
}

func testAccVolumeGroupV2ConfigWithNoName() string {
	return `
		data "nutanix_clusters" "clusters" {}
	
		locals{
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}
	
		resource "nutanix_volume_group_v2" "test" {
			cluster_reference                  = local.cluster1
		  }
	`
}

func testAccVolumeGroupV2ConfigWithNoClusterReference(name string) string {
	return fmt.Sprintf(`
	resource "nutanix_volume_group_v2" "test" {
		name                               = "%s"	
	  }
`, name)
}
