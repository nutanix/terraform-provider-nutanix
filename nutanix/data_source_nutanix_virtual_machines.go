package nutanix

import (
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixVirtualMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixVirtualMachinesRead,

		Schema: getDataSourceVMSSchema(),
	}
}

func dataSourceNutanixVirtualMachinesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*NutanixClient).API

	metadata := &v3.VMListMetadata{}

	if v, ok := d.GetOk("metadata"); ok {
		m := v.(map[string]interface{})
		metadata.Kind = utils.String("vm")
		if mv, mok := m["sort_attribute"]; mok {
			metadata.SortAttribute = utils.String(mv.(string))
		}
		if mv, mok := m["filter"]; mok {
			metadata.Filter = utils.String(mv.(string))
		}
		if mv, mok := m["length"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Length = utils.Int64(int64(i))
		}
		if mv, mok := m["sort_order"]; mok {
			metadata.SortOrder = utils.String(mv.(string))
		}
		if mv, mok := m["offset"]; mok {
			i, err := strconv.Atoi(mv.(string))
			if err != nil {
				return err
			}
			metadata.Offset = utils.Int64(int64(i))
		}
	}

	// Make request to the API
	resp, err := conn.V3.ListVM(metadata)
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})
		// set metadata values
		metadata := make(map[string]interface{})
		metadata["last_update_time"] = utils.TimeValue(v.Metadata.LastUpdateTime).String()
		metadata["kind"] = utils.StringValue(v.Metadata.Kind)
		metadata["uuid"] = utils.StringValue(v.Metadata.UUID)
		metadata["creation_time"] = utils.TimeValue(v.Metadata.CreationTime).String()
		metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(v.Metadata.SpecVersion)))
		metadata["spec_hash"] = utils.StringValue(v.Metadata.SpecHash)
		metadata["name"] = utils.StringValue(v.Metadata.Name)
		entity["metadata"] = metadata

		entity["categories"] = v.Metadata.Categories
		entity["api_version"] = utils.StringValue(v.APIVersion)

		pr := make(map[string]interface{})
		pr["kind"] = utils.StringValue(v.Metadata.ProjectReference.Kind)
		pr["name"] = utils.StringValue(v.Metadata.ProjectReference.Name)
		pr["uuid"] = utils.StringValue(v.Metadata.ProjectReference.UUID)

		entity["project_reference"] = pr

		or := make(map[string]interface{})
		or["kind"] = utils.StringValue(v.Metadata.OwnerReference.Kind)
		or["name"] = utils.StringValue(v.Metadata.OwnerReference.Name)
		or["uuid"] = utils.StringValue(v.Metadata.OwnerReference.UUID)
		entity["owner_reference"] = or
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)

		// set availability zone reference values
		availabilityZoneReference := make(map[string]interface{})
		if v.Status.AvailabilityZoneReference != nil {
			availabilityZoneReference["kind"] = utils.StringValue(v.Status.AvailabilityZoneReference.Kind)
			availabilityZoneReference["name"] = utils.StringValue(v.Status.AvailabilityZoneReference.Name)
			availabilityZoneReference["uuid"] = utils.StringValue(v.Status.AvailabilityZoneReference.UUID)
		}
		entity["availability_zone_reference"] = availabilityZoneReference
		// set cluster reference values
		clusterReference := make(map[string]interface{})
		clusterReference["kind"] = utils.StringValue(v.Status.ClusterReference.Kind)
		clusterReference["name"] = utils.StringValue(v.Status.ClusterReference.Name)
		clusterReference["uuid"] = utils.StringValue(v.Status.ClusterReference.UUID)
		entity["cluster_reference"] = availabilityZoneReference
		entity["state"] = *v.Status.State
		entity["num_vnuma_nodes"] = utils.Int64Value(v.Status.Resources.VnumaConfig.NumVnumaNodes)

		// set nic list value
		nics := v.Status.Resources.NicList
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

				entity["ip_address"] = utils.StringValue(v.IPEndpointList[0].IP)

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
			entity["nic_list"] = nicLists
		} else {
			entity["nic_list"] = make([]map[string]interface{}, 0)
		}
		// set host_reference value
		hostRef := make(map[string]interface{})
		if v.Status.Resources.HostReference != nil {
			hostRef["kind"] = utils.StringValue(v.Status.Resources.HostReference.Kind)
			hostRef["name"] = utils.StringValue(v.Status.Resources.HostReference.Name)
			hostRef["uuid"] = utils.StringValue(v.Status.Resources.HostReference.UUID)
		}
		entity["host_reference"] = hostRef
		entity["guest_os_id"] = utils.StringValue(v.Status.Resources.GuestOsID)
		entity["power_state"] = utils.StringValue(v.Status.Resources.PowerState)

		if v.Status.Resources.GuestTools != nil {
			tools := v.Status.Resources.GuestTools.NutanixGuestTools
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
			entity["nutanix_guest_tools"] = nutanixGuestTools
		} else {
			entity["nutanix_guest_tools"] = make(map[string]interface{})
		}
		entity["num_vcpus_per_socket"] = utils.Int64Value(v.Status.Resources.NumVcpusPerSocket)
		entity["num_sockets"] = utils.Int64Value(v.Status.Resources.NumSockets)
		if v.Status.Resources.GpuList != nil {
			gpuList := make([]map[string]interface{}, len(v.Status.Resources.GpuList))
			for k, v := range v.Status.Resources.GpuList {
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
			entity["gpu_list"] = gpuList

		} else {
			entity["gpu_list"] = make([]map[string]interface{}, 0)
		}
		if v.Status.Resources.ParentReference != nil {
			parentRef := make(map[string]interface{})
			parentRef["kind"] = utils.StringValue(v.Status.Resources.ParentReference.Kind)
			parentRef["name"] = utils.StringValue(v.Status.Resources.ParentReference.Name)
			parentRef["uuid"] = utils.StringValue(v.Status.Resources.ParentReference.UUID)
			// set parent_reference value
			entity["parent_reference"] = parentRef
		} else {
			entity["parent_reference"] = make(map[string]interface{})
		}
		entity["memory_size_mib"] = utils.Int64Value(v.Status.Resources.MemorySizeMib)

		if v.Status.Resources.BootConfig != nil {
			boots := make([]string, len(v.Status.Resources.BootConfig.BootDeviceOrderList))
			for k, v := range v.Status.Resources.BootConfig.BootDeviceOrderList {
				boots[k] = utils.StringValue(v)
			}
			// set boot_device_order_list value
			entity["boot_device_order_list"] = boots

			bootDevice := make(map[string]interface{})
			disk := make([]map[string]interface{}, 1)
			diskAddress := make(map[string]interface{})
			if v.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
				diskAddress["device_index"] = utils.Int64Value(v.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)
				diskAddress["adapter_type"] = utils.StringValue(v.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
			}
			disk[0] = diskAddress

			bootDevice["disk_address"] = disk
			bootDevice["mac_address"] = utils.StringValue(v.Status.Resources.BootConfig.BootDevice.MacAddress)
			// set boot_device value
			entity["boot_device"] = bootDevice
		} else {
			entity["boot_device_order_list"] = make([]string, 0)
			entity["boot_device"] = make(map[string]interface{})
		}
		entity["hardware_clock_timezone"] = utils.StringValue(v.Status.Resources.HardwareClockTimezone)

		if v.Status.Resources.GuestCustomization != nil {
			cloudInit := make(map[string]interface{})
			if v.Status.Resources.GuestCustomization.CloudInit != nil {
				cloudInit["meta_data"] = utils.StringValue(v.Status.Resources.GuestCustomization.CloudInit.MetaData)
				cloudInit["user_data"] = utils.StringValue(v.Status.Resources.GuestCustomization.CloudInit.UserData)
				cloudInit["custom_key_values"] = v.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues
			}
			// set guest_customization_cloud_init value
			entity["guest_customization_cloud_init"] = cloudInit
			entity["guest_customization_is_overridable"] = utils.BoolValue(v.Status.Resources.GuestCustomization.IsOverridable)

			sysprep := make(map[string]interface{})
			if v.Status.Resources.GuestCustomization.Sysprep != nil {
				sysprep["install_type"] = utils.StringValue(v.Status.Resources.GuestCustomization.Sysprep.InstallType)
				sysprep["unattend_xml"] = utils.StringValue(v.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
				sysprep["custom_key_values"] = v.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues
			}
			// set guest_customization_sysprep value
			entity["guest_customization_sysprep"] = sysprep
			entity["should_fail_on_script_failure"] = utils.BoolValue(v.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure)
			entity["enable_script_exec"] = utils.BoolValue(v.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec)
			entity["power_state_mechanism"] = utils.StringValue(v.Status.Resources.PowerStateMechanism.Mechanism)
			entity["vga_console_enabled"] = utils.BoolValue(v.Status.Resources.VgaConsoleEnabled)

			if v.Status.Resources.DiskList != nil {
				diskList := make([]map[string]interface{}, len(v.Status.Resources.DiskList))
				for k, v := range v.Status.Resources.DiskList {
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
				entity["disk_list"] = diskList
			} else {
				entity["disk_list"] = make([]map[string]interface{}, 0)
			}
		}
		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}

func getDataSourceVMSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"sort_attribute": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"filter": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"length": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"sort_order": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"offset": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"entities": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}
