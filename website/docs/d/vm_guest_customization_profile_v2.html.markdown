---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_guest_customization_profile_v2"
sidebar_current: "docs-nutanix-datasource-vm-guest-customization-profile-v2"
description: |-
  Retrieves the VM Guest Customization Profile configuration of the provided VM Guest Customization Profile external identifier.
---

# nutanix_vm_guest_customization_profile_v2

Retrieves the VM Guest Customization Profile configuration of the provided VM Guest Customization Profile external identifier.

## Example

```hcl
data "nutanix_vm_guest_customization_profile_v2" "example" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) A globally unique identifier of a VM Guest Customization Profile in UUID format.

## Attribute Reference

The following attributes are exported:

* `name` - Name of the VM Guest Customization Profile.
* `description` - VM Guest Customization Profile description.
* `config` - Configuration of the VM Guest Customization Profile.
* `create_time` - VM Guest Customization Profile creation time.
* `update_time` - VM Guest Customization Profile last updated time.
* `created_by` - The user who created the profile.
* `updated_by` - The user who last updated the profile.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.

### config

* `sysprep_config` - Configuration for customization of a Windows guest operating system.

### sysprep_config

* `customization` - Either specify the values for the parameters or an unattend XML file.

### customization

* `sysprep_params` - A set of various unattended settings supported by Windows.
* `answer_file` - The answer file (unattend.xml) for sysprep.

### sysprep_params

* `first_logon_commands` - List of commands to be executed automatically when a user logs in for the first time after Windows setup.
* `general_settings` - A set of general unattended settings supported by Windows.
* `locale_settings` - Language and locale settings for the system and the user.
* `network_settings` - Network settings to apply to the NICs attached to the VM.
* `workgroup_or_domain_info` - JoinWorkgroup or JoinDomain settings of the computer.

### general_settings

* `administrator_password` - Password to be configured for built-in Administrator account.
* `auto_logon_settings` - Autologon settings.
* `computer_name` - Mechanism to use to generate the computer name of the VM.
* `registered_organization` - Name of the organization of the end user.
* `registered_owner` - Full name of the end user.
* `timezone` - The computer's time zone in string format.
* `windows_product_key` - The product key to use to install and activate Windows.

### auto_logon_settings

* `logon_count` - The number of automatic logons allowed for the computer.

### computer_name

* `use_vm_name` - Whether to use the VM name as the computer name.
* `must_provide_during_deployment` - Whether the computer name must be provided during deployment.

### locale_settings

* `system_locale` - Default language to use for non-Unicode programs.
* `ui_language` - Default system language to use to display UI items.
* `user_locale` - Per-user settings for formatting dates, times, currency, and numbers.

### network_settings

* `nic_config_list` - List of NIC configurations to be applied to the NICs attached to the VM.

### nic_config_list

* `dns_config` - DNS configuration to be applied to the NIC.
* `ipv4_config` - Mechanism to configure IPv4 settings of the NIC.

### dns_config

* `alternate_dns_server_addresses` - List of IPv4 addresses to look for after preferred DNS server.
* `preferred_dns_server_address` - An IPv4 address preferred to search first for the DNS server.

### ipv4_config

* `use_dhcp` - Whether to use DHCP for IPv4 configuration.
* `must_provide_during_deployment` - Whether IPv4 configuration must be provided during deployment.

### workgroup_or_domain_info

* `workgroup` - Workgroup settings.
* `domain_settings` - Domain settings.

### workgroup

* `name` - Name of workgroup to be applied to the computer.

### domain_settings

* `credentials` - Credentials for the domain account.

### credentials

* `domain_name` - The name of the domain.
* `password` - The password of the domain user account.
* `username` - Name of the domain user account.

### answer_file

* `unattend_xml` - The unattend XML file as a string value.

See detailed information in [Nutanix VM Guest Customization Profiles V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmGuestCustomizationProfiles/operation/getVmGuestCustomizationProfileById)
