---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_database"
sidebar_current: "docs-nutanix-resource-ndb-database"
description: |-
  This operation submits a request to create, update and delete database instance in Nutanix database service (NDB).
  Note: For 1.8.0 release, only postgress database type is qualified and officially supported.
---

# nutanix_ndb_database

Provides a resource to create database instance based on the input parameters. For 1.8.0 release, only postgress database type is qualified and officially supported.

## Example Usage

### NDB database resource with new database server VM

```hcl
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
        networkprofileid= "<network-profile-uuid>"
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
```


### NDB database resource to provision HA instance with new database server VM

```hcl
resource "nutanix_ndb_database" "dbp" {
    databasetype = "postgres_database"
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
```

### NDB database resource with registered database server VM

```hcl
  resource "nutanix_ndb_database" "dbp" {

    // name of database type
    databasetype = "postgres_database"

    // required name of db instance
    name = "test-inst"
    description = "add description"

    // adding the profiles details
    dbparameterprofileid = "{{ db_parameter_profile_id }}"

    // required dbserver id 
    dbserverId = "{{ dbserver_id }}"
    createdbserver = false

    // postgreSQL Info
    postgresql_info{
        listener_port = "{{ listner_port }}"

        database_size= "{{ 200 }}"

        db_password =  "password"

        database_names= "testdb1"
    }
	actionarguments{
		name = "host_ip"
		value = "{{ hostIP }}"
	}

    // node for single instance
    nodes{
        // id of dbserver vm 
        dbserverid = "{{ dbserver_id }}"
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
```

## Argument Reference

The following arguments are supported:

* `name`: - (Required) Name of the instance.
* `description`: - (Optional) The description
* `databasetype`: - (Required) Type of database. Valid values: postgres_database
* `softwareprofileid`: - (Optional) ID of software profile
* `softwareprofileversionid`: - (Optional) ID of version in software profile
* `computeprofileid`: - (Optional) ID of compute profile
* `networkprofileid`: - (Optional) ID of network profile
* `dbparameterprofileid`: - (Optional) DB parameters profile ID
* `newdbservertimezone`: - (Optional) Timezone of new DB server VM
* `nxclusterid`: - (Optional) Cluster ID for DB server VM
* `sshpublickey`: - (Optional) public key for ssh access to DB server VM
* `createdbserver`: - (Optional) Set this to create new DB server VM. Default: true
* `dbserverid`: - (Optional) DB server VM ID for creating instance on registered DB server VM
* `clustered`: - (Optional) If clustered database. Default: false
* `autotunestagingdrive`: - (Optional) Enable auto tuning of staging drive. Default: true
* `nodecount`: - (Optional) No. of nodes/db server vms. Default: 1
* `vm_password`: - (Optional) password for DB server VM and era drive user
* `actionarguments`: - (Optional) action arguments for database. For postgress, you can use postgresql_info
* `timemachineinfo`: - (Optional) time machine config
* `nodes`: - (Optional) nodes info
* `postgresql_info`: - (Optional) action arguments for postgress type database.

* `delete`:- (Optional) Delete the database from the VM. Default value is true
* `remove`:- (Optional) Unregister the database from NDB. Default value is true
* `soft_remove`:- (Optional) Soft remove. Default will be false
* `forced`:- (Optional) Force delete of instance. Default is false
* `delete_time_machine`:- (Optional) Delete the database's Time Machine (snapshots/logs) from the NDB. Default value is true
* `delete_logical_cluster`:- (Optional) Delete the logical cluster. Default is true

### actionarguments

Structure for each action argument in actionarguments list:

* `name`: - (Required) name of argument
* `value`: - (Required) value for argument

### nodes

Each node in nodes supports the following:

* `properties`: - (Optional) list of additional properties
* `vmname`: - (Required) name of vm
* `networkprofileid`: - (Required) network profile ID
* `dbserverid`: - (Optional) Database server ID required for existing VM
* `ip_infos` :- (Optional) IP infos for custom network profile. 
* `computeprofileid` :- (Optional) compute profile id
* `nx_cluster_id` :- (Optional) cluster id. 

### timemachineinfo

The timemachineinfo attribute supports the following:

