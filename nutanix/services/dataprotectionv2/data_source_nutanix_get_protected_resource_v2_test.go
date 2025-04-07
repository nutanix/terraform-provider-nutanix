package dataprotectionv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameGetProtectedResource = "data.nutanix_protected_resource_v2.test"

func TestAccV2NutanixPromoteProtectedResourceDatasource_GetProtectedVm(t *testing.T) {
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

resource "nutanix_virtual_machine_v2" "test" {
  name                 = "%[2]s"
  description          = "%[3]s"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
  }
  categories {
    ext_id = nutanix_category_v2.synchronous-pp-category.id
  }
  power_state = "OFF"
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

func testGetProtectedResourceVgConfig() string {
	return `

data "nutanix_protected_resource_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
}

`
}
