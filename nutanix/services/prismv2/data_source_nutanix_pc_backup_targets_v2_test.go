package prismv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameBackupTargets = "data.nutanix_pc_backup_targets_v2.test"

func TestAccV2NutanixBackupTargetsDatasource_Basic(t *testing.T) {
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
			// List backup targets
			{
				Config: testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkAttributeLength(datasourceNameBackupTargets, "backup_targets", 1),
					resource.TestCheckResourceAttrSet(datasourceNameBackupTargets, "backup_targets.0.ext_id"),
					// check the name and ext_id of cluster location in backup target by comparing with the cluster name and ext_id
					resource.TestCheckResourceAttrPair(datasourceNameBackupTargets, "backup_targets.0.location.0.cluster_location.0.config.0.ext_id","data.nutanix_cluster_v2.test","id"),
					resource.TestCheckResourceAttrPair(datasourceNameBackupTargets, "backup_targets.0.location.0.cluster_location.0.config.0.name","data.nutanix_cluster_v2.test","name"),
				),
			},
		},
	})
}
