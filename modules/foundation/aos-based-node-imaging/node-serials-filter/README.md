# node-serials-filter

This module is used to image nodes given only node serials and ipmi creds. All other info for node imaging like cvm_ip, hypervisor_ip, etc is obtained internaly by data_sources and passed as node imaging input. Also we can give override these fields as well, in case we don't want to use existing network configuration.

Note : This module can only reimage nodes which are discoverable and cvm is running on it.

## Resources & Modules used

1. foundation/discover-nodes-network-details module
2. nutanix_foundation_image_nodes resource

## Usage

Basic example of usage. This gets all network & hypervisor information from data sources and uses them for imaging. We only provide ipmi creds and defaults.

```hcl
module "batch1" {

    // source where module code is present in local machine
    source = "<local-path-to-nutanix-terraform-provider-repo>/terraform-provider-nutanix/modules/foundation/aos-based-node-imaging/node-serials-filter/"
    timeout = 120
    cvm_netmask = "xx.xx.xx.xx"
    cvm_gateway = "xx.xx.xx.xx"
    hypervisor_gateway = "xx.xx.xx.xx"
    hypervisor_netmask = "xx.xx.xx.xx"
    nos_package = data.nutanix_foundation_nos_packages.nos.entities[0]

    defaults = {
        ipmi_user : "<ipmi-user-1>"
        cvm_gb_ram : 50
        hypervisor : "kvm"
    }

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

```
