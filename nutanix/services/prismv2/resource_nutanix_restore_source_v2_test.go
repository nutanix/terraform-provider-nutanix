package prismv2_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRestoreSourceClusterLocation = "nutanix_restore_source_v2.cluster-location"
const resourceNameRestoreSourceObjectStoreLocation = "nutanix_restore_source_v2.object-store-location"

func TestAccV2NutanixRestoreSourceResource_ClusterLocation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create the restore source, cluster location
			{
				Config: testAccRestoreSourceResourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						aJson, _ := json.MarshalIndent(s.RootModule().Resources[resourceNameRestoreSourceClusterLocation].Primary.Attributes, "", "  ")
						fmt.Println("############################################")
						fmt.Println(fmt.Sprintf("Resource Attributes: \n%v", string(aJson)))
						fmt.Println("############################################")

						return nil
					}),
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
			// Create the restore source, Object store location
			{
				Config: testAccRestoreSourceResourceObjectStoreLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceObjectStoreLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.backup_policy.0.rpo_in_minutes", "60"),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
		},
	})
}

func TestAccV2NutanixRestoreSourceResource_ClusterLocationAndObjectStoreLocation(t *testing.T) {
	bucket := testVars.Prism.Bucket

	if bucket.Name == "" || bucket.AccessKey == "" || bucket.SecretKey == "" {
		t.Skip("Skipping test due to missing bucket configuration")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create the restore source, Object store location
			{
				Config: testAccRestoreSourceResourceObjectStoreLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceObjectStoreLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceObjectStoreLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.backup_policy.0.rpo_in_minutes", "60"),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.bucket_name", testVars.Prism.Bucket.Name),
					resource.TestCheckResourceAttr(resourceNameRestoreSourceObjectStoreLocation, "location.0.object_store_location.0.provider_config.0.region", testVars.Prism.Bucket.Region),
				),
			},
			// Create the restore source, cluster location
			{
				Config: testAccRestoreSourceResourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "location.0.cluster_location.0.config.0.name"),
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

data "nutanix_clusters_v2" "clusters" {
  provider = nutanix
}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_restore_source_v2" "cluster-location" {
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

resource "nutanix_backup_target_v2" "object-store-location" {
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
