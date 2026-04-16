#######################################################################
#   _  _  ___ ___    _____   ___  _         _  _   _ _____    _  _  ___  _  _   _ _____  __   ____  __
#  | \| |/ __|_  )__| __\ \ / / \| |  ___  | \| | /_\_   _|__| \| |/ _ \| \| | /_\_   _|_\ \ / /  \/  |___
#  | .` | (__ / /___| _| \ V /| .` | |___| | .` |/ _ \| ||___| .` | (_) | .` |/ _ \| ||___\ V /| |\/| (_-<
#  |_|\_|\___/___|  |_|   \_/ |_|\_|       |_|\_/_/ \_\_|    |_|\_|\___/|_|\_/_/ \_\_|     \_/ |_|  |_/__/
#
#######################################################################
# Name: nc2-fvn_nat-nonat-vms
# Description: This module creates NAT & no-NAT VPCs, subnets, downloads an image and creates VMs from it
# Author: Jonoas Werner
# Date: 2025-06-06
# Version: 1.0.0
# Usage: terraform init && terraform plan && terraform apply
# Inputs:
#   NUTANIX_USERNAME: The username for Prism Central
#   NUTANIX_PASSWORD: The password for Prism Central
#   NUTANIX_ENDPOINT: The IP / FQDN for Prism Central
#   NUTANIX_PORT: The port for Prism Central
#   NUTANIX_INSECURE: Enable if the Prism Central instance has no official certificate
#   SSH_PUBLIC_KEY: The public SSH key to be used for VM access
#######################################################################

variable "SSH_PUBLIC_KEY" {
  description = "The public SSH key to be used for VM access"
  type        = string
}

#################################################
# NUTANIX PROVIDER DEFINITION
#################################################

terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.2.0"
    }
  }
}

provider "nutanix" {
  username = var.NUTANIX_USERNAME
  password = var.NUTANIX_PASSWORD
  endpoint = var.NUTANIX_ENDPOINT
  port     = var.NUTANIX_PORT
  insecure = var.NUTANIX_INSECURE
}

#################################################
# GET CLUSTERS DATA
#################################################

data "nutanix_clusters" "clusters" {}

locals {
  cluster1 = [
    for cluster in data.nutanix_clusters.clusters.entities :
    cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
  ][0]
  
  vpc_a_uuid = resource.nutanix_vpc.vpc_a.metadata.uuid
  vpc_b_uuid = resource.nutanix_vpc.vpc_b.metadata.uuid
  
  # Set DNS to match AWS .2 resolver based on the Prism Central endpoint IP
  dns_server = replace(var.NUTANIX_ENDPOINT, "/^([0-9]+\\.[0-9]+)\\.[0-9]+\\.[0-9]+$/", "$1.0.2")

  # Define subnet configurations
  subnet_configs = {
    vpc_a = {
      subnet1 = {
        name = "VPC-A_Subnet-1"
        ip = "192.168.10.0"
        gateway = "192.168.10.1"
        pool_start = "192.168.10.10"
        pool_end = "192.168.10.20"
      }
      subnet2 = {
        name = "VPC-A_Subnet-2"
        ip = "192.168.20.0"
        gateway = "192.168.20.1"
        pool_start = "192.168.20.10"
        pool_end = "192.168.20.20"
      }
    }
    vpc_b = {
      subnet1 = {
        name = "VPC-B_Subnet-1"
        ip = "192.168.30.0"
        gateway = "192.168.30.1"
        pool_start = "192.168.30.10"
        pool_end = "192.168.30.20"
      }
      subnet2 = {
        name = "VPC-B_Subnet-2"
        ip = "192.168.40.0"
        gateway = "192.168.40.1"
        pool_start = "192.168.40.10"
        pool_end = "192.168.40.20"
      }
    }
  }

  # Define VM configurations
  vm_configs = {
    vpc_a = {
      subnet1 = {
        name_prefix = "VPC-A_Sub-01_VM"
        subnet_key = "subnet1"
      }
      subnet2 = {
        name_prefix = "VPC-A_Sub-02_VM"
        subnet_key = "subnet2"
      }
    }
    vpc_b = {
      subnet1 = {
        name_prefix = "VPC-B_Sub-01_VM"
        subnet_key = "subnet1"
      }
      subnet2 = {
        name_prefix = "VPC-B_Sub-02_VM"
        subnet_key = "subnet2"
      }
    }
  }
}

#################################################
# GET EXISTING EXTERNAL NAT SUBNET
#################################################

data "nutanix_subnet" "external-subnet" {
  subnet_name = "overlay-external-subnet-nat"
}

#################################################
# GET TRANSIT VPC
#################################################

data "nutanix_vpc" "transit-vpc" {
  vpc_name = "transit-vpc"
}

#################################################
# CREATE NO-NAT SUBNET IN TRANSIT VPC
#################################################

