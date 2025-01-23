package prismv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameDeployPC = "nutanix_deploy_pc_v2.test"
const datasourceNameFetchPC = "data.nutanix_pc_v2.test"

const resourceNameBackupTarget = "nutanix_backup_target_v2.test"
const datasourceNameListBackupTargets = "data.nutanix_backup_targets_v2.test"
const datasourceNameFetchBackupTarget = "data.nutanix_backup_target_v2.test"
const resourceNameRestoreSource = "nutanix_restore_source_v2.test"
const datasourceNameFetchRestoreSource = "data.nutanix_restore_source_v2.test"
const resourceNameRestorePC = "nutanix_restore_pc_v2.test"
const resourceNameUnregisterPC = "nutanix_unregister_cluster_v2.test"

func TestAccV2NutanixDeployPcResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-deploy-pc-%d", r)

	// config strings
	backupTargetConfig := testAccDeployPCConfig(name) + testAccBackupTargetResourceConfig()
	backupTargetUpdateConfig := testAccDeployPCConfig(name) + testAccBackupTargetResourceUpdateConfig()
	restoreSourceConfig := backupTargetConfig + testAccRestoreSourceResourceConfig()
	restorePCConfig := restoreSourceConfig + testAccRestorePCResourceConfig()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// deploy pc
			{
				Config: testAccDeployPCConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDeployPC, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "name", name),
				),
			},
			// List pcs
			{
				Config: testAccDeployPCConfig(name) + testAccListPCConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameListPCs, "pcs.#"),
				),
			},
			// Fetch pc
			{
				Config: testAccDeployPCConfig(name) + testAccFetchPCConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameFetchPC, "ext_id"),
				),
			},
			// Create backup target
			{
				Config: backupTargetConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "ext_id"),
				),
			},
			// List backup targets
			{
				Config: backupTargetConfig + testAccListBackupTargetsDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameListBackupTargets, "backup_targets.#"),
				),
			},
			// Fetch backup target
			{
				Config: backupTargetConfig + testAccFetchBackupTargetDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameFetchBackupTarget, "ext_id"),
				),
			},
			// Create restore source
			{
				Config: restoreSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreSource, "ext_id"),
				),
			},
			// Fetch restore source
			{
				Config: restoreSourceConfig + testAccFetchRestoreSourceDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameFetchRestoreSource, "ext_id"),
				),
			},
			// restore pc
			{
				Config: restorePCConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestorePC, "ext_id"),
				),
			},

			// update backup target
			{
				Config: backupTargetUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameBackupTarget, "ext_id"),
				),
			},
			// unregister pc
			{
				Config: testAccDeployPCConfig(name) + testAccUnregisterPCResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUnregisterPC, "ext_id"),
				),
			},
		},
	})
}

func testAccDeployPCConfig(name string) string {
	return fmt.Sprintf(`
 resource "nutanix_deploy_pc_v2" "test" {
  config {
    build_info {
      version = "5.17.0"
    }
    size = "SMALL"
    name = "%[1]s"
  }
  network {
    external_address {
      ipv4 {
        value = ""
      }
    }
    ntp_servers {
      ipv4 {
        value = ""
      }
    }
    name_servers {
      ipv4 {
        value = ""
      }
    }
  }
}
 
 `, name, filepath)
}

// Backup Target
func testAccBackupTargetResourceConfig() string {
	return `

resource "nutanix_backup_target_v2" "test" {
  domain_manager_ext_id = nutanix_deploy_pc_v2.test.id
  location {
    cluster_location {
      config {
        ext_id = "cluster uuid"
      }
    }
    object_store_location {
      provider_config {
        bucket_name = "bucket name"
        region      = "region"
        credentials {
          access_key_id     = ""
          secret_access_key = ""
        }
      }
      backup_policy {
        rpo_in_minutes = 0
      }
    }
  }
}

`
}

func testAccBackupTargetResourceUpdateConfig() string {
	return `

resource "nutanix_backup_target_v2" "test" {
  domain_manager_ext_id = nutanix_deploy_pc_v2.test.id
  location {
    cluster_location {
      config {
        ext_id = "cluster uuid"
      }
    }
    object_store_location {
      provider_config {
        bucket_name = "bucket name"
        region      = "region"
        credentials {
          access_key_id     = ""
          secret_access_key = ""
        }
      }
      backup_policy {
        rpo_in_minutes = 0
      }
    }
  }
}

`
}

// restore source
func testAccRestoreSourceResourceConfig() string {
	return `
resource "nutanix_restore_source_v2" "test" {
  location {
    cluster_location {
      config {
        ext_id = "cluster uuid"
      }
    }
    object_store_location {
      provider_config {
        bucket_name = "bucket name"
        region      = "region"
        credentials {
          access_key_id     = ""
          secret_access_key = ""
        }
      }
      backup_policy {
        rpo_in_minutes = 0
      }
    }
  }
}
`
}

func testAccFetchRestoreSourceDatasourceConfig() string {
	return `
data "nutanix_restore_source_v2" "test" {
  ext_id = nutanix_restore_source_v2.test.id
}
`
}

// restore pc

func testAccRestorePCResourceConfig() string {
	return `
resource "nutanix_restore_pc_v2" "test" {
  restorable_domain_manager_ext_id = nutanix_deploy_pc_v2.test.id
  restore_source_ext_id            = nutanix_restore_source_v2.test.id
  ext_id                           = nutanix_restore_pc_v2.test.id
  domain_manager {
    config {
      name = ""
      size = ""
    }
    network {
      external_address {
        ipv4 {
          value = ""
        }
      }
      ntp_servers {
        ipv4 {
          value = ""
        }
      }
      name_servers {
        ipv4 {
          value = ""
        }
      }
    }
    should_enable_high_availability = false
  }
}
`
}

// unregister

func testAccUnregisterPCResourceConfig() string {
	return `
resource "nutanix_unregister_cluster_v2" "test" {
    pc_ext_id = nutanix_deploy_pc_v2.test.id
    ext_id = "cluster uuid"
}
`
}
