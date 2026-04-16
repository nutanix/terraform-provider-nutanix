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




#pull all clusters data
data "nutanix_clusters_v2" "clusters" {}

#create local variable pointing to desired cluster
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

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


# create PBR with vpc name with any source or destination or protocol with permit action
resource "nutanix_pbr_v2" "pbr1" {
  name        = "routing_policy_any_source_destination"
  description = "routing policy with any source and destination"
  vpc_ext_id  = nutanix_vpc_v2.vpc.id
  priority    = 11
  policies {
    policy_match {
      source {
        address_type = "ANY"
      }
      destination {
        address_type = "ANY"
      }
      protocol_type = "UDP"
    }
    policy_action {
      action_type = "PERMIT"
    }
  }
}



# create PBR with vpc uuid with source external
resource "nutanix_pbr_v2" "pbr2" {
  name        = "routing_policy_external_source"
  description = "routing policy with external source"
  vpc_ext_id  = nutanix_vpc_v2.vpc.id
  priority    = 12
  policies {
    policy_match {
      source {
        address_type = "EXTERNAL"
      }
      destination {
        address_type = "SUBNET"
        subnet_prefix {
          ipv4 {
            ip {
              value         = "10.10.10.0"
              prefix_length = 24
            }
          }
        }
      }
      protocol_type = "ANY"
    }
    policy_action {
      action_type = "FORWARD"
      nexthop_ip_address {
        ipv4 {
          value = "10.10.10.10"
        }
      }
    }
  }
}


#create PBR with vpc name with source Any and destination external
resource "nutanix_pbr_v2" "pbr3" {
  name        = "routing_policy_any_source_external_destination"
  description = "routing policy with any source and external destination"
  vpc_ext_id  = nutanix_vpc_v2.vpc.id
  priority    = 13
  policies {
    policy_match {
      source {
        address_type = "ANY"
      }
      destination {
        address_type = "EXTERNAL"
      }
      protocol_type = "UDP"
    }
    policy_action {
      action_type = "PERMIT"
    }
  }
}

# list pbrs
data "nutanix_pbrs_v2" "list-pbrs" {
  depends_on = [nutanix_pbr_v2.pbr1, nutanix_pbr_v2.pbr2, nutanix_pbr_v2.pbr3]
}

# get an entity with pbr uuid
data "nutanix_pbr_v2" "get-pbr" {
  ext_id = nutanix_pbr_v2.pbr1.id
}


