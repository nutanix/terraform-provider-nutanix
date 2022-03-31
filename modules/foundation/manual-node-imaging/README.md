# manual-node-imaging

This module is used to image nodes by giving all required details. Defaults can be used to apply certian details all nodes in one go.

Note : This module can only image all kinds of node - cvm running, discovery os installed or bare metal.

## Resources & DataSources used

1. nutanix_foundation_image_nodes resource

## Usage

Basic example of usage. 

```hcl
module "batch1" {
    source = "<local-path-to-nutanix-terraform-provider-repo>/terraform-provider-nutanix/modules/foundation/manual-node-imaging/"
    timeout = 120
    cvm_netmask = "xx.xx.xx.xx"
    cvm_gateway = "xx.xx.xx.xx"
    hypervisor_gateway = "xx.xx.xx.xx"
    hypervisor_netmask = "xx.xx.xx.xx"

    defaults = {
        ipmi_netmask : "xx.xx.xx.xx"
        ipmi_gateway : "xx.xx.xx.xx"
        ipmi_user : "<default-ipmi-username>"
        cvm_gb_ram : 50
        hypervisor : "kvm"
    }

    nos_package = data.nutanix_foundation_nos_packages.nos.entities[0]

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
                    ipmi_password: "<ipmi-password-3>"
                }
            ]
        }
    ]
    
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

```
