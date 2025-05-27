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

## resource to create tag in ndb

resource "nutanix_ndb_tag" "name" {
  name= "testst-up"
  description = "this is desc ok"
  entity_type = "DATABASE"
  required=true
}

## resource to deprecate the tag
resource "nutanix_ndb_tag" "name" {
  name= "testst-up"
  description = "this is desc ok"
  entity_type = "DATABASE"
  required=true
  status = "DEPRECATED"
}