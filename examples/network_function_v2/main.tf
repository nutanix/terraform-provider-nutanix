terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.4.1"
    }
  }
}

provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = var.nutanix_port
  insecure = true
}

data "nutanix_clusters_v2" "clusters" {}

data "nutanix_images_v2" "nf_vm_image" {
  filter = "name eq '${var.image_name}'"
  limit  = 1
}

locals {
  candidate_clusters = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities : cluster
    if try(cluster.config[0].cluster_function[0], "") != "PRISM_CENTRAL" && (
      var.cluster_name == "" || cluster.name == var.cluster_name
    )
  ]

  cluster_ext_id = local.candidate_clusters[0].ext_id

  nf_vm_names = {
    primary = "${var.network_function_name}-primary"
    standby = "${var.network_function_name}-standby"
  }

  nf_vm_cloud_init = <<-EOT
    #cloud-config
    chpasswd:
      list: |
        ubuntu:${var.nf_vm_admin_password}
      expire: false
    disable_root: false
    ssh_pwauth: true
    package_update: true
    packages:
      - bridge-utils
    runcmd:
      - iface1=$(ls /sys/class/net/ | grep -E '^e' | sort | head -1)
      - iface2=$(ls /sys/class/net/ | grep -E '^e' | sort | head -2 | tail -1)
      - ip link set dev "$iface1" up
      - ip link set dev "$iface2" up
      - brctl addbr br0
      - brctl addif br0 "$iface1"
      - brctl addif br0 "$iface2"
      - ip link set dev br0 up
  EOT
}

resource "nutanix_subnet_v2" "management" {
  name              = var.management_subnet_name
  description       = "Management subnet for the network function VMs"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = var.management_subnet_vlan_id
  # NF VMs mix NETWORK_FUNCTION_NIC and NORMAL_NIC interfaces, which requires
  # an advanced-networking subnet for the regular management NIC.
  is_advanced_networking = true

  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = var.management_subnet_network
        }
        prefix_length = var.management_subnet_prefix_length
      }

      default_gateway_ip {
        value = var.management_subnet_gateway
      }

      pool_list {
        start_ip {
          value = var.management_subnet_pool_start
        }
        end_ip {
          value = var.management_subnet_pool_end
        }
      }
    }
  }
}

resource "nutanix_virtual_machine_v2" "nf_vm" {
  for_each = local.nf_vm_names

  name                         = each.value
  description                  = "Network function VM (${each.key}) managed by Terraform"
  num_cores_per_socket         = var.vm_num_cores_per_socket
  num_sockets                  = var.vm_num_sockets
  memory_size_bytes            = var.vm_memory_size_bytes
  is_agent_vm                  = false
  hardware_clock_timezone      = "UTC"
  is_memory_overcommit_enabled = false

  cluster {
    ext_id = local.cluster_ext_id
  }

  apc_config {
    is_apc_enabled = false
  }

  disks {
    backing_info {
      vm_disk {
        disk_size_bytes = var.vm_disk_size_bytes
        data_source {
          reference {
            image_reference {
              image_ext_id = data.nutanix_images_v2.nf_vm_image.images[0].ext_id
            }
          }
        }
      }
    }

    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
  }

  guest_customization {
    config {
      cloud_init {
        cloud_init_script {
          user_data {
            value = base64encode(local.nf_vm_cloud_init)
          }
        }
      }
    }
  }

  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "INGRESS"
      }
    }
  }

  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type                  = "NETWORK_FUNCTION_NIC"
        network_function_nic_type = "EGRESS"
      }
    }
  }

  nics {
    nic_network_info {
      virtual_ethernet_nic_network_info {
        nic_type = "NORMAL_NIC"
        subnet {
          ext_id = nutanix_subnet_v2.management.id
        }
        ipv4_config {
          should_assign_ip = true
        }
      }
    }
  }

  power_state = "ON"

  lifecycle {
    ignore_changes = [
      cd_roms,
      guest_customization,
      nics.2.nic_network_info.0.virtual_ethernet_nic_network_info.0.ipv4_config,
    ]
  }
}

locals {
  nf_vm_details = {
    for name, vm in nutanix_virtual_machine_v2.nf_vm :
    name => {
      vm_ext_id = vm.id
      ingress_nic_ext_id = [
        for nic in vm.nics : nic.ext_id
        if try(nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type, "") == "INGRESS"
      ][0]
      egress_nic_ext_id = [
        for nic in vm.nics : nic.ext_id
        if try(nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].network_function_nic_type, "") == "EGRESS"
      ][0]
      management_ips = flatten([
        for nic in vm.nics : [
          for address in try(nic.nic_network_info[0].virtual_ethernet_nic_network_info[0].ipv4_info[0].learned_ip_addresses, []) : address.value
        ]
      ])
    }
  }
}

resource "nutanix_network_function_v2" "nf" {
  name                    = var.network_function_name
  description             = var.network_function_description
  high_availability_mode  = "ACTIVE_PASSIVE"
  failure_handling        = "FAIL_CLOSE"
  traffic_forwarding_mode = "INLINE"

  data_plane_health_check_config {
    failure_threshold = 3
    interval_secs     = 5
    success_threshold = 3
    timeout_secs      = 2
  }

  dynamic "nic_pairs" {
    for_each = local.nf_vm_details

    content {
      ingress_nic_reference = nic_pairs.value.ingress_nic_ext_id
      egress_nic_reference  = nic_pairs.value.egress_nic_ext_id
      vm_reference          = nic_pairs.value.vm_ext_id
      is_enabled            = true
    }
  }
}

data "nutanix_network_function_v2" "nf" {
  ext_id = nutanix_network_function_v2.nf.ext_id
}

