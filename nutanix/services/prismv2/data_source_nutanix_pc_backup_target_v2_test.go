package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameBackupTarget = "data.nutanix_pc_backup_target_v2.test"

func TestAccV2NutanixBackupTargetDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and Create if backup target not exists
			{
				Config: testAccCheckBackupTargetExistAndCreateIfNotExistsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkClusterLocationBackupTargetExistAndCreateIfNotExists(),
				),
			},
			// Create backup target, cluster location
			{
				Config: testAccFetchBackupTargetDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameBackupTarget, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameBackupTarget, "domain_manager_ext_id"),
					// check the name and ext_id of cluster location in backup target by comparing with the cluster name and ext_id
					resource.TestCheckResourceAttrPair(datasourceNameBackupTarget, "location.0.cluster_location.0.config.0.ext_id", "data.nutanix_cluster_v2.test", "id"),
					resource.TestCheckResourceAttrPair(datasourceNameBackupTarget, "location.0.cluster_location.0.config.0.name", "data.nutanix_cluster_v2.test", "name"),
				),
			},
		},
	})
}

func testAccCheckBackupTargetExistAndCreateIfNotExistsConfig() string {
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

`
}

func testAccFetchBackupTargetDatasourceConfig() string {
	return `
data "nutanix_clusters_v2" "pcs" {
	filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  domainManagerExtId = data.nutanix_clusters_v2.pcs.cluster_entities.0.ext_id
}

data "nutanix_pc_backup_targets_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
}

data "nutanix_pc_backup_target_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
  ext_id = data.nutanix_pc_backup_targets_v2.test.backup_targets.0.ext_id
}

# Get Cluster By Id to get the cluster name and ext_id
data "nutanix_cluster_v2" "test" {
  ext_id = data.nutanix_pc_backup_target_v2.test.location.0.cluster_location.0.config.0.ext_id
}
`
}
