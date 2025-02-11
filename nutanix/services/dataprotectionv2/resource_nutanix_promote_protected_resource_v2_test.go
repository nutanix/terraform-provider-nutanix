package dataprotectionv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePromoteProtectedResource = "nutanix_promote_protected_resource_v2.test"

const maxRetries = 60
const retryInterval = 10 * time.Second
const sleepTime = 5 * time.Minute

func TestAccV2NutanixPromoteProtectedResourceResource_PromoteVm(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-protected-vm-promote-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-vm-%d", r)
	description := "create a new protected vm and promote it"

	vmResourceName := "nutanix_virtual_machine_v2.test"

	remotePcIP := testVars.DataProtection.RemotePcIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResource,
		Steps: []resource.TestStep{
			//// create protection policy and protected vm
			{
				Config: testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//promote protected vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Promote Protected Resource")
				},

				Config: testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description) +
					testPromoteProtectedResourceVMConfig(remotePcIP),
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

				Config: testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description) +
					testPromoteProtectedResourceVMConfig(remotePcIP),
				Check: resource.ComposeTestCheckFunc(
					deletePromotedVM(),
				),
			},
		},
	})
}

func TestAccV2NutanixPromoteProtectedResourceResource_PromoteVG(t *testing.T) {
	r := acctest.RandInt()
	vgName := fmt.Sprintf("tf-test-protected-vg-promote-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-vg-%d", r)
	description := "create a new protected VG and promote it"

	//vgResourceName := "nutanix_virtual_machine_v2.test"

	remotePcIP := testVars.DataProtection.RemotePcIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResource,
		Steps: []resource.TestStep{
			//// create protection policy and protected vm
			{
				Config: testPromoteProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description),
				Check:  resource.ComposeTestCheckFunc(
				//resource.TestCheckResourceAttrSet(vgResourceName, "id"),
				//waitForVMToBeProtected(vgResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//promote protected vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Promote Protected Resource")
				},

				Config: testPromoteProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description) +
					testPromoteProtectedResourceVGConfig(vgName, remotePcIP),
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

				Config: testPromoteProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description) +
					testPromoteProtectedResourceVGConfig(vgName, remotePcIP),
				Check: resource.ComposeTestCheckFunc(
					deletePromotedVM(),
				),
			},
		},
	})
}

func testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs" {
}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {
}

locals {
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
		ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
	}
    categories {
	  ext_id = nutanix_category_v2.test.id
    }
	power_state = "OFF"
	depends_on = [nutanix_protection_policy_v2.test]
}


	`, filepath, vmName, description, ppName)
}

func testPromoteProtectedResourceVMConfig(remotePcIP string) string {
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

resource "nutanix_promote_protected_resource_v2" "test" {
  provider = nutanix-2
  ext_id = nutanix_virtual_machine_v2.test.id
}

`, remoteHostProviderConfig)
}

func testPromoteProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs" {
}

# List categories
data "nutanix_categories_v2" "categories" {
}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {
}

locals {	
	category1 = data.nutanix_categories_v2.categories.categories.5.ext_id
	category2 = data.nutanix_categories_v2.categories.categories.6.ext_id
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

 category_ids = [local.category1,local.category2]
 depends_on = [data.nutanix_categories_v2.categories]
}

resource "nutanix_volume_group_v2" "test" {
  name                               = "%[2]s"
  description                        = "%[3]s"
  should_load_balance_vm_attachments = false
  sharing_status                     = "SHARED"
  iscsi_features {
    target_secret			 = "1234567891011"
    enabled_authentications  = "CHAP"
  }
  storage_features {
    flash_mode {
      is_enabled = false
    }
  }
  usage_type = "USER"
  is_hidden = false
  lifecycle {
    ignore_changes = [
      iscsi_features[0].target_secret
    ]
  }
}

resource "nutanix_associate_category_to_volume_group_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  categories{
    ext_id = local.category1
  }
}


	`, filepath, vgName, description, ppName)
}

func testPromoteProtectedResourceVGConfig(name, remotePcIP string) string {
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

data "nutanix_volume_groups_v2" "test" {
	provider = nutanix
	filter = "name eq '%[1]s'"
}

%[2]s

resource "nutanix_promote_protected_resource_v2" "test" {
  provider = nutanix-2
  ext_id = data.nutanix_volume_groups_v2.test.volumes[0].ext_id
  depends_on = [data.nutanix_volume_groups_v2.test]
}

data "nutanix_volume_groups_v2" "p-vg" {
	provider = nutanix-2 
	filter = "name eq '%[1]s'"
	depends_on = [nutanix_promote_protected_resource_v2.test]
}

`, name, remoteHostProviderConfig)
}
