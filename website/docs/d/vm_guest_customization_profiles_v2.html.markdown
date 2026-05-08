---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_guest_customization_profiles_v2"
sidebar_current: "docs-nutanix-datasource-vm-guest-customization-profiles-v2"
description: |-
  Lists VM Guest Customization Profiles.
---

# nutanix_vm_guest_customization_profiles_v2

Lists VM Guest Customization Profiles.

## Example

```hcl
data "nutanix_vm_guest_customization_profiles_v2" "profiles" {}

data "nutanix_vm_guest_customization_profiles_v2" "filtered_profiles" {
  limit  = 10
  filter = "name eq 'my-profile'"
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) A URL query parameter that specifies the page number of the result set.
* `limit` - (Optional) A URL query parameter that specifies the total number of records returned in the result set.
* `filter` - (Optional) A URL query parameter that allows clients to filter a collection of resources.
* `order_by` - (Optional) A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.
* `select` - (Optional) A URL query parameter that allows clients to request a specific set of properties for each entity.

## Attribute Reference

The following attributes are exported:

* `vm_guest_customization_profiles` - List of VM Guest Customization Profiles. Each entry has the same attributes as the `nutanix_vm_guest_customization_profile_v2` datasource.

See detailed information in [Nutanix List VM Guest Customization Profiles V4](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.2#tag/VmGuestCustomizationProfiles/operation/listVmGuestCustomizationProfiles)
