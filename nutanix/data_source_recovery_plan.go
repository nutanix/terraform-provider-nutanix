package nutanix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNutanixRecoveryPlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixRecoveryPlanRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"recovery_plan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
	}
}

func dataSourceNutanixRecoveryPlanRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	recoveryPlanID := d.Get("recovery_plan_id").(string)
	resp, err := conn.V3.GetRecoveryPlan(recoveryPlanID)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("name", resp.Spec.Name); err != nil {
		return err
	}
	if err := d.Set("stage_list", flattenStageList(resp.Spec.Resources.StageList)); err != nil {
		return err
	}
	if err := d.Set("parameters", flattenParameters(resp.Spec.Resources.Parameters)); err != nil {
		return err
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}
