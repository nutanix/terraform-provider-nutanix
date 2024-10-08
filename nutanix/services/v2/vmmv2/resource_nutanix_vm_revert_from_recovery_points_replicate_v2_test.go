package vmmv2_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRecoveryPointRestore = "nutanix_recovery_point_restore_v2.test"
const resourceNameRecoveryPoint = "nutanix_recovery_points_v2.test"
const resourceNameRevertVm = "nutanix_vm_revert_v2.test"

func TestAccNutanixRecoveryPointRestoreV2Resource_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create Recovery Point
			{
				Config: testRecoveryPointsResourceConfigWithVmRecoveryPoints(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoint, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoint, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoint, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoint, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoint, "recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoint, "vm_recovery_points.0.vm_ext_id"),
				),
			},
			// Create Recovery Point Restore
			{
				Config: testRecoveryPointResourceConfig(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPointRestore, "vm_ext_ids.#"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPointRestore, "ext_id"),
				),
			},
			// VM Revert
			{
				Config: testRevertVmResourceConfig(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRevertVm, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRevertVm, "vm_recovery_point_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRevertVm, "status", "SUCCEEDED"),
				),
			},
		},
	})
}

func testRecoveryPointsResourceConfigWithVmRecoveryPoints(name, expirationTime string) string {
	var path, _ = os.Getwd()
	var filepath = path + "/../../../../test_config_v2.json"

	return fmt.Sprintf(`
	data "nutanix_clusters" "clusters" {} 
	locals{
		cluster1 = [
			for cluster in data.nutanix_clusters.clusters.entities :
			cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	  	][0]
		config = (jsondecode(file("%[3]s")))
		data_protection = local.config.data_protection			
	}
	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = local.data_protection.vm_ext_id[0]  
		}
	}`, name, expirationTime, filepath)
}

func testRecoveryPointResourceConfig(name, expirationTime string) string {
	return testRecoveryPointsResourceConfigWithVmRecoveryPoints(name, expirationTime) + `
	resource "nutanix_recovery_point_restore_v2" "test" {
	  ext_id         = nutanix_recovery_points_v2.test.id
	  cluster_ext_id = local.cluster1
	  vm_recovery_point_restore_overrides {
		vm_recovery_point_ext_id = nutanix_recovery_points_v2.test.vm_recovery_points[0].ext_id
	  }
	  depends_on = [nutanix_recovery_points_v2.test]
	}`
}

func testRevertVmResourceConfig(name, expirationTime string) string {
	return testRecoveryPointResourceConfig(name, expirationTime) + `
		resource "nutanix_vm_revert_v2" "test" {
		  ext_id = nutanix_recovery_point_restore_v2.test.vm_ext_ids[0]
		  vm_recovery_point_ext_id = nutanix_recovery_points_v2.test.vm_recovery_points[0].ext_id
		  depends_on = [nutanix_recovery_point_restore_v2.test]
		}
     `
}
