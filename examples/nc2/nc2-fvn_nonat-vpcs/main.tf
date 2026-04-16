#######################################################################
#   _  _  ___ ___    _____   ___  _         _  _     _  _   _ _____  __   _____  ___    
#  | \| |/ __|_  )__| __\ \ / / \| |  ___  | \| |___| \| | /_\_   _| \ \ / / _ \/ __|___
#  | .` | (__ / /___| _| \ V /| .` | |___| | .` / _ \ .` |/ _ \| |    \ V /|  _/ (__(_-<
#  |_|\_|\___/___|  |_|   \_/ |_|\_|       |_|\_\___/_|\_/_/ \_\_|     \_/ |_|  \___/__/
#                                                                             
#######################################################################
# Name: nc2-fvn_nonat-vpcs
# Description: This module creates VPCs and subnets with no-NAT connectivity
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
#######################################################################


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
  
  vpc_uuid = [resource.nutanix_vpc.vpc_tf[0].metadata.uuid,resource.nutanix_vpc.vpc_tf[1].metadata.uuid]
  
  # Set DNS to match AWS .2 resolver based on the Prism Central endpoint IP
  dns_server = replace(var.NUTANIX_ENDPOINT, "/^([0-9]+\\.[0-9]+)\\.[0-9]+\\.[0-9]+$/", "$1.0.2")
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
# CREATE VPCS
#################################################

# Create two VPCs with reference to external no-NAT subnet
resource "nutanix_vpc" "vpc_tf" {
  count = length(var.VPC)
  name = var.VPC[count.index]
  external_subnet_reference_name = [
    nutanix_subnet_v2.external-subnet-nonat.name
  ]
  
  # Add ERP configuration for each VPC's own subnets
  externally_routable_prefix_list {
    ip = "192.168.${count.index + 1}0.0"
    prefix_length = 24
  }
  externally_routable_prefix_list {
    ip = "192.168.${count.index + 3}0.0"
    prefix_length = 24
  }
}


#################################################
# CREATE OVERLAY SUBNETS
#################################################


# Create two Overlay Subnet-A attached to each VPC
resource "nutanix_subnet" "subnetOverlay" {
  count = length(var.SUBNET_A)
  name                 = var.SUBNET_A[count.index]
  subnet_type                = "OVERLAY"
  subnet_ip                  = "192.168.${count.index + 1}0.0"
  prefix_length              = 24
  default_gateway_ip         = "192.168.${count.index + 1}0.1"
  ip_config_pool_list_ranges = ["192.168.${count.index + 1}0.10 192.168.${count.index + 1}0.20"]
  dhcp_domain_name_server_list = [local.dns_server]
  vpc_reference_uuid = "${local.vpc_uuid[count.index]}"
  depends_on = [
    nutanix_vpc.vpc_tf
  ]
}

# Create two Overlay Subnet-B attached to each VPC
resource "nutanix_subnet" "subnetOverlayB" {
  count = length(var.SUBNET_B)
  name                 = var.SUBNET_B[count.index]
  subnet_type                = "OVERLAY"
  subnet_ip                  = "192.168.${count.index + 3}0.0"
  prefix_length              = 24
  default_gateway_ip         = "192.168.${count.index + 3}0.1"
  ip_config_pool_list_ranges = ["192.168.${count.index + 3}0.10 192.168.${count.index + 3}0.20"]
  dhcp_domain_name_server_list = [local.dns_server]
  vpc_reference_uuid = "${local.vpc_uuid[count.index]}"
  depends_on = [
    nutanix_vpc.vpc_tf
  ]
}

