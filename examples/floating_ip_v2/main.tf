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

# Example1 :  create Floating IP with External Subnet
# create external subnet
resource "nutanix_subnet_v2" "ext-subnet" {
  name              = "tf-example-subnet-floating-ip"
  description       = "example subnet managed by Terraform with IP pool"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = 129
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

# create VPC
resource "nutanix_vpc_v2" "vpc" {
  name        = "tf-vpc-floating-ip"
  description = "example vpc managed by Terraform"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.ext-subnet.id
  }
  common_dhcp_options {
    domain_name_servers {
      ipv4 {
        value         = "8.8.8.9"
        prefix_length = 32
      }
    }
    domain_name_servers {
      ipv4 {
        value         = "8.8.8.8"
        prefix_length = 32
      }
    }
  }

}

# create Floating IP with External Subnet UUID
resource "nutanix_floating_ip_v2" "fip-ext-subnet" {
  name                      = "example-fip"
  description               = "example fip  description"
  external_subnet_reference = nutanix_subnet_v2.ext-subnet.id
  depends_on                = [nutanix_vpc_v2.vpc]
}


# Example2 :  create Floating IP with External Subnet with vm association
resource "nutanix_subnet_v2" "external-nat-subnet" {
  name              = "tf-external-nat-subnet"
  description       = "terraform"
  cluster_reference = local.clusterExtId
  subnet_type       = "VLAN"
  network_id        = 208
  is_external       = true
  is_nat_enabled    = true
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "10.44.3.192"
        }
        prefix_length = 27
      }
      default_gateway_ip {
        value = "10.44.3.193"
      }
      pool_list {
        start_ip {
          value = "10.44.3.198"
        }
        end_ip {
          value = "10.44.3.207"
        }
      }
    }
  }
}

resource "nutanix_vpc_v2" "vm-vpc" {
  name        = "tf-fip-vpc"
  description = "example vpc managed by Terraform"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.external-nat-subnet.id
  }
}

resource "nutanix_subnet_v2" "overlay-subnet" {
  name        = "tf-overlay-subnet"
  subnet_type = "OVERLAY"

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value         = "192.168.1.0"
          prefix_length = 32
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value         = "192.168.1.1"
        prefix_length = 32
      }
    }
  }
  vpc_reference = nutanix_vpc_v2.vm-vpc.id
}

resource "nutanix_virtual_machine_v2" "vm" {
  name              = "tf-example-vm-floating-ip"
  is_agent_vm       = false
  num_sockets       = 1
  memory_size_bytes = 4 * 1024 * 1024 * 1024
  cluster {
    ext_id = local.clusterExtId
  }
  boot_config {
    legacy_boot {
      boot_order = ["CDROM", "DISK", "NETWORK"]
    }
  }
  nics {
    nic_backing_info {
      virtual_ethernet_nic {
        is_connected = true
      }
    }
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        ipv4_config {
          ip_address {
            value = "192.168.1.15"
          }
          should_assign_ip = true
        }
        subnet {
          ext_id = nutanix_subnet_v2.overlay-subnet.id
        }
        vlan_mode = "ACCESS"
      }
    }
  }
  power_state = "OFF"
  lifecycle {
    ignore_changes = [nics.0.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config.0.should_assign_ip]
  }
  depends_on = [nutanix_vpc_v2.vm-vpc]
}

resource "nutanix_floating_ip_v2" "fip-ext-subnet-vm" {
  name                      = "example-fip"
  description               = "example fip  description"
  external_subnet_reference = nutanix_subnet_v2.external-nat-subnet.id
  association {
    vm_nic_association {
      vm_nic_reference = nutanix_virtual_machine_v2.vm.nics[0].ext_id
    }
  }
  depends_on = [nutanix_vpc_v2.vm-vpc]
}

# Example3 :  fetch floating IP data source and list of floating IPs

# get floating IP
data "nutanix_floating_ip_v2" "get-fip" {
  ext_id = nutanix_floating_ip_v2.fip-ext-subnet.ext_id
}

# list of floating IPs
data "nutanix_floating_ips_v2" "list-fips" {
  depends_on = [nutanix_floating_ip_v2.fip-ext-subnet, nutanix_floating_ip_v2.fip-vpc-ext-subnet, nutanix_floating_ip_v2.fip-ext-subnet-vm]
}

# filter floating IPs
data "nutanix_floating_ips_v2" "filter-fips" {
  filter = "name eq '${nutanix_floating_ip_v2.fip-ext-subnet.name}'"
}
