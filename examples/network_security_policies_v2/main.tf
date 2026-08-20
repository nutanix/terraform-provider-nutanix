terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
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

# To create a Security Policy, please enable Flow in the Prism Central UI. go to Settings > Microsegmentation > check Enable Microsegmentation box

#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}


#list all categories
data "nutanix_categories_v2" "category-list" {}


#creating subnet without IP pool
resource "nutanix_subnet_v2" "vlan-112" {
  name              = "vlan-112"
  description       = "subnet VLAN 112 managed by Terraform with IP pool"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = 122
  is_external       = true
  ip_config {

    ipv4 {
      ip_subnet {
        ip {
          value = "192.168.0.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "192.168.0.1"
      }
      pool_list {
        start_ip {
          value = "192.168.0.20"
        }
        end_ip {
          value = "192.168.0.30"
        }
      }
    }
  }
}

// creating VPC
resource "nutanix_vpc_v2" "vpc" {
  name        = "tf-vpc-example"
  description = "VPC "
  external_subnets {
    subnet_reference = nutanix_subnet_v2.vlan-112.id
  }
  externally_routable_prefixes {
    ipv4 {
      ip {
        value         = "172.30.0.0"
        prefix_length = 32
      }
      prefix_length = 16
    }
  }
}

# Network Security Policy TWO_ENV_ISOLATION Rule
resource "nutanix_network_security_policy_v2" "isolation-nsp" {
  name        = "isolation_policy"
  description = "isolation policy example"
  state       = "SAVE"
  type        = "ISOLATION"
  rules {
    type = "TWO_ENV_ISOLATION"
    spec {
      two_env_isolation_rule_spec {
        first_isolation_group = [
          data.nutanix_categories_v2.category-list.categories.0.ext_id,
        ]
        second_isolation_group = [
          data.nutanix_categories_v2.category-list.categories.1.ext_id,
        ]
      }
    }
  }
  is_hitlog_enabled = true
}

# Network Security Policy APPLICATION Rule with GLOBAL scope (category-based across all VPCs)
resource "nutanix_network_security_policy_v2" "global-application-nsp" {
  name        = "global_application_policy"
  description = "Application policy with GLOBAL scope - VMs resolved by category across all VPCs"
  type        = "APPLICATION"
  state       = "SAVE"
  scope       = "GLOBAL"
  rules {
    description = "global application rule"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.category-list.categories.0.ext_id,
          data.nutanix_categories_v2.category-list.categories.1.ext_id
        ]
        src_category_references = [
          data.nutanix_categories_v2.category-list.categories.2.ext_id
        ]
        is_all_protocol_allowed = true
      }
    }
  }
  is_hitlog_enabled = false
}

# Network Security Policy APPLICATION Rule and INTRA_GROUP Rule (VPC_LIST scope)
resource "nutanix_network_security_policy_v2" "application-nsp" {
  name        = "application_policy"
  description = "application policy example"
  type        = "APPLICATION"
  state       = "SAVE"
  scope       = "VPC_LIST"
  rules {
    description = "test"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.category-list.categories.0.ext_id,
          data.nutanix_categories_v2.category-list.categories.1.ext_id
        ]
        src_category_references = [
          data.nutanix_categories_v2.category-list.categories.2.ext_id
        ]
        is_all_protocol_allowed = true
      }
    }
  }
  rules {
    description = "test22"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          data.nutanix_categories_v2.category-list.categories.3.ext_id,
          data.nutanix_categories_v2.category-list.categories.4.ext_id
        ]
        dest_category_references = [
          data.nutanix_categories_v2.category-list.categories.5.ext_id
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
          data.nutanix_categories_v2.category-list.categories.6.ext_id,
          data.nutanix_categories_v2.category-list.categories.7.ext_id
        ]
        secured_group_action = "ALLOW"
      }
    }
  }

  vpc_reference = [
    nutanix_vpc_v2.vpc.id
  ]
  is_hitlog_enabled = false
}

