package volumesv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupIscsiClients = "data.nutanix_volume_group_iscsi_clients_v2.vg_iscsi_test"

func TestAccNutanixVolumeGroupIscsiClientsV2_Basic(t *testing.T) {
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
				Config: testAccVolumeGroupIscsiClientsV2Config(filepath, name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupIscsiClients, "iscsi_clients.#"),
					testAccCheckResourceAttrListNotEmpty(dataSourceVolumeGroupIscsiClients, "iscsi_clients", "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupIscsiClients, "iscsi_clients.0.cluster_reference"),
					resource.TestCheckResourceAttrSet(dataSourceVolumeGroupIscsiClients, "iscsi_clients.0.ext_id"),
				),
			},
		},
	})
}

func testAccVolumeGroupIscsiClientsV2Config(filepath, name, desc string) string {
	return testAccVolumeGroupResourceConfig(filepath, name, desc) + `	
	  resource "nutanix_volume_group_iscsi_client_v2" "vg_iscsi_test" {
		vg_ext_id = resource.nutanix_volume_group_v2.test.id
		ext_id     = local.volumes.iscsi_client.ext_id
		iscsi_initiator_name = local.volumes.iscsi_client.initiator_name
		depends_on = [ resource.nutanix_volume_group_v2.test ]
	  }

	  data "nutanix_volume_group_iscsi_clients_v2" "vg_iscsi_test" {  
		ext_id= resource.nutanix_volume_group_v2.test.id
		depends_on = [ resource.nutanix_volume_group_iscsi_client_v2.vg_iscsi_test , resource.nutanix_volume_group_v2.test]  
	  }
`
}
