---
layout: "nutanix"
page_title: "NUTANIX: nutanix_eula_v2"
sidebar_current: "docs-nutanix-datasource-eula-v2"
description: |-
  API to fetch active End User License Agreement.
---

# nutanix_accept_eula_v2

API to fetch active End User License Agreement.

## Example Usage

```hcl
data "nutanix_eula_v2" "get_eula" {}
```

## Attributes Reference
* `tenant_id`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `content`: - Textual contents of the end user license agreement.
* `updated_time`: - EULA update time since epoch in ISO date time.
* `version`: - Version of the EULA.
* `is_enabled`: - Indicates whether this is the current EULA of the cluster or not.
* `acceptances`: - List of users accepting the EULA along with acceptance time for each.

### acceptances
The `acceptances` object contains the following attributes:

* `accepted_by`: - Model containing the EULA User Details attributes
* `acceptance_time`: - Date-time at which EULA was accepted.

### accepted_by
The `accepted_by` object contains the following attributes:

* `user_name`: - User name of the user accepting the EULA.
* `login_id`: - Login ID of the user accepting the EULA.
* `job_title`: - Job title of the user accepting the EULA.
* `company_name`: - Company name of the user accepting the EULA.


See detailed information in [Get End User License Agreement](https://developers.nutanix.com/api-reference?namespace=licensing&version=v4.1#tag/EndUserLicenseAgreement/operation/getEula)