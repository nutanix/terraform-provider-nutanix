terraform {
    experiments = [module_variable_optional_attrs]
}

// [Required] : Hypervisor netmask common for all nodes
variable "hypervisor_netmask" {
    description = "hypervisor netmask ip"
    type = string
}

// [Optional] : custom timeout in minutes form image_nodes resource
variable "timeout" {
  description = "custom timeout in minutes for image_nodes resource"
  default = null
  type = number
}

/*
[Optional] : this will be used if there is no info provided in node spec of blocks
*/
variable "defaults" {
    description = "default spec for nodes"
    type = object({
        ipmi_netmask = optional(string)
        ipmi_gateway = optional(string)
        ipmi_user = optional(string)
        ipmi_password = optional(string)
        hypervisor = optional(string)
        cvm_gb_ram = optional(string)
        cvm_num_vcpus = optional(string)
        current_cvm_vlan_tag = optional(string)
    })
    default = {
      ipmi_netmask : null
      ipmi_gateway : null
      current_cvm_vlan_tag : null
      cvm_gb_ram : null
      cvm_num_vcpus : null
      hypervisor : null
      ipmi_password : null
      ipmi_user : null
    }
}

// [Required] : Hypervisor gateway common for all nodes
variable "hypervisor_gateway" {
    description = "hypervisor gateway ip"
    type = string
}

// [Required] : cvm_netmask common for all nodes
variable "cvm_netmask" {
    description = "cvm netmask ip"
    type = string
}

// [Required] : cvm_gateway common for all nodes
variable "cvm_gateway" {
    description = "cvm gateway ip"
    type = string
}

// [Required] : nos_package file name for ex. nos_image.tar
variable "nos_package" {
    description = "nos package file name"
    type = string
}

// [Optional] : xs_master_label
variable "xs_master_label" {
    description = "xs_master_label for node imaging"
    type = string
    default = ""
}

// [Optional] : layout_egg_uuid
variable "layout_egg_uuid" {
    description = "layout_egg_uuid for node imaging"
    type = string
    default = ""
}

// [Optional] : hyperv_external_vnic
variable "hyperv_external_vnic" {
    description = "hyperv_external_vnic for node imaging"
    type = string
    default = ""
}

// [Optional] : xen_config_type
variable "xen_config_type" {
    description = "xen_config_type for node imaging"
    type = string
    default = ""
}

// [Optional] : ucsm_ip
variable "ucsm_ip" {
    description = "ucsm_ip for node imaging"
    type = string
    default = ""
}

// [Optional] : ucsm_password
variable "ucsm_password" {
    description = "ucsm_password for node imaging"
    type = string
    default = ""
}

// [Optional] : xs_master_password
variable "xs_master_password" {
    description = "xs_master_password for node imaging"
    type = string
    default = ""
}

// [Optional] : xs_master_ip
variable "xs_master_ip" {
    description = "xs_master_ip for node imaging"
    type = string
    default = ""
}

// [Optional] : hyperv_external_vswitch
variable "hyperv_external_vswitch" {
    description = "hyperv_external_vswitch for node imaging"
    type = string
    default = ""
}

// [Optional] : hypervisor_nameserver
variable "hypervisor_nameserver" {
    description = "hypervisor_nameserver for node imaging"
    type = string
    default = ""
}

// [Optional] : hyperv_product_key
variable "hyperv_product_key" {
    description = "hyperv_product_key for node imaging"
    type = string
    default = ""
}

// [Optional] : unc_username
variable "unc_username" {
    description = "unc_username for node imaging"
    type = string
    default = ""
}

// [Optional] : install_script
variable "install_script" {
    description = "install_script for node imaging"
    type = string
    default = ""
}

// [Optional] : hypervisor_password
variable "hypervisor_password" {
    description = "hypervisor_password for node imaging"
    type = string
    default = ""
}

// [Optional] : unc_password
variable "unc_password" {
    description = "unc_password for node imaging"
    type = string
    default = ""
}

// [Optional] : xs_master_username
variable "xs_master_username" {
    description = "xs_master_username for node imaging"
    type = string
    default = ""
}

// [Optional] : skip_hypervisor
variable "skip_hypervisor" {
    description = "skip_hypervisor for node imaging"
    type = bool
    default = null
}

// [Optional] : ucsm_user
variable "ucsm_user" {
    description = "ucsm_user for node imaging"
    type = string
    default = ""
}

// [Optional] : foundation central settings
variable "fc_settings" {
    description = "foundation central settings for node imaging"
    type = object({
        fc_metadata = object({
            fc_ip = string
            api_key = string
        })
        foundation_central = bool
    })
    default = null
}

