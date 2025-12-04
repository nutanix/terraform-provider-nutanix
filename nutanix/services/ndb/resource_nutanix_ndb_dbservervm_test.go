package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameDBServer = "nutanix_ndb_dbserver_vm.acctest-managed"

func TestAccEra_DBServerVMbasic(t *testing.T) {
	r := acc.RandIntBetween(21, 30)
	name := fmt.Sprintf("test-dbserver-%d", r)
	desc := "this is desc"
	sshKey := testVars.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseServerConfig(name, desc, sshKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDBServer, "name", name),
					resource.TestCheckResourceAttr(resourceNameDBServer, "description", desc),
					resource.TestCheckResourceAttr(resourceNameDBServer, "status", "UP"),
					resource.TestCheckResourceAttr(resourceNameDBServer, "type", "DBSERVER"),
					resource.TestCheckResourceAttrSet(resourceNameDBServer, "properties.#"),
				),
			},
		},
	})
}

func TestAccEra_DBServerVMbasicWithTimeMachine(t *testing.T) {
	r := acc.RandIntBetween(161, 170)
	name := fmt.Sprintf("test-dbserver-%d", r)
	desc := "this is desc"
	sshKey := testVars.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseServerTMSConfig(name, desc, sshKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDBServer, "name", name),
					resource.TestCheckResourceAttr(resourceNameDBServer, "description", desc),
					resource.TestCheckResourceAttr(resourceNameDBServer, "status", "UP"),
					resource.TestCheckResourceAttr(resourceNameDBServer, "type", "DBSERVER"),
					resource.TestCheckResourceAttrSet(resourceNameDBServer, "properties.#"),
				),
			},
		},
	})
}

func testAccEraDatabaseServerConfig(name, desc, sshKey string) string {
	return fmt.Sprintf(`
	data "nutanix_ndb_profiles" "p"{
	}
	data "nutanix_ndb_slas" "slas"{}
	data "nutanix_ndb_clusters" "clusters"{}

	locals {
		profiles_by_type = {
			for p in data.nutanix_ndb_profiles.p.profiles : p.type => p...
		}
		storage_profiles = {
			for p in local.profiles_by_type.Storage: p.name => p
		}
		compute_profiles = {
			for p in local.profiles_by_type.Compute: p.name => p
		}
		network_profiles = {
			for p in local.profiles_by_type.Network: p.name => p
		}
		software_profiles = {
			for p in local.profiles_by_type.Software: p.name => p
		}
		clusters = {
			for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
		}
	}

	resource nutanix_ndb_dbserver_vm acctest-managed {
		database_type = "postgres_database"
		software_profile_id = local.software_profiles["POSTGRES_15.6_ROCKY_LINUX_8_OOB"].id
		software_profile_version_id =  local.software_profiles["POSTGRES_15.6_ROCKY_LINUX_8_OOB"].latest_version_id
		compute_profile_id =  local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
		network_profile_id = local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
		nx_cluster_id = local.clusters.NDBCluster.id
		vm_password = "pass"
		postgres_database {
			vm_name = "%[1]s"
			client_public_key = "%[3]s"
		}
		description = "%[2]s"

	}
	`, name, desc, sshKey)
}

func testAccEraDatabaseServerTMSConfig(name, desc, sshKey string) string {
	return fmt.Sprintf(`
	data "nutanix_ndb_profiles" "p"{
	}
	data "nutanix_ndb_slas" "slas"{}
	data "nutanix_ndb_clusters" "clusters"{}

	locals {
		profiles_by_type = {
			for p in data.nutanix_ndb_profiles.p.profiles : p.type => p...
		}
		storage_profiles = {
			for p in local.profiles_by_type.Storage: p.name => p
		}
		compute_profiles = {
			for p in local.profiles_by_type.Compute: p.name => p
		}
		network_profiles = {
			for p in local.profiles_by_type.Network: p.name => p
		}
		software_profiles = {
			for p in local.profiles_by_type.Software: p.name => p
		}
		clusters = {
			for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
		}
	}

	data "nutanix_ndb_time_machines" "test1" {}

	resource nutanix_ndb_dbserver_vm acctest-managed {
		database_type = "postgres_database"
		time_machine_id = data.nutanix_ndb_time_machines.test1.time_machines.0.id
		compute_profile_id =  local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
		network_profile_id = local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
		nx_cluster_id = local.clusters.NDBCluster.id
		vm_password = "pass"
		postgres_database {
			vm_name = "%[1]s"
			client_public_key = "%[3]s"
		}
		description = "%[2]s"

	}
	`, name, desc, sshKey)
}
