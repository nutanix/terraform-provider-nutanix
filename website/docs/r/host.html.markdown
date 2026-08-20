---
layout: "nutanix"
page_title: "NUTANIX: nutanix_host"
sidebar_current: "docs-nutanix-resource-nutanix-host"
description: |-
  This operation manages the categories assigned to an existing Nutanix host.
---

# nutanix_host

Manages the categories assigned to an existing Nutanix host.

Hosts are physical infrastructure that Prism Central discovers automatically - they are
never created or destroyed through the Nutanix API. This resource therefore "adopts" an
existing host, identified by `host_id`, and manages the category assignment on it, the
same way `categories` is managed on [`nutanix_subnet`](subnet.html) and other resources.

All other host attributes exposed by this resource are read-only and mirror the
[`nutanix_host`](/docs/providers/nutanix/d/host.html) data source; they describe the
adopted host but cannot be changed via Terraform.

~> **NOTE** Destroying this resource does not delete the underlying host. It only clears
the categories Terraform assigned to it and removes the resource from Terraform state.

## Example Usage

```hcl
data "nutanix_hosts" "hosts" {}

resource "nutanix_host" "managed-host" {
  host_id = data.nutanix_hosts.hosts.entities.0.metadata.uuid

  categories {
    name  = "Environment"
    value = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `host_id`: - (Required) The UUID of the existing host to adopt and manage. Changing this
  forces a new resource (it re-targets a different host rather than mutating the current one).
* `categories`: - (Optional) The categories to assign to the host. See below.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

## Attribute Reference

The following attributes are exported (all read-only, mirroring the current state of the
adopted host):

* `name`: - The host's name.
* `api_version` - The API version.
* `metadata`: - The host kind metadata.
* `project_reference`: - The reference to a project.
* `owner_reference`: - The reference to a user.
* `cluster_reference`: - Reference to the cluster the host belongs to.
* `gpu_driver_version`: - Host GPU driver version.
* `failover_cluster`: - Hyper-V failover cluster.
* `ipmi`: - Host IPMI info.
* `cpu_model`: - Host CPU model.
* `host_nics_id_list`: - Host NICs.
* `num_cpu_sockets`: - Number of CPU sockets.
* `windows_domain`: - The name of the node to be renamed to during domain-join. If not given, a new name will be automatically assigned.
* `gpu_list`: - List of GPUs on the host.
* `serial_number`: - Node serial number.
* `cpu_capacity_hz`: - Host CPU capacity.
* `memory_capacity_mib`: - Host memory capacity in MiB.
* `host_disks_reference_list`: - The reference to a disk.
* `monitoring_state`: - Host monitoring status.
* `hypervisor`: - Host Hypervisor information.
* `host_type`: - Host type.
* `num_cpu_cores`: - Number of CPU cores on the host.
* `rackable_unit_reference`: - The reference to a rackable_unit.
* `controller_vm`: - Host controller VM information.
* `block`: - Host block config info.

### Reference

The `project_reference`, `owner_reference`, and `cluster_reference` attributes support the following:

* `kind`: - The kind name.
* `name`: - the name.
* `uuid`: - the uuid.

## Importing

An existing `nutanix_host` resource can be imported using its host UUID:

```shell
terraform import nutanix_host.managed-host <HOST-UUID>
```

See detailed information in [Nutanix Host](https://www.nutanix.dev/api_references/prism-central-v3/#/5ef25d36e143a-get-a-existing-host).