resource "nutanix_subnet_v2" "external-subnet-nonat" {
  name              = "overlay-external-subnet-nonat"
  description       = "External subnet without NAT"
  subnet_type       = "OVERLAY"
  is_nat_enabled    = false
  is_external       = true
  vpc_reference     = data.nutanix_vpc.transit-vpc.metadata.uuid

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = "100.64.10.0"
        }
        prefix_length = 24
      }
      default_gateway_ip {
        value = "100.64.10.1"
      }
      pool_list {
        start_ip {
          value = "100.64.10.10"
        }
        end_ip {
          value = "100.64.10.100"
        }
      }
    }
  }
}


#################################################
# CREATE NAT VPC
#################################################

# Create single VPC with reference to external NAT subnet and set default route
resource "nutanix_vpc" "vpc_a" {
  name = "VPC-A"
  external_subnet_reference_name = [
    data.nutanix_subnet.external-subnet.name
  ]
}

#################################################
# CREATE NO-NAT VPC
#################################################

# Create single VPC with reference to external no-NAT subnet and set default route
resource "nutanix_vpc" "vpc_b" {
  name = "VPC-B"
  external_subnet_reference_name = [
    nutanix_subnet_v2.external-subnet-nonat.name
  ]

  # Add ERP configuration for VPC-B's subnets
  externally_routable_prefix_list {
    ip = "192.168.30.0"
    prefix_length = 24
  }
  externally_routable_prefix_list {
    ip = "192.168.40.0"
    prefix_length = 24
  }
}


#################################################
# CREATE OVERLAY SUBNETS
#################################################

# Create subnets for VPC-A
resource "nutanix_subnet" "vpc_a_subnets" {
  for_each = local.subnet_configs.vpc_a
  
  name                 = each.value.name
  subnet_type          = "OVERLAY"
  subnet_ip            = each.value.ip
  prefix_length        = 24
  default_gateway_ip   = each.value.gateway
  ip_config_pool_list_ranges = ["${each.value.pool_start} ${each.value.pool_end}"]
  dhcp_domain_name_server_list = [local.dns_server]
  vpc_reference_uuid   = local.vpc_a_uuid
  depends_on = [
    nutanix_vpc.vpc_a
  ]
}

# Create subnets for VPC-B
resource "nutanix_subnet" "vpc_b_subnets" {
  for_each = local.subnet_configs.vpc_b
  
  name                 = each.value.name
  subnet_type          = "OVERLAY"
  subnet_ip            = each.value.ip
  prefix_length        = 24
  default_gateway_ip   = each.value.gateway
  ip_config_pool_list_ranges = ["${each.value.pool_start} ${each.value.pool_end}"]
  dhcp_domain_name_server_list = [local.dns_server]
  vpc_reference_uuid   = local.vpc_b_uuid
  depends_on = [
    nutanix_vpc.vpc_b
  ]
}

# Add a delay to ensure VPC and subnets are fully created
resource "time_sleep" "wait_for_vpc_and_subnets" {
  depends_on = [
    nutanix_vpc.vpc_a,
    nutanix_vpc.vpc_b,
    nutanix_subnet.vpc_a_subnets,
    nutanix_subnet.vpc_b_subnets
  ]
  create_duration = "30s"
}

# Get VPC route tables after VPC and subnets are created
data "nutanix_route_tables_v2" "vpc_a_route_tables" {
  filter = "vpcReference eq '${nutanix_vpc.vpc_a.metadata.uuid}'"
  depends_on = [
    time_sleep.wait_for_vpc_and_subnets
  ]
}

data "nutanix_route_tables_v2" "vpc_b_route_tables" {
  filter = "vpcReference eq '${nutanix_vpc.vpc_b.metadata.uuid}'"
  depends_on = [
    time_sleep.wait_for_vpc_and_subnets
  ]
}

locals {
  vpc_a_route_table = data.nutanix_route_tables_v2.vpc_a_route_tables.route_tables[0].ext_id
  vpc_b_route_table = data.nutanix_route_tables_v2.vpc_b_route_tables.route_tables[0].ext_id
}

