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

## resource to refresh clone with snapshot id

resource "nutanix_ndb_clone_refresh" "acctest-managed"{
    clone_id = "{{ clone_id }}"
    snapshot_id = "{{ snapshot_id }}"
    timezone = "Asia/Calcutta"
}

## resource to refresh clone with userpitr timestamp

resource "nutanix_ndb_clone_refresh" "acctest-managed"{
    clone_id = "{{ clone_id }}"
    user_pitr_timestamp = "{{ timestamp }}"
    timezone = "Asia/Calcutta"
}
