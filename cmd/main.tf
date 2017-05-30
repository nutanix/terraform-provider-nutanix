provider "example" {
    username = "admin"
    password = "Nutanix.1"
    endpoint = "https://10.5.68.6:9440/api/nutanix/v3/vms"
}

resource "example_server" "my-server" {
    name = "kritagya_test1"
    Spec_Resources_NumVCPUsPerSocket = 1
    Spec_Resources_NumSockets = 1
    Spec_Resources_MemorySizeMib = 1024
    Spec_Resources_PowerState = "POWERED_ON"
    APIversion = "3.0"
}
