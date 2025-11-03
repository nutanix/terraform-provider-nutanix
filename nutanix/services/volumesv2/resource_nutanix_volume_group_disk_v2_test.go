package volumesv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceVolumeGroupDisk = "nutanix_volume_group_disk_v2.test"

func TestAccV2NutanixVolumeGroupDiskResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group disk description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDiskResourceConfig(filepath, name, desc, int(diskSizeBytes)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "description", desc),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "disk_size_bytes", strconv.Itoa(int(diskSizeBytes))),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "index", "1"),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "disk_storage_features.0.flash_mode.0.is_enabled", "false"),
				),
			},
		},
	})
}

func TestAccV2NutanixVolumeGroupDiskResource_BasicUpdateDiskSize(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group disk description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDiskResourceConfig(filepath, name, desc, int(diskSizeBytes)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "description", desc),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "disk_size_bytes", strconv.Itoa(int(diskSizeBytes))),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "index", "1"),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "disk_storage_features.0.flash_mode.0.is_enabled", "false"),
				),
			},
			{
				Config: testAccVolumeGroupsDiskResourceConfig(filepath, name, desc, int(updatedDiskSizeBytes)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "description", desc),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "disk_size_bytes", strconv.Itoa(int(updatedDiskSizeBytes))),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "index", "1"),
					resource.TestCheckResourceAttr(resourceVolumeGroupDisk, "disk_storage_features.0.flash_mode.0.is_enabled", "false"),
				),
			},
		},
	})
}

func testAccVolumeGroupsDiskResourceConfig(filepath, name, desc string, diskSizeBytes int) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes)
}
