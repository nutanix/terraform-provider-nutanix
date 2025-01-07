package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameProfile = "nutanix_ndb_profile.acctest-managed-profile"

func TestAccEraProfile_ByCompute(t *testing.T) {
	name := "test-compute-tf"
	desc := "this is compute desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileConfigByCompute(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameProfile, "versions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameProfile, "compute_profile.0.cpus", "1"),
					resource.TestCheckResourceAttr(resourceNameProfile, "compute_profile.0.core_per_cpu", "2"),
					resource.TestCheckResourceAttr(resourceNameProfile, "compute_profile.0.memory_size", "2"),
				),
			},
		},
	})
}

func TestAccEraProfile_BySoftware(t *testing.T) {
	t.Skip()
	name := "test-software-tf"
	desc := "this is software desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileConfigBySoftware(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameProfile, "versions.#", "1"),
				),
			},
		},
	})
}

func TestAccEraProfile_ByDatabaseParams(t *testing.T) {
	name := "test-software-tf"
	desc := "this is software desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileConfigByDatabaseParams(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameProfile, "versions.#", "1"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.max_connections", "100"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.max_replication_slots", "10"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.max_wal_senders", "10"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.max_wal_size", "1GB"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.wal_buffers", "-1"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.random_page_cost", "4"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.autovacuum_freeze_max_age", "200000000"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.checkpoint_completion_target", "0.5"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.checkpoint_timeout", "5min"),
					resource.TestCheckResourceAttr(resourceNameProfile, "database_parameter_profile.0.postgres_database.0.max_worker_processes", "8"),
				),
			},
		},
	})
}

func TestAccEraProfile_ByNetwork(t *testing.T) {
	name := "test-network-tf"
	desc := "this is network desc"
	subnet := testVars.SubnetName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileConfigByNetwork(name, desc, subnet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameProfile, "versions.#", "1"),
				),
			},
		},
	})
}

func TestAccEraProfile_ByNetworkHAPostgres(t *testing.T) {
	name := "test-network-tf"
	desc := "this is network desc for HA postgres"
	subnet := testVars.SubnetName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileConfigByNetworkHA(name, desc, subnet),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameProfile, "versions.#", "1"),
				),
			},
		},
	})
}

func testAccEraProfileConfigByCompute(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_ndb_profile" "acctest-managed-profile" {
			name = "%[1]s"
			description = "%[2]s"
			compute_profile{
			cpus = 1
			core_per_cpu = 2
			memory_size = 2
			}
			published= true
		}
		`, name, desc)
}

func testAccEraProfileConfigBySoftware(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_clusters" "clusters"{}

		locals{
			clusters = {
				for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
			}
		}
		resource "nutanix_ndb_profile" "name12" {
			name= "%[1]s"
			description = "%[2]s"
			engine_type = "postgres_database"
			software_profile {
				topology = "single"
				postgres_database{
					source_dbserver_id = ""
					base_profile_version_name = "test1"
					base_profile_version_description= "test1 desc"
				}
				available_cluster_ids= [local.clusters.NDBCluster.id]
			}
			published = true
		}
	`, name, desc)
}

func testAccEraProfileConfigByNetwork(name, desc, subnet string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_clusters" "clusters"{}

		locals{
			clusters = {
				for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
			}
		}
		resource "nutanix_ndb_profile" "acctest-managed-profile" {
			name = "%[1]s"
			description = "%[2]s"
			engine_type = "postgres_database"
			network_profile{
				topology = "single"
				postgres_database{
					single_instance{
						vlan_name = "%[3]s"
					}
				}
				version_cluster_association{
					nx_cluster_id = local.clusters.NDBCluster.id
				}
			}
			published = true
		}
	`, name, desc, subnet)
}

func testAccEraProfileConfigByDatabaseParams(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_clusters" "clusters"{}

		locals{
			clusters = {
				for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
			}
		}
		resource "nutanix_ndb_profile" "acctest-managed-profile" {
			name = "%[1]s"
			description = "%[2]s"
			engine_type = "postgres_database"
			database_parameter_profile {
				postgres_database {
			       	max_connections = "100"
				    max_replication_slots = "10"
				}
			}
			published = true
		}
	`, name, desc)
}

func testAccEraProfileConfigByNetworkHA(name, desc, subnet string) string {
	return fmt.Sprintf(`
		data "nutanix_ndb_clusters" "clusters"{}

		locals{
			clusters = {
				for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
			}
		}
		resource "nutanix_ndb_profile" "acctest-managed-profile" {
			name = "%[1]s"
			description = "%[2]s"
			engine_type = "postgres_database"
			network_profile{
			  topology = "cluster"
			  postgres_database{
				ha_instance{
				 num_of_clusters= "1"
				 vlan_name = ["%[3]s"]
				 cluster_name = [local.clusters.NDBCluster.name]
				}
			  }
			}
			published = true
		}
	`, name, desc, subnet)
}
