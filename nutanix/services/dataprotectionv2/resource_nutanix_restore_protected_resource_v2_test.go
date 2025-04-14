package dataprotectionv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRestoreProtectedResource = "nutanix_restore_protected_resource_v2.test"

func TestAccV2NutanixRestoreProtectedResourceResource_RestoreVm(t *testing.T) {
	r := acctest.RandIntRange(1, 100)
	vmName := fmt.Sprintf("tf-test-protected-vm-restore-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-restore-vm-%d", r)
	description := "create a new protected vm and restore it"

	vmResourceName := "nutanix_virtual_machine_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResourceAndCleanup,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testRestoreProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					resource.TestCheckResourceAttr(vmResourceName, "name", vmName),
					waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
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
					// Clean up the promoted vm
					deleteRestoredVM(vmName),
				),
			},
		},
	})
}

func TestAccV2NutanixRestoreProtectedResourceResource_RestoreVG(t *testing.T) {
	r := acctest.RandIntRange(1, 100)
	vgName := fmt.Sprintf("tf-test-protected-vg-restore-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-vg-%d", r)
	description := "create a new protected vg and promote it"

	vgResourceName := "nutanix_volume_group_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResourceAndCleanup,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testRestoreProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vgResourceName, "id"),
					resource.TestCheckResourceAttr(vgResourceName, "name", vgName),
					//wait 7 minutes for the VG to be protected
					func(s *terraform.State) error {
						// wait 7 min for the VG to be protected
						time.Sleep(7 * time.Minute)
						return nil
					},
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

resource "nutanix_category_v2" "test" {
  key = "tf-test-category-pp-restore-vm-%[5]d"
  value = "tf_test_category_pp_restore_vm_%[5]d"
  description = "category for protection policy and protected vm"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[4]s"
  description = "%[3]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
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
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
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
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target"
    is_primary            = false
  }

  category_ids = [nutanix_category_v2.test.id]
}

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
	power_state = "OFF"
	depends_on = [nutanix_protection_policy_v2.test]
}

	`, filepath, vmName, description, ppName, r)
}

func testRestoreProtectedResourceVMConfig() string {
	remotePcIP := testVars.AvailabilityZone.RemotePcIP
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	remoteHostProviderConfig := fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}

`, username, password, remotePcIP, insecure, port)

	return fmt.Sprintf(
		`

%[1]s

resource "nutanix_restore_protected_resource_v2" "test" {
  provider = nutanix-2
  ext_id = nutanix_virtual_machine_v2.test.id
  cluster_ext_id = local.availability_zone.cluster_ext_id
}


`, remoteHostProviderConfig)
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

resource "nutanix_category_v2" "test" {
  key = "tf-test-category-pp-restore-vg-%[5]d"
  value = "tf_test_category_pp_restore_vg_%[5]d"
  description = "category for protection policy and protected vm"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[4]s"
  description = "%[3]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
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
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 300
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
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target"
    is_primary            = false
  }

  category_ids = [nutanix_category_v2.test.id]
}

resource "nutanix_volume_group_v2" "test" {
  name                               = "%[2]s"
  description                        = "%[3]s"
  cluster_reference                  = local.clusterExtId
}

resource "nutanix_associate_category_to_volume_group_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  categories {
    ext_id = nutanix_category_v2.test.id
  }
}


	`, filepath, vgName, description, ppName, r)
}

func testRestoreProtectedResourceVGConfig() string {
	remotePcIP := testVars.AvailabilityZone.RemotePcIP
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	remoteHostProviderConfig := fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}

`, username, password, remotePcIP, insecure, port)

	return fmt.Sprintf(
		`

%[1]s

resource "nutanix_restore_protected_resource_v2" "test" {
  provider = nutanix-2
  ext_id = nutanix_volume_group_v2.test.id
  cluster_ext_id = local.availability_zone.cluster_ext_id
}

`, remoteHostProviderConfig)
}