# Network Security Policy MULTI_ENV_ISOLATION Rule
resource "nutanix_network_security_policy_v2" "multi-env-isolation-nsp" {
  name        = "multi_env_isolation_policy"
  description = "multi env isolation policy example"
  type        = "ISOLATION"
  state       = "SAVE"
  rules {
    description = "test"
    type        = "MULTI_ENV_ISOLATION"
    spec {
      multi_env_isolation_rule_spec {
        spec {
          all_to_all_isolation_group {
            isolation_group {
              group_category_references = [
                data.nutanix_categories_v2.category-list.categories.0.ext_id,
                data.nutanix_categories_v2.category-list.categories.1.ext_id
              ]
            }
            isolation_group {
              group_category_references = [
                data.nutanix_categories_v2.category-list.categories.2.ext_id,
                data.nutanix_categories_v2.category-list.categories.3.ext_id
              ]
            }
          }
        }
      }
    }
  }

  vpc_reference = [
    nutanix_vpc_v2.vpc.id
  ]
  is_hitlog_enabled = false
  depends_on        = [nutanix_vpc_v2.vpc]
}


# Network Security Policy FLEX rule (rule-centric / SMSP "flex" mode).
# Requires Flow flex mode (SMSP) enabled on Prism Central. Flex rules reference
# entity groups (not raw categories) and are only allowed in WORKLOAD/CRITICAL/
# COREINFRASTRUCTURE/ZONE policies, which require a policy `priority`.

# Entity groups the flex rule applies to / matches as source / destination.
resource "nutanix_entity_group_v2" "flex_applied" {
  name        = "tf-flex-applied"
  description = "applied-to entity group for the flex policy"
  allowed_config {
    entities {
      type              = "VM"
      selected_by       = "CATEGORY_EXT_ID"
      reference_ext_ids = [data.nutanix_categories_v2.category-list.categories.0.ext_id]
    }
  }
}

resource "nutanix_entity_group_v2" "flex_src" {
  name        = "tf-flex-src"
  description = "source entity group for the flex policy"
  allowed_config {
    entities {
      type              = "VM"
      selected_by       = "CATEGORY_EXT_ID"
      reference_ext_ids = [data.nutanix_categories_v2.category-list.categories.1.ext_id]
    }
  }
}

resource "nutanix_entity_group_v2" "flex_dst" {
  name        = "tf-flex-dst"
  description = "destination entity group for the flex policy"
  allowed_config {
    entities {
      type              = "VM"
      selected_by       = "CATEGORY_EXT_ID"
      reference_ext_ids = [data.nutanix_categories_v2.category-list.categories.2.ext_id]
    }
  }
}

resource "nutanix_network_security_policy_v2" "flex-nsp" {
  name        = "flex_policy"
  description = "flex policy example"
  type        = "WORKLOAD"
  state       = "SAVE"
  # Mandatory for flex policy types; 1-349 for user-defined WORKLOAD policies
  # (350 is reserved for the system catch-all). Lower value = higher precedence.
  priority = 300
  rules {
    name               = "allow-app-to-db"
    description        = "flex rule example"
    type               = "FLEX"
    is_logging_enabled = true
    spec {
      flex_rule_spec {
        action                             = "ALLOW"
        direction                          = "IN_OUT"
        priority                           = 1000
        applied_to_entity_group_references = [nutanix_entity_group_v2.flex_applied.id]
        src_entity_group_references        = [nutanix_entity_group_v2.flex_src.id]
        dest_entity_group_references       = [nutanix_entity_group_v2.flex_dst.id]
        is_all_protocol_allowed            = true
      }
    }
  }
  depends_on = [
    nutanix_entity_group_v2.flex_applied,
    nutanix_entity_group_v2.flex_src,
    nutanix_entity_group_v2.flex_dst,
  ]
}

# get network security policies
data "nutanix_network_security_policies_v2" "list-nsps" {
  depends_on = [
    nutanix_network_security_policy_v2.application-nsp,
    nutanix_network_security_policy_v2.isolation-nsp,
    nutanix_network_security_policy_v2.multi-env-isolation-nsp,
    nutanix_network_security_policy_v2.global-application-nsp,
  ]
}

# get network security policies with filter
data "nutanix_network_security_policies_v2" "filtered-nsps" {
  filter = "name eq '${nutanix_network_security_policy_v2.application-nsp.name}'"
}

# get network security policy data by id
data "nutanix_network_security_policy_v2" "get-nsp" {
  ext_id = nutanix_network_security_policy_v2.multi-env-isolation-nsp.id
}

# get network security policy rules
data "nutanix_network_security_policy_rules_v2" "get-nsp-rules" {
  policy_ext_id = nutanix_network_security_policy_v2.multi-env-isolation-nsp.id
}

output "network_security_policy_rules" {
  value = data.nutanix_network_security_policy_rules_v2.get-nsp-rules
}
