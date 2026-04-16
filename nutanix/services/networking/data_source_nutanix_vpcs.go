package networking

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixVPCs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixVPCsRead,
		Schema: map[string]*schema.Schema{
			//Computed attributes
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"offset": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"total_matches": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sort_attribute": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resources": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"external_subnet_list": {
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
															"external_ip_list": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"active_gateway_node": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"host_reference": {
																			Type:     schema.TypeMap,
																			Required: true,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"ip_address": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
															"active_gateway_count": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"externally_routable_prefix_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Required: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"common_domain_name_server_ip_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
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
												"external_subnet_list": {
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
												"externally_routable_prefix_list": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Required: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Required: true,
															},
														},
													},
												},
												"common_domain_name_server_ip_list": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Optional: true,
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
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixVPCsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")
	if filtersOk {
		request = buildDataSourceListMetadata(metadata.(*schema.Set))
	}

	resp, err := conn.V3.ListAllVPC(ctx, utils.StringValue(request.Filter))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entities", flattenVPCEntities(resp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("metadata", flattenVPCMetadata(resp.Metadata)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVPCEntities(ent []*v3.VPCIntentResponse) []map[string]interface{} {
	if len(ent) > 0 {
		entList := make([]map[string]interface{}, len(ent))

		for k, v := range ent {
			ents := make(map[string]interface{})
			ents["status"] = flattenStatusVPC(v.Status)
			ents["spec"] = flattenSpecVPC(v.Spec)

			m, _ := setRSEntityMetadata(v.Metadata)
			ents["metadata"] = m

			entList[k] = ents
		}
		return entList
	}
	return nil
}

func flattenStatusVPC(stat *v3.VPCDefStatus) []interface{} {
	statList := make([]interface{}, 0)
	if stat != nil {
		stats := make(map[string]interface{})

		stats["state"] = stat.State
		stats["name"] = stat.Name
		stats["resources"] = flattenResourcesVPC(stat.Resources)
		stats["execution_context"] = flattenExecutionContext(stat.ExecutionContext)

		statList = append(statList, stats)
	}
	return statList
}

func flattenSpecVPC(vpc *v3.VPC) []interface{} {
	vpcList := make([]interface{}, 0)

	if vpc != nil {
		vpcs := make(map[string]interface{})

		vpcs["name"] = vpc.Name
		vpcs["resources"] = flattenVPCResources(vpc.Resources)

		vpcList = append(vpcList, vpcs)
	}
	return vpcList
}

func flattenResourcesVPC(res *v3.VpcResources) []interface{} {
	resList := make([]interface{}, 0)

	if res != nil {
		ress := make(map[string]interface{})

		ress["common_domain_name_server_ip_list"] = flattenCommonDNSIPList(res.CommonDomainNameServerIPList)
		ress["external_subnet_list"] = flattenExtSubnetListStatus(res.ExternalSubnetList)
		ress["externally_routable_prefix_list"] = flattenExtRoutableList(res.ExternallyRoutablePrefixList)

		resList = append(resList, ress)
	}
	return resList
}

func flattenExecutionContext(exe *v3.ExecutionContext) []interface{} {
	exec := make([]interface{}, 0)

	if exe != nil {
		execdata := make(map[string]interface{})

		execdata["task_uuid"] = exe.TaskUUID

		exec = append(exec, execdata)
	}

	return exec
}

func flattenVPCResources(res *v3.VpcResources) []interface{} {
	resList := make([]interface{}, 0)

	if res != nil {
		ress := make(map[string]interface{})

		ress["common_domain_name_server_ip_list"] = flattenCommonDNSIPList(res.CommonDomainNameServerIPList)
		ress["external_subnet_list"] = flattenExtSubnetList(res.ExternalSubnetList)
		ress["externally_routable_prefix_list"] = flattenExtRoutableList(res.ExternallyRoutablePrefixList)

		resList = append(resList, ress)
	}
	return resList
}

func flattenVPCMetadata(met *v3.ListMetadataOutput) []interface{} {
	metList := make([]interface{}, 0)

	if met != nil {
		mets := make(map[string]interface{})

		mets["total_matches"] = met.TotalMatches
		mets["kind"] = met.Kind
		mets["length"] = met.Length
		mets["offset"] = met.Offset
		mets["filter"] = met.Filter

		metList = append(metList, mets)
	}
	return metList
}
