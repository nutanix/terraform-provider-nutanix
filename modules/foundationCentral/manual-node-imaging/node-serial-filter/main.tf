//resource for node imaging and cluster creation
resource "nutanix_foundation_central_image_cluster" "this" {
  //required fields
  aos_package_url   = var.aos_package_url
  cluster_name      = var.cluster_name
  redundancy_factor = var.redundancy_factor

  //Optional fields
  storage_node_count    = var.storage_node_count
  cluster_size          = var.cluster_size
  aos_package_sha256sum = var.aos_package_sha256sum
  timezone              = var.timezone
  cluster_external_ip   = var.cluster_external_ip
  skip_cluster_creation = var.skip_cluster_creation

  // define common network settings
  dynamic "common_network_settings" {
    for_each = [var.common_network_settings]
    content {
      cvm_dns_servers        = common_network_settings.value.cvm_dns_servers
      hypervisor_dns_servers = common_network_settings.value.hypervisor_dns_servers
      cvm_ntp_servers        = common_network_settings.value.cvm_ntp_servers
      hypervisor_ntp_servers = common_network_settings.value.hypervisor_ntp_servers
    }
  }

  // define node list 
  dynamic "node_list" {
    for_each = var.node_list
    content {
      cvm_gateway                   = node_list.value.cvm_gateway
      cvm_netmask                   = node_list.value.cvm_netmask
      cvm_ip                        = node_list.value.cvm_ip
      hypervisor_gateway            = node_list.value.hypervisor_gateway
      hypervisor_netmask            = node_list.value.hypervisor_netmask
      hypervisor_ip                 = node_list.value.hypervisor_ip
      hypervisor_hostname           = node_list.value.hypervisor_hostname
      imaged_node_uuid              = node_list.value.imaged_node_uuid
      use_existing_network_settings = node_list.value.use_existing_network_settings
      ipmi_gateway                  = node_list.value.ipmi_gateway
      ipmi_netmask                  = node_list.value.ipmi_netmask
      ipmi_ip                       = node_list.value.ipmi_ip
      image_now                     = node_list.value.image_now
      hypervisor_type               = node_list.value.hypervisor_type
    }
  }
}