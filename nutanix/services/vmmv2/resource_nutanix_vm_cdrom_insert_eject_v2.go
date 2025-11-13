package vmmv2

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmsCdRomsInsertEjectV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmsCdRomsInsertEjectV2Create,
		ReadContext:   ResourceNutanixVmsCdRomsInsertEjectV2Read,
		UpdateContext: ResourceNutanixVmsCdRomsInsertEjectV2Update,
		DeleteContext: ResourceNutanixVmsCdRomsInsertEjectV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				const expectedPartsCount = 2
				parts := strings.Split(d.Id(), "/")
				if len(parts) != expectedPartsCount {
					return nil, fmt.Errorf("invalid import uuid (%q), expected vm_ext_id/cdrom_ext_id", d.Id())
				}
				d.Set("vm_ext_id", parts[0])
				d.Set("ext_id", parts[1])
				d.Set("action", "insert")
				d.SetId(resource.UniqueId())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cdrom_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bus_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"backing_info": {
				Type:     schema.TypeList,
				Optional: true,
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
									},
								},
							},
						},
						"storage_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_flash_mode_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"data_source": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reference": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image_reference": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"image_ext_id": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"vm_disk_reference": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_ext_id": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"disk_address": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bus_type": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"index": {
																			Type:     schema.TypeInt,
																			Optional: true,
																		},
																	},
																},
															},
															"vm_reference": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ext_id": {
																			Type:     schema.TypeString,
																			Optional: true,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "insert",
				ValidateFunc: validation.StringInSlice([]string{"insert", "eject"}, false),
			},
		},
	}
}

func ResourceNutanixVmsCdRomsInsertEjectV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixVmsCdRomsInsertEjectV2Create : Inserting ISO into the CD-ROM %s of the VM %s", d.Get("ext_id").(string), d.Get("vm_ext_id").(string))
	if action, ok := d.GetOk("action"); ok && action.(string) == "insert" {
		conn := meta.(*conns.Client).VmmAPI

		vmExtID := d.Get("vm_ext_id")
		extID := d.Get("ext_id")
		body := config.CdRomInsertParams{}

		if backInfo, ok := d.GetOk("backing_info"); ok {
			body.BackingInfo = expandVMDisk(backInfo)
		}

		readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
		if err != nil {
			return diag.Errorf("error while reading vm : %v", err)
		}
		// Extract E-Tag Header
		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readResp, conn)

		resp, err := conn.VMAPIInstance.InsertCdRomById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)), &body, args)
		if err != nil {
			return diag.Errorf("error while inserting cd-rom : %v", err)
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
			return diag.Errorf("error waiting for cd-rom (%s) to insert: %s", utils.StringValue(taskUUID), errWaitTask)
		}

		d.SetId(resource.UniqueId())
		return ResourceNutanixVmsCdRomsInsertEjectV2Read(ctx, d, meta)
	}
	return diag.Errorf("Action %s is not supported for CD-ROM Insert", d.Get("action").(string))
}

func ResourceNutanixVmsCdRomsInsertEjectV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixVmsCdRomsInsertEjectV2Read : Reading CD-ROM %s of the VM %s", d.Get("ext_id").(string), d.Get("vm_ext_id").(string))
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")

	readResp, err := conn.VMAPIInstance.GetCdRomById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while reading cd-rom : %v", err)
	}

	getResp := readResp.Data.GetValue().(config.CdRom)

	if err := d.Set("iso_type", getResp.IsoType.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_address", flattenCdRomAddress(getResp.DiskAddress)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backing_info", flattenVMDisk(getResp.BackingInfo)); err != nil {
		return diag.FromErr(err)
	}
	// set the cdrom ext id to the state file, same as ext_id
	// This attribute is used to eject, ejectCdromISO is a common function for both ngt ISO and Other ISOs
	if err := d.Set("cdrom_ext_id", extID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixVmsCdRomsInsertEjectV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if action, ok := d.GetOk("action"); ok && action.(string) == "eject" {
		log.Printf("[DEBUG] ResourceNutanixVmsCdRomsInsertEjectV2Update : Action %s", action.(string))
		diags := ejectCdromISO(ctx, d, meta)
		if diags.HasError() {
			// Ejection failed, set the action to INSERT to avoid Terraform from saving "EJECT" in state
			d.Set("action", "insert")
			return diags
		}
		return ResourceNutanixVmsCdRomsInsertEjectV2Read(ctx, d, meta)
	}
	return nil
}

func ResourceNutanixVmsCdRomsInsertEjectV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixVmsCdRomsInsertEjectV2Delete : Ejecting ISO from the CD-ROM %s of the VM %s", d.Get("ext_id").(string), d.Get("vm_ext_id").(string))
	if action, ok := d.GetOk("action"); ok && action.(string) == "eject" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "ISO is not inserted on the CD-ROM of the VM or ejected earlier using an action, Ignoring the request to eject the ISO",
		}}
	}
	return ejectCdromISO(ctx, d, meta)
}
