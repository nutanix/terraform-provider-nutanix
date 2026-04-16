---
layout: "nutanix"
page_title: "NUTANIX: nutanix_authorization_policy_v2"
sidebar_current: "docs-nutanix-datasource-authorization-policy-v2"
description: |-
   Provides a datasource to retrieve authorization policy with  authorization policy uuid .
---
# nutanix_authorization_policy_v2

Provides a datasource to retrieve authorization policy with  authorization policy uuid .

## Example Usage

```hcl
data "nutanix_authorization_policy_v2" "get-acp"{
  ext_id = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
}

```

## Argument Reference

The following arguments are supported:

* `ext_id`: (Required) Authorization Policy UUID

## Attribute Reference

The following attributes are exported:
* `ext_id`: ext_id of Authorization policy.
* `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `display_name`: Name of the Authorization Policy.
* `description`: Description of the Authorization Policy.
* `client_name`: Client that created the entity.
* `identities`: The identities for which the Authorization Policy is created.
* `entities`: The entities being qualified by the Authorization Policy.
* `role`: The Role associated with the Authorization Policy.
* `created_time`: The creation time of the Authorization Policy.
* `last_updated_time`: The time when the Authorization Policy was last updated.
* `created_by`: User or Service Name that created the Authorization Policy.
* `is_system_defined`: Flag identifying if the Authorization Policy is system defined or not.
* `authorization_policy_type`: Type of Authorization Policy.
    * `PREDEFINED_READ_ONLY` : System-defined read-only ACP, i.e. no modifications allowed.
    * `SERVICE_DEFINED_READ_ONLY` : Read-only ACP defined by a service.
    * `PREDEFINED_UPDATE_IDENTITY_ONLY` : System-defined ACP prohibiting any modifications from customer.
    * `SERVICE_DEFINED` : ACP defined by a service.
    * `USER_DEFINED` : ACP defined by an User.


### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object


See detailed information in [Nutanix Get Authorization Policy v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/AuthorizationPolicies/operation/getAuthorizationPolicyById).
