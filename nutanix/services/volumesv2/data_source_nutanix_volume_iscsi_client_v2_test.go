package volumesv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeIscsiClient = "data.nutanix_volume_iscsi_client_v2.v_iscsi"

func TestAccV2NutanixVolumeIscsiClientDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeIscsiClientV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeIscsiClient, "iscsi_initiator_name"),
					resource.TestCheckResourceAttrSet(dataSourceVolumeIscsiClient, "ext_id"),
				),
			},
		},
	})
}

func testAccVolumeIscsiClientV2Config() string {
	return `

	data "nutanix_volume_iscsi_clients_v2" "test" {}

	data "nutanix_volume_iscsi_client_v2" "v_iscsi" {
		ext_id = data.nutanix_volume_iscsi_clients_v2.test.iscsi_clients.0.ext_id
	}`
}
