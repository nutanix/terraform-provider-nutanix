package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/locations"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLocationsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLocationsV2Read,
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
			"locations": {
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
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"coordinates": schemaForGeographicCoordinates(),
						"links":       schemaForLinks(),
					},
				},
			},
		},
	}
}

func DatasourceNutanixLocationsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NCDomainManagementAPI

	listReq := import1.ListLocationsRequest{}
	if v, ok := d.GetOk("page"); ok {
		listReq.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listReq.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listReq.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listReq.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listReq.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.Locations.ListLocations(ctx, &listReq)
	if err != nil {
		return diag.Errorf("error while listing Locations: %s", err)
	}

	if resp == nil || resp.Data == nil {
		if err := d.Set("locations", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of locations.",
		}}
	}

	var flatLocations []map[string]interface{}
	switch v := resp.Data.GetValue().(type) {
	case []config.Location:
		flatLocations = flattenLocations(v)
	case []config.LocationProjection:
		flatLocations = flattenLocationProjections(v)
	default:
		return diag.Errorf("unexpected type %T for List Locations response", resp.Data.GetValue())
	}

	if len(flatLocations) == 0 {
		if err := d.Set("locations", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of locations.",
		}}
	}

	if err := d.Set("locations", flatLocations); err != nil {
		return diag.Errorf("error setting locations: %s", err)
	}
	d.SetId(utils.GenUUID())
	return nil
}

func flattenLocations(locations []config.Location) []map[string]interface{} {
	if len(locations) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(locations))
	for _, loc := range locations {
		out = append(out, map[string]interface{}{
			"ext_id":      utils.StringValue(loc.ExtId),
			"name":        utils.StringValue(loc.Name),
			"country":     utils.StringValue(loc.Country),
			"state":       utils.StringValue(loc.State),
			"tenant_id":   utils.StringValue(loc.TenantId),
			"coordinates": flattenGeographicCoordinates(loc.Coordinates),
			"links":       flattenLinks(loc.Links),
		})
	}
	return out
}

func flattenLocationProjections(locations []config.LocationProjection) []map[string]interface{} {
	if len(locations) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(locations))
	for _, loc := range locations {
		out = append(out, map[string]interface{}{
			"ext_id":      utils.StringValue(loc.ExtId),
			"name":        utils.StringValue(loc.Name),
			"country":     utils.StringValue(loc.Country),
			"state":       utils.StringValue(loc.State),
			"tenant_id":   utils.StringValue(loc.TenantId),
			"coordinates": flattenGeographicCoordinates(loc.Coordinates),
			"links":       flattenLinks(loc.Links),
		})
	}
	return out
}

func schemaForGeographicCoordinates() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"latitude": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"longitude": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
			},
		},
	}
}

func flattenGeographicCoordinates(coordinates *config.GeographicCoordinates) []map[string]interface{} {
	if coordinates == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"latitude":  utils.Float64Value(coordinates.Latitude),
			"longitude": utils.Float64Value(coordinates.Longitude),
		},
	}
}
