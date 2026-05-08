---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_guest_customization_profile_v2"
sidebar_current: "docs-nutanix-resource-vm-guest-customization-profile-v2"
description: |-
  Creates a new VM Guest Customization profile with the provided configuration.
---

# nutanix_vm_guest_customization_profile_v2

Provides a resource to create, read, update, and delete VM Guest Customization Profiles. A VM Guest Customization Profile provides configuration for the customization of either Windows or Linux guest operating systems during VM deployment.

## Example Usage

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
          }
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

* `sysprep_config` - (Optional) Sysprep configuration for Windows guest OS customization.

### sysprep_config

* `customization` - (Optional) Either specify the values for the parameters or an unattend XML file.

### customization

* `sysprep_params` - (Optional) Sysprep parameters for Windows customization.
* `answer_file` - (Optional) Answer file configuration.

### sysprep_params

* `first_logon_commands` - (Optional) List of commands to be executed automatically when a user logs in for the first time after Windows setup. This is an ordered list.
* `general_settings` - (Optional) General settings for Windows customization.
* `locale_settings` - (Optional) Locale settings for Windows customization.
* `network_settings` - (Optional) Network settings for Windows customization.
* `workgroup_or_domain_info` - (Optional) JoinWorkgroup or JoinDomain settings of the computer.

### general_settings

* `administrator_password` - (Optional, Sensitive) Password to be configured for built-in Administrator account.
* `auto_logon_settings` - (Optional) Auto logon settings.
  * `logon_count` - (Optional) The number of automatic logons allowed for the computer using the specified local account.
* `computer_name` - (Optional) Mechanism to use to generate the computer name of the VM. Either UseVmName or MustProvideDuringDeployment should be provided.
  * `must_provide_during_deployment` - (Optional) If true, the user must provide the value for the computer name during the VM deployment.
  * `use_vm_name` - (Optional) If true, the name of the VM is used as the computer name during deployment.
* `registered_organization` - (Optional) Name of the organization of the end user.
* `registered_owner` - (Optional) Full name of the end user.
* `timezone` - (Optional) The computer's time zone in string format.
* `windows_product_key` - (Optional, Sensitive) The product key to use to install and activate Windows.

### locale_settings

* `system_locale` - (Optional) Default language to use for non-Unicode programs (e.g., en-US, fr-FR).
* `ui_language` - (Optional) Default system language to use to display user interface (UI) items (e.g., en-US, fr-FR).
* `user_locale` - (Optional) Per-user settings to be used for formatting dates, times, currency, and numbers (e.g., en-US, fr-FR).

### network_settings

* `nic_config_list` - (Optional) List of NIC configurations to be applied to the NICs attached to the VM in serial order.

### nic_config_list

* `dns_config` - (Optional) DNS configuration for the NIC.
  * `alternate_dns_server_addresses` - (Optional) List of IPv4 addresses to look for after preferred DNS server.
  * `preferred_dns_server_address` - (Optional) An IPv4 address preferred to search first when searching for the DNS server.
* `ipv4_config` - (Optional) Mechanism to configure IPv4 settings of the NIC.
  * `use_dhcp` - (Optional) If true, DhcpEnabled is set to True for the interface.
  * `must_provide_during_deployment` - (Optional) If true, the user must provide the IPv4 address during the deployment.

### workgroup_or_domain_info

* `workgroup` - (Optional) Workgroup settings.
  * `name` - (Optional) Name of workgroup to be applied to the computer when joining the workgroup.
* `domain_settings` - (Optional) Domain settings.
  * `credentials` - (Optional) Domain credentials.
    * `domain_name` - (Optional) The name of the domain to use for authentication.
    * `password` - (Optional, Sensitive) The password of the domain user account.
    * `username` - (Optional) Name of the domain user account with permission to add the computer to a domain.

### answer_file

* `unattend_xml` - (Optional) The unattend XML file as a string value.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.
* `create_time` - VM Guest Customization Profile creation time.
* `update_time` - VM Guest Customization Profile last updated time.
* `created_by` - The user who created the profile.
  * `ext_id` - The external ID (UUID) of the user.
* `updated_by` - The user who last updated the profile.
  * `ext_id` - The external ID (UUID) of the user.

## Import

VM Guest Customization Profiles can be imported using the `ext_id`:

```hcl
terraform import nutanix_vm_guest_customization_profile_v2.example <ext_id>
```

See detailed information in [Nutanix VM Guest Customization Profiles V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmGuestCustomizationProfiles)
