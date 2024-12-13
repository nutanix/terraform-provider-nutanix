terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.3.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = var.nutanix_port
  insecure = true
}

#pull cluster data
data "nutanix_clusters" "clusters" {}

#creating subnet
data "nutanix_subnet" "test" {
  subnet_name = "nutanix-subnet"
  additional_filter {
    name = "name"
    values = ["vlan.15421", "nutanix-subnet"]
  }

  additional_filter {
    name = "vlan_id"
    values = ["15421", "123"]
  }
  
  additional_filter {
    name = "cluster_reference.uuid"
    values = ["0005d504-660e-43c3-3507-ac1f6b60292f"]
  }
}

#output the subnet created above
output "subnet" {
  value = data.nutanix_subnet.test
}