#Here we will create a vm using an existing image and ngt install configuration
#Note - Replace appropriate values of variables in terraform.tfvars & resources/ubuntu.tpl
terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = "1.3.0"
    }
  }
}

#nutanix provider configuration
provider "nutanix" {
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_endpoint
  port         = var.nutanix_port
  insecure     = true
  wait_timeout = 10
}

#pull cluster data
data "nutanix_clusters" "clusters" {
}

#pull desired cluster data from setup
locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

#pull desired image data
data "nutanix_image" "test-ubuntu-cloud-init" {
  image_name = "test-ubuntu-cloud-init"
}

# data "nutanix_subnet" "dhcp-lan" {
#     subnet_name = "Lab_DHCP_VLAN_41"
# }

# data "nutanix_subnet" "test-managed-network-1234" {
#     subnet_name = "test-managed-network-1234"
# }

#creating managed subnet
resource "nutanix_subnet" "test-managed-network-0" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information
  name        = "test-managed-network-0"
  vlan_id     = 0
  subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
  subnet_ip = "10.xx.xx.xx"

  default_gateway_ip = "10.xx.xx.xx"
  prefix_length      = 22

  dhcp_options = {
    boot_file_name   = "bootfile"
    domain_name      = "lab"
    tftp_server_name = "10.xx.xx.xx"
  }

  dhcp_server_address = {
    ip = "10.xx.xx.xx"
  }

  dhcp_domain_name_server_list = ["10.xx.xx.xx"]
  dhcp_domain_search_list      = ["lab.local"]
  ip_config_pool_list_ranges   = ["10.xx.xx.xx 10.xx.xx.xx"] 
}

#creating virtual_machine
resource "nutanix_virtual_machine" "test-terra" {
  name                 = var.vm_name
  description          = "test_description"
  num_vcpus_per_socket = 4
  num_sockets          = 1
  memory_size_mib      = 2048
  cluster_uuid = local.cluster1

  guest_customization_cloud_init_user_data = base64encode(templatefile("${path.module}/resources/cloud-init/ubuntu.tpl", { hostname = var.vm_name }))


  # nic_list {
  #     subnet_uuid = data.nutanix_subnet.dhcp-lan.id
  # }

  nic_list {
    #   subnet_uuid = data.nutanix_subnet.test-managed-network-1234.id
      subnet_uuid = nutanix_subnet.test-managed-network-0.id
      ip_endpoint_list {
                ip   = "10.xx.xx.xx"
                type = "ASSIGNED"
    }
  }
  
  disk_list {    
    data_source_reference = {
        kind = "image"
        uuid = data.nutanix_image.test-ubuntu-cloud-init.id
      }
      

    device_properties {
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }

      device_type = "DISK"
    }
  }
  disk_list {
    disk_size_mib   = 100000
    disk_size_bytes = 104857600000
  }

  disk_list {
    disk_size_bytes = 0

    data_source_reference = {}

    device_properties {
      device_type = "CDROM"
      disk_address = {
        device_index = "1"
        adapter_type = "SATA"
      }
    }
  }

  nutanix_guest_tools = {
    state           = "ENABLED",
    # ngt_state       = "INSTALLED",
    iso_mount_state = "MOUNTED"
  }

  ngt_enabled_capability_list = [
    "SELF_SERVICE_RESTORE",
    "VSS_SNAPSHOT"
  ]

  # ngt_credentials = {
  #   username = var.ngt_credentials_username_linux
  #   password = var.ngt_credentials_password_linux
  # }

}

#output the ip address of above created vm
output "ip_address" {
  value = nutanix_virtual_machine.test-terra.nic_list_status[0].ip_endpoint_list[0].ip
}
