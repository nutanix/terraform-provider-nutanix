terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0"
    }
  }
}

#definig nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

# List all the hosts
data "nutanix_host_v2" "hosts"{
  filter = "cluster/name eq '<cluster name>'"
}

# Get the host details
data "nutanix_host_v2" "host"{
  cluster_ext_id = "<cluster uuid>"
  ext_id = "<host uuid>"
}