#######################################################################
#   _  _  ___ ___    _____   ___  _         _  _   _ _____  __   ____  __        ___ ___ __  __ ___ _    ___ 
#  | \| |/ __|_  )__| __\ \ / / \| |  ___  | \| | /_\_   _|_\ \ / /  \/  |______/ __|_ _|  \/  | _ \ |  | __|
#  | .` | (__ / /___| _| \ V /| .` | |___| | .` |/ _ \| ||___\ V /| |\/| (_-<___\__ \| || |\/| |  _/ |__| _| 
#  |_|\_|\___/___|  |_|   \_/ |_|\_|       |_|\_/_/ \_\_|     \_/ |_|  |_/__/   |___/___|_|  |_|_| |____|___|
#                                                                             
#######################################################################
# Name: nc2-fvn_nat-vms-simple
# Description: This module creates NAT enabled VPCs, subnets, downloads an image and creates VMs from it
# Author: Jonoas Werner
# Date: 2025-06-05
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
  
  vpc_uuid = resource.nutanix_vpc.vpc_tf.metadata.uuid
  
  # Set DNS to match AWS .2 resolver based on the Prism Central endpoint IP
  dns_server = replace(var.NUTANIX_ENDPOINT, "/^([0-9]+\\.[0-9]+)\\.[0-9]+\\.[0-9]+$/", "$1.0.2")
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
# CREATE VPC
#################################################

# Create single VPC with reference to external NAT subnet and set default route
resource "nutanix_vpc" "vpc_tf" {
  name = "VPC-A"
  external_subnet_reference_name = [
    data.nutanix_subnet.external-subnet.name
  ]
}

#################################################
# CREATE OVERLAY SUBNETS
#################################################

# Create first Overlay Subnet (192.168.10.0/24)
resource "nutanix_subnet" "subnetOverlay" {
  name                 = "VPC-A_Subnet-1"
  subnet_type          = "OVERLAY"
  subnet_ip            = "192.168.10.0"
  prefix_length        = 24
  default_gateway_ip   = "192.168.10.1"
  ip_config_pool_list_ranges = ["192.168.10.10 192.168.10.20"]
  dhcp_domain_name_server_list = [local.dns_server]
  vpc_reference_uuid   = local.vpc_uuid
  depends_on = [
    nutanix_vpc.vpc_tf
  ]
}

# Create second Overlay Subnet (192.168.20.0/24)
resource "nutanix_subnet" "subnetOverlayB" {
  name                 = "VPC-A_Subnet-2"
  subnet_type          = "OVERLAY"
  subnet_ip            = "192.168.20.0"
  prefix_length        = 24
  default_gateway_ip   = "192.168.20.1"
  ip_config_pool_list_ranges = ["192.168.20.10 192.168.20.20"]
  dhcp_domain_name_server_list = [local.dns_server]
  vpc_reference_uuid   = local.vpc_uuid
  depends_on = [
    nutanix_vpc.vpc_tf
  ]
}

# Add a delay to ensure VPC and subnets are fully created
resource "time_sleep" "wait_for_vpc_and_subnets" {
  depends_on = [
    nutanix_vpc.vpc_tf,
    nutanix_subnet.subnetOverlay,
    nutanix_subnet.subnetOverlayB
  ]
  create_duration = "30s"
}

# Get VPC route table after VPC and subnets are created
data "nutanix_route_tables_v2" "vpc_a_route_tables" {
  filter = "vpcReference eq '${nutanix_vpc.vpc_tf.metadata.uuid}'"
  depends_on = [
    time_sleep.wait_for_vpc_and_subnets
  ]
}

locals {
  vpc_route_table = data.nutanix_route_tables_v2.vpc_a_route_tables.route_tables[0].ext_id
}

resource "nutanix_routes_v2" "vpc_a_default_route" {
  name = "vpc-a-default-route"
  description = "Default route for VPC-A via external subnet"
  vpc_reference = nutanix_vpc.vpc_tf.metadata.uuid
  route_table_ext_id = local.vpc_route_table
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
    nutanix_vpc.vpc_tf,
    nutanix_subnet.subnetOverlay,
    nutanix_subnet.subnetOverlayB,
    data.nutanix_subnet.external-subnet,
    data.nutanix_route_tables_v2.vpc_a_route_tables
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
# CREATE VM
#################################################

locals {
  subnet_uuid_a = resource.nutanix_subnet.subnetOverlay.metadata.uuid
  subnet_uuid_b = resource.nutanix_subnet.subnetOverlayB.metadata.uuid
}

# Create two VMs in Subnet-1
resource "nutanix_virtual_machine" "vm_subnet1" {
  count = 2
  name                 = "VPC-A_Sub-01_VM-${count.index + 1}"
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  cluster_uuid         = local.cluster1
  
  nic_list {
     subnet_uuid = local.subnet_uuid_a
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

# Create two VMs in Subnet-2
resource "nutanix_virtual_machine" "vm_subnet2" {
  count = 2
  name                 = "VPC-A_Sub-02_VM-${count.index + 1}"
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  cluster_uuid         = local.cluster1
  
  nic_list {
     subnet_uuid = local.subnet_uuid_b
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

#################################################
# CREATE FLOATING IP
#################################################

locals {
  vpc_reference_uuid = resource.nutanix_vpc.vpc_tf.metadata.uuid
  vm_nic_reference_uuid = [
    resource.nutanix_virtual_machine.vm_subnet1[0].nic_list[0].uuid,
    resource.nutanix_virtual_machine.vm_subnet1[1].nic_list[0].uuid,
    resource.nutanix_virtual_machine.vm_subnet2[0].nic_list[0].uuid,
    resource.nutanix_virtual_machine.vm_subnet2[1].nic_list[0].uuid
  ]
}

resource "nutanix_floating_ip_v2" "fip" {
  count = length(local.vm_nic_reference_uuid)
  name = "fip-${count.index + 1}"
  description = "Floating IP for VPC-A VM ${count.index + 1}"
  external_subnet_reference = data.nutanix_subnet.external-subnet.metadata.uuid
  vpc_reference = local.vpc_reference_uuid

  association {
    vm_nic_association {
      vm_nic_reference = local.vm_nic_reference_uuid[count.index]
    }
  }

  depends_on = [
    nutanix_virtual_machine.vm_subnet1,
    nutanix_virtual_machine.vm_subnet2
  ]
}


