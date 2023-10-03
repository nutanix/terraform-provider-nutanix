// resources/datasources used in this file were introduced in nutanix/nutanix version 1.5.0-beta
terraform{
    required_providers{
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.5.0-beta"
        }
    }
}

// default foundation_port is 8000 so can be ignored
provider "nutanix" {
    // foundation_port = 8000
    foundation_endpoint = "10.xx.xx.xx"
}

/*
Description:
- Here we will discover nodes within ipv6 network of foundation vm & retrieve
  node network details all nodes which are not part of cluster.
- Nodes discovered having configured parameter false are not part of any cluster
*/

//discovery of nodes
data "nutanix_foundation_discover_nodes" "nodes"{}

//Get all unconfigured node's ipv6 addresses
locals {
  ipv6_addresses = flatten([
      for block in data.nutanix_foundation_discover_nodes.nodes.entities:
        [
          for node in block.nodes: 
            node.ipv6_address if node.configured==false
        ]
  ])
}

//Get node network details as per the ipv6 addresses collected
data "nutanix_foundation_node_network_details" "ntw_details" {
  ipv6_addresses = local.ipv6_addresses
}

//create map of node_serial => node_networ_details of each node
locals {
    ipv6_node_network_details_map = tomap({
        for node in data.nutanix_foundation_node_network_details.ntw_details.nodes:
        "${node.node_serial}" => node
        if node.node_serial != ""
    })
}

output "nodes" {
    value = local.ipv6_node_network_details_map
}
