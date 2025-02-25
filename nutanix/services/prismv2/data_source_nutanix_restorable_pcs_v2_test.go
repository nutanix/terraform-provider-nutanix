package prismv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameListRestorablePCs = "data.nutanix_restorable_pcs_v2.test"

func TestAccV2NutanixRestorablePcsDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List backup targets and Create if backup target not exists
			{
				Config: testAccCheckBackupTargetExistAndCreateIfNotExistsConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkBackupTargetExistAndCreateIfNotExists(),
				),
			},
			// Create the restore source, cluster location
			{
				Config: testRestoreSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreSourceClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
				),
			},
			// List Restorable pcs
			{
				PreConfig: func() {
					fmt.Printf("Step 4: List Restorable pcs\n")
				},
				Config: testRestoreSourceConfig() + testAccListRestorablePCConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restore_source_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.#"),
					checkAttributeLength(datasourceNameListRestorablePCs, "restorable_pcs", 1),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.config.0.name"),
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_pcs.0.network.0.external_address.0.ipv4.0.value"),
				),
			},
		},
	})
}

func testRestoreSourceConfig() string {
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

output "restore_source" {
   value = nutanix_restore_source_v2.cluster-location.id
}

`, username, password, endpoint, insecure, port)
}

func testAccListRestorablePCConfig() string {
	return `

data "nutanix_restorable_pcs_v2" "test" {
  provider = nutanix-2
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.ext_id
}

`
}
