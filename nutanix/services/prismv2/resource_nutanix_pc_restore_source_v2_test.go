package prismv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRestoreSourceClusterLocation = "nutanix_pc_restore_source_v2.cluster-location"
const resourceNameRestoreSourceObjectStoreLocation = "nutanix_pc_restore_source_v2.object-store-location"

func TestAccV2NutanixRestoreSourceResource_ClusterLocation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and Create if backup target not exists
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkClusterLocationBackupTargetExistAndCreateIfNotExists(),
				),
			},
			{
				Config: testAccRestoreSourceResourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixRestoreSourceResource_ObjectStoreLocation(t *testing.T) {
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
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkObjectStoreLocationBackupTargetExistAndCreateIfNotExists(),
				),
			},
			// Create the restore source, Object store location
			{
				Config: testAccRestoreSourceResourceObjectStoreLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
		},
	})
}

func testAccRestoreSourceResourceClusterLocationConfig() string {
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


# list Clusters
data "nutanix_clusters_v2" "clusters" {}

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

`, username, password, endpoint, insecure, port)
}

func testAccRestoreSourceResourceObjectStoreLocationConfig() string {
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

locals {
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

`, filepath, username, password, endpoint, insecure, port)
}
