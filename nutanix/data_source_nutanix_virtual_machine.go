package nutanix

import (
	"fmt"
	"strconv"

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
	conn := meta.(*NutanixClient).API

	vm, ok := d.GetOk("vm_id")

	if !ok {
		return fmt.Errorf("please provide the required attribute vm_id")
	}

	// Make request to the API
	resp, err := conn.V3.GetVM(vm.(string))
	if err != nil {
		return err
	}

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
	metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
	metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
	metadata["creation_time"] = resp.Metadata.CreationTime.String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(resp.Metadata.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
	metadata["name"] = utils.StringValue(resp.Metadata.Name)
	if err := d.Set("metadata", metadata); err != nil {
		return err
	}
	if err := d.Set("categories", resp.Metadata.Categories); err != nil {
		return err
	}
	if resp.Metadata.ProjectReference != nil {
		pr := make(map[string]interface{})
		pr["kind"] = utils.StringValue(resp.Metadata.ProjectReference.Kind)
		pr["name"] = utils.StringValue(resp.Metadata.ProjectReference.Name)
		pr["uuid"] = utils.StringValue(resp.Metadata.ProjectReference.UUID)
		if err := d.Set("project_reference", pr); err != nil {
			return err
		}
	} else {
		if err := d.Set("project_reference", make(map[string]interface{})); err != nil {
			return err
		}
	}
	if resp.Metadata.OwnerReference != nil {
		or := make(map[string]interface{})
		or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
		or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
		or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)
		if err := d.Set("owner_reference", or); err != nil {
			return err
		}
	} else {
		if err := d.Set("owner_reference", make(map[string]interface{})); err != nil {
			return err
		}
	}
	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}
	if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
		return err
	}
	if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
		return err
	}
	// set availability zone reference values
	if resp.Status.AvailabilityZoneReference != nil {
		availabilityZoneReference := make(map[string]interface{})
		if resp.Status.AvailabilityZoneReference != nil {
			availabilityZoneReference["kind"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Kind)
			availabilityZoneReference["name"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Name)
			availabilityZoneReference["uuid"] = utils.StringValue(resp.Status.AvailabilityZoneReference.UUID)
		}
		if err := d.Set("availability_zone_reference", availabilityZoneReference); err != nil {
			return err
		}
	} else {
		if err := d.Set("availability_zone_reference", make(map[string]interface{})); err != nil {
			return err
		}
	}
	// set cluster reference values
	if resp.Status.ClusterReference != nil {
		clusterReference := make(map[string]interface{})
		clusterReference["kind"] = utils.StringValue(resp.Status.ClusterReference.Kind)
		clusterReference["name"] = utils.StringValue(resp.Status.ClusterReference.Name)
		clusterReference["uuid"] = utils.StringValue(resp.Status.ClusterReference.UUID)
		if err := d.Set("cluster_reference", clusterReference); err != nil {
			return err
		}
	} else {
		if err := d.Set("cluster_reference", make(map[string]interface{})); err != nil {
			return err
		}
	}
	// set message list values
	if resp.Status.MessageList != nil {
		messages := make([]map[string]interface{}, len(resp.Status.MessageList))
		for k, v := range resp.Status.MessageList {
			message := make(map[string]interface{})
			message["message"] = utils.StringValue(v.Message)
			message["reason"] = utils.StringValue(v.Reason)
			message["details"] = v.Details
			messages[k] = message
		}
		if err := d.Set("message_list", messages); err != nil {
			return err
		}
	}
	// set state value
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}
	// set vnuma_config value
	if err := d.Set("num_vnuma_nodes", utils.Int64Value(resp.Status.Resources.VnumaConfig.NumVnumaNodes)); err != nil {
		return err
	}
	// set nic list value
	nics := resp.Status.Resources.NicList
	if nics != nil {
		nicLists := make([]map[string]interface{}, len(nics))
		for k, v := range nics {
			nic := make(map[string]interface{})
			// simple firts
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["floating_ip"] = utils.StringValue(v.FloatingIP)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)

			// set ip lists value
			ipEndpointList := make([]map[string]interface{}, len(v.IPEndpointList))
			for k1, v1 := range v.IPEndpointList {
				ipEndpoint := make(map[string]interface{})
				ipEndpoint["ip"] = utils.StringValue(v1.IP)
				ipEndpoint["type"] = utils.StringValue(v1.Type)
				ipEndpointList[k1] = ipEndpoint
			}
			nic["ip_endpoint_list"] = ipEndpointList

			// set network_function_chain_reference value
			netFnChainRef := make(map[string]interface{})
			if v.NetworkFunctionChainReference != nil {
				netFnChainRef["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
				netFnChainRef["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
				netFnChainRef["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
			}
			nic["network_function_chain_reference"] = netFnChainRef

			// set subnet_reference value
			subtnetRef := make(map[string]interface{})
			if v.SubnetReference != nil {
				subtnetRef["kind"] = utils.StringValue(v.SubnetReference.Kind)
				subtnetRef["name"] = utils.StringValue(v.SubnetReference.Name)
				subtnetRef["uuid"] = utils.StringValue(v.SubnetReference.UUID)
			}
			nic["subnet_reference"] = subtnetRef

			nicLists[k] = nic
		}
		if err := d.Set("nic_list", nicLists); err != nil {
			return err
		}
	} else {
		if err := d.Set("nic_list", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
	}
	// set host_reference value
	if resp.Status.Resources.HostReference != nil {
		hostRef := make(map[string]interface{})
		hostRef["kind"] = utils.StringValue(resp.Status.Resources.HostReference.Kind)
		hostRef["name"] = utils.StringValue(resp.Status.Resources.HostReference.Name)
		hostRef["uuid"] = utils.StringValue(resp.Status.Resources.HostReference.UUID)
		if err := d.Set("host_reference", hostRef); err != nil {
			return err
		}
	} else {
		if err := d.Set("host_reference", make(map[string]interface{})); err != nil {
			return err
		}
	}
	// set guest_os_id value
	if err := d.Set("guest_os_id", resp.Status.Resources.GuestOsID); err != nil {
		return err
	}
	// set power_state value
	if err := d.Set("power_state", resp.Status.Resources.PowerState); err != nil {
		return err
	}

	if resp.Status.Resources.GuestTools != nil {
		tools := resp.Status.Resources.GuestTools.NutanixGuestTools
		nutanixGuestTools := make(map[string]interface{})
		nutanixGuestTools["available_version"] = utils.StringValue(tools.AvailableVersion)
		nutanixGuestTools["iso_mount_state"] = utils.StringValue(tools.IsoMountState)
		nutanixGuestTools["state"] = utils.StringValue(tools.State)
		nutanixGuestTools["version"] = utils.StringValue(tools.Version)
		nutanixGuestTools["guest_os_version"] = utils.StringValue(tools.GuestOsVersion)

		capList := make([]string, len(tools.EnabledCapabilityList))
		for k, v := range tools.EnabledCapabilityList {
			capList[k] = *v
		}
		nutanixGuestTools["enabled_capability_list"] = capList
		nutanixGuestTools["vss_snapshot_capable"] = utils.BoolValue(tools.VSSSnapshotCapable)
		nutanixGuestTools["is_reachable"] = utils.BoolValue(tools.IsReachable)
		nutanixGuestTools["vm_mobility_drivers_installed"] = utils.BoolValue(tools.VMMobilityDriversInstalled)

		// set nutanix_guest_tools value
		if err := d.Set("nutanix_guest_tools", nutanixGuestTools); err != nil {
			return err
		}
	} else {
		if err := d.Set("nutanix_guest_tools", make(map[string]interface{})); err != nil {
			return err
		}
	}
	// set num_vcpus_per_socket value
	if err := d.Set("num_vcpus_per_socket", resp.Status.Resources.NumVcpusPerSocket); err != nil {
		return err
	}
	// set num_sockets value
	if err := d.Set("num_sockets", resp.Status.Resources.NumSockets); err != nil {
		return err
	}
	if resp.Status.Resources.GpuList != nil {
		gpuList := make([]map[string]interface{}, len(resp.Status.Resources.GpuList))
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
		// set gpu_list value
		if err := d.Set("gpu_list", gpuList); err != nil {
			return err
		}
	} else {
		if err := d.Set("gpu_list", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
	}

	if resp.Status.Resources.ParentReference != nil {
		parentRef := make(map[string]interface{})
		parentRef["kind"] = utils.StringValue(resp.Status.Resources.ParentReference.Kind)
		parentRef["name"] = utils.StringValue(resp.Status.Resources.ParentReference.Name)
		parentRef["uuid"] = utils.StringValue(resp.Status.Resources.ParentReference.UUID)
		// set parent_reference value
		if err := d.Set("parent_reference", parentRef); err != nil {
			return err
		}
	}
	if err := d.Set("parent_reference", make(map[string]interface{})); err != nil {
		return err
	}
	// set memory_size_mib value
	if err := d.Set("memory_size_mib", resp.Status.Resources.MemorySizeMib); err != nil {
		return err
	}
	if resp.Status.Resources.BootConfig != nil {
		boots := make([]string, len(resp.Status.Resources.BootConfig.BootDeviceOrderList))
		for k, v := range resp.Status.Resources.BootConfig.BootDeviceOrderList {
			boots[k] = utils.StringValue(v)
		}
		// set boot_device_order_list value
		if err := d.Set("boot_device_order_list", boots); err != nil {
			return err
		}

		bootDevice := make(map[string]interface{})
		disk := make([]map[string]interface{}, 1)
		diskAddress := make(map[string]interface{})
		if resp.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
			diskAddress["device_index"] = utils.Int64Value(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)
			diskAddress["adapter_type"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
		}
		disk[0] = diskAddress

		bootDevice["disk_address"] = disk
		bootDevice["mac_address"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.MacAddress)
		// set boot_device value
		if err := d.Set("boot_device", bootDevice); err != nil {
			return err
		}
	} else {
		if err := d.Set("boot_device", make(map[string]interface{})); err != nil {
			return err
		}
	}
	// set hardware_clock_timezone value
	if err := d.Set("hardware_clock_timezone", resp.Status.Resources.HardwareClockTimezone); err != nil {
		return err
	}
	if resp.Status.Resources.GuestCustomization != nil {
		cloudInit := make(map[string]interface{})
		if resp.Status.Resources.GuestCustomization.CloudInit != nil {
			cloudInit["meta_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.MetaData)
			cloudInit["user_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.UserData)
			cloudInit["custom_key_values"] = resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues
		}
		// set guest_customization_cloud_init value
		if err := d.Set("guest_customization_cloud_init", cloudInit); err != nil {
			return err
		}
		// set guest_customization_is_overridable value
		if err := d.Set("guest_customization_is_overridable", utils.BoolValue(resp.Status.Resources.GuestCustomization.IsOverridable)); err != nil {
			return err
		}
		sysprep := make(map[string]interface{})
		if resp.Status.Resources.GuestCustomization.Sysprep != nil {
			sysprep["install_type"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.InstallType)
			sysprep["unattend_xml"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
			sysprep["custom_key_values"] = resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues
		}
		// set guest_customization_sysprep value
		if err := d.Set("guest_customization_sysprep", sysprep); err != nil {
			return err
		}
	}

	// set power_state_guest_transition_config value
	if err := d.Set("should_fail_on_script_failure", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure)); err != nil {
		return err
	}
	if err := d.Set("enable_script_exec", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec)); err != nil {
		return err
	}
	// set power_state_mechanism value
	if err := d.Set("power_state_mechanism", utils.StringValue(resp.Status.Resources.PowerStateMechanism.Mechanism)); err != nil {
		return err
	}
	// set vga_console_enabled value
	if err := d.Set("vga_console_enabled", utils.BoolValue(resp.Status.Resources.VgaConsoleEnabled)); err != nil {
		return err
	}
	if resp.Status.Resources.DiskList != nil {
		diskList := make([]map[string]interface{}, len(resp.Status.Resources.DiskList))
		for k, v := range resp.Status.Resources.DiskList {
			disk := make(map[string]interface{})
			disk["uuid"] = *v.UUID
			disk["disk_size_bytes"] = *v.DiskSizeBytes
			disk["disk_size_mib"] = *v.DiskSizeMib

			ds := make([]map[string]interface{}, 1)
			dsourceRef := make(map[string]interface{})
			if v.DataSourceReference != nil {
				dsourceRef["kind"] = utils.StringValue(v.DataSourceReference.Kind)
				dsourceRef["name"] = utils.StringValue(v.DataSourceReference.Name)
				dsourceRef["uuid"] = utils.StringValue(v.DataSourceReference.UUID)
			}
			ds[0] = dsourceRef

			disk["data_source_reference"] = ds

			vr := make([]map[string]interface{}, 1)
			volumeRef := make(map[string]interface{})
			if v.VolumeGroupReference != nil {
				volumeRef["kind"] = utils.StringValue(v.VolumeGroupReference.Kind)
				volumeRef["name"] = utils.StringValue(v.VolumeGroupReference.Name)
				volumeRef["uuid"] = utils.StringValue(v.VolumeGroupReference.UUID)
			}
			vr[0] = volumeRef

			disk["volume_group_reference"] = vr

			dp := make([]map[string]interface{}, 1)
			deviceProps := make(map[string]interface{})
			deviceProps["device_type"] = utils.StringValue(v.DeviceProperties.DeviceType)
			dp[0] = deviceProps

			da := make([]map[string]interface{}, 1)
			diskAddress := make(map[string]interface{})
			if v.DeviceProperties.DiskAddress != nil {
				diskAddress["device_index"] = utils.Int64Value(v.DeviceProperties.DiskAddress.DeviceIndex)
				diskAddress["adapter_type"] = utils.StringValue(v.DeviceProperties.DiskAddress.AdapterType)
			}
			da[0] = diskAddress
			deviceProps["disk_address"] = da

			disk["device_properties"] = dp

			diskList[k] = disk
		}
		// set disk_list value
		if err := d.Set("disk_list", diskList); err != nil {
			return err
		}
	} else {
		if err := d.Set("disk_list", make([]map[string]interface{}, 0)); err != nil {
			return err
		}
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func getDataSourceVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vm_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"metadata": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"categories": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
		},
		"project_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"owner_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		// COMPUTED
		"message_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"hypervisor_type": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},

		// RESOURCES ARGUMENTS

		"num_vnuma_nodes": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"nic_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"floating_ip": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"network_function_nic_type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"mac_address": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"ip_endpoint_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"network_function_chain_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"subnet_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"guest_os_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"power_state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"nutanix_guest_tools": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"available_version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"iso_mount_state": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"guest_os_version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"enabled_capability_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"vss_snapshot_capable": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: true,
					},
					"is_reachable": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: true,
					},
					"vm_mobility_drivers_installed": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"num_vcpus_per_socket": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"num_sockets": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gpu_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frame_buffer_size_mib": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vendor": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"pci_address": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"fraction": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
					"mode": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"num_virtual_display_heads": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guest_driver_version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_id": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"parent_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"memory_size_mib": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"boot_device_order_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"boot_device": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_address": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_index": &schema.Schema{
									Type:     schema.TypeInt,
									Computed: true,
								},
								"adapter_type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"mac_address": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"hardware_clock_timezone": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"guest_customization_cloud_init": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"meta_data": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"user_data": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"custom_key_values": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"guest_customization_is_overridable": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"guest_customization_sysprep": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"install_type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"unattend_xml": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"custom_key_values": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"should_fail_on_script_failure": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"enable_script_exec": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"power_state_mechanism": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"vga_console_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"disk_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"disk_size_bytes": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
					"disk_size_mib": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
					"device_properties": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"disk_address": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_index": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
											"adapter_type": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"data_source_reference": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"volume_group_reference": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
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
