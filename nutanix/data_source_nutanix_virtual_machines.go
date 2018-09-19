package nutanix

import (
	"strconv"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	// CDROM ...
	CDROM = "CDROM"
)

func dataSourceNutanixVirtualMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixVirtualMachinesRead,

		Schema: map[string]*schema.Schema{
			//"metadata": getDSMetadataSchema(),
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
								},
							},
						},
						"cluster_reference_name": {
							Type:     schema.TypeString,
							Computed: true,
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
									"subnet_reference_name": {
										Type:     schema.TypeString,
										Computed: true,
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
						"boot_device_disk_address": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_index": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"adapter_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"boot_device_mac_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"hardware_clock_timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guest_customization_cloud_init_meta_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guest_customization_cloud_init_user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guest_customization_cloud_init_custom_key_values": {
							Type:     schema.TypeMap,
							Computed: true,
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
								},
							},
						},
						"guest_customization_sysprep_custom_key_values": {
							Type:     schema.TypeMap,
							Computed: true,
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
					},
				},
			},
		},
	}
}

func dataSourceNutanixVirtualMachinesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	resp, err := conn.V3.ListAllVM()
	if err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))
	for k, v := range resp.Entities {
		entity := make(map[string]interface{})

		m, c := setRSEntityMetadata(v.Metadata)

		entity["metadata"] = m
		entity["project_reference"] = getReferenceValues(v.Metadata.ProjectReference)
		entity["owner_reference"] = getReferenceValues(v.Metadata.OwnerReference)
		entity["categories"] = c
		entity["api_version"] = utils.StringValue(v.APIVersion)
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["description"] = utils.StringValue(v.Status.Description)
		entity["availability_zone_reference"] = getReferenceValues(v.Status.AvailabilityZoneReference)
		entity["cluster_reference"] = getClusterReferenceValues(v.Status.ClusterReference)
		entity["cluster_reference_name"] = utils.StringValue(v.Status.ClusterReference.Name)
		entity["state"] = utils.StringValue(v.Status.State)
		entity["num_vnuma_nodes"] = utils.Int64Value(v.Status.Resources.VnumaConfig.NumVnumaNodes)
		entity["nic_list"] = flattenNicList(v.Status.Resources.NicList)
		entity["host_reference"] = getReferenceValues(v.Status.Resources.HostReference)
		entity["guest_os_id"] = utils.StringValue(v.Status.Resources.GuestOsID)
		entity["power_state"] = utils.StringValue(v.Status.Resources.PowerState)
		entity["nutanix_guest_tools"] = setNutanixGuestTools(v.Status.Resources.GuestTools)
		entity["num_vcpus_per_socket"] = utils.Int64Value(v.Status.Resources.NumVcpusPerSocket)
		entity["num_sockets"] = utils.Int64Value(v.Status.Resources.NumSockets)
		entity["gpu_list"] = flattenGPUList(v.Status.Resources.GpuList)
		entity["parent_reference"] = getReferenceValues(v.Status.Resources.ParentReference)
		entity["memory_size_mib"] = utils.Int64Value(v.Status.Resources.MemorySizeMib)

		diskAddress := make(map[string]interface{})
		mac := ""
		b := make([]string, 0)

		if v.Status.Resources.BootConfig != nil {
			if v.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
				i := strconv.Itoa(int(utils.Int64Value(v.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)))
				diskAddress["device_index"] = i
				diskAddress["adapter_type"] = utils.StringValue(v.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
			}
			if v.Status.Resources.BootConfig.BootDeviceOrderList != nil {
				b = utils.StringValueSlice(v.Status.Resources.BootConfig.BootDeviceOrderList)
			}
			mac = utils.StringValue(v.Status.Resources.BootConfig.BootDevice.MacAddress)
		}

		entity["boot_device_order_list"] = b
		entity["boot_device_disk_address"] = diskAddress
		entity["boot_device_mac_address"] = mac
		entity["hardware_clock_timezone"] = utils.StringValue(v.Status.Resources.HardwareClockTimezone)

		cloudInitUser := ""
		cloudInitMeta := ""
		cloudInitCV := make(map[string]string)
		sysprep := make(map[string]interface{})
		sysprepCV := make(map[string]string)
		isOv := false
		if v.Status.Resources.GuestCustomization != nil {
			isOv = utils.BoolValue(v.Status.Resources.GuestCustomization.IsOverridable)

			if v.Status.Resources.GuestCustomization.CloudInit != nil {
				cloudInitMeta = utils.StringValue(v.Status.Resources.GuestCustomization.CloudInit.MetaData)
				cloudInitUser = utils.StringValue(v.Status.Resources.GuestCustomization.CloudInit.UserData)
				if v.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues != nil {
					for k, v := range v.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues {
						cloudInitCV[k] = v
					}
				}
			}
			if v.Status.Resources.GuestCustomization.Sysprep != nil {
				sysprep["install_type"] = utils.StringValue(v.Status.Resources.GuestCustomization.Sysprep.InstallType)
				sysprep["unattend_xml"] = utils.StringValue(v.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
				if v.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues != nil {
					for k, v := range v.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues {
						sysprepCV[k] = v
					}
				}
			}
		}

		entity["guest_customization_cloud_init_custom_key_values"] = cloudInitCV
		entity["guest_customization_sysprep_custom_key_values"] = sysprepCV
		entity["guest_customization_is_overridable"] = isOv
		entity["guest_customization_cloud_init_user_data"] = cloudInitUser
		entity["guest_customization_cloud_init_meta_data"] = cloudInitMeta
		entity["guest_customization_sysprep"] = sysprep
		entity["should_fail_on_script_failure"] = utils.BoolValue(
			v.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure)
		entity["enable_script_exec"] = utils.BoolValue(v.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec)
		entity["power_state_mechanism"] = utils.StringValue(v.Status.Resources.PowerStateMechanism.Mechanism)
		entity["vga_console_enabled"] = utils.BoolValue(v.Status.Resources.VgaConsoleEnabled)
		entity["disk_list"] = setDiskList(v.Status.Resources.DiskList, v.Status.Resources.GuestCustomization)

		entities[k] = entity
	}

	d.SetId(resource.UniqueId())
	d.Set("api_version", utils.StringValue(resp.APIVersion))

	return d.Set("entities", entities)
}

