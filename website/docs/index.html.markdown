---
layout: "nutanix"
page_title: "Provider: Nutanix"
sidebar_current: "docs-nutanix-index"
description: |-
  The provider is used to interact with the many resources supported by Nutanix. The provider needs to be configured with the proper credentials before it can be used.
---

# Nutanix Provider

The provider is used to interact with the
many resources and data sources supported by Nutanix, using either Prism Central or Prism Element as the provider endpoint.

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
- `username` - **(Required)** This is the username for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_USERNAME` environment variable.
- `password` - **(Required)** This is the password for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_PASSWORD` environment variable.
- `endpoint` - **(Required)** This is the endpoint for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_ENDPOINT` environment variable.
- `insecure` - (Optional) This specifies whether to allow verify ssl certificates. This can also be specified with `NUTANIX_INSECURE`. Defaults to `false`.
- `port` - (Optional) This is the port for the Prism Elements or Prism Central instance. This can also be specified with the `NUTANIX_PORT` environment variable. Defaults to `9440`.
- `session_auth` - (Optional) This specifies whether to use [session authentication](#session-based-authentication). This can also be specified with `NUTANIX_SESSION_AUTH`. Defaults to `true`
- `wait_timeout` - (Optional) This specifies the timeout on all resource operations in the provider in minutes. This can also be specified with `NUTANIX_WAIT_TIMEOUT`. Defaults to `1`. Also see [resource timeouts](#resource-timeouts).
- `proxy_url` - (Optional) This specifies the url to proxy through to access the Prism Elements or Prism Central endpoint.

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
Currently, the only way to set a timeout is using the `wait_timeout` attribute or `NUTANIX_WAIT_TIMEOUT` environment variable. This will set a timeout for all operations on all resources. This provider currently doesn't support specifying [operation timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts).