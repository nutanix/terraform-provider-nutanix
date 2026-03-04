package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/images/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/imageplacementpolicies"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixImagePlacementsV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixImagePlacementsV4Read,
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
			"placement_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"placement_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_entity_filter": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category_ext_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"cluster_entity_filter": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category_ext_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"enforcement_state": {
							Type:     schema.TypeString,
							Computed: true,
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
					},
				},
			},
		},
	}
}

func DatasourceNutanixImagePlacementsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	listPlacementPoliciesRequest := import2.ListPlacementPoliciesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listPlacementPoliciesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listPlacementPoliciesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listPlacementPoliciesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listPlacementPoliciesRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listPlacementPoliciesRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.ImagesPlacementAPIInstance.ListPlacementPolicies(ctx, &listPlacementPoliciesRequest)
	if err != nil {
		return diag.Errorf("error while fetching image placement policies : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("placement_policies", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of placement policies.",
		}}
	}

	policies := resp.Data.GetValue().([]import7.PlacementPolicy)

	if err := d.Set("placement_policies", flattenPlacementPolicyEntities(policies)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenPlacementPolicyEntities(pr []import7.PlacementPolicy) []interface{} {
	if len(pr) > 0 {
		policies := make([]interface{}, len(pr))

		for k, v := range pr {
			policy := make(map[string]interface{})

			if v.ExtId != nil {
				policy["ext_id"] = v.ExtId
			}
			if v.Name != nil {
				policy["name"] = v.Name
			}
			if v.Description != nil {
				policy["description"] = v.Description
			}
			if v.PlacementType != nil {
				policy["placement_type"] = flattenPlacementType(v.PlacementType)
			}
			if v.ImageEntityFilter != nil {
				policy["image_entity_filter"] = flattenEntityFilter(v.ImageEntityFilter)
			}
			if v.ClusterEntityFilter != nil {
				policy["cluster_entity_filter"] = flattenEntityFilter(v.ClusterEntityFilter)
			}
			if v.EnforcementState != nil {
				policy["enforcement_state"] = flattenEnforcementState(v.EnforcementState)
			}
			if v.CreateTime != nil {
				t := v.CreateTime
				policy["create_time"] = t.String()
			}
			if v.LastUpdateTime != nil {
				t := v.LastUpdateTime
				policy["last_update_time"] = t.String()
			}
			if v.OwnerExtId != nil {
				policy["owner_ext_id"] = v.OwnerExtId
			}
			policies[k] = policy
		}
		return policies
	}
	return nil
}
