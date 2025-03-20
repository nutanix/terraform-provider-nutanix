package prismv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRestorePoint = "data.nutanix_pc_restore_point_v2.test"

func TestAccV2NutanixRestorePointDatasource_FetchRestorePoint(t *testing.T) {
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
				Config: testAccCheckBackupTargetLastSyncTimeConfig(),
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
			// Fetch Restore Point
			{
				PreConfig: func() {
					fmt.Printf("Step 5: Fetch Restore Point\n")
				},
				Config: testRestoreSourceConfig() + testAccListRestorePointsConfig() +
					testAccFetchRestorePointConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRestorePoint, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameRestorePoint, "restorable_domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameRestorePoint, "restore_source_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameRestorePoint, "domain_manager.#"),
					checkAttributeLength(datasourceNameRestorePoint, "domain_manager", 1),
					resource.TestCheckResourceAttrSet(datasourceNameRestorePoint, "domain_manager.0.config.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameRestorePoint, "domain_manager.0.config.0.build_info.0.version"),
					resource.TestCheckResourceAttrSet(datasourceNameRestorePoint, "domain_manager.0.config.0.size"),
				),
			},
		},
	})
}

func testAccCheckBackupTargetLastSyncTimeConfig() string {
	return `
# list Clusters
data "nutanix_clusters_v2" "cls" {
	filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

data "nutanix_clusters_v2" "clusters" {}


locals {
  domainManagerExtId = data.nutanix_clusters_v2.cls.cluster_entities.0.ext_id
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

data "nutanix_pc_backup_targets_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
}

output "domainManagerExtID" {
  value = local.domainManagerExtId
}

output "clusterExtID" {
  value = local.clusterExtId
}


data "nutanix_pc_backup_target_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
  ext_id = data.nutanix_pc_backup_targets_v2.test.backup_targets.0.ext_id
}

	`
}

func testAccFetchRestorePointConfig() string {
	return `


data "nutanix_pc_restore_point_v2" "test" {
  provider = nutanix-2
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
  restore_source_ext_id = nutanix_pc_restore_source_v2.cluster-location.id
  ext_id = data.nutanix_pc_restore_points_v2.test.restore_points.0.ext_id
}

`
}
