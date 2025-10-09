package networkingv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRoutesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRoutesV2Read,
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
			"route_table_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixRouteV2(),
			},
		},
	}
}

func DatasourceNutanixRoutesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DatasourceNutanixRoutesV2Read")

	conn := meta.(*conns.Client).NetworkingAPI

	routeTableExtID := d.Get("route_table_ext_id").(string)

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
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}

	resp, err := conn.Routes.ListRoutesByRouteTableId(&routeTableExtID, page, limit, filter, orderBy)
	if err != nil {
		return diag.Errorf("error while fetching routes : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("routes", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of routes.",
		}}
	}

	getResp := resp.Data.GetValue().([]config.Route)

	if err := d.Set("routes", flattenRoutes(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRoutes(routes []config.Route) []interface{} {
	log.Printf("[DEBUG] flattenRoutes")

	if routes == nil {
		log.Printf("[DEBUG] flattenRoutes: routes is empty")
		return make([]interface{}, 0)
	}

	log.Printf("[DEBUG] flattenRoutes: routes is not empty")

	result := make([]interface{}, len(routes))

	for i, route := range routes {
		log.Printf("[DEBUG] flattenRoutes: route[%d]: %v", i, route)
		result[i] = map[string]interface{}{
			"tenant_id":                         route.TenantId,
			"ext_id":                            route.ExtId,
			"links":                             flattenLinks(route.Links),
			"metadata":                          flattenMetadata(route.Metadata),
			"name":                              route.Name,
			"description":                       route.Description,
			"destination":                       flattenDestination(route.Destination),
			"next_hop":                          flattenNextHop(route.Nexthop),
			"route_table_reference":             route.RouteTableReference,
			"vpc_reference":                     route.VpcReference,
			"external_routing_domain_reference": route.ExternalRoutingDomainReference,
			"route_type":                        flattenRouteType(route.RouteType),
			"is_active":                         route.IsActive,
			"priority":                          route.Priority,
		}
	}

	log.Printf("[DEBUG] flattenRoutes: result: %v", result)
	return result
}
