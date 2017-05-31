provider "example" {
    username = "admin"
    password = "Nutanix123#"
    endpoint = "https://10.5.68.6:9440/api/nutanix/v3/vms"
}

resource "example_server" "my-server" {
    name = "kritagya_test1"
    spec {
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 1
            memory_size_mb = 1024
            power_state = "POWERED_ON"
  
        }
    }
    metadata = {
        kind = "vm"
        spec_version = 0
        name = "kritagya1"
    }
    api_version = "3.0"
}
