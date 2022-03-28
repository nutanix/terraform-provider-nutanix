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

//Get node network details as per the ipv6 addresses
data "nutanix_foundation_node_network_details" "ntw_details" {
  ipv6_addresses = local.ipv6_addresses
}

//create map of ipv6_address => node_networ_details of each node
locals {
    ipv6_node_network_details_map = tomap({
        for node in data.nutanix_foundation_node_network_details.ntw_details.nodes:
        "${node.ipv6_address}" => node
        if node.ipv6_address != ""
    })
}
