# manual-node-imaging

This module is used to image nodes and create cluster by providing all the required details.

## Resource used
1. nutanix_foundation_central_image_cluster resource

## Usage

Basic example of usage. 

```hcl
module "batch"{
    source = "<local-path-to-nutanix-terraform-provider-repo>/terraform-provider-nutanix/modules//modulesFc/FoundationCentral/manual_node_imaging"

    cluster_name = "test_cluster"
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
    node_list = [{
      cvm_gateway="xx.x.xxx.x"
      cvm_netmask="xxx.xxx.xxx.x"
      cvm_ip="xx.xx.xx.xx"
      hypervisor_gateway="xx.xx.xx.xx"
      hypervisor_netmask="xxx.xxx.xxx.x"
      hypervisor_ip="xx.xx.xx.xx"
      hypervisor_hostname="HOST-3"
      imaged_node_uuid="<imaged_node_uuid>"
      use_existing_network_settings=false
      ipmi_gateway="xx.xx.xx.xx"
      ipmi_netmask="xxx.xxx.xxx.x"
      ipmi_ip="xx.xx.xx.xx"
      image_now=true
      hypervisor_type="kvm"
    },
    {
       cvm_gateway="xx.x.xxx.x"
      cvm_netmask="xxx.xxx.xxx.x"
      cvm_ip="xx.xx.xx.xx"
      hypervisor_gateway="xx.xx.xx.xx"
      hypervisor_netmask="xxx.xxx.xxx.x"
      hypervisor_ip="xx.xx.xx.xx"
      hypervisor_hostname="HOST-2"
      imaged_node_uuid="<imaged_node_uuid>"
      use_existing_network_settings=false
      ipmi_gateway="xx.xx.xx.xx"
      ipmi_netmask="xxx.xxx.xxx.x"
      ipmi_ip="xx.xx.xx.xx"
      image_now=true
      hypervisor_type="kvm"
    }]
    redundancy_factor = 2
    aos_package_url="<aos_package_url>"
}

```