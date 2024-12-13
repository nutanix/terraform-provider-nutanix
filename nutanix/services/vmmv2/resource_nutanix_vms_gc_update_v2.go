package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVMGCUpdateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVMGCUpdateV2Create,
		ReadContext:   ResourceNutanixVMGCUpdateV2Read,
		UpdateContext: ResourceNutanixVMGCUpdateV2Update,
		DeleteContext: ResourceNutanixVMGCUpdateV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysprep": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"install_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"PREPARED", "FRESH"}, false),
									},
									"sysprep_script": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unattend_xml": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"custom_key_values": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key_value_pairs": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"value": {
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
											},
										},
									},
								},
							},
						},
						"cloud_init": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datasource_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "CONFIG_DRIVE_V2",
									},
									"metadata": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"cloud_init_script": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_data": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"custom_keys": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key_value_pairs": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"value": {
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

func ResourceNutanixVMGCUpdateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	vmExtID := d.Get("ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	body := &config.GuestCustomizationParams{}

	if configData, ok := d.GetOk("config"); ok {
		body.Config = expandOneOfGuestCustomizationParamsConfig(configData)
	}

	resp, err := conn.VMAPIInstance.CustomizeGuestVm(utils.StringPtr(vmExtID.(string)), body, args)
	if err != nil {
		return diag.Errorf("error while creating Vm's Customize Guest : %v", err)
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
		return diag.Errorf("error waiting for virtual Machine (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching vm UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId

	d.SetId(*uuid)

	return nil
}

func ResourceNutanixVMGCUpdateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVMGCUpdateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVMGCUpdateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
