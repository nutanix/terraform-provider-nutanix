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

// resource to create snapshot with time machine id

resource "nutanix_ndb_database_snapshot" "name" {
  time_machine_id = "{{ tms_ID }}"
  name = "test-snap"
  remove_schedule_in_days = 1
}

// resource to craete snapshot with time machine name

resource "nutanix_ndb_database_snapshot" "name" {
  time_machine_name = "{{ tms_name }}"
  name = "test-snap"
  remove_schedule_in_days = 1
}