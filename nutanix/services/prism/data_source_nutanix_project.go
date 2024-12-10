package prism

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixProjectRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"project_name"},
			},
			"project_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"project_id"},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"resource_domain": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"units": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"limit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"account_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"environment_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"default_subnet_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"external_user_group_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"subnet_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"external_network_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"default_environment_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"acp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_reference_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"user_group_reference_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"role_reference": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"context_filter_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope_filter_expression_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"left_hand_side": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"right_hand_side": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"collection": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"categories": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"value": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"uuid_list": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
									"entity_filter_expression_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"left_hand_side_entity_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"right_hand_side": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"collection": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"categories": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"value": {
																			Type:     schema.TypeSet,
																			Computed: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"uuid_list": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceNutanixProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	id, iok := d.GetOk("project_id")
	name, nOk := d.GetOk("project_name")

	if !iok && !nOk {
		return diag.Errorf("please provide `project_id` or `project_name`")
	}

	var err error
	var er error
	var project *v3.ProjectInternalIntentResponse
	var projectID string
	if iok {
		project, err = conn.V3.GetProjectInternal(ctx, id.(string))
	}
	if nOk {
		projectID, er = findProjectByName(conn, name.(string))
		if er != nil {
			return diag.FromErr(er)
		}
		project, err = conn.V3.GetProjectInternal(ctx, projectID)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	m, c := setRSEntityMetadata(project.Metadata)

	if err := d.Set("name", project.Status.ProjectStatus.Name); err != nil {
		return diag.Errorf("error setting `name` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("description", project.Status.ProjectStatus.Description); err != nil {
		return diag.Errorf("error setting `description` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("state", project.Status.State); err != nil {
		return diag.Errorf("error setting `state` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("is_default", project.Status.ProjectStatus.Resources.IsDefault); err != nil {
		return diag.Errorf("error setting `is_default` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("resource_domain", flattenResourceDomain(project.Spec.ProjectDetail.Resources.ResourceDomain)); err != nil {
		return diag.Errorf("error setting `resource_domain` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("account_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.AccountReferenceList)); err != nil {
		return diag.Errorf("error setting `account_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("environment_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.EnvironmentReferenceList)); err != nil {
		return diag.Errorf("error setting `environment_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("default_subnet_reference", flattenReference(project.Spec.ProjectDetail.Resources.DefaultSubnetReference)); err != nil {
		return diag.Errorf("error setting `default_subnet_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("user_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.UserReferenceList)); err != nil {
		return diag.Errorf("error setting `user_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("external_user_group_reference_list",
		flattenReferenceList(project.Spec.ProjectDetail.Resources.ExternalUserGroupReferenceList)); err != nil {
		return diag.Errorf("error setting `external_user_group_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("subnet_reference_list", flattenReferenceList(project.Status.ProjectStatus.Resources.SubnetReferenceList)); err != nil {
		return diag.Errorf("error setting `subnet_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("external_network_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.ExternalNetworkList)); err != nil {
		return diag.Errorf("error setting `external_network_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting `metadata` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("project_reference", flattenReferenceValues(project.Metadata.ProjectReference)); err != nil {
		return diag.Errorf("error setting `project_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("owner_reference", flattenReferenceValues(project.Metadata.OwnerReference)); err != nil {
		return diag.Errorf("error setting `owner_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("categories", c); err != nil {
		return diag.Errorf("error setting `categories` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("api_version", project.APIVersion); err != nil {
		return diag.Errorf("error setting `api_version` for Project(%s): %s", d.Id(), err)
	}

	if err := d.Set("tunnel_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.TunnelReferenceList)); err != nil {
		return diag.Errorf("error setting `tunnel_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("vpc_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.VPCReferenceList)); err != nil {
		return diag.Errorf("error setting `vpc_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("cluster_reference_list", flattenReferenceList(project.Spec.ProjectDetail.Resources.ClusterReferenceList)); err != nil {
		return diag.Errorf("error setting `cluster_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("default_environment_reference", flattenReferenceValuesList(project.Spec.ProjectDetail.Resources.DefaultEnvironmentReference)); err != nil {
		return diag.Errorf("error setting `default_environment_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("acp", flattenProjectAcp(project.Status.AccessControlPolicyListStatus)); err != nil {
		return diag.Errorf("error setting `acp` for Project(%s): %s", d.Id(), err)
	}

	d.SetId(utils.StringValue(project.Metadata.UUID))

	return nil
}

func findProjectByName(conn *v3.Client, name string) (string, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllProject(filter)
	if err != nil {
		return "nil", err
	}

	entities := resp.Entities

	found := make([]*v3.Project, 0)
	for _, v := range entities {
		if v.Spec.Name == name {
			found = append(found, v)
		}
	}

	if len(found) > 1 {
		return "nil", fmt.Errorf("your query returned more than one result. Please use project_id argument instead")
	}

	if len(found) == 0 {
		return "nil", fmt.Errorf("project with the given name, not found")
	}

	return utils.StringValue(found[0].Metadata.UUID), nil
}
