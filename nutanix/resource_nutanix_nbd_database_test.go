package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceNameDB = "nutanix_ndb_database.acctest-managed"

func TestAccNutanixEra_basic(t *testing.T) {
	name := "test-pg-inst-tf"
	desc := "this is desc"
	vmName := "testvm12"
	sshKey := testVars.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixEraDatabaseConfig(name, desc, vmName, sshKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDB, "name", name),
					resource.TestCheckResourceAttr(resourceNameDB, "description", desc),
				),
			},
		},
	})
}

func testAccNutanixEraDatabaseConfig(name, desc, vmName, sshKey string) string {
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
		softwareprofileid = local.software_profiles["POSTGRES_10.4_OOB"].id
		softwareprofileversionid =  local.software_profiles["POSTGRES_10.4_OOB"].latest_version_id
		computeprofileid =  local.compute_profiles["DEFAULT_OOB_SMALL_COMPUTE"].id
		networkprofileid = local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
		dbparameterprofileid = local.database_parameter_profiles.DEFAULT_POSTGRES_PARAMS.id
	
		postgresql_info{
			listener_port = "5432"
			database_size= "200"
			db_password =  "password"
			database_names= "testdb1"
		}
		nxclusterid= local.clusters.EraCluster.id
		sshpublickey= "%[4]s"
		nodes{
				vmname= "%[3]s"
				networkprofileid= local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
			}
		timemachineinfo {
			name= "test-pg-inst"
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
	`, name, desc, vmName, sshKey)
}
