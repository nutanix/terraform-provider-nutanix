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

# Network Security Policy APPLICATION Rule and INTRA_GROUP Rule
resource "nutanix_network_security_policy_v2" "application-nsp" {
  name        = "application_policy"
  description = "application policy example"
  type        = "APPLICATION"
  state       = "SAVE"
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
  depends_on        = [ nutanix_vpc_v2.vpc]
}


# get network security policies
data "nutanix_network_security_policies_v2" "list-nsps" {
  depends_on = [nutanix_network_security_policy_v2.application-nsp, nutanix_network_security_policy_v2.isolation-nsp, nutanix_network_security_policy_v2.multi-env-isolation-nsp]
}

# get network security policies with filter
data "nutanix_network_security_policies_v2" "filtered-nsps" {
  filter = "name eq '${nutanix_network_security_policy_v2.application-nsp.name}'"
}

# get network security policy data by id
data "nutanix_network_security_policy_v2" "get-nsp" {
  ext_id = nutanix_network_security_policy_v2.multi-env-isolation-nsp.id
}
