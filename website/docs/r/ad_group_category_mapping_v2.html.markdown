---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ad_group_category_mapping_v2"
sidebar_current: "docs-nutanix-resource-ad-group-category-mapping-v2"
description: |-
  Creates the mapping between a group in Active Directory and the Category.
---

# nutanix_ad_group_category_mapping_v2

Creates and manages the mapping between a group in Active Directory and a Prism Central category. Flow Network Security uses these mappings to dynamically tag VDI machines with the category of the logged-in user's Active Directory group.

~> **Note:** A directory server must be configured before creating category mappings. Use `depends_on` to enforce this ordering. Directory Server can configured using resource `nutanix_directory_server_config_v2`

## Example Usage

### Fetch the objectGUID and create a mapping

```hcl
# Map an Active Directory group to a Prism Central category.
# A Directory Server Config must be configured before creating category mappings.
resource "nutanix_ad_group_category_mapping_v2" "example" {
  name           = "developer_group_infra_team"
  # Recommendation: use the name of the Active Directory group as the Category Mapping name.

  category_value = "infra_team"
  # This is the value of the category that will be mapped to the Active Directory group.

  category_name  = "ADGroup"
  # By default it is ADGroup.
  # Recommendation: use the same key for all AD group category mappings.

  ad_info {
    directory_service_reference = "<Directory Service ext_id>
    object_identifier           = "<User Group objectGUID>"
  }
  # object_identifier is the objectGUID for the Active Directory group. Cannot be updated once created.

  depends_on = [nutanix_directory_server_config_v2.example]
}
```


## Argument Reference

* `name` - (Required) Name of Category Mapping. Recommendation is to use the name of the Active Directory group.
* `category_name` - (Optional) The name (key) for the category that this mapping is for. Defaults to `ADGroup`.
* `category_value` - (Required) The value for the category that this mapping is for.
* `ad_info` - (Required) The Active Directory object information for this mapping.
* `project_ext_id` - (Optional) The external identifier of the project associated with the Category Mapping.

### ad_info

* `directory_service_reference` - (Required) The ExtID of the directory service that will be used for mapping.
* `object_identifier` - (Required) The objectGUID for the object in AD. Use `nutanix_directory_service_users_search_v2` to look up this value. Cannot be updated once created.
* `object_path` - (Computed) The path for the mapped object in AD. Returned by the API.
* `status` - (Optional) The mapping status of AD Mapping. Valid values: `USABLE`, `DELETED`, `DIRECTORY_NOT_CONFIGURED`.

## Attributes Reference

* `ext_id` - A globally unique identifier of an instance that is suitable for external consumption.
* `links` - A HATEOAS style link for the response.
* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity.

## Import

Category Mapping can be imported using its `ext_id`:

```bash
terraform import nutanix_ad_group_category_mapping_v2.example <categoryMappingUUID>
```

See detailed information in [Nutanix Directory Server Configs V4](https://developers.nutanix.com/api-reference?namespace=microseg&version=v4.3).
