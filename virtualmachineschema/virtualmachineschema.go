package virtualmachineschema

import (
	"github.com/hashicorp/terraform/helper/schema"
)
// VMSchema is Schema for VM
func VMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"spec": &schema.Schema{
			Optional: true,
			Type: schema.TypeSet,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cluster_reference": &schema.Schema{
						Optional: true,
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Optional: true,
									Type: schema.TypeString,
								},
								"name": &schema.Schema{
									Optional: true,
									Type: schema.TypeString,
								},
								"uuid": &schema.Schema{
									Optional: true,
									Type: schema.TypeString,
								},
							},
						},
					},
					"resources": &schema.Schema{
						Optional: true,
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"power_state": &schema.Schema{
									Optional: true,
									Type: schema.TypeString,
								},
								"num_sockets": &schema.Schema{
									Optional: true,
									Type: schema.TypeInt,
								},
								"memory_size_mb": &schema.Schema{
									Optional: true,
									Type: schema.TypeInt,
								},
								"gpu_list": &schema.Schema{
									Optional: true,
									Type: schema.TypeList,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"vendor": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"mode": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"device_id": &schema.Schema{
												Optional: true,
												Type: schema.TypeInt,
											},
										},
									},
								},
								"guest_customization": &schema.Schema{
									Optional: true,
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"cloud_init": &schema.Schema{
												Optional: true,
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"meta_data": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"user_data": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
													},
												},
											},
											"sysprep": &schema.Schema{
												Optional: true,
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"unattend_xml": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"install_type": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
													},
												},
											},
										},
									},
								},
								"boot_config": &schema.Schema{
									Optional: true,
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"disk_address": &schema.Schema{
												Optional: true,
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"device_index": &schema.Schema{
															Optional: true,
															Type: schema.TypeInt,
														},
														"adapter_type": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
													},
												},
											},
											"mac_address": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
										},
									},
								},
								"nic_list": &schema.Schema{
									Optional: true,
									Type: schema.TypeList,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"nic_type": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"subnet_reference": &schema.Schema{
												Optional: true,
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"uuid": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"kind": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"name": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
													},
												},
											},
											"network_function_nic_type": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"mac_address": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"ip_endpoint_list": &schema.Schema{
												Optional: true,
												Type: schema.TypeList,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"type": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"address": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
													},
												},
											},
											"network_function_chain_reference": &schema.Schema{
												Optional: true,
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"kind": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"name": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"uuid": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
													},
												},
											},
										},
									},
								},
								"num_vcpus_per_socket": &schema.Schema{
									Optional: true,
									Type: schema.TypeInt,
								},
								"parent_reference": &schema.Schema{
									Optional: true,
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"kind": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"uuid": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"name": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
										},
									},
								},
								"disk_list": &schema.Schema{
									Optional: true,
									Type: schema.TypeList,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"data_source_reference": &schema.Schema{
												Optional: true,
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"kind": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"uuid": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"name": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
													},
												},
											},
											"device_properties": &schema.Schema{
												Optional: true,
												Type: schema.TypeSet,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"device_type": &schema.Schema{
															Optional: true,
															Type: schema.TypeString,
														},
														"disk_address": &schema.Schema{
															Optional: true,
															Type: schema.TypeSet,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"device_index": &schema.Schema{
																		Optional: true,
																		Type: schema.TypeInt,
																	},
																	"adapter_type": &schema.Schema{
																		Optional: true,
																		Type: schema.TypeString,
																	},
																},
															},
														},
													},
												},
											},
											"uuid": &schema.Schema{
												Optional: true,
												Type: schema.TypeString,
											},
											"disk_size_mib": &schema.Schema{
												Optional: true,
												Type: schema.TypeInt,
											},
										},
									},
								},
							},
						},
					},
					"name": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
					"description": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
				},
			},
		},
		"api_version": &schema.Schema{
			Optional: true,
			Type: schema.TypeString,
		},
		"metadata": &schema.Schema{
			Optional: true,
			Type: schema.TypeSet,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
					"parent_reference": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
					"name": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
					"last_update_time": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
					"uuid": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
					"creation_time": &schema.Schema{
						Optional: true,
						Type: schema.TypeString,
					},
					"categories": &schema.Schema{
						Optional: true,
						Type: schema.TypeMap,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"owner_reference": &schema.Schema{
						Optional: true,
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Optional: true,
									Type: schema.TypeString,
								},
								"name": &schema.Schema{
									Optional: true,
									Type: schema.TypeString,
								},
								"uuid": &schema.Schema{
									Optional: true,
									Type: schema.TypeString,
								},
							},
						},
					},
					"entity_version": &schema.Schema{
						Optional: true,
						Type: schema.TypeInt,
					},
				},
			},
		},
	}
}
