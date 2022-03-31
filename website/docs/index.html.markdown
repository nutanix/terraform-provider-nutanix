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
**Note : Foundation feature in nutanix provider is in beta mode**

Foundation based examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/foundation/

Foundation based modules & examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/modules/foundation/

## Provider configuration required details

Going from 1.5.0-beta release of nutanix provider, fields inside provider configuration would be mandatory as per the usecase : 

* `Prism Central & Karbon` : For prism central and karbon related resources and data sources, `username`, `password` & `endpoint` are manadatory.
* `Foundation` : For foundation related resources and data sources, `foundation_endpoint` in manadatory.
