# IAM v2 examples: Entity datasource, SAML Identity Provider datasources and resource.
# Set variables in terraform.tfvars or environment.

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = ">= 2.4.1"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

# Get a single entity by ext_id (e.g. from authorization policy entities)
data "nutanix_iam_entity_v2" "example" {
  count  = var.entity_ext_id != "" ? 1 : 0
  ext_id = var.entity_ext_id
}


# List IAM entities (with optional filter/pagination)
data "nutanix_iam_entities_v2" "examples" {
  limit   = var.entities_limit
  filter  = var.entities_filter
  order_by = var.entities_order_by
}


output "entity_name" {
  value       = try(data.nutanix_iam_entity_v2.example[0].name, null)
  description = "Name of the entity fetched by ext_id"
}

output "entity_display_name" {
  value       = try(data.nutanix_iam_entity_v2.example[0].display_name, null)
  description = "Display name of the entity"
}

output "entities_count" {
  value       = length(data.nutanix_iam_entities_v2.examples.entities)
  description = "Number of IAM entities returned by list"
}
