package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v16/models/prism/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/prism/v4/config"
	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/models/vmm/v4/content"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixTemplatesVersionV4() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTemplatesVersionV4Create,
		ReadContext:   ResourceNutanixTemplatesVersionV4Read,
		UpdateContext: ResourceNutanixTemplatesVersionV4Update,
		DeleteContext: ResourceNutanixTemplatesVersionV4Delete,
		Schema: map[string]*schema.Schema{
			"template_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_version_spec": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version_description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vm_spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"entity_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"num_sockets": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"num_cores_per_socket": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"num_threads_per_core": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"num_numa_nodes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"memory_size_bytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_vcpu_hard_pinning_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_cpu_passthrough_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enabled_cpu_features": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"is_memory_overcommit_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_gpu_console_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"generation_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"bios_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"categories": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"ownership_info": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"owner": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ext_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"host": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"cluster": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									// not present in API reference
									// "availability_zone":   {},
									"guest_customization": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"config": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"sysprep": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"install_type": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"sysprep_script": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"unattend_xml": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"value": {
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"custom_key_values": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"key_value_pairs": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"name": {
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																											"value": {
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
															"cloud_init": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"datasource_type": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"metadata": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"cloud_init_script": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"user_data": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"value": {
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"custom_keys": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"key_value_pairs": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"name": {
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																											"value": {
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
														},
													},
												},
											},
										},
									},
									"guest_tools": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_installed": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_iso_inserted": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"available_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"guest_os_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_reachable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_vss_snapshot_capable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_vm_mobility_drivers_installed": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"capabilities": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"hardware_clock_timezone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_branding_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"boot_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"legacy_boot": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"boot_device": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"boot_device_disk": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"disk_address": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bus_type": {
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																								"index": {
																									Type:     schema.TypeInt,
																									Computed: true,
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"boot_device_nic": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mac_address": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																	},
																},
															},
															"boot_order": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"uefi_boot": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_secure_boot_enabled": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"nvram_device": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"backing_storage_info": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"disk_ext_id": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"disk_size_bytes": {
																						Type:     schema.TypeInt,
																						Computed: true,
																					},
																					"storage_container": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"ext_id": {
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"storage_config": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"is_flash_mode_enabled": {
																									Type:     schema.TypeBool,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"data_source": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"reference": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"image_reference": {
																												Type:     schema.TypeList,
																												Computed: true,
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"image_ext_id": {
																															Type:     schema.TypeString,
																															Computed: true,
																														},
																													},
																												},
																											},
																											"vm_disk_reference": {
																												Type:     schema.TypeList,
																												Computed: true,
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"disk_ext_id": {
																															Type:     schema.TypeString,
																															Computed: true,
																														},
																														"disk_address": {
																															Type:     schema.TypeList,
																															Computed: true,
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"bus_type": {
																																		Type:     schema.TypeString,
																																		Computed: true,
																																	},
																																	"index": {
																																		Type:     schema.TypeInt,
																																		Computed: true,
																																	},
																																},
																															},
																														},
																														"vm_reference": {
																															Type:     schema.TypeList,
																															Computed: true,
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"ext_id": {
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
																					"is_migration_in_progress": {
																						Type:     schema.TypeBool,
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
									"is_vga_console_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"machine_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"power_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vtpm_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_vtpm_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"version": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"is_agent_vm": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"apc_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_apc_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"cpu_model": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ext_id": {
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
											},
										},
									},
									"storage_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_flash_mode_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"qos_config": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"throttled_iops": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"disks": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"disk_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bus_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"index": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"backing_info": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"vm_disk": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"disk_ext_id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																		"disk_size_bytes": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"storage_container": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"ext_id": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"storage_config": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"is_flash_mode_enabled": {
																						Type:     schema.TypeBool,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"data_source": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"reference": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_reference": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"image_ext_id": {
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																										},
																									},
																								},
																								"vm_disk_reference": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"disk_ext_id": {
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																											"disk_address": {
																												Type:     schema.TypeList,
																												Computed: true,
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bus_type": {
																															Type:     schema.TypeString,
																															Computed: true,
																														},
																														"index": {
																															Type:     schema.TypeInt,
																															Computed: true,
																														},
																													},
																												},
																											},
																											"vm_reference": {
																												Type:     schema.TypeList,
																												Computed: true,
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"ext_id": {
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
																		"is_migration_in_progress": {
																			Type:     schema.TypeBool,
																			Computed: true,
																		},
																	},
																},
															},
															"adfs_volume_group_reference": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"volume_group_ext_id": {
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
									"cd_roms": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"disk_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bus_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"index": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"backing_info": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_ext_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"disk_size_bytes": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"storage_container": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ext_id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
															"storage_config": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_flash_mode_enabled": {
																			Type:     schema.TypeBool,
																			Computed: true,
																		},
																	},
																},
															},
															"data_source": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"reference": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_reference": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_ext_id": {
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"vm_disk_reference": {
																						Type:     schema.TypeList,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"disk_ext_id": {
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																								"disk_address": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bus_type": {
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																											"index": {
																												Type:     schema.TypeInt,
																												Computed: true,
																											},
																										},
																									},
																								},
																								"vm_reference": {
																									Type:     schema.TypeList,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"ext_id": {
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
															"is_migration_in_progress": {
																Type:     schema.TypeBool,
																Computed: true,
															},
														},
													},
												},
												"iso_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"nics": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"backing_info": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"model": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"mac_address": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"is_connected": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"num_queues": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"network_info": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"nic_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"network_function_chain": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ext_id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
															"network_function_nic_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"subnet": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ext_id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
															"vlan_mode": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"trunked_vlans": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Schema{
																	Type: schema.TypeInt,
																},
															},
															"should_allow_unknown_macs": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"ipv4_config": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"should_assign_ip": {
																			Type:     schema.TypeBool,
																			Computed: true,
																		},
																		"ip_address": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"value": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"prefix_length": {
																						Type:     schema.TypeInt,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"secondary_ip_address_list": {
																			Type:     schema.TypeList,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"value": {
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"prefix_length": {
																						Type:     schema.TypeInt,
																						Computed: true,
																					},
																				},
																			},
																		},
																	},
																},
															},
															// not visible in API reference
															// "ipv4Info":   {},
														},
													},
												},
											},
										},
									},
									"gpus": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mode": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"device_id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"vendor": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"pci_address": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"segment": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"bus": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"device": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"func": {
																Type:     schema.TypeInt,
																Computed: true,
															},
														},
													},
												},
												"guest_driver_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"frame_buffer_size_bytes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"num_virtual_display_heads": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"fraction": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"serial_ports": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"is_connected": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"index": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"protection_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protection_policy_state": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"policy": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ext_id": {
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
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"idp_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"first_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"middle_initial": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"locale": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_force_reset_password_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"additional_attributes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"version_source": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"template_version_reference": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"override_vm_config": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"num_sockets": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"num_cores_per_socket": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"num_threads_per_core": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"memory_size_bytes": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"nics": {
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
																		"backing_info": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Computed: true,
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
																			Type:     schema.TypeList,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"nic_type": {
																						Type:     schema.TypeString,
																						Optional: true,
																						Computed: true,
																						ValidateFunc: validation.StringInSlice([]string{"SPAN_DESTINATION_NIC",
																							"NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC"}, false),
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
																						ValidateFunc: validation.StringInSlice([]string{"TAP", "EGRESS",
																							"INGRESS"}, false),
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
																												Computed: true,
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
																												Computed: true,
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					// not visible in API reference
																					// "ipv4Info":   {},
																				},
																			},
																		},
																	},
																},
															},
															"guest_customization": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"config": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"sysprep": {
																						Type:     schema.TypeList,
																						Optional: true,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"install_type": {
																									Type:     schema.TypeString,
																									Optional: true,
																									Computed: true,
																								},
																								"sysprep_script": {
																									Type:     schema.TypeList,
																									Optional: true,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"unattend_xml": {
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
																													},
																												},
																											},
																											"custom_key_values": {
																												Type:     schema.TypeList,
																												Optional: true,
																												Computed: true,
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"key_value_pairs": {
																															Type:     schema.TypeList,
																															Optional: true,
																															Computed: true,
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"name": {
																																		Type:     schema.TypeString,
																																		Optional: true,
																																		Computed: true,
																																	},
																																	"value": {
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
																						},
																					},
																					"cloud_init": {
																						Type:     schema.TypeList,
																						Optional: true,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"datasource_type": {
																									Type:     schema.TypeString,
																									Optional: true,
																									Default:  "CONFIG_DRIVE_V2",
																								},
																								"metadata": {
																									Type:     schema.TypeString,
																									Optional: true,
																									Computed: true,
																								},
																								"cloud_init_script": {
																									Type:     schema.TypeList,
																									Optional: true,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"user_data": {
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
																													},
																												},
																											},
																											"custom_keys": {
																												Type:     schema.TypeList,
																												Optional: true,
																												Computed: true,
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"key_value_pairs": {
																															Type:     schema.TypeList,
																															Optional: true,
																															Computed: true,
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"name": {
																																		Type:     schema.TypeString,
																																		Optional: true,
																																		Computed: true,
																																	},
																																	"value": {
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
								},
							},
						},
						"version_source_discriminator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_active_version": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"is_gc_override_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixTemplatesVersionV4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	templateExtID := d.Get("template_ext_id").(string)

	readResp, err := conn.TemplatesAPIInstance.GetTemplateById(utils.StringPtr(templateExtID))
	if err != nil {
		return diag.Errorf("error while fetching template : %v", err)
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	respTemplates := readResp.Data.GetValue().(import5.Template)

	updateSpec := &import5.Template{}
	updateSpec.ExtId = respTemplates.ExtId
	updateSpec.TemplateVersionSpec = respTemplates.TemplateVersionSpec

	if d.HasChange("template_version_spec") {
		updateSpec.TemplateVersionSpec = expandTemplateVersionSpec(d.Get("template_version_spec"))
	}

	respUpdate, err := conn.TemplatesAPIInstance.UpdateTemplateById(utils.StringPtr(templateExtID), updateSpec, args)
	if err != nil {
		return diag.Errorf("error while creating template version : %v", err)
	}

	TaskRef := respUpdate.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Template to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template version(%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching vm UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[1].ExtId
	d.SetId(*uuid)
	return ResourceNutanixTemplatesVersionV4Read(ctx, d, meta)
}

func ResourceNutanixTemplatesVersionV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	templateExtID := d.Get("template_ext_id").(string)
	resp, err := conn.TemplatesAPIInstance.GetTemplateVersionById(utils.StringPtr(templateExtID), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching template version: %v", err)
	}

	getResp := resp.Data.GetValue().(import5.TemplateVersionSpec)

	if err := d.Set("template_version_spec", flattenTemplateVersionSpec(&getResp)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixTemplatesVersionV4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixTemplatesVersionV4Create(ctx, d, meta)
}

func ResourceNutanixTemplatesVersionV4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	templateExtID := d.Get("template_ext_id").(string)
	resp, err := conn.TemplatesAPIInstance.DeleteTemplateVersionById(utils.StringPtr(templateExtID), utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting template verison : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template version(%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}
