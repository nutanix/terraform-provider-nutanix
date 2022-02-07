terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "1.3"
    }
  }
}

provider "nutanix" {
  username = "admin"
  password = "password"
  endpoint = "pc-ip"
  insecure = true
  port     = 9440
}

data "nutanix_clusters" "clusters" {}

output "subnet" {
  value = data.nutanix_subnet.test
}

data "nutanix_subnet" "test" {
  subnet_name = "nutanix-subnet"
  additional_filter {
    name = "name"
    values = ["vlan.154", "nutanix-subnet"]
  }

  additional_filter {
    name = "vlan_id"
    values = ["154", "123"]
  }
  
  additional_filter {
    name = "cluster_reference.uuid"
    values = ["0005d304-e60e-43c3-3507-ac1f6b60292f"]
  }
}

