package volumesv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupsDisks = "data.nutanix_volume_group_disks_v2.test"

func TestAccNutanixVolumeGroupsDisksV2DataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{

		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDisksDataSourceConfig(filepath, name, desc),
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

func TestAccNutanixVolumeGroupsDisksV2DataSource_WithLimit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	limit := 1
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDisksDataSourceWithLimit(filepath, name, desc, limit),
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

func testAccVolumeGroupsDisksDataSourceConfig(filepath, name, desc string) string {
	return testAccVolumeGroupResourceConfig(filepath, name, desc) + testAccVolumeGroupDiskResourceConfig(filepath, name, desc) +
		fmt.Sprintf(`
	resource "nutanix_volume_group_disk_v2" "test-2" {
		volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
		index               = 2
		description         = "%[1]s"
		disk_size_bytes     = %[2]d
		disk_data_source_reference {
		  name        = "terraform-test-disk_data_source_reference-disk-2"
		  ext_id      = local.vg_disk.disk_data_source_reference.ext_id
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

func testAccVolumeGroupsDisksDataSourceWithLimit(filepath, name, desc string, limit int) string {
	return testAccVolumeGroupResourceConfig(filepath, name, desc) + testAccVolumeGroupDiskResourceConfig(filepath, name, desc) +
		fmt.Sprintf(`	  
	  	resource "nutanix_volume_group_disk_v2" "test-2" {
		    volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			index               = 2
			description         = "%[1]s"
			disk_size_bytes     = %[2]d		  
			disk_data_source_reference {
			  name        = "terraform-test-disk_data_source_reference-disk-2"
			  ext_id      = local.vg_disk.disk_data_source_reference.ext_id
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
