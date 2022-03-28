locals {

    // list of required details to check if this details are present or not for node imaging
    required_details = [ "hypervisor", "ipmi_user", "ipmi_password", "ipmi_netmask", "ipmi_gateway"]

    // create error messages incase required details are not present in node_info/discover_nodes/node_network_details for a particular node
    node_info_validation_messages = flatten([
        for block in var.blocks:
            [
                for node in block.nodes: [
                    for attr in local.required_details:
                        format("%s for node having cvm_ip %s is missing. ", attr, node.cvm_ip)
                        if node[attr] == null && var.defaults[attr] == null
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
        for_each = var.blocks
        content {
            block_id = blocks.value.block_id
            dynamic "nodes" {
                for_each = blocks.value.nodes
                content { 

                    // set required fields
                    ipmi_netmask = nodes.value.ipmi_netmask != null ? nodes.value.ipmi_netmask : var.defaults.ipmi_netmask
                    ipmi_gateway = nodes.value.ipmi_gateway != null ? nodes.value.ipmi_gateway : var.defaults.ipmi_gateway
                    hypervisor_hostname = nodes.value.hypervisor_hostname
                    hypervisor_ip = nodes.value.hypervisor_ip
                    hypervisor = nodes.value.hypervisor != null ? nodes.value.hypervisor : var.defaults.hypervisor
                    image_now = true
                    ipmi_ip = nodes.value.ipmi_ip
                    cvm_ip = nodes.value.cvm_ip
                    node_position = nodes.value.node_position 
                    ipmi_user = nodes.value.ipmi_user
                    ipmi_password = nodes.value.ipmi_password
                    node_serial = nodes.value.node_serial

                    // set optional fields
                    ipv6_address = nodes.value.ipv6_address                    
                    image_delay = nodes.value.image_delay    
                    cvm_gb_ram = nodes.value.cvm_gb_ram != null ? nodes.value.cvm_gb_ram : var.defaults.cvm_gb_ram
                    device_hint = nodes.value.device_hint    
                    bond_mode = nodes.value.bond_mode    
                    rdma_passthrough = nodes.value.rdma_passthrough
                    cluster_id = nodes.value.cluster_id    
                    ucsm_node_serial = nodes.value.ucsm_node_serial    
                    ipmi_configure_now = nodes.value.ipmi_configure_now    
                    cvm_num_vcpus = nodes.value.cvm_num_vcpus != null ? nodes.value.cvm_num_vcpus : var.defaults.cvm_num_vcpus
                    image_successful = nodes.value.image_successful    
                    ipv6_interface = nodes.value.ipv6_interface    
                    ipmi_mac = nodes.value.ipmi_mac
                    rdma_mac_addr = nodes.value.rdma_mac_addr    
                    bond_uplinks = nodes.value.bond_uplinks    
                    current_network_interface = nodes.value.current_network_interface    
                    bond_lacp_rate = nodes.value.bond_lacp_rate
                    ucsm_managed_mode = nodes.value.ucsm_managed_mode
                    current_cvm_vlan_tag = nodes.value.current_cvm_vlan_tag != null ? nodes.value.current_cvm_vlan_tag : var.defaults.current_cvm_vlan_tag
                    exlude_boot_serial = nodes.value.exlude_boot_serial    
                    mitigate_low_boot_space = nodes.value.mitigate_low_boot_space

                    // set vswitches if given
                    dynamic "vswitches" {
                        for_each = nodes.value.vswitches != null ? [nodes.value.vswitches] : []
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
                        for_each = nodes.value.ucsm_params != null ? [nodes.value.ucsm_params] : []
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
