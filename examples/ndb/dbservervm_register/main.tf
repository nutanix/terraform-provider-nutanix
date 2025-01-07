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


## resource to register dbserver vm
resource "nutanix_ndb_register_dbserver" "name" {
  vm_ip= "{{ vmip }}"
  nxcluster_id = "{{ cluster_id }}"
  username= "{{ era_driver_user}}"
  password="{{ password }}"
  database_type = "postgres_database"
  postgres_database{
	listener_port  = 5432
    // directory where the PostgreSQL database software is installed
	postgres_software_home= "{{ directory }}"
  }
}