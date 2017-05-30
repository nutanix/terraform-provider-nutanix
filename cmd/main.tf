provider "example" {
    username = "admin"
    password = "Nutanix.1"
    endpoint = "https://10.5.68.6:9440/api/nutanix/v3/vms"
}

resource "example_server" "my-server" {
    name = "kritagya_test1"
    Spec {
        Resources = {
            NumVCPUsPerSocket = 1
            NumSockets = 1
            MemorySizeMib = 1024
            PowerState = "POWERED_ON"
        }
    }
    APIversion = "3.0"
}
