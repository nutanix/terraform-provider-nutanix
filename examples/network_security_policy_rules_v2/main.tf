terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.1"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# Create two temporary categories to use as isolation groups in the example.
resource "nutanix_category_v2" "first_group" {
  key         = "example-category-1"
  value       = "first-group"
  description = "Temporary category for first isolation group"
}

resource "nutanix_category_v2" "second_group" {
  key         = "example-category-2"
  value       = "second-group"
  description = "Temporary category for second isolation group"
}

# Create a network security policy with rules
resource "nutanix_network_security_policy_v2" "example" {
  name        = "two-env-isolation-example"
  description = "Example policy for listing rules via nutanix_network_security_policy_rules_v2"
  state       = "SAVE"
  type        = "ISOLATION"
  rules {
    type = "TWO_ENV_ISOLATION"
    spec {
      two_env_isolation_rule_spec {
        first_isolation_group = [
          nutanix_category_v2.first_group.id,
        ]
        second_isolation_group = [
          nutanix_category_v2.second_group.id,
        ]
      }
    }
  }
}

output "policy_id" {
  value = nutanix_network_security_policy_v2.example.id
}

data "nutanix_network_security_policy_rules_v2" "rules" {
  policy_ext_id = nutanix_network_security_policy_v2.example.ext_id
}

output "network_security_policy_rules" {
  value = data.nutanix_network_security_policy_rules_v2.rules
}
