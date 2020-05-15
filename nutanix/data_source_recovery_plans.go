package nutanix

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func dataSourceNutanixRecoveryPlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixRecoveryPlansRead,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"last_update_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"spec_version": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"spec_hash": {
										Type:     schema.TypeString,
										Optional: true,
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
						"categories": categoriesSchema(),
						"owner_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"project_reference": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stage_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"stage_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"delay_time_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"stage_work": {
										Type:     schema.TypeList,
										MinItems: 1,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"recover_entities": {
													Type:     schema.TypeList,
													MinItems: 1,
													MaxItems: 1,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"entity_info_list": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"any_entity_reference_kind": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"any_entity_reference_uuid": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"any_entity_reference_name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"categories": categoriesSchema(),
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
						},
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"floating_ip_assignment_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"availability_zone_url": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"vm_ip_assignment_list": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"test_floating_ip_config": {
																Type:     schema.TypeMap,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"should_allocate_dynamically": {
																			Type:     schema.TypeBool,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"recovery_floating_ip_config": {
																Type:     schema.TypeMap,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"should_allocate_dynamically": {
																			Type:     schema.TypeBool,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"vm_reference": {
																Type:     schema.TypeMap,
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
															"vm_nic_information": {
																Type:     schema.TypeMap,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"uuid": {
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
									"network_mapping_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"availability_zone_network_mapping_list": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"availability_zone_url": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"recovery_network": {
																Type:     schema.TypeList,
																MinItems: 1,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"virtual_network_reference": {
																			Type:     schema.TypeMap,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"kind": {
																						Type:     schema.TypeString,
																						Optional: true,
																						Computed: true,
																					},
																					"uuid": {
																						Type:     schema.TypeString,
																						Optional: true,
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
																		"subnet_list": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"gateway_ip": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"external_connectivity_state": {
																						Type:     schema.TypeString,
																						Optional: true,
																						Computed: true,
																					},
																					"prefix_length": {
																						Type:     schema.TypeInt,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"name": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
															"test_network": {
																Type:     schema.TypeList,
																MinItems: 1,
																MaxItems: 1,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"virtual_network_reference": {
																			Type:     schema.TypeMap,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"kind": {
																						Type:     schema.TypeString,
																						Optional: true,
																						Computed: true,
																					},
																					"uuid": {
																						Type:     schema.TypeString,
																						Optional: true,
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
																		"subnet_list": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"gateway_ip": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"external_connectivity_state": {
																						Type:     schema.TypeString,
																						Optional: true,
																						Computed: true,
																					},
																					"prefix_length": {
																						Type:     schema.TypeInt,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"name": {
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
								},
							},
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixRecoveryPlansRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	resp, err := conn.V3.ListAllRecoveryPlans()
	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}
	if err := d.Set("entities", flattenRecoveryPlanEntities(resp.Entities)); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRecoveryPlanEntities(protectionRules []*v3.RecoveryPlanResponse) []map[string]interface{} {
	entities := make([]map[string]interface{}, len(protectionRules))

	for i, recoveryPlan := range protectionRules {
		metadata, categories := setRSEntityMetadata(recoveryPlan.Metadata)

		entities[i] = map[string]interface{}{
			"name":              recoveryPlan.Status.Name,
			"metadata":          metadata,
			"categories":        categories,
			"project_reference": flattenReferenceValues(recoveryPlan.Metadata.ProjectReference),
			"owner_reference":   flattenReferenceValues(recoveryPlan.Metadata.OwnerReference),
			"stage_list":        flattenStageList(recoveryPlan.Status.Resources.StageList),
			"parameters":        flattenParameters(recoveryPlan.Spec.Resources.Parameters),
			"state":             recoveryPlan.Status.State,
			"api_version":       recoveryPlan.APIVersion,
		}
	}
	return entities
}
