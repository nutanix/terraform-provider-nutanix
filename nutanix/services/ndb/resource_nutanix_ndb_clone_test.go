package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceClone = "nutanix_ndb_clone.acctest-managed"

func TestAccEra_Clonebasic(t *testing.T) {
	r := acc.RandIntBetween(25, 35)
	name := fmt.Sprintf("test-pg-inst-tf-clone-%d", r)
	desc := "this is desc"
	vmName := fmt.Sprintf("testvm-%d", r)
	sshKey := testVars.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraCloneConfig(name, desc, vmName, sshKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceClone, "name", name),
					resource.TestCheckResourceAttr(resourceClone, "description", desc),
					resource.TestCheckResourceAttr(resourceClone, "clone", "true"),
					resource.TestCheckResourceAttrSet(resourceClone, "date_created"),
					resource.TestCheckResourceAttrSet(resourceClone, "database_name"),
					resource.TestCheckResourceAttrSet(resourceClone, "database_nodes.#"),
					resource.TestCheckResourceAttrSet(resourceClone, "linked_databases.#"),
					resource.TestCheckResourceAttrSet(resourceClone, "time_machine.#"),
				),
			},
		},
	})
}

func testAccEraCloneConfig(name, desc, vmName, sshKey string) string {
	return fmt.Sprintf(`
	data "nutanix_ndb_profiles" "p"{
	}
	data "nutanix_ndb_clusters" "clusters"{}

	locals {
		profiles_by_type = {
			for p in data.nutanix_ndb_profiles.p.profiles : p.type => p...
		}
		compute_profiles = {
			for p in local.profiles_by_type.Compute: p.name => p
		}
		network_profiles = {
			for p in local.profiles_by_type.Network: p.name => p
		}
		database_parameter_profiles = {
			for p in local.profiles_by_type.Database_Parameter: p.name => p
		}

		clusters = {
			for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
		}
	}

	data "nutanix_ndb_time_machines" "test1" {}

	data "nutanix_ndb_time_machine" "test"{
		time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
	}

	data "nutanix_ndb_tms_capability" "test"{
		time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
	}

	resource "nutanix_ndb_clone" "acctest-managed" {
		time_machine_id = data.nutanix_ndb_time_machine.test.id
		name = "%[1]s"
		description = "%[2]s"
		nx_cluster_id = local.clusters.NDBCluster.id
		ssh_public_key = "%[4]s"
		snapshot_id = data.nutanix_ndb_tms_capability.test.last_continuous_snapshot.0.id
		create_dbserver = true
		compute_profile_id =  local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
		network_profile_id = local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
		database_parameter_profile_id = local.database_parameter_profiles.DEFAULT_POSTGRES_PARAMS.id
		nodes{
		  	vm_name="%[3]s"
		  	compute_profile_id =  local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
			network_profile_id = local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			nx_cluster_id = local.clusters.NDBCluster.id
		}
		postgresql_info{
		  vm_name="%[3]s"
		  db_password= "pass"
		  # dbserver_description = "des"
		}
	}
	`, name, desc, vmName, sshKey)
}
