/*
Description:
- Here we will create an automation which picks up config.json, containing the blocks of nodes & 
  cluster info, create spec dynamically as per the struture of config.json and do imaging of nodes.
- You can download this config.json from install.nutanix.com or create one using the sample config.json given.
- Here we are following Bare metal workflow, imaging nodes using IPMI.
*/

/*
[IMPORTANT]
- Please note that this is just an example and the spec generated dynamically is only having minimal fields which are required
  for imaging the blocks of nodes and creating cluster. You can add more fields as per the requirements. 
- This example was created as per reference to config.json downloaded from install.nutanix.com -> Foundation Preconfiguration.
- Keep config.json in same directory where this .tf file is kept
*/



// pull provider plugin with appropriate version
terraform {
  required_providers{
      nutanix = {
          source = "nutanix/nutanix"
          version = "1.5.0-beta"
      }
  }
}

// give provider config
provider "nutanix" {
    foundation_endpoint = "10.xx.xx.xx"
}

// import config.json . Replace the file location if required.
locals{
    config = (jsondecode(file("config.json"))).config
}

// pull nos packages info
data "nutanix_foundation_nos_packages" "nos"{}

resource "nutanix_foundation_image_nodes" "batch1" {

  // give required info
  ipmi_netmask = local.config.ipmi_netmask
  ipmi_gateway = local.config.ipmi_gateway
  cvm_netmask = local.config.cvm_netmask
  cvm_gateway = local.config.cvm_gateway
  hypervisor_netmask = local.config.hypervisor_netmask
  hypervisor_gateway = local.config.hypervisor_gateway

  // use nos package info from data source
  nos_package = data.nutanix_foundation_nos_packages.nos.entities[0]

  // this will dynamically create multiple blocks of multiple nodes spec using array of blocks in config
  dynamic "blocks" {
      for_each = local.config.blocks
      content{
        block_id = blocks.value.block_id
        dynamic "nodes" {
            for_each = blocks.value.nodes
            content {
                ipmi_ip = nodes.value.ipmi_ip
                ipmi_user = nodes.value.ipmi_user
                ipmi_password = nodes.value.ipmi_password
                cvm_ip = nodes.value.cvm_ip
                image_now = true
                hypervisor_ip = nodes.value.hypervisor_ip
                hypervisor = "kvm"
                hypervisor_hostname = nodes.value.hypervisor_hostname
                node_position = nodes.value.node_position
            }
        }
      }
  }

  // this will create multiple clusters spec as per array of cluster in config file
  dynamic "clusters"{
      for_each = local.config.clusters
      content{
        cluster_name = clusters.value.cluster_name
        redundancy_factor = clusters.value.redundancy_factor
        cluster_external_ip = clusters.value.cluster_external_ip
        cluster_members = clusters.value.cluster_members
        single_node_cluster = length(clusters.value.cluster_members) > 1 ? false : true
        cluster_init_now = true
      }
  }
}