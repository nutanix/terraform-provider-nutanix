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
data "nutanix_vm_guest_customization_profile_v2" "profile" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}

output "profile_name" {
  value = data.nutanix_vm_guest_customization_profile_v2.profile.name
}
```

## Argument Reference

The following arguments are supported:

* `ext_id` - (Required) A globally unique identifier of a VM Guest Customization Profile in UUID format.

## Attribute Reference

The following attributes are exported:

* `name` - Name of the VM Guest Customization Profile.
* `description` - VM Guest Customization Profile description.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
* `links` - A HATEOAS style link for the response.
* `config` - Configuration of the VM Guest Customization Profile.
* `create_time` - VM Guest Customization Profile creation time.
* `update_time` - VM Guest Customization Profile last updated time.
* `created_by` - The user who created the profile.
  * `ext_id` - The external ID (UUID) of the user.
* `updated_by` - The user who last updated the profile.
  * `ext_id` - The external ID (UUID) of the user.

### config

* `sysprep_config` - Sysprep configuration for Windows guest OS customization.
  * `customization` - Either specify the values for the parameters or an unattend XML file.
    * `sysprep_params` - Sysprep parameters for Windows customization.
      * `first_logon_commands` - List of commands to be executed automatically when a user logs in for the first time after Windows setup.
      * `general_settings` - General settings for Windows customization.
        * `administrator_password` - Password to be configured for built-in Administrator account.
        * `auto_logon_settings` - Auto logon settings.
          * `logon_count` - The number of automatic logons allowed for the computer.
        * `computer_name` - Mechanism to use to generate the computer name of the VM.
          * `must_provide_during_deployment` - Whether the user must provide the computer name during deployment.
          * `use_vm_name` - Whether to use the VM name as the computer name.
        * `registered_organization` - Name of the organization of the end user.
        * `registered_owner` - Full name of the end user.
        * `timezone` - The computer's time zone in string format.
        * `windows_product_key` - The product key to use to install and activate Windows.
      * `locale_settings` - Locale settings.
        * `system_locale` - Default language to use for non-Unicode programs.
        * `ui_language` - Default system language to use to display user interface (UI) items.
        * `user_locale` - Per-user settings to be used for formatting dates, times, currency, and numbers.
      * `network_settings` - Network settings.
        * `nic_config_list` - List of NIC configurations.
          * `dns_config` - DNS configuration.
            * `alternate_dns_server_addresses` - List of IPv4 addresses to look for after preferred DNS server.
            * `preferred_dns_server_address` - An IPv4 address preferred to search first.
          * `ipv4_config` - IPv4 configuration.
            * `use_dhcp` - Whether to use DHCP.
            * `must_provide_during_deployment` - Whether the user must provide the IPv4 address during deployment.
      * `workgroup_or_domain_info` - JoinWorkgroup or JoinDomain settings.
        * `workgroup` - Workgroup settings.
          * `name` - Name of workgroup.
        * `domain_settings` - Domain settings.
          * `credentials` - Domain credentials.
            * `domain_name` - The name of the domain.
            * `password` - The password of the domain user account.
            * `username` - Name of the domain user account.
    * `answer_file` - Answer file configuration.
      * `unattend_xml` - The unattend XML file as a string value.

See detailed information in [Nutanix Get VM Guest Customization Profile V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmGuestCustomizationProfiles/operation/getVmGuestCustomizationProfileById)
