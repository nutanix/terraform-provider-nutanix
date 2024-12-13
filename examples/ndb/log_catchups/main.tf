terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.8.0-beta.2"
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


## resource to perform log catchup

resource "nutanix_ndb_log_catchups" "name" {
  time_machine_id = "{{ timeMachineID }}"
}