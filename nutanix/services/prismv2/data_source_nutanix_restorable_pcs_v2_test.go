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
			// List backup targets and delete if backup target exists
			{
				PreConfig: func() {
					fmt.Printf("Step 1: List backup targets and delete if backup target exists\n")
				},
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkBackupTargetExist(),
				),
			},
			// Create backup target, cluster location
			{
				PreConfig: func() {
					fmt.Printf("Step 2: Create backup target, cluster location\n")
				},
				Config: testAccBackupTargetResourceClusterLocationConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.name"),
				),
			},
			// Create the restore source, cluster location
			{
				PreConfig: func() {
					fmt.Printf("Step 3: Create the restore source, cluster location\n")
				},
				Config: testAccBackupTargetResourceClusterLocationConfig() +
					testRestoreSourceConfig(),
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
				Config: testAccBackupTargetResourceClusterLocationConfig() +
					testRestoreSourceConfig() +
					testAccListRestorablePCConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameListRestorablePCs, "restorable_source_ext_id"),
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

func testAccListRestorablePCConfig() string {
	return `

data "nutanix_restorable_pcs_v2" "test" {
  provider = nutanix-2
  restorable_source_ext_id = nutanix_restore_source_v2.cluster-location.id
}

`
}
