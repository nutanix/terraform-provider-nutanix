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

func DatasourceNutanixVmsCdRomsV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmsCdRomsV4Read,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cdroms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func DatasourceNutanixVmsCdRomsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	// initialize query params
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	vmExtID := d.Get("vm_ext_id")
	resp, err := conn.VMAPIInstance.ListCdRomsByVmId(utils.StringPtr(vmExtID.(string)), page, limit)
	if err != nil {
		return diag.Errorf("error while fetching cd-roms : %v", err)
	}
	getResp := resp.Data

	if getResp != nil {
		cdRomResp := getResp.GetValue().([]config.CdRom)
		if err := d.Set("cdroms", flattenCdRomsEntities(cdRomResp)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(resource.UniqueId())
	return nil
}

func flattenCdRomsEntities(pr []config.CdRom) []interface{} {
	if len(pr) > 0 {
		cdRomList := make([]interface{}, len(pr))

		for k, v := range pr {
			cdRom := make(map[string]interface{})

			if v.ExtId != nil {
				cdRom["ext_id"] = v.ExtId
			}
			if v.TenantId != nil {
				cdRom["tenant_id"] = v.ExtId
			}
			if v.Links != nil {
				cdRom["links"] = flattenApiLink(v.Links)
			}
			if v.DiskAddress != nil {
				cdRom["disk_address"] = flattenCdRomAddress(v.DiskAddress)
			}
			if v.BackingInfo != nil {
				cdRom["backing_info"] = flattenVmDisk(v.BackingInfo)
			}
			if v.IsoType != nil {
				cdRom["iso_type"] = flattenIsoType(v.IsoType)
			}

			cdRomList[k] = cdRom
		}
		return cdRomList
	}
	return nil
}
