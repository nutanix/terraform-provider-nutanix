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

# create PBR with vpc name with any source or destination or protocol with permit action

resource "nutanix_pbr_v2" "pbr1" {
  name        = "routing_policy"
  description = "routing policy"
  vpc_ext_id  = var.vpc_reference_uuid
  priority    = 14
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
  name        = "routing_policy"
  description = "routing policy"
  vpc_ext_id  = var.vpc_reference_uuid
  priority    = 11
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
  name        = "routing_policy"
  description = "routing policy"
  vpc_ext_id  = var.vpc_reference_uuid
  priority    = 14
  policies {
    policy_match {
      source {
        address_type = "ALL"
      }
      destination {
        address_type = "INTERNET"
      }
      protocol_type = "UDP"
    }
    policy_action {
      action_type = "PERMIT"
    }
  }
}

# list pbr 

data "nutanix_pbrs_v2" "pbrs4" {
}



# get an entity with pbr uuid
data "nutanix_pbr_v2" "pbr5" {
  ext_id = nutanix_pbr_v2.pbr1.id
  depends_on = [
    nutanix_pbr_v2.pbr1
  ]
}

