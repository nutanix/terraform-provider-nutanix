package ndb_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

const resourceNameDB = "nutanix_ndb_database.acctest-managed"

func TestAccEra_basic(t *testing.T) {
	r := acc.RandIntBetween(1, 10)
	name := fmt.Sprintf("test-pg-inst-tf-%d", r)
	desc := "this is desc"
	vmName := fmt.Sprintf("testvm-%d", r)
	sshKey := testVars.SSHKey
	updatedName := fmt.Sprintf("test-pg-inst-tf-updated-%d", r)
	updatedesc := "this is updated desc"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccEraPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixNDBDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseConfig(name, desc, vmName, sshKey, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDB, "name", name),
					resource.TestCheckResourceAttr(resourceNameDB, "description", desc),
					resource.TestCheckResourceAttr(resourceNameDB, "databasetype", "postgres_database"),
					resource.TestCheckResourceAttr(resourceNameDB, "database_nodes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameDB, "time_machine_id"),
					resource.TestCheckResourceAttrSet(resourceNameDB, "time_machine.#"),
				),
			},
			{
				Config: testAccEraDatabaseConfig(updatedName, updatedesc, vmName, sshKey, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDB, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameDB, "description", updatedesc),
					resource.TestCheckResourceAttr(resourceNameDB, "databasetype", "postgres_database"),
					resource.TestCheckResourceAttr(resourceNameDB, "database_nodes.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameDB, "time_machine_id"),
					resource.TestCheckResourceAttrSet(resourceNameDB, "time_machine.#"),
				),
			},
		},
	})
}

func TestAccEraDatabaseProvisionHA(t *testing.T) {
	r := acc.RandIntBetween(11, 25)
	name := fmt.Sprintf("test-pg-inst-HA-tf-%d", r)
	desc := "this is desc"
	sshKey := testVars.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccEraPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixNDBDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseHAConfig(name, desc, sshKey, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDB, "name", name),
					resource.TestCheckResourceAttr(resourceNameDB, "description", desc),
					resource.TestCheckResourceAttr(resourceNameDB, "databasetype", "postgres_database"),
					resource.TestCheckResourceAttr(resourceNameDB, "database_nodes.#", "3"),
					resource.TestCheckResourceAttr(resourceNameDB, "linked_databases.#", "4"),
					resource.TestCheckResourceAttrSet(resourceNameDB, "time_machine_id"),
				),
			},
		},
	})
}

func TestAccEra_SchemaValidationwithCreateDBserver(t *testing.T) {
	r := acc.RandIntBetween(1, 10)
	name := fmt.Sprintf("test-pg-inst-tf-%d", r)
	desc := "this is desc"
	vmName := fmt.Sprintf("testvm-%d", r)
	sshKey := testVars.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccEraPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixNDBDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccEraDatabaseSchemaValidationConfig(name, desc, vmName, sshKey, r),
				ExpectError: regexp.MustCompile(`missing required fields are \[softwareprofileid softwareprofileversionid computeprofileid networkprofileid dbparameterprofileid\] for ndb_provision_database`),
			},
		},
	})
}

func TestAccEra_SchemaValidationwithCreateDBserverFalse(t *testing.T) {
	r := acc.RandIntBetween(1, 10)
	name := fmt.Sprintf("test-pg-inst-tf-%d", r)
	desc := "this is desc"
	vmName := fmt.Sprintf("testvm-%d", r)
	sshKey := testVars.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccEraPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixNDBDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccEraDatabaseSchemaValidationConfigWithoutCreateDBserver(name, desc, vmName, sshKey, r),
				ExpectError: regexp.MustCompile(`missing required fields are \[dbparameterprofileid timemachineinfo\] for ndb_provision_database`),
			},
		},
	})
}

