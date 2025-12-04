#######################################################################
#  _  _  ___ ___    _____   ___  _         _  _     _  _   _ _____ 
# | \| |/ __|_  )__| __\ \ / / \| |  ___  | \| |___| \| | /_\_   _|
# | .` | (__ / /___| _| \ V /| .` | |___| | .` / _ \ .` |/ _ \| |  
# |_|\_|\___/___|  |_|   \_/ |_|\_|       |_|\_\___/_|\_/_/ \_\_|  
#                                                                  
#######################################################################
# Name: nc2-fvn_nonat-only
# Description: This module creates a no-NAT subnet in the transit VPC
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

# Export the following environment variables:
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
}


#################################################
# GET EXISTING TRANSIT VPC
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

