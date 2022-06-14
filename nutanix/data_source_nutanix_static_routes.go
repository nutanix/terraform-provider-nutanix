package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func dataSourceNutanixStaticRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixStaticRouteRead,
		Schema: map[string]*schema.Schema{
			"vpc_reference_uuid": {
				Type:     schema.TypeString,
				Required: true,
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
												"nexthop": {
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
														},
													},
												},
												"destination": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
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
												"nexthop": {
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
														},
													},
												},
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
	conn := meta.(*Client).API

	vpc, ok := d.GetOk("vpc_reference_uuid")
	if !ok {
		return diag.Errorf("please provide one of vpc_reference_uuid attributes")
	}

	// Get request to static Routes

	resp, err := conn.V3.GetStaticRoute(ctx, vpc.(string))
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
		stats := make(map[string]interface{}, 0)

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
		stats := make(map[string]interface{}, 0)

		stats["static_routes_list"] = flattenSpecStaticRouteList(pr.StaticRoutesList)

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

		nhs["external_subnet_reference"] = flattenReferenceValues(nh.ExternalSubnetReference)

		nhList = append(nhList, nhs)

		return nhList
	}
	return nil
}
