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
const sleepTime = 120 * time.Second

func TestAccV2NutanixPromoteProtectedResourceResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-protected-vm-promote-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-%d", r)
	description := "create a new protected vm and promote it"

	vmResourceName := "nutanix_virtual_machine_v2.test"
	//vmPromoteResourceName := "data.nutanix_virtual_machine_v2.promote-vm"

	remotePcIP := testVars.DataProtection.RemotePcIP

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},

		CheckDestroy: testCheckDestroyPromoteProtectedResource,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testPromoteProtectedResourceVmAndProtectionPolicyConfig(vmName, ppName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					waitForVmToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			// promote protected vm
			{
				//Taint: []string{"data.nutanix_domain_managers_v2.pcs", "data.nutanix_categories_v2.categories", "data.nutanix_clusters_v2.clusters"},
				Config: testPromoteProtectedResourceVmAndProtectionPolicyConfig(vmName, ppName, description) +
					testPromoteProtectedResourceConfig(remotePcIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNamePromoteProtectedResource, "ext_id"),
					//resource.TestCheckResourceAttrSet(vmPromoteResourceName, "ext_id"),
				),
			},
		},
	})
}

//

func testPromoteProtectedResourceVmAndProtectionPolicyConfig(vmName, ppName, description string) string {
	//username := os.Getenv("NUTANIX_USERNAME")
	//password := os.Getenv("NUTANIX_PASSWORD")
	//port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))

	return fmt.Sprintf(`
provider "nutanix" {
  alias    = "host1"
  username = "admin"
  password = "Nutanix.123"
  endpoint = "10.44.76.58"
  insecure = true
  port     = 9440
}

# List domain Managers
data "nutanix_domain_managers_v2" "pcs" {
	provider = nutanix.host1
}

# List categories
data "nutanix_categories_v2" "categories" {
	provider = nutanix.host1
}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {
	provider = nutanix.host1
}

locals {	
	config = jsondecode(file("%[1]s"))
  	data_policies = local.config.data_policies
}

resource "nutanix_protection_policy_v2" "test" {
 provider = nutanix.host1
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
	provider = nutanix.host1
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

resource "time_sleep" "delay" {
  depends_on = [nutanix_virtual_machine_v2.test]

  create_duration = "5m"
}

provider "nutanix" {
  alias    = "host2"
  username = "admin"
  password = "Nutanix.123"
  endpoint = "10.44.76.117"
  insecure = true
  port     = 9440
}

resource "nutanix_promote_protected_resource_v2" "test" {
  provider = nutanix.host2
  ext_id = nutanix_virtual_machine_v2.test.id
  depends_on = [nutanix_virtual_machine_v2.test, time_sleep.delay]
}
	`, filepath, vmName, description, ppName)
}

func testPromoteProtectedResourceConfig(remotePcIp string) string {

	//username := os.Getenv("NUTANIX_USERNAME")
	//password := os.Getenv("NUTANIX_PASSWORD")
	//port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	//
	//fmt.Printf("username: %s\n", username)
	//fmt.Printf("password: %s\n", password)
	//fmt.Printf("endpoint: %s\n", remotePcIp)

	//return fmt.Sprintf(
	return `

provider "nutanix" {
  username = "admin"
  password = "Nutanix.123"
  endpoint = "10.44.76.117"
  insecure = true
  port     = 9440
}

resource "nutanix_promote_protected_resource_v2" "test" {
  provider = nutanix-host2
  ext_id = nutanix_virtual_machine_v2.test.id
  depends_on = [nutanix_virtual_machine_v2.test]
}

`
	//, username, password, remotePcIp, port)
}
