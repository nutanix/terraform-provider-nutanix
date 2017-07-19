provider "nutanix" {
    username = ""
    password = ""
    endpoint = "10.5.68.6"
    insecure =  true
}

resource "nutanix_virtual_machine" "my-machine" {
    name = "kritagya_testupdate1"
    spec {
        name = "kritagya_newvm"
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 1
            memory_size_mb = 2048
            power_state = "On"
            nic_list = [
                {
                    nic_type = "NORMAL_NIC"
                    subnet_reference = {
                        kind = "subnet"
                        uuid = "c03ecf8f-aa1c-4a07-af43-9f2f198713c0"
                    }
                    network_function_nic_type = "INGRESS"
                }
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
                }
            ]
        }
    }
    metadata = {
        categories = {
            "Project" =  "nucalm"
        }
    }
    provisioner "remote-exec"{
        inline = [
        ]
        connection {
            type = "ssh"
            user = "root"
            password = "nutanix/4u"
            host = "${nutanix_virtual_machine.my-machine.ip_address}"
        }
    }
}

output "ip" {
    value = "${nutanix_virtual_machine.my-machine.ip_address}"
}
