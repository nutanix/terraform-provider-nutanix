package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import5 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixImagesV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixImagesV4Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checksum": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"size_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url_source": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"should_allow_insecure_url": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"basic_auth": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"username": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"password": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"vm_disk_source": {
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
									"object_lite_source": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"category_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cluster_location_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"placement_policy_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"placement_policy_ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"compliance_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enforcement_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy_cluster_ext_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"enforced_cluster_ext_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"conflicting_policy_ext_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
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

func DatasourceNutanixImagesV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	// initialize query params
	var filter, orderBy, selects *string
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
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}
	resp, err := conn.ImagesAPIInstance.ListImages(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching images : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("images", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of images.",
		}}
	}

	getResp := resp.Data.GetValue().([]import5.Image)

	if err := d.Set("images", flattenImagesEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenImagesEntities(pr []import5.Image) []interface{} {
	if len(pr) > 0 {
		imgs := make([]interface{}, len(pr))

		for k, v := range pr {
			img := make(map[string]interface{})

			if v.ExtId != nil {
				img["ext_id"] = v.ExtId
			}
			if v.Name != nil {
				img["name"] = v.Name
			}
			if v.Description != nil {
				img["description"] = v.Description
			}
			if v.Type != nil {
				img["type"] = flattenImageType(v.Type)
			}
			if v.Checksum != nil {
				img["checksum"] = flattenOneOfImageChecksum(v.Checksum)
			}
			if v.SizeBytes != nil {
				img["size_bytes"] = v.SizeBytes
			}
			if v.Source != nil {
				img["source"] = flattenOneOfImageSource(v.Source)
			}
			if v.CategoryExtIds != nil {
				img["category_ext_ids"] = v.CategoryExtIds
			}
			if v.ClusterLocationExtIds != nil {
				img["cluster_location_ext_ids"] = v.ClusterLocationExtIds
			}
			if v.CreateTime != nil {
				t := v.CreateTime
				img["create_time"] = t.String()
			}
			if v.LastUpdateTime != nil {
				t := v.LastUpdateTime
				img["last_update_time"] = t.String()
			}
			if v.OwnerExtId != nil {
				img["owner_ext_id"] = v.OwnerExtId
			}
			if v.PlacementPolicyStatus != nil {
				img["placement_policy_status"] = flattenImagePlacementStatus(v.PlacementPolicyStatus)
			}
			imgs[k] = img
		}
		return imgs
	}
	return nil
}
