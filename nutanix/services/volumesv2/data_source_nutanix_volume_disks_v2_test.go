package volumesv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupsDisks = "data.nutanix_volume_group_disks_v2.test"

func TestAccV2NutanixVolumeGroupsDisksDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDisksDataSourceConfig(filepath, name, desc, int(diskSizeBytes)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAttrListNotEmpty(dataSourceVolumeGroupsDisks, "disks", "index"),
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupsDisks, "disks.#"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.#", "2"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.0.disk_size_bytes", strconv.Itoa(int(diskSizeBytes))),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.0.disk_storage_features.0.flash_mode.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.0.description", desc),
				),
			},
		},
	})
}

func TestAccV2NutanixVolumeGroupsDisksDataSource_WithLimit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"

	limit := 1
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDisksDataSourceWithLimit(filepath, name, desc, int(diskSizeBytes), limit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAttrListNotEmpty(dataSourceVolumeGroupsDisks, "disks", "index"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.#", "1"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.0.disk_size_bytes", strconv.Itoa(int(diskSizeBytes))),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.0.disk_storage_features.0.flash_mode.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.0.description", desc),
				),
			},
		},
	})
}

func TestAccV2NutanixVolumeGroupsDisksDataSource_WithInvalidFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDisksDataSourceWithInvalidFilter(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisks, "disks.#", "0"),
				),
			},
		},
	})
}

func testAccVolumeGroupsDisksDataSourceConfig(filepath, name, desc string, diskSizeBytes int) string {
	return testAccVolumeGroupResourceConfig(name, desc) + testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		fmt.Sprintf(`


	resource "nutanix_volume_group_disk_v2" "test-2" {
		volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
		index               = 2
		description         = "%[1]s"
		disk_size_bytes     = %[2]d
		disk_data_source_reference {
		  name        = "terraform-test-disk_data_source_reference-disk-2"
		  ext_id      =  data.nutanix_storage_containers_v2.test.storage_containers[0].ext_id
		  entity_type = "STORAGE_CONTAINER"
		  uris        = ["uri3","uri4"]
		}
		disk_storage_features {
		  flash_mode {
			is_enabled = false
		  }
		}
		lifecycle {
		  ignore_changes = [
			disk_data_source_reference
		  ]
		}
		depends_on = [resource.nutanix_volume_group_v2.test]
	}

	data "nutanix_volume_group_disks_v2" "test" {
		volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
		depends_on          = [resource.nutanix_volume_group_v2.test ,resource.nutanix_volume_group_disk_v2.test, resource.nutanix_volume_group_disk_v2.test-2]
	}
	`, desc, diskSizeBytes)
}

func testAccVolumeGroupsDisksDataSourceWithLimit(filepath, name, desc string, diskSizeBytes int, limit int) string {
	return testAccVolumeGroupResourceConfig(name, desc) + testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		fmt.Sprintf(`

	  	resource "nutanix_volume_group_disk_v2" "test-2" {
		    volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			index               = 2
			description         = "%[1]s"
			disk_size_bytes     = %[2]d
			disk_data_source_reference {
			  name        = "terraform-test-disk_data_source_reference-disk-2"
			  ext_id      = data.nutanix_storage_containers_v2.test.storage_containers[0].ext_id
			  entity_type = "STORAGE_CONTAINER"
			  uris        = ["uri3","uri4"]
			}
			disk_storage_features {
			  flash_mode {
				is_enabled = false
			  }
			}
			lifecycle {
			  ignore_changes = [
				disk_data_source_reference
			  ]
			}
			depends_on = [resource.nutanix_volume_group_v2.test]
		}

		data "nutanix_volume_group_disks_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			limit  = %[3]d
			depends_on          = [resource.nutanix_volume_group_v2.test ,resource.nutanix_volume_group_disk_v2.test, resource.nutanix_volume_group_disk_v2.test-2]
		}

		`, desc, diskSizeBytes, limit,
		)
}

func testAccVolumeGroupsDisksDataSourceWithInvalidFilter(name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		`
		data "nutanix_volume_group_disks_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			filter = "storageContainerId eq 'invalid'"
		}
	`
}