resource "nutanix_routes_v2" "vpc_a_default_route" {
  name = "vpc-a-default-route"
  description = "Default route for VPC-A via external subnet"
  vpc_reference = nutanix_vpc.vpc_a.metadata.uuid
  route_table_ext_id = local.vpc_a_route_table
  route_type = "STATIC"

  destination {
    ipv4 {
      ip {
        value = "0.0.0.0"
      }
      prefix_length = 0
    }
  }

  next_hop {
    next_hop_type = "EXTERNAL_SUBNET"
    next_hop_reference = data.nutanix_subnet.external-subnet.metadata.uuid
  }

  depends_on = [
    time_sleep.wait_for_vpc_and_subnets,
    nutanix_vpc.vpc_a,
    nutanix_subnet.vpc_a_subnets,
    data.nutanix_subnet.external-subnet,
    data.nutanix_route_tables_v2.vpc_a_route_tables
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "nutanix_routes_v2" "vpc_b_default_route" {
  name = "vpc-b-default-route"
  description = "Default route for VPC-B via external subnet"
  vpc_reference = nutanix_vpc.vpc_b.metadata.uuid
  route_table_ext_id = local.vpc_b_route_table
  route_type = "STATIC"

  destination {
    ipv4 {
      ip {
        value = "0.0.0.0"
      }
      prefix_length = 0
    }
  }

  next_hop {
    next_hop_type = "EXTERNAL_SUBNET"
    next_hop_reference = nutanix_subnet_v2.external-subnet-nonat.id
  }

  depends_on = [
    time_sleep.wait_for_vpc_and_subnets,
    nutanix_vpc.vpc_b,
    nutanix_subnet.vpc_b_subnets,
    nutanix_subnet_v2.external-subnet-nonat,
    data.nutanix_route_tables_v2.vpc_b_route_tables
  ]

  lifecycle {
    create_before_destroy = true
  }
}

#################################################
# IMPORT OS IMAGES
#################################################

# Ubuntu 24.04
resource "nutanix_image" "ubuntu-24_04_noble-numbat" {
  name = "ubuntu-24_04_noble-numbat"
  source_uri  = "https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
}

data "nutanix_image" "ubuntu-24_04_noble-numbat" {
  image_id = nutanix_image.ubuntu-24_04_noble-numbat.id
}

#################################################
# CREATE VMs
#################################################

# Create VMs for VPC-A
resource "nutanix_virtual_machine" "vpc_a_vms" {
  for_each = {
    for idx, config in flatten([
      for subnet_key, subnet_config in local.vm_configs.vpc_a : [
        for i in range(2) : {
          key = "${subnet_key}-${i}"
          name = "${subnet_config.name_prefix}-${i + 1}"
          subnet_uuid = nutanix_subnet.vpc_a_subnets[subnet_config.subnet_key].metadata.uuid
        }
      ]
    ]) : config.key => config
  }

  name                 = each.value.name
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  cluster_uuid         = local.cluster1
  
  nic_list {
     subnet_uuid = each.value.subnet_uuid
  }

  guest_customization_cloud_init_user_data = base64encode(templatefile("${path.module}/templates/cloud-init-app.tpl", {
    ssh_public_key = var.SSH_PUBLIC_KEY
  }))

  disk_list {
    data_source_reference = {
      kind = "image"
      uuid = nutanix_image.ubuntu-24_04_noble-numbat.id
    }
    device_properties {
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }
      device_type = "DISK"
    }
    disk_size_mib = 5120  # 5GB in MiB
  }

  disk_list {
    disk_size_bytes = 0
    data_source_reference = {}
    device_properties {
      device_type = "CDROM"
      disk_address = {
        device_index = "1"
        adapter_type = "SATA"
      }
    }
  }
}

# Create VMs for VPC-B
resource "nutanix_virtual_machine" "vpc_b_vms" {
  for_each = {
    for idx, config in flatten([
      for subnet_key, subnet_config in local.vm_configs.vpc_b : [
        for i in range(2) : {
          key = "${subnet_key}-${i}"
          name = "${subnet_config.name_prefix}-${i + 1}"
          subnet_uuid = nutanix_subnet.vpc_b_subnets[subnet_config.subnet_key].metadata.uuid
        }
      ]
    ]) : config.key => config
  }

  name                 = each.value.name
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  cluster_uuid         = local.cluster1
  
  nic_list {
     subnet_uuid = each.value.subnet_uuid
  }

  guest_customization_cloud_init_user_data = base64encode(templatefile("${path.module}/templates/cloud-init-db.tpl", {
    ssh_public_key = var.SSH_PUBLIC_KEY
  }))

  disk_list {
    data_source_reference = {
      kind = "image"
      uuid = nutanix_image.ubuntu-24_04_noble-numbat.id
    }
    device_properties {
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }
      device_type = "DISK"
    }
    disk_size_mib = 5120  # 5GB in MiB
  }

  disk_list {
    disk_size_bytes = 0
    data_source_reference = {}
    device_properties {
      device_type = "CDROM"
      disk_address = {
        device_index = "1"
        adapter_type = "SATA"
      }
    }
  }
}

#################################################
# CREATE FLOATING IP
#################################################

locals {
  vpc_a_vm_nics = flatten([
    for vm in nutanix_virtual_machine.vpc_a_vms : [
      vm.nic_list[0].uuid
    ]
  ])
}

resource "nutanix_floating_ip_v2" "fip" {
  for_each = {
    for idx, nic_uuid in local.vpc_a_vm_nics : "vm-${idx}" => nic_uuid
  }
  
  name = "fip-${each.key}"
  description = "Floating IP for VPC-A VM ${each.key}"
  external_subnet_reference = data.nutanix_subnet.external-subnet.metadata.uuid
  vpc_reference = local.vpc_a_uuid
  
  association {
    vm_nic_association {
      vm_nic_reference = each.value
    }
  }
  
  depends_on = [
    nutanix_virtual_machine.vpc_a_vms
  ]
}
