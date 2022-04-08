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
Description :
Here we will image nutanix imaged nodes having cvm running on them and create cluster out of it using ipmi based imaging.
We have used node serial based module. So using given node serials, module will pull info using
node discovery and get its network details internally. We can mention details in every node spec which we 
want to ovveride from network details info we obtained. Then collectively it will use them to image nodes

For ex. here the node is discoverable and we get cvm_ip, hypervisor_ip, etc internaly using data sources,
but for some nodes we can give this details and module will override it, else it will use info coming from 
network detail data sources as imaging input.
nodes_info = {
    "<node-serial>" : {
        ipmi_password = "xxxx"
        cvm_ip = "xx.xx.xx.xx"
        hypervisor_ip = "xx.xx.xx.xx"
    }
}
*/

// use aos-based-node-imaging based on given node-serials-filter
module "batch1" {

    // source where module code is present in local machine
    source = "../../../modules/foundation/aos-based-node-imaging/node-serials-filter/"
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
        ipmi_user : "<ipmi-user-1>"
        cvm_gb_ram : 50
        hypervisor : "kvm"
    }
   
    /* 
    Mention nodes_serial => info of nodes map which needs to be imaged.
    Module will use union of info from node network details, discover nodes details and user given fields 
    with priority to user given fields > node network details for image nodes input.
    Below info get all info like cvm_ip, hypervisor_ip, etc. from discover_nodes and node_network_details data sources internaly. 
    Note : ipmi_user & ipmi_password needs to be defined either at defaults or in node spec (if every node have different creds)
    */
    nodes_info = {
        
        "<node-serial-1>" : {
            ipmi_password : "<node-serial-1>"
        }
        "<node-serial-2>" : {
            ipmi_password : "<node-serial-2>"
        }
        "<node-serial-3>" : {
            ipmi_password : "<node-serial-3>"
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