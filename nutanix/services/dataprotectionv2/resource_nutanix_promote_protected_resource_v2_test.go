package dataprotectionv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"os"
	"strconv"
	"testing"
	"time"
)

const resourceNamePromoteProtectedResource = "nutanix_promote_protected_resource_v2.test"

const maxRetries = 60
const retryInterval = 10 * time.Second
const sleepTime = 120 * time.Second

func TestAccV2NutanixPromoteProtectedResourceResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-protected-vm-promote-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-%d", r)
	description := "create a new protected vm and promote it"

	vmResourceName := "nutanix_virtual_machine_v2.test"

	remotePcIP := testVars.DataProtection.RemotePcIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyPromoteProtectedResource,
		Steps: []resource.TestStep{
			//// create protection policy and protected vm
			{
				Config: testPromoteProtectedResourceVmAndProtectionPolicyConfig(vmName, ppName, description),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					waitForVmToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//promote protected vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Promote Protected Resource")
				},

				Config: testPromoteProtectedResourceConfig(vmName, remotePcIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNamePromoteProtectedResource, "promoted_vm_ext_id"),
				),
			},
			// Clean up the promoted vm
			{
				PreConfig: func() {
					fmt.Println("Step 3: Clean up")
					time.Sleep(5 * time.Second)
				},

				Config: testPromoteProtectedResourceConfig(vmName, remotePcIP),
				Check: resource.ComposeTestCheckFunc(
					deletePromotedVm(),
				),
			},
		},
	})
}

func testPromoteProtectedResourceVmAndProtectionPolicyConfig(vmName, ppName, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_domain_managers_v2" "pcs" {
}

# List categories
data "nutanix_categories_v2" "categories" {
}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {
}

locals {	
	config = jsondecode(file("%[1]s"))
  	data_policies = local.config.data_policies
}

resource "nutanix_protection_policy_v2" "test" {
 name        = "%[4]s"
 description = "%[3]s"

 replication_configurations {
   source_location_label = "source"
   remote_location_label = "target"
   schedule {
     recovery_point_objective_time_seconds         = 0
     recovery_point_type                           = "CRASH_CONSISTENT"
     sync_replication_auto_suspend_timeout_seconds = 10
   }
 }
 replication_configurations {
   source_location_label = "target"
   remote_location_label = "source"
   schedule {
     recovery_point_objective_time_seconds         = 0
     recovery_point_type                           = "CRASH_CONSISTENT"
     sync_replication_auto_suspend_timeout_seconds = 10
   }
 }

 replication_locations {
   domain_manager_ext_id = data.nutanix_domain_managers_v2.pcs.domain_managers[0].ext_id
   label                 = "source"
   is_primary            = true
 }
 replication_locations {
   domain_manager_ext_id = local.data_policies.domain_manager_ext_id
   label                 = "target"
   is_primary            = false
 }

 category_ids = [data.nutanix_categories_v2.categories.categories.5.ext_id,data.nutanix_categories_v2.categories.categories.6.ext_id]
}

resource "nutanix_virtual_machine_v2" "test"{
	name= "%[2]s"
	description =  "%[3]s"
	num_cores_per_socket = 1
	num_sockets = 1
	cluster {
		ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
	}
    categories {
	  ext_id = data.nutanix_categories_v2.categories.categories.5.ext_id
    }
	power_state = "OFF"
	depends_on = [nutanix_protection_policy_v2.test]
}


	`, filepath, vmName, description, ppName)
}

func testPromoteProtectedResourceConfig(name, remotePcIP string) string {
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

data "nutanix_virtual_machines_v2" "test" {
	provider = nutanix
	filter = "name eq '%[1]s'"
}

%[2]s

resource "nutanix_promote_protected_resource_v2" "test" {
  provider = nutanix-2
  ext_id = data.nutanix_virtual_machines_v2.test.vms[0].ext_id
  depends_on = [data.nutanix_virtual_machines_v2.test]
}

data "nutanix_virtual_machines_v2" "p-vm" {
	provider = nutanix-2 
	filter = "name eq '%[1]s'"
	depends_on = [nutanix_promote_protected_resource_v2.test]
}

`, name, remoteHostProviderConfig)
}
