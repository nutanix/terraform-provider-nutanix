package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameBackupTarget = "data.nutanix_backup_target_v2.test"

func TestAccV2NutanixBackupTargetDatasource_Basic(t *testing.T) {
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
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameBackupTargetClusterLocation, "location.0.cluster_location.0.config.0.name"),
				),
			},
			// Fetch backup target
			{
				Config: testAccBackupTargetResourceClusterLocationConfig() + testAccFetchBackupTargetDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameBackupTarget, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameBackupTarget, "domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameBackupTarget, "location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameBackupTarget, "location.0.cluster_location.0.config.0.name"),
				),
			},
		},
	})
}

func testAccFetchBackupTargetDatasourceConfig() string {
	return `
data "nutanix_backup_target_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
  ext_id = nutanix_backup_target_v2.cluster-location.id 
}
`
}
