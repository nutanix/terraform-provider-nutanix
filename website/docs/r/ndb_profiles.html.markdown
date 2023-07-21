---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_profile"
sidebar_current: "docs-nutanix-resource-ndb-profile"
description: |-
  This operation submits a request to create, update and delete profiles in Nutanix database service (NDB).
  Note: For 1.8.0 release, only postgress database type is qualified and officially supported.
---

# nutanix_ndb_profile

Provides a resource to create profiles (Software, Network, Database Parameter, Compute) based on the input parameters. 

## Example Usage

```hcl

    // resource to create compute profile

    resource "nutanix_ndb_profile" "computeProfile" {
        name = "compute-tf"
        description = "tf added compute"
        compute_profile{
            cpus = 1
            core_per_cpu = 2
            memory_size = 2
        }
        published= true
    }


    // resource to create database parameter profile

    resource "nutanix_ndb_database_parameter_profile" "dbProfile" {
        name=  "dbParams-tf"
        description = "database description"
        engine_type = "postgres_database"

        // optional args for engine type else will set to default values
        postgres_database {
            max_connections = "100"
            max_replication_slots = "10"
        }
    }


    // resource to create Postgres Database Single Instance  Network profile
    
    resource "nutanix_ndb_profile" "networkProfile" {
    name = "tf-net"
        description = "terraform created"
        engine_type = "postgres_database"
        network_profile{
            topology = "single"
            postgres_database{  
                single_instance{
                    vlan_name = "vlan.154"
                }
            }
        }
        published = true
    }


    // resource to create Postgres Database HA Instance  Network profile

    resource "nutanix_ndb_profile" "networkProfile" {
        name = "tf-net"
        description = "terraform created"
        engine_type = "postgres_database"
        network_profile{
            topology = "cluster"
            postgres_database{  
                ha_instance{
                    num_of_clusters= "1"
                    vlan_name = ["{{ vlanName }}"]
                    cluster_name = ["{{ ClusterName }}"]
                }
            }
        }
        published = true
    }
    

    // resource to create Software Profile

    resource "nutanix_ndb_profile" "softwareProfile" {
        name= "test-software"
        description = "description"
        engine_type = "postgres_database"
        software_profile {
            topology = "single"
            postgres_database{
                source_dbserver_id = "{{ source_dbserver_id }}"
                base_profile_version_name = "test1"
                base_profile_version_description= "test1 desc"
            }
            available_cluster_ids= ["{{ cluster_ids }}"]
        }
        published = true
    }
``` 

## Argument Reference
* `name` : (Required) Name of profile
* `description` : (Optional) Description of profile
* `engine_type` : Engine Type of database
* `published` : (Optional) Publish for all users 
* `compute_profile`: (Optional) Compute Profile
* `software_profile` : (Optional) Software Profile
* `network_profile` : (Optional) Network Profile
* `database_parameter_profile`:  (Optional) Database Parameter Profile

### compute_profile
The following arguments are supported to create a compute profile.

* `cpus`: (Optional) number of vCPUs for the database server VM.
* `core_per_cpu`: (Optional) number of cores per vCPU for the database server VM.
* `memory_size`: (Optional) amount of memory for the database server VM. 

### software_profile
Ensure that you have registered an existing PostgreSQL database server VM with NDB. NDB uses a registered database server VM to create a software profile.

* `topology`: (Required) Topology of software profile. Allowed values are "cluster" and "single"

* `postgres_database`: (Optional) Software profile info about postgres database.
* `postgres_database.source_dbserver_id`: source dbserver id where postgress software will be installed. 
* `postgres_database.base_profile_version_name`: name for the software profile version.
* `postgres_database.base_profile_version_description`: description for the software profile version.
* `postgres_database.os_notes`: a note to provide additional information about the operating system
* `postgres_database.db_software_notes`: a note to provide additional information about the database software.

* `available_cluster_ids`: specify Nutanix clusters where this profile is available.


### network_profile
A network profile specifies the VLAN for the new database server VM. You can add one or more NICs to segment the network traffic of the database server VM or server cluster.

* `topology`: (Required) Topology supported for network profile. Allowed values are "cluster" and "single"

* `postgres_database`: (Optional) Postgres Info to create network profile

* `postgres_database.single_instance`: (Optional) Info for postgres database to create single instance network profile.
* `postgres_database.single_instance.vlan_name`: (Required) specify the VLAN to provide the IP address used to connect the database from the public network.
* `postgres_database.single_instance.enable_ip_address_selection`: (Optional) If Advanced Network Segmentation is enabled, then this vLAN needs to be a static vLAN and needs to be true.

* `postgres_database.ha_instance`: (Optional) Info for craeting Network profile for HA instance
* `postgres_database.ha_instance.vlan_name`: (Required) specify the VLANs for network
* `postgres_database.ha_instance.cluster_name`: (Required) specify the cluster name associated with given VLANs
* `postgres_database.ha_instance.cluster_id`: (Optional) specify the cluster ids associated with given VLANs
* `postgres_database.ha_instance.num_of_clusters`: (Required) number of cluster attached to network profile

* `version_cluster_association`: (Optional) cluster associated with VLAN. this is used with Single instance for postgres database.
* `version_cluster_association.nx_cluster_id`: (Required) cluster id for associated VLAN. 


