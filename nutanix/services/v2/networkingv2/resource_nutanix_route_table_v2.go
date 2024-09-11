package networkingv2

//
//import (
//	"context"
//	"fmt"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
//
//	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/networking/v4/config"
//	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v16/models/prism/v4/config"
//
//	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
//	"github.com/terraform-providers/terraform-provider-nutanix/utils"
//)
//
//func ResourceNutanixRouteTablesV2() *schema.Resource {
//	return &schema.Resource{
//		CreateContext: ResourceNutanixRouteTablesV2Create,
//		ReadContext:   ResourceNutanixRouteTablesV2Read,
//		UpdateContext: ResourceNutanixRouteTablesV2Update,
//		DeleteContext: ResourceNutanixRouteTablesV2Delete,
//		Schema: map[string]*schema.Schema{
//			"vpc_ext_id": {
//				Type:     schema.TypeString,
//				Required: true,
//			},
//			"metadata": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem: &schema.Resource{
//					Schema: DatasourceMetadataSchemaV4(),
//				},
//			},
//			"vpc_reference": {
//				Type:     schema.TypeString,
//				Optional: true,
//				Computed: true,
//			},
//			"external_routing_domain_reference": {
//				Type:     schema.TypeString,
//				Optional: true,
//				Computed: true,
//			},
//			"static_routes": {
//				Type:     schema.TypeList,
//				Optional: true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"destination": {
//							Type:     schema.TypeList,
//							Required: true,
//							Elem: &schema.Resource{
//								Schema: map[string]*schema.Schema{
//									"ipv4": {
//										Type:     schema.TypeList,
//										Optional: true,
//										Elem: &schema.Resource{
//											Schema: map[string]*schema.Schema{
//												"ip": SchemaForValuePrefixLength(),
//												"prefix_length": {
//													Type:     schema.TypeInt,
//													Optional: true,
//												},
//											},
//										},
//									},
//									"ipv6": {
//										Type:     schema.TypeList,
//										Optional: true,
//										Elem: &schema.Resource{
//											Schema: map[string]*schema.Schema{
//												"ip": SchemaForValuePrefixLength(),
//												"prefix_length": {
//													Type:     schema.TypeInt,
//													Optional: true,
//												},
//											},
//										},
//									},
//								},
//							},
//						},
//						"next_hop_type": {
//							Type:     schema.TypeString,
//							Required: true,
//							ValidateFunc: validation.StringInSlice([]string{"INTERNAL_SUBNET", "DIRECT_CONNECT_VIF",
//								"VPN_CONNECTION", "IP_ADDRESS", "EXTERNAL_SUBNET"}, false),
//						},
//						"next_hop_reference": {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						"next_hop_ip_address": {
//							Type:     schema.TypeList,
//							Optional: true,
//							Elem: &schema.Resource{
//								Schema: map[string]*schema.Schema{
//									"ipv4": SchemaForValuePrefixLength(),
//									"ipv6": SchemaForValuePrefixLength(),
//								},
//							},
//						},
//						"next_hop_name": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"source": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"is_active": {
//							Type:     schema.TypeBool,
//							Computed: true,
//						},
//						"priority": {
//							Type:     schema.TypeInt,
//							Computed: true,
//						},
//					},
//				},
//			},
//			"dynamic_routes": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem: &schema.Resource{
//					Schema: DatasourceRoutesSchemaV4(),
//				},
//			},
//			"local_routes": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem: &schema.Resource{
//					Schema: DatasourceRoutesSchemaV4(),
//				},
//			},
//			"tenant_id": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			"links": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"href": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"rel": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//					},
//				},
//			},
//			"ext_id": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//		},
//	}
//}
//
//func ResourceNutanixRouteTablesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	conn := meta.(*conns.Client).NetworkingAPI
//	extID := ""
//
//	vpcExtID := d.Get("vpc_ext_id").(string)
//
//	filter := fmt.Sprintf("vpcReference eq '%s'", vpcExtID)
//	vpcResp, err := conn.RoutesTable.ListRouteTables(nil, nil, &filter, nil, nil)
//	if err != nil {
//		return diag.Errorf("error while fetching vpc : %v", err)
//	}
//
//	vpcGetResp := vpcResp.Data.GetValue().([]import1.RouteTable)
//
//	extID = *vpcGetResp[0].ExtId
//
//	req := &import1.RouteTable{}
//
//	if vpcRef, ok := d.GetOk("vpc_reference"); ok {
//		req.VpcReference = utils.StringPtr(vpcRef.(string))
//	}
//	if extRouting, ok := d.GetOk("external_routing_domain_reference"); ok {
//		req.ExternalRoutingDomainReference = utils.StringPtr(extRouting.(string))
//	}
//	if staticRoute, ok := d.GetOk("static_routes"); ok {
//		req.StaticRoutes = expandRoute(staticRoute.([]interface{}))
//	}
//
//	// Get Etag
//	routeResp, err := conn.RoutesTable.GetRouteTable(utils.StringPtr(extID))
//	if err != nil {
//		return diag.Errorf("error while fetching vpcs : %v", err)
//	}
//	// Extract E-Tag Header
//	etagValue := conn.APIClientInstance.GetEtag(routeResp)
//
//	args := make(map[string]interface{})
//	args["If-Match"] = etagValue
//
//	resp, err := conn.RoutesTable.UpdateRouteTable(req, &extID, args)
//	if err != nil {
//		return diag.Errorf("error while updating route tables : %v", err)
//	}
//
//	taskRef := resp.Data.GetValue().(import4.TaskReference)
//	taskUUID := taskRef.ExtId
//
//	taskconn := meta.(*conns.Client).PrismAPI
//	// Wait for the route table to be available
//	stateConf := &resource.StateChangeConf{
//		Pending: []string{"QUEUED", "QUEUED"},
//		Target:  []string{"SUCCEEDED"},
//		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
//		Timeout: d.Timeout(schema.TimeoutCreate),
//	}
//
//	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
//		return diag.Errorf("error waiting for route table (%s) to perform: %s", utils.StringValue(taskUUID), errWaitTask)
//	}
//
//	d.SetId(extID)
//	return ResourceNutanixRouteTablesV2Read(ctx, d, meta)
//}
//
//func ResourceNutanixRouteTablesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	conn := meta.(*conns.Client).NetworkingAPI
//
//	resp, err := conn.RoutesTable.GetRouteTable(utils.StringPtr(d.Id()))
//	if err != nil {
//		return diag.Errorf("error while fetching route table : %v", err)
//	}
//
//	getResp := resp.Data.GetValue().(import1.RouteTable)
//
//	if err := d.Set("ext_id", getResp.ExtId); err != nil {
//		return diag.FromErr(err)
//	}
//
//	if err := d.Set("links", flattenLinksExternalNetworkingApi(getResp.Links)); err != nil {
//		return diag.FromErr(err)
//	}
//	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
//		return diag.FromErr(err)
//	}
//	if err := d.Set("metadata", flattenMetadataExternalNetworkingApi(getResp.Metadata)); err != nil {
//		return diag.FromErr(err)
//	}
//
//	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
//		return diag.FromErr(err)
//	}
//	if err := d.Set("external_routing_domain_reference", getResp.ExternalRoutingDomainReference); err != nil {
//		return diag.FromErr(err)
//	}
//	if err := d.Set("static_routes", flattenRoute(getResp.StaticRoutes)); err != nil {
//		return diag.FromErr(err)
//	}
//	if err := d.Set("dynamic_routes", flattenRoute(getResp.DynamicRoutes)); err != nil {
//		return diag.FromErr(err)
//	}
//	if err := d.Set("local_routes", flattenRoute(getResp.LocalRoutes)); err != nil {
//		return diag.FromErr(err)
//	}
//	return nil
//}
//
//func ResourceNutanixRouteTablesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	conn := meta.(*conns.Client).NetworkingAPI
//	updateSpec := &import1.RouteTable{}
//	resp, err := conn.RoutesTable.GetRouteTable(utils.StringPtr(d.Id()))
//	if err != nil {
//		return diag.Errorf("error while fetching route table: %v", err)
//	}
//
//	getResp := resp.Data.GetValue().(import1.RouteTable)
//
//	// updateSpec = &getResp
//
//	if d.HasChange("vpc_reference") {
//		updateSpec.VpcReference = utils.StringPtr(d.Get("vpc_reference").(string))
//	}
//	if d.HasChange("external_routing_domain_reference") {
//		updateSpec.ExternalRoutingDomainReference = utils.StringPtr(d.Get("external_routing_domain_reference").(string))
//	}
//	if d.HasChange("static_routes") {
//		updateSpec.StaticRoutes = expandRoute(d.Get("static_routes").([]interface{}))
//	}
//
//	if len(updateSpec.StaticRoutes) == 0 {
//		updateSpec.StaticRoutes = []import1.Route{}
//	}
//
//	// Extract E-Tag Header
//	etagValue := conn.APIClientInstance.GetEtag(getResp)
//
//	args := make(map[string]interface{})
//	args["If-Match"] = etagValue
//
//	updateResp, err := conn.RoutesTable.UpdateRouteTable(updateSpec, utils.StringPtr(d.Id()), args)
//	if err != nil {
//		return diag.Errorf("error while updating route tables : %v", err)
//	}
//
//	taskRef := updateResp.Data.GetValue().(import4.TaskReference)
//	taskUUID := taskRef.ExtId
//
//	taskconn := meta.(*conns.Client).PrismAPI
//	// Wait for the route table to be available
//	stateConf := &resource.StateChangeConf{
//		Pending: []string{"QUEUED", "QUEUED"},
//		Target:  []string{"SUCCEEDED"},
//		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
//		Timeout: d.Timeout(schema.TimeoutCreate),
//	}
//
//	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
//		return diag.Errorf("error waiting for route table (%s) to perform: %s", utils.StringValue(taskUUID), errWaitTask)
//	}
//	return ResourceNutanixRouteTablesV2Read(ctx, d, meta)
//}
//
//func ResourceNutanixRouteTablesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	return nil
//}
//
//func expandRoute(pr []interface{}) []import1.Route {
//	if len(pr) > 0 {
//		routes := make([]import1.Route, len(pr))
//
//		for k, v := range pr {
//			route := import1.Route{}
//			val := v.(map[string]interface{})
//
//			if dest, ok := val["destination"]; ok && len(dest.([]interface{})) > 0 {
//				route.Destination = expandIPSubnetObject(dest)
//			}
//			if nextHop, ok := val["next_hop_type"]; ok && len(nextHop.(string)) > 0 {
//				const two, three, four, five, six = 2, 3, 4, 5, 6
//				hopTypeMaps := map[string]interface{}{
//					"IP_ADDRESS":         two,
//					"DIRECT_CONNECT_VIF": three,
//					"INTERNAL_SUBNET":    four,
//					"EXTERNAL_SUBNET":    five,
//					"VPN_CONNECTION":     six,
//				}
//				opVal := hopTypeMaps[nextHop.(string)]
//				op := import1.NexthopType(opVal.(int))
//				route.NexthopType = &op
//			}
//			if nextHopRef, ok := val["next_hop_reference"]; ok && len(nextHopRef.(string)) > 0 {
//				route.NexthopReference = utils.StringPtr(nextHopRef.(string))
//			}
//			if nextHopIPAdd, ok := val["next_hop_ip_address"]; ok && len(nextHopIPAdd.([]interface{})) > 0 {
//				route.NexthopIpAddress = expandIPAddressObject(nextHopIPAdd)
//			}
//
//			routes[k] = route
//		}
//		return routes
//	}
//	return nil
//}
