provider "nutanix" {
    username = "admin"
    password = "Nutanix123#"
    endpoint = "10.5.68.6"
    insecure =  true 
}

resource "nutanix_virtual_machine" "my-machine" {
    name = "kritagya_test1"
    spec {
        name = "kritagya_testvm"
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 1
            memory_size_mb = 1024
            power_state = "POWERED_ON"
            nic_list = [
                { 
                    nic_type = "NORMAL_NIC"
                    subnet_reference = {
                        kind = "subnet"
                        uuid = "c03ecf8f-aa1c-4a07-af43-9f2f198713c0"
                    }
                    network_function_nic_type = "INGRESS"
                },
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
                }
            ]
        }
    }
    metadata = {
        kind = "vm"
        spec_version = 0
        name = "kritagya1"
        categories = {
            "Project" =  "nucalm"
        }
    }
    api_version = "3.0"
    
}

output "ip" {
    value = "${nutanix_virtual_machine.my-machine.ip_address}"
}
