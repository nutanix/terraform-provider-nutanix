package dataprotectionv2_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameGetProtectedResource = "nutanix_protected_resource_v2.test"

func TestAccV2NutanixPromoteProtectedResourceDatasource_GetProtectedVm(t *testing.T) {
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-protected-vm-get-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-get-vm-%d", r)
	description := "create a new protected vm and get it"

	vmResourceName := "nutanix_virtual_machine_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResource,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testCreateProtectedResourceVMConfig(vmName, ppName, description),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//Get protected vm
			{

				Config: testGetProtectedResourceVmConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "entity_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameGetProtectedResource, "entity_type", "VM"),
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[dataSourceNameGetProtectedResource].Primary.Attributes, "", "  ")
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

func TestAccV2NutanixPromoteProtectedResourceDatasource_GetProtectedVG(t *testing.T) {
	r := acctest.RandInt()
	vgName := fmt.Sprintf("tf-test-protected-vg-get-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-get-vg-%d", r)
	description := "create a new protected vg and get it"

	vgResourceName := "nutanix_volume_group_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccFoundationPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCheckDestroyProtectedResource,
		Steps: []resource.TestStep{
			// create protection policy and protected VG
			{
				Config: testCreateProtectedResourceVgConfig(vgName, ppName, description),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttrSet(vgResourceName, "id"),
					waitForVgToBeProtected(vgResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//Get protected VG
			{

				Config: testGetProtectedResourceVgConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameGetProtectedResource, "entity_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameGetProtectedResource, "entity_type", "VOLUME_GROUP"),
					func(s *terraform.State) error {
						aJSON, _ := json.MarshalIndent(s.RootModule().Resources[dataSourceNameGetProtectedResource].Primary.Attributes, "", "  ")
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

func testCreateProtectedResourceVMConfig(vmName, ppName, description string) string {
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

resource "nutanix_virtual_machine_v2" "test"{
	name= "%[2]s"
	description =  "%[3]s"
	num_cores_per_socket = 1
	num_sockets = 1
	cluster {
		ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
	}
    categories {
	  ext_id = local.category1
    }
	power_state = "OFF"
	depends_on = [nutanix_protection_policy_v2.test]
}


	`, filepath, vmName, description, ppName)
}

func testGetProtectedResourceVmConfig() string {
	return `

data "nutanix_protected_resource_v2" "test" {
	  ext_id = nutanix_virtual_machine_v2.test.id
}


`
}

func testCreateProtectedResourceVgConfig(vgName, ppName, description string) string {
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

func testGetProtectedResourceVgConfig() string {
	return `

data "nutanix_protected_resource_v2" "test" {
	  ext_id = nutanix_volume_group_v2.test.id
}


`
}
