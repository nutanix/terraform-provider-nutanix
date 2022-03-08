// [Required] : This list will be used to get node's information and image them
variable "node_serials" {
    description = "list of node serial numbers which wanted to be imaged"
    type = list(string)
}

// [Required] : Hypervisor netmask common for all nodes
variable "hypervisor_netmask" {
    description = "hypervisor netmask ip"
    type = string
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

// [Optional] : hypverisor type if given will override all nodes except nodes not having hypervisor defined in node_info_override
// The preference order is hypervisor(var) > hypervisor(node_info_override) > hypervisor(node's existing)
variable "hypervisor" {
    description = "default hypervisor type"
    type = string
    default = ""
}

// [Required] : nos_package file name for ex. nos_image.tar
variable "nos_package" {
    description = "nos package file name"
    type = string
}

// [Required] : default ipmi_user for all nodes. ipmi_user mentioned in node_info_override will override this for particular node
variable "ipmi_user" {
    description = "default ipmi username"
    type = string
}

// [Required] : default ipmi_password for all nodes. ipmi_password mentioned in node_info_override will override this for particular node
variable "ipmi_password"{
    description = "default ipmi password"
    type = string
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
        ] (list of strings)
    }
]
*/
variable "clusters" {
    description = "a list of map having info about cluster"
    type = any
    default = null
}

/*
[Optional] : node_info_override would have details for particular node that needs to be override
Format (Just mention things that needs override over default or existing values. Skip fields which doesn't need to be overriden over default or existing value) :
node_info_override = {
    <node1_serial_number> : {
        cvm_ip              : "10.xx.xx.xx"
        hypervisor          : "kvm"
        hypervisor_hostname : "batman-100"
        hypervisor_ip       : "10.xx.xx.xx"
        ipmi_ip             : "10.xx.xx.xx"
        ipmi_password       : "<password>"
        ipmi_user           : "<username>"
        node_position       : "A"
    },
    <node2_serial_number> : {
        cvm_ip              : "10.xx.xx.xx"
        hypervisor          : "kvm"
        hypervisor_hostname : "batman-100"
        hypervisor_ip       : "10.xx.xx.xx"
        ipmi_ip             : "10.xx.xx.xx"
        ipmi_password       : "<password>"
        ipmi_user           : "<username>"
        node_position       : "A"
    },
}
*/
variable "node_info_override" {
    description = "a map of node serial (key) to the info (value) for specific node related info"
    type = any
    default = {}
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
    type = any
    default = null
}

