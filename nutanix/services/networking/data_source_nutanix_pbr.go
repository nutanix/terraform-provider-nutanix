package networking

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixPbr() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixPbrRead,
		Schema: map[string]*schema.Schema{
			"pbr_uuid": {
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
		},
	}
}

func dataSourceNutanixPbrRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	pbrUUID, ok := d.GetOk("pbr_uuid")
	if !ok {
		return diag.Errorf("please provide pbr reference uuid")
	}

	// make call to the GetPBR API
	resp, err := conn.V3.GetPBR(ctx, pbrUUID.(string))
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

	if err := d.Set("status", flattenPbrStatus(resp.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("spec", flattenPbrSpec(resp.Spec)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Metadata.UUID)
	return nil
}

func flattenPbrStatus(pr *v3.PbrDefStatus) []interface{} {
	if pr != nil {
		res := make([]interface{}, 0)

		prs := make(map[string]interface{})
		prs["state"] = pr.State
		prs["execution_context"] = flattenExecutionContext(pr.ExecutionContext)
		prs["resources"] = flattenPbrResources(pr.Resources)

		res = append(res, prs)

		return res
	}
	return nil
}

func flattenPbrSpec(pr *v3.PbrSpec) []interface{} {
	if pr != nil {
		res := make([]interface{}, 0)

		prs := make(map[string]interface{})

		prs["name"] = pr.Name
		prs["resources"] = flattenPbrResources(pr.Resources)

		res = append(res, prs)
		return res
	}
	return nil
}

func flattenPbrResources(pr *v3.PbrResources) []interface{} {
	if pr != nil {
		res := make([]interface{}, 0)

		prs := make(map[string]interface{})

		prs["is_bidirectional"] = pr.IsBidirectional
		prs["priority"] = pr.Priority
		prs["protocol_type"] = pr.ProtocolType
		prs["vpc_reference"] = flattenReferenceValues(pr.VpcReference)
		prs["source"] = flattenSourceDest(pr.Source)
		prs["destination"] = flattenSourceDest(pr.Destination)

		if pr.RoutingPolicyCounters != nil {
			rpc := make([]interface{}, 0)

			rp := make(map[string]interface{})
			rp["packet_count"] = pr.RoutingPolicyCounters.PacketCount
			rp["byte_count"] = pr.RoutingPolicyCounters.ByteCount

			rpc = append(rpc, rp)

			prs["routing_policy_counters"] = rpc
		}

		if pr.Action != nil {
			act := make([]interface{}, 0)

			ac := make(map[string]interface{})

			ac["action"] = pr.Action.Action
			if pr.Action.ServiceIPList != nil {
				ac["service_ip_list"] = utils.StringSlice(pr.Action.ServiceIPList)
			}

			act = append(act, ac)
			prs["action"] = act
		}

		prs["protocol_parameters"] = flattenProtocolParams(pr.ProtocolParameters)

		res = append(res, prs)
		return res
	}
	return nil
}

func portRangeSchemaForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"end_port": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"start_port": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}
