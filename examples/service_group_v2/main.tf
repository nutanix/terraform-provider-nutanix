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

# Add Service  group.
resource "nutanix_service_groups_v2" "example_1" {
  name        = "service_group"
  description = "service group description"
  tcp_services {
    start_port = "232"
    end_port   = "232"
  }
  udp_services {
    start_port = "232"
    end_port   = "232"
  }
}

# service group with ICMP
resource "nutanix_service_groups_v2" "example_2" {
  name        = "service_group"
  description = "service group description"
  icmp_services {
    type = 8
    code = 0
  }
}

# service group with All
resource "nutanix_service_groups_v2" "example_3" {
  name        = "service_group"
  description = "service group description"
  tcp_services {
    start_port = "232"
    end_port   = "232"
  }
  udp_services {
    start_port = "232"
    end_port   = "232"
  }
  icmp_services {
    type = 8
    code = 0
  }
}


# get service group by ext_id
data "nutanix_service_group_v2" "test" {
  ext_id = nutanix_service_groups_v2.example_3.ext_id
}

# list all service groups
data "nutanix_service_groups_v2" "test" {}

# list all service groups with filter
data "nutanix_service_groups_v2" "test" {
  filter = "name eq 'service_group'"
}