func testAccEraDatabaseConfig(name, desc, vmName, sshKey string, r int) string {
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
		database_parameter_profiles = {
			for p in local.profiles_by_type.Database_Parameter: p.name => p
		}
		software_profiles = {
			for p in local.profiles_by_type.Software: p.name => p
		}
		slas = {
			for p in data.nutanix_ndb_slas.slas.slas: p.name => p
		}
		clusters = {
			for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
		}
		// Use NDBCluster if it exists, otherwise pick the first available cluster
		target_cluster = contains(keys(local.clusters), "NDBCluster") ? local.clusters.NDBCluster : data.nutanix_ndb_clusters.clusters.clusters[0]
	}

	resource "nutanix_ndb_database" "acctest-managed" {
		databasetype = "postgres_database"
		name = "%[1]s"
		description = "%[2]s"
		softwareprofileid = local.software_profiles["POSTGRES_15.6_ROCKY_LINUX_8_OOB"].id
		softwareprofileversionid =  local.software_profiles["POSTGRES_15.6_ROCKY_LINUX_8_OOB"].latest_version_id
		computeprofileid =  local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
		networkprofileid = local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
		dbparameterprofileid = local.database_parameter_profiles.DEFAULT_POSTGRES_PARAMS.id

		postgresql_info{
			listener_port = "5432"
			database_size= "200"
			db_password =  "password"
			database_names= "testdb1"
		}
		nxclusterid= local.target_cluster.id
		sshpublickey= "%[4]s"
		nodes{
				vmname= "%[3]s"
				networkprofileid= local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			}
		timemachineinfo {
			name= "test-pg-inst-%[5]d"
			description=""
			slaid=local.slas["DEFAULT_OOB_BRONZE_SLA"].id
			schedule {
				snapshottimeofday{
					hours= 16
					minutes= 0
					seconds= 0
				}
			continuousschedule{
					enabled=true
					logbackupinterval= 30
					snapshotsperday=1
				}
			weeklyschedule{
					enabled=true
					dayofweek= "WEDNESDAY"
				}
			monthlyschedule{
					enabled = true
					dayofmonth= "27"
				}
			quartelyschedule{
					enabled=true
					startmonth="JANUARY"
					dayofmonth= 27
				}
			yearlyschedule{
					enabled= false
					dayofmonth= 31
					month="DECEMBER"
				}
			}
	  }
	}
	`, name, desc, vmName, sshKey, r)
}

func testAccEraDatabaseHAConfig(name, desc, sshKey string, r int) string {
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
		database_parameter_profiles = {
			for p in local.profiles_by_type.Database_Parameter: p.name => p
		}
		software_profiles = {
			for p in local.profiles_by_type.Software: p.name => p
		}
		slas = {
			for p in data.nutanix_ndb_slas.slas.slas: p.name => p
		}
		clusters = {
			for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
		}
		// Use NDBCluster if it exists, otherwise pick the first available cluster
		target_cluster = contains(keys(local.clusters), "NDBCluster") ? local.clusters.NDBCluster : data.nutanix_ndb_clusters.clusters.clusters[0]
	}

	resource "nutanix_ndb_database" "acctest-managed" {
		databasetype = "postgres_database"
		name = "%[1]s"
		description = "%[2]s"
		softwareprofileid = local.software_profiles["POSTGRES_15.6_ROCKY_LINUX_8_OOB"].id
		softwareprofileversionid =  local.software_profiles["POSTGRES_15.6_ROCKY_LINUX_8_OOB"].latest_version_id
		computeprofileid =  local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
		networkprofileid = local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
		dbparameterprofileid = local.database_parameter_profiles.DEFAULT_POSTGRES_PARAMS.id

		createdbserver = true
		nodecount= 4
		clustered = true

		postgresql_info{
			listener_port = "5432"
			database_size= "200"
			db_password =  "password"
			database_names= "testdb1"
			ha_instance{
				proxy_read_port= "5001"

				proxy_write_port = "5000"

				cluster_name= "ha-cls-%[4]d"

				patroni_cluster_name = "ha-patroni-cluster"
			}
		}
		nxclusterid= local.target_cluster.id
		sshpublickey= "%[3]s"
		nodes{
			properties{
				name =  "node_type"
				value = "haproxy"
			}
			vmname =  "ha-cls_haproxy-%[4]d"
			nx_cluster_id =  local.target_cluster.id
		}
		nodes{
			properties{
				name= "role"
				value=  "Primary"
			}
			properties{
				name= "failover_mode"
				value=  "Automatic"
			}
			properties{
				name= "node_type"
				value=  "database"
			}
			vmname = "ha-cls-1%[4]d"
			networkprofileid=local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			computeprofileid= local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
			nx_cluster_id=  local.target_cluster.id
		}
		nodes{
			properties{
				name= "role"
				value=  "Secondary"
			}
			properties{
				name= "failover_mode"
				value=  "Automatic"
			}
			properties{
				name= "node_type"
				value=  "database"
			}
			vmname = "ha-cls-2%[4]d"
			networkprofileid=local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			computeprofileid= local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
			nx_cluster_id=  local.target_cluster.id
		}

		nodes{
			properties{
				name= "role"
				value=  "Secondary"
			}
			properties{
				name= "failover_mode"
				value=  "Automatic"
			}
			properties{
				name= "node_type"
				value=  "database"
			}
			vmname = "ha-cls-3%[4]d"
			networkprofileid=local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			computeprofileid= local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
			nx_cluster_id= local.target_cluster.id
		}
		timemachineinfo {
			name= "test-pg-inst-%[4]d"
			description=""

			sla_details{
				primary_sla{
				  sla_id= local.slas["DEFAULT_OOB_BRONZE_SLA"].id
				  nx_cluster_ids=  [
					local.target_cluster.id
				  ]
				}
			  }
			schedule {
				snapshottimeofday{
					hours= 16
					minutes= 0
					seconds= 0
				}
			continuousschedule{
					enabled=true
					logbackupinterval= 30
					snapshotsperday=1
				}
			weeklyschedule{
					enabled=true
					dayofweek= "WEDNESDAY"
				}
			monthlyschedule{
					enabled = true
					dayofmonth= "27"
				}
			quartelyschedule{
					enabled=true
					startmonth="JANUARY"
					dayofmonth= 27
				}
			yearlyschedule{
					enabled= false
					dayofmonth= 31
					month="DECEMBER"
				}
			}
	  }
	}
	`, name, desc, sshKey, r)
}

