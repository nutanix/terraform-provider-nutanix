package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/ahv/config"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVmsCdRomV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmsCdRomV4Read,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_container": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"storage_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_flash_mode_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"data_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reference": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image_reference": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"image_ext_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"vm_disk_reference": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_ext_id": {
																Type:     schema.TypeString,
																Computed: true,
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
															"vm_reference": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ext_id": {
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
		},
	}
}

func DatasourceNutanixVmsCdRomV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	vmExtID := d.Get("vm_ext_id")
	extID := d.Get("ext_id")

	resp, err := conn.VMAPIInstance.GetCdRomById(utils.StringPtr(vmExtID.(string)), utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching cd-rom : %v", err)
	}

	cdResp := resp.Data

	if cdResp != nil {
		cdRespData := cdResp.GetValue().(config.CdRom)
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
		d.SetId(*cdRespData.ExtId)
		return nil
	}
	d.SetId(resource.UniqueId())
	return nil
}
