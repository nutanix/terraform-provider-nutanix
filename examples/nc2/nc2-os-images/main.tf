#######################################################################
#  _  _  ___ ___    _____   ___  _          ___  ___   ___ __  __   _   ___ ___ ___ 
# | \| |/ __|_  )__| __\ \ / / \| |  ___   / _ \/ __| |_ _|  \/  | /_\ / __| __/ __|
# | .` | (__ / /___| _| \ V /| .` | |___| | (_) \__ \  | || |\/| |/ _ \ (_ | _|\__ \
# |_|\_|\___/___|  |_|   \_/ |_|\_|        \___/|___/ |___|_|  |_/_/ \_\___|___|___/
#                                                                                                                                                                     
#######################################################################
# Name: nc2-os-images
# Description: This module imports OS images into NC2 from web sources
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
}

#################################################
# IMPORT OS IMAGES
#################################################

# Ubuntu 24.04
resource "nutanix_image" "ubuntu-24_04_noble-numbat" {
  name = "ubuntu-24_04_noble-numbat"
  source_uri  = "https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
}

# Ubuntu 22.04
resource "nutanix_image" "ubuntu-22_04_jammy-jellyfish" {
  name = "ubuntu-22_04_jammy-jellyfish"
  source_uri  = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}

# Ubuntu 20.04
resource "nutanix_image" "ubuntu-20_04_focal-fossa" {
  name = "ubuntu-20_04_focal-fossa"
  source_uri  = "https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img"
}

# CentOS 10
resource "nutanix_image" "centos-10_cloud-image" {
  name = "centos-10_cloud-image"
  source_uri  = "https://cloud.centos.org/centos/10-stream/x86_64/images/CentOS-Stream-GenericCloud-10-20241118.0.x86_64.qcow2"
}

# CentOS 9
resource "nutanix_image" "centos-9_cloud-image" {
  name = "centos-9_cloud-image"
  source_uri  = "https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-20240527.0.x86_64.qcow2"
}

