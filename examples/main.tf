provider "nutanix" {
  username = "admin"
  password = "Nutanix/1234"
  endpoint = "10.5.81.134"
  insecure = true
  port     = 9440
}

variable cluster1 {
  default = "000567f3-1921-c722-471d-0cc47ac31055"
}

variable ip_haproxy {
  default = "1.2.3.5"
}

variable ip_app {
  default = "1.2.3.6"
}

variable ip_db {
  default = "1.2.3.7"
}

resource "nutanix_image" "centos-lamp-app" {
  name        = "CentOS-LAMP-APP.qcow2"
  description = "CentOS LAMP - App"
  source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-APP.qcow2"

  metadata = {
    kind = "image"
  }
}

resource "nutanix_image" "centos-lamp-db" {
  name        = "CentOS-LAMP-DB.qcow2"
  description = "CentOS LAMP - DB"
  source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-DB.qcow2"

  metadata = {
    kind = "image"
  }
}

resource "nutanix_image" "centos-lamp-haproxy" {
  name        = "CentOS-LAMP-HAPROXY.qcow2"
  description = "CentOS LAMP - HAProxy"
  source_uri  = "http://filer.dev.eng.nutanix.com:8080/GoldImages/NuCalm/AHV-UVM-Images/CentOS-LAMP-HAProxy.qcow2"

  metadata = {
    kind = "image"
  }
}

resource "nutanix_subnet" "next-lamp-subnet" {
  metadata = {
    kind = "subnet"
  }

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.cluster1}"
  }

  name        = "next-lamp"
  description = "lamp lamp lampy lamp vlan 0"
  vlan_id     = 202
  subnet_type = "VLAN"

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]

  dhcp_domain_search_list = ["nutanix.com", "eng.nutanix.com"]
}

resource "nutanix_virtual_machine" "demo-01-web" {
  metadata {
    kind = "vm"
  }

  name                 = "demo-01-web"
  description          = "demo Frontend Web Server"
  num_vcpus_per_socket = 2
  num_sockets          = 1
  memory_size_mib      = 4096
  power_state          = "ON"

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.cluster1}"
  }

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-lamp-subnet.id}"
    }

    ip_endpoint_list = {
      ip   = "${var.ip_haproxy}"
      type = "ASSIGNED"
    }
  }]

  disk_list = [{
    data_source_reference = [{
      kind = "image"
      name = "Centos7"
      uuid = "${nutanix_image.centos-lamp-haproxy.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}

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
    uuid = "${var.cluster1}"
  }

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-lamp-subnet.id}"
    }

    ip_endpoint_list = {
      ip   = "${var.ip_app}"
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

resource "nutanix_virtual_machine" "demo-01-db" {
  metadata {
    kind = "vm"
  }

  name                 = "demo-01-db"
  description          = "demo MySQL Database Server"
  num_vcpus_per_socket = 4
  num_sockets          = 1
  memory_size_mib      = 16384
  power_state          = "ON"

  cluster_reference = {
    kind = "cluster"
    uuid = "${var.cluster1}"
  }

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.next-lamp-subnet.id}"
    }

    ip_endpoint_list = {
      ip   = "${var.ip_db}"
      type = "ASSIGNED"
    }
  }]

  disk_list = [{
    data_source_reference = [{
      kind = "image"
      name = "Centos7"
      uuid = "${nutanix_image.centos-lamp-db.id}"
    }]

    device_properties = [{
      device_type = "DISK"
    }]

    disk_size_mib = 5000
  }]
}
