package volumesv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeIscsiClients = "data.nutanix_volume_iscsi_clients_v2.v_iscsi"

func TestAccV2NutanixVolumeIscsiClientsDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeIscsiClientsV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeIscsiClients, "iscsi_clients.#"),
					testAccCheckResourceAttrListNotEmpty(dataSourceVolumeIscsiClients, "iscsi_clients", "iscsi_initiator_name"),
				),
			},
		},
	})
}

func testAccVolumeIscsiClientsV2Config() string {
	return `
	data "nutanix_volume_iscsi_clients_v2" "v_iscsi" {}`
}
