terraform{
    required_providers {
        nutanix = {
            source = "nutanix/nutanix"
            version = "1.6.0"
        }
    }
}

#definig nutanix configuration
provider "nutanix"{
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port = 9440
  insecure = true
}


# create Floating IP with External Subnet UUID

resource "nutanix_floating_ip" "fip1" {
    external_subnet_reference_uuid = "{{ext_sub_uuid}}"
}

# create Floating IP with vpc UUID with external subnet uuid

resource "nutanix_floating_ip" "fip2" {
    external_subnet_reference_uuid = "{{ext_sub_uuid}}"
    vpc_reference_uuid= "{{vpc_uuid}}"
    private_ip = "{{ip_address}}"
}

# create Floating IP with External Subnet with vm

resource "nutanix_floating_ip" "fip3" {
    external_subnet_reference_uuid = "{{ext_sub_uuid}}"
    vm_nic_reference_uuid = "{{vm_uuid}}"
}