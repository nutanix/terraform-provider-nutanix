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

// pull hypervisors and nos packages info
data "nutanix_foundation_hypervisor_isos" "isos"{}
data "nutanix_foundation_nos_packages" "nos"{}

/*
Description:
Incase of baremetal nodes or manually added nodes, there is no information obtained using 
data sources. So here we have to provide each and every info required to image nodes.

Here we can define defaults which will be added to all nodes spec while imaging 
*/
module "batch1" {
    source = "../../../modules/foundation/manual-node-imaging/"
    timeout = 120
    cvm_netmask = "xx.xx.xx.xx"
    cvm_gateway = "xx.xx.xx.xx"
    hypervisor_gateway = "xx.xx.xx.xx"
    hypervisor_netmask = "xx.xx.xx.xx"

    // this defaults will be added to every node spec have less priority then info given inside node spec
    // Allowed params : ipmi_user, ipmi_password, hypervisor, cvm_gb_ram and current_cvm_vlan_tag
    defaults = {
        ipmi_netmask : "xx.xx.xx.xx"
        ipmi_gateway : "xx.xx.xx.xx"
        ipmi_user : "<default-ipmi-username>"
        cvm_gb_ram : 50
        hypervisor : "kvm"
    }

    nos_package = data.nutanix_foundation_nos_packages.nos.entities[0]

    // check variables.tf for all available options
    blocks = [
        {
            nodes : [
                {
                    node_position: "D",
                    hypervisor_hostname: "superman-4",
                    hypervisor_ip: "xx.xx.xx.xx",
                    ipmi_ip: "xx.xx.xx.xx",
                    cvm_ip: "xx.xx.xx.xx",
                    ipmi_password: "<ipmi-password-1>"
                },
                {
                    node_position: "C",
                    hypervisor_hostname: "superman-3",
                    hypervisor_ip: "xx.xx.xx.xx",
                    ipmi_ip: "xx.xx.xx.xx",
                    cvm_ip: "xx.xx.xx.xx",
                    ipmi_password: "<ipmi-password-2>"
                },
                {
                    node_position: "B",
                    hypervisor_hostname: "xx.xx.xx.xx",
                    hypervisor_ip: "xx.xx.xx.xx",
                    ipmi_ip: "xx.xx.xx.xx",
                    cvm_ip: "xx.xx.xx.xx",
                    ipmi_password: "<ipmi-password-2>"
                }
            ]
        }
    ]

    // create clusters
    clusters = [
        {
            cluster_external_ip : "xx.xx.xx.xx"
            redundancy_factor : 2
            cluster_name : "test_cluster"
            cluster_members : [
                "xx.xx.xx.xx", "xx.xx.xx.xx", "xx.xx.xx.xx"
            ]
        }
    ]
}

output "check" {
    value = module.batch1
}