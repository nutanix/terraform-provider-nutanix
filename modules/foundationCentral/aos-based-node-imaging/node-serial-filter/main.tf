// datasource to get list of nodes
data "nutanix_foundation_central_imaged_nodes_list" "nodelist" {}


locals {
  // create list of node serials
  nodeinfo = keys(var.node_info)
  // get the imaged_node_uuid from node_serials
  nodedata = flatten([
    for node in local.nodeinfo :
    [
      for nl in data.nutanix_foundation_central_imaged_nodes_list.nodelist.imaged_nodes :
      nl.imaged_node_uuid if nl.node_serial == node
    ]
  ])
}

// datasource to get the node details
data "nutanix_foundation_central_imaged_node_details" "nodedetails" {
  for_each         = toset(local.nodedata)
  imaged_node_uuid = each.value
}

// resource for node imaging and cluster creation
resource "nutanix_foundation_central_image_cluster" "this" {

  //Required field to be taken as input from module
  aos_package_url   = var.aos_package_url
  redundancy_factor = var.redundancy_factor
  cluster_name      = var.cluster_name

  //Optional fields
  storage_node_count    = var.storage_node_count
  cluster_size          = var.cluster_size
  aos_package_sha256sum = var.aos_package_sha256sum
  timezone              = var.timezone
  cluster_external_ip   = var.cluster_external_ip
  skip_cluster_creation = var.skip_cluster_creation


  // define common network setting as input given in module
  dynamic "common_network_settings" {
    for_each = [var.common_network_settings]
    content {
      cvm_dns_servers        = common_network_settings.value.cvm_dns_servers
      hypervisor_dns_servers = common_network_settings.value.hypervisor_dns_servers
      cvm_ntp_servers        = common_network_settings.value.cvm_ntp_servers
      hypervisor_ntp_servers = common_network_settings.value.hypervisor_ntp_servers
    }
  }

  // Dynamically define  nodes as input given in module
  dynamic "node_list" {
    for_each = data.nutanix_foundation_central_imaged_node_details.nodedetails
    content {
      cvm_gateway                   = var.node_info[node_list.value.node_serial].cvm_gateway != null ? var.node_info[node_list.value.node_serial].cvm_gateway : node_list.value.cvm_gateway
      cvm_netmask                   = var.node_info[node_list.value.node_serial].cvm_netmask != null ? var.node_info[node_list.value.node_serial].cvm_netmask : node_list.value.cvm_netmask
      cvm_ip                        = var.node_info[node_list.value.node_serial].cvm_ip != null ? var.node_info[node_list.value.node_serial].cvm_ip : node_list.value.cvm_ip
      hypervisor_gateway            = var.node_info[node_list.value.node_serial].hypervisor_gateway != null ? var.node_info[node_list.value.node_serial].hypervisor_gateway : node_list.value.hypervisor_gateway
      hypervisor_netmask            = var.node_info[node_list.value.node_serial].hypervisor_netmask != null ? var.node_info[node_list.value.node_serial].hypervisor_netmask : node_list.value.hypervisor_netmask
      hypervisor_ip                 = var.node_info[node_list.value.node_serial].hypervisor_ip != null ? var.node_info[node_list.value.node_serial].cvm_gateway : node_list.value.hypervisor_ip
      hypervisor_hostname           = var.node_info[node_list.value.node_serial].hypervisor_hostname != null ? var.node_info[node_list.value.node_serial].hypervisor_hostname : node_list.value.hypervisor_hostname
      imaged_node_uuid              = var.node_info[node_list.value.node_serial].imaged_node_uuid != null ? var.node_info[node_list.value.node_serial].imaged_node_uuid : node_list.value.imaged_node_uuid
      use_existing_network_settings = var.node_info[node_list.value.node_serial].use_existing_network_settings != null ? var.node_info[node_list.value.node_serial].use_existing_network_settings : null
      ipmi_gateway                  = var.node_info[node_list.value.node_serial].ipmi_gateway != null ? var.node_info[node_list.value.node_serial].ipmi_gateway : node_list.value.ipmi_gateway
      ipmi_netmask                  = var.node_info[node_list.value.node_serial].ipmi_netmask != null ? var.node_info[node_list.value.node_serial].ipmi_netmask : node_list.value.ipmi_netmask
      ipmi_ip                       = var.node_info[node_list.value.node_serial].ipmi_ip != null ? var.node_info[node_list.value.node_serial].ipmi_ip : node_list.value.ipmi_ip
      image_now                     = true
      cvm_vlan_id                   = var.node_info[node_list.value.node_serial].cvm_vlan_id != null ? var.node_info[node_list.value.node_serial].cvm_vlan_id : node_list.value.cvm_vlan_id
      hardware_attributes_override  = var.node_info[node_list.value.node_serial].hardware_attributes_override != null ? var.node_info[node_list.value.node_serial].hardware_attributes_override : node_list.value.hardware_attributes
      cvm_ram_gb                    = var.node_info[node_list.value.node_serial].cvm_ram_gb != null ? var.node_info[node_list.value.node_serial].cvm_gateway : 0
      hypervisor_type               = var.node_info[node_list.value.node_serial].hypervisor_type != null ? var.node_info[node_list.value.node_serial].hypervisor_type : node_list.value.hypervisor_type
    }
  }
  // define hypervisor iso as input given
  dynamic "hypervisor_iso_details" {
    for_each = var.hypervisor_iso_details != null ? [var.hypervisor_iso_details] : []
    content {
      hyperv_sku         = hypervisor_iso_details.value.hyperv_sku
      url                = hypervisor_iso_details.value.url
      hyperv_product_key = hypervisor_iso_details.value.hyperv_product_key
      sha256sum          = hypervisor_iso_details.value.sha256sum
    }
  }
}