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

# resource "nutanix_virtual_machine" "vm1" {
#   metadata {
#     kind = "vm"
#     name = "metadata-name-test-dou"
#   }

#   name = "test-dou"

#   cluster_reference = {
#     kind = "cluster"
#     uuid = "${var.clusterid}"
#   }

#   num_vcpus_per_socket = 1
#   num_sockets          = 1
#   memory_size_mib      = 2048
#   power_state          = "ON"

#   nic_list = [{
#     subnet_reference = {
#       kind = "subnet"
#       uuid = "${nutanix_subnet.test.id}"
#     }

#     ip_endpoint_list = {
#       ip   = "192.168.0.10"
#       type = "ASSIGNED"
#     }
#   }]

#   disk_list = [{
#     data_source_reference = [{
#       kind = "image"
#       name = "Centos7"
#       uuid = "${nutanix_image.test.id}"
#     }]

#     device_properties = [{
#       device_type = "DISK"
#     }]

#     disk_size_mib = 5000
#   }]
# }

# resource "nutanix_subnet" "test" {
#   metadata = {
#     kind = "subnet"
#   }

#   name        = "dou_vlan0_test"
#   description = "Dou Vlan 0"

#   cluster_reference = {
#     kind = "cluster"
#     uuid = "${var.clusterid}"
#   }

#   vlan_id     = 201
#   subnet_type = "VLAN"

#   prefix_length      = 24
#   default_gateway_ip = "192.168.0.1"
#   subnet_ip          = "192.168.0.0"

#   dhcp_options {
#     boot_file_name   = "bootfile"
#     tftp_server_name = "192.168.0.252"
#     domain_name      = "nutanix"
#   }

#   dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
#   dhcp_domain_search_list      = ["nutanix.com", "calm.io"]
# }

# resource "nutanix_image" "test" {
#   metadata = {
#     kind = "image"
#   }

#   name        = "dou_image_%d"
#   description = "Dou Image Test"
#   name        = "CentOS7-ISO"
#   source_uri  = "http://10.7.1.7/data1/ISOs/CentOS-7-x86_64-Minimal-1503-01.iso"

#   checksum = {
#     checksum_algorithm = "SHA_256"
#     checksum_value     = "a9e4e0018c98520002cd7cf506e980e66e31f7ada70b8fc9caa4f4290b019f4f"
#   }
# }

# data "nutanix_virtual_machine" "nutanix_virtual_machine" {
#   vm_id = "${nutanix_virtual_machine.vm1.id}"
# }

resource "nutanix_image" "centos73-install-iso" {
  name        = "iso_CentOS-7.3-x86_64-Minimal-1611"
  description = "Here is a CentOS 7.3 Install CD from Endor filer"
  source_uri  = "http://endor.dyn.nutanix.com/isos/linux/centos/7/CentOS-7.3-x86_64-Minimal-1611.iso"

  checksum = {
    checksum_algorithm = "SHA_256"
    checksum_value     = "27bd866242ee058b7a5754e83d8ee8403e216b93d130d800852a96f41c34d86a"
  }

  metadata = {
    kind = "image"
  }
}

resource "nutanix_image" "centos-lamp-app" {
  name        = "CentOS-LAMP-APP.qcow2"
  description = "CentOS LAMP - App"
  image_type  = "DISK_IMAGE"
  source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-APP.qcow2"

  metadata = {
    kind = "image"
  }
}

# resource "nutanix_image" "centos-lamp-db" {
#   name        = "CentOS-LAMP-DB.qcow2"
#   description = "CentOS LAMP - DB"
#   source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-DB.qcow2"

#   metadata = {
#     kind = "image"
#   }
# }

# resource "nutanix_image" "centos-lamp-haproxy" {
#   name        = "CentOS-LAMP-HAPROXY.qcow2"
#   description = "CentOS LAMP - HAProxy"
#   source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-HAProxy.qcow2"

#   metadata = {
#     kind = "image"
#   }
# }

resource "nutanix_subnet" "next-lamp-subnet" {
  metadata = {
    kind = "subnet"
  }

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.clusterid}"
  }

  name               = "next-lamp"
  description        = "lamp lamp lampy lamp vlan 0"
  vlan_id            = 301
  subnet_type        = "VLAN"
  prefix_length      = 24
  default_gateway_ip = "1.2.3.1"
  subnet_ip          = "1.2.3.0"

  dhcp_options {
    boot_file_name   = "bootfile"
    tftp_server_name = "1.2.3.200"
    domain_name      = "nutanix"
  }

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "eng.nutanix.com"]
}

# resource "nutanix_virtual_machine" "demo-01-web" {
#   metadata {
#     kind = "vm"
#   }

#   name                 = "demo-01-web"
#   description          = "demo Frontend Web Server"
#   num_vcpus_per_socket = 2
#   num_sockets          = 1
#   memory_size_mib      = 4096
#   power_state          = "ON"

#   cluster_reference = {
#     kind = "cluster"
#     uuid = "${var.clusterid}"
#   }

#   nic_list = [{
#     subnet_reference = {
#       kind = "subnet"
#       uuid = "${nutanix_subnet.next-lamp-subnet.id}"
#     }

#     ip_endpoint_list = {
#       ip   = "192.168.0.10"
#       type = "ASSIGNED"
#     }
#   }]

#   disk_list = [{
#     data_source_reference = [{
#       kind = "image"
#       name = "Centos7"
#       uuid = "${nutanix_image.centos-lamp-haproxy.id}"
#     }]

#     device_properties = [{
#       device_type = "DISK"
#     }]

#     disk_size_mib = 5000
#   }]
# }

resource "nutanix_virtual_machine" "demo-01-app" {
  metadata {
    kind = "vm"
  }

  name                 = "demo-01-app"
  description          = "Demo Java middleware App server"
  num_vcpus_per_socket = 2
  num_sockets          = 1
  memory_size_mib      = 8192
  power_state          = "ON"

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.clusterid}"
  }

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-lamp-subnet.id}"
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
      uuid = "${nutanix_image.centos-lamp-app.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}

# resource "nutanix_virtual_machine" "demo-01-db" {
#   metadata {
#     kind = "vm"
#   }


#   name                 = "demo-01-db"
#   description          = "demo MySQL Database Server"
#   num_vcpus_per_socket = 4
#   num_sockets          = 1
#   memory_size_mib      = 16384
#   power_state          = "ON"


#   cluster_reference = {
#     kind = "cluster"
#     uuid = "${var.clusterid}"
#   }


#   nic_list = [{
#     subnet_reference = {
#       kind = "subnet"
#       uuid = "${nutanix_subnet.next-lamp-subnet.id}"
#     }


#     ip_endpoint_list = {
#       ip   = "192.168.0.10"
#       type = "ASSIGNED"
#     }
#   }]


#   disk_list = [{
#     data_source_reference = [{
#       kind = "image"
#       name = "Centos7"
#       uuid = "${nutanix_image.centos-lamp-db.id}"
#     }]


#     device_properties = [{
#       device_type = "DISK"
#     }]


#     disk_size_mib = 5000
#   }]
# }

