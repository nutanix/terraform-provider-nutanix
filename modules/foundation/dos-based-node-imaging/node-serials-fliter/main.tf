module "discovered_nodes_network_details" {
    source = "../../discover-nodes-network-details/"
}

locals {
    //create list of node serials
    node_serials = keys(var.nodes_info)

    //Merge the network details and discovery of nodes response for given node serials
    block_details = [
        for block in module.discovered_nodes_network_details.discovered_node_details:
            {
                "nodes" = [
                    for node in block.nodes:
                        {
                            "node_serial" = tomap({
                                "val" = node.node_serial != "" ? node.node_serial : (lookup(module.discovered_nodes_network_details.ipv6_to_node_network_details_map, node.ipv6_address, null) != null ? module.discovered_nodes_network_details.ipv6_to_node_network_details_map[node.ipv6_address].node_serial : "")
                            })
                            "node" = node
                            "network_details" = lookup(module.discovered_nodes_network_details.ipv6_to_node_network_details_map,lookup(node,"ipv6_address",""), null)
                        } if lookup(module.discovered_nodes_network_details.ipv6_to_node_network_details_map, node.ipv6_address, null) != null ? ( contains(local.node_serials, node.node_serial) || contains(local.node_serials,module.discovered_nodes_network_details.ipv6_to_node_network_details_map[node.ipv6_address].node_serial)) : false
                ]
                "block_id" = lookup(block, "block_id", "")
            }
    ]

    //Remove not required data
    filtered_node_details = [
        for block in local.block_details:
            block if length(block.nodes)>0
    ]

    // list of required details to check if this details are present or not for node imaging
    // source can be "network_details" (node network details info) or "node" (normal node info from discover nodes)
    // global defines the params which can be declared common for all nodes in module input, ex. ipmi_user, etc. 
    // not_allowed are used to validate to info == not_allowed. This is to avoid "", 0, false, etc.
    required_details = [
        {
            attribute: "ipmi_ip",
            source: "network_details",
            global: false
            not_allowed : ""
        },
        {
            attribute: "cvm_ip",
            source: "network_details",
            global: false
            not_allowed : ""
        },
        {
            attribute: "node_position",
            source: "node",
            global: false
            not_allowed : ""
        },
        {
            attribute: "hypervisor",
            source: "node",
            global: true
            not_allowed : ""
        },{
            attribute: "hypervisor",
            source: "node",
            global: true
            not_allowed : "pheonix"
        },
        {
            attribute: "ipmi_user",
            source: "",
            global: true
            not_allowed : ""
        },
        {
            attribute: "ipmi_password",
            source: "",
            global: true
            not_allowed : ""
        }
    ]
        
    // create error messages incase required details are not present in node_info/discover_nodes/node_network_details for a particular node
    node_info_validation_messages = flatten([
        for block in local.block_details:
            [
                for node in block.nodes: [
                    for attr_details in local.required_details:
                        format("%s for node serial %s is missing. ", attr_details.attribute, node.node_serial.val)
                        if try(var.nodes_info[node.node_serial.val][attr_details.attribute], null) == null && try(node[attr_details.source][attr_details.attribute] == attr_details.not_allowed, true)  && (attr_details.global? var.defaults[attr_details.attribute] == null : true)
                ]

            ]
    ])
}


// Internal assert helper checking for error messages from above operations and errors out if present
data "nutanix_assert_helper" "checks" {
    
    dynamic "checks" {
        for_each = local.node_info_validation_messages
        content{
            condition = false
            error_message = checks.value
        }
    }
}

