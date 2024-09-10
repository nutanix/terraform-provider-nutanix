package networkingv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRouteTablesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRouteTablesV2Read,
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
			"route_tables": {
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
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceMetadataSchemaV4(),
							},
						},
						"vpc_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_routing_domain_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"static_routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceRoutesSchemaV4(),
							},
						},
						"dynamic_routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceRoutesSchemaV4(),
							},
						},
						"local_routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceRoutesSchemaV4(),
							},
						},
					},
				},
			},
		},
	}
}

func DatasourceMetadataSchemaV4() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"owner_reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"owner_user_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"category_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func DatasourceNutanixRouteTablesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// initialize query params
	var filter, orderBy *string
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
	if order_byf, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order_byf.(string))
	} else {
		orderBy = nil
	}

	resp, err := conn.RoutesTable.ListRouteTables(page, limit, filter, orderBy)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching route tables : %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().([]import1.RouteTable)

	if err := d.Set("route_tables", flattenRouteTableEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRouteTableEntities(pr []import1.RouteTable) []interface{} {
	if len(pr) > 0 {
		routes := make([]interface{}, len(pr))

		for k, v := range pr {
			route := make(map[string]interface{})

			route["ext_id"] = v.ExtId
			route["tenant_id"] = v.TenantId
			route["links"] = flattenLinks(v.Links)
			route["metadata"] = flattenMetadata(v.Metadata)
			route["vpc_reference"] = v.VpcReference
			route["external_routing_domain_reference"] = v.ExternalRoutingDomainReference
			routes[k] = route
		}
		return routes
	}
	return nil
}
