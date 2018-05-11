provider "nutanix" {
    username = "jon"
    password = "superSecretStuff/1234"
    endpoint = "10.5.80.30"
    insecure = true
}

resource "nutanix_virtual_machine" "ThisOldCloud-TF-Windows" {
    name = "ThisOldCloud-TF-Windows"
    spec {
        description = "Beep Boop I'm a VM"
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
                        uuid = "4cf6d903-6e91-46a4-90b2-4d0c0ba3955f"
                    }
                }
            ]

        }
    }

    metadata = {
        kind = "vm"
    }
}

output "name" {
    value = "${nutanix_virtual_machine.ThisOldCloud-TF-Windows.name}"
}

output "ip" {
    value = "${nutanix_virtual_machine.ThisOldCloud-TF-Windows.ip_address}"
}

output "UUID" {
    value = "${nutanix_virtual_machine.ThisOldCloud-TF-Windows.id}"
}
