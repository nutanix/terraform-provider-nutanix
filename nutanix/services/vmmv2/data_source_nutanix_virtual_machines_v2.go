package vmmv2

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVirtualMachinesV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVirtualMachinesV4Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vms": {
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
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
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
						"memorysizebytes": {
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
						"is_cpu_hotplug_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_scsi_controller_enabled": {
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
						"availability_zone": {
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
																		"disk_size_bytes": {
																			Type:     schema.TypeInt,
																			Computed: true,
																		},
																		"disk_ext_id": {
																			Type:     schema.TypeString,
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
												"ipv4_info": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"learned_ip_addresses": {
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
		},
	}
}

func DatasourceNutanixVirtualMachinesV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	// initialize query params
	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}
	resp, err := conn.VMAPIInstance.ListVms(page, limit, filter, orderBy, selects)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching vms : %v", errorMessage["message"])
	}
	getResp := resp.Data.GetValue().([]config.Vm)

	if err := d.Set("vms", flattenVMEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVMEntities(vms []config.Vm) []interface{} {
	if len(vms) > 0 {
		vmsList := make([]interface{}, len(vms))

		for k, v := range vms {
			vm := make(map[string]interface{})

			if v.ExtId != nil {
				vm["ext_id"] = v.ExtId
			}
			if v.Name != nil {
				vm["name"] = v.Name
			}
			if v.Description != nil {
				vm["description"] = v.Description
			}
			if v.CreateTime != nil {
				t := v.CreateTime
				vm["create_time"] = t.String()
			}
			if v.UpdateTime != nil {
				t := v.UpdateTime
				vm["update_time"] = t.String()
			}
			if v.Source != nil {
				vm["source"] = flattenVMSourceReference(v.Source)
			}
			if v.NumSockets != nil {
				vm["num_sockets"] = v.NumSockets
			}
			if v.NumCoresPerSocket != nil {
				vm["num_cores_per_socket"] = v.NumCoresPerSocket
			}
			if v.NumThreadsPerCore != nil {
				vm["num_threads_per_core"] = v.NumThreadsPerCore
			}
			if v.NumNumaNodes != nil {
				vm["num_numa_nodes"] = v.NumNumaNodes
			}
			if v.MemorySizeBytes != nil {
				vm["memorysizebytes"] = v.MemorySizeBytes
			}
			if v.IsVcpuHardPinningEnabled != nil {
				vm["is_vcpu_hard_pinning_enabled"] = v.IsVcpuHardPinningEnabled
			}
			if v.IsCpuPassthroughEnabled != nil {
				vm["is_cpu_passthrough_enabled"] = v.IsCpuPassthroughEnabled
			}
			if v.EnabledCpuFeatures != nil {
				vm["enabled_cpu_features"] = flattenCPUFeature(v.EnabledCpuFeatures)
			}
			if v.IsMemoryOvercommitEnabled != nil {
				vm["is_memory_overcommit_enabled"] = v.IsMemoryOvercommitEnabled
			}
			if v.IsGpuConsoleEnabled != nil {
				vm["is_gpu_console_enabled"] = v.IsGpuConsoleEnabled
			}
			if v.IsCpuHotplugEnabled != nil {
				vm["is_cpu_hotplug_enabled"] = v.IsCpuHotplugEnabled
			}
			if v.IsScsiControllerEnabled != nil {
				vm["is_scsi_controller_enabled"] = v.IsScsiControllerEnabled
			}
			if v.GenerationUuid != nil {
				vm["generation_uuid"] = v.GenerationUuid
			}
			if v.BiosUuid != nil {
				vm["bios_uuid"] = v.BiosUuid
			}
			if v.Categories != nil {
				vm["categories"] = flattenCategoryReference(v.Categories)
			}
			if v.OwnershipInfo != nil {
				vm["ownership_info"] = flattenOwnershipInfo(v.OwnershipInfo)
			}
			if v.Host != nil {
				vm["host"] = flattenHostReference(v.Host)
			}
			if v.Cluster != nil {
				vm["cluster"] = flattenClusterReference(v.Cluster)
			}
			if v.GuestCustomization != nil {
				vm["guest_customization"] = flattenGuestCustomizationParams(v.GuestCustomization)
			}
			if v.GuestTools != nil {
				vm["guest_tools"] = flattenGuestTools(v.GuestTools)
			}
			if v.HardwareClockTimezone != nil {
				vm["hardware_clock_timezone"] = v.HardwareClockTimezone
			}
			if v.IsBrandingEnabled != nil {
				vm["is_branding_enabled"] = v.IsBrandingEnabled
			}
			if v.BootConfig != nil {
				vm["boot_config"] = flattenOneOfVMBootConfig(v.BootConfig)
			}
			if v.IsVgaConsoleEnabled != nil {
				vm["is_vga_console_enabled"] = v.IsVgaConsoleEnabled
			}
			if v.MachineType != nil {
				vm["machine_type"] = flattenMachineType(v.MachineType)
			}
			if v.PowerState != nil {
				vm["power_state"] = flattenPowerState(v.PowerState)
			}
			if v.VtpmConfig != nil {
				vm["vtpm_config"] = flattenVtpmConfig(v.VtpmConfig)
			}
			if v.IsAgentVm != nil {
				vm["is_agent_vm"] = v.IsAgentVm
			}
			if v.ApcConfig != nil {
				vm["apc_config"] = flattenApcConfig(v.ApcConfig)
			}
			if v.StorageConfig != nil {
				vm["storage_config"] = flattenADSFVmStorageConfig(v.StorageConfig)
			}
			if v.Disks != nil {
				vm["disks"] = flattenDisk(v.Disks)
			}
			if v.CdRoms != nil {
				vm["cd_roms"] = flattenCdRom(v.CdRoms)
			}
			if v.Nics != nil {
				vm["nics"] = flattenNic(v.Nics)
			}
			if v.Gpus != nil {
				vm["gpus"] = flattenGpu(v.Gpus)
			}
			if v.SerialPorts != nil {
				vm["serial_ports"] = flattenSerialPort(v.SerialPorts)
			}
			if v.ProtectionType != nil {
				vm["protection_type"] = flattenProtectionType(v.ProtectionType)
			}
			if v.ProtectionPolicyState != nil {
				vm["protection_policy_state"] = flattenProtectionPolicyState(v.ProtectionPolicyState)
			}

			vmsList[k] = vm
		}
		return vmsList
	}
	return nil
}
