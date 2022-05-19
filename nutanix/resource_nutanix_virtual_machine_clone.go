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
			},
			"num_sockets": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"num_vcpus_per_socket": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"num_threads_per_core": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory_size_mib": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nic_list": {
				Type:     schema.TypeList,
				Optional: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"boot_device_disk_address": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"boot_device_mac_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"boot_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UEFI", "LEGACY", "SECURE_BOOT"}, false),
			},
			"guest_customization_cloud_init_user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"guest_customization_cloud_init_meta_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"guest_customization_cloud_init_custom_key_values": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"guest_customization_is_overridable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"guest_customization_sysprep": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"guest_customization_sysprep_custom_key_values": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"task_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_vm_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixVirtualMachineCloneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	var id string
	vm_uuid, nok := d.GetOk("vm_uuid")
	if nok {
		id = *utils.StringPtr(vm_uuid.(string))
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

	task_info, errWaitTask := stateConf.WaitForStateContext(ctx)
	if errWaitTask != nil {
		return diag.Errorf("error waiting for task (%s) to create: %s", taskUUID, errWaitTask)
	}

	// Get the cloned VM UUID
	var cloneVmUUID string
	task_details, ok := task_info.(*v3.TasksResponse)
	if ok {
		cloneVmUUID = *task_details.EntityReferenceList[0].UUID
	}

	// State Changed to Power ON
	if err := changePowerState(ctx, conn, cloneVmUUID, "ON"); err != nil {
		return diag.Errorf("internal error: cannot turn ON the VM with UUID(%s): %s", cloneVmUUID, err)
	}

	// Wait for IP available
	waitIPConf := &resource.StateChangeConf{
		Pending:    []string{WAITING},
		Target:     []string{"AVAILABLE"},
		Refresh:    waitForIPRefreshFunc(conn, cloneVmUUID),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	vmIntentResponse, err := waitIPConf.WaitForStateContext(ctx)
	if err != nil {
		log.Printf("[WARN] could not get the IP for VM(%s): %s", cloneVmUUID, err)
	} else {
		vm := vmIntentResponse.(*v3.VMIntentResponse)

		if len(vm.Status.Resources.NicList) > 0 && len(vm.Status.Resources.NicList[0].IPEndpointList) != 0 {
			d.SetConnInfo(map[string]string{
				"type": "ssh",
				"host": *vm.Status.Resources.NicList[0].IPEndpointList[0].IP,
			})
		}
	}

	d.SetId(cloneVmUUID)
	return nil
	//return resourceNutanixVirtualMachineCloneRead(ctx, d, meta)
}

func resourceNutanixVirtualMachineCloneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixVirtualMachineCloneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
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
