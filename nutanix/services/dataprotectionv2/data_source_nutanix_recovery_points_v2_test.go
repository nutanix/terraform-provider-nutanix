package dataprotectionv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRecoveryPoints = "data.nutanix_recovery_points_v2.test"

func TestAccV2NutanixRecoveryPointsDatasource_Basic(t *testing.T) {
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
				Config: testVMConfigRecovery(vmName) + testRecoveryPointsDatasourceConfig(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.#"),
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.0.name"),
				),
			},
		},
	})
}

func TestAccV2NutanixRecoveryPointsDatasource_WithFilter(t *testing.T) {
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
				Config: testVMConfigRecovery(vmName) + testRecoveryPointsDatasourceConfigWithFilter(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.0.name", name),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.0.status", "COMPLETE"),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.0.expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.0.vm_recovery_points.0.vm_ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixRecoveryPointsDatasource_WithLimit(t *testing.T) {
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
				Config: testVMConfigRecovery(vmName) + testRecoveryPointsDatasourceConfigWithLimit(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.0.ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixRecoveryPointsDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryPointsDatasourceConfigWithInvalidFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.#", "0"),
				),
			},
		},
	})
}

func testRecoveryPointsDatasourceConfig(name, expirationTime string) string {
	return fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-1.id
		}
	}
	data "nutanix_recovery_points_v2" "test"{
		depends_on = [ nutanix_recovery_points_v2.test ]
	}

`, name, expirationTime)
}

func testRecoveryPointsDatasourceConfigWithFilter(name, expirationTime string) string {
	return fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-1.id
		}
	}

	data "nutanix_recovery_points_v2" "test"{
		filter = "name eq '${nutanix_recovery_points_v2.test.name}'"
		depends_on = [ nutanix_recovery_points_v2.test ]
	}

`, name, expirationTime)
}

func testRecoveryPointsDatasourceConfigWithLimit(name, expirationTime string) string {
	return fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-1.id
		}
	}

	data "nutanix_recovery_points_v2" "test"{
		limit = 1
	}

`, name, expirationTime)
}

func testRecoveryPointsDatasourceConfigWithInvalidFilter() string {
	return `
	data "nutanix_recovery_points_v2" "test"{
		filter = "name eq 'invalid_filter'"
	}

`
}
