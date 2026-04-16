package networking

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
)

func DataSourceNutanixStaticRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixStaticRouteRead,
		Schema: map[string]*schema.Schema{
			"vpc_reference_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_name"},
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_reference_uuid"},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"static_routes_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"nexthop": NexthopSpecSchemaForDataSource(),
												"destination": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"default_route_nexthop": NexthopSpecSchemaForDataSource(),
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"static_routes_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"nexthop": NexthopStatusSchemaForDataSource(),
												"destination": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_active": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"default_route": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"nexthop": NexthopStatusSchemaForDataSource(),
												"destination": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_active": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"local_routes_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"nexthop": NexthopStatusSchemaForDataSource(),
												"destination": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_active": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"dynamic_routes_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"nexthop": NexthopStatusSchemaForDataSource(),
												"destination": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_active": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"execution_context": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_uuid": {
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

func dataSourceNutanixStaticRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	vpc, ok := d.GetOk("vpc_reference_uuid")
	vName, nok := d.GetOk("vpc_name")
	if !ok && !nok {
		return diag.Errorf("please provide one of vpc_reference_uuid or vpc_name attributes")
	}

	var vpcName string
	if ok {
		vpcName = vpc.(string)
	} else {
		var reqErr error
		var resp *v3.VPCIntentResponse
		resp, reqErr = findVPCByName(ctx, conn, vName.(string))
		if reqErr != nil {
			return diag.FromErr(reqErr)
		}
		vpcName = *resp.Metadata.UUID
	}

	// Get request to static Routes

	resp, err := conn.V3.GetStaticRoute(ctx, vpcName)
	if err != nil {
		return diag.FromErr(err)
	}

	m, _ := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("spec", flattenStaticRouteSpec(resp.Spec)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", flattenStaticRouteStatus(resp.Status)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Metadata.UUID)
	return nil
}

func flattenStaticRouteSpec(stat *v3.StaticRouteSpec) []interface{} {
	statList := make([]interface{}, 0)
	if stat != nil {
		stats := make(map[string]interface{})

		stats["name"] = stat.Name
		stats["resources"] = flattenStaticRouteSpecResources(stat.Resources)

		statList = append(statList, stats)
		return statList
	}
	return nil
}

func flattenStaticRouteStatus(srs *v3.StaticRouteDefStatus) []interface{} {
	srsList := make([]interface{}, 0)
	if srs != nil {
		sr := make(map[string]interface{})

		sr["state"] = srs.State
		sr["resources"] = flattenStaticRouteStatusResources(srs.Resources)
		sr["execution_context"] = flattenExecutionContext(srs.ExecutionContext)

		srsList = append(srsList, sr)
		return srsList
	}
	return nil
}

func flattenStaticRouteSpecResources(pr *v3.StaticRouteResources) []interface{} {
	prList := make([]interface{}, 0)
	if pr != nil {
		stats := make(map[string]interface{})

		stats["static_routes_list"] = flattenSpecStaticRouteList(pr.StaticRoutesList)
		stats["default_route_nexthop"] = flattenNextHop(pr.DefaultRouteNexthop)

		prList = append(prList, stats)
		return prList
	}
	return nil
}

func flattenStaticRouteStatusResources(pr *v3.StaticRouteResources) []interface{} {
	prList := make([]interface{}, 0)
	if pr != nil {
		prs := make(map[string]interface{})

		prs["static_routes_list"] = flattenStatusStaticRouteList(pr.StaticRoutesList)
		prs["default_route"] = flattenStatusDefaultRouteList(pr.DefaultRoute)
		prs["local_routes_list"] = flattenStatusStaticRouteList(pr.LocalRoutesList)
		prs["dynamic_routes_list"] = flattenStatusStaticRouteList(pr.DynamicRoutesList)
		prList = append(prList, prs)
		return prList
	}
	return nil
}

func flattenSpecStaticRouteList(sr []*v3.StaticRoutesList) []map[string]interface{} {
	srList := make([]map[string]interface{}, len(sr))

	if len(sr) > 0 {
		for k, v := range sr {
			srs := make(map[string]interface{})

			srs["destination"] = v.Destination
			srs["nexthop"] = flattenNextHop(v.NextHop)

			srList[k] = srs
		}
		return srList
	}
	return nil
}

func flattenStatusDefaultRouteList(sr *v3.StaticRoutesList) []map[string]interface{} {
	srList := make([]map[string]interface{}, 0)
	if sr != nil {
		srs := make(map[string]interface{})

		srs["destination"] = sr.Destination
		srs["nexthop"] = flattenNextHop(sr.NextHop)
		srs["is_active"] = sr.IsActive
		srs["priority"] = sr.Priority

		srList = append(srList, srs)

		return srList
	}
	return nil
}

func flattenStatusStaticRouteList(sr []*v3.StaticRoutesList) []map[string]interface{} {
	srList := make([]map[string]interface{}, len(sr))

	if len(sr) > 0 {
		for k, v := range sr {
			srs := make(map[string]interface{})

			srs["destination"] = v.Destination
			srs["nexthop"] = flattenNextHop(v.NextHop)
			srs["is_active"] = v.IsActive
			srs["priority"] = v.Priority

			srList[k] = srs
		}
		return srList
	}
	return nil
}

func flattenNextHop(nh *v3.NextHop) []interface{} {
	nhList := make([]interface{}, 0)

	if nh != nil {
		nhs := make(map[string]interface{})

		if nh.ExternalSubnetReference != nil {
			nhs["external_subnet_reference"] = flattenReferenceValues(nh.ExternalSubnetReference)
		}
		if nh.LocalSubnetReference != nil {
			nhs["local_subnet_reference"] = flattenReferenceValues(nh.LocalSubnetReference)
		}

		if nh.DirectConnectVirtualInterfaceReference != nil {
			nhs["direct_connect_virtual_interface_reference"] = flattenReferenceValues(nh.DirectConnectVirtualInterfaceReference)
		}

		if nh.VpnConnectionReference != nil {
			nhs["vpn_connection_reference"] = flattenReferenceValues(nh.DirectConnectVirtualInterfaceReference)
		}

		if nh.NexthopIPAddress != nil {
			nhs["nexthop_ip_address"] = nh.NexthopIPAddress
		}
		nhList = append(nhList, nhs)

		return nhList
	}
	return nil
}

func NexthopStatusSchemaForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"external_subnet_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"direct_connect_virtual_interface_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"local_subnet_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"vpn_connection_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"nexthop_ip_address": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func NexthopSpecSchemaForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"external_subnet_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"direct_connect_virtual_interface_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"local_subnet_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"vpn_connection_reference": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}
