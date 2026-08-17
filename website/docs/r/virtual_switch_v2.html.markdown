---
layout: "nutanix"
page_title: "NUTANIX: nutanix_virtual_switch_v2"
sidebar_current: "docs-nutanix-resource-virtual-switch-v2"
description: |-
  Create, update, and delete Virtual Switches and their underlying physical network configurations in the Nutanix v4 API.
---

# nutanix_virtual_switch_v2

The `nutanix_virtual_switch_v2` resource allows you to create, manage, update, and delete Virtual Switches across your Nutanix clusters. 

This resource exposes two distinct create paths through one schema. The
provider picks the path based on whether `clusters[].existing_bridge_name`
is set in your configuration. See [Migration vs standard create](#migration-vs-standard-create)
below for details before reaching for the second example.

## Example Usage

### Standard Create (Build a Virtual Switch from scratch)

Use this method when you want Nutanix to allocate a new internal bridge (e.g., `br1`, `br2`) and explicitly bind physical NICs to it.

```hcl
resource "nutanix_virtual_switch_v2" "example" {
  name        = "example-vs"
  description = "Standard Virtual Switch spanning specific hosts"
  bond_mode   = "ACTIVE_BACKUP"
  mtu         = 1500

  clusters {
    ext_id = "00000000-0000-0000-0000-000000000000"
    
    hosts {
      ext_id    = "00000000-0000-0000-0000-000000000001"
      host_nics = ["eth0", "eth1"]
      active_uplink = "eth0"
    }
  }

  igmp_spec {
    is_snooping_enabled = true
    snooping_timeout    = 300
    querier_spec {
      is_querier_enabled = false
    }
  }
}
```

### Migrate an Existing OVS Bridge into a Virtual Switch

When `clusters[].existing_bridge_name` is set, the provider routes the
create call to `POST /api/networking/v4.3/config/virtual-switches/$actions/migrate`
instead of the standard create endpoint. This is the only way to bind a
new Virtual Switch to a specific, pre-existing OVS bridge -- the standard
create endpoint silently auto-allocates the next free `brN` and ignores
any pre-existing bridge hint.

The bridge itself must already exist on the host. There is no Terraform
primitive for OVS bridge creation; create it out-of-band, for example:

```bash
ssh nutanix@<PE-IP> "source /etc/profile; \
  /usr/local/nutanix/cluster/bin/manage_ovs --bridge_name br1 create_single_bridge"
```

Once the OVS bridge exists, you can consume and convert it using Terraform:

```hcl
resource "nutanix_virtual_switch_v2" "from_existing_bridge" {
  name        = "example-vs-from-bridge"
  description = "Virtual Switch built by migrating existing br1"

  clusters {
    ext_id               = "00000000-0000-0000-0000-000000000000"
    existing_bridge_name = "br1"
    
    hosts {
      ext_id = "00000000-0000-0000-0000-000000000001"
    }
  }
}
```

## Migration vs Standard Create

| Path | Trigger Condition | API Endpoint | Primary Use Case |
| --- | --- | --- | --- |
| **Standard Create** | `clusters[].existing_bridge_name` is **not** set | `POST /api/networking/v4.3/config/virtual-switches` | Build a completely new Virtual Switch with explicit host-NIC bindings, IPs, bond modes, and IGMP settings. |
| **Migrate from Bridge** | `clusters[].existing_bridge_name` is set | `POST /api/networking/v4.3/config/virtual-switches/$actions/migrate` | Convert a pre-existing physical OVS bridge into a managed Virtual Switch. The switch automatically inherits the bridge's existing NIC bindings. |

### Important Migration Limitations

When the migrate path is triggered, the Nutanix API **silently ignores** the following top-level fields during creation, even if they are defined in your Terraform schema: 
* `bond_mode`, `mtu`, `igmp_spec`, `is_quick_mode`, `shared_with_projects`, and `owner_type`.

The provider will emit a `[WARN]` log for each ignored field. To apply these configurations to a migrated bridge, you must run `terraform apply` a second time. The standard *update* endpoint (used during the second apply) will successfully process these fields.

Additionally:
* `existing_bridge_name` is a **create-time-only** parameter. Because the API never returns this field on reads, Terraform preserves the initial value in the state file to prevent constant drift during `terraform plan`.
* Changing `existing_bridge_name` after creation has no effect; treat it as immutable. 
* Setting `existing_bridge_name` on multiple `clusters[]` blocks will be rejected by the provider at apply time. Define it only on `clusters[0]`.

## Argument Reference

* `name` - (Required) The human-readable name of the Virtual Switch.
* `description` - (Optional) A descriptive text providing context or details about the Virtual Switch's purpose.
* `bond_mode` - (Optional) The uplink bonding policy type. Valid values are `ACTIVE_BACKUP`, `BALANCE_SLB` (Source Load Balancing), `BALANCE_TCP`, or `NONE`.
* `mtu` - (Optional) The Maximum Transmission Unit (MTU) for the virtual switch, in bytes.
* `is_quick_mode` - (Optional) When set to `true`, host nodes are not put into maintenance mode during Virtual Switch updates, expediting the process.
* `project_ext_id` - (Optional) The UUID of the Nutanix Project that owns this entity, utilized for Role-Based Access Control (RBAC). This can only be set at creation time; changing it on an existing virtual switch is not supported and will return an error.
* `shared_with_projects` - (Optional) A list of project UUIDs that this virtual switch is shared with, allowing them to utilize the network. Sharing is reconciled through the dedicated share/unshare endpoints: projects added to this list are shared on `apply`, projects removed from it are unshared, and all projects are unshared automatically before the Virtual Switch is destroyed. (Note: the migrate path ignores this field at create time; run `terraform apply` a second time to share a migrated switch.)
* `igmp_spec` - (Optional) A block to configure Internet Group Management Protocol (IGMP) settings. See [`igmp_spec`](#igmp_spec) below.
* `clusters` - (Optional) A block defining the configuration of the Virtual Switch on specific clusters. See [`clusters`](#clusters) below.

### clusters

The `clusters` block supports the following arguments:

* `ext_id` - (Optional) The UUID of the cluster where this Virtual Switch configuration is applied. When using the migrate path, this acts as the `clusterReference`.
* `existing_bridge_name` - (Optional, Create Only) The name of an existing OVS bridge on the host (e.g., `br1`) to convert into this Virtual Switch. Setting this triggers the migration path.
* `gateway_ip_address` - (Optional) The Gateway IP address for the virtual switch on this cluster. *(Only applicable during Standard Create; ignored during Migration).*
* `vlan_identifier` - (Optional) The VLAN Identifier for this virtual switch within the cluster. Set to `0` to explicitly remove VLAN tagging.
* `hosts` - (Optional) A list of host-specific configurations within the cluster. See [`hosts`](#hosts) below.

### hosts

The `hosts` block details how the Virtual Switch connects to physical nodes:

* `ext_id` - (Optional) The UUID of the physical host.
* `host_nics` - (Optional) A list of physical Network Interface Cards (e.g., `["eth0", "eth1"]`) to bind to this virtual switch as uplinks. *(Only applicable during Standard Create; ignored during Migration).*
* `active_uplink` - (Optional) Specifies the active uplink interface when using the `ACTIVE_BACKUP` bond mode. *(Standard create only).*
* `internal_bridge_name` - (Optional) The internal logical bridge name (e.g., `br0`). In practice, this is **read-only**. The standard create endpoint auto-allocates this and ignores user input. To target a specific bridge, use `clusters[].existing_bridge_name` instead.
* `ip_address` - (Optional) The IPv4 subnet address for the bridge interface on this host. *(Only applicable during Standard Create; ignored during Migration).*
* `route_table` - (Optional) The internal route table number utilized by this host.

### igmp_spec

The `igmp_spec` block manages multicast traffic settings. *(Note: The migrate endpoint ignores `igmp_spec` during creation. Apply a second time to set these values on a migrated switch).*

* `is_snooping_enabled` - (Optional) Enables IGMP snooping to optimize multicast traffic delivery.
* `snooping_timeout` - (Optional) The timeout value in seconds before an inactive multicast group membership expires.
* `querier_spec` - (Optional) A block defining IGMP querier settings. See [`querier_spec`](#querier_spec) below.

### querier_spec

* `is_querier_enabled` - (Optional) Enables the Virtual Switch to act as an IGMP querier.
* `vlan_id_list` - (Optional) A list of VLAN IDs on which the Virtual Switch will actively send IGMP queries.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `ext_id` - The globally unique identifier (UUID) of the created Virtual Switch instance.
* `is_default` - A boolean indicating whether this is a default system Virtual Switch (which generally cannot be deleted).
* `has_deployment_error` - A boolean indicating if the configuration failed to deploy consistently across all cluster nodes.
* `has_update_in_progress` - A boolean indicating if the virtual switch is currently undergoing an update operation.
* `has_delete_in_progress` - A boolean indicating if the virtual switch is currently in the process of being deleted.
* `owner_type` - Identifies the type of entity that created or owns the Virtual Switch.
* `tenant_id` - A globally unique identifier representing the tenant that owns this entity.
* `links` - A list of HATEOAS-style links related to the API response.
* `metadata` - A block containing administrative and organizational metadata (such as categories).

## Import

Virtual Switches can be imported into Terraform using their `ext_id`:

```shell
terraform import nutanix_virtual_switch_v2.example <ext_id>
```

**Note on Imports:** Imported Virtual Switches will always reflect `clusters[].existing_bridge_name = ""` in Terraform state, regardless of whether they were originally created via the standard or migration paths. This is expected behavior because the Nutanix API does not record or return the original migration hint.

See detailed information in [Nutanix Create Virtual Switch V4](https://developers.nutanix.com/api-reference?namespace=networking&version=v4.3#tag/VirtualSwitches/operation/createVirtualSwitch).
