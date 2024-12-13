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

## database_restore with Point in Time

resource "nutanix_ndb_database_restore" "name" {
    database_id= "{{ database_id }}"
    user_pitr_timestamp = "2022-12-28 00:54:30"
    time_zone_pitr = "Asia/Calcutta"
}

## database_restore with snapshot uuid

resource "nutanix_ndb_database_restore" "name" {
    database_id= "{{ database_id }}"
    snapshot_id= "{{ snapshot id }}"
}