### database_parameter_profile
A database parameter profile is a template of custom database parameters that you want to apply to your database

* `postgres_database`: (Optional) Database parameters suuported for postgress.
* `postgres_database.max_connections`: (Optional) Determines the maximum number of concurrent connections to the database server. The default is set to 100
* `postgres_database.max_replication_slots`: (Optional) Specifies the maximum number of replication slots that the server can support. The default is zero. wal_level must be set to archive or higher to allow replication slots to be used. Setting it to a lower value than the number of currently existing replication slots will prevent the server from starting.
* `postgres_database.effective_io_concurrency`: (Optional) Sets the number of concurrent disk I/O operations that PostgreSQL expects can be executed simultaneously. Raising this value will increase the number of I/O operations that any individual PostgreSQL session attempts to initiate in parallel.
* `postgres_database.timezone`: (Optional) Sets the time zone for displaying and interpreting time stamps. Defult is UTC .
* `postgres_database.max_prepared_transactions`: (Optional) Sets the maximum number of transactions that can be in the prepared state simultaneously. Setting this parameter to zero (which is the default) disables the prepared-transaction feature. 
* `postgres_database.max_locks_per_transaction`: (Optional) This parameter controls the average number of object locks allocated for each transaction; individual transactions can lock more objects as long as the locks of all transactions fit in the lock table. Default is 64. 
* `postgres_database.max_wal_senders`: (Optional) Specifies the maximum number of concurrent connections from standby servers or streaming base backup clients (i.e., the maximum number of simultaneously running WAL sender processes). The default is 10. 
* `postgres_database.max_worker_processes`: (Optional) Sets the maximum number of background processes that the system can support. The default is 8. 
* `postgres_database.min_wal_size`: (Optional) As long as WAL disk usage stays below this setting, old WAL files are always recycled for future use at a checkpoint, rather than removed. This can be used to ensure that enough WAL space is reserved to handle spikes in WAL usage, for example when running large batch jobs. The default is 80 MB.
* `postgres_database.max_wal_size`: (Optional) Maximum size to let the WAL grow to between automatic WAL checkpoints. The default is 1 GB
* `postgres_database.checkpoint_timeout`: (Optional) Sets the maximum time between automatic WAL checkpoints . High Value gives Good Performance, but takes More Recovery Time, Reboot time. can reduce the I/O load on your system, especially when using large values for shared_buffers. Default is 5min
* `postgres_database.autovacuum`: (Optional) Controls whether the server should run the autovacuum launcher daemon. This is on by default; however, track_counts must also be enabled for autovacuum to work.
* `postgres_database.checkpoint_completion_target`: (Optional) 	
Specifies the target of checkpoint completion, as a fraction of total time between checkpoints. Time spent flushing dirty buffers during checkpoint, as fraction of checkpoint interval . Formula - (checkpoint_timeout - 2min) / checkpoint_timeout. The default is 0.5.
* `postgres_database.autovacuum_freeze_max_age`: (Optional) Age at which to autovacuum a table to prevent transaction ID wraparound. Default is 200000000
* `postgres_database.autovacuum_vacuum_threshold`: (Optional) Min number of row updates before vacuum. Minimum number of tuple updates or deletes prior to vacuum. Take value in KB. Default is 50 .
* `postgres_database.autovacuum_vacuum_scale_factor`: (Optional) Number of tuple updates or deletes prior to vacuum as a fraction of reltuples. Default is 0.2 
* `postgres_database.autovacuum_work_mem`: (Optional) Sets the maximum memory to be used by each autovacuum worker process. Unit is in KB. Default is -1
* `postgres_database.autovacuum_max_workers`: (Optional) Sets the maximum number of simultaneously running autovacuum worker processes. Default is 3
* `postgres_database.autovacuum_vacuum_cost_delay`: (Optional) Vacuum cost delay in milliseconds, for autovacuum. Specifies the cost delay value that will be used in automatic VACUUM operation. Default is 2ms
* `postgres_database.wal_buffers`: (Optional) 
Sets the number of disk-page buffers in shared memory for WAL. The amount of shared memory used for WAL data that has not yet been written to disk. The default is -1.
* `postgres_database.synchronous_commit`: (Optional) Sets the current transaction's synchronization level. Specifies whether transaction commit will wait for WAL records to be written to disk before the command returns a success indication to the client. Default is on.
* `postgres_database.random_page_cost`: (Optional) Sets the planner's estimate of the cost of a nonsequentially fetched disk page. Sets the planner's estimate of the cost of a non-sequentially-fetched disk page. The default is 4.0. 
* `postgres_database.wal_keep_segments`: (Optional) Sets the number of WAL files held for standby servers, Specifies the minimum number of past log file segments kept in the pg_wal directory. Default is 700 .

## Attributes Reference

* `status`: status of profile
* `owner`: owner  of profile
* `latest_version`: latest version of profile 
* `latest_version_id`: latest version id of profile
* `versions`: versions of associated profiles
* `nx_cluster_id`: cluster on which profile created
* `assoc_databases`: associated databases of profiles
* `assoc_db_servers`: associated database servers for associated profiles
* `cluster_availability`: cluster availability of profile


See detailed information in [NDB Profiles](https://www.nutanix.dev/api_references/ndb/#/467d68a88c0d2-create-a-profile) .
