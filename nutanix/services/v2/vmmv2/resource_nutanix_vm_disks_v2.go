package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsDisksV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsDisksV4Create,
		ReadContext:   ResourceNutanixVmsDisksV4Read,
		UpdateContext: ResourceNutanixVmsDisksV4Update,
		DeleteContext: ResourceNutanixVmsDisksV4Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_address": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bus_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR", "PCI", "IDE", "SATA"}, false),
						},
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"backing_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_disk": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size_bytes": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"storage_container": {
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
									"storage_config": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_flash_mode_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"data_source": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"reference": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"image_reference": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_ext_id": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"vm_disk_reference": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"disk_ext_id": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"disk_address": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bus_type": {
																						Type:     schema.TypeString,
																						Optional: true,
																						Computed: true,
																						ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR", "PCI",
																							"IDE", "SATA"}, false),
																					},
																					"index": {
																						Type:     schema.TypeInt,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"vm_reference": {
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
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"is_migration_in_progress": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"adfs_volume_group_reference": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_group_ext_id": {
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
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixVmsDisksV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	diskInput := config.Disk{}
	var busType string
	var idx int

	if extID, ok := d.GetOk("ext_id"); ok {
		diskInput.ExtId = utils.StringPtr(extID.(string))
	}
	if diskAddress, ok := d.GetOk("disk_address"); ok {
		diskInput.DiskAddress = expandDiskAddress(diskAddress)
		busType = flattenDiskBusType(diskInput.DiskAddress.BusType)
		idx = utils.IntValue(diskInput.DiskAddress.Index)
	} else {
		diskInput.DiskAddress = nil
	}

	if backingInfo, ok := d.GetOk("backing_info"); ok {
		diskInput.BackingInfo = expandOneOfDiskBackingInfo(backingInfo)
	} else {
		diskInput.BackingInfo = nil
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.CreateDisk(utils.StringPtr(vmExtID.(string)), &diskInput, args)
	if err != nil {
		return diag.Errorf("error while creating Disk : %v", err)
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
		return diag.Errorf("error waiting for disk (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// reading again VM to fetch disk UUID

	SubResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	out := SubResp.Data.GetValue().(config.Vm)

	diskOut := out.Disks
	diskUUID := ""
	for _, v := range diskOut {
		respBusType := flattenDiskBusType(v.DiskAddress.BusType)

		if respBusType == (busType) && utils.IntValue(v.DiskAddress.Index) == (idx) {
			diskUUID = *v.ExtId
		}
	}
	d.SetId(diskUUID)
	return ResourceNutanixVmsDisksV4Read(ctx, d, meta)
}

func ResourceNutanixVmsDisksV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	resp, err := conn.VMAPIInstance.GetDiskById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching disk : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Disk)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenApiLink(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_address", flattenDiskAddress(getResp.DiskAddress)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backing_info", flattenOneOfDiskBackingInfo(getResp.BackingInfo)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVmsDisksV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	resp, err := conn.VMAPIInstance.GetDiskById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching disk : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Disk)

	updateSpec := getResp

	if d.HasChange("disk_address") {
		updateSpec.DiskAddress = expandDiskAddress(d.Get("disk_address"))
	}
	if d.HasChange("backing_info") {
		updateSpec.BackingInfo = expandOneOfDiskBackingInfo(d.Get("backing_info"))
	}

	updateResp, err := conn.VMAPIInstance.UpdateDiskById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), &updateSpec)
	if err != nil {
		return diag.Errorf("error while updating disk : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for disk (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixVmsDisksV4Read(ctx, d, meta)
}

func ResourceNutanixVmsDisksV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.DeleteDiskById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while fetching disk : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for disk (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
