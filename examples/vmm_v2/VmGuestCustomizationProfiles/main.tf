terraform {
  required_providers {
    nutanix = {
      source  = "nutanix/nutanix"
      version = "2.0.0"
    }
  }
}

#defining nutanix configuration
provider "nutanix" {
  username = var.nutanix_username
  password = var.nutanix_password
  endpoint = var.nutanix_endpoint
  port     = 9440
  insecure = true
}

resource "nutanix_vm_guest_customization_profile_v2" "example" {
  name        = "example-gc-profile"
  description = "Example VM Guest Customization Profile"
  config {
    sysprep_config {
      customization {
        sysprep_params {
          general_settings {
            computer_name {
              use_vm_name = true
            }
            timezone = "Pacific Standard Time"
          }
          locale_settings {
            ui_language   = "en-US"
            system_locale = "en-US"
            user_locale   = "en-US"
          }
          network_settings {
            nic_config_list {
              ipv4_config {
                use_dhcp = true
              }
            }
          }
          workgroup_or_domain_info {
            workgroup {
              name = "WORKGROUP"
            }
          }
        }
      }
    }
  }
}

data "nutanix_vm_guest_customization_profile_v2" "example" {
  ext_id = nutanix_vm_guest_customization_profile_v2.example.id
}

data "nutanix_vm_guest_customization_profiles_v2" "example" {
  depends_on = [nutanix_vm_guest_customization_profile_v2.example]
}
