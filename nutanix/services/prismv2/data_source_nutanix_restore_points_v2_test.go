package prismv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	datasourceNameListRestorePoints = "data.nutanix_restore_points_v2.test"
	retries                         = 120
	delay                           = 30 * time.Second
)

func TestAccV2NutanixRestorePointsDatasource_ListRestorePoints(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and create if backup target not exists
			{
				PreConfig: func() {
					fmt.Printf("Step 1: List backup targets and create if backup target not exists\n")
				},
				Config: testAccCheckBackupTargetExistAndCreateIfNotExistsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkBackupTargetExistAndCreateIfNotExists(),
				),
			},
			// Check last sync time of backup target
			{
				PreConfig: func() {
					fmt.Printf("Step 2: Check last sync time of backup target\n")
				},
				Config: testAccCheckBackupTargetExistAndCreateIfNotExistsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkLastSyncTimeBackupTarget(retries, delay),
				),
			},
			// Create the restore source, cluster location
			{
				PreConfig: func() {
					fmt.Printf("Step 3: Create the restore source, cluster location\n")
				},
				Config: testRestoreSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
				),
			},
			// List Points
			{
				PreConfig: func() {
					fmt.Printf("Step 4: List Restore Points\n")
				},
				Config: testRestoreSourceConfig() + testAccListRestorePointsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkAttributeLength(datasourceNameListRestorePoints, "restore_points", 1),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorePoints, "restore_points.0.creation_time"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorePoints, "restore_points.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorePoints, "restore_points.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorePoints, "restore_source_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorePoints, "restorable_domain_manager_ext_id"),
				),
			},
		},
	})
}

func testAccListRestorePointsConfig() string {
	return `

data "nutanix_restorable_pcs_v2" "test" {
  provider = nutanix-2
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.id
}

data "nutanix_restore_points_v2" "test" {
  provider = nutanix-2
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.id
}

output "restore_point" {
  value = data.nutanix_restore_points_v2.test.restore_points.0.ext_id
}

output "restorable_pc_ext_id" {
  value = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
}

data "nutanix_pc_v2" "test" {
  ext_id = local.domainManagerExtId
}

output "pc_details" {
  value = data.nutanix_pc_v2.test
}
`
}
