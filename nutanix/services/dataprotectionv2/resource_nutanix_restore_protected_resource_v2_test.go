package dataprotectionv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRestoreProtectedResource = "nutanix_restore_protected_resource_v2.test"

func TestAccV2NutanixRestoreProtectedResourceResource_RestoreVm(t *testing.T) {
	// if the test is running using NUTANIX_API_KEY, skip the test
	if os.Getenv("NUTANIX_API_KEY") != "" {
		t.Skip("Skipping test as it not supported using NUTANIX_API_KEY")
	}
	r := acctest.RandIntRange(1, 100)
	vmName := fmt.Sprintf("tf-test-protected-vm-restore-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-restore-vm-%d", r)
	description := "Protected VM for restore"

	vmResourceName := "nutanix_virtual_machine_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testCheckDestroyProtectedResourceAndCleanup,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testRestoreProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					resource.TestCheckResourceAttr(vmResourceName, "name", vmName),
					// Ansible: wait until protection_type == RULE_PROTECTED, then pause 300s.
					waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, 0),
					waitForProtectedResourceAvailable(vmResourceName, maxRetries, retryInterval, sleepTime),
				),
			},
			//restore protected vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Restore Protected Resource")
				},

				Config: testRestoreProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description, r) +
					testRestoreProtectedResourceVMConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "ext_id"),
					// Clean up the restored vm
					deleteRestoredVM(vmName),
				),
			},
		},
	})
}

func TestAccV2NutanixRestoreProtectedResourceResource_RestoreVG(t *testing.T) {
	// if the test is running using NUTANIX_API_KEY, skip the test
	if os.Getenv("NUTANIX_API_KEY") != "" {
		t.Skip("Skipping test as it not supported using NUTANIX_API_KEY")
	}
	r := acctest.RandIntRange(1, 100)
	vgName := fmt.Sprintf("tf-test-protected-vg-restore-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-restore-vg-%d", r)
	description := "Volume group for restore"

	vgResourceName := "nutanix_volume_group_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testCheckDestroyProtectedResourceAndCleanup,
		Steps: []resource.TestStep{
			// create protection policy and protected vg (Ansible: Create Volume group to be restored)
			{
				Config: testRestoreProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vgResourceName, "id"),
					resource.TestCheckResourceAttr(vgResourceName, "name", vgName),
					// Ansible: pause 300s after protect before restore.
					waitForProtectedResourceAvailable(vgResourceName, maxRetries, retryInterval, sleepTime),
				),
			},
			//restore protected vg
			{
				PreConfig: func() {
					fmt.Println("Step 2: Restore Protected Resource")
				},

				Config: testRestoreProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description, r) +
					testRestoreProtectedResourceVGConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "ext_id"),
					// Clean up the restored vg
					deleteRestoredVg(vgName),
				),
			},
		},
	})
}

func testRestoreProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description string, r int) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}


# list Clusters
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
	clusterExtId = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
	config = jsondecode(file("%[1]s"))
  	availability_zone = local.config.availability_zone
}

data "nutanix_storage_containers_v2" "ssc" {
  filter = "clusterExtId eq '${local.clusterExtId}' and name eq 'SelfServiceContainer'"
  limit  = 1
}

locals {
  sourceStorageContainerExtId = data.nutanix_storage_containers_v2.ssc.storage_containers[0].ext_id
}

# Mirrors Ansible category_ext_id_2 + "Create auto retention protection policy"
# (ntnx_protection_policies_v2/tasks/protection_policies.yml) used for restore VM/VG.
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-pp-restore-vm-%[5]d"
  value = "tf_test_category_pp_restore_vm_%[5]d"
  description = "category for restore auto-retention protection policy"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[4]s"
  description = "%[3]s"

  # Ansible auto-retention PP: RPO=60 (minutely), CRASH_CONSISTENT, DAILY rollup, suspend=300.
  # Restore VM uses this category (not sync RPO=0 — that is promote-only in Ansible).
  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_type                   = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds = 60
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
      sync_replication_auto_suspend_timeout_seconds = 300
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_type                   = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds = 60
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
      sync_replication_auto_suspend_timeout_seconds = 300
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target"
    is_primary            = false
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.availability_zone.cluster_ext_id]
      }
    }
  }

  category_ids = [nutanix_category_v2.test.id]
}

