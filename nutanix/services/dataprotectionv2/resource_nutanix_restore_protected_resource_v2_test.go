package dataprotectionv2_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRestoreProtectedResource = "nutanix_restore_protected_resource_v2.test"

func TestAccV2NutanixRestoreProtectedResourceResource_RestoreVm(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-protected-vm-restore-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-vm-%d", r)
	description := "create a new protected vm and promote it"

	vmResourceName := "nutanix_virtual_machine_v2.test"
	remotePcIP := testVars.DataProtection.RemotePcIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResource,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testRestoreProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description),
				Check: resource.ComposeTestCheckFunc(
					waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//restore protected vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Restore Protected Resource")
				},

				Config: testRestoreProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description) +
					testRestoreProtectedResourceVMConfig(remotePcIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "ext_id"),
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[resourceNameRestoreProtectedResource].Primary.Attributes, "", "  ")
						fmt.Printf("############################################\n")
						fmt.Printf(fmt.Sprintf("Resource Attributes: \n%v", string(aJSON)))
						fmt.Printf("############################################\n")
						return nil
					},
				),
			},
		},
	})
}

func TestAccV2NutanixRestoreProtectedResourceResource_RestoreVG(t *testing.T) {
	r := acctest.RandInt()
	vgName := fmt.Sprintf("tf-test-protected-vg-restore-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-vg-%d", r)
	description := "create a new protected vg and promote it"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResource,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testRestoreProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description),
				Check:  resource.ComposeTestCheckFunc(
				//waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//restore protected vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Restore Protected Resource")
				},

				Config: testRestoreProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description) +
					testRestoreProtectedResourceVGConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRestoreProtectedResource, "ext_id"),
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[resourceNameRestoreProtectedResource].Primary.Attributes, "", "  ")
						fmt.Printf("############################################\n")
						fmt.Printf(fmt.Sprintf("Resource Attributes: \n%v", string(aJSON)))
						fmt.Printf("############################################\n")
						return nil
					},
				),
			},
		},
	})
}

func testRestoreProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs" {
}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
	clusterExtId = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
	config = jsondecode(file("%[1]s"))
  	data_policies = local.config.data_policies
}

resource "nutanix_category_v2" "test" {
  key = "tf-test-category-pp"
  value = "tf_test_category_pp"
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
      sync_replication_auto_suspend_timeout_seconds = 20
      start_time                                    = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
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
      sync_replication_auto_suspend_timeout_seconds = 30
      start_time                                    = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.data_policies.domain_manager_ext_id
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

	`, filepath, vmName, description, ppName)
}

func testRestoreProtectedResourceVMConfig(remotePcIP string) string {
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
  cluster_ext_id = "00062c47-ac15-ee40-185b-ac1f6b6f97e2"
}


`, remoteHostProviderConfig)
}

func testRestoreProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs" {
}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
	clusterExtId = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
	config = jsondecode(file("%[1]s"))
  	data_policies = local.config.data_policies
}

resource "nutanix_category_v2" "test" {
  key = "tf-test-category-pp"
  value = "tf_test_category_pp"
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
      sync_replication_auto_suspend_timeout_seconds = 20
      start_time                                    = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
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
      sync_replication_auto_suspend_timeout_seconds = 30
      start_time                                    = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.data_policies.domain_manager_ext_id
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

	`, filepath, vgName, description, ppName)
}

func testRestoreProtectedResourceVGConfig() string {
	return `

resource "nutanix_restore_protected_resource_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  cluster_ext_id = local.clusterExtId
}


`
}
