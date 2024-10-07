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

func TestAccNutanixRecoveryPointV2Datasource_VmRecoveryPoints(t *testing.T) {
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
				Config: testRecoveryPointDatasourceConfigWithVmRecoveryPoints(name, expirationTimeFormatted),
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

func testRecoveryPointDatasourceConfigWithVmRecoveryPoints(name, expirationTime string) string {
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
	data "nutanix_recovery_point_v2" "test"{
		ext_id = nutanix_recovery_points_v2.test.id
	}

`, name, expirationTime, filepath)
}
