package volumesv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceVolumeGroupIscsiClient = "nutanix_volume_group_iscsi_client_v2.test"

func TestAccNutanixVolumeGroupIscsiClientV2_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-volume-group-%d", r)
	desc := "test volume group Iscsi Client description"
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupResourceConfig(filepath, name, desc) + testAccVolumeGroupIscsiClientResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVolumeGroupIscsiClient, "ext_id", testVars.Volumes.IscsiClient.ExtId),
				),
			},
		},
	})
}

func testAccVolumeGroupIscsiClientResourceConfig() string {

	return `		
		resource "nutanix_volume_group_iscsi_client_v2" "test" {
			vg_ext_id = resource.nutanix_volume_group_v2.test.id
			ext_id     = local.volumes.iscsi_client.ext_id
			iscsi_initiator_name = local.volumes.iscsi_client.initiator_name
			depends_on = [ resource.nutanix_volume_group_v2.test ]
		}		
	`
}
