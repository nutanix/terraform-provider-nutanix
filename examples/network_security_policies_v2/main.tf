terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# get the categories
data "nutanix_categories_v2" "cat" {}

# get the vpcs 
data "nutanix_vpcs_v2" "vpcs" {}

# Application Rule
resource "nutanix_network_security_policy_v2" "example" {
  name        = "network_security_policy"
  description = "network security policy example"
  type        = "APPLICATION"
  state       = "SAVE"
  rules {
    description = "test"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.cat.0.ext_id,
          data.nutanix_categories_v2.cat.1.ext_id
        ]
        src_category_references = [
          data.nutanix_categories_v2.cat.2.ext_id
        ]
        is_all_protocol_allowed = true
      }
    }
  }

  rules {
    type = "INTRA_GROUP"
    spec {
      intra_entity_group_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.cat.6.ext_id,
          data.nutanix_categories_v2.cat.7.ext_id
        ]
        secured_group_action = "ALLOW"
      }
    }
  }

  vpc_reference = [
    data.nutanix_vpcs_v2.test.vpcs.0.ext_id
  ]
  is_hitlog_enabled = false
}

# Isolation Rule

resource "nutanix_network_security_policy_v2" "example-2" {
  name        = "network_security_policy_isolation"
  description = "network security policy example"
  state       = "SAVE"
  type        = "ISOLATION"
  rules {
    type = "TWO_ENV_ISOLATION"
    spec {
      two_env_isolation_rule_spec {
        first_isolation_group = [
          data.nutanix_categories_v2.cat.0.ext_id,
        ]
        second_isolation_group = [
          data.nutanix_categories_v2.cat.1.ext_id,
        ]
      }
    }
  }

  rules {
    description = "test"
    type        = "MULTI_ENV_ISOLATION"
    spec {
      multi_env_isolation_rule_spec {
        spec {
          all_to_all_isolation_group {
            isolation_group {
              group_category_references = [
                data.nutanix_categories_v2.cat.0.ext_id,
                data.nutanix_categories_v2.cat.1.ext_id
              ]
            }
            isolation_group {
              group_category_references = [
                data.nutanix_categories_v2.cat.2.ext_id,
                data.nutanix_categories_v2.cat.3.ext_id
              ]
            }
          }
        }
      }
    }
  }

  is_hitlog_enabled = true
}

# get network security policies
data "nutanix_network_security_policies_v2" "nsps-1" {}

# get network security policies with filter
data "nutanix_network_security_policies_v2" "nsps-2" {
  filter = "name eq '${nutanix_network_security_policy_v2.test.name}'"
}
# get network security policy data by id
data "nutanix_network_security_policy_v2" "nsp" {
  ext_id = nutanix_network_security_policy_v2.example-2.id
}