// [Optional] : svm_rescue_args
variable "svm_rescue_args" {
    description = "svm_rescue_args for node imaging"
    type = list(string)
    default = null
}

// [Optional] : eos_metadata
variable "eos_metadata" {
    description = "eos_metadata for node imaging"
    type = object({
        config_id = string
        account_name = string
        email = string
    })
    default = null
}

// [Optional] : tests
variable "tests" {
    description = "tests params for node imaging"
    type = object({
        run_syscheck = bool
        run_ncc = bool
    })
    default = null
}

/*
[Optional] : to create cluster out of imaged nodes
Format (this are required for cluster creation):
clusters = [
    {
        cluster_external_ip = "10.xx.xx.xx" (string)
        redundancy_factor = xx (number)
        cluster_name = "cluster-1" (string)
        cluster_members = [
            "10.xx.xx.xx", "10.xx.xx.xx"
        ]
    }
]
*/
variable "clusters" {
    description = "a list of map having info about cluster"
    type = list(object({
        enable_ns = optional(bool)
        backplane_subnet = optional(string)
        backplane_netmask = optional(string)
        redundancy_factor = number
        backplane_vlan = optional(string)
        cluster_name = string
        cluster_external_ip = optional(string)
        cvm_ntp_servers = optional(string)
        single_node_cluster = optional(bool)
        cluster_members = list(string)
        cvm_dns_servers = optional(string)
        hypervisor_ntp_servers = optional(string)

    }))
    default = []
}

/*
[Required] : variable to declare block of nodes
Fields like node_position, ipmi_ip, cvm_ip, hypervisor_ip & hypervisor_hostname are mandatory
hypervisor needs to be mentioned in defaults or in every node spec here
*/
variable "blocks" {
    description = "a map of node serial (key) to the info (value) for specific node related info"
    type = list(object({
        block_id = optional(string)
        nodes = list(object(
            {
                ipmi_netmask = optional(string)
                ipmi_gateway = optional(string)
                ipv6_address = optional(string)
                node_position = string
                image_delay = optional(number)
                ucsm_params = optional(object({
                    native_vlan = number
                    keep_ucsm_settings = number
                    mac_pool = string
                    vlan_name = string

                }))
                hypervisor_hostname = string
                cvm_gb_ram = optional(number)
                device_hint = optional(string)
                bond_mode = optional(string)
                rdma_passthrough = optional(bool)
                cluster_id = optional(string)
                ucsm_node_serial = optional(string)
                hypervisor_ip = string
                node_serial = optional(string)
                ipmi_configure_now = optional(bool)
                cvm_num_vcpus = optional(number)
                image_successful = optional(bool)
                ipv6_interface = optional(string)
                ipmi_mac = optional(string)
                rdma_mac_addr = optional(string)
                bond_uplinks = optional(list(string))
                current_network_interface = optional(string)
                hypervisor = optional(string)
                vswitches = optional(object({
                    lacp = string
                    bond_mode = string
                    name = string
                    uplinks = list(string)
                    other_config = list(string)
                    mtu = number
                }))
                bond_lacp_rate = optional(string)
                ucsm_managed_mode = optional(string)
                ipmi_ip = string
                current_cvm_vlan_tag = optional(number)
                cvm_ip = string
                exlude_boot_serial = optional(string)
                mitigate_low_boot_space = optional(bool)
                ipmi_user = optional(string)
                ipmi_password = optional(string)
            }
        ))
    })) 
}

/*
[Optional] : It is only optional when nos package bundled hypervisor is need to be used
             and hypervisor = "kvm" needs to be mentioned. If wants to use other hypervisor,
             details needs to be mentioned using below format.
Format (skip types which doesn't needs to be used) :
hypervisor_isos = {
    kvm : {
        filename : "xyz.iso" (required)
        checksum : "xyz" (optional)
    },
    esx : {
        filename : "xyz.iso" (required)
        checksum : "xyz" (optional)
    },
    hyperv : {
        filename : "xyz.iso" (required)
        checksum : "xyz" (optional)
    },
    xen : {
        filename : "xyz.iso" (required)
        checksum : "xyz" (optional)
    },
}
*/
variable "hypervisor_isos" {
    description = "a map of hypervisor type to file name"
    type = object({
        kvm = optional(object({
            filename = string
            checksum = string
        }))
        esx = optional(object({
            filename = string
            checksum = string
        }))
        hyperv = optional(object({
            filename = string
            checksum = string
        }))
        xen = optional(object({
            filename = string
            checksum = string
        }))
    })
    default = null
}
