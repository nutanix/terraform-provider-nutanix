provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}

variable clusterid {
  default = "000567f3-1921-c722-471d-0cc47ac31055"
}

resource "nutanix_virtual_machine" "vm1" {
  metadata {
    kind = "vm"
    name = "metadata-name-test-dou"
  }

  name = "test-dou"

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.clusterid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  power_state          = "ON"

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.test.id}"
    }

    ip_endpoint_list = {
      ip   = "192.168.0.10"
      type = "ASSIGNED"
    }
  }]

  disk_list = [{
    data_source_reference = [{
      kind = "image"
      name = "Centos7"
      uuid = "${nutanix_image.test.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}

resource "nutanix_subnet" "test" {
  metadata = {
    kind = "subnet"
  }

  name        = "dou_vlan0_test"
  description = "Dou Vlan 0"

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.clusterid}"
  }

  vlan_id     = 201
  subnet_type = "VLAN"

  prefix_length      = 24
  default_gateway_ip = "192.168.0.1"
  subnet_ip          = "192.168.0.0"

  dhcp_options {
    boot_file_name   = "bootfile"
    tftp_server_name = "192.168.0.252"
    domain_name      = "nutanix"
  }

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "calm.io"]
}

resource "nutanix_image" "test" {
  metadata = {
    kind = "image"
  }

  name        = "dou_image_%d"
  description = "Dou Image Test"
  name        = "CentOS7-ISO"
  source_uri  = "http://10.7.1.7/data1/ISOs/CentOS-7-x86_64-Minimal-1503-01.iso"

  checksum = {
    checksum_algorithm = "SHA_256"
    checksum_value     = "a9e4e0018c98520002cd7cf506e980e66e31f7ada70b8fc9caa4f4290b019f4f"
  }
}

data "nutanix_virtual_machine" "nutanix_virtual_machine" {
  vm_id = "${nutanix_virtual_machine.vm1.id}"
}
