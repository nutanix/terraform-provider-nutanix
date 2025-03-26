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

func TestAccV2NutanixVolumeGroupResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
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
		},
	})
}

func TestAccV2NutanixVolumeGroupResource_RequiredAttr(t *testing.T) {
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

func TestAccV2NutanixVolumeGroupResource_WithNoName(t *testing.T) {
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

func TestAccV2NutanixVolumeGroupResource_WithNoClusterReference(t *testing.T) {
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

func TestAccV2NutanixVolumeGroupResource_WithAttachmentTypeAndProtocolAndDisks(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group description with attachment type and protocol and disks"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupResourceConfigWithAttachmentTypeAndProtocolAndDisks(name, desc),
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
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "attachment_type", "DIRECT"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "protocol", "ISCSI"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_size_bytes", "10737418240"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.index", "1"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_storage_features.0.flash_mode.0.is_enabled", "false"),
				),
			},
		},
	})
}

// VG just required attributes
func testAccVolumeGroupV2RequiredAttributes(name string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals{
		cluster1 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
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
		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster1 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
					cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
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

func testAccVolumeGroupResourceConfigWithAttachmentTypeAndProtocolAndDisks(name string, desc string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster1 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

    data "nutanix_storage_containers_v2" "test" {
	  filter = "clusterExtId eq '${local.cluster1}'"
	  limit  = 1
    }

	resource "nutanix_volume_group_v2" "test" {
		name                               = "%[1]s"
		description                        = "%[2]s"
		should_load_balance_vm_attachments = false
		sharing_status                     = "SHARED"
		created_by 						   = "admin"
		cluster_reference                  = local.cluster1
		iscsi_features {
			target_secret			 = "1234567891011"
			enabled_authentications  = "CHAP"
		}
		storage_features {
		  flash_mode {
			is_enabled = true
		  }
		}
		usage_type = "USER"
		attachment_type = "DIRECT"
		protocol = "ISCSI"
		disks {
			disk_size_bytes = 10 * 1024 * 1024 * 1024
			index = 1
			disk_data_source_reference {
			  name        = "vg-disk-%[1]s"
			  ext_id      = data.nutanix_storage_containers_v2.test.storage_containers[0].ext_id
			  entity_type = "STORAGE_CONTAINER"
			  uris        = ["uri1","uri2"]
			}
			disk_storage_features {
				flash_mode {
					is_enabled = false
				}
			}
		}
		is_hidden = false
		lifecycle {
			ignore_changes = [
			  iscsi_features[0].target_secret
			]
		}
	  }
	`, name, desc)
}
