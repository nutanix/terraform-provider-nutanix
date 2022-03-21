//discovery of nodes
data "nutanix_foundation_discover_nodes" "nodes"{}

//Get all unconfigured node's ipv6 addresses
locals {
  node_serials = keys(var.nodes_info)
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
        if lookup(var.nodes_info, node.node_serial, null) != null
        
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

    // Set optional fields if given
    xs_master_label = var.xs_master_label
    layout_egg_uuid = var.layout_egg_uuid
    hyperv_external_vnic = var.hyperv_external_vnic
    xen_config_type = var.xen_config_type
    ucsm_ip = var.ucsm_ip
    ucsm_password = var.ucsm_password
    xs_master_password = var.xs_master_password
    xs_master_ip = var.xs_master_ip
    hyperv_external_vswitch = var.hyperv_external_vswitch
    hypervisor_nameserver = var.hypervisor_nameserver
    hyperv_product_key = var.hyperv_product_key
    unc_username = var.unc_username
    install_script = var.install_script
    hypervisor_password = var.hypervisor_password
    unc_password = var.unc_password
    xs_master_username = var.xs_master_username
    skip_hypervisor = var.skip_hypervisor
    ucsm_user = var.ucsm_user
    svm_rescue_args = var.svm_rescue_args

    // foundation central settings 
    dynamic "fc_settings" {
        for_each = var.fc_settings != null ? [var.fc_settings] : []
        content {
            fc_metadata {
                fc_ip = fc_settings.value.fc_metadata.fc_ip
                api_key = fc_settings.value.fc_metadata.api_key
            }
            foundation_central = fc_settings.value.foundation_central
        }
    }

    // eos metadata
    dynamic "eos_metadata" {
        for_each = var.eos_metadata != null ? [var.eos_metadata] : []
        content {
            config_id = eos_metadata.value.config_id
            account_name = eos_metadata.value.account_name
            email = eos_metadata.value.email
        }
    }
    
    // node check tests
    dynamic "tests" {
        for_each = var.tests != null ? [var.tests] : []
        content {
            run_syscheck = eos_metadata.value.run_syscheck
            run_ncc = eos_metadata.value.run_ncc
        }
    }

    //Dynamically mention hypervisor information for every type
    dynamic "hypervisor_iso" {
        for_each = var.hypervisor_isos!=null ? [var.hypervisor_isos] : []
        content { 
            dynamic "kvm"{
                for_each = lookup(hypervisor_iso.value,"kvm",null) != null ? [hypervisor_iso.value.kvm] : []
                content {
                    filename = kvm.value.filename
                    checksum = kvm.value.checksum
                }
            }
            dynamic "esx"{
                for_each = lookup(hypervisor_iso.value,"esx",null) != null ? [hypervisor_iso.value.esx] : []
                content {
                    filename = esx.value.filename
                    checksum = kvm.value.checksum
                }
            }
            dynamic "hyperv"{
                for_each = lookup(hypervisor_iso.value,"hyperv",null) != null ? [hypervisor_iso.value.hyperv] : []
                content {
                    filename = hyperv.value.filename
                    checksum = kvm.value.checksum
                }
            }
            dynamic "xen"{
                for_each = lookup(hypervisor_iso.value,"xen",null) != null ? [hypervisor_iso.value.xen] : []
                content {
                    filename = xen.value.filename
                    checksum = kvm.value.checksum
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

                    // set required fields
                    hypervisor_hostname = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"hypervisor_hostname","")!=""? var.nodes_info[nodes.value.network_details.node_serial].hypervisor_hostname : nodes.value.network_details.hypervisor_hostname) : nodes.value.network_details.hypervisor_hostname
                    hypervisor_ip = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"hypervisor_ip","")!=""? var.nodes_info[nodes.value.network_details.node_serial].hypervisor_ip : nodes.value.network_details.hypervisor_ip) : nodes.value.network_details.hypervisor_ip
                    hypervisor = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"hypervisor","")!=""? var.nodes_info[nodes.value.network_details.node_serial].hypervisor : nodes.value.node.hypervisor) : nodes.value.node.hypervisor
                    image_now = true
                    ipmi_ip = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"ipmi_ip","")!=""? var.nodes_info[nodes.value.network_details.node_serial].ipmi_ip : (var.hypervisor != "" ? var.hypervisor : nodes.value.network_details.ipmi_ip)) : (var.hypervisor != "" ? var.hypervisor : nodes.value.network_details.ipmi_ip)
                    cvm_ip = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"cvm_ip","")!=""? var.nodes_info[nodes.value.network_details.node_serial].cvm_ip : nodes.value.network_details.cvm_ip) : nodes.value.network_details.cvm_ip
                    node_position = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"node_position","")!=""? var.nodes_info[nodes.value.network_details.node_serial].node_position : nodes.value.node.node_position) : nodes.value.node.node_position
                    ipmi_user = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"ipmi_user","")!=""? var.nodes_info[nodes.value.network_details.node_serial].ipmi_user : var.ipmi_user) : var.ipmi_user
                    ipmi_password = lookup(var.nodes_info,nodes.value.network_details.node_serial,null) != null ? (lookup(var.nodes_info[nodes.value.network_details.node_serial],"ipmi_password","")!=""? var.nodes_info[nodes.value.network_details.node_serial].ipmi_password : var.ipmi_password) : var.ipmi_password
                
                    node_serial = nodes.value.network_details.node_serial

                    // set optional fields
                    ipv6_address = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"ipv6_address",null) : null                    
                    image_delay = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"image_delay",null) : null
                    cvm_gb_ram = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"cvm_gb_ram",null) : null
                    device_hint = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"device_hint",null) : null
                    bond_mode = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"bond_mode",null) : null
                    rdma_passthrough = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"rdma_passthrough",null) : null
                    cluster_id = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"cluster_id",null) : null
                    ucsm_node_serial = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"ucsm_node_serial",null) : null
                    ipmi_configure_now = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"ipmi_configure_now",null) : null
                    cvm_num_vcpus = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"cvm_num_vcpus",null) : null
                    image_successful = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"image_successful",null) : null
                    ipv6_interface = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"ipv6_interface",null) : null
                    ipmi_mac = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"ipmi_mac",null) : null
                    rdma_mac_addr = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"rdma_mac_addr",null) : null
                    bond_uplinks = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"bond_uplinks",null) : null
                    current_network_interface = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"current_network_interface",null) : null
                    bond_lacp_rate = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"bond_lacp_rate",null) : null
                    ucsm_managed_mode = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"ucsm_managed_mode",null) : null
                    current_cvm_vlan_tag = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"current_cvm_vlan_tag",null) : null
                    exlude_boot_serial = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"exlude_boot_serial",null) : null
                    mitigate_low_boot_space = lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"mitigate_low_boot_space",null) : null

                    // set vswitches if given
                    dynamic "vswitches" {
                        for_each = (lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"vswitches",null) : null) != null ? [var.nodes_info.nodes.value.network_details.node_serial.vswitches] : []
                        content {
                            lacp = vswitches.value.lacp
                            bond_mode = vswitches.value.bond_mode
                            name = vswitches.value.name
                            uplinks = vswitches.value.uplinks
                            other_config = vswitches.value.other_config
                            mtu = vswitches.value.mtu
                        }
                    } 

                    // set ucsm_params of given
                    dynamic "ucsm_params" {
                        for_each = (lookup(var.nodes_info, nodes.nodes.value.network_details.node_serial,null) != null ? lookup(var.nodes_info[nodes.value.network_details.node_serial],"ucsm_params",null) : null) != null ? [var.nodes_info.nodes.value.network_details.node_serial.ucsm_params] : []
                        content {
                            native_vlan = ucsm_params.value.native_vlan
                            keep_ucsm_settings = ucsm_params.value.keep_ucsm_settings
                            mac_pool = ucsm_params.value.mac_pool
                            vlan_name = ucsm_params.value.vlan_name
                        }
                    } 
                    

                }
            }
        }
    }

    //Dynamically define cluster blocks
	  dynamic "clusters" {
        for_each = var.clusters != null ? var.clusters : []
        content {   
            // set required fields
            cluster_init_successful = true
            cluster_external_ip = clusters.value.cluster_external_ip
            redundancy_factor = clusters.value.redundancy_factor
            cluster_name = clusters.value.cluster_name
            cluster_members = clusters.value.cluster_members
            cluster_init_now = true

            // set optional fields
            enable_ns = lookup(clusters.value, "enable_ns", null)
            backplane_subnet = lookup(clusters.value, "backplane_subnet", null)
            backplane_netmask = lookup(clusters.value, "backplane_netmask", null)
            backplane_vlan = lookup(clusters.value, "backplane_vlan", null)
            cvm_ntp_servers = lookup(clusters.value, "cvm_ntp_servers", null)
            single_node_cluster = lookup(clusters.value, "single_node_cluster", null)
            cvm_dns_servers = lookup(clusters.value, "cvm_dns_servers", null)
            hypervisor_ntp_servers = lookup(clusters.value, "hypervisor_ntp_servers", null)

        }
    }
}
