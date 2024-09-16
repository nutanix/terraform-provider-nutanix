package volumesv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceVolumeGroupVm = "nutanix_volume_group_vm_v2.test"

func TestAccNutanixVolumeGroupVmV2_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-volume-group-%d", r)
	desc := "test volume group Vm Attachment description"
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v4.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupVmConfig(filepath, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVolumeGroupVm, "vm_ext_id", testVars.Volumes.VmExtId),
				),
			},
		},
	})
}

func testAccVolumeGroupVmConfig(filepath, name, desc string) string {

	return testAccVolumeGroupResourceConfig(filepath, name, desc) + `		  
		  resource "nutanix_volume_group_vm_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			vm_ext_id           = local.volumes.vm_ext_id
			depends_on          = [resource.nutanix_volume_group_v2.test]
		  }
		
	`
}
