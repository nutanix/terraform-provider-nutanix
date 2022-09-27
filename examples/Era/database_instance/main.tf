terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.8.0"
        }
    }
}

#definig nutanix configuration
provider "nutanix"{
  era_username = var.era_username
  era_password = var.era_password
  era_endpoint = var.era_endpoint
  insecure = true
}

## provision PostgreSQL database with single instance

resource "nutanix_era_database_provision" "dbp" {

    // name of database type
    databasetype = "postgres_database"

    // required name of db instance
    name = "test-inst"
    description = "add description"

    // adding the profiles details
    softwareprofileid = "{{ software_profile_id }}"
    softwareprofileversionid =  "{{ software_profile_version_id }}"
    computeprofileid =  "{{ compute_profile_id }}"
    networkprofileid = "{{ network_profile_id }}"
    dbparameterprofileid = "{{ db_parameter_profile_id }}"

    // postgreSQL Info
    postgresql_info{
        listener_port = "{{ listner_port }}"

        database_size= "{{ 200 }}"

        db_password =  "password"

        database_names= "testdb1"
    }

    // era cluster id
    nxclusterid= local.clusters.EraCluster.id

    // ssh-key
    sshpublickey= "{{ ssh-public-key }}"

    // node for single instance
    nodes{
        // name of dbserver vm 
        vmname= "test-era-vm1"

        // network profile id
        networkprofileid= local.network_profiles.DEFAULT_OOB_POSTGRESQL_NETWORK.id
    }

    // time machine info 
    timemachineinfo {
        name= "test-pg-inst"
        description="description of time machine"
        slaid= "{{ sla_id }}"

        // schedule info fields are optional.
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
