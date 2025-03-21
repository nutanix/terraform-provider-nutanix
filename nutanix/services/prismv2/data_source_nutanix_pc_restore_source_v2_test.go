package prismv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameRestoreSourceClusterLocation = "data.nutanix_pc_restore_source_v2.cluster-location"
const datasourceNameRestoreSourceObjectStoreLocation = "data.nutanix_pc_restore_source_v2.object-store-location"

func TestAccV2NutanixRestoreSourceDatasource_ClusterLocation(t *testing.T) {
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
			// Create the restore source, cluster location
			{
				Config: testAccRestoreSourceDatasourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRestoreSourceClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameRestoreSourceClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixRestoreSourceDatasource_ObjectStoreLocation(t *testing.T) {
	bucket := testVars.Prism.Bucket

	if bucket.Name == "" || bucket.AccessKey == "" || bucket.SecretKey == "" {
		t.Skip("Skipping test due to missing bucket configuration")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and Create if backup target not exists
			{
				Config: testAccCheckBackupTargetExistAndCreateIfNotExistsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkObjectStoreLocationBackupTargetExistAndCreateIfNotExists(),
				),
			},
			// Create the restore source, Object Store Location
			{
				Config: testAccRestoreSourceDatasourceObjectStoreLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameRestoreSourceObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(datasourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
		},
	})
}

func testAccRestoreSourceDatasourceClusterLocationConfig() string {
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	endpoint := testVars.Prism.RestoreSource.PeIP

	return fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}

data "nutanix_clusters_v2" "clusters" {
  provider = nutanix
}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_pc_restore_source_v2" "cluster-location" {
  provider = nutanix-2
  location {
    cluster_location {
      config {
        ext_id = local.clusterExtId
      }
    }
  }
}

data "nutanix_pc_restore_source_v2" "cluster-location" {
	provider = nutanix-2
	ext_id = nutanix_pc_restore_source_v2.cluster-location.id
}

`, username, password, endpoint, insecure, port)
}

func testAccRestoreSourceDatasourceObjectStoreLocationConfig() string {
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	endpoint := testVars.Prism.RestoreSource.PeIP

	return fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[2]s"
  password = "%[3]s"
  endpoint = "%[4]s"
  insecure = %[5]t
  port     = %[6]d
}

data "nutanix_clusters_v2" "clusters" {
  provider = nutanix
}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  config = jsondecode(file("%[1]s"))
  bucket = local.config.prism.bucket
}

resource "nutanix_pc_restore_source_v2" "object-store-location" {
  provider = nutanix-2
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
    }
  }
  lifecycle {
    ignore_changes = [
      location[0].object_store_location[0].provider_config[0].credentials
    ]
  }
}

data "nutanix_pc_restore_source_v2" "object-store-location" {
  provider = nutanix-2
  ext_id = nutanix_pc_restore_source_v2.object-store-location.id
}

`, filepath, username, password, endpoint, insecure, port)
}
