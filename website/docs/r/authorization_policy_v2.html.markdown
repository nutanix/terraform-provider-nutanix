---
layout: "nutanix"
page_title: "NUTANIX: authorization_policy_v2"
sidebar_current: "docs-nutanix-resource-authorization-policy-v2"
description: |-
  Create Virtual Private Cloud .
---

# authorization_policy_v2

Provides Nutanix resource to create authorization policy.

## Example

```hcl

resource "nutanix_authorization_policy_v2" "ap-example"{
  role                      = "ba250e3e-1db1-4950-917f-a9e2ea35b8e3"
  display_name              = "auth_policy_example"
  description               = "authorization policy example"
  authorization_policy_type = "USER_DEFINED"
  identities {
    reserved = "{\"user\":{\"uuid\":{\"anyof\":[\"00000000-0000-0000-0000-000000000000\"]}}}"
  }
  entities {
    reserved = "{\"images\":{\"*\":{\"eq\":\"*\"}}}"
  }
  entities {
    reserved = "{\"marketplace_item\":{\"owner_uuid\":{\"eq\":\"SELF_OWNED\"}}}"
  }
}
```

## Argument Reference

The following arguments are supported:

- `display_name`: Name of the Authorization Policy.
- `description`: Description of the Authorization Policy.
- `client_name`: Client that created the entity.
- `identities`: The identities for which the Authorization Policy is created.
- `entities`: The entities being qualified by the Authorization Policy.
- `role`: The Role associated with the Authorization Policy.
- `authorization_policy_type`: Type of Authorization Policy.
  - `PREDEFINED_READ_ONLY` : System-defined read-only ACP, i.e. no modifications allowed.
  - `SERVICE_DEFINED_READ_ONLY` : Read-only ACP defined by a service.
  - `PREDEFINED_UPDATE_IDENTITY_ONLY` : System-defined ACP prohibiting any modifications from customer.
  - `SERVICE_DEFINED` : ACP defined by a service.
  - `USER_DEFINED` : ACP defined by an User.

## Attribute Reference

The following attributes are exported:

- `ext_id`: ext_id of Authorization policy.
- `links`: A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `display_name`: Name of the Authorization Policy.
- `description`: Description of the Authorization Policy.
- `client_name`: Client that created the entity.
- `identities`: The identities for which the Authorization Policy is created.
- `entities`: The entities being qualified by the Authorization Policy.
- `role`: The Role associated with the Authorization Policy.
- `created_time`: The creation time of the Authorization Policy.
- `last_updated_time`: The time when the Authorization Policy was last updated.
- `created_by`: User or Service Name that created the Authorization Policy.
- `is_system_defined`: Flag identifying if the Authorization Policy is system defined or not.
- `authorization_policy_type`: Type of Authorization Policy.

### Links

The links attribute supports the following:

- `href`: - The URL at which the entity described by the link can be accessed.
- `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object

## Import

This helps to manage existing entities which are not created through terraform. authorization policy can be imported using the `UUID`. (ext_id in v4 API context).  eg,
```hcl
// create its configuration in the root module. For example:
resource "nutanix_authorization_policy_v2" "import_policy" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_authorization_policies_v2" "fetch_policies"{}
terraform import nutanix_authorization_policy_v2.import_policy <UUID>
```

See detailed information in [Nutanix Authorization Policy v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/AuthorizationPolicies/operation/createAuthorizationPolicy).
