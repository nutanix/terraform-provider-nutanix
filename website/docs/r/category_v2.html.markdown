---
layout: "nutanix"
page_title: "NUTANIX: nutanix_category_v2"
sidebar_current: "docs-nutanix-datasource-category-v2"
description: |-
  Create, Update and Delete category (key - value pair).

# nutanix_category_v2
Create, Update and Delete category (key - value pair).


## Example

```hcl

resource "nutanix_category_v2" "example" {
  key         = "category_example_key"
  value       = "category_example_value"
  description = "category example description"
}

```


## Argument Reference

The following arguments are supported:

* `key`: -(Required) The key of a category when it is represented in key:value format. Constraints applicable when field is given in the payload during create and update:

  - A string of maxlength of 64
  - Character at the start cannot be `$`
  - Character `/` is not allowed anywhere.

  It is a mandatory field in the payload of `createCategory` and `updateCategoryById` APIs.
  This field can't be updated through `updateCategoryById` API.
* `value`: -(Required) The value of a category when it is represented in key:value format.  Constraints applicable when field is given in the payload during create and update:

  - A string of maxlength of 64
  - Character at the start cannot be `$`
  - Character `/` is not allowed anywhere.

  It is a mandatory field in the payload of `createCategory` and `updateCategoryById` APIs.
  This field can't be updated through `updateCategoryById` API.
  Updating the value will not change the extId of the category.
* `type`: -(Required) Denotes the type of a category.
    Valid values are:
     - `SYSTEM` Predefined categories contained in the system to be used by workflows visible in the UI that involve categories. System-defined categories can't be created through the Categories API. They are predefined in a configuration file and are created at PC boot-up time. System-defined categories can't be updated or deleted.
     - `INTERNAL` Predefined categories contained in the system to be used by internal services, APIs and workflows that involve categories. These categories will not be visible in the UI. However, these categories will be returned in the response of `listCategories` and `getCategoryById` APIs, and are available for filtering as well. Internal categories can't be created through the Categories API. They are predefined in a configuration file and are created at PC boot-up time. Internal categories can't be updated or deleted.
     - `USER` These categories get created by users through the invocation of `createCategory` API. User-defined categories can be updated or deleted after creation.

* `description`: -(Optional) A string consisting of the description of the category as defined by the user.
Description can be optionally provided in the payload of createCategory and updateCategoryById APIs.
Description field can be updated through updateCategoryById API.
The server does not validate this value nor does it enforce the uniqueness or any other constraints.
It is the responsibility of the user to ensure that any semantic or syntactic constraints are retained when mutating this field.

* `owner_uuid`: -(Optional) This field contains the UUID of a user who owns the category.
This field will be ignored if given in the payload of createCategory API. Hence, when a category is created, the logged-in user automatically becomes the owner of the category.
This field can be updated through updateCategoryById API, in which case, should be provided, UUID of a valid user is present in the system.
Validity of the user UUID can be checked by invoking the API: authn/users/{extId} in the 'Identity and Access Management' or 'IAM' namespace.
It is used for enabling RBAC access to self-owned categories.

## Attributes Reference

The following attributes are exported:

* `key`: The key of a category when it is represented in key:value format.
* `value`: The value of a category when it is represented in key:value format
* `type`: Denotes the type of category.
There are three types of categories: SYSTEM, INTERNAL, and USER.
* `description`: A string consisting of the description of the category as defined by the user.
* `owner_uuid`: This field contains the UUID of a user who owns the category.
* `associations`: This field gives basic information about resources that are associated with the category.
The results present under this field summarize the counts of various kinds of resources associated with the category.
For more detailed information about the UUIDs of the resources, please look into the field detailedAssociations.
This field will be ignored, if given in the payload of updateCategoryById or createCategory APIs.
This field will not be present by default in listCategories API, unless the parameter $expand=associations is present in the URL.
* `detailed_associations`: This field gives detailed information about the resources which are associated with the category.
The results present under this field contain the UUIDs of the entities and policies of various kinds associated with the category.
This field will be ignored, if given in the payload of updateCategoryById or createCategory APIs.
This field will not be present by default in listCategories or getCategoryById APIs, unless the parameter $expand=detailedAssociations is present in the URL.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.


