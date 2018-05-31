provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}

resource "nutanix_image" "test" {
  name        = "Ubuntu"
  description = "Ubuntu Server Mini ISO"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}

data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou"

  cluster_reference = {
    kind = "cluster"
    uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-iac-managed.id}"
    }

    ip_endpoint_list = {
      ip   = "10.6.80.10"
      type = "ASSIGNED"
    }
  },
    {
      subnet_reference = {
        kind = "subnet"
        uuid = "${nutanix_subnet.next-iac-managed2.id}"
      }

      ip_endpoint_list = {
        ip   = "10.5.80.10"
        type = "ASSIGNED"
      }
    },
  ]

  #What disk/cdrom configuration will this have?
  disk_list = [{
    data_source_reference = [{
      kind = "image"
      name = "ubuntu"
      uuid = "${nutanix_image.test.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}

resource "nutanix_subnet" "next-iac-managed" {
  cluster_reference = {
    kind = "cluster"
    uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  name        = "next-iac-managed-%d"
  vlan_id     = 3
  subnet_type = "VLAN"

  prefix_length = 20

  default_gateway_ip = "10.6.80.1"
  subnet_ip          = "10.6.80.0"

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "eng.nutanix.com"]
}

resource "nutanix_subnet" "next-iac-managed2" {
  cluster_reference = {
    kind = "cluster"
    uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  name        = "next-iac-managed-%d"
  vlan_id     = 4
  subnet_type = "VLAN"

  prefix_length = 20

  default_gateway_ip = "10.5.80.1"
  subnet_ip          = "10.5.80.0"

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "eng.nutanix.com"]
}
