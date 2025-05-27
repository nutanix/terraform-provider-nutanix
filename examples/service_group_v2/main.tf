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


# Add Service  group. with TCP and UDP
resource "nutanix_service_groups_v2" "tcp-udp-service" {
  name        = "service_group_tcp_udp"
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
resource "nutanix_service_groups_v2" "icmp-service" {
  name        = "service_group_icmp"
  description = "service group description"
  icmp_services {
    type = 8
    code = 0
  }
}

# service group with All TCP, UDP and ICMP
resource "nutanix_service_groups_v2" "all-service" {
  name        = "service_group_udp_tcp_icmp"
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
  ext_id = nutanix_service_groups_v2.all-service.ext_id
}

# list all service groups
data "nutanix_service_groups_v2" "list-all-sg" {
  depends_on = [nutanix_service_groups_v2.tcp-udp-service, nutanix_service_groups_v2.icmp-service, nutanix_service_groups_v2.all-service]
}

# list all service groups with filter
data "nutanix_service_groups_v2" "filtered-sg" {
  filter     = "name eq 'service_group_udp_tcp_icmp'"
  depends_on = [nutanix_service_groups_v2.tcp-udp-service, nutanix_service_groups_v2.icmp-service, nutanix_service_groups_v2.all-service]
}

# list all service groups with filter and limit
data "nutanix_service_groups_v2" "filtered-sg-limit" {
  filter     = "startswith(name,'service_group')"
  limit      = 1
  depends_on = [nutanix_service_groups_v2.tcp-udp-service, nutanix_service_groups_v2.icmp-service, nutanix_service_groups_v2.all-service]
}
