package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const resourceRegisterDB = "nutanix_ndb_database.acctest-managed"

func TestAccEra_Registerbasic(t *testing.T) {
	name := "test-pg-inst-tf"
	desc := "this is desc"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEraPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraDatabaseRegisterConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRegisterDB, "name", name),
					resource.TestCheckResourceAttr(resourceRegisterDB, "description", desc),
				),
			},
		},
	})
}

func testAccEraDatabaseRegisterConfig(name, desc string) string {
	return fmt.Sprintf(`
	data "nutanix_ndb_profiles" "p"{
	}
	data "nutanix_ndb_slas" "slas"{}
	data "nutanix_ndb_clusters" "clusters"{}
	
	locals {
		slas = {
			for p in data.nutanix_ndb_slas.slas.slas: p.name => p
		}
		clusters = {
			for p in data.nutanix_ndb_clusters.clusters.clusters: p.name => p
		}  
	}
	
	resource "nutanix_ndb_register_database" "name" {
		database_type = "postgres_database"
		database_name=  "%[1]s"
		description = "%[2]s"
		vm_username = "era"
		vm_password = "pass"
		vm_ip = "10.51.144.226"
		nx_cluster_id = local.clusters.EraCluster.id
		time_machine_info {
		  name= "test-pg-inst-regis"
		  description="tms by terraform"
		  slaid=local.slas["DEFAULT_OOB_BRONZE_SLA"].id
		  schedule {
			snapshottimeofday{
			  hours= 13
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
		postgress_info{
		  listener_port= "5432"
		  db_user= "postgres"
		//   postgres_software_home= "/usr/pgsql-10.4"
		//   software_home= "/usr/pgsql-10.4"
		  db_password ="pass"
		  db_name= "testdb1"
		}
	  }
	`, name, desc)
}
