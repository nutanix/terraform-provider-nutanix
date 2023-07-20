---
layout: "nutanix"
page_title: "NUTANIX: nutanix_host"
sidebar_current: "docs-nutanix-datasource-host"
description: |-
 Describes a host
---

# nutanix_host

Describes a Host

## Example Usage

```hcl
data "nutanix_host" "host" {
   host_id = "<YOUR-host-ID>"
}`
```

## Argument Reference

The following arguments are supported:

* `host_id`: Represents hosts uuid

### Metadata

The metadata attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when image was last updated.
* `uuid`: - image uuid.
* `creation_time`: - UTC date and time in RFC-3339 format when image was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from server.
* `name`: - image name.
* `should_force_translate`: - Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.

### Categories

The categories attribute supports the following:

* `name`: - the key name.
* `value`: - value of the key.

## Attribute Reference

The following attributes are exported:

* `name`: -  The name for the image.
* `categories`: - Categories for the image.
* `project_reference`: - The reference to a project.
* `owner_reference`: - The reference to a user.
* `cluster_reference`: - Reference to a kind. Either one of (kind, uuid) or url needs to be specified.
* `api_version` - The API version.
* `gpu_driver_version`: - Host GPU driver version.
* `failover_cluster`: - Hyper-V failover cluster.
* `ipmi`: - Host IPMI info.
* `cpu_model`: - Host CPU model.
* `host_nics_id_list`: - Host NICs.
* `num_cpu_sockets`: - Number of CPU sockets.
* `windows_domain`: - The name of the node to be renamed to during domain-join. If not given,a new name will be automatically assigned.
* `gpu_list`: - List of GPUs on the host.
* `serial_number`: - Node serial number.
* `cpu_capacity_hz`: - Host CPU capacity.
* `memory_capacity_mib`: - Host memory capacity in MiB.
* `host_disks_reference_list`: - The reference to a disk.
* `monitoring_state`: - Host monitoring status.
* `hypervisor`: - Host Hypervisor information.
* `host_type`: - Host type.
* `num_cpu_cores`: - Number of CPU cores on Host.
* `rackable_unit_reference`: - The reference to a rackable_unit.
* `controller_vm`: - Host controller vm information.
* `block`: - Host block config info.

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, attributes supports the following:

* `kind`: - The kind name (Default value: project).
* `name`: - the name.
* `uuid`: - the uuid.

### Version

The version attribute supports the following:

* `product_name`: - Name of the producer/distribution of the image. For example windows or red hat.
* `product_version`: - Version string for the disk image.

See detailed information in [Nutanix Host](https://www.nutanix.dev/api_references/prism-central-v3/#/5ef25d36e143a-get-a-existing-host).