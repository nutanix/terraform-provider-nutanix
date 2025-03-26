package dataprotectionv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVMRecoveryPoint = "data.nutanix_vm_recovery_point_info_v2.test"

func TestAccV2NutanixVmRecoveryPointDatasource_VmRecoveryPoint(t *testing.T) {
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
				Config: testVMConfigRecovery(vmName) + testVMRecoveryPointDatasourceConfigWithVMRecoveryPoint(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "location_agnostic_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "recovery_point_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "vm_ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmRecoveryPointDatasource_VmRecoveryPointWithAppConsProps(t *testing.T) {
	t.Skip("Skipping this test case as it is failing due to missing app consistent properties in get request")
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
				Config: testVMConfigRecovery(vmName) + testVMRecoveryPointsDatasourceConfigWithAppConsProps(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "location_agnostic_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "recovery_point_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameVMRecoveryPoint, "vm_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "application_consistent_properties.0.backup_type", "FULL_BACKUP"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "application_consistent_properties.0.should_include_writers", "true"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "application_consistent_properties.0.should_store_vss_metadata", "true"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "application_consistent_properties.0.object_type", "dataprotection.v4.common.VssProperties"),
				),
			},
		},
	})
}

func testVMRecoveryPointDatasourceConfigWithVMRecoveryPoint(name, expirationTime string) string {
	return testRecoveryPointsResourceConfigWithVMRecoveryPoints(name, expirationTime) + `
	data "nutanix_vm_recovery_point_info_v2" "test" {
	  recovery_point_ext_id = nutanix_recovery_points_v2.test.ext_id
	  ext_id                = nutanix_recovery_points_v2.test.vm_recovery_points[0].ext_id
	  depends_on            = [nutanix_recovery_points_v2.test]
	}
`
}

func testVMRecoveryPointsDatasourceConfigWithAppConsProps(name, expirationTime string) string {
	return testRecoveryPointsResourceConfigWithVMRecoveryPointsWithAppConsProps(name, expirationTime) + `
		data "nutanix_vm_recovery_point_info_v2" "test" {
		  recovery_point_ext_id = nutanix_recovery_points_v2.test.ext_id
		  ext_id                = nutanix_recovery_points_v2.test.vm_recovery_points[0].ext_id
		  depends_on            = [nutanix_recovery_points_v2.test]
		}
	`
}
