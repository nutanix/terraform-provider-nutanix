provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}

resource "nutanix_virtual_machine" "vm1" {
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
    subnet_reference = {
      kind = "subnet"
      uuid = "7206a75c-717a-4e72-b91e-16352971a25a"
    }
  }]
}
