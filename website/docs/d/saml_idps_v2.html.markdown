---
layout: "nutanix"
page_title: "NUTANIX: nutanix_saml_identity_providers_v2"
sidebar_current: "docs-nutanix-datasource-saml-identity-providers-v2"
description: |-
  Provides a datasource to retrieve all the  all SAML Identity Provider(s).
---

# nutanix_saml_identity_providers_v2

Provides a datasource to retrieve all the SAML Identity Provider(s).

## Example Usage

```hcl
data "nutanix_saml_identity_providers_v2" "idps-list"{}

# list saml identity providers
data "nutanix_saml_identity_providers_v2" "filtered-idps"{
  filter = "name eq 'idp_example_name'"
  limit  = 2
}

```

##  Argument Reference

The following arguments are supported:

* `page`: - A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
* `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
* `filter` :A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:
  - createdBy
  - extId
  - name
* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:
  - createdTime
  - lastUpdatedTime
  - name
* `select` : A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., *), then all properties on the matching resource will be returned. The select can be applied to the following fields:
  - createdBy
  - createdTime
  - customAttributes
  - emailAttribute
  - entityIssuer
  - extId
  - groupsAttribute
  - groupsDelim
  - idpMetadata/certificate
  - idpMetadata/entityId
  - idpMetadata/errorUrl
  - idpMetadata/loginUrl
  - idpMetadata/logoutUrl
  - idpMetadata/nameIdPolicyFormat
  - isSignedAuthnReqEnabled
  - lastUpdatedTime
  - links
  - name
  - tenantId
  - usernameAttribute

## Attributes Reference
The following attributes are exported:

* `identity_providers` : List all SAML Identity Provider(s).

### Identity Providers

The identity_providers  attribute element contains the following attributes:

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
* `name_id_policy_format`: - Name ID Policy format. supported values:
    * `emailAddress`: -  Uses email address as NameID format
    * `encrypted`: -  Uses encrypted as NameID format.
    * `unspecified`: -  NameID format is left to individual implementations.
    * `transient`: -  	Uses identifier with transient semantics as NameID format.
    * `WindowsDomainQualifiedName`: -  Uses Windows domain qualified name as NameID format.
    * `X509SubjectName`: -  	Uses X509SubjectName as NameID format.
    * `kerberos`: -  	Uses kerberos principal name as NameID format.
    * `persistent`: -  Uses persistent name identifier as NameID format.
    * `entity`: -  Uses identifier of an entity as NameID format.

See detailed information in [Nutanix List SAML identity providers v4](https://developers.nutanix.com/api-reference?namespace=iam&version=v4.0#tag/SAMLIdentityProviders/operation/listSamlIdentityProviders).