# Ansible: Create a Protected VM to be restored (category_ext_id_2, disk 26843545600).
resource "nutanix_virtual_machine_v2" "test"{
	name= "%[2]s"
	description =  "%[3]s"
	num_cores_per_socket = 1
	num_sockets = 1
	cluster {
		ext_id = local.clusterExtId
	}
  categories {
	  ext_id = nutanix_category_v2.test.id
  }
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        disk_size_bytes = 26843545600
        storage_container {
          ext_id = local.sourceStorageContainerExtId
        }
      }
    }
  }
  power_state = "ON"

	depends_on = [nutanix_protection_policy_v2.test]
}

	`, filepath, vmName, description, ppName, r)
}

func testRestoreProtectedResourceVMConfig() string {
	return fmt.Sprintf(
		`

%[1]s

data "nutanix_protected_resource_v2" "restore_ready" {
  ext_id = nutanix_virtual_machine_v2.test.id

  depends_on = [
    nutanix_virtual_machine_v2.test
  ]
}

resource "nutanix_restore_protected_resource_v2" "test" {
  provider = nutanix-2
  ext_id = nutanix_virtual_machine_v2.test.id
  cluster_ext_id = local.availability_zone.cluster_ext_id

  depends_on = [
    data.nutanix_protected_resource_v2.restore_ready
  ]
}


`, remotePCProviderConfig(testVars.AvailabilityZone.RemotePcIP))
}

func testRestoreProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description string, r int) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}


# list Clusters
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
	clusterExtId = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
	config = jsondecode(file("%[1]s"))
  	availability_zone = local.config.availability_zone
}

# Mirrors Ansible category_ext_id_2 + auto-retention PP used for restore VG
# (same policy as restore VM; promote VG uses sync category4 instead).
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-pp-restore-vg-%[5]d"
  value = "tf_test_category_pp_restore_vg_%[5]d"
  description = "category for restore auto-retention protection policy"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[4]s"
  description = "%[3]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_type                   = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds = 60
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
      sync_replication_auto_suspend_timeout_seconds = 300
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_type                   = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds = 60
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
      sync_replication_auto_suspend_timeout_seconds = 300
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target"
    is_primary            = false
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.availability_zone.cluster_ext_id]
      }
    }
  }

  category_ids = [nutanix_category_v2.test.id]
}

# Ansible: Create Volume group to be restored (name/description/cluster_reference only).
resource "nutanix_volume_group_v2" "test" {
  name              = "%[2]s"
  description       = "%[3]s"
  cluster_reference = local.clusterExtId
  depends_on        = [nutanix_protection_policy_v2.test]
}

# Ansible: Associate category with second VG (category_ext_id_2, entity_type CATEGORY).
resource "nutanix_associate_category_to_volume_group_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  categories {
    ext_id      = nutanix_category_v2.test.id
    entity_type = "CATEGORY"
  }
  depends_on = [
    nutanix_protection_policy_v2.test,
    nutanix_volume_group_v2.test,
  ]
}


	`, filepath, vgName, description, ppName, r)
}

func testRestoreProtectedResourceVGConfig() string {
	return fmt.Sprintf(
		`

%[1]s

data "nutanix_protected_resource_v2" "restore_ready" {
  ext_id = nutanix_volume_group_v2.test.id

  depends_on = [
    nutanix_associate_category_to_volume_group_v2.test
  ]
}

resource "nutanix_restore_protected_resource_v2" "test" {
  provider       = nutanix-2
  ext_id         = nutanix_volume_group_v2.test.id
  cluster_ext_id = local.availability_zone.cluster_ext_id

  depends_on = [
    data.nutanix_protected_resource_v2.restore_ready
  ]
}

`, remotePCProviderConfig(testVars.AvailabilityZone.RemotePcIP))
}
