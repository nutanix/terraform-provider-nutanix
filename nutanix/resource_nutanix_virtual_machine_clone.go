package nutanix

import (
	"context"
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixVirtualMachineClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixVirtualMachineCloneCreate,
		UpdateContext: resourceNutanixVirtualMachineCloneUpdate,
		ReadContext:   resourceNutanixVirtualMachineCloneRead,
		DeleteContext: resourceNutanixVirtualMachineCloneDelete,
		Schema: map[string]*schema.Schema{
			"vm_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"num_sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"num_vcpus_per_socket": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"num_threads_per_core": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"memory_size_mib": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"nic_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_function_nic_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_endpoint_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"network_function_chain_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"num_queues": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"subnet_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_connected": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "true",
						},
					},
				},
			},

			"boot_device_order_list": {
				Type:     schema.TypeList,
				Optional: true,
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"UEFI", "LEGACY", "SECURE_BOOT"}, false),
			},
			"guest_customization_cloud_init_user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"guest_customization_cloud_init_meta_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"guest_customization_cloud_init_custom_key_values": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"guest_customization_is_overridable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"guest_customization_sysprep": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"guest_customization_sysprep_custom_key_values": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"task_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed  Resource Argument
			"cloud_init_cdrom_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"categories": categoriesSchema(),
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_reference": {
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
			"nic_list_status": {
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

			// RESOURCES ARGUMENTS

			"enable_cpu_passthrough": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_vcpu_hard_pinned": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"use_hot_add": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"num_vnuma_nodes": {
				Type:     schema.TypeInt,
				Computed: true,
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
			"ngt_credentials": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"ngt_enabled_capability_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
							Optional: true,
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
							Optional: true,
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
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"parent_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"machine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hardware_clock_timezone": {
				Type:     schema.TypeString,
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
							Optional: true,
							Computed: true,
						},
						"disk_size_bytes": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"disk_size_mib": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"storage_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flash_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"storage_container_reference": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"kind": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "storage_container",
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Optional: true,
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
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "DISK",
									},
									"disk_address": {
										Type:     schema.TypeMap,
										Optional: true,
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
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"volume_group_reference": {
							Type:     schema.TypeMap,
							Optional: true,
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
							Optional: true,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixVirtualMachineCloneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	var id string
	vmUUID, nok := d.GetOk("vm_uuid")
	if nok {
		id = *utils.StringPtr(vmUUID.(string))
	}
	spec := &v3.VMCloneInput{}

	spec.Metadata = getMetadataCloneAttributes(d)
	spec.OverrideSpec = expandOverrideSpec(d)

	// Make request to the API
	resp, err := conn.V3.CloneVM(id, spec)
	if err != nil {
		return diag.FromErr(err)
	}
	taskUUID := *resp.TaskUUID

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	taskInfo, errWaitTask := stateConf.WaitForStateContext(ctx)
	if errWaitTask != nil {
		return diag.Errorf("error waiting for task (%s) to create: %s", taskUUID, errWaitTask)
	}

	// Get the cloned VM UUID
	var cloneVMUUID string
	taskDetails, ok := taskInfo.(*v3.TasksResponse)
	if ok {
		cloneVMUUID = *taskDetails.EntityReferenceList[0].UUID
	}

	// State Changed to Power ON
	if er := changePowerState(ctx, conn, cloneVMUUID, "ON"); er != nil {
		return diag.Errorf("internal error: cannot turn ON the VM with UUID(%s): %s", cloneVMUUID, err)
	}

	// Wait for IP available
	waitIPConf := &resource.StateChangeConf{
		Pending:    []string{WAITING},
		Target:     []string{"AVAILABLE"},
		Refresh:    waitForIPRefreshFunc(conn, cloneVMUUID),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	vmIntentResponse, ero := waitIPConf.WaitForStateContext(ctx)
	if ero != nil {
		log.Printf("[WARN] could not get the IP for VM(%s): %s", cloneVMUUID, err)
	} else {
		vm := vmIntentResponse.(*v3.VMIntentResponse)

		if len(vm.Status.Resources.NicList) > 0 && len(vm.Status.Resources.NicList[0].IPEndpointList) != 0 {
			d.SetConnInfo(map[string]string{
				"type": "ssh",
				"host": *vm.Status.Resources.NicList[0].IPEndpointList[0].IP,
			})
		}
	}

	d.SetId(cloneVMUUID)
	//return nil
	return resourceNutanixVirtualMachineCloneRead(ctx, d, meta)
}

func resourceNutanixVirtualMachineCloneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixVirtualMachineCloneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNutanixVirtualMachineRead(ctx, d, meta)
}

func resourceNutanixVirtualMachineCloneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func getMetadataCloneAttributes(d *schema.ResourceData) (out *v3.Metadata) {
	resourceData, ok := d.GetOk("metadata")
	if !ok {
		return nil
	}

	meta := resourceData.([]interface{})[0].(map[string]interface{})

	if name, ok := meta["uuid"]; ok {
		out.UUID = utils.StringPtr(name.(string))
	}

	if name, ok := meta["entity_version"]; ok {
		out.EntityVersion = utils.StringPtr(name.(string))
	}
	return out
}

func expandOverrideSpec(d *schema.ResourceData) *v3.OverrideSpec {
	res := &v3.OverrideSpec{}

	if name, ok := d.GetOk("name"); ok {
		res.Name = utils.StringPtr(name.(string))
	}

	if numSockets, sok := d.GetOk("num_sockets"); sok {
		res.NumSockets = utils.IntPtr(numSockets.(int))
	}

	if vcpuSock, vok := d.GetOk("num_vcpus_per_socket"); vok {
		res.NumVcpusPerSocket = utils.IntPtr(vcpuSock.(int))
	}

	if numThreads, ok := d.GetOk("num_threads_per_core"); ok {
		res.NumThreadsPerCore = utils.IntPtr(numThreads.(int))
	}

	if memorySize, mok := d.GetOk("memory_size_mib"); mok {
		res.MemorySizeMib = utils.IntPtr(memorySize.(int))
	}

	if _, nok := d.GetOk("nic_list"); nok {
		res.NicList = expandNicList(d)
	}

	guestCustom := &v3.GuestCustomization{}
	cloudInit := &v3.GuestCustomizationCloudInit{}

	if v, ok := d.GetOk("guest_customization_cloud_init_user_data"); ok {
		cloudInit.UserData = utils.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("guest_customization_cloud_init_meta_data"); ok {
		cloudInit.MetaData = utils.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("guest_customization_cloud_init_custom_key_values"); ok {
		cloudInit.CustomKeyValues = utils.ConvertMapString(v.(map[string]interface{}))
	}

	if !reflect.DeepEqual(*cloudInit, (v3.GuestCustomizationCloudInit{})) {
		guestCustom.CloudInit = cloudInit
	}

	if v, ok := d.GetOk("guest_customization_is_overridable"); ok {
		guestCustom.IsOverridable = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("guest_customization_sysprep"); ok {
		guestCustom.Sysprep = &v3.GuestCustomizationSysprep{}
		spi := v.(map[string]interface{})
		if v2, ok2 := spi["install_type"]; ok2 {
			guestCustom.Sysprep.InstallType = utils.StringPtr(v2.(string))
		}
		if v2, ok2 := spi["unattend_xml"]; ok2 {
			guestCustom.Sysprep.UnattendXML = utils.StringPtr(v2.(string))
		}
	}

	if v, ok := d.GetOk("guest_customization_sysprep_custom_key_values"); ok {
		if guestCustom.Sysprep == nil {
			guestCustom.Sysprep = &v3.GuestCustomizationSysprep{}
		}
		guestCustom.Sysprep.CustomKeyValues = v.(map[string]string)
	}

	if !reflect.DeepEqual(*guestCustom, (v3.GuestCustomization{})) {
		res.GuestCustomization = guestCustom
	}
	bootConfig := &v3.VMBootConfig{}

	if v, ok := d.GetOk("boot_device_order_list"); ok {
		bootConfig.BootDeviceOrderList = expandStringList(v.([]interface{}))
		res.BootConfig = bootConfig
	}

	bd := &v3.VMBootDevice{}
	da := &v3.DiskAddress{}
	if v, ok := d.GetOk("boot_device_disk_address"); ok {
		dai := v.(map[string]interface{})

		if value3, ok3 := dai["device_index"]; ok3 {
			if i, err := strconv.ParseInt(value3.(string), 10, 64); err == nil {
				da.DeviceIndex = utils.Int64Ptr(i)
			}
		}
		if value3, ok3 := dai["adapter_type"]; ok3 {
			da.AdapterType = utils.StringPtr(value3.(string))
		}
		bd.DiskAddress = da
		bootConfig.BootDevice = bd
		res.BootConfig = bootConfig
	}

	if bdmac, ok := d.GetOk("boot_device_mac_address"); ok {
		bd.MacAddress = utils.StringPtr(bdmac.(string))
		res.BootConfig.BootDevice = bd
	}

	if bootType, ok := d.GetOk("boot_type"); ok {
		bootConfig.BootType = utils.StringPtr(bootType.(string))
		res.BootConfig = bootConfig
	}
	return res
}
