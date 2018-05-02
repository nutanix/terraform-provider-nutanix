provider "nutanix" {
    username = "jon"
    password = "Nutanix/1234"
    endpoint = "10.5.80.30"
    insecure = true
}

resource "nutanix_image" "centos73-minimal-iso" {
    name = "centos73-minimal-iso"
    source_uri = "http://earth.corp.nutanix.com/isos/linux/centos/7/CentOS-7.3-x86_64-Minimal-1611.iso"
    description = "here is my centos73 image from earth filer"
}

resource "nutanix_image" "nutanix-virtio-111-iso" {
    name = "nutanix-virtio-111-iso"
    source_uri = "http://endor.dyn.nutanix.com/GoldImages/virtio/1.1.1/Nutanix-VirtIO-1.1.1.iso"
    description = "here is my Nutanix-VirtIO-1.1.1.iso image"
}

resource "nutanix_image" "windows2016-iso" {
    name = "windows2016-iso"
    source_uri = "http://earth.corp.nutanix.com/isos/microsoft/server/2016/en_windows_server_2016_x64_dvd_9327751.iso"
    description = "heres a windows iso"
}

resource "nutanix_image" "cirros-034-disk" {
    name = "cirros-034-disk"
    source_uri = "http://endor.dyn.nutanix.com/acro_images/DISKs/cirros-0.3.4-x86_64-disk.img"
    description = "heres a tiny linux image, not an iso, but a real disk!"
}

resource "nutanix_virtual_machine" "tf-cirros" {
    name = "tf-cirros"
    spec {
        description = "Beep Boop I run cirros"
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 2
            memory_size_mib = 2048
            power_state = "ON"
            nic_list = [
                {
                    subnet_reference = {
                        kind = "subnet"
                        uuid = "bf1168dd-9355-4dc2-b3eb-18c65615bcba"
                    }
                }
            ]
            disk_list = [
                {
                    data_source_reference = {
                        kind = "image"
                        uuid = "${nutanix_image.cirros-034-disk.id}"
                    }
                }
            ]
        }
    }
}

resource "nutanix_virtual_machine" "tf-windows" {
    name = "tf-windows"
    spec {
        description = "Beep Boop I run windows 2016"
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 2
            memory_size_mib = 2048
            power_state = "ON"
            nic_list = [
                {
                    subnet_reference = {
                        kind = "subnet"
                        uuid = "bf1168dd-9355-4dc2-b3eb-18c65615bcba"
                    }
                }
            ]
            disk_list = [
                {
                    data_source_reference = {
                        kind = "image"
                        uuid = "${nutanix_image.windows2016-iso.id}"
                    }
                },
                {
                    data_source_reference = {
                        kind = "image"
                        uuid = "${nutanix_image.nutanix-virtio-111-iso.id}"
                    }
                },
                {
                    disk_size_mib = 50000
                }
            ]
        }
    }
}

resource "nutanix_virtual_machine" "tf-centos" {
    name = "tf-centos"
    spec {
        description = "Beep Boop I run centos73"
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 2
            memory_size_mib = 2048
            power_state = "ON"
            nic_list = [
                {
                    subnet_reference = {
                        kind = "subnet"
                        uuid = "bf1168dd-9355-4dc2-b3eb-18c65615bcba"
                    }
                }
            ]
            disk_list = [
                {
                    data_source_reference = {
                        kind = "image"
                        uuid = "${nutanix_image.centos73-minimal-iso.id}"
                    }
                },
                {
                    disk_size_mib = 50000
                }
            ]
        }
    }
}