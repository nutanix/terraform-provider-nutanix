package prismv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	datasourceNameListRestorePoints = "data.nutanix_pc_restore_points_v2.test"
	retries                         = 120
	delay                           = 30 * time.Second
)

func TestAccV2NutanixRestorePointsDatasource_ListRestorePointsClusterLocation(t *testing.T) {
	var backupTargetExtID, domainManagerExtID = new(string), new(string)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and create if backup target not exists
			// Check last sync time of backup target to ensure that the restore points are available
			{
				PreConfig: func() {
					fmt.Printf("Step 1: List backup targets and create if backup target not exists\n")
				},
				Config: testAccCheckBackupTargetExistAndCreateIfNotExistsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkClusterLocationBackupTargetExistAndCreateIfNot(backupTargetExtID, domainManagerExtID),
					checkLastSyncTimeBackupTarget(domainManagerExtID, backupTargetExtID, retries, delay),
				),
			},
			// Create the restore source, cluster location
			{
				PreConfig: func() {
					fmt.Printf("Step 2: Create the restore source, cluster location\n")
				},
				Config: testClusterLocationRestoreSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
				),
			},
			// List Points
			{
				PreConfig: func() {
					fmt.Printf("Step 3: List Restore Points\n")
				},
				Config: testClusterLocationRestoreSourceConfig() + testAccListRestorePointsClusterLocationConfig(),
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

func TestAccV2NutanixRestorePointsDatasource_ListRestorePointsObjectStoreLocation(t *testing.T) {
	bucket := testVars.Prism.Bucket

	if bucket.Name == "" || bucket.AccessKey == "" || bucket.SecretKey == "" {
		t.Skip("Skipping test due to missing bucket configuration")
	}

	var backupTargetExtID, domainManagerExtID = new(string), new(string)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and create if backup target not exists
			// Check last sync time of backup target to ensure that the restore points are available
			{
				PreConfig: func() {
					fmt.Printf("Step 1: List backup targets and create if backup target not exists\n")
				},
				Config: testAccCheckBackupTargetExistAndCreateIfNotExistsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkObjectRestoreLocationBackupTargetExistAndCreateIfNot(backupTargetExtID, domainManagerExtID),
					checkLastSyncTimeBackupTarget(domainManagerExtID, backupTargetExtID, retries, delay),
				),
			},
			// Create the restore source, cluster location
			{
				PreConfig: func() {
					fmt.Printf("Step 2: Create the restore source, cluster location\n")
				},
				Config: testObjectStoreLocationRestoreSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
			// List Points
			{
				PreConfig: func() {
					fmt.Printf("Step 3: List Restore Points\n")
				},
				Config: testObjectStoreLocationRestoreSourceConfig() + testAccListRestorePointsObjectStoreLocationConfig(),
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

func testAccListRestorePointsClusterLocationConfig() string {
	return `

data "nutanix_restorable_pcs_v2" "test" {
  provider = nutanix-2
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.id
  filter = "extId eq ${local.domainManagerExtId}"
}

data "nutanix_pc_restore_points_v2" "test" {
  provider = nutanix-2
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.id
}

output "restore_point" {
  value = data.nutanix_pc_restore_points_v2.test.restore_points.0.ext_id
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

func testAccListRestorePointsObjectStoreLocationConfig() string {
	return `

data "nutanix_restorable_pcs_v2" "test" {
  provider = nutanix-2
  restore_source_ext_id = nutanix_pc_restore_source_v2.object-store-location.id
  filter = "extId eq ${local.domainManagerExtId}"
}

data "nutanix_pc_restore_points_v2" "test" {
  provider = nutanix-2
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
  restore_source_ext_id = nutanix_pc_restore_source_v2.object-store-location.id
}

output "restore_point" {
  value = data.nutanix_pc_restore_points_v2.test.restore_points.0.ext_id
}

output "restorable_pc_ext_id" {
  value = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
}

`
}
