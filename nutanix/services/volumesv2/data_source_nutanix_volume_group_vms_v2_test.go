package volumesv2_test

// Test plan for the nutanix_volume_group_vms_v2 (ListVmAttachmentsByVolumeGroupId) datasource.
//
// The datasource wraps the deprecated ListVmAttachmentsByVolumeGroupId API which
// returns the list of VM attachments (extId + optional SCSI bus index) for a
// Volume Group identified by {volumeGroupExtId}.
//
// Because this is a List datasource (no Create/Update/Delete), the datasource
// steps live in this file and compose the shared prerequisite configs:
//   - a Volume Group (nutanix_volume_group_v2)
//   - a VM (nutanix_virtual_machine_v2)
//   - a VM attachment (nutanix_volume_group_vm_v2)
//
// Scenarios covered:
//   1. WithAttachment  — create VG + VM + attachment, then read the list
//      datasource by volume_group_ext_id and assert:
//        * attachments.# == "1"
//        * attachments.0.ext_id is set (the attached VM ext_id)
//        * volume_group_ext_id matches the created VG
//   2. WithInvalidVolumeGroupID — query the datasource for a random,
//      non-existent volume group ext_id and expect the API to error out
//      (not-found), validating the error path of the read handler.

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceVolumeGroupVms = "data.nutanix_volume_group_vms_v2.test"

func TestAccV2NutanixVolumeGroupVmsDatasource_WithAttachment(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-volume-group-vms-%d", r)
	desc := "test volume group VM attachments list datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupVmsDatasourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceVolumeGroupVms, "volume_group_ext_id"),
					resource.TestCheckResourceAttr(datasourceVolumeGroupVms, "attachments.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceVolumeGroupVms, "attachments.0.ext_id"),
					resource.TestCheckResourceAttrPair(datasourceVolumeGroupVms, "attachments.0.ext_id",
						"nutanix_virtual_machine_v2.test", "id"),
					resource.TestCheckResourceAttrPair(datasourceVolumeGroupVms, "volume_group_ext_id",
						"nutanix_volume_group_v2.test", "id"),
				),
			},
		},
	})
}

func TestAccV2NutanixVolumeGroupVmsDatasource_WithInvalidVolumeGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVolumeGroupVmsDatasourceInvalidVGConfig(),
				ExpectError: regexp.MustCompile(`error while fetching VM attachments for the volume group`),
			},
		},
	})
}

func testAccVolumeGroupVmsDatasourceConfig(name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) + fmt.Sprintf(`
		resource "nutanix_virtual_machine_v2" "test" {
			name                 = "tf-test-vg-vms-%[1]s"
			description          = "%[2]s"
			num_cores_per_socket = 1
			num_sockets          = 1
			cluster {
				ext_id = local.cluster1
			}
			lifecycle {
				ignore_changes = [
					disks
				]
			}
		}

		resource "nutanix_volume_group_vm_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			vm_ext_id           = resource.nutanix_virtual_machine_v2.test.id
			index               = 1
			depends_on          = [resource.nutanix_volume_group_v2.test]
		}

		data "nutanix_volume_group_vms_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			depends_on          = [resource.nutanix_volume_group_vm_v2.test]
		}
	`, name, desc)
}

func testAccVolumeGroupVmsDatasourceInvalidVGConfig() string {
	return `
		data "nutanix_volume_group_vms_v2" "test" {
			volume_group_ext_id = "00000000-0000-0000-0000-000000000000"
		}
	`
}
