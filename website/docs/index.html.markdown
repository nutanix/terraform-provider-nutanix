---
layout: "nutanix"
page_title: "Provider: Nutanix"
sidebar_current: "docs-nutanix-index"
description: |-
  The provider is used to interact with the many resources supported by Nutanix. The provider needs to be configured with the proper credentials before it can be used.
---

# Nutanix Provider

The provider is used to interact with the many resources and data sources supported by Nutanix, using Prism Central as the provider endpoint.

Use the navigation on the left to read about the available resources and data sources this provider can use.

~> **Important Notice:** Upcoming Deprecation of Legacy Nutanix Terraform Provider Resources. Starting with the Nutanix Terraform Provider release planned for Q4-CY2026, legacy resources which are based on v0.8,v1,v2 and v3 APIs will be deprecated and no longer supported. For more information, visit [Legacy API Deprecation Announcement](https://portal.nutanix.com/page/documents/eol/list?type=announcement) [Legacy API Deprecation - FAQs](https://portal.nutanix.com/page/documents/kbs/details?targetId=kA0VO0000005rgP0AQ). Nutanix strongly encourages you to migrate your scripts and applications to the latest v2 version of the Nutanix Terraform Provider resources, which are built on our v4 APIs/SDKs. By adopting the latest v2 version based on v4 APIs and SDKs, our users can leverage the enhanced capabilities and latest innovations from Nutanix. We understand that this transition may require some effort, and we are committed to supporting you throughout the process. Please refer to our documentation and support channels for guidance and assistance.

## Support

Terraform Nutanix Provider leverages the community-supported model. See [Open Source Support](https://portal.nutanix.com/page/documents/kbs/details?targetId=kA07V000000LdWPSA0) for more information about its support policy.

-> **Note:** Update!! 
We now have a brand new developer-centric Support Program designed for organizations that require a deeper level of developer support to manage their Nutanix environment and build applications quickly and efficiently. As part of this new Advanced API/SDK Support Program, you will get access to trusted technical advisors who specialize in developer tools including Nutanix Terraform Provider and receive support for your unique development needs and custom integration queries.Visit our Support Portal - Premium Add-On Support Programs  to learn more about this program. **Contributions to open-source Nutanix Terraform Provider repository will continue to leverage a community-supported model. Visit https://portal.nutanix.com/kb/13424  for more details. 


## Example Usage

### Terraform 0.12 and below

```terraform
provider "nutanix" {
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_endpoint
  port         = var.nutanix_port
  insecure     = true
  wait_timeout = 10
}
```

### Terraform 0.13+

```terraform
terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = "1.2.0"
    }
  }
}

provider "nutanix" {
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_endpoint
  port         = var.nutanix_port
  insecure     = true
  wait_timeout = 10
}
```

## Argument Reference

The following arguments are used to configure the Nutanix Provider:
* `username` - **(Required)** This is the username for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_USERNAME` environment variable.
* `password` - **(Required)** This is the password for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_PASSWORD` environment variable.
* `endpoint` - **(Required)** This is the endpoint for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_ENDPOINT` environment variable.
* `insecure` - (Optional) This specifies whether to allow verify ssl certificates. This can also be specified with `NUTANIX_INSECURE`. Defaults to `false`.
* `port` - (Optional) This is the port for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_PORT` environment variable. Defaults to `9440`.
* `session_auth` - (Optional) This specifies whether to use [session authentication](#session-based-authentication). This can also be specified with the `NUTANIX_SESSION_AUTH` environment variable. Defaults to `true`
* `wait_timeout` - (Optional) This specifies the timeout on all resource operations in the provider in minutes. This can also be specified with the `NUTANIX_WAIT_TIMEOUT` environment variable. Defaults to `1`. Also see [resource timeouts](#resource-timeouts).
* `proxy_url` - (Optional) This specifies the url to proxy through to access the Prism Elements or Prism Central endpoint. This can also be specified with the `NUTANIX_PROXY_URL` environment variable.

### Session based Authentication

Session based authentication can be used which authenticates only once with basic authentication and uses a cookie for all further attempts.
The main benefit is a reduction in the time API calls take to complete. Sessions are only valid for 15 minutes.

Usage:

```terraform
provider "nutanix" {
  ...
  session_auth = true
  ...
}
```

## Notes

### Resource Timeouts
Currently, the only way to set a timeout is using the `wait_timeout` argument or `NUTANIX_WAIT_TIMEOUT` environment variable. This will set a timeout for all operations on all resources. This provider currently doesn't support specifying [operation timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts).

## Nutanix Foundation (>=v1.5.0-beta)

Going from 1.5.0-beta release of nutanix provider, two more params are added to provider configuration to support foundation components :

* `foundation_endpoint` - (Optional) This is the endpoint for foundation vm. This can also be specified with the `FOUNDATION_ENDPOINT` environment variable.
* `foundation_port` - (Optional) This is the port for foundation vm. This can also be specified with the `FOUNDATION_PORT` environment variable. Default is `8000`.

```terraform
terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = ">=1.5.0-beta"
    }
  }
}

provider "nutanix" {
  username            = var.nutanix_username
  password            = var.nutanix_password
  endpoint            = var.nutanix_endpoint
  port                = var.nutanix_port
  insecure            = true
  wait_timeout        = 10
  foundation_endpoint = var.foundation_endpoint
  foundation_port     = var.foundation_port
}
```

Foundation based examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/foundation/

Foundation based modules & examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/modules/foundation/

## Nutanix Database Service (NDB) (>=v1.8.0)

Going from 1.8.0 release of nutanix provider, some params are added to provider configuration to support Nutanix Database Service (NDB) components :

* `ndb_username` - (Optional) This is the username for the NDB instance. This can also be specified with the `NDB_USERNAME` environment variable.
* `ndb_password` - (Optional) This is the password for the NDB instance. This can also be specified with the `NDB_PASSWORD` environment variable.
* `ndb_endpoint` - (Optional) This is the endpoint for the NDB instance. This can also be specified with the `NDB_ENDPOINT` environment variable.

```terraform
terraform {
  required_providers {
    nutanix = {
      source = "nutanix/nutanix"
      version = ">=1.8.0"
    }
  }
}

provider "nutanix" {
  username            = var.nutanix_username
  password            = var.nutanix_password
  endpoint            = var.nutanix_endpoint
  port                = var.nutanix_port
  insecure            = true
  wait_timeout        = 10
  ndb_endpoint        = var.ndb_endpoint 
  ndb_username        = var.ndb_username
  ndb_password        = var.ndb_password
}
```

NDB based examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/ndb/

## Provider configuration required details

Going from 1.8.0-beta release of nutanix provider, fields inside provider configuration would be mandatory as per the usecase : 

* `Prism Central & Karbon` : For prism central and karbon related resources and data sources, `username`, `password` & `endpoint` are manadatory.
* `Foundation` : For foundation related resources and data sources, `foundation_endpoint` in manadatory.
* `NDB` : For Nutanix Database Service (NDB) related resources and data sources. 

