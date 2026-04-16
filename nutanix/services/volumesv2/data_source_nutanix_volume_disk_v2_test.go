package volumesv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupsDisk = "data.nutanix_volume_group_disk_v2.test"

func TestAccV2NutanixVolumeGroupsDiskDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-disk-%d", r)
	desc := "terraform test volume group disk description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDiskDataSourceConfig(filepath, name, desc, int(diskSizeBytes)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisk, "index", "1"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisk, "disk_size_bytes", strconv.Itoa(int(diskSizeBytes))),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisk, "disk_storage_features.0.flash_mode.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceVolumeGroupsDisk, "description", desc),
				),
			},
		},
	})
}

func testAccVolumeGroupsDiskDataSourceConfig(filepath, name, desc string, diskSizeBytes int) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		`		  
		  data "nutanix_volume_group_disk_v2" "test" {
			volume_group_ext_id = nutanix_volume_group_v2.test.id
			ext_id              = resource.nutanix_volume_group_disk_v2.test.id
		    depends_on = [resource.nutanix_volume_group_disk_v2.test]
		  }	  		  
	`
}
