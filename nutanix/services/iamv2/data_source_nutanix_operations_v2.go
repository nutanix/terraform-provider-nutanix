package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authz"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixOperationsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixOperationsV4Read,
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
			"operations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"related_operation_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"associated_endpoint_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"endpoint_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"http_method": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixOperationsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

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

	resp, err := conn.OperationsAPIInstance.ListOperations(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching operations : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("operations", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of operations.",
		}}
	}

	operations := resp.Data.GetValue().([]import1.Operation)

	if err := d.Set("operations", flattenPermissionEntities(operations)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenPermissionEntities(pr []import1.Operation) []interface{} {
	if len(pr) > 0 {
		operations := make([]interface{}, len(pr))

		for k, v := range pr {
			permission := make(map[string]interface{})

			if v.ExtId != nil {
				permission["ext_id"] = v.ExtId
			}
			if v.DisplayName != nil {
				permission["display_name"] = v.DisplayName
			}
			if v.Description != nil {
				permission["description"] = v.Description
			}
			if v.EntityType != nil {
				permission["entity_type"] = v.EntityType
			}
			if v.OperationType != nil {
				permission["operation_type"] = flattenOperationType(v.OperationType)
			}
			if v.ClientName != nil {
				permission["client_name"] = v.ClientName
			}
			if v.RelatedOperationList != nil {
				permission["related_operation_list"] = utils.StringSlice(v.RelatedOperationList)
			}

			if v.AssociatedEndpointList != nil {
				permission["associated_endpoint_list"] = flattenAssociatedEndpointList(v.AssociatedEndpointList)
			}
			if v.CreatedTime != nil {
				t := v.CreatedTime
				permission["created_time"] = t.String()
			}
			if v.LastUpdatedTime != nil {
				t := v.LastUpdatedTime
				permission["last_updated_time"] = t.String()
			}
			operations[k] = permission
		}
		return operations
	}
	return nil
}
