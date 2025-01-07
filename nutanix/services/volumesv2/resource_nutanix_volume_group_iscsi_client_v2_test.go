package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceVolumeGroupIscsiClient = "nutanix_volume_group_iscsi_client_v2.test"

func TestAccV2NutanixVolumeGroupIscsiClientResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-volume-group-%d", r)
	desc := "test volume group Iscsi Client description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupResourceConfig(name, desc) + testAccVolumeGroupIscsiClientResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceVolumeGroupIscsiClient, "ext_id"),
				),
			},
		},
	})
}

func testAccVolumeGroupIscsiClientResourceConfig() string {
	return `	
		data "nutanix_volume_iscsi_clients_v2" "test" {}
		resource "nutanix_volume_group_iscsi_client_v2" "test" {
			vg_ext_id = resource.nutanix_volume_group_v2.test.id
			ext_id     =  data.nutanix_volume_iscsi_clients_v2.test.iscsi_clients.0.ext_id
			iscsi_initiator_name = data.nutanix_volume_iscsi_clients_v2.test.iscsi_clients.0.iscsi_initiator_name
			depends_on = [ resource.nutanix_volume_group_v2.test ]
		}		
	`
}
