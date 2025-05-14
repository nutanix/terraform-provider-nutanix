terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.8.0"
        }
    }
}

#defining nutanix configuration
provider "nutanix"{
  ndb_username = var.ndb_username
  ndb_password = var.ndb_password
  ndb_endpoint = var.ndb_endpoint
  insecure = true
}

## provision PostgreSQL database with single instance

resource "nutanix_ndb_database" "dbp" {

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


## provision HA instance 

resource "nutanix_ndb_database" "dbp" {
    // database type
    databasetype = "postgres_database"

    // database name & descriptio
    name = "test-pg-inst-HA-tf"
    description = "adding description"

    // adding the profiles details
    softwareprofileid = "{{ software_profile_id }}"
    softwareprofileversionid =  "{{ software_profile_version_id }}"
    computeprofileid =  "{{ compute_profile_id }}"
    networkprofileid = "{{ network_profile_id }}"
    dbparameterprofileid = "{{ db_parameter_profile_id }}"

    // required for HA instance
    createdbserver = true
    clustered = true

    // node count (with haproxy server node)
    nodecount= 4 

    // min required details for provisioning HA instance
    postgresql_info{
      listener_port = "5432"

      database_size= "200"

      db_password =  "{{ database password}}"

      database_names= "testdb1"

      ha_instance{
      proxy_read_port= "5001"

      proxy_write_port = "5000"

      cluster_name= "{{ cluster_name }}"

      patroni_cluster_name = " {{ patroni_cluster_name }}"
      }
    }
  
  nxclusterid= "1c42ca25-32f4-42d9-a2bd-6a21f925b725"
  sshpublickey= "{{ ssh_public_key }}"
  
  // nodes are required.

  // HA proxy node 
  nodes{
    properties{
      name =  "node_type"
      value = "haproxy"
    }
    vmname =  "{{ vm name }}"
    nx_cluster_id =  "{{ nx_cluster_id }}"
  }

  // Primary node for read/write ops
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

    vmname = "{{ name of vm }}"
    networkprofileid="{{ network_profile_id }}"
    computeprofileid= "{{ compute_profile_id }}"
    nx_cluster_id=  "{{ nx_cluster_id }}"
  }

  // secondary nodes for read ops
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
    vmname = "{{ name of vm }}"
    networkprofileid="{{ network_profile_id }}"
    computeprofileid= "{{ compute_profile_id }}"
    nx_cluster_id=  "{{ nx_cluster_id }}"
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
    
    vmname = "{{ name of vm }}"
    networkprofileid="{{ network_profile_id }}"
    computeprofileid= "{{ compute_profile_id }}"
    nx_cluster_id=  "{{ nx_cluster_id }}"
  }

  // time machine required 
  timemachineinfo {
    name= "test-pg-inst-HA"
    description=""
    sla_details{
      primary_sla{
        sla_id= "{{ required SLA}}0"
        nx_cluster_ids=  [
          "{{ nx_cluster_id}}"
        ]
      }
    }
    // schedule fields are optional
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
  
  vm_password= "{{ vm_password}}"
  autotunestagingdrive= true
}