func testAccEraDatabaseSchemaValidationConfig(name, desc, vmName, sshKey string, r int) string {
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
		database_parameter_profiles = {
			for p in local.profiles_by_type.Database_Parameter: p.name => p
		}
		software_profiles = {
			for p in local.profiles_by_type.Software: p.name => p
		}
		slas = {
			for p in data.nutanix_ndb_slas.slas.slas: p.name => p
		}
		clusters = {
			for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
		}
		// Use NDBCluster if it exists, otherwise pick the first available cluster
		target_cluster = contains(keys(local.clusters), "NDBCluster") ? local.clusters.NDBCluster : data.nutanix_ndb_clusters.clusters.clusters[0]
	}

	resource "nutanix_ndb_database" "acctest-managed" {
		databasetype = "postgres_database"
		name = "%[1]s"
		description = "%[2]s"

		postgresql_info{
			listener_port = "5432"
			database_size= "200"
			db_password =  "password"
			database_names= "testdb1"
		}
		nxclusterid= local.target_cluster.id
		sshpublickey= "%[4]s"
		nodes{
				vmname= "%[3]s"
				networkprofileid= local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			}
		timemachineinfo {
			name= "test-pg-inst-%[5]d"
			description=""
			slaid=local.slas["DEFAULT_OOB_BRONZE_SLA"].id
			schedule {
				snapshottimeofday{
					hours= 16
					minutes= 0
					seconds= 0
				}
			continuousschedule{
					enabled=true
					logbackupinterval= 30
					snapshotsperday=1
				}
			weeklyschedule{
					enabled=true
					dayofweek= "WEDNESDAY"
				}
			monthlyschedule{
					enabled = true
					dayofmonth= "27"
				}
			quartelyschedule{
					enabled=true
					startmonth="JANUARY"
					dayofmonth= 27
				}
			yearlyschedule{
					enabled= false
					dayofmonth= 31
					month="DECEMBER"
				}
			}
	  }
	}
	`, name, desc, vmName, sshKey, r)
}

func testAccEraDatabaseSchemaValidationConfigWithoutCreateDBserver(name, desc, vmName, sshKey string, r int) string {
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
			database_parameter_profiles = {
				for p in local.profiles_by_type.Database_Parameter: p.name => p
			}
			software_profiles = {
				for p in local.profiles_by_type.Software: p.name => p
			}
			slas = {
				for p in data.nutanix_ndb_slas.slas.slas: p.name => p
			}
			clusters = {
				for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
			}
		}

		resource "nutanix_ndb_database" "acctest-managed" {
			databasetype = "postgres_database"
			name = "%[1]s"
			description = "%[2]s"

			postgresql_info{
				listener_port = "5432"
				database_size= "200"
				db_password =  "password"
				database_names= "testdb1"
			}
			nxclusterid= local.clusters.NDBCluster.id
			sshpublickey= "%[4]s"
			nodes{
				vmname= "%[3]s"
				networkprofileid= local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			}
			createdbserver=false
		}
	`, name, desc, vmName, sshKey, r)
}

func testAccCheckNutanixNDBDatabaseDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).Era

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_ndb_database" {
			continue
		}

		// 1. Check if Database still exists
		_, err := conn.Service.GetDatabaseInstance(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("NDB Database %s still exists", rs.Primary.ID)
		}

		// 2. Check if associated DB Server VMs still exist
		for k, v := range rs.Primary.Attributes {
			if strings.Contains(k, "dbserver_id") && v != "" {
				_, err := conn.Service.ReadDBServerVM(context.Background(), v)
				if err == nil {
					return fmt.Errorf("NDB DB Server VM %s still exists", v)
				}
			}
			// 3. Check if associated DB Server Cluster (DPC) still exists
			if k == "dbserver_logical_cluster_id" && v != "" {
				_, err := conn.Service.GetDBServerVM(context.Background(), &era.DBServerFilterRequest{
					DbserverClusterID: &v,
				})
				// If GetDBServerVM doesn't return an error for a DPC, it might still exist
				if err == nil {
					return fmt.Errorf("NDB DB Server Cluster %s still exists", v)
				}
			}
		}
	}

	return nil
}
