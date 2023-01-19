---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_profile"
sidebar_current: "docs-nutanix-resource-ndb-profile"
description: |-
  This operation submits a request to create, update and delete profiles in Nutanix database service (NDB).
  Note: For 1.8.0-beta.2 release, only postgress database type is qualified and officially supported.
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


See detailed information in [NDB Profiles](https://www.nutanix.dev/api_references/ndb/#/a626231269b79-create-a-profile) .
