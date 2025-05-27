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

## resource to add cluster in time machines for snapshot destination with SLAs

resource "nutanix_ndb_tms_cluster" "cls" {
  time_machine_id = "{{ tms_id }}"
  nx_cluster_id = "{{ cluster_id }}"
  sla_id = "{{ sla_id }}"
}