// resources/datasources used in this file were introduced in nutanix/nutanix version >1.4.1
terraform{
    required_providers{
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.4.2"
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
Description :
Here we will image nutanix imaged nodes having discovery os running on them and create cluster out of it using ipmi based imaging.
We have used node serial based module. So using given node serials, module will pull info using
node discovery and get its network details internally. We can mention details in every node spec which we 
want to ovveride from network details info we obtained. Then collectively it will use them to image nodes

For ex. here the node is discoverable and we get cvm_ip, ipmi_ip, etc. internaly using data sources,
but for some nodes we can give this details and module will override it, else it will use info coming from 
network detail data sources.
nodes_info = {
    "<node-serial>" : {
        ipmi_password = "xxxx"
        cvm_ip = "xx.xx.xx.xx" (required)
        hypervisor_ip = "xx.xx.xx.xx" (required)
        hypervisor_hostname = "xyz" (required)
    }
}

Note : hypervisor realated attributes are must in either defaults or node spec
*/

// use dos-based-node-imaging based on given node-serials-filter
module "batch1" {

    // source where module code is present in local machine
    source = "../../../modules/foundation/aos-based-node-imaging/node-serials-fliter/"
    timeout = 120
    cvm_netmask = "xx.xx.xx.xx"
    cvm_gateway = "xx.xx.xx.xx"
    hypervisor_gateway = "xx.xx.xx.xx"
    hypervisor_netmask = "xx.xx.xx.xx"
    nos_package = data.nutanix_foundation_nos_packages.nos.entities[0]
    
    
    # hypervisor_isos = {
    #     kvm : {
    #         filename : "xyz.iso"
    #         checksum : "xyz"
    #     },
    #     esx : {
    #         filename : "xyz.iso"
    #         checksum : "xyz"
    #     },
    # }

    // this defaults will be added to every node spec have less priority then info given inside node spec
    // Allowed params : ipmi_user, ipmi_password, hypervisor, cvm_gb_ram and current_cvm_vlan_tag
    defaults = {
        ipmi_user : "<ipmi-username>"
        cvm_gb_ram : 50
        hypervisor : "kvm"
    }
   
    /* 
    Mention nodes_serial => info of nodes map which needs to be imaged.
    Module will use union of node network details, discover nodes details and user given fields 
    with priority to user given fields > node network details for image nodes input.
    Module internally gets all info like ipmi_ip, ipv6_address, etc. from discover_nodes and node_network_details data sources internaly. 
    */
    nodes_info = {
        
        "<node-serial-1>" : {
            ipmi_password : "<ipmi-password-1>"
            hypervisor_ip : "xx.xx.xx.xx"
            hypervisor_hostname : "xx.xx.xx.xx"
            cvm_ip : "xx.xx.xx.xx"
        }
         "<node-serial-2>" : {
            ipmi_password : "<ipmi-password-2>"
            hypervisor_ip : "xx.xx.xx.xx"
            hypervisor_hostname : "xx.xx.xx.xx"
            cvm_ip : "xx.xx.xx.xx"
        }
        "<node-serial-3>" : {
            ipmi_password : "<ipmi-password-3>"
            hypervisor_ip : "xx.xx.xx.xx"
            hypervisor_hostname : "xx.xx.xx.xx"
            cvm_ip : "xx.xx.xx.xx"
        }
    }

    // give cluster creation info
    clusters = [
        {
            redundancy_factor : 2
            cluster_name : "test_cluster"
            cluster_members : [
                "xx.xx.xx.xx", "xx.xx.xx.xx", "xx.xx.xx.xx"
            ]
        }
    ]
}

output "node-imaging" {
    value = module.batch1
}