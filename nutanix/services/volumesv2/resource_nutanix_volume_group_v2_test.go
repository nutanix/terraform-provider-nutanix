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

// TestAccV2NutanixVolumeGroupResource_Update exercises UpdateVolumeGroupById:
// it creates a Volume Group and then updates mutable fields (name, description,
// sharing_status, usage_type) in a subsequent step, asserting the new literal
// values are round-tripped back through the Read path.
func TestAccV2NutanixVolumeGroupResource_Update(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group description"
	updatedName := fmt.Sprintf("tf-test-volume-group-updated-%d", r)
	updatedDesc := "test volume group description updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupV2UpdateConfig(name, desc, "SHARED", "USER", "ISCSI", "DIRECT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "sharing_status", "SHARED"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "usage_type", "USER"),
					resource.TestCheckResourceAttrSet(resourceNameVolumeGroup, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "protocol", "ISCSI"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "attachment_type", "DIRECT"),
				),
			},
			{
				Config: testAccVolumeGroupV2UpdateConfig(updatedName, updatedDesc, "NOT_SHARED", "TEMPORARY", "NVMF", "EXTERNAL"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "sharing_status", "NOT_SHARED"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "usage_type", "TEMPORARY"),
					resource.TestCheckResourceAttrSet(resourceNameVolumeGroup, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "protocol", "NVMF"),
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

	gb := 1024 * 1024 * 1024
	disk1Size := 10 * gb
	disk1 := fmt.Sprintf(`
		disks {
			disk_size_bytes = %[2]d
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
	`, name, disk1Size)
	disk2 := fmt.Sprintf(`
		disks {
			disk_size_bytes = %[2]d
			index = 2
			disk_data_source_reference {
			  name        = "vg-disk-2-%[1]s"
			  ext_id      = data.nutanix_storage_containers_v2.test.storage_containers[0].ext_id
			  entity_type = "STORAGE_CONTAINER"
			  uris        = ["uri3","uri4"]
			}
			disk_storage_features {
				flash_mode {
					is_enabled = false
				}
			}
		}
	`, name, 20*gb)
	disk1Updated := fmt.Sprintf(`
		disks {
			disk_size_bytes = %[2]d
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
	`, name, 15*gb)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupResourceConfigWithAttachmentTypeAndProtocolAndDisks(name, desc, disk1),
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
			{
				Config: testAccVolumeGroupResourceConfigWithAttachmentTypeAndProtocolAndDisksUpdate(name, desc, disk1, disk2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "should_load_balance_vm_attachments", "false"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "sharing_status", "SHARED"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "attachment_type", "EXTERNAL"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "protocol", "NVMF"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_size_bytes", "10737418240"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.index", "1"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_storage_features.0.flash_mode.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.1.disk_size_bytes", "21474836480"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.1.index", "2"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.1.disk_storage_features.0.flash_mode.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "is_hidden", "false"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "usage_type", "USER"),
				),
			},
			{
				Config: testAccVolumeGroupResourceConfigWithAttachmentTypeAndProtocolAndDisks(name, desc, disk1Updated),
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
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_size_bytes", "16106127360"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.index", "1"),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_storage_features.0.flash_mode.0.is_enabled", "false"),
				),
			},
		},
	})
}

// TestAccV2NutanixVolumeGroupResource_UpdateDiskSize verifies that increasing
// disk_size_bytes on an existing inline disk is actually applied to the Volume
// Group (previously the update was a no-op: apply succeeded but the disk was
// never resized and state silently reported the new size).
func TestAccV2NutanixVolumeGroupResource_UpdateDiskSize(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group disk resize"

	const initialSize = "10737418240" // 10 GiB
	const updatedSize = "21474836480" // 20 GiB

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupResourceConfigWithDiskSize(name, desc, "10"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_size_bytes", initialSize),
				),
			},
			{
				Config: testAccVolumeGroupResourceConfigWithDiskSize(name, desc, "20"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "name", name),
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_size_bytes", updatedSize),
				),
			},
		},
	})
}

// TestAccV2NutanixVolumeGroupResource_UpdateDiskSizeShrinkFails verifies that
// attempting to shrink an existing disk is rejected, matching the fabric
// constraint that disk size can only be increased.
func TestAccV2NutanixVolumeGroupResource_UpdateDiskSizeShrinkFails(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group disk shrink"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupResourceConfigWithDiskSize(name, desc, "20"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVolumeGroup, "disks.0.disk_size_bytes", "21474836480"),
				),
			},
			{
				Config:      testAccVolumeGroupResourceConfigWithDiskSize(name, desc, "10"),
				ExpectError: regexp.MustCompile("disk size can only be increased"),
			},
		},
	})
}

