---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_ipmi_config"
sidebar_current: "docs-nutanix-resource-foundation-ipmi-config"
description: |-
  Configures IPMI IP address on BMC of nodes.
---

# nutanix_foundation_ipmi_config

Configures IPMI IP address on BMC of nodes.

## Example Usage

```hcl
resource "nutanix_foundation_ipmi_config" "impi-1" {
  ipmi_user = "username"
  ipmi_netmask = "10.xx.xx.xx"
  blocks{
    nodes {
          ipmi_mac = "ff:ff:ff:ff:ff:ff"
          ipmi_configure_now =  true
          ipmi_ip = "10.xx.xx.xx"
    }
    nodes {
          ipmi_mac = "ff:ff:ff:ff:ff:ff"
          ipmi_configure_now = true
          ipmi_ip = "10.xx.xx.xx"
    }
    block_id = "xyz"
  }
  ipmi_gateway = "10.xx.xx.xx"
  ipmi_password = "XXXXX"
}
```

## Argument Reference

The following arguments are supported:

* `ipmi_user`: - (Required) IPMI username.
* `ipmi_password`: - (Required) IPMI password.
* `ipmi_netmask`: - (Required) IPMI netmask.
* `ipmi_gateway`: - (Required) IPMI gateway.
* `blocks`: - (Required) List of blocks.

### blocks

The blocks attribute's each element supports following :

* `nodes`: - (Required) array of nodes for ipmi config.
* `block_id`: - (Optional) Block Id

### nodes

The nodes attribute's each element supports following :

* `ipmi_mac`: - (Required) IPMI mac address.
* `ipmi_configure_now`: - (Required) Whether IPMI should be configured. Should be kept true to configure
* `ipmi_ip`: - IPMI IP address.

## Attributes Reference

The following attributes are exported:

* `ipmi_user`: -  IPMI username.
* `ipmi_password`: - IPMI password.
* `ipmi_netmask`: - IPMI netmask.
* `ipmi_gateway`: - IPMI gateway.
* `blocks`: - List of blocks.

### blocks

The blocks attribute's each element exports following :

* `nodes`: - array of nodes for ipmi config.
* `block_id`: - Block Id

### nodes

The nodes attribute's each element exports following :

* `ipmi_mac`: - IPMI mac address.
* `ipmi_configure_now`: - Whether IPMI should be configured.
* `ipmi_ip`: - IPMI IP address.
* `ipmi_configure_successful` : - Whether IPMI was successfully configured.
* `ipmi_message` : - IPMI configuration status message if any.

## Error 

Incase of error, terraform will error out and display error for every failed ipmi configuration respective to its ipmi_ip.

## lifecycle

* `Update` : - Resource will trigger new resource create call for any kind of update in resource block.
* `Delete` : - Delete will be a soft delete.

See detailed information in [Nutanix Foundation IPMI Config](https://www.nutanix.dev/api_references/foundation/#/b3A6MjIyMjMzNzI-configure-bmc-i-pv4).
