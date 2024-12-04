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

## resource to create Compute Profile

resource "nutanix_ndb_profile" "computeProfile" {
  name = "compute-tf"
  description = "compute description"
  compute_profile{
    cpus = 1
    core_per_cpu = 2
    memory_size = 2
  }
  // optional
  published= true
}

## resource to create Database parameters Profile

resource "nutanix_ndb_database_parameter_profile" "dbProfile" {
  name=  "dbParams-tf"
  description = "database description"
  // required engine type
  engine_type = "postgres_database"

  // optional args for engine type else will set to default values
  postgres_database {
    max_connections = "100"
    max_replication_slots = "10"
  }
}

## resource to create Network Profile

### Postgres Database Single Instance profile
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

### Postgres Database HA Instance profile
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

## resource to create Software Profile

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