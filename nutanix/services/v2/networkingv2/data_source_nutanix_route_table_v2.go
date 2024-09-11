package networkingv2

// import (
// 	"context"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/networking/v4/config"

// 	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
// 	"github.com/terraform-providers/terraform-provider-nutanix/utils"
// )

// func DatasourceNutanixRouteTableV2() *schema.Resource {
// 	return &schema.Resource{
// 		ReadContext: DatasourceNutanixRouteTableV2Read,
// 		Schema: map[string]*schema.Schema{
// 			"ext_id": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"tenant_id": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"links": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"href": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"rel": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 					},
// 				},
// 			},
// 			"metadata": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: DatasourceMetadataSchemaV4(),
// 				},
// 			},
// 			"vpc_reference": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"external_routing_domain_reference": {
// 				Type:     schema.TypeString,
// 				Computed: true,
// 			},
// 			"static_routes": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: DatasourceRoutesSchemaV4(),
// 				},
// 			},
// 			"dynamic_routes": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: DatasourceRoutesSchemaV4(),
// 				},
// 			},
// 			"local_routes": {
// 				Type:     schema.TypeList,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: DatasourceRoutesSchemaV4(),
// 				},
// 			},
// 		},
// 	}
// }

// func DatasourceNutanixRouteTableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	conn := meta.(*conns.Client).NetworkingAPI

// 	extID := d.Get("ext_id")
// 	resp, err := conn.RoutesTable.GetRouteTableById(utils.StringPtr(extID.(string)))
// 	if err != nil {
// 		return diag.Errorf("error while fetching route table : %v", err)
// 	}

// 	getResp := resp.Data.GetValue().(import1.RouteTable)

// 	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("external_routing_domain_reference", getResp.ExternalRoutingDomainReference); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("static_routes", flattenRoute(getResp.StaticRoutes)); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("dynamic_routes", flattenRoute(getResp.DynamicRoutes)); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("local_routes", flattenRoute(getResp.LocalRoutes)); err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId(*getResp.ExtId)
// 	return nil
// }

// func flattenRoute(pr []import1.Route) []interface{} {
// 	if len(pr) > 0 {
// 		routes := make([]interface{}, len(pr))

// 		for k, v := range pr {
// 			route := make(map[string]interface{})

// 			route["is_active"] = v.IsActive
// 			route["priority"] = v.Priority
// 			route["destination"] = flattenIPSubnet(v.Destination)
// 			route["next_hop_type"] = flattenNexthopType(v.NexthopType)
// 			if v.NexthopReference != nil {
// 				route["next_hop_reference"] = v.NexthopReference
// 			}
// 			if v.NexthopIpAddress != nil {
// 				route["next_hop_ip_address"] = flattenNodeIPAddress(v.NexthopIpAddress)
// 			}
// 			if v.NexthopName != nil {
// 				route["next_hop_name"] = v.NexthopName
// 			}
// 			if v.Source != nil {
// 				route["source"] = v.Source
// 			}

// 			routes[k] = route
// 		}
// 		return routes
// 	}
// 	return nil
// }

// func DatasourceRoutesSchemaV4() map[string]*schema.Schema {
// 	return map[string]*schema.Schema{
// 		"is_active": {
// 			Type:     schema.TypeBool,
// 			Computed: true,
// 		},
// 		"priority": {
// 			Type:     schema.TypeInt,
// 			Computed: true,
// 		},
// 		"destination": {
// 			Type:     schema.TypeList,
// 			Computed: true,
// 			Elem: &schema.Resource{
// 				Schema: map[string]*schema.Schema{
// 					"ipv4": {
// 						Type:     schema.TypeList,
// 						Computed: true,
// 						Elem: &schema.Resource{
// 							Schema: map[string]*schema.Schema{
// 								"ip": SchemaForValuePrefixLength(),
// 								"prefix_length": {
// 									Type:     schema.TypeInt,
// 									Computed: true,
// 								},
// 							},
// 						},
// 					},
// 					"ipv6": {
// 						Type:     schema.TypeList,
// 						Computed: true,
// 						Elem: &schema.Resource{
// 							Schema: map[string]*schema.Schema{
// 								"ip": SchemaForValuePrefixLength(),
// 								"prefix_length": {
// 									Type:     schema.TypeInt,
// 									Computed: true,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		"next_hop_type": {
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		},
// 		"next_hop_reference": {
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		},
// 		"next_hop_ip_address": {
// 			Type:     schema.TypeList,
// 			Computed: true,
// 			Elem: &schema.Resource{
// 				Schema: map[string]*schema.Schema{
// 					"ipv4": SchemaForValuePrefixLength(),
// 					"ipv6": SchemaForValuePrefixLength(),
// 				},
// 			},
// 		},
// 		"next_hop_name": {
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		},
// 		"source": {
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		},
// 	}
// }

// func flattenNexthopType(pr *import1.NexthopType) string {
// 	if pr != nil {
// 		const two, three, four, five, six = 2, 3, 4, 5, 6

// 		if *pr == import1.NexthopType(two) {
// 			return "IP_ADDRESS"
// 		}
// 		if *pr == import1.NexthopType(three) {
// 			return "DIRECT_CONNECT_VIF"
// 		}
// 		if *pr == import1.NexthopType(four) {
// 			return "INTERNAL_SUBNET"
// 		}
// 		if *pr == import1.NexthopType(five) {
// 			return "EXTERNAL_SUBNET"
// 		}
// 		if *pr == import1.NexthopType(six) {
// 			return "VPN_CONNECTION"
// 		}
// 	}
// 	return "UNKNOWN"
// }
