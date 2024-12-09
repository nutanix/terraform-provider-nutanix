package networking

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
)

func DataSourceNutanixVPC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixVPCRead,
		Schema: map[string]*schema.Schema{
			"vpc_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_name"},
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_uuid"},
			},
			//  Computed attributes

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
		},
	}
}

func dataSourceNutanixVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	uuid, iok := d.GetOk("vpc_uuid")
	name, nok := d.GetOk("vpc_name")

	if !iok && !nok {
		return diag.Errorf("please provide one of vpc_uuid or vpc_name attributes")
	}

	var reqErr error
	var resp *v3.VPCIntentResponse

	if iok {
		resp, reqErr = findVPCByUUID(ctx, conn, uuid.(string))
	} else {
		resp, reqErr = findVPCByName(ctx, conn, name.(string))
	}

	if reqErr != nil {
		if strings.Contains(fmt.Sprint(reqErr), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return diag.Errorf("error reading user with error %s", reqErr)
	}

	m, _ := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", flattenStatusVPC(resp.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("spec", flattenSpecVPC(resp.Spec)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func findVPCByUUID(ctx context.Context, conn *v3.Client, uuid string) (*v3.VPCIntentResponse, error) {
	return conn.V3.GetVPC(ctx, uuid)
}

func findVPCByName(ctx context.Context, conn *v3.Client, name string) (*v3.VPCIntentResponse, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllVPC(ctx, filter)
	if err != nil {
		return nil, err
	}

	entities := resp.Entities

	found := make([]*v3.VPCIntentResponse, 0)
	for _, v := range entities {
		if *v.Status.Name == name {
			found = append(found, v)
		}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please use uuid argument instead")
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("vpc with the given name, not found")
	}

	return findVPCByUUID(ctx, conn, *found[0].Metadata.UUID)
}
