---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_image_cluster"
sidebar_current: "docs-nutanix-resource-foundation-central-image-cluster"
description: |-
  Image Nodes and Create a cluster out of nodes registered with Foundation Central.
---

# nutanix_foundation_central_image_cluster

Image Nodes and Create a cluster out of nodes registered with Foundation Central.

## Example Usage

```hcl
resource "nutanix_foundation_central_image_cluster" "img2"{
  cluster_name = "test-FC"
  cluster_external_ip = "<CLUSTER-IP>"
  common_network_settings{
    cvm_dns_servers=[
        "xx.x.xx.xx"
    ]
    hypervisor_dns_servers=[
        "xx.x.xx.xx"
    ]
    cvm_ntp_servers=[
        "<cvm-ntp>"
    ]
    hypervisor_ntp_servers=[
        "<hypervisor-ntp>"
    ]
  }
    redundancy_factor = 2
    node_list{
      cvm_gateway="10.xx.xx.xx"
      cvm_netmask="xx.xx.xx.xx"
      cvm_ip="10.x.xx.xx"
      hypervisor_gateway="10.x.x.xx"
      hypervisor_netmask="xx.xx.xx.xx"
      hypervisor_ip="10.x.xx.xx"
      hypervisor_hostname="HOST-1"
      imaged_node_uuid="<NODE-UUID>"
      use_existing_network_settings=false
      ipmi_gateway="10.x.xx.xx"
      ipmi_netmask="10.x.xx.xx"
      ipmi_ip="10.x.xx.xx"
      image_now=true
      hypervisor_type="kvm"
      hardware_attributes_override = {
        default_workload="vdi"
        lcm_family= "smc_gen_10"
        maybe_1GbE_only= true
        robo_mixed_hypervisor= true
      }
    }
    node_list{
        cvm_gateway="10.xx.xx.xx"
        cvm_netmask="xx.xx.xx.xx"
        cvm_ip="10.x.xx.xx"
        hypervisor_gateway="10.x.x.xx"
        hypervisor_netmask="xx.xx.xx.xx"
        hypervisor_ip="10.x.xx.xx"
        hypervisor_hostname="HOST-2"
        imaged_node_uuid="<NODE-UUID>"
        use_existing_network_settings=false
        ipmi_gateway="10.x.xx.xx"
        ipmi_netmask="10.x.xx.xx"
        ipmi_ip="10.x.xx.xx"
        image_now=true
        hypervisor_type="kvm"
    }
    node_list{
        cvm_gateway="10.xx.xx.xx"
        cvm_netmask="xx.xx.xx.xx"
        cvm_ip="10.x.xx.xx"
        hypervisor_gateway="10.x.x.xx"
        hypervisor_netmask="xx.xx.xx.xx"
        hypervisor_ip="10.x.xx.xx"
        hypervisor_hostname="HOST-3"
        imaged_node_uuid="<NODE-UUID>"
        use_existing_network_settings=false
        ipmi_gateway="10.x.xx.xx"
        ipmi_netmask="10.x.xx.xx"
        ipmi_ip="10.x.xx.xx"
        image_now=true
        hypervisor_type="kvm"
    }
    aos_package_url="<URL>"

    // required for deploying AOS >= v6.8
    hypervisor_isos{
      url="<hypervisor-installer-link>"
      sha256sum="<hypervisor-installer-checksum>"
      hypervisor_type = "kvm"
    }
    //pass true to skip cluster creation
    skip_cluster_creation = true


}

```


## Argument Reference

The following arguments are supported:

