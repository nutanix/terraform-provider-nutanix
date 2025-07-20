package vmm

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixOVAImage() *schema.Resource {
	return &schema.Resource{
		ReadContext:   DataSourceNutanixOVAImageRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"ova_image_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ova_image_name"},
			},
			"ova_image_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ova_image_id"},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hypervisor_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// RESOURCES ARGUMENTS
			"enable_cpu_passthrough": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_vcpu_hard_pinned": {
				Type:     schema.TypeBool,
				Computed: true,
			},
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"num_queues": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"subnet_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_connected": {
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ngt_enabled_capability_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ngt_credentials": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"num_vcpus_per_socket": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parent_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"boot_device_mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"boot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
						"storage_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flash_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage_container_reference": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:     schema.TypeString,
													Computed: true,
												},
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
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"data_source_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"volume_group_reference": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"serial_port_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceNutanixOVAImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*conns.Client).API

	imageID, iok := d.GetOk("ova_image_id")
	imageName, nok := d.GetOk("ova_image_name")

	if !iok && !nok {
		return diag.Errorf("please provide one of ova_image_id or ova_image_name attributes")
	}

	var reqErr error
	var resp *v3.OVAImageIntentResponse

	if nok {
		imageID = findOVAImageByName(ctx, conn, imageName.(string))
	}
	resp, reqErr = findOVAImageByUUID(conn, imageID.(string))

	if reqErr != nil {
		return diag.FromErr(reqErr)
	}

	m, c := setRSEntityMetadata(resp.VMSpec.Metadata)
	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("categories", c); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("project_reference", flattenReferenceValues(resp.VMSpec.Metadata.ProjectReference)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("availability_zone_reference", flattenReferenceValues(resp.VMSpec.Spec.AvailabilityZoneReference)); err != nil {
		return diag.FromErr(err)
	}

	if err := flattenClusterReference(resp.VMSpec.Spec.ClusterReference, d); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nic_list", flattenNicList(resp.VMSpec.Spec.Resources.NicList)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_reference", flattenReferenceValues(resp.VMSpec.Spec.Resources.ParentReference)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("disk_list", flattenDiskListHelper(resp.VMSpec.Spec.Resources.DiskList, "", false)); err != nil {
		return diag.FromErr(err)
	}

	diskAddress := make(map[string]interface{})
	mac := ""
	bootType := ""
	machineType := ""
	b := make([]string, 0)

	if resp.VMSpec.Spec.Resources.BootConfig != nil {
		if resp.VMSpec.Spec.Resources.BootConfig.BootDevice != nil {
			if resp.VMSpec.Spec.Resources.BootConfig.BootDevice.DiskAddress != nil {
				i := strconv.Itoa(int(utils.Int64Value(resp.VMSpec.Spec.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)))
				diskAddress["device_index"] = i
				diskAddress["adapter_type"] = utils.StringValue(resp.VMSpec.Spec.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
			}
			mac = utils.StringValue(resp.VMSpec.Spec.Resources.BootConfig.BootDevice.MacAddress)
		}
		if resp.VMSpec.Spec.Resources.BootConfig.BootDeviceOrderList != nil {
			b = utils.StringValueSlice(resp.VMSpec.Spec.Resources.BootConfig.BootDeviceOrderList)
		}
		if resp.VMSpec.Spec.Resources.BootConfig.BootType != nil {
			bootType = utils.StringValue(resp.VMSpec.Spec.Resources.BootConfig.BootType)
		}
	} else {
	}

	if resp.VMSpec.Spec.Resources.MachineType != nil {
		machineType = utils.StringValue(resp.VMSpec.Spec.Resources.MachineType)
	} else {
	}

	d.Set("boot_device_order_list", b)

	d.Set("boot_device_disk_address", diskAddress)

	d.Set("boot_device_mac_address", mac)

	d.Set("boot_type", bootType)

	d.Set("machine_type", machineType)

	sysprep := make(map[string]interface{})
	sysrepCV := make(map[string]string)
	cloudInitUser := ""
	cloudInitMeta := ""
	cloudInitCV := make(map[string]string)
	isOv := false

	if resp.VMSpec.Spec.Resources.GuestCustomization != nil {
		isOv = utils.BoolValue(resp.VMSpec.Spec.Resources.GuestCustomization.IsOverridable)

		if resp.VMSpec.Spec.Resources.GuestCustomization.CloudInit != nil {
			cloudInitMeta = utils.StringValue(resp.VMSpec.Spec.Resources.GuestCustomization.CloudInit.MetaData)
			cloudInitUser = utils.StringValue(resp.VMSpec.Spec.Resources.GuestCustomization.CloudInit.UserData)

			if resp.VMSpec.Spec.Resources.GuestCustomization.CloudInit.CustomKeyValues != nil {
				for k, v := range resp.VMSpec.Spec.Resources.GuestCustomization.CloudInit.CustomKeyValues {
					cloudInitCV[k] = v
				}
			} else {
			}
		} else {
		}

		if resp.VMSpec.Spec.Resources.GuestCustomization.Sysprep != nil {
			sysprep["install_type"] = utils.StringValue(resp.VMSpec.Spec.Resources.GuestCustomization.Sysprep.InstallType)
			sysprep["unattend_xml"] = utils.StringValue(resp.VMSpec.Spec.Resources.GuestCustomization.Sysprep.UnattendXML)

			if resp.VMSpec.Spec.Resources.GuestCustomization.Sysprep.CustomKeyValues != nil {
				for k, v := range resp.VMSpec.Spec.Resources.GuestCustomization.Sysprep.CustomKeyValues {
					sysrepCV[k] = v
				}
			} else {
			}
		} else {
		}
	} else {
	}

	if err := d.Set("guest_customization_cloud_init_custom_key_values", cloudInitCV); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("guest_customization_sysprep_custom_key_values", sysrepCV); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("guest_customization_sysprep", sysprep); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("serial_port_list", resp.VMSpec.Spec.Resources.SerialPortList); err != nil {
		return diag.FromErr(err)
	}

	d.Set("guest_customization_cloud_init_user_data", cloudInitUser)

	d.Set("guest_customization_cloud_init_meta_data", cloudInitMeta)

	d.Set("hardware_clock_timezone", utils.StringValue(resp.VMSpec.Spec.Resources.HardwareClockTimezone))

	d.Set("api_version", utils.StringValue(resp.VMSpec.APIVersion))

	d.Set("name", utils.StringValue(resp.VMSpec.Spec.Name))

	d.Set("description", utils.StringValue(resp.VMSpec.Spec.Description))

	d.Set("state", nil) // Setting nil, consider if this is desired.

	d.Set("enable_cpu_passthrough", utils.BoolValue(resp.VMSpec.Spec.Resources.EnableCPUPassthrough))

	d.Set("is_vcpu_hard_pinned", utils.BoolValue(resp.VMSpec.Spec.Resources.EnableCPUPinning))

	d.Set("guest_os_id", utils.StringValue(resp.VMSpec.Spec.Resources.GuestOsID))

	d.Set("power_state", utils.StringValue(resp.VMSpec.Spec.Resources.PowerState))

	d.Set("num_vcpus_per_socket", utils.Int64Value(resp.VMSpec.Spec.Resources.NumVcpusPerSocket))

	d.Set("num_sockets", utils.Int64Value(resp.VMSpec.Spec.Resources.NumSockets))

	d.Set("memory_size_mib", utils.Int64Value(resp.VMSpec.Spec.Resources.MemorySizeMib))

	d.Set("guest_customization_is_overridable", isOv)

	d.Set("vga_console_enabled", utils.BoolValue(resp.VMSpec.Spec.Resources.VgaConsoleEnabled))

	d.SetId(imageID.(string))

	return nil
}

func findOVAImageByUUID(conn *v3.Client, uuid string) (*v3.OVAImageIntentResponse, error) {
	return conn.V3.GetOVAImage(uuid)
}

func findOVAImageByName(ctx context.Context, conn *v3.Client, name string) string {
	entityType := "ova"

	request := v3.GroupsGetEntitiesRequest{
		EntityType:     &entityType,
		FilterCriteria: fmt.Sprintf(`name==%s`, name),
	}

	var response *v3.GroupsGetEntitiesResponse
	response, err := conn.V3.GroupsGetEntities(ctx, &request)

	if err != nil {
		if response != nil {
			log.Printf("Partial response: %+v", response)
		}
	} else {
		groupResults := response.GroupResults
		if len(groupResults) > 0 {
			entityList := groupResults[0].EntityResults
			if len(entityList) > 0 {
				log.Printf("Found OVA image by name: %v", &entityList[0].Data[0].Values[0])
				return entityList[0].EntityID
			}
		}
	}

	return ""
}
