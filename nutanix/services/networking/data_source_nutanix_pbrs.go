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

func DataSourceNutanixPbrs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixPbrsRead,
		Schema: map[string]*schema.Schema{
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
						"sort_attribute": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"total_matches": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
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
												"is_bidirectional": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"vpc_reference": {
													Type:     schema.TypeMap,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"destination": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"subnet_ip": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"subnet_ip": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"routing_policy_counters": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"packet_count": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"byte_count": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"protocol_parameters": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tcp": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_port_range_list": portRangeSchemaForDataSource(),
																		"source_port_range_list":      portRangeSchemaForDataSource(),
																	},
																},
															},
															"udp": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_port_range_list": portRangeSchemaForDataSource(),
																		"source_port_range_list":      portRangeSchemaForDataSource(),
																	},
																},
															},
															"icmp": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"icmp_type": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"icmp_code": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"protocol_number": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"action": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"service_ip_list": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"protocol_type": {
													Type:     schema.TypeString,
													Computed: true,
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
												"is_bidirectional": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"vpc_reference": {
													Type:     schema.TypeMap,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"destination": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"subnet_ip": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"source": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"address_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"subnet_ip": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"priority": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"protocol_parameters": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tcp": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_port_range_list": portRangeSchemaForDataSource(),
																		"source_port_range_list":      portRangeSchemaForDataSource(),
																	},
																},
															},
															"udp": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"destination_port_range_list": portRangeSchemaForDataSource(),
																		"source_port_range_list":      portRangeSchemaForDataSource(),
																	},
																},
															},
															"icmp": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"icmp_type": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"icmp_code": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																	},
																},
															},
															"protocol_number": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"action": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"service_ip_list": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"protocol_type": {
													Type:     schema.TypeString,
													Computed: true,
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

func dataSourceNutanixPbrsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")
	if filtersOk {
		request = buildDataSourceListMetadata(metadata.(*schema.Set))
	}

	resp, err := conn.V3.ListAllPBR(ctx, utils.StringValue(request.Filter))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("api_version", resp.APIVersion)

	if err := d.Set("metadata", flattenPbrMetadata(resp.Metadata)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entities", flattenPbrEntities(resp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenPbrEntities(pr []*v3.PbrIntentResponse) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ents := make(map[string]interface{})
			ents["spec"] = flattenPbrSpec(v.Spec)
			ents["status"] = flattenPbrStatus(v.Status)

			m, _ := setRSEntityMetadata(v.Metadata)
			ents["metadata"] = m

			res[k] = ents
		}
		return res
	}
	return nil
}

func flattenPbrMetadata(met *v3.ListMetadataOutput) []interface{} {
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
