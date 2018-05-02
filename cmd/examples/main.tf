provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}

resource "nutanix_virtual_machine" "my-machine" {
  metadata {
    kind = "vm"
    name = "metadata-name-test-dou-%d"
  }

  name = "name-test-dou-%d"

  cluster_reference = {
    kind = "cluster"
    uuid = "000567f3-1921-c722-471d-0cc47ac31055"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  power_state          = "ON"

  nic_list = [{
    nic_type = "NORMAL_NIC"

    subnet_reference = {
      kind = "subnet"
      uuid = "af25e4af-12bd-4b1a-98ca-efcb0e9c06f8"
    }

    network_function_nic_type = "INGRESS"
  }]

  disk_list = [
    {
      data_source_reference = [{
        kind = "image"
        name = "Centos7"
        uuid = "9eabbb39-1baf-4872-beaf-adedcb612a0b"
      }]

      device_properties = [{
        device_type = "DISK"
      }]

      disk_size_mib = 1
    },
  ]

  disk_list = [
    {
      data_source_reference = {
        kind = "image"
        name = "Centos7"
        uuid = "9eabbb39-1baf-4872-beaf-adedcb612a0b"
      }

      device_properties = {
        device_type = "DISK"
      }

      disk_size_mib = 1
    },
    {
      data_source_reference = {
        kind = "image"
        name = "Centos7"
        uuid = "9eabbb39-1baf-4872-beaf-adedcb612a0b"
      }

      device_properties = {
        device_type = "DISK"
      }

      disk_size_mib = 1
    },
    {
      data_source_reference = {
        kind = "image"
        name = "Centos7"
        uuid = "9eabbb39-1baf-4872-beaf-adedcb612a0b"
      }

      device_properties = {
        device_type = "DISK"
      }

      disk_size_mib = 1
    },
  ]
}

output "ip" {
  value = "${nutanix_virtual_machine.my-machine.ip_address}"
}
