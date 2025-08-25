---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ova_image"
sidebar_current: "docs-nutanix-datasource-ova-image"
description: |-
  Describes an OVA Image
---

# nutanix_ova_image

Describes an OVA Image

## Example Usage

```hcl
data "nutanix_ova_image" "test" {
  ova_image_id = "your-ova-image-uuid"
}

data "nutanix_ova_image" "test_name" {
  ova_image_name = "your-ova-image-name"
}
```

## Argument Reference

The following arguments are supported:

* `ova_image_id`: (Optional) Represents the OVA image UUID. Conflicts with `ova_image_name`.
* `ova_image_name`: (Optional) Represents the OVA image name. Conflicts with `ova_image_id`.

## Attribute Reference

The following attributes are exported:

* `api_version`: - The API version.
* `name`: - The name of the OVA image.
* `description`: - A description for the OVA image.
* `state`: - The state of the OVA image.
* `categories`: - Categories for the OVA image.
* `project_reference`: - The reference to a project.
* `availability_zone_reference`: - The reference to an availability zone.
* `cluster_uuid`: - The UUID of the cluster where the OVA image resides.
* `cluster_name`: - The name of the cluster where the OVA image resides.
* `message_list`: - List of messages associated with the OVA image.
* `enable_cpu_passthrough`: - Indicates if CPU passthrough is enabled.
* `is_vcpu_hard_pinned`: - Indicates if vCPU is hard pinned.
* `num_vnuma_nodes`: - The number of vNUMA nodes.
* `nic_list`: - List of NICs attached to the OVA image.
* `guest_os_id`: - The guest OS ID.
* `power_state`: - The power state of the OVA image.
* `nutanix_guest_tools`: - Nutanix Guest Tools (NGT) information.
* `ngt_enabled_capability_list`: - List of NGT enabled capabilities.
* `ngt_credentials`: - NGT credentials.
* `num_vcpus_per_socket`: - Number of vCPUs per socket.
* `num_sockets`: - Number of sockets.
* `parent_reference`: - The reference to the parent.
* `memory_size_mib`: - Memory size in MiB.
* `boot_device_order_list`: - The boot device order list.
* `boot_device_disk_address`: - The disk address of the boot device.
* `boot_device_mac_address`: - The MAC address of the boot device.
* `boot_type`: - The boot type.
* `machine_type`: - The machine type.
* `hardware_clock_timezone`: - The hardware clock timezone.
* `guest_customization_cloud_init_meta_data`: - Cloud-init meta-data for guest customization.
* `guest_customization_cloud_init_user_data`: - Cloud-init user-data for guest customization.
* `guest_customization_cloud_init_custom_key_values`: - Custom key-values for cloud-init guest customization.
* `guest_customization_is_overridable`: - Indicates if guest customization is overridable.
* `guest_customization_sysprep`: - Sysprep settings for guest customization.
* `guest_customization_sysprep_custom_key_values`: - Custom key-values for Sysprep guest customization.
* `should_fail_on_script_failure`: - Indicates if the script should fail on script failure.
* `enable_script_exec`: - Indicates if script execution is enabled.
* `power_state_mechanism`: - The power state mechanism.
* `vga_console_enabled`: - Indicates if VGA console is enabled.
* `disk_list`: - List of disks attached to the OVA image.
* `serial_port_list`: - List of serial ports.
* `host_reference`: - Host reference.

---

### Metadata

The `metadata` attribute exports the following:

* `last_update_time`: - UTC date and time in RFC-3339 format when the OVA image was last updated.
* `UUID`: - OVA image UUID.
* `creation_time`: - UTC date and time in RFC-3339 format when the OVA image was created.
* `spec_version`: - Version number of the latest spec.
* `spec_hash`: - Hash of the spec. This will be returned from the server.
* `name`: - OVA image name.

---

### Reference

The `project_reference`, `owner_reference`, `availability_zone_reference`, `cluster_reference`, `parent_reference`, `data_source_reference`, `volume_group_reference` and `network_function_chain_reference` attributes support the following:

* `kind`: - The kind name (Default value: project).
* `name`: - The name.
* `uuid`: - The UUID.

---

### Nic List

The `nic_list` attribute exports the following:

* `nic_type`: - The type of NIC.
* `uuid`: - The UUID of the NIC.
* `floating_ip`: - The floating IP address.
* `model`: - The NIC model.
* `network_function_nic_type`: - The network function NIC type.
* `mac_address`: - The MAC address of the NIC.
* `ip_endpoint_list`: - List of IP endpoints.
* `network_function_chain_reference`: - The reference to the network function chain.
* `num_queues`: - The number of queues.
* `subnet_uuid`: - The UUID of the subnet.
* `subnet_name`: - The name of the subnet.
* `is_connected`: - Indicates if the NIC is connected.

---

### IP Endpoint List

The `ip_endpoint_list` attribute exports the following:

* `ip`: - The IP address.
* `type`: - The type of IP endpoint.

---

### Disk List

The `disk_list` attribute exports the following:

* `uuid`: - The UUID of the disk.
* `disk_size_bytes`: - The size of the disk in bytes.
* `disk_size_mib`: - The size of the disk in MiB.
* `storage_config`: - Storage configuration for the disk.
* `device_properties`: - Device properties of the disk.
* `data_source_reference`: - Reference to a data source.
* `volume_group_reference`: - Reference to a volume group.

---

### Storage Config

The `storage_config` attribute exports the following:

* `flash_mode`: - The flash mode.
* `storage_container_reference`: - Reference to the storage container.

---

### Device Properties

The `device_properties` attribute exports the following:

* `device_type`: - The device type.
* `disk_address`: - The disk address.

---

### Disk Address

The `disk_address` attribute exports the following:

* `device_index`: - The device index.
* `adapter_type`: - The adapter type.

---

### Serial Port List

The `serial_port_list` attribute exports the following:

* `index`: - The index of the serial port.
* `is_connected`: - Indicates if the serial port is connected.

---

### Guest Customization Sysprep

The `guest_customization_sysprep` attribute exports the following:

* `install_type`: - The Sysprep installation type.
* `unattend_xml`: - The Unattend XML content.

See detailed information in [Nutanix OVA Image](https://www.nutanix.dev/api_reference/apis/prism_v3.html#tag/ovas/paths/~1ovas~1%7Buuid%7D/get)