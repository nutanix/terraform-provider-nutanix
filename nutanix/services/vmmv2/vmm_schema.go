package vmmv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func nicsElemSchemaV2() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// TODO: keep in sync with v2 VM resource NIC schema.
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backing_info": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Deprecated: "Use `nic_backing_info` instead. This field will be removed in a future release.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"VIRTIO", "E1000"}, false),
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"num_queues": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
			"network_info": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Deprecated: "Use `nic_backing_info` instead. This field will be removed in a future release.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SPAN_DESTINATION_NIC",
								"NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC",
							}, false),
						},
						"network_function_chain": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"network_function_nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"TAP", "EGRESS",
								"INGRESS",
							}, false),
						},
						"subnet": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"vlan_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"TRUNK", "ACCESS"}, false),
						},
						"trunked_vlans": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"should_allow_unknown_macs": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"ipv4_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"should_assign_ip": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Default:  defaultValue,
												},
											},
										},
									},
									"secondary_ip_address_list": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Default:  defaultValue,
												},
											},
										},
									},
								},
							},
						},
						// not visible in API reference
						"ipv4_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"learned_ip_addresses": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
													Default:  defaultValue,
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
			// new field added in v2.4.1 since backing_info and network_info are deprecated
			"nic_backing_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_ethernet_nic": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"VIRTIO", "E1000"}, false),
									},
									"mac_address": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"is_connected": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"num_queues": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
								},
							},
						},
						"sriov_nic": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sriov_profile_reference": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"host_pcie_device_reference": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"is_connected": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"mac_address": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"dp_offload_nic": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dp_offload_profile_reference": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"host_pcie_device_reference": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"is_connected": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"mac_address": {
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
			"nic_network_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_ethernet_nic_network_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nic_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											"SPAN_DESTINATION_NIC",
											"NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC",
										}, false),
									},
									"network_function_chain": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"network_function_nic_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											"TAP", "EGRESS",
											"INGRESS",
										}, false),
									},
									"subnet": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"vlan_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"TRUNK", "ACCESS"}, false),
									},
									"trunked_vlans": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"should_allow_unknown_macs": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ipv4_config": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"should_assign_ip": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"ip_address": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
															},
														},
													},
												},
												"secondary_ip_address_list": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
															},
														},
													},
												},
											},
										},
									},
									// not visible in API reference
									"ipv4_info": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"learned_ip_addresses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
															},
														},
													},
												},
											},
										},
									},
									// not visible in API reference
									"ipv6_info": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"learned_ipv6_addresses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
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
						"sriov_nic_network_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vlan_id": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"dp_offload_nic_network_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"vlan_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"TRUNK", "ACCESS"}, false),
									},
									"trunked_vlans": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"should_allow_unknown_macs": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ipv4_config": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"should_assign_ip": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
												"ip_address": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
															},
														},
													},
												},
												"secondary_ip_address_list": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
															},
														},
													},
												},
											},
										},
									},
									// not visible in API reference
									"ipv4_info": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"learned_ip_addresses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
															},
														},
													},
												},
											},
										},
									},
									// not visible in API reference
									"ipv6_info": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"learned_ipv6_addresses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Required: true,
															},
															"prefix_length": {
																Type:     schema.TypeInt,
																Optional: true,
																Default:  defaultValue,
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
		},
	}
}

func nicsElemSchemaV2WithTenantLinks() *schema.Resource {
	base := nicsElemSchemaV2()

	merged := make(map[string]*schema.Schema, len(base.Schema)+2)
	merged["tenant_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	merged["links"] = schemaForLinks()

	for key, value := range base.Schema {
		merged[key] = value
	}

	return &schema.Resource{Schema: merged}
}