func testAccVolumeGroupResourceConfigWithDiskSize(name, desc, sizeGiB string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster1 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	data "nutanix_storage_containers_v2" "test" {
	  filter = "clusterExtId eq '${local.cluster1}' and startswith(name,'default-container-')"
	  limit  = 1
	}

	resource "nutanix_volume_group_v2" "test" {
		name              = "%[1]s"
		description       = "%[2]s"
		cluster_reference = local.cluster1

		disks {
			disk_size_bytes = %[3]s * 1024 * 1024 * 1024
			index           = 1
			disk_data_source_reference {
			  name        = "vg-disk-%[1]s"
			  ext_id      = data.nutanix_storage_containers_v2.test.storage_containers[0].ext_id
			  entity_type = "STORAGE_CONTAINER"
			}
		}
	  }
	`, name, desc, sizeGiB)
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

func TestAccV2NutanixVolumeGroupResource_ProjectAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-vg-projassoc-%d", r)
	desc := "volume group project association test"
	projectName := fmt.Sprintf("tf-vg-pa-proj-%d", r)
	dataSourceNameVolumeGroupsV2 := "data.nutanix_volume_groups_v2.list_volume_groups"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVolumeGroupProjectAssociationConfig(name, desc, projectName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("nutanix_volume_group_v2.create", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_volume_group_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameVolumeGroupsV2, "volumes.#"),
					resource.TestCheckResourceAttr(dataSourceNameVolumeGroupsV2, "volumes.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceNameVolumeGroupsV2, "volumes.0.project_ext_id", "nutanix_project_v2.test", "ext_id"),
				),
			},
			{
				Config:      testVolumeGroupProjectAssociationConfig(name, desc, projectName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
			},
		},
	})
}

func vgProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testVolumeGroupProjectAssociationConfig(name, desc, projectName, projectExtIDOverride string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}
	data "nutanix_storage_containers_v2" "storage_container" {}
	locals {
		cluster0 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
		storage_container0 = [
			for storage_container in data.nutanix_storage_containers_v2.storage_container.storage_containers :
			storage_container.ext_id if storage_container.name == "SelfServiceContainer"
		][0]
	}

	resource "nutanix_project_v2" "test" {
		name        = "%[3]s"
		project_id  = "%[3]s"
		description = "project association test"
	}

	resource "nutanix_resource_group_v2" "rg" {
		name           = "tf-vg-pa-rg-%[3]s"
		project_ext_id = nutanix_project_v2.test.ext_id
		placement_targets {
			cluster_ext_id = local.cluster0
			storage_containers {
				ext_id = local.storage_container0
			}
		}
	}

	resource "nutanix_volume_group_v2" "create" {
		name              = "%[1]s"
		description       = "%[2]s"
		cluster_reference = local.cluster0
		%[4]s
		depends_on = [nutanix_project_v2.test, nutanix_resource_group_v2.rg]
	}

	data "nutanix_volume_group_v2" "test" {
		ext_id     = nutanix_volume_group_v2.create.id
		depends_on = [nutanix_volume_group_v2.create]
	}

	data "nutanix_volume_groups_v2" "list_volume_groups" {
		filter = "projectExtId eq '${nutanix_project_v2.test.ext_id}'"
		depends_on = [nutanix_volume_group_v2.create]
	}
	`, name, desc, projectName, vgProjectExtIDLine(projectExtIDOverride))
}

// testAccVolumeGroupV2UpdateConfig builds a Volume Group config parameterized by
// the mutable fields exercised by the UpdateVolumeGroupById path.
func testAccVolumeGroupV2UpdateConfig(name, desc, sharingStatus, usageType, protocol, attachmentType string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster1 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_volume_group_v2" "test" {
		name              = "%[1]s"
		description       = "%[2]s"
		cluster_reference = local.cluster1
		sharing_status    = "%[3]s"
		usage_type        = "%[4]s"
		created_by        = "admin"
		protocol          = "%[5]s"
		attachment_type   = "%[6]s"
	}
	`, name, desc, sharingStatus, usageType, protocol, attachmentType)
}

func testAccVolumeGroupResourceConfigWithAttachmentTypeAndProtocolAndDisks(name string, desc string, disk1 string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster1 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

    data "nutanix_storage_containers_v2" "test" {
	  filter = "clusterExtId eq '${local.cluster1}' and startswith(name,'default-container-')"
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
		%[3]s
		is_hidden = false
		lifecycle {
			ignore_changes = [
			  iscsi_features[0].target_secret
			]
		}
	  }
	`, name, desc, disk1)
}

func testAccVolumeGroupResourceConfigWithAttachmentTypeAndProtocolAndDisksUpdate(name string, desc string, disk1 string, disk2 string) string {

	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster1 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

    data "nutanix_storage_containers_v2" "test" {
	  filter = "clusterExtId eq '${local.cluster1}' and startswith(name,'default-container-')"
	  limit  = 1
    }


	resource "nutanix_volume_group_v2" "test" {
		name                               = "%[1]s"
		description                        = "%[2]s"
		should_load_balance_vm_attachments = false
		sharing_status                     = "SHARED"
		created_by 						   = "admin"
		cluster_reference                  = local.cluster1
		attachment_type = "EXTERNAL"
		protocol = "NVMF"

		%[3]s

		%[4]s

		is_hidden = false
		lifecycle {
			ignore_changes = [
			  iscsi_features[0].target_secret
			]
		}
	  }
	`, name, desc, disk1, disk2)
}
