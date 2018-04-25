provider "nutanix" {
  username = ""
  password = ""
  endpoint = "10.5.68.6"
  insecure = true
}

resource "nutanix_virtual_machine" "vm1" {
  metadata {
    categories {
      "Project" = "nucalm"
    }
  }

  name = "test 1"

  resource {
    nic_list = [{
      nic_type                  = "NORMAL_NIC"
      network_function_nic_type = "INGRESS"

      subnet_reference = {
        kind = "subnet"
        uuid = "c03ecf8f-aa1c-4a07-af43-9f2f198713c0"
      }
    }]

    num_vcpus_per_socket = 1
    num_sockets          = 1
    memory_size_mb       = 2048
    power_state          = "On"

    disk_list = [{
      data_source_reference = {
        kind = "image"
        name = "Centos7"
        uuid = "9eabbb39-1baf-4872-beaf-adedcb612a0b"
      }

      device_properties = {
        device_type = "DISK"
      }

      disk_size_mib = 1
    }]
  }
}
