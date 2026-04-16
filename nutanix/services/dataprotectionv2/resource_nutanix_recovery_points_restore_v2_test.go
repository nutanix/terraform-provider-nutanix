package dataprotectionv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRecoveryPointRestore = "nutanix_recovery_point_restore_v2.test"

func TestAccV2NutanixRecoveryPointRestoreResource_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-recovery-point-%d", r)
	vmName := fmt.Sprintf("tf-test-vm-%d", r)
	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMConfigRecovery(vmName) + testVMConfig(vmName) +
					testRecoveryPointRestoreResourceConfig(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPointRestore, "volume_group_ext_ids.#"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPointRestore, "vm_ext_ids.#"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPointRestore, "ext_id"),
				),
			},
		},
	})
}

func testRecoveryPointRestoreResourceConfig(name, expirationTime string) string {
	return testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPointsWithMultipleVMAndVGs(name, expirationTime) + `
	resource "nutanix_recovery_point_restore_v2" "test" {
	  ext_id         = nutanix_recovery_points_v2.test.id
	  cluster_ext_id = local.cluster_ext_id
	  vm_recovery_point_restore_overrides {
		vm_recovery_point_ext_id = nutanix_recovery_points_v2.test.vm_recovery_points[0].ext_id
	  }
	  vm_recovery_point_restore_overrides {
		vm_recovery_point_ext_id = nutanix_recovery_points_v2.test.vm_recovery_points[1].ext_id
	  }
	  volume_group_recovery_point_restore_overrides {
		volume_group_recovery_point_ext_id = nutanix_recovery_points_v2.test.volume_group_recovery_points[0].ext_id
		volume_group_override_spec {
		  name = "vg-1-test-restore"
		}
	  }
	  volume_group_recovery_point_restore_overrides {
		volume_group_recovery_point_ext_id = nutanix_recovery_points_v2.test.volume_group_recovery_points[1].ext_id
		volume_group_override_spec {
		  name = "vg-2-test-restore"
		}
	  }
	  depends_on = [nutanix_recovery_points_v2.test]
	}`
}
