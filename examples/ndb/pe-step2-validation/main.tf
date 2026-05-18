terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
    }
  }
}

provider "nutanix" {
  ndb_endpoint = var.ndb_endpoint
  ndb_username = var.ndb_username
  ndb_password = var.ndb_password
  insecure     = true
}

variable "ndb_endpoint" {
  type = string
}

variable "ndb_username" {
  type = string
}

variable "ndb_password" {
  type      = string
  sensitive = true
}

variable "pe_name" {
  type = string
}

variable "pe_cluster_ip" {
  type = string
}

variable "pe_username" {
  type = string
}

variable "pe_password" {
  type      = string
  sensitive = true
}

resource "nutanix_ndb_cluster" "step2" {
  name       = var.pe_name
  cluster_ip = var.pe_cluster_ip
  username   = var.pe_username
  password   = var.pe_password
}
