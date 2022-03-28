output "discovered_node_details" {
    value = data.nutanix_foundation_discover_nodes.nodes.entities
}

output "node_network_details_unconfigured_nodes" {
    value = data.nutanix_foundation_node_network_details.ntw_details
}

output "ipv6_to_node_network_details_map" {
    value = local.ipv6_node_network_details_map
}