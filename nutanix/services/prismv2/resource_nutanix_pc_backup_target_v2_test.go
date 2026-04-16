package prismv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameBackupTargetClusterLocation = "nutanix_pc_backup_target_v2.cluster-location"
const resourceNameBackupTargetObjectStoreLocation = "nutanix_pc_backup_target_v2.object-store-location"

func TestAccV2NutanixBackupTargetResource_ClusterLocation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and delete if backup target exists
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkClusterLocationBackupTargetExistAndDeleteIfExists(),
				),
			},
			// Create backup target, cluster location
			{
				Config: testAccBackupTargetResourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.name"),
					// check the name and ext_id of cluster location in backup target by comparing with the cluster name and ext_id
					resource.TestCheckResourceAttrPair(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.ext_id", "data.nutanix_cluster_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.name", "data.nutanix_cluster_v2.test", "name"),
				),
			},
		},
	})
}

func TestAccV2NutanixBackupTargetResource_ObjectStoreLocation(t *testing.T) {
	bucket := testVars.Prism.Bucket

	if bucket.Name == "" || bucket.AccessKey == "" || bucket.SecretKey == "" {
		t.Skip("Skipping test due to missing bucket configuration")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and delete if backup target exists
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkObjectStoreLocationBackupTargetExistAndDeleteIfExists(),
				),
			},
			// Create backup target, Object store location
			{
				Config: testAccBackupTargetResourceObjectStoreLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetObjectStoreLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.backup_policy.0.rpo_in_minutes", "60"),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
			// Update Backup target, Object store location
			{
				Config: testAccBackupTargetResourceObjectStoreLocationUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetObjectStoreLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.backup_policy.0.rpo_in_minutes", "120"),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
		},
	})
}

func TestAccV2NutanixBackupTargetResource_ClusterLocationAndObjectStoreLocation(t *testing.T) {
	bucket := testVars.Prism.Bucket

	if bucket.Name == "" || bucket.AccessKey == "" || bucket.SecretKey == "" {
		t.Skip("Skipping test due to missing bucket configuration")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and delete if backup target exists
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkClusterLocationBackupTargetExistAndDeleteIfExists(),
					checkObjectStoreLocationBackupTargetExistAndDeleteIfExists(),
				),
			},
			// Create backup target, Object store location
			{
				Config: testAccBackupTargetResourceObjectStoreLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetObjectStoreLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.backup_policy.0.rpo_in_minutes", "60"),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameBackupTargetObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
			// Create backup target, cluster location
			{
				Config: testAccBackupTargetResourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.name"),
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

resource "nutanix_pc_backup_target_v2" "cluster-location" {
  domain_manager_ext_id = local.domainManagerExtId
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtId
      }
    }
  }
}

# Get Cluster By Id to get the cluster name and ext_id
data "nutanix_cluster_v2" "test" {
  ext_id = nutanix_pc_backup_target_v2.cluster-location.location.0.cluster_location.0.config.0.ext_id
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

resource "nutanix_pc_backup_target_v2" "object-store-location" {
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

func testAccBackupTargetResourceObjectStoreLocationUpdateConfig() string {
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

resource "nutanix_pc_backup_target_v2" "object-store-location" {
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
        rpo_in_minutes = 120
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
