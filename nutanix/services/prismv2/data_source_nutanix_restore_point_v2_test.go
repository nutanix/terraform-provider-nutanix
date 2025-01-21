package prismv2_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameListRestorePoint = "data.nutanix_restore_point_v2.test"

func TestAccV2NutanixRestorePointDatasource_Basic(t *testing.T) {
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
			// List Points
			{
				PreConfig: func() {
					fmt.Printf("Step 4: List Restore Points\n")
				},
				Config: testAccBackupTargetResourceClusterLocationConfig() +
					testRestoreSourceConfig() +
					testAccListRestorePointConfig(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						aJson, _ := json.MarshalIndent(s.RootModule().Resources[datasourceNameListRestorePoint].Primary.Attributes, "", "  ")
						fmt.Println("############################################")
						fmt.Println(fmt.Sprintf("Resource Attributes: \n%v", string(aJson)))
						fmt.Println("############################################")

						return nil
					},
					resource.TestCheckResourceAttrSet(datasourceNameListRestorePoint, "ext_id"),
				),
			},
		},
	})
}

func testAccListRestorePointConfig() string {
	return `

data "nutanix_restorable_pcs_v2" "test" {
  provider = nutanix-2
  restorable_source_ext_id = nutanix_restore_source_v2.cluster-location.id
}

data "nutanix_restore_points_v2" "test" {
  provider = nutanix-2
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.id
}

data "nutanix_restore_point_v2" "test" {
  provider = nutanix-2
  restorable_domain_manager_ext_id = data.nutanix_restorable_pcs_v2.test.restorable_pcs.0.ext_id
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.id
  ext_id = data.nutanix_restore_points_v2.test.restore_points.0.ext_id
}

`
}
