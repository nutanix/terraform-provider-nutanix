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

## resource to create sla 

resource "nutanix_ndb_sla" "sla" {
    // name is required
    name= "test-sla"
    // desc is optional
    description = "here goes description"
    // Rentention args are optional with default values
    continuous_retention = 30
    daily_retention = 3
    weekly_retention = 2
    monthly_retention= 1
    quarterly_retention=1
}

## data source sla with sla_name
data "nutanix_ndb_sla" "sla"{
    sla_name = "{{ SLA_NAME }}"
}

output "salO1" {
  value = data.nutanix_ndb_sla.sal
}

## data source sla with sla_id
data "nutanix_ndb_sla" "sla"{
      sla_id = "{{ SLA_ID }}"
}

output "salO1" {
  value = data.nutanix_ndb_sla.sla
}

## List SLAs

data "nutanix_ndb_slas" "sla"{}

output "salO" {
  value = data.nutanix_ndb_slas.sla
}