* `cluster_external_ip`: External management ip of the cluster.
* `common_network_settings`: Common network settings across the nodes in the cluster. 
* `hypervisor_iso_details`: Details of the hypervisor iso. (Deprecated)
* `hypervisor_isos`: Details of the hypervisor iso. Required for deploying node with AOS >= 6.8
* `storage_node_count`: Number of storage only nodes in the cluster. AHV iso for storage node will be taken from aos package.
* `redundancy_factor`: Redundancy factor of the cluster.
* `cluster_name`: Name of the cluster.
* `aos_package_url`: URL to download AOS package. Required only if imaging is needed.
* `cluster_size`: Number of nodes in the cluster.
* `aos_package_sha256sum`: Sha256sum of AOS package.
* `timezone`: Timezone to be set on the cluster.
* `nodes_list`: List of details of nodes out of which the cluster needs to be created.

### common network settings
* `cvm_dns_servers`: List of dns servers for the cvms in the cluster.
* `hypervisor_dns_servers`: List of dns servers for the hypervisors in the cluster.
* `cvm_ntp_servers`: List of ntp servers for the cvms in the cluster.
* `hypervisor_ntp_servers`: List of ntp servers for the hypervisors in the cluster.

### hypervisor iso details
* `hyperv_sku`: SKU of hyperv to be installed if hypervisor_type is hyperv.
* `url`: (Required) URL to download hypervisor iso. Required only if imaging is needed. 
* `hyperv_product_key`: Product key for hyperv isos. Required only if the hypervisor type is hyperv and product key is mandatory (ex: for volume license).
* `sha256sum`: sha256sum of the hypervisor iso.

### hypervisor isos
* `hyperv_sku`: SKU of hyperv to be installed if hypervisor_type is hyperv.
* `url`: (Required) URL to download hypervisor iso. Required only if imaging is needed. 
* `hyperv_product_key`: Product key for hyperv isos. Required only if the hypervisor type is hyperv and product key is mandatory (ex: for volume license).
* `sha256sum`: sha256sum of the hypervisor iso.
* `hypervisor_type`: Hypervisor type. Only supports "kvm", "esx" and "hyperv"

### node list
* `cvm_gateway`: Gateway of the cvm.
* `ipmi_netmask`: Netmask of the ipmi.
* `rdma_passthrough`: Passthrough RDMA nic to CVM if possible, default to false.
* `imaged_node_uuid`: UUID of the node.
* `cvm_vlan_id`: Vlan tag of the cvm, if the cvm is on a vlan.
* `hypervisor_type`: Type of hypervisor to be installed. Must be one of {kvm, esx, hyperv}.
* `image_now`: True, if the node should be imaged, False, otherwise.
* `hypervisor_hostname`: Name to be set for the hypervisor host.
* `hypervisor_netmask`: Netmask of the hypervisor.
* `cvm_netmask`: Netmask of the cvm.
* `ipmi_ip`: IP address to be set for the ipmi of the node.
* `hypervisor_gateway`: Gateway of the hypervisor.
* `hardware_attributes_override`: Hardware attributes override json for the node.
* `cvm_ram_gb`: Amount of memory to be assigned for the cvm.
* `cvm_ip`: IP address to be set for the cvm on the node.
* `hypervisor_ip`: IP address to be set for the hypervisor on the node.
* `use_existing_network_settings`: Decides whether to use the existing network settings for the node. If True, the existing network settings of the node will be used during cluster creation. If False, then client must provide new network settings. If all nodes are booted in phoenix, this field is, by default, considered to be False.
* `ipmi_gateway`: Gateway of the ipmi.


## Attributes Reference

The following attributes are exported:

* `imaged_cluster_uuid`: Unique id of the cluster.

## Error 

Incase of error in any individual node or cluster, terraform will error our after full imaging process is completed. Error will be shown for every failed node and cluster.

## lifecycle

* `Update` : - Resource will trigger new resource create call for any kind of update in resource config.
* `delete` : - Resource will be deleted from Foundation Central deployment history. For Actual Cluster delete , manually destroy the cluster.   

See detailed information in [Nutanix Foundation Central Create a Cluster](https://www.nutanix.dev/api_references/foundation-central/#/cba507f282927-request-to-create-a-cluster).