func setDiskList(disk []*v3.VMDisk, hasCloudInit *v3.GuestCustomizationStatus) []map[string]interface{} {
	var diskList []map[string]interface{}
	if len(disk) > 0 {
		for _, v1 := range disk {

			if hasCloudInit != nil {
				if hasCloudInit.CloudInit != nil && utils.StringValue(v1.DeviceProperties.DeviceType) == CDROM {
					continue
				}
			}

			disk := make(map[string]interface{})
			disk["uuid"] = utils.StringValue(v1.UUID)
			disk["disk_size_bytes"] = utils.Int64Value(v1.DiskSizeBytes)
			disk["disk_size_mib"] = utils.Int64Value(v1.DiskSizeMib)
			disk["data_source_reference"] = []map[string]interface{}{getClusterReferenceValues(v1.DataSourceReference)}
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

			diskList = append(diskList, disk)
		}
	}

	if diskList == nil {
		return make([]map[string]interface{}, 0)
	}

	return diskList
}

func flattenGPUList(gpu []*v3.VMGpuOutputStatus) []map[string]interface{} {
	gpuList := make([]map[string]interface{}, 0)
	if gpu != nil {
		gpuList = make([]map[string]interface{}, len(gpu))
		for k, v := range gpu {
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
	return gpuList
}

func setNutanixGuestTools(guest *v3.GuestToolsStatus) map[string]interface{} {
	nutanixGuestTools := make(map[string]interface{})
	if guest != nil {
		tools := guest.NutanixGuestTools
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
	return nutanixGuestTools
}

func flattenNicList(nics []*v3.VMNicOutputStatus) []map[string]interface{} {
	nicLists := make([]map[string]interface{}, 0)
	if nics != nil {
		nicLists = make([]map[string]interface{}, len(nics))
		for k, v := range nics {
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
			nic["subnet_reference"] = getClusterReferenceValues(v.SubnetReference)
			nic["subnet_reference_name"] = utils.StringValue(v.SubnetReference.Name)

			nicLists[k] = nic
		}
	}

	return nicLists
}

func getDSMetadataSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"kind": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sort_attribute": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"filter": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"length": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sort_order": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"offset": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}
