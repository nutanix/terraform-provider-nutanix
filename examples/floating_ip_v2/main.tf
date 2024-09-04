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
resource "nutanix_floating_ip_v2" "fip1" {
  name = "example-fip"
  description = "example fip  description"
  external_subnet_reference = "{{ext_sub_uuid}}"
}


# create Floating IP with vpc UUID with external subnet uuid

resource "nutanix_floating_ip_v2" "fip2" {
    name = "example-fip"
    description = "example fip  description"
    external_subnet_reference_uuid = "{{ext_sub_uuid}}"
    vpc_reference_uuid= "{{vpc_uuid}}"
    association{
      private_ip_association{
        vpc_reference = "{{vpc_uuid}}"
        private_ip{
          ipv4{
            value = "10.44.44.7"
          }
        }
      }
    }
}

# create Floating IP with External Subnet with vm

resource "nutanix_floating_ip" "fip3" {
    name = "example-fip"
    description = "example fip  description"
    external_subnet_reference_uuid = "{{ext_sub_uuid}}"
    association{
    vm_nic_association{
      vm_nic_reference = "{{vm_nic_uuid}}"
    }
  }
}

# data source floating IP

data "nutanix_floating_ip_v2" "fip4"{
    floating_ip_uuid = "{{floating_ip_uuid}}"
}

# list of floating IPs

data "nutanix_floating_ips_v2" "fip5"{ }

output "csf1" {
  value = data.nutanix_floating_ips_v2.fip5
}



data "nutanix_floating_ips_v2" "fip6"{
	metadata{
		filter = "name eq 'example-fip'"
	}
}

output "csf2" {
  value = data.nutanix_floating_ips_v2.fip6
}