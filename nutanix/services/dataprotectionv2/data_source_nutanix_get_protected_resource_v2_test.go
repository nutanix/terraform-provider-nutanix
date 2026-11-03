package dataprotectionv2_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameGetProtectedResource = "data.nutanix_protected_resource_v2.test"

func TestAccV2NutanixPromoteProtectedResourceDatasource_GetProtectedVm(t *testing.T) {
	// if the test is running using NUTANIX_API_KEY, skip the test
	if os.Getenv("NUTANIX_API_KEY") != "" {
		t.Skip("Skipping test as it not supported using NUTANIX_API_KEY")
	}
	r := acctest.RandIntRange(1, 99)
	vmName := fmt.Sprintf("tf-test-protected-vm-get-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-get-vm-%d", r)
	description := "create a new protected vm and get it"

	vmResourceName := "nutanix_virtual_machine_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResourceAndCleanup,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				PreConfig: func() {
					fmt.Printf("Step 1: Create protection policy and protected vm\n")
				},
				Config: testCreateProtectedResourceVMConfig(vmName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					resource.TestCheckResourceAttr(vmResourceName, "name", vmName),
					waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//Get protected vm
			{
				PreConfig: func() {
					fmt.Printf("Step 2: Get protected vm details\n")
					// Step 2: Add a PreConfig to pause before Terraform runs the refresh/apply for the data source
					time.Sleep(90 * time.Second)
				},
				Config: testGetProtectedResourceVMConfig() +
					testCreateProtectedResourceVMConfig(vmName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "ext_id"),
					resource.TestCheckResourceAttrPair(dataSourceNameGetProtectedResource, "entity_ext_id", vmResourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "replication_states.0.target_site_reference.0.cluster_ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "site_protection_info.0.location_reference.0.cluster_ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "source_site_reference.0.cluster_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameGetProtectedResource, "entity_type", "VM"),
					resource.TestCheckResourceAttr(dataSourceNameGetProtectedResource, "replication_states.0.replication_status", "IN_SYNC"),
				),
			},
		},
	})
}

func TestAccV2NutanixPromoteProtectedResourceDatasource_GetProtectedVG(t *testing.T) {
	// if the test is running using NUTANIX_API_KEY, skip the test
	if os.Getenv("NUTANIX_API_KEY") != "" {
		t.Skip("Skipping test as it not supported using NUTANIX_API_KEY")
	}
	r := acctest.RandIntRange(1, 99)
	vgName := fmt.Sprintf("tf-test-protected-vg-get-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-get-vg-%d", r)
	description := "create a new protected vg and get it"

	vgResourceName := "nutanix_volume_group_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResourceAndCleanup,
		Steps: []resource.TestStep{
			// create protection policy and protected VG
			{
				Config: testCreateProtectedResourceVgConfig(vgName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vgResourceName, "id"),
					resource.TestCheckResourceAttr(vgResourceName, "name", vgName),
				),
			},
			//Get protected VG
			{
				PreConfig: func() {
					fmt.Printf("Step 2: Get protected VG details\n")
					//delay 7 minutes to allow the VG to be protected
					time.Sleep(7 * time.Minute)
				},
				Config: testCreateProtectedResourceVgConfig(vgName, ppName, description, r) +
					testGetProtectedResourceVgConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "ext_id"),
					resource.TestCheckResourceAttrPair(dataSourceNameGetProtectedResource, "entity_ext_id", vgResourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceNameGetProtectedResource, "entity_type", "VOLUME_GROUP"),
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "site_protection_info.0.location_reference.0.cluster_ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "site_protection_info.0.location_reference.0.mgmt_cluster_ext_id"),
				),
			},
		},
	})
}

func testCreateProtectedResourceVMConfig(vmName, ppName, description string, r int) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

# list Clusters
data "nutanix_clusters_v2" "clusters" {}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][
  0
  ]
  config = jsondecode(file("%[1]s"))
  availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "synchronous-pp-category" {
  key   = "tf-synchronous-pp-%[5]d"
  value = "tf_synchronous_pp_%[5]d"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[4]s"
  description = "%[3]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
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

  category_ids = [nutanix_category_v2.synchronous-pp-category.id]
}

data "nutanix_storage_containers_v2" "ngt-sc" {
	filter = "clusterExtId eq '${local.clusterExtId}' and startswith(name,'default-container-')"
	limit = 1
}

resource "nutanix_virtual_machine_v2" "test" {
  name                 = "%[2]s"
  description          = "%[3]s"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.clusterExtId
  }
  categories {
    ext_id = nutanix_category_v2.synchronous-pp-category.id
  }
  power_state = "ON"

  disks {
    disk_address{
			bus_type = "SCSI"
			index = 0
		}
		backing_info{
			vm_disk{
				disk_size_bytes = "1073741824" # 10 GB
				storage_container{
					ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
				}
			}
		}
  }

  cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 0
    }
  }

  depends_on = [nutanix_protection_policy_v2.test]
}


	`, filepath, vmName, description, ppName, r)
}

func testGetProtectedResourceVMConfig() string {
	return `

data "nutanix_protected_resource_v2" "test" {
   ext_id = nutanix_virtual_machine_v2.test.id
}
`
}

func testCreateProtectedResourceVgConfig(vgName, ppName, description string, r int) string {
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
  key = "tf-test-category-pp-get-vg-%[5]d"
  value = "category_pp_protected_vg_%[5]d"
  description = "category for protection policy and protected vg"
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

data "nutanix_storage_containers_v2" "vg-sc" {
  filter = "clusterExtId eq '${local.clusterExtId}' and startswith(name,'default-container-')"
  limit  = 1
}

resource "nutanix_volume_group_v2" "test" {
  name                               = "%[2]s"
  description                        = "%[3]s"
  cluster_reference                  = local.clusterExtId
}

# A Volume Group must contain at least one vdisk to be protectable; an empty VG
# has no storage objects to replicate and is skipped, staying UNPROTECTED.
resource "nutanix_volume_group_disk_v2" "test" {
  volume_group_ext_id = nutanix_volume_group_v2.test.id
  index               = 0
  description         = "disk for protected vg"
  disk_size_bytes     = 1073741824
  disk_data_source_reference {
    name        = "%[2]s-disk"
    ext_id      = data.nutanix_storage_containers_v2.vg-sc.storage_containers[0].ext_id
    entity_type = "STORAGE_CONTAINER"
  }
  # The v4 GET does not return disk_data_source_reference (it is a create-time-only
  # source reference), so it would otherwise show a perpetual add-in-place diff.
  lifecycle {
    ignore_changes = [disk_data_source_reference]
  }
}

resource "nutanix_associate_category_to_volume_group_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  categories {
    ext_id = nutanix_category_v2.test.id
  }
  depends_on = [nutanix_volume_group_disk_v2.test]
}


	`, filepath, vgName, description, ppName, r)
}

func testGetProtectedResourceVgConfig() string {
	return `

data "nutanix_protected_resource_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
}

`
}