### associations
* `category_id`: External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.
* `resource_type`: An enum denoting the associated resource types. Resource types are further grouped into 2 types - entity or a policy.
   Valid values are:
     - `APP`:  A resource of type application.
     - `PROTECTION_RULE`:  A policy or rule of type protection rule.
     - `IMAGE_RATE_LIMIT`: A resource of type rate limit.
     - `MH_VM`: A resource of type Virtual Machine.
     - `BLUEPRINT`:  A resource of type blueprint.
     - `HOST`:  A resource representing the underlying host, the machine hosting the hypervisors and VMs.
     - `IMAGE`:  A resource of type image.
     - `VM_VM_ANTI_AFFINITY_POLICY`:  A policy of type VM-VM anti-affinity; This policy decides that the specified set of VMs are running on different hosts.
     - `ACCESS_CONTROL_POLICY`: A policy or rule of type access control policy or ACP; the rules that decide authorization of users to access an API.
     - `VM_HOST_AFFINITY_POLICY`: A policy of type VM host affinity; The policy decides the affinity between a set of VMs to be run only a specified set of hosts
     - `NGT_POLICY`:  A policy or rule of type NGT policy.
     - `RECOVERY_PLAN`:  A policy or rule of type recovery plan.
     - `MARKETPLACE_ITEM`: A resource of type marketplace item.
     - `CLUSTER`:  A resource of type cluster, usually refers to a PE cluster.
     - `NETWORK_SECURITY_RULE`:  A rule of type network security.
     - `HOST_NIC`:  A resource of type Physical NIC.
     - `ACTION_RULE`:  A policy of type Playbook.
     - `VOLUMEGROUP`:  A resource of type volume group.
     - `REPORT`:  A resource of type report.
     - `STORAGE_POLICY`: A policy or rule of type storage policy.
     - `BUNDLE`:  A resource of type bundle.
     - `QOS_POLICY`: A policy or rule of type QoS policy.
     - `SUBNET`:  A resource of type network subnets.
     - `VM`: A resource of type Virtual Machine.
     - `NETWORK_SECURITY_POLICY`:  A policy of type network security.
     - `POLICY_SCHEMA`:  Policies like user-defined-alerts.
     - `IMAGE_PLACEMENT_POLICY`:  A policy of type image placement.
* `resource_group`: An enum denoting the resource group.
Resources can be organized into either an entity or a policy.
  Valid values are:
     - `POLICY`:  A ResourceGroup denoting a nutanix policy like VM host affinity policy, image placement policy, access control policy, and so on.<br> A category is generally associated with many entities.<br> The policy which is associated with this category, is then applied to those entities which are also associated with the same category.
     - `ENTITY`: A ResourceGroup denoting a nutanix entity like VM, cluster, host, image, and so on.<br> A category is generally associated with many entities.<br> A policy is then applied to these entities through the category.
* `count`: Count of associations of a particular type of entity or policy

### detailed_associations
* `category_id`: External identifier for the given category, used across all v4 apis/entities/resources where categories are referenced.
* `resource_type`: An enum denoting the associated resource types. Resource types are further grouped into 2 types - entity or a policy.
   Valid values are:
     - `APP`:  A resource of type application.
     - `PROTECTION_RULE`:  A policy or rule of type protection rule.
     - `IMAGE_RATE_LIMIT`: A resource of type rate limit.
     - `MH_VM`: A resource of type Virtual Machine.
     - `BLUEPRINT`:  A resource of type blueprint.
     - `HOST`:  A resource representing the underlying host, the machine hosting the hypervisors and VMs.
     - `IMAGE`:  A resource of type image.
     - `VM_VM_ANTI_AFFINITY_POLICY`:  A policy of type VM-VM anti-affinity; This policy decides that the specified set of VMs are running on different hosts.
     - `ACCESS_CONTROL_POLICY`: A policy or rule of type access control policy or ACP; the rules that decide authorization of users to access an API.
     - `VM_HOST_AFFINITY_POLICY`: A policy of type VM host affinity; The policy decides the affinity between a set of VMs to be run only a specified set of hosts
     - `NGT_POLICY`:  A policy or rule of type NGT policy.
     - `RECOVERY_PLAN`:  A policy or rule of type recovery plan.
     - `MARKETPLACE_ITEM`: A resource of type marketplace item.
     - `CLUSTER`:  A resource of type cluster, usually refers to a PE cluster.
     - `NETWORK_SECURITY_RULE`:  A rule of type network security.
     - `HOST_NIC`:  A resource of type Physical NIC.
     - `ACTION_RULE`:  A policy of type Playbook.
     - `VOLUMEGROUP`:  A resource of type volume group.
     - `REPORT`:  A resource of type report.
     - `STORAGE_POLICY`: A policy or rule of type storage policy.
     - `BUNDLE`:  A resource of type bundle.
     - `QOS_POLICY`: A policy or rule of type QoS policy.
     - `SUBNET`:  A resource of type network subnets.
     - `VM`: A resource of type Virtual Machine.
     - `NETWORK_SECURITY_POLICY`:  A policy of type network security.
     - `POLICY_SCHEMA`:  Policies like user-defined-alerts.
     - `IMAGE_PLACEMENT_POLICY`:  A policy of type image placement.
* `resource_group`: An enum denoting the resource group.
Resources can be organized into either an entity or a policy.
  Valid values are:
     - `POLICY`:  A ResourceGroup denoting a nutanix policy like VM host affinity policy, image placement policy, access control policy, and so on.<br> A category is generally associated with many entities.<br> The policy which is associated with this category, is then applied to those entities which are also associated with the same category.
     - `ENTITY`: A ResourceGroup denoting a nutanix entity like VM, cluster, host, image, and so on.<br> A category is generally associated with many entities.<br> A policy is then applied to these entities through the category.
* `resource_id`: The UUID of the entity or policy associated with the particular category.

## Import
This helps to manage existing entities which are not created through terraform. Category (key - value pair) can be imported using the `UUID` (ext_id in v4 terms).  eg,

`
terraform import nutanix_category_v2.<resource_name> <UUID>
`

Note: 
We have two resources separately for category key (nutanix_category_key) and value (nutanix_category_key). Using v4 API, `nutanix_category_v2` represents category key value pair as one entity. 

Please use datasources (nutanix_categories_v2) to fetch uuids (ext_id) of all category key valye pairs to import them.


See detailed information in [Nutanix Create Category v4](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/Categories/operation/createCategory).
