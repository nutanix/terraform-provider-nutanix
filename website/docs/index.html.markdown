---
layout: "nutanix"
page_title: "Provider: Nutanix"
sidebar_current: "docs-nutanix-index"
description: |-
  The provider is used to interact with the many resources supported by Nutanix. The provider needs to be configured with the proper credentials before it can be used.
---

# Nutanix Provider

The provider is used to interact with the
many resources supported by Nutanix. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```terraform
provider "nutanix" {
  username     = "xxxx"
  password     = "xxxx"
  endpoint     = "xxxx"
  insecure     = true
  port         = 9440
  wait_timeout = 10
}
```

### Provider Variables

| Variable Name  | Description                                                                  | Default |
| -------------- | ---------------------------------------------------------------------------- | ------- |
| `username`     | Prism Central/Elements username                                              | N/A     |
| `password`     | Prism Central/Elements password                                              | N/A     |
| `endpoint`     | Prism Central/Elements URL                                                   | N/A     |
| `insecure`     | Whether to verify the SSL certificate on https connections                   | `false` |
| `port`         | Port of Prism Central/Elements                                               | `9440`  |
| `wait_timeout` | Timeout waiting for resource to create/update/destroy (in minutes)           | `1`     |
| `proxy_url`    | Proxy to use when connecting to Prism Central/Elements                       | N/A     |
| `session_auth` | Whether to use [session-based authentication](#session-based-authentication) | `false` |

## Authentication

The Nutanix provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

### Static credentials

Static credentials can be provided by adding the fowlloing attributes in-line in the Nutanix provider block:

Usage:

```terraform
provider "nutanix" {
  username     = "xxxx"
  password     = "xxxx"
  endpoint     = "xxxx"
  insecure     = true
  port         = 9440
  wait_timeout = 10 //Optional
  proxy_url    =  "xxxx" //Optional
}
```

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

### Environment variables

You can provide your credentials via environment variables, representing your Nutanix
authentication.

```hcl
provider "nutanix" {}
```

Usage:

``` bash
$ export NUTANIX_USERNAME="xxxx"
$ export NUTANIX_PASSWORD="xxxx"
$ export NUTANIX_INSECURE="xxxx"
$ export NUTANIX_PORT="xxxx"
$ export NUTANIX_ENDPOINT="xxxx"
$ export NUTANIX_WAIT_TIMEOUT = "xxx"
$ export NUTANIX_PROXY_URL = "xxx"

$ terraform plan
```
