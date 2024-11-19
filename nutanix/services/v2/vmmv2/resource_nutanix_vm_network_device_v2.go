package vmmv2

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVMNetworkDeviceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVMNetworkDeviceV2Create,
		ReadContext:   ResourceNutanixVMNetworkDeviceV2Read,
		UpdateContext: ResourceNutanixVMNetworkDeviceV2Update,
		DeleteContext: ResourceNutanixVMNetworkDeviceV2Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backing_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"VIRTIO", "E1000"}, false),
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(
									"^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"),
								"MAC address should be in format xx:xx:xx:xx:xx:xx or xx-xx-xx-xx-xx-xx"),
						},
						"is_connected": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"num_queues": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Default:      1,
						},
					},
				},
			},
			"network_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{"SPAN_DESTINATION_NIC",
								"NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC"}, false),
						},
						"network_function_chain": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"network_function_nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{"TAP", "EGRESS",
								"INGRESS"}, false),
						},
						"subnet": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"vlan_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"TRUNK", "ACCESS"}, false),
						},
						"trunked_vlans": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntAtLeast(0),
							},
						},
						"should_allow_unknown_macs": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"ipv4_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"should_assign_ip": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Default:  32,
												},
											},
										},
									},
									"secondary_ip_address_list": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Default:  32,
												},
											},
										},
									},
								},
							},
						},
						"ipv4_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"learned_ip_addresses": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Default:  32,
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

func ResourceNutanixVMNetworkDeviceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	body := &config.Nic{}

	vmExtID := d.Get("vm_ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	if extID, ok := d.GetOk("ext_id"); ok {
		body.ExtId = utils.StringPtr(extID.(string))
	}
	if backing_info, ok := d.GetOk("backing_info"); ok {
		body.BackingInfo = expandEmulatedNic(backing_info)
	}
	if network_info, ok := d.GetOk("network_info"); ok {
		body.NetworkInfo = expandNicNetworkInfo(network_info)
	}

	resp, err := conn.VMAPIInstance.CreateNic(utils.StringPtr(vmExtID.(string)), body, args)
	if err != nil {
		return diag.Errorf("error while creating vm's nic : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vm's nic (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// reading again VM to fetch nic UUID

	SubResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	out := SubResp.Data.GetValue().(config.Vm)

	nicOut := out.Nics
	nicUUID := ""

	if len(nicOut) == 1 {
		nicUUID = *nicOut[0].ExtId
	}

	d.SetId(nicUUID)

	return ResourceNutanixVMNetworkDeviceV2Read(ctx, d, meta)
}

func ResourceNutanixVMNetworkDeviceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	resp, err := conn.VMAPIInstance.GetNicById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching network device : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Nic)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backing_info", flattenEmulatedNic(getResp.BackingInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_info", flattenNicNetworkInfo(getResp.NetworkInfo)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixVMNetworkDeviceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	resp, err := conn.VMAPIInstance.GetNicById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching vm : %v", errorMessage["message"])
	}

	respNicss := resp.Data.GetValue().(config.Nic)
	updateSpec := respNicss

	if d.HasChange("backing_info") {
		updateSpec.BackingInfo = expandEmulatedNic(d.Get("backing_info"))
	}
	if d.HasChange("network_info") {
		updateSpec.NetworkInfo = expandNicNetworkInfo(d.Get("network_info"))
	}

	updateResp, err := conn.VMAPIInstance.UpdateNicById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		var errordata map[string]interface{}

		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}

		return diag.Errorf("error while updating vm's nic : %v", e)
	}

	TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vm's nic (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func ResourceNutanixVMNetworkDeviceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	readResp, err := conn.VMAPIInstance.GetNicById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}

		return diag.Errorf("error while reading vm's nic : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.DeleteNicById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting vm's nic : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vm's nic (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
