package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceNameProfile = "nutanix_ndb_profile.acctest-managed-profile"

func TestAccEra_ByCompute(t *testing.T) {
	name := "test-compute-tf"
	desc := "this is compute desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileConfigByCompute(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameProfile, "versions.#", "1"),
				),
			},
		},
	})
}
func TestAccEra_BySoftware(t *testing.T) {
	t.Skip()
	name := "test-software-tf"
	desc := "this is software desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
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

func TestAccEra_ByDatabaseParams(t *testing.T) {
	name := "test-software-tf"
	desc := "this is software desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraProfileConfigByDatabaseParams(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProfile, "name", name),
					resource.TestCheckResourceAttr(resourceNameProfile, "description", desc),
					resource.TestCheckResourceAttr(resourceNameProfile, "versions.#", "1"),
				),
			},
		},
	})
}

func TestAccEra_ByNetwork(t *testing.T) {
	name := "test-network-tf"
	desc := "this is network desc"
	subnet := testVars.SubnetName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
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
					source_dbserver_id = "d2f12bd9-bc08-4c17-bd00-c0f7d1a48f5c"
					base_profile_version_name = "test1"
					base_profile_version_description= "test1 desc"
				}
				available_cluster_ids= [local.clusters.EraCluster.id]
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
					nx_cluster_id = local.clusters.EraCluster.id
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
