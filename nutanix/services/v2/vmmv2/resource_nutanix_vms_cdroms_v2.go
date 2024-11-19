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

func ResourceNutanixVmsCdRomsV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsCdRomsV4Create,
		ReadContext:   ResourceNutanixVmsCdRomsV4Read,
		UpdateContext: ResourceNutanixVmsCdRomsV4Update,
		DeleteContext: ResourceNutanixVmsCdRomsV4Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
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
							ValidateFunc: validation.StringInSlice([]string{"IDE", "SATA"}, false),
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
													Computed: true,
													Optional: true,
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
																			ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR",
																				"PCI", "IDE", "SATA"}, false),
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
			"iso_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION"}, false),
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

func ResourceNutanixVmsCdRomsV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	cdRomInput := config.CdRom{}
	var busType string
	var idx int

	if extID, ok := d.GetOk("ext_id"); ok {
		cdRomInput.ExtId = utils.StringPtr(extID.(string))
	}
	if diskAddress, ok := d.GetOk("disk_address"); ok && len(diskAddress.([]interface{})) > 0 {
		cdRomInput.DiskAddress = expandCdRomAddress(diskAddress)
		busType = flattenCdRomBusType(cdRomInput.DiskAddress.BusType)
		idx = utils.IntValue(cdRomInput.DiskAddress.Index)
	} else {
		cdRomInput.DiskAddress = nil
	}

	if backingInfo, ok := d.GetOk("backing_info"); ok && len(backingInfo.([]interface{})) > 0 {
		cdRomInput.BackingInfo = expandVmDisk(backingInfo)
	} else {
		cdRomInput.BackingInfo = nil
	}

	if isoType, ok := d.GetOk("iso_type"); ok {
		subMap := map[string]interface{}{
			"OTHER":               2,
			"GUEST_TOOLS":         3,
			"GUEST_CUSTOMIZATION": 4,
		}
		pVal := subMap[isoType.(string)]
		p := config.IsoType(pVal.(int))
		cdRomInput.IsoType = &p
	}

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.CreateCdRom(utils.StringPtr(vmExtID.(string)), &cdRomInput, args)
	if err != nil {
		return diag.Errorf("error while creating cd-rom : %v", err)
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
		return diag.Errorf("error waiting for cd-rom (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// reading again VM to fetch disk UUID

	SubResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	out := SubResp.Data.GetValue().(config.Vm)

	cdOut := out.CdRoms
	cdRomUUID := ""
	for _, v := range cdOut {
		respBusType := flattenCdRomBusType(v.DiskAddress.BusType)

		if respBusType == (busType) && utils.IntValue(v.DiskAddress.Index) == (idx) {
			cdRomUUID = *v.ExtId
		}
	}
	d.SetId(cdRomUUID)
	return ResourceNutanixVmsCdRomsV4Read(ctx, d, meta)
}

func ResourceNutanixVmsCdRomsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	resp, err := conn.VMAPIInstance.GetCdRomById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching cd-rom : %v", err)
	}

	cdRespData := resp.Data.GetValue().(config.CdRom)
	if err := d.Set("tenant_id", cdRespData.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenApiLink(cdRespData.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_address", flattenCdRomAddress(cdRespData.DiskAddress)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backing_info", flattenVmDisk(cdRespData.BackingInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iso_type", flattenIsoType(cdRespData.IsoType)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVmsCdRomsV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVmsCdRomsV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.DeleteCdRomById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while fetching cd-rom : %v", err)
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
		return diag.Errorf("error waiting for cd-rom (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
