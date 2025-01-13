package dataprotectionv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"testing"
	"time"
)

const resourceNamePromoteProtectedResource = "nutanix_promote_protected_resource_v2.test"

const maxRetries = 60
const retryInterval = 10 * time.Second

func TestAccV2NutanixPromoteProtectedResourceResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-%d", r)
	ppName := fmt.Sprintf("tf-test-promote-protected-resource-%d", r)
	description := "create a new protected vm and promote it"
	vmResourceName := "nutanix_virtual_machine_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPromoteProtectedResourceResourceConfig(vmName, ppName, description),
				Check: resource.ComposeTestCheckFunc(
					waitForVmToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval),
					resource.TestCheckResourceAttrSet(resourceNamePromoteProtectedResource, "ext_id"),
				),
			},
		},
	})
}

func testPromoteProtectedResourceResourceConfig(vmName, ppName, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_domain_managers_v2" "pcs" {}

# List categories
data "nutanix_categories_v2" "categories" {}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {}

locals {
	clusterExtId = [
		  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	][0]
	categoryExtId = data.nutanix_categories_v2.categories.categories.3.ext_id
	config = jsondecode(file("%[1]s"))
  	data_policies = local.config.data_policies
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
	  ext_id = local.categoryExtId
    }
}

data "nutanix_virtual_machine_v2" "test" {
  ext_id = nutanix_virtual_machine_v2.test.id
  depends_on = [nutanix_virtual_machine_v2.test]
}


	`, filepath, vmName, description, ppName)
}

//`
//resource "nutanix_protection_policy_v2" "test" {
//  name        = "%[4]s"
//  description = "%[3]s"
//
//  replication_configurations {
//    source_location_label = "source"
//    remote_location_label = "target"
//    schedule {
//      recovery_point_objective_time_seconds         = 0
//      recovery_point_type                           = "CRASH_CONSISTENT"
//      sync_replication_auto_suspend_timeout_seconds = 10
//    }
//  }
//  replication_configurations {
//    source_location_label = "target"
//    remote_location_label = "source"
//    schedule {
//      recovery_point_objective_time_seconds         = 0
//      recovery_point_type                           = "CRASH_CONSISTENT"
//      sync_replication_auto_suspend_timeout_seconds = 10
//    }
//  }
//
//  replication_locations {
//    domain_manager_ext_id = data.nutanix_domain_managers_v2.pcs.domain_managers[0].ext_id
//    label                 = "source"
//    is_primary            = true
//  }
//  replication_locations {
//    domain_manager_ext_id = local.data_policies.domain_manager_ext_id
//    label                 = "target"
//    is_primary            = false
//  }
//
//  category_ids = [data.nutanix_categories_v2.categories.categories.3.ext_id]
//}
//
//resource "nutanix_promote_protected_resource_v2" "test" {
//  ext_id = nutanix_protected_resource_v2.test.id
//}
//`
