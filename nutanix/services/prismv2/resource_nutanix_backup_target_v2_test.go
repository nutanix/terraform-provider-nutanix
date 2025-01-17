package prismv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"testing"
)

const resourceNameBackupTarget = "nutanix_backup_target_v2.test"

func TestAccV2NutanixBackupTargetResource_ClusterLocation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and delete if backup target exists
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkBackupTargetExist(),
				),
			},
			// Create backup target, cluster location
			{
				Config: testAccBackupTargetResourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "location.0.cluster_location.0.config.0.name"),
				),
			},
		},
	})
}

func TestAccV2NutanixBackupTargetResource_ObjectStoreLocation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and delete if backup target exists
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkBackupTargetExist(),
				),
			},
			// Create backup target, Object store location
			{
				Config: testAccBackupTargetResourceObjectStoreLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameBackupTarget, "location.0.object_store_location.0.backup_policy.0.rpo_in_minutes", "60"),
					resource.TestCheckResourceAttr(resourceNameBackupTarget, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameBackupTarget, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
		},
	})
}

func testAccListBackupTargetsDatasourceConfig() string {
	return `

# list Clusters
data "nutanix_clusters_v2" "cls" {
	filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  domainManagerExtId = data.nutanix_clusters_v2.cls.cluster_entities.0.ext_id
}

data "nutanix_backup_targets_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
}
`
}

func testAccBackupTargetResourceClusterLocationConfig() string {
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

resource "nutanix_backup_target_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtId
      }
    }
  }
}

`
}

func testAccBackupTargetResourceObjectStoreLocationConfig() string {
	return fmt.Sprintf(`
# list Clusters
data "nutanix_clusters_v2" "cls" {
	filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

locals {
  domainManagerExtId = data.nutanix_clusters_v2.cls.cluster_entities.0.ext_id
  config = jsondecode(file("%[1]s"))
  bucket = local.config.prism.bucket 
}

resource "nutanix_backup_target_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
  location {
    object_store_location {
      provider_config {
        bucket_name = local.bucket.name
        region      = local.bucket.region
        credentials {
          access_key_id     = local.bucket.access_key
          secret_access_key = local.bucket.secret_key
        }
      }
      backup_policy {
        rpo_in_minutes = 60
      }
    }
  }
  lifecycle {
    ignore_changes = [
      location[0].object_store_location[0].provider_config[0].credentials
    ]
  }
}

`, filepath)
}