* `name`: - (Required) name of time machine
* `description`: - (Optional) description of time machine
* `slaid`: - (Optional) SLA ID for single instance 
* `sla_details`:-  (optional) SLA details for HA instance
* `autotunelogdrive`: - (Optional) enable auto tune log drive. Default: true
* `schedule`: - (Optional) schedule for snapshots
* `tags`: - (Optional) tags

### sla_details

* `primary_sla`:- (Required) primary sla details
* `primary_sla.sla_id` :- (Required) sla id
* `primary_sla.nx_cluster_ids` -: (Optioanl) cluster ids


### schedule

The schedule attribute supports the following:

* `snapshottimeofday`: - (Optional) daily snapshot config
* `continuousschedule`: - (Optional) snapshot freq and log config
* `weeklyschedule`: - (Optional) weekly snapshot config
* `monthlyschedule`: - (Optional) monthly snapshot config
* `quartelyschedule`: - (Optional) quaterly snapshot config
* `yearlyschedule`: - (Optional) yearly snapshot config

### snapshottimeofday

The snapshottimeofday attribute supports the following for HH:MM:SS time snapshot needs to be taken:

* `hours`: - (Required) hours
* `minutes`: - (Required) minutes
* `seconds`: - (Required) seconds

### continuousschedule

The continuousschedule attribute supports the following:

* `enabled`: - (Required) to enable
* `logbackupinterval`: - (Required) log catchup interval for database
* `snapshotsperday`: - (Required) num of snapshots per day

### weeklyschedule

The weeklyschedule attribute supports the following:

* `enabled`: - (Required) to enable
* `dayofweek`: - (Required) day of week to take snaphsot. Eg. "WEDNESDAY"

### monthlyschedule

The monthlyschedule attribute supports the following:

* `enabled`: - (Required) to enable
* `dayofmonth`: - (Required) day of month to take snapshot

### quartelyschedule

The quartelyschedule attribute supports the following:

* `enabled`: - (Required) to enable
* `startmonth`: - (Required) quarter start month
* `dayofmonth`: - (Required) month's day for snapshot

### yearlyschedule

The yearlyschedule attribute supports the following:

* `enabled`: - (Required) to enable
* `month`: - (Required) month for snapshot 
* `dayofmonth`: - (Required) day of month to take snapshot

### postgresql_info

The postgresql_info attribute supports the following:

* `listener_port`: - (Required) listener port for database instance
* `database_size`: - (Required) initial database size
* `auto_tune_staging_drive`: - (Optional) enable auto tuning of staging drive. Default: false
* `allocate_pg_hugepage`: - (Optional) allocate huge page. Default: false
* `cluster_database`: - (Optional) if clustered database. Default: false
* `auth_method`: - (Optional) auth methods. Default: md5
* `database_names`: - (Required) name of initial database to be created
* `db_password`: - (Required) database instance password
* `pre_create_script`: - (Optional) pre instance create script
* `post_create_script`: - (Optional) post instance create script
* `ha_instance` :- (Optional) High Availability instance

### ha_instance

* `cluster_name` :- (Required) cluster name
* `patroni_cluster_name` :- (Required) patroni cluster name
* `proxy_read_port` :-  (Required) proxy read port
* `proxy_write_port` :- (Required) proxy write port
* `provision_virtual_ip` :- (Optional) provisional virtual ip. Default is set to true
* `deploy_haproxy` :- (Optional) HA proxy node. Default is set to false
* `enable_synchronous_mode` :- (Optional) enable synchronous mode. Default is set to true
* `failover_mode` :- (Optional) failover mode of nodes. 
* `node_type` :- (Optional) node type of instance. Default is set to database
* `archive_wal_expire_days` :- (Optional) archive wal expire days. Default is set to -1 
* `backup_policy` :- (Optional) backup policy for instance. Default is "primary_only"
* `enable_peer_auth` :- (Optional) enable peer auth . Default is set to false. 

## lifecycle

* `Update` : - Currently only update of instance's name and description is supported using this resource

See detailed information in [NDB Database Instance](https://www.nutanix.dev/api_references/ndb/#/9d9eee4304496-provision-a-database).
