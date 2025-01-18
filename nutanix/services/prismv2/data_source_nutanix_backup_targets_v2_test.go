package prismv2_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"testing"
)

const datasourceNameBackupTargets = "data.nutanix_backup_targets_v2.test"

func TestAccV2NutanixBackupTargetsDatasource_Basic(t *testing.T) {
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
			// List backup targets
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameBackupTargets, "backup_targets.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameBackupTargets, "backup_targets.0.location.0.cluster_location.0.config.0.ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameBackupTargets, "backup_targets.0.location.0.cluster_location.0.config.0.name"),
				),
			},
		},
	})
}
