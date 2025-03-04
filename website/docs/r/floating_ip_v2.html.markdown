---
layout: "nutanix"
page_title: "NUTANIX: nutanix_floating_ip_v2"
sidebar_current: "docs-nutanix-resource-floating-ip-v2"
description: |-
  Create Floating IPs .
---

# nutanix_floating_ip_v2

Provides Nutanix resource to create Floating IPs.

##  Example1 :  create Floating IP with External Subnet

```hcl

#pull all clusters data
data "nutanix_clusters_v2" "clusters"{}

#create local variable pointing to desired cluster
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

# create external subnet
resource "nutanix_subnet_v2" "ext-subnet"{
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
resource "nutanix_vpc_v2" "vpc"{
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
resource "nutanix_floating_ip_v2" "fip-ext-subnet"{
  name                      = "example-fip"
  description               = "example fip  description"
  external_subnet_reference = nutanix_subnet_v2.ext-subnet.id
  depends_on                = [nutanix_vpc_v2.vpc]
}

```

## Example2 :  create Floating IP with External Subnet with vm association

```hcl
#pull all clusters data
data "nutanix_clusters_v2" "clusters"{}

#create local variable pointing to desired cluster
locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_subnet_v2" "external-nat-subnet"{
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

resource "nutanix_subnet_v2" "overlay-subnet"{
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

resource "nutanix_virtual_machine_v2" "vm"{
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
    backing_info {
      is_connected = true
    }
    network_info {
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
  power_state = "OFF"
  lifecycle {
    ignore_changes = [nics.0.network_info.0.ipv4_config.0.should_assign_ip]
  }
  depends_on = [nutanix_vpc_v2.vm-vpc]
}

resource "nutanix_floating_ip_v2" "fip-ext-subnet-vm"{
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
```

## Argument Reference

The following arguments are supported:

- `name`: (Required) Name of the floating IP.
- `description`: (Optional) Description for the Floating IP.
- `association`: (Optional) Association of the Floating IP with either NIC or Private IP
- `floating_ip`: (Optional) Floating IP address.
- `external_subnet_reference`: (Optional) External subnet reference for the Floating IP to be allocated in on-prem only.
- `vpc_reference`: (Optional) VPC reference UUID
- `vm_nic_reference`: (Optional) VM NIC reference.

### association

- `vm_nic_association`: (Optional) Association of Floating IP with nic
- `vm_nic_association.vm_nic_reference`: (Required) VM NIC reference.
- `vm_nic_association.vpc_reference`: (Optional) VPC reference to which the VM NIC subnet belongs.

- `private_ip_association`: (Optional) Association of Floating IP with private IP
- `private_ip_association.vpc_reference`: (Required) VPC in which the private IP exists.
- `private_ip_association.private_ip`: (Required) An unique address that identifies a device on the internet or a local network in IPv4 or IPv6 format.

### floating_ip

- `ipv4`: Reference to IP Configuration
- `ipv6`: Reference to IP Configuration

### ipv4, ipv6 (Reference to IP Configuration)

- `value`: value of address
- `prefix_length`: Prefix length of the network to which this host IPv4 address belongs. Default value is 32.

## Attributes Reference

The following attributes are exported:

- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
- `metadata`: Metadata associated with this resource.
- `association_status`: Association status of floating IP.
- `external_subnet`: Networking common base object
- `vpc`: Networking common base object
- `vm_nic`: Virtual NIC for projections

See detailed information in [Nutanix Floating IP v4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.0#tag/FloatingIps/operation/createFloatingIp).
