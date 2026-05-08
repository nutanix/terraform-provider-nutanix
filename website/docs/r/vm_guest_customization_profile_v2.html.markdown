---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_guest_customization_profile_v2"
sidebar_current: "docs-nutanix-resource-vm-guest-customization-profile-v2"
description: |-
  Provides a Nutanix VM Guest Customization Profile resource to create and manage guest customization profiles.
---

# nutanix_vm_guest_customization_profile_v2

Provides a resource to create, read, update, and delete VM Guest Customization Profiles. A VM Guest Customization Profile provides configuration for the customization of either Windows or Linux guest operating systems.

## Example Usage

### Sysprep Params Configuration

```hcl
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
```

### Answer File Configuration

```hcl
resource "nutanix_vm_guest_customization_profile_v2" "answer_file" {
  name        = "example-gc-profile-answer-file"
  description = "Example with answer file"
  config {
    sysprep_config {
      customization {
        answer_file {
          unattend_xml = "<unattend xmlns='urn:schemas-microsoft-com:unattend'></unattend>"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the VM Guest Customization Profile.
* `description` - (Optional) VM Guest Customization Profile description.
* `config` - (Optional) Configuration of the VM Guest Customization Profile. A Configuration should be provided for the customization of either Windows or Linux guest operating system.

### config

* `sysprep_config` - (Optional) Configuration for customization of a Windows guest operating system.

### sysprep_config

* `customization` - (Optional) Either specify the values for the parameters or an unattend XML file.

### customization

* `sysprep_params` - (Optional) A set of various unattended settings supported by Windows.
* `answer_file` - (Optional) The answer file (unattend.xml) for sysprep.

### sysprep_params

* `first_logon_commands` - (Optional) List of commands to be executed automatically when a user logs in for the first time after Windows setup.
* `general_settings` - (Optional) A set of general unattended settings supported by Windows.
* `locale_settings` - (Optional) Language and locale settings for the system and the user.
* `network_settings` - (Optional) Network settings to apply to the NICs attached to the VM.
* `workgroup_or_domain_info` - (Optional) JoinWorkgroup or JoinDomain settings of the computer.

### general_settings

* `administrator_password` - (Optional, Sensitive) Password to be configured for built-in Administrator account.
* `auto_logon_settings` - (Optional) Autologon settings.
* `computer_name` - (Optional) Mechanism to use to generate the computer name of the VM. Either `use_vm_name` or `must_provide_during_deployment` should be set to `true`.
* `registered_organization` - (Optional) Name of the organization of the end user.
* `registered_owner` - (Optional) Full name of the end user.
* `timezone` - (Optional) The computer's time zone in string format.
* `windows_product_key` - (Optional, Sensitive) The product key to use to install and activate Windows.

### auto_logon_settings

* `logon_count` - (Required) The number of automatic logons allowed for the computer.

### computer_name

* `use_vm_name` - (Optional) Whether to use the VM name as the computer name.
* `must_provide_during_deployment` - (Optional) Whether the computer name must be provided during deployment.

### locale_settings

* `system_locale` - (Optional) Default language to use for non-Unicode programs.
* `ui_language` - (Optional) Default system language to use to display UI items.
* `user_locale` - (Optional) Per-user settings for formatting dates, times, currency, and numbers.

### network_settings

* `nic_config_list` - (Required) List of NIC configurations to be applied to the NICs attached to the VM in serial order.

### nic_config_list

* `dns_config` - (Optional) DNS configuration to be applied to the NIC.
* `ipv4_config` - (Required) Mechanism to configure IPv4 settings of the NIC.

### dns_config

* `alternate_dns_server_addresses` - (Optional) List of IPv4 addresses to look for after preferred DNS server.
* `preferred_dns_server_address` - (Required) An IPv4 address preferred to search first for the DNS server.

### ipv4_config

* `use_dhcp` - (Optional) Whether to use DHCP for IPv4 configuration.
* `must_provide_during_deployment` - (Optional) Whether IPv4 configuration must be provided during deployment.

### workgroup_or_domain_info

* `workgroup` - (Optional) Workgroup settings.
* `domain_settings` - (Optional) Domain settings.

### workgroup

* `name` - (Required) Name of workgroup to be applied to the computer. It must be a valid NetBIOS name.

### domain_settings

* `credentials` - (Required) Credentials for the domain account.

### credentials

* `domain_name` - (Required) The name of the domain.
* `password` - (Required, Sensitive) The password of the domain user account.
* `username` - (Required) Name of the domain user account with permission to add the computer to a domain.

### answer_file

* `unattend_xml` - (Required) The unattend XML file as a string value. Note that double quotes in the XML file need to be escaped.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `create_time` - VM Guest Customization Profile creation time.
* `update_time` - VM Guest Customization Profile last updated time.
* `created_by` - The user who created the profile.
* `updated_by` - The user who last updated the profile.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.

## Import

VM Guest Customization Profiles can be imported using the `ext_id`.

```hcl
terraform import nutanix_vm_guest_customization_profile_v2.example <ext_id>
```

See detailed information in [Nutanix VM Guest Customization Profiles V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmGuestCustomizationProfiles)
