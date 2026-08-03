---
layout: "nutanix"
page_title: "NUTANIX: nutanix_vm_startup_policy_vm_compliance_states_v2"
sidebar_current: "docs-nutanix-datasource-vm-startup-policy-vm-compliance-states-v2"
description: |-
  List VM compliances of the provided VM startup policy external identifier.
---

# nutanix_vm_startup_policy_vm_compliance_states_v2

List VM compliances of the provided VM startup policy external identifier.

## Example Usage

```hcl
data "nutanix_vm_startup_policy_vm_compliance_states_v2" "example" {
  vm_startup_policy_ext_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

* `vm_startup_policy_ext_id` - (Required) The external ID of the VM startup policy.

## Attributes Reference

* `vm_compliance_states` - List of VM compliance states.
  * `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
  * `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.
  * `associated_categories` - The categories through which the VM is associated with the policy.
    * `ext_id` - The globally unique identifier of an instance of type UUID.
  * `cluster` - The cluster reference.
    * `ext_id` - The globally unique identifier of an instance of type UUID.
  * `compliance_status` - The compliance status of the VM. Possible values: `COMPLIANT`, `NON_COMPLIANT`, `PENDING`.
  * `non_compliance_reason` - The reason why the VM is non-compliant. Only populated when `compliance_status` is `NON_COMPLIANT`. Possible values: `NGT_NOT_ENABLED`, `HA_NOT_SUPPORTED`, `CLUSTER_NOT_SUPPORTED`.
  * `links` - A HATEOAS style link for the response.
    * `href` - The URL at which the entity described by the link can be accessed.
    * `rel` - A name that identifies the relationship of the link to the object that is returned by the URL.

See detailed information in [Nutanix VM Startup Policy V2](https://developers.nutanix.com/).
