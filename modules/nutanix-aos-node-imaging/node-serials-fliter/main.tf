//discovery of nodes
data "nutanix_foundation_discover_nodes" "nodes"{}

//Get all unconfigured node's ipv6 addresses
locals {
  ipv6Addresses = flatten([
      for block in data.nutanix_foundation_discover_nodes.nodes.entities:
        [
          for node in block.nodes: 
            node.ipv6_address if node.configured==false
        ]
  ])
}

//Get node network details as per the ipv6 addresses
data "nutanix_foundation_node_network_details" "ntwDetails" {
  ipv6_addresses = local.ipv6Addresses
}

locals {
    //Create map of node_serial => network details as per input node serials
    nodeSerialNtwDetails = tomap({
        for node in data.nutanix_foundation_node_network_details.ntwDetails.nodes:
        "${node.ipv6_address}" => node
        if contains(var.node_serials, node.node_serial) == true
        
    })

    //Merge the network details and discovery of nodes response for given node serials
    blockDetails = [
        for block in data.nutanix_foundation_discover_nodes.nodes.entities:
            {
                "nodes" = [
                    for node in block.nodes:
                        {
                            "node" = node
                            "network_details" = lookup(local.nodeSerialNtwDetails,lookup(node,"ipv6_address",""), null)
                        } if lookup(local.nodeSerialNtwDetails, lookup(node,"ipv6_address",""), null) != null
                ]
                "block_id" = lookup(block, "block_id", "")
            }
    ]

    //Remove not required data
    filteredNodeDetails = [
        for block in local.blockDetails:
            block if length(block.nodes)>0
    ]
}

//Resource block for imaging the nodes
resource "nutanix_foundation_image_nodes" "this"{

    // Required fields to be taken from module input
    nos_package = var.nos_package
    ipmi_user = var.ipmi_user
    ipmi_password = var.ipmi_password
    hypervisor_netmask = var.hypervisor_netmask
    hypervisor_gateway = var.hypervisor_gateway
    cvm_netmask = var.cvm_netmask
    cvm_gateway = var.hypervisor_gateway

    
    //Dynamically mention hypervisor information for every type
    dynamic "hypervisor_iso" {
        for_each = var.hypervisor_isos!=null ? [var.hypervisor_isos] : []
        content { 
            dynamic "kvm"{
                for_each = lookup(hypervisor_iso.value,"kvm",null) != null ? [hypervisor_iso.value.kvm] : []
                content {
                    filename = kvm.value.filename
                    checksum = lookup(kvm.value, "checksum", "")
                }
            }
            dynamic "esx"{
                for_each = lookup(hypervisor_iso.value,"esx",null) != null ? [hypervisor_iso.value.esx] : []
                content {
                    filename = esx.value.filename
                    checksum = lookup(esx.value, "checksum", "")

                }
            }
            dynamic "hyperv"{
                for_each = lookup(hypervisor_iso.value,"hyperv",null) != null ? [hypervisor_iso.value.hyperv] : []
                content {
                    filename = hyperv.value.filename
                    checksum = lookup(hyperv.value, "checksum", "")

                }
            }
            dynamic "xen"{
                for_each = lookup(hypervisor_iso.value,"xen",null) != null ? [hypervisor_iso.value.xen] : []
                content {
                    filename = xen.value.filename
                    checksum = lookup(xen.value, "checksum", "")

                }
            }
        }
    }

    //Dynamically define blocks and its nodes as per the input variables of module
    dynamic "blocks"{
        for_each = local.filteredNodeDetails
        content {
            block_id = blocks.value.block_id
            dynamic "nodes" {
                for_each = blocks.value.nodes
                content { 
                    hypervisor_hostname = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"hypervisor_hostname","")!=""? var.node_info_override[nodes.value.network_details.node_serial].hypervisor_hostname : nodes.value.network_details.hypervisor_hostname) : nodes.value.network_details.hypervisor_hostname
                    hypervisor_ip = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"hypervisor_ip","")!=""? var.node_info_override[nodes.value.network_details.node_serial].hypervisor_ip : nodes.value.network_details.hypervisor_ip) : nodes.value.network_details.hypervisor_ip
                    hypervisor = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"hypervisor","")!=""? var.node_info_override[nodes.value.network_details.node_serial].hypervisor : nodes.value.node.hypervisor) : nodes.value.node.hypervisor
                    image_now = true
                    ipmi_ip = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"ipmi_ip","")!=""? var.node_info_override[nodes.value.network_details.node_serial].ipmi_ip : (var.hypervisor != "" ? var.hypervisor : nodes.value.network_details.ipmi_ip)) : (var.hypervisor != "" ? var.hypervisor : nodes.value.network_details.ipmi_ip)
                    cvm_ip = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"cvm_ip","")!=""? var.node_info_override[nodes.value.network_details.node_serial].cvm_ip : nodes.value.network_details.cvm_ip) : nodes.value.network_details.cvm_ip
                    node_position = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"node_position","")!=""? var.node_info_override[nodes.value.network_details.node_serial].node_position : nodes.value.node.node_position) : nodes.value.node.node_position
                    ipmi_user = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"ipmi_user","")!=""? var.node_info_override[nodes.value.network_details.node_serial].ipmi_user : var.ipmi_user) : var.ipmi_user
                    ipmi_password = lookup(var.node_info_override,nodes.value.network_details.node_serial,null) != null ? (lookup(var.node_info_override[nodes.value.network_details.node_serial],"ipmi_password","")!=""? var.node_info_override[nodes.value.network_details.node_serial].ipmi_password : var.ipmi_password) : var.ipmi_password
                }
            }
        }
    }

    //Dynamically define cluster blocks
	  dynamic "clusters" {
        for_each = var.clusters != null ? var.clusters : []
        content {   
            cluster_init_successful = true
            cluster_external_ip = clusters.value.cluster_external_ip
            redundancy_factor = clusters.value.cluster_external_ip.redundancy_factor
            cluster_name = clusters.value.cluster_external_ip.cluster_name
            cluster_members = clusters.value.cluster_external_ip.cluster_members
            cluster_init_now = true
        }
    }
}
