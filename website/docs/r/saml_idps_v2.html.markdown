---
layout: "nutanix"
page_title: "NUTANIX: nutanix_saml_identity_providers_v2"
sidebar_current: "docs-nutanix-resource-saml-identity-providers-v2"
description: |-
  Create a SAML Identity Provider.
---

# nutanix_saml_identity_providers_v2

Provides a resource to Create a SAML Identity Provider.

## Example Usage

```hcl
resource "nutanix_saml_identity_providers_v2" "idp" {
  name                        = "example_idp_name"
  idp_metadata {
    entity_id = "entity_id"
    login_url = "login_url"
    logout_url = "logout_url"
    error_url = "error_url"
    certificate = "certificate"
  }
  username_attribute          = "username"
  email_attribute             = "email"
  groups_attribute            = "groups"
  groups_delim                = "," # such as ',' or ';'
  idp_metadata_xml            = "<IDENTITY_PROVIDER_METADATA_XML content>"
  entity_issuer               = "entity_issuer_issuer"
  is_signed_authn_req_enabled = true
  custom_attributes           = ["custom1", "custom2"]
}
```

##  Argument Reference

The following arguments are supported:

* `ext_id`: -(Optional) External identifier of the SAML Identity Provider.
* `idp_metadata_url`: -(Optional) Metadata url that provides IDP details.
* `idp_metadata_xml`: -(Optional) Base64 encoded metadata in XML format with IDP details.
* `idp_metadata`: -(Optional) Type of the User Group. LDAP (User Group belonging to a Directory Service (Open LDAP/AD)),  SAML (User Group belonging to a SAML IDP.)
* `name`: -(Required) Unique name of the IDP.
* `username_attr`: -(Optional) SAML assertion Username attribute element.
* `email_attr`: -(Optional) SAML assertion email attribute element.
* `groups_attr`: -(Optional) SAML assertion groups attribute element.
* `groups_delim`: -(Optional) Delimiter is used to split the value of attribute into multiple groups.
* `custom_attr`: -(Optional) SAML assertions for list of custom attribute elements.
* `entity_issuer`: -(Optional) It will be used as Issuer in SAML authnRequest.
* `is_signed_authn_req_enabled`: -(Optional) Flag indicating signing of SAML authnRequests.

### Idp Metadata

The idp_metadata attribute supports the following:

* `entity_id`: -(Required) Entity Identifier of Identity provider.
* `login_url`: -(Required) Login URL of the Identity provider.
* `logout_url`: -(Optional) Logout URL of the Identity provider.
* `error_url`: - (Optional) Error URL of the Identity provider.
* `certificate`: -(Required) Certificate for verification.
* `name_id_policy_format`: -(Optional) Name ID Policy format.
  * supported values:
    * `emailAddress`: -  Uses email address as NameID format
    * `encrypted`: -  Uses encrypted as NameID format.
    * `unspecified`: -  NameID format is left to individual implementations.
    * `transient`: -  	Uses identifier with transient semantics as NameID format.
    * `WindowsDomainQualifiedName`: -  Uses Windows domain qualified name as NameID format.
    * `X509SubjectName`: -  	Uses X509SubjectName as NameID format.
    * `kerberos`: -  	Uses kerberos principal name as NameID format.
    * `persistent`: -  Uses persistent name identifier as NameID format.
    * `entity`: -  Uses identifier of an entity as NameID format.

## Attributes Reference
The following attributes are exported:

* `tenant_id` - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this Id to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
* `ext_id` - The External Identifier of the User Group.
* `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
* `idp_metadata`: - Type of the User Group. LDAP (User Group belonging to a Directory Service (Open LDAP/AD)),  SAML (User Group belonging to a SAML IDP.)
* `name`: - Unique name of the IDP.
* `username_attr`: - SAML assertion Username attribute element.
* `email_attr`: - SAML assertion email attribute element.
* `groups_attr`: - SAML assertion groups attribute element.
* `groups_delim`: - Delimiter is used to split the value of attribute into multiple groups.
* `custom_attr`: - SAML assertions for list of custom attribute elements.
* `entity_issuer`: - It will be used as Issuer in SAML authnRequest.
* `is_signed_authn_req_enabled`: - Flag indicating signing of SAML authnRequests.
* `created_time`: - Creation time of the SAML Identity Provider.
* `last_updated_time`: - Last updated time of the SAML Identity Provider.
* `created_by`: - User or Service who created the SAML Identity Provider.


### Links

The links attribute supports the following:

* `href`: - The URL at which the entity described by the link can be accessed.
* `rel`: - A name that identifies the relationship of the link to the object that is returned by the URL. The unique value of "self" identifies the URL for the object.

### Idp Metadata

The idp_metadata attribute supports the following:

* `entity_id`: - Entity Identifier of Identity provider.
* `login_url`: - Login URL of the Identity provider.
* `logout_url`: - Logout URL of the Identity provider.
* `error_url`: - Error URL of the Identity provider.
* `certificate`: - Certificate for verification.
* `name_id_policy_format`: - Name ID Policy format.
  * supported values:
    * `emailAddress`: -  Uses email address as NameID format
    * `encrypted`: -  Uses encrypted as NameID format.
    * `unspecified`: -  NameID format is left to individual implementations.
    * `transient`: -  	Uses identifier with transient semantics as NameID format.
    * `WindowsDomainQualifiedName`: -  Uses Windows domain qualified name as NameID format.
    * `X509SubjectName`: -  	Uses X509SubjectName as NameID format.
    * `kerberos`: -  	Uses kerberos principal name as NameID format.
    * `persistent`: -  Uses persistent name identifier as NameID format.
    * `entity`: -  Uses identifier of an entity as NameID format.



## Import

This helps to manage existing entities which are not created through terraform. SAML Identity Provider can be imported using the `UUID`. (ext_id in v4 API context).  eg,

```hcl
// create its configuration in the root module. For example:
resource "nutanix_saml_identity_providers_v2" "import_saml_idp" {}

// execute the below command. UUID can be fetched using datasource. Example: data "nutanix_saml_identity_providers_v2" "fetch_saml_idps"{}
terraform import nutanix_saml_identity_providers_v2.import_saml_idp <UUID>
```


See detailed information in [Nutanix SAML Identity Providers v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0).
