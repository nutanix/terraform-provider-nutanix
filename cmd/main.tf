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
