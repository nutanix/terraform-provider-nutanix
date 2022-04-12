# node-serials-filter

This module is used to image nodes and create cluster given only node serials. All other params for node imaging like ipmi_ip, cvm_ip, hypervisor_ip, etc is obtained internaly by data_sources and passed as node imaging input. Also we can give override these fields as well. 

## dataresource and Resource used
1. nutanix_foundation_central_imaged_nodes_list
2. nutanix_foundation_central_imaged_node_details
3. nutanix_foundation_central_image_cluster resource

## Usage

Basic example of usage. This gets nodes information from data sources and uses them for imaging. We only provide node_serials and common_network_setting .

```hcl
module "batch1"{
    source = "<local-path-to-nutanix-terraform-provider-repo>/terraform-provider-nutanix/modules/modulesFc/FoundationCentral/aos_based_imaging/"

    cluster_name = "test_cluster"
    
    node_info = {
        "Node-Serial-1": {}

        "Node-Serial-2": {}

        "Node-Serial-3": {}
    }

    common_network_settings={
        cvm_dns_servers:[
            "xx.x.x.xx"
        ]
        hypervisor_dns_servers:[
            "xx.x.x.xx"
        ]
        cvm_ntp_servers:[
            "server-ip"
        ]
        hypervisor_ntp_servers:[
            "server-ip"
        ]
    }
    aos_package_url="package_url"
}
```