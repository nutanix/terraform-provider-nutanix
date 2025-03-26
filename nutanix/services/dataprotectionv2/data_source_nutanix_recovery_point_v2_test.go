package dataprotectionv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRecoveryPoint = "data.nutanix_recovery_point_v2.test"

func TestAccV2NutanixRecoveryPointDatasource_VmRecoveryPoints(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-recovery-point-%d", r)
	vmName := fmt.Sprintf("tf-test-rp-vm-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMConfigRecovery(vmName) + testRecoveryPointDatasourceConfigWithVMRecoveryPoints(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoint, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoint, "name", name),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoint, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoint, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoint, "recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoint, "vm_recovery_points.0.vm_ext_id"),
				),
			},
		},
	})
}

func testRecoveryPointDatasourceConfigWithVMRecoveryPoints(name, expirationTime string) string {
	return fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-1.id
		}

		depends_on = [ nutanix_virtual_machine_v2.test-1 ]
	}

	data "nutanix_recovery_point_v2" "test"{
		ext_id = nutanix_recovery_points_v2.test.id
		depends_on = [ nutanix_recovery_points_v2.test ]
	}

`, name, expirationTime, filepath)
}

func testVMConfigRecovery(name string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster_ext_id = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = (jsondecode(file("%[2]s")))
			availability_zone = local.config.availability_zone
		}

		resource "nutanix_virtual_machine_v2" "test-1"{
			name= "%[1]s"
			description =  "test recovery point vm 1"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster_ext_id
			}
		}
`, name, filepath)
}

func testVMConfig(name string) string {
	return fmt.Sprintf(`

		resource "nutanix_virtual_machine_v2" "test-2"{
			name= "%[1]s-vm2"
			description =  "test recovery point vm 2"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster_ext_id
			}
		}
`, name)
}
