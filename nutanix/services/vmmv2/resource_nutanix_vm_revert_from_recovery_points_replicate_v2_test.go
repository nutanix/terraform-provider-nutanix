package vmmv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameRecoveryPointRestore = "nutanix_recovery_point_restore_v2.test"
	resourceNameRecoveryPoint        = "nutanix_recovery_points_v2.test"
	resourceNameRevertVM             = "nutanix_vm_revert_v2.test"
)

func TestAccV2NutanixRecoveryPointRestoreResource_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create Recovery Point
			{
				Config: testRecoveryPointsResourceConfigWithVMRecoveryPoints(name, expirationTimeFormatted),
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
				Config: testRevertVMResourceConfig(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRevertVM, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRevertVM, "vm_recovery_point_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRevertVM, "status", "SUCCEEDED"),
				),
			},
		},
	})
}

func testRecoveryPointsResourceConfigWithVMRecoveryPoints(name, expirationTime string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}
	locals{
		cluster1 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		config = (jsondecode(file("%[3]s")))
		availability_zone = local.config.availability_zone
	}

	resource "nutanix_virtual_machine_v2" "test"{
		name= "tf-vm-%[1]s"
		description =  "terraform test vm for RP %[2]s"
		num_cores_per_socket = 1
		num_sockets = 1
		cluster {
			ext_id = local.cluster1
		}
		boot_config{
			legacy_boot{
			  boot_order = ["CDROM", "DISK","NETWORK"]
			}
		}
		cd_roms{
			disk_address{
				bus_type = "IDE"
				index= 0
			}
		}
		power_state = "OFF"
	}

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test.id
		}
		depends_on = [nutanix_virtual_machine_v2.test]
	}`, name, expirationTime, filepath)
}

func testRecoveryPointResourceConfig(name, expirationTime string) string {
	return testRecoveryPointsResourceConfigWithVMRecoveryPoints(name, expirationTime) + `
	resource "nutanix_recovery_point_restore_v2" "test" {
	  ext_id         = nutanix_recovery_points_v2.test.id
	  cluster_ext_id = local.cluster1
	  vm_recovery_point_restore_overrides {
		vm_recovery_point_ext_id = nutanix_recovery_points_v2.test.vm_recovery_points[0].ext_id
	  }
	  depends_on = [nutanix_recovery_points_v2.test, nutanix_virtual_machine_v2.test]
	}`
}

func testRevertVMResourceConfig(name, expirationTime string) string {
	return testRecoveryPointResourceConfig(name, expirationTime) + `
		resource "nutanix_vm_revert_v2" "test" {
		  ext_id                   = nutanix_virtual_machine_v2.test.id
		  vm_recovery_point_ext_id = nutanix_recovery_points_v2.test.vm_recovery_points[0].ext_id
		  depends_on = [nutanix_recovery_point_restore_v2.test]
		}
     `
}
