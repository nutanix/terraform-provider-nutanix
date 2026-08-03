package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/request/routetables"
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
					},
				},
			},
		},
	}
}

func DatasourceNutanixRouteTablesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	listRouteTablesRequest := import2.ListRouteTablesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listRouteTablesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listRouteTablesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listRouteTablesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listRouteTablesRequest.Orderby_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.RoutesTable.ListRouteTables(ctx, &listRouteTablesRequest)
	if err != nil {
		return diag.Errorf("error while fetching route tables : %v", err)
	}

	if resp.Data == nil {
		log.Printf("[DEBUG] DatasourceNutanixRouteTablesV2Read: No data found.")
		if err := d.Set("route_tables", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of route tables.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.RouteTable)
	aJSON, _ := json.Marshal(getResp)
	log.Printf("[DEBUG] DatasourceNutanixRouteTablesV2Read: %v", string(aJSON))

	if err := d.Set("route_tables", flattenRouteTableEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRouteTableEntities(pr []import1.RouteTable) []interface{} {
	if pr == nil {
		return make([]interface{}, 0)
	}

	routeTables := make([]interface{}, len(pr))
	for i, v := range pr {
		routeTables[i] = map[string]interface{}{
			"ext_id":                            v.ExtId,
			"tenant_id":                         v.TenantId,
			"links":                             flattenLinks(v.Links),
			"metadata":                          flattenMetadata(v.Metadata),
			"vpc_reference":                     v.VpcReference,
			"external_routing_domain_reference": v.ExternalRoutingDomainReference,
		}
	}

	aJSON, _ := json.Marshal(routeTables)
	log.Printf("[DEBUG] flattenRouteTableEntities: %v", string(aJSON))
	return routeTables
}
