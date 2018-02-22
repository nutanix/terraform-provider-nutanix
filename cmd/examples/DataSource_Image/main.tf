provider "nutanix" { 
    username = "admin"
    password = "Nutanix/1234"
    endpoint = "10.5.80.30"
    insecure = true
}

data "nutanix_image" "my-image" {
    name = "Sarath_Centos7"
}

resource "nutanix_virtual_machine" "my-machine" {
    name = "Sarath_Test_Terraform"
    spec {
        name = "Sarath_Test_Terraform"
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 1
            memory_size_mib = 2048
            power_state = "On"
            nic_list = [
                {
                    nic_type = "NORMAL_NIC"
                    subnet_reference = {
                        kind = "subnet"
                        uuid = "d4ff3f77-c70b-43f5-b7af-f5a62e37014d"
                    }
                    network_function_nic_type = "INGRESS"
                }
            ]
            disk_list = [
                {
                    data_source_reference = {
                        kind = "image"
                        uuid = "${data.nutanix_image.my-image.uuid}"
                    }
                    device_properties = {
                        device_type = "DISK"
                    }
                    disk_size_mib = 1
                }
            ]
        }
    }
    provisioner "remote-exec"{
        inline = [
	"ip addr"
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
