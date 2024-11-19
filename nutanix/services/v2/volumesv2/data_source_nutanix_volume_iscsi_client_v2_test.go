package volumesv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeIscsiClient = "data.nutanix_volume_iscsi_client_v2.v_iscsi"

func TestAccNutanixVolumeIscsiClientV2_Basic(t *testing.T) {
	path, _ := os.Getwd()
	filepath := path + "/../../../../test_config_v2.json"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeIscsiClientV2Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceVolumeIscsiClient, "iscsi_initiator_name"),
					resource.TestCheckResourceAttr(dataSourceVolumeIscsiClient, "ext_id", testVars.Volumes.IscsiClient.ExtId),
				),
			},
		},
	})
}

func testAccVolumeIscsiClientV2Config(filepath string) string {
	return fmt.Sprintf(`
	locals {
		config = (jsondecode(file("%s")))
		volumes = local.config.volumes
	}

	data "nutanix_volume_iscsi_client_v2" "v_iscsi" {
		ext_id = local.volumes.iscsi_client.ext_id
	}`, filepath)
}
