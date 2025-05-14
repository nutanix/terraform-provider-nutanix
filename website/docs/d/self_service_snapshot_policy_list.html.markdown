---
layout: "nutanix"
page_title: "NUTANIX: nutanix_self_service_snapshot_policy_list"
sidebar_current: "docs-nutanix_self_service_app"
description: |-
  Describes snapshot policies present in a blueprint.
---

# nutanix_self_service_snapshot_policy_list

Describes snapshot policies present in a blueprint. Environment containing snapshot policy must be added in blueprint for this data source to work.

## Example Usage

```hcl
data "nutanix_self_service_snapshot_policy_list" "test" {
    bp_name = "NAME OF BLUEPRINT"
    length = 250
    offset = 0
}
```

## Argument Reference

The following arguments are supported:

* `bp_name`: - (Optional) The UUID of the blueprint for which snapshot policies should be listed.
* `bp_uuid`: - (Optional) The name of the blueprint for which snapshot policies should be listed.

Both (`bp_name` and `bp_uuid`) are optional but atleast one of them to be provided for this data source to work.

## Attribute Reference

The following attributes are exported:

* `length`: - (Required) The number of snapshot policy records to return.
* `offset`: - (Required) The index of the first snapshot policy to return Used for pagination. Default value: 0

### policy_list

The policy_list block contains a list of snapshot policies. Each item in the list includes the following attributes:

* `policy_name`: -  The name of the snapshot policy.
* `policy_uuid`: - The UUID of the snapshot policy.
* `policy_expiry_days`: -  The number of days after which the snapshot policy expires.
* `snapshot_config_name`: - The name of the associated snapshot configuration.
* `snapshot_config_uuid`: - The UUID of the associated snapshot configuration.


See detailed information in [List App Protection Policy](https://www.nutanix.dev/api_reference/apis/self-service.html#tag/AppProtectionPolicies/paths/~1blueprints~1%7Buuid%7D~1app_profile~1%7Bapp_profile_uuid%7D~1config_spec~1%7Bconfig_uuid%7D~1app_protection_policies~1list/post).
