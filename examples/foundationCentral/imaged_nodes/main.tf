// resources/datasources used in this file were introduced in nutanix/nutanix version 1.5.0-beta.2
terraform {
    required_providers {
      nutanix = {
          source = "nutanix/nutanix"
          version = ">1.5.0-beta.2"
      }
    }
}

provider "nutanix" {
    username  = "user"
    password  = "pass"
    endpoint  = "10.x.xx.xx"
    insecure  = true
    port      = 9440
}

// datasource to List all the nodes registered with Foundation Central.
data "nutanix_foundation_central_imaged_nodes_list" "img"{}

output "img1"{
    value = data.nutanix_foundation_central_imaged_nodes_list.img
}

// datasource to Get the details of a single node given its UUID.
data "nutanix_foundation_central_imaged_node_details" "imgdet"{
    imaged_node_uuid = "<imaged_node_uuid>"

output "imgdetails"{
    value = data.nutanix_foundation_central_imaged_node_details.imgdet
}