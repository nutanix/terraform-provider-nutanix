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

func TestAccNutanixRecoveryPointsV2Datasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryPointsDatasourceConfig(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.#"),
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.0.name"),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPointsV2Datasource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryPointsDatasourceConfigWithFilter(name, expirationTimeFormatted),
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

func TestAccNutanixRecoveryPointsV2Datasource_WithLimit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryPointsDatasourceConfigWithLimit(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameRecoveryPoints, "recovery_points.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceNameRecoveryPoints, "recovery_points.0.ext_id"),
				),
			},
		},
	})
}

func testRecoveryPointsDatasourceConfig(name, expirationTime string) string {
	return fmt.Sprintf(`
	locals{
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
	}
	data "nutanix_recovery_points_v2" "test"{}

`, name, expirationTime, filepath)
}

func testRecoveryPointsDatasourceConfigWithFilter(name, expirationTime string) string {
	return fmt.Sprintf(`
	locals{
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
	}
	data "nutanix_recovery_points_v2" "test"{
		filter = "name eq '%[1]s'"
		depends_on = [ nutanix_recovery_points_v2.test ]
	}

`, name, expirationTime, filepath)
}

func testRecoveryPointsDatasourceConfigWithLimit(name, expirationTime string) string {
	return fmt.Sprintf(`
	locals{
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
	}

	data "nutanix_recovery_points_v2" "test"{
		limit = 1
	}

`, name, expirationTime, filepath)
}
