package nutanix

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixVirtualMachineRead,

		Schema: getDataSourceVMSchema(),
	}
}

func dataSourceNutanixVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	vm, ok := d.GetOk("vm_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute vm_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetVM(vm.(string))
	if err != nil {
		return err
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", getReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", getReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_reference", getReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	if err := d.Set("cluster_reference", getClusterReferenceValues(resp.Status.ClusterReference)); err != nil {
		return err
	}

	d.Set("cluster_reference_name", utils.StringValue(resp.Status.ClusterReference.Name))
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("num_vnuma_nodes", utils.Int64Value(resp.Status.Resources.VnumaConfig.NumVnumaNodes))

	nicLists := make([]map[string]interface{}, 0)

	if resp.Status.Resources.NicList != nil {
		nicLists = make([]map[string]interface{}, len(resp.Status.Resources.NicList))
		for k, v := range resp.Status.Resources.NicList {
			nic := make(map[string]interface{})
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["floating_ip"] = utils.StringValue(v.FloatingIP)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)

			ipEndpointList := make([]map[string]interface{}, len(v.IPEndpointList))
			for k1, v1 := range v.IPEndpointList {
				ipEndpoint := make(map[string]interface{})
				ipEndpoint["ip"] = utils.StringValue(v1.IP)
				ipEndpoint["type"] = utils.StringValue(v1.Type)
				ipEndpointList[k1] = ipEndpoint
			}
			nic["ip_endpoint_list"] = ipEndpointList
			nic["network_function_chain_reference"] = getReferenceValues(v.NetworkFunctionChainReference)
			nic["subnet_reference"] = getReferenceValues(v.SubnetReference)

			nicLists[k] = nic
		}
	}
	if err := d.Set("nic_list", nicLists); err != nil {
		return err
	}
	if err := d.Set("host_reference", getReferenceValues(resp.Status.Resources.HostReference)); err != nil {
		return err
	}
	if err := d.Set("guest_os_id", utils.StringValue(resp.Status.Resources.GuestOsID)); err != nil {
		return err
	}
	if err := d.Set("power_state", utils.StringValue(resp.Status.Resources.PowerState)); err != nil {
		return err
	}

	nutanixGuestTools := make(map[string]interface{})
	if resp.Status.Resources.GuestTools != nil {
		tools := resp.Status.Resources.GuestTools.NutanixGuestTools
		nutanixGuestTools["available_version"] = utils.StringValue(tools.AvailableVersion)
		nutanixGuestTools["iso_mount_state"] = utils.StringValue(tools.IsoMountState)
		nutanixGuestTools["state"] = utils.StringValue(tools.State)
		nutanixGuestTools["version"] = utils.StringValue(tools.Version)
		nutanixGuestTools["guest_os_version"] = utils.StringValue(tools.GuestOsVersion)
		nutanixGuestTools["enabled_capability_list"] = utils.StringValueSlice(tools.EnabledCapabilityList)
		nutanixGuestTools["vss_snapshot_capable"] = utils.BoolValue(tools.VSSSnapshotCapable)
		nutanixGuestTools["is_reachable"] = utils.BoolValue(tools.IsReachable)
		nutanixGuestTools["vm_mobility_drivers_installed"] = utils.BoolValue(tools.VMMobilityDriversInstalled)
	}
	if err := d.Set("nutanix_guest_tools", nutanixGuestTools); err != nil {
		return err
	}

	d.Set("num_vcpus_per_socket", utils.Int64Value(resp.Status.Resources.NumVcpusPerSocket))
	d.Set("num_sockets", utils.Int64Value(resp.Status.Resources.NumSockets))

	gpuList := make([]map[string]interface{}, 0)
	if resp.Status.Resources.GpuList != nil {
		gpuList = make([]map[string]interface{}, len(resp.Status.Resources.GpuList))
		for k, v := range resp.Status.Resources.GpuList {
			gpu := make(map[string]interface{})
			gpu["frame_buffer_size_mib"] = utils.Int64Value(v.FrameBufferSizeMib)
			gpu["vendor"] = utils.StringValue(v.Vendor)
			gpu["uuid"] = utils.StringValue(v.UUID)
			gpu["name"] = utils.StringValue(v.Name)
			gpu["pci_address"] = utils.StringValue(v.PCIAddress)
			gpu["fraction"] = utils.Int64Value(v.Fraction)
			gpu["mode"] = utils.StringValue(v.Mode)
			gpu["num_virtual_display_heads"] = utils.Int64Value(v.NumVirtualDisplayHeads)
			gpu["guest_driver_version"] = utils.StringValue(v.GuestDriverVersion)
			gpu["device_id"] = utils.Int64Value(v.DeviceID)

			gpuList[k] = gpu
		}
	}
	if err := d.Set("gpu_list", gpuList); err != nil {
		return err
	}
	if err := d.Set("parent_reference", getReferenceValues(resp.Status.Resources.ParentReference)); err != nil {
		return err
	}
	d.Set("memory_size_mib", utils.Int64Value(resp.Status.Resources.MemorySizeMib))

	bootDevice := make(map[string]interface{})
	if resp.Status.Resources.BootConfig != nil {
		boots := make([]string, len(resp.Status.Resources.BootConfig.BootDeviceOrderList))
		for k, v := range resp.Status.Resources.BootConfig.BootDeviceOrderList {
			boots[k] = utils.StringValue(v)
		}
		if err := d.Set("boot_device_order_list", boots); err != nil {
			return err
		}

		disk := make([]map[string]interface{}, 1)
		diskAddress := make(map[string]interface{})
		if resp.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
			diskAddress["device_index"] = utils.Int64Value(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)
			diskAddress["adapter_type"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
		}
		disk[0] = diskAddress

		bootDevice["disk_address"] = disk
		bootDevice["mac_address"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.MacAddress)
	}
	if err := d.Set("boot_device", bootDevice); err != nil {
		return err
	}

	d.Set("hardware_clock_timezone", utils.StringValue(resp.Status.Resources.HardwareClockTimezone))

	sysprep := make(map[string]interface{})
	cloudInit := make(map[string]interface{})
	isOv := false
	if resp.Status.Resources.GuestCustomization != nil {
		isOv = utils.BoolValue(resp.Status.Resources.GuestCustomization.IsOverridable)
		if resp.Status.Resources.GuestCustomization.CloudInit != nil {
			cloudInit["meta_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.MetaData)
			cloudInit["user_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.UserData)
			cloudInit["custom_key_values"] = resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues
		}
		if resp.Status.Resources.GuestCustomization.Sysprep != nil {
			sysprep["install_type"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.InstallType)
			sysprep["unattend_xml"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
			sysprep["custom_key_values"] = resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues
		}
	}
	if err := d.Set("guest_customization_sysprep", sysprep); err != nil {
		return err
	}
	if err := d.Set("guest_customization_cloud_init", cloudInit); err != nil {
		return err
	}

	d.Set("guest_customization_is_overridable", isOv)
	d.Set("should_fail_on_script_failure", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure))
	d.Set("enable_script_exec", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec))
	d.Set("power_state_mechanism", utils.StringValue(resp.Status.Resources.PowerStateMechanism.Mechanism))
	d.Set("vga_console_enabled", utils.BoolValue(resp.Status.Resources.VgaConsoleEnabled))

	diskList := make([]map[string]interface{}, 0)
	if resp.Status.Resources.DiskList != nil {
		diskList = make([]map[string]interface{}, len(resp.Status.Resources.DiskList))
		for k, v1 := range resp.Status.Resources.DiskList {
			disk := make(map[string]interface{})
			disk["uuid"] = utils.StringValue(v1.UUID)
			disk["disk_size_bytes"] = utils.Int64Value(v1.DiskSizeBytes)
			disk["disk_size_mib"] = utils.Int64Value(v1.DiskSizeMib)
			disk["data_source_reference"] = []map[string]interface{}{getReferenceValues(v1.DataSourceReference)}
			disk["volume_group_reference"] = []map[string]interface{}{getReferenceValues(v1.VolumeGroupReference)}

			dp := make([]map[string]interface{}, 1)
			deviceProps := make(map[string]interface{})
			deviceProps["device_type"] = utils.StringValue(v1.DeviceProperties.DeviceType)
			dp[0] = deviceProps

			da := make([]map[string]interface{}, 1)
			diskAddress := make(map[string]interface{})
			if v1.DeviceProperties.DiskAddress != nil {
				diskAddress["device_index"] = utils.Int64Value(v1.DeviceProperties.DiskAddress.DeviceIndex)
				diskAddress["adapter_type"] = utils.StringValue(v1.DeviceProperties.DiskAddress.AdapterType)
			}
			da[0] = diskAddress
			deviceProps["disk_address"] = da

			disk["device_properties"] = dp

			diskList[k] = disk
		}
	}

	d.SetId(*resp.Metadata.UUID)

	return d.Set("disk_list", diskList)
}

func getDataSourceVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vm_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"categories": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"api_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		// COMPUTED
		"message_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": {
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"hypervisor_type": {
			Type:     schema.TypeString,
			Computed: true,
		},

		// RESOURCES ARGUMENTS

		"num_vnuma_nodes": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"nic_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"floating_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network_function_nic_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"mac_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"ip_endpoint_list": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"subnet_reference": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"guest_os_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"power_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nutanix_guest_tools": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"available_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"iso_mount_state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guest_os_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"enabled_capability_list": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"vss_snapshot_capable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"is_reachable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"vm_mobility_drivers_installed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"num_vcpus_per_socket": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"num_sockets": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gpu_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frame_buffer_size_mib": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vendor": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pci_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"fraction": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"mode": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"num_virtual_display_heads": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guest_driver_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"parent_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"memory_size_mib": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"boot_device_order_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"boot_device": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_address": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_index": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"adapter_type": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"mac_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"hardware_clock_timezone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"guest_customization_cloud_init": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"meta_data": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"user_data": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"custom_key_values": {
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"guest_customization_is_overridable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"guest_customization_sysprep": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"install_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"unattend_xml": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"custom_key_values": {
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"should_fail_on_script_failure": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"enable_script_exec": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"power_state_mechanism": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vga_console_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"disk_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"disk_size_bytes": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"disk_size_mib": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"device_properties": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_type": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"disk_address": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_index": {
												Type:     schema.TypeInt,
												Computed: true,
											},
											"adapter_type": {
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"data_source_reference": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"volume_group_reference": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
