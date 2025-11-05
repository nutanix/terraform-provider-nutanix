#######################################################################
#  _  _  ___ ___    _____   ___  _         _  _   _ _____  __   _____  ___    
# | \| |/ __|_  )__| __\ \ / / \| |  ___  | \| | /_\_   _| \ \ / / _ \/ __|___
# | .` | (__ / /___| _| \ V /| .` | |___| | .` |/ _ \| |    \ V /|  _/ (__(_-<
# |_|\_|\___/___|  |_|   \_/ |_|\_|       |_|\_/_/ \_\_|     \_/ |_|  \___/__/
#                                                                             
#######################################################################
# Name: nc2-fvn_nat-vpcs
# Description: This module creates NAT enabled VPCs and subnets
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
# CREATE VPCS
#################################################

# Create two VPCs with reference to external NAT subnet
resource "nutanix_vpc" "vpc_tf" {
  count = length(var.VPC)
  name = var.VPC[count.index]
  external_subnet_reference_name = [
    data.nutanix_subnet.external-subnet.name
  ]
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