//Resource block for imaging the nodes
resource "nutanix_foundation_image_nodes" "this"{

    // Required fields to be taken from module input
    nos_package = var.nos_package
    ipmi_user = var.defaults.ipmi_user
    ipmi_password = var.defaults.ipmi_password
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
                for_each = hypervisor_iso.value.kvm != null ? [hypervisor_iso.value.kvm] : []
                content {
                    filename = kvm.value.filename
                    checksum = kvm.value.checksum
                }
            }
            dynamic "esx"{
                for_each = hypervisor_iso.value.esx != null ? [hypervisor_iso.value.esx] : []
                content {
                    filename = esx.value.filename
                    checksum = esx.value.checksum
                }
            }
            dynamic "hyperv"{
                for_each = hypervisor_iso.value.hyperv != null ? [hypervisor_iso.value.hyperv] : []
                content {
                    filename = hyperv.value.filename
                    checksum = hyperv.value.checksum
                }
            }
            dynamic "xen"{
                for_each = hypervisor_iso.value.xen != null ? [hypervisor_iso.value.xen] : []
                content {
                    filename = xen.value.filename
                    checksum = xen.value.checksum
                }
            }
        }
    }

    //Dynamically define blocks and its nodes as per the input variables of module
    dynamic "blocks"{
        for_each = local.filtered_node_details
        content {
            block_id = blocks.value.block_id
            dynamic "nodes" {
                for_each = blocks.value.nodes
                content { 

                    // set required fields
                    hypervisor_hostname = var.nodes_info[nodes.value.node_serial.val].hypervisor_hostname
                    hypervisor_ip = var.nodes_info[nodes.value.node_serial.val].hypervisor_ip
                    hypervisor = var.nodes_info[nodes.value.node_serial.val].hypervisor != null ? var.nodes_info[nodes.value.node_serial.val].hypervisor : var.defaults.hypervisor
                    image_now = true
                    ipmi_ip = var.nodes_info[nodes.value.node_serial.val].ipmi_ip != null ? var.nodes_info[nodes.value.node_serial.val].ipmi_ip : nodes.value.network_details.ipmi_ip
                    cvm_ip = var.nodes_info[nodes.value.node_serial.val].cvm_ip != null ? var.nodes_info[nodes.value.node_serial.val].cvm_ip : nodes.value.network_details.cvm_ip
                    node_position = var.nodes_info[nodes.value.node_serial.val].node_position != null ? var.nodes_info[nodes.value.node_serial.val].node_position : nodes.value.node.node_position
                    ipmi_user = var.nodes_info[nodes.value.node_serial.val].ipmi_user
                    ipmi_password = var.nodes_info[nodes.value.node_serial.val].ipmi_password
                    node_serial = nodes.value.node_serial.val

                    // set optional fields
                    ipv6_address = var.nodes_info[nodes.value.node_serial.val].ipv6_address != null ? var.nodes_info[nodes.value.node_serial.val].ipv6_address : nodes.value.node.ipv6_address                    
                    image_delay = var.nodes_info[nodes.value.node_serial.val].image_delay
                    cvm_gb_ram = var.nodes_info[nodes.value.node_serial.val].cvm_gb_ram != null ? var.nodes_info[nodes.value.node_serial.val].cvm_gb_ram : var.defaults.cvm_gb_ram
                    device_hint = var.nodes_info[nodes.value.node_serial.val].device_hint
                    bond_mode = var.nodes_info[nodes.value.node_serial.val].bond_mode
                    rdma_passthrough = var.nodes_info[nodes.value.node_serial.val].rdma_passthrough
                    cluster_id = var.nodes_info[nodes.value.node_serial.val].cluster_id
                    ucsm_node_serial = var.nodes_info[nodes.value.node_serial.val].ucsm_node_serial
                    ipmi_configure_now = var.nodes_info[nodes.value.node_serial.val].ipmi_configure_now
                    cvm_num_vcpus = var.nodes_info[nodes.value.node_serial.val].cvm_num_vcpus != null ? var.nodes_info[nodes.value.node_serial.val].cvm_num_vcpus : var.defaults.cvm_num_vcpus
                    image_successful = var.nodes_info[nodes.value.node_serial.val].image_successful
                    ipv6_interface = var.nodes_info[nodes.value.node_serial.val].ipv6_interface
                    ipmi_mac = var.nodes_info[nodes.value.node_serial.val].ipmi_mac
                    rdma_mac_addr = var.nodes_info[nodes.value.node_serial.val].rdma_mac_addr
                    bond_uplinks = var.nodes_info[nodes.value.node_serial.val].bond_uplinks
                    current_network_interface = var.nodes_info[nodes.value.node_serial.val].current_network_interface != null ? var.nodes_info[nodes.value.node_serial.val].current_network_interface : nodes.value.node.current_network_interface
                    bond_lacp_rate = var.nodes_info[nodes.value.node_serial.val].bond_lacp_rate
                    ucsm_managed_mode = var.nodes_info[nodes.value.node_serial.val].ucsm_managed_mode
                    current_cvm_vlan_tag = var.nodes_info[nodes.value.node_serial.val].current_cvm_vlan_tag != null ? var.nodes_info[nodes.value.node_serial.val].current_cvm_vlan_tag : (var.defaults.current_cvm_vlan_tag != null ? var.defaults.current_cvm_vlan_tag : nodes.value.node.current_cvm_vlan_tag)
                    exlude_boot_serial = var.nodes_info[nodes.value.node_serial.val].exlude_boot_serial
                    mitigate_low_boot_space = var.nodes_info[nodes.value.node_serial.val].mitigate_low_boot_space

                    // set vswitches if given
                    dynamic "vswitches" {
                        for_each = var.nodes_info[nodes.value.node_serial.val].vswitches != null ? [var.nodes_info[nodes.value.node_serial.val].vswitches] : []
                        content{
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
                        for_each = var.nodes_info[nodes.value.node_serial.val].ucsm_params != null ? [var.nodes_info[nodes.value.node_serial.val].ucsm_params] : []
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
        for_each = var.clusters
        content {   
            // set required fields
            cluster_init_successful = true
            cluster_external_ip = clusters.value.cluster_external_ip
            redundancy_factor = clusters.value.redundancy_factor
            cluster_name = clusters.value.cluster_name
            cluster_members = clusters.value.cluster_members
            cluster_init_now = true

            // set optional fields
            enable_ns = clusters.value.enable_ns
            backplane_subnet = clusters.value.backplane_subnet
            backplane_netmask = clusters.value.backplane_netmask
            backplane_vlan = clusters.value.backplane_vlan
            cvm_ntp_servers = clusters.value.cvm_ntp_servers
            single_node_cluster = clusters.value.single_node_cluster
            cvm_dns_servers = clusters.value.cvm_dns_servers
            hypervisor_ntp_servers = clusters.value.hypervisor_ntp_servers

        }
    }
}
