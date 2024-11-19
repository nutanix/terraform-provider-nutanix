package volumesv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceVolumeGroupDisk = "nutanix_volume_group_disk_v2.test"

func TestAccNutanixVolumeGroupDiskV2_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-volume-group-%d", r)
	desc := "test volume group disk description"
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDiskResourceConfig(filepath, name, desc),
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

func testAccVolumeGroupsDiskResourceConfig(filepath, name, desc string) string {
	return testAccVolumeGroupResourceConfig(filepath, name, desc) + testAccVolumeGroupDiskResourceConfig(filepath, name, desc)
}
