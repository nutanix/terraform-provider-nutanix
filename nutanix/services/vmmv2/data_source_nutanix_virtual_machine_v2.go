package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import4 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVirtualMachineV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVirtualMachineV4Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func DatasourceNutanixVirtualMachineV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	extID := d.Get("ext_id")

	resp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(extID.(string)))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching vm : %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().(config.Vm)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if getResp.CreateTime != nil {
		t := getResp.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.UpdateTime != nil {
		t := getResp.UpdateTime
		if err := d.Set("update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("source", flattenVMSourceReference(getResp.Source)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_sockets", getResp.NumSockets); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_cores_per_socket", getResp.NumCoresPerSocket); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_threads_per_core", getResp.NumThreadsPerCore); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_numa_nodes", getResp.NumNumaNodes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memorysizebytes", getResp.MemorySizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_vcpu_hard_pinning_enabled", getResp.IsVcpuHardPinningEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_cpu_passthrough_enabled", getResp.IsCpuPassthroughEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled_cpu_features", flattenCPUFeature(getResp.EnabledCpuFeatures)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_memory_overcommit_enabled", getResp.IsMemoryOvercommitEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_gpu_console_enabled", getResp.IsGpuConsoleEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_cpu_hotplug_enabled", getResp.IsCpuHotplugEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_scsi_controller_enabled", getResp.IsScsiControllerEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("generation_uuid", getResp.GenerationUuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("bios_uuid", getResp.BiosUuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", flattenCategoryReference(getResp.Categories)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ownership_info", flattenOwnershipInfo(getResp.OwnershipInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host", flattenHostReference(getResp.Host)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster", flattenClusterReference(getResp.Cluster)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_customization", flattenGuestCustomizationParams(getResp.GuestCustomization)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_tools", flattenGuestTools(getResp.GuestTools)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hardware_clock_timezone", getResp.HardwareClockTimezone); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_branding_enabled", getResp.IsBrandingEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("boot_config", flattenOneOfVMBootConfig(getResp.BootConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_vga_console_enabled", getResp.IsVgaConsoleEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("machine_type", flattenMachineType(getResp.MachineType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("power_state", flattenPowerState(getResp.PowerState)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vtpm_config", flattenVtpmConfig(getResp.VtpmConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_agent_vm", getResp.IsAgentVm); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("apc_config", flattenApcConfig(getResp.ApcConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_config", flattenADSFVmStorageConfig(getResp.StorageConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disks", flattenDisk(getResp.Disks)); err != nil {
		log.Printf("[ERROR] error while setting disks : %v", err)
		return diag.FromErr(err)
	}
	if err := d.Set("cd_roms", flattenCdRom(getResp.CdRoms)); err != nil {
		log.Printf("[ERROR] error while setting cd_roms : %v", err)
		return diag.FromErr(err)
	}
	if err := d.Set("nics", flattenNic(getResp.Nics)); err != nil {
		log.Printf("[ERROR] error while setting nics : %v", err)
		return diag.FromErr(err)
	}
	if err := d.Set("gpus", flattenGpu(getResp.Gpus)); err != nil {
		log.Printf("[ERROR] error while setting gpus : %v", err)
		return diag.FromErr(err)
	}
	if err := d.Set("serial_ports", flattenSerialPort(getResp.SerialPorts)); err != nil {
		log.Printf("[ERROR] error while setting serial_ports : %v", err)
		return diag.FromErr(err)
	}
	if err := d.Set("protection_type", flattenProtectionType(getResp.ProtectionType)); err != nil {
		log.Printf("[ERROR] error while setting protection_type : %v", err)
		return diag.FromErr(err)
	}
	if err := d.Set("protection_policy_state", flattenProtectionPolicyState(getResp.ProtectionPolicyState)); err != nil {
		log.Printf("[ERROR] error while setting protection_policy_state : %v", err)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully read vm with ext_id %s", extID)
	d.SetId(*getResp.ExtId)
	return nil
}

func flattenVMSourceReference(ref *config.VmSourceReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.EntityType != nil {
			refs["entity_type"] = flattenVMSourceReferenceEntityType(ref.EntityType)
		}
		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenVMSourceReferenceEntityType(ent *config.VmSourceReferenceEntityType) string {
	if ent != nil {
		const two, three = 2, 3
		if *ent == config.VmSourceReferenceEntityType(two) {
			return "VM"
		}
		if *ent == config.VmSourceReferenceEntityType(three) {
			return "VM_RECOVERY_POINT"
		}
	}
	return "UNKNOWN"
}

func flattenCPUFeature(cfg []config.CpuFeature) []interface{} {
	if len(cfg) > 0 {
		cfgList := make([]interface{}, len(cfg))
		const two = 2
		for _, v := range cfg {
			if v == config.CpuFeature(two) {
				cfgList = append(cfgList, "HARDWARE_VIRTUALIZATION")
			}
		}
		return cfgList
	}
	return nil
}

func flattenCategoryReference(ctg []config.CategoryReference) []interface{} {
	if len(ctg) > 0 {
		ctgList := make([]interface{}, len(ctg))

		for k, v := range ctg {
			ctgs := make(map[string]interface{})

			if v.ExtId != nil {
				ctgs["ext_id"] = v.ExtId
			}
			ctgList[k] = ctgs
		}
		return ctgList
	}
	return nil
}

func flattenOwnershipInfo(own *config.OwnershipInfo) []map[string]interface{} {
	if own != nil {
		ownList := make([]map[string]interface{}, 0)

		owns := make(map[string]interface{})

		if own.Owner != nil {
			owns["owner"] = flattenOwnerReference(own.Owner)
		}

		ownList = append(ownList, owns)
		return ownList
	}
	return nil
}

func flattenOwnerReference(ref *config.OwnerReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenHostReference(ref *config.HostReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenClusterReference(ref *config.ClusterReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenAvailabilityZoneReference(zone *config.AvailabilityZoneReference) []map[string]interface{} {
	if zone == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ext_id": utils.StringValue(zone.ExtId),
		},
	}
}

func flattenGuestCustomizationParams(gst *config.GuestCustomizationParams) []map[string]interface{} {
	if gst != nil {
		gstList := make([]map[string]interface{}, 0)

		gsts := make(map[string]interface{})

		if gst.Config != nil {
			gsts["config"] = flattenOneOfGuestCustomizationParamsConfig(gst.Config)
		}

		gstList = append(gstList, gsts)
		aJSON, _ := json.Marshal(gstList)
		log.Printf("[DEBUG] flattenGuestCustomizationParams[config]: %s", string(aJSON))
		return gstList
	}
	return nil
}

func flattenOneOfGuestCustomizationParamsConfig(cfg *config.OneOfGuestCustomizationParamsConfig) []map[string]interface{} {
	if cfg != nil {
		sysCfg := make(map[string]interface{})
		sysCfgList := make([]map[string]interface{}, 0)
		cloudCfg := make(map[string]interface{})
		cloudCfgList := make([]map[string]interface{}, 0)

		if *cfg.ObjectType_ == "vmm.v4.ahv.config.Sysprep" {
			sysObj := cfg.GetValue().(config.Sysprep)
			sysprepObj := make(map[string]interface{})
			sysprepObjList := make([]map[string]interface{}, 0)

			sysprepObj["install_type"] = flattenInstallType(sysObj.InstallType)

			sysprepObj["sysprep_script"] = flattenOneOfSysprepSysprepScript(sysObj.SysprepScript)

			aJSON, _ := json.Marshal(sysprepObj)
			log.Printf("[DEBUG] flattenOneOfGuestCustomizationParamsConfig sysprep: %s", string(aJSON))

			sysprepObjList = append(sysprepObjList, sysprepObj)
			sysCfg["sysprep"] = sysprepObjList

			sysCfgList = append(sysCfgList, sysCfg)

			return sysCfgList
		}
		cloudInitObj := make(map[string]interface{})
		cloudInitObjList := make([]map[string]interface{}, 0)
		cloudObj := cfg.GetValue().(config.CloudInit)

		cloudInitObj["datasource_type"] = flattenCloudInitDataSourceType(cloudObj.DatasourceType)
		cloudInitObj["metadata"] = cloudObj.Metadata
		cloudInitObj["cloud_init_script"] = flattenOneOfCloudInitCloudInitScript(cloudObj.CloudInitScript)

		cloudInitObjList = append(cloudInitObjList, cloudInitObj)
		cloudCfg["cloud_init"] = cloudInitObjList

		cloudCfgList = append(cloudCfgList, cloudCfg)

		return cloudCfgList
	}
	return nil
}

func flattenInstallType(pr *config.InstallType) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == config.InstallType(two) {
			return "FRESH"
		}
		if *pr == config.InstallType(three) {
			return "PREPARED"
		}
	}
	return "UNKNOWN"
}

func flattenOneOfSysprepSysprepScript(cfg *config.OneOfSysprepSysprepScript) []map[string]interface{} {
	if cfg != nil {
		unattendCfg := make(map[string]interface{})
		unattendCfgList := make([]map[string]interface{}, 0)
		customKeyValCfg := make(map[string]interface{})
		customKeyValCfgList := make([]map[string]interface{}, 0)

		if *cfg.ObjectType_ == "vmm.v4.ahv.config.Unattendxml" {
			unattendXML := make(map[string]interface{})
			unattendXMLList := make([]map[string]interface{}, 0)
			xmlCfg := cfg.GetValue().(config.Unattendxml)

			unattendXML["value"] = xmlCfg.Value

			unattendXMLList = append(unattendXMLList, unattendXML)

			unattendCfg["unattend_xml"] = unattendXMLList

			unattendCfgList = append(unattendCfgList, unattendCfg)

			return unattendCfgList
		}
		customObj := make(map[string]interface{})
		customObjList := make([]map[string]interface{}, 0)
		customCfg := cfg.GetValue().(config.CustomKeyValues)

		customObj["key_value_pairs"] = flattenCustomKVPair(customCfg.KeyValuePairs)

		customObjList = append(customObjList, customObj)
		customKeyValCfg["custom_key_values"] = customObjList

		customKeyValCfgList = append(customKeyValCfgList, customKeyValCfg)
		return customKeyValCfgList
	}
	return nil
}

func flattenCloudInitDataSourceType(dsType *config.CloudInitDataSourceType) string {
	if dsType != nil {
		const two = 2

		if *dsType == config.CloudInitDataSourceType(two) {
			return "CONFIG_DRIVE_V2"
		}
	}
	return "UNKNOWN"
}

func flattenOneOfCloudInitCloudInitScript(cfg *config.OneOfCloudInitCloudInitScript) []map[string]interface{} {
	if cfg != nil {
		cfgList := make([]map[string]interface{}, 0)

		userDataCfg := make(map[string]interface{})
		customKeyValCfg := make(map[string]interface{})

		if *cfg.ObjectType_ == "vmm.v4.ahv.config.Userdata" {
			userdataObj := make(map[string]interface{})
			userdataObjList := make([]map[string]interface{}, 0)
			userdataVals := cfg.GetValue().(config.Userdata)

			userdataObj["value"] = userdataVals.Value

			userdataObjList = append(userdataObjList, userdataObj)
			userDataCfg["user_data"] = userdataObjList
			cfgList = append(cfgList, userDataCfg)
		} else {
			customObj := make(map[string]interface{})
			customObjList := make([]map[string]interface{}, 0)

			customCfg := cfg.GetValue().(config.CustomKeyValues)

			customObj["key_value_pairs"] = flattenCustomKVPair(customCfg.KeyValuePairs)

			customObjList = append(customObjList, customObj)
			customKeyValCfg["custom_key_values"] = customObjList
			cfgList = append(cfgList, customKeyValCfg)
		}

		return cfgList
	}
	return nil
}

func flattenGuestTools(pr *config.GuestTools) []map[string]interface{} {
	if pr != nil {
		toolsList := make([]map[string]interface{}, 0)

		tools := make(map[string]interface{})

		if pr.Version != nil {
			tools["version"] = pr.Version
		}
		if pr.IsInstalled != nil {
			tools["is_installed"] = pr.IsInstalled
		}
		if pr.IsIsoInserted != nil {
			tools["is_iso_inserted"] = pr.IsIsoInserted
		}
		if pr.AvailableVersion != nil {
			tools["available_version"] = pr.AvailableVersion
		}
		if pr.GuestOsVersion != nil {
			tools["guest_os_version"] = pr.GuestOsVersion
		}
		if pr.IsReachable != nil {
			tools["is_reachable"] = pr.IsReachable
		}
		if pr.IsVssSnapshotCapable != nil {
			tools["is_vss_snapshot_capable"] = pr.IsVssSnapshotCapable
		}
		if pr.IsVmMobilityDriversInstalled != nil {
			tools["is_vm_mobility_drivers_installed"] = pr.IsVmMobilityDriversInstalled
		}
		if pr.IsEnabled != nil {
			tools["is_enabled"] = pr.IsEnabled
		}
		if pr.Capabilities != nil {
			tools["capabilities"] = flattenNgtCapability(pr.Capabilities)
		}

		toolsList = append(toolsList, tools)

		return toolsList
	}
	return nil
}

func flattenNgtCapability(pr []config.NgtCapability) []interface{} {
	if len(pr) > 0 {
		ngtList := make([]interface{}, len(pr))
		const two, three = 2, 3
		for k, v := range pr {
			ss := new(string)

			if v == config.NgtCapability(two) {
				ss = utils.StringPtr("SELF_SERVICE_RESTORE")
			}
			if v == config.NgtCapability(three) {
				ss = utils.StringPtr("VSS_SNAPSHOT")
			}
			ngtList[k] = ss
		}
		return ngtList
	}
	return nil
}

func flattenOneOfVMBootConfig(pr *config.OneOfVmBootConfig) []map[string]interface{} {
	if pr != nil {
		legacyBootCfg := make(map[string]interface{})
		legacyBootCfgList := make([]map[string]interface{}, 0)
		uefiBootCfg := make(map[string]interface{})
		uefiBootCfgList := make([]map[string]interface{}, 0)

		if *pr.ObjectType_ == "vmm.v4.ahv.config.LegacyBoot" {
			legacyObj := make(map[string]interface{})
			legacyObjList := make([]map[string]interface{}, 0)
			legacyVals := pr.GetValue().(config.LegacyBoot)

			legacyObj["boot_device"] = flattenOneOfLegacyBootBootDevice(legacyVals.BootDevice)
			legacyObj["boot_order"] = flattenBootDeviceType(legacyVals.BootOrder)

			legacyBootCfgList = append(legacyBootCfgList, legacyObj)
			legacyBootCfg["legacy_boot"] = legacyBootCfgList

			legacyObjList = append(legacyObjList, legacyBootCfg)
			return legacyObjList
		}
		uefiObj := make(map[string]interface{})
		uefiObjList := make([]map[string]interface{}, 0)
		uefiVals := pr.GetValue().(config.UefiBoot)

		uefiObj["is_secure_boot_enabled"] = uefiVals.IsSecureBootEnabled
		uefiObj["nvram_device"] = flattenNvramDevice(uefiVals.NvramDevice)

		uefiObjList = append(uefiObjList, uefiObj)
		uefiBootCfg["uefi_boot"] = uefiObjList
		uefiBootCfgList = append(uefiBootCfgList, uefiBootCfg)

		return uefiBootCfgList
	}
	return nil
}

func flattenOneOfLegacyBootBootDevice(cfg *config.OneOfLegacyBootBootDevice) []map[string]interface{} {
	if cfg != nil {
		// bootDeviceList := make([]map[string]interface{}, 0)
		deviceDisk := make(map[string]interface{})
		deviceDiskList := make([]map[string]interface{}, 0)
		deviceNic := make(map[string]interface{})
		deviceNicList := make([]map[string]interface{}, 0)

		if *cfg.ObjectType_ == "vmm.v4.ahv.config.BootDeviceDisk" {
			deviceDiskObj := make(map[string]interface{})
			deviceDiskObjList := make([]map[string]interface{}, 0)
			deviceDiskVal := cfg.GetValue().(config.BootDeviceDisk)

			if deviceDiskVal.DiskAddress != nil {
				deviceDiskObj["disk_address"] = flattenDiskAddress(deviceDiskVal.DiskAddress)
			}

			deviceDiskObjList = append(deviceDiskObjList, deviceDiskObj)

			deviceDisk["boot_device_disk"] = deviceDiskObjList
			deviceDiskList = append(deviceDiskList, deviceDisk)
			return deviceDiskList
		}
		deviceNicObj := make(map[string]interface{})
		deviceNicObjList := make([]map[string]interface{}, 0)
		deviceNicVal := cfg.GetValue().(config.BootDeviceNic)

		deviceNicObj["mac_address"] = deviceNicVal.MacAddress
		deviceNicObjList = append(deviceNicObjList, deviceNicObj)

		deviceNic["boot_device_nic"] = deviceNicObjList

		deviceNicList = append(deviceNicList, deviceNic)

		return deviceNicList
		// return bootDeviceList
	}
	return nil
}

func flattenDiskAddress(pr *config.DiskAddress) []map[string]interface{} {
	if pr != nil {
		diskList := make([]map[string]interface{}, 0)

		disk := make(map[string]interface{})

		if pr.BusType != nil {
			disk["bus_type"] = flattenDiskBusType(pr.BusType)
		}
		if pr.Index != nil {
			disk["index"] = pr.Index
		}

		diskList = append(diskList, disk)
		return diskList
	}
	return nil
}

func flattenDiskBusType(pr *config.DiskBusType) string {
	const two, three, four, five, six = 2, 3, 4, 5, 6
	if pr != nil {
		if *pr == config.DiskBusType(two) {
			return "SCSI"
		}
		if *pr == config.DiskBusType(three) {
			return "IDE"
		}
		if *pr == config.DiskBusType(four) {
			return "PCI"
		}
		if *pr == config.DiskBusType(five) {
			return "SATA"
		}
		if *pr == config.DiskBusType(six) {
			return "SPAPR"
		}
	}
	return "UNKNOWN"
}

func flattenBootDeviceType(pr []config.BootDeviceType) []interface{} {
	if len(pr) > 0 {
		bootDeviceList := make([]interface{}, len(pr))
		const two, three, four = 2, 3, 4
		for k, v := range pr {
			ss := new(string)

			if v == config.BootDeviceType(two) {
				ss = utils.StringPtr("CDROM")
			}
			if v == config.BootDeviceType(three) {
				ss = utils.StringPtr("DISK")
			}
			if v == config.BootDeviceType(four) {
				ss = utils.StringPtr("NETWORK")
			}
			bootDeviceList[k] = ss
		}
		return bootDeviceList
	}
	return nil
}

func flattenNvramDevice(pr *config.NvramDevice) []map[string]interface{} {
	if pr != nil {
		ramList := make([]map[string]interface{}, 0)

		rams := make(map[string]interface{})

		if pr.BackingStorageInfo != nil {
			rams["backing_storage_info"] = flattenVMDisk(pr.BackingStorageInfo)
		}
		ramList = append(ramList, rams)
		return ramList
	}
	return nil
}

func flattenVMDisk(pr *config.VmDisk) []map[string]interface{} {
	if pr != nil {
		vmDiskList := make([]map[string]interface{}, 0)

		disks := make(map[string]interface{})

		if pr.DiskSizeBytes != nil {
			disks["disk_size_bytes"] = pr.DiskSizeBytes
		}
		if pr.StorageContainer != nil {
			disks["storage_container"] = flattenVMDiskContainerReference(pr.StorageContainer)
		}
		if pr.StorageConfig != nil {
			disks["storage_config"] = flattenVMDiskStorageConfig(pr.StorageConfig)
		}
		if pr.DataSource != nil {
			disks["data_source"] = flattenDataSource(pr.DataSource)
		}
		if pr.DiskExtId != nil {
			disks["disk_ext_id"] = pr.DiskExtId
		}
		if pr.IsMigrationInProgress != nil {
			disks["is_migration_in_progress"] = pr.IsMigrationInProgress
		}

		vmDiskList = append(vmDiskList, disks)
		return vmDiskList
	}
	return nil
}

func flattenVMDiskContainerReference(ref *config.VmDiskContainerReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenVMDiskStorageConfig(ref *config.VmDiskStorageConfig) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.IsFlashModeEnabled != nil {
			refs["is_flash_mode_enabled"] = ref.IsFlashModeEnabled
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenDataSource(ref *config.DataSource) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.Reference != nil {
			refs["reference"] = flattenOneOfDataSourceReference(ref.Reference)
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenOneOfDataSourceReference(pr *config.OneOfDataSourceReference) []map[string]interface{} {
	if pr != nil {
		vmDiskRef := make(map[string]interface{})
		vmDiskRefList := make([]map[string]interface{}, 0)
		imageRef := make(map[string]interface{})
		imageRefList := make([]map[string]interface{}, 0)

		if *pr.ObjectType_ == "vmm.v4.ahv.config.VmDiskReference" {
			vmDiskObj := make(map[string]interface{})
			vmDiskObjList := make([]map[string]interface{}, 0)
			vmDiskVal := pr.GetValue().(config.VmDiskReference)

			vmDiskObj["disk_address"] = flattenDiskAddress(vmDiskVal.DiskAddress)
			vmDiskObj["disk_ext_id"] = vmDiskVal.DiskExtId
			vmDiskObj["vm_reference"] = flattenVMReference(vmDiskVal.VmReference)

			vmDiskObjList = append(vmDiskObjList, vmDiskObj)
			vmDiskRef["vm_disk_reference"] = vmDiskObjList
			vmDiskRefList = append(vmDiskRefList, vmDiskRef)

			return vmDiskRefList
		}
		imageObj := make(map[string]interface{})
		imageObjList := make([]map[string]interface{}, 0)
		imageVal := pr.GetValue().(config.ImageReference)

		imageObj["image_ext_id"] = imageVal.ImageExtId

		imageObjList = append(imageObjList, imageObj)
		imageRef["image_reference"] = imageObjList
		imageRefList = append(imageRefList, imageRef)

		return imageRefList
	}
	return nil
}

func flattenVMReference(ref *config.VmReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenMachineType(pr *config.MachineType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == config.MachineType(two) {
			return "PC"
		}
		if *pr == config.MachineType(three) {
			return "PSERIES"
		}
		if *pr == config.MachineType(four) {
			return "Q35"
		}
	}
	return "UNKNOWN"
}

func flattenVtpmConfig(pr *config.VtpmConfig) []map[string]interface{} {
	if pr != nil {
		vtpmList := make([]map[string]interface{}, 0)
		vtpm := make(map[string]interface{})

		if pr.IsVtpmEnabled != nil {
			vtpm["is_vtpm_enabled"] = pr.IsVtpmEnabled
		}
		if pr.Version != nil {
			vtpm["version"] = pr.Version
		}
		vtpmList = append(vtpmList, vtpm)
		return vtpmList
	}
	return nil
}

func flattenApcConfig(pr *config.ApcConfig) []map[string]interface{} {
	if pr != nil {
		cfgList := make([]map[string]interface{}, 0)
		cfg := make(map[string]interface{})

		if pr.IsApcEnabled != nil {
			cfg["is_apc_enabled"] = pr.IsApcEnabled
		}
		if pr.CpuModel != nil {
			cfg["cpu_model"] = flattenCPUModelReference(pr.CpuModel)
		}

		cfgList = append(cfgList, cfg)
		return cfgList
	}
	return nil
}

func flattenCPUModelReference(ref *config.CpuModelReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		if ref.Name != nil {
			refs["name"] = ref.Name
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenADSFVmStorageConfig(pr *config.ADSFVmStorageConfig) []map[string]interface{} {
	if pr != nil {
		cfgList := make([]map[string]interface{}, 0)

		cfg := make(map[string]interface{})

		if pr.IsFlashModeEnabled != nil {
			cfg["is_flash_mode_enabled"] = pr.IsFlashModeEnabled
		}
		if pr.QosConfig != nil {
			cfg["qos_config"] = flattenQosConfig(pr.QosConfig)
		}

		cfgList = append(cfgList, cfg)
		return cfgList
	}
	return nil
}

func flattenQosConfig(pr *config.QosConfig) []map[string]interface{} {
	if pr != nil {
		cfgList := make([]map[string]interface{}, 0)

		cfg := make(map[string]interface{})

		if pr.ThrottledIops != nil {
			cfg["throttled_iops"] = pr.ThrottledIops
		}

		cfgList = append(cfgList, cfg)
		return cfgList
	}
	return nil
}

func flattenDisk(pr []config.Disk) []interface{} {
	if len(pr) > 0 {
		diskList := make([]interface{}, len(pr))

		for k, v := range pr {
			disk := make(map[string]interface{})
			if v.TenantId != nil {
				disk["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				disk["links"] = flattenAPILink(v.Links)
			}
			if v.ExtId != nil {
				disk["ext_id"] = v.ExtId
			}
			if v.DiskAddress != nil {
				disk["disk_address"] = flattenDiskAddress(v.DiskAddress)
			}
			if v.BackingInfo != nil {
				disk["backing_info"] = flattenOneOfDiskBackingInfo(v.BackingInfo)
			}

			log.Printf("[DEBUG] disk: %v", disk)
			diskList[k] = disk
		}
		log.Printf("[DEBUG] diskList: %v", diskList)
		return diskList
	}
	return nil
}

func flattenOneOfDiskBackingInfo(pr *config.OneOfDiskBackingInfo) []map[string]interface{} {
	if pr != nil {
		backingInfoList := make([]map[string]interface{}, 0)
		vmDiskInfo := make(map[string]interface{})
		// vmDiskInfoList := make([]map[string]interface{}, 0)
		volumeGroupInfo := make(map[string]interface{})
		volumeGroupInfoList := make([]map[string]interface{}, 0)

		backingInfoObj := make(map[string]interface{})

		if *pr.ObjectType_ == "vmm.v4.ahv.config.VmDisk" {
			vmDiskObj := make(map[string]interface{})
			vmDiskObjList := make([]map[string]interface{}, 0)
			vmDiskVal := pr.GetValue().(config.VmDisk)

			if vmDiskVal.DiskSizeBytes != nil {
				vmDiskObj["disk_size_bytes"] = vmDiskVal.DiskSizeBytes
			}
			if vmDiskVal.StorageContainer != nil {
				vmDiskObj["storage_container"] = flattenVMDiskContainerReference(vmDiskVal.StorageContainer)
			}
			if vmDiskVal.StorageConfig != nil {
				vmDiskObj["storage_config"] = flattenVMDiskStorageConfig(vmDiskVal.StorageConfig)
			}
			if vmDiskVal.DataSource != nil {
				vmDiskObj["data_source"] = flattenDataSource(vmDiskVal.DataSource)
			}
			if vmDiskVal.DiskExtId != nil {
				vmDiskObj["disk_ext_id"] = vmDiskVal.DiskExtId
			}
			if vmDiskVal.IsMigrationInProgress != nil {
				vmDiskObj["is_migration_in_progress"] = vmDiskVal.IsMigrationInProgress
			}
			log.Printf("[DEBUG] vmDiskObj: %v", vmDiskObj)

			vmDiskObjList = append(vmDiskObjList, vmDiskObj)
			backingInfoObj["vm_disk"] = vmDiskObjList
			log.Printf("[DEBUG] backingInfoObj: %v", vmDiskInfo)
			backingInfoList = append(backingInfoList, backingInfoObj)
			return backingInfoList
		}
		volumeGroupObj := make(map[string]interface{})
		volumeGroupObjList := make([]map[string]interface{}, 0)
		volumeGroupVal := pr.GetValue().(config.ADSFVolumeGroupReference)

		if volumeGroupVal.VolumeGroupExtId != nil {
			volumeGroupObj["volume_group_ext_id"] = volumeGroupVal.VolumeGroupExtId
		}

		volumeGroupObjList = append(volumeGroupObjList, volumeGroupObj)

		volumeGroupInfo["adfs_volume_group_reference"] = volumeGroupObjList
		volumeGroupInfoList = append(volumeGroupInfoList, volumeGroupInfo)

		backingInfoList = volumeGroupInfoList
		return backingInfoList
	}
	return nil
}

func flattenCdRom(pr []config.CdRom) []interface{} {
	if len(pr) > 0 {
		cdRomList := make([]interface{}, len(pr))

		for k, v := range pr {
			cd := make(map[string]interface{})

			if v.TenantId != nil {
				cd["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				cd["links"] = flattenAPILink(v.Links)
			}
			if v.ExtId != nil {
				cd["ext_id"] = v.ExtId
			}
			if v.DiskAddress != nil {
				cd["disk_address"] = flattenCdRomAddress(v.DiskAddress)
			}
			if v.BackingInfo != nil {
				cd["backing_info"] = flattenVMDisk(v.BackingInfo)
			}
			if v.IsoType != nil {
				cd["iso_type"] = flattenIsoType(v.IsoType)
			}
			cdRomList[k] = cd
		}
		return cdRomList
	}
	return nil
}

func flattenIsoType(pr *config.IsoType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == config.IsoType(two) {
			return "OTHER"
		}
		if *pr == config.IsoType(three) {
			return "GUEST_TOOLS"
		}
		if *pr == config.IsoType(four) {
			return "GUEST_CUSTOMIZATION"
		}
	}
	return "UNKNOWN"
}

func flattenCdRomAddress(pr *config.CdRomAddress) []map[string]interface{} {
	if pr != nil {
		cdromList := make([]map[string]interface{}, 0)
		cdrom := make(map[string]interface{})

		if pr.BusType != nil {
			cdrom["bus_type"] = flattenCdRomBusType(pr.BusType)
		}
		if pr.Index != nil {
			cdrom["index"] = pr.Index
		}
		cdromList = append(cdromList, cdrom)
		return cdromList
	}
	return nil
}

func flattenCdRomBusType(pr *config.CdRomBusType) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == config.CdRomBusType(two) {
			return "IDE"
		}
		if *pr == config.CdRomBusType(three) {
			return "SATA"
		}
	}
	return "UNKNOWN"
}

func flattenNic(nic []config.Nic) []interface{} {
	if len(nic) > 0 {
		nicList := make([]interface{}, len(nic))

		for k, v := range nic {
			nics := make(map[string]interface{})
			if v.TenantId != nil {
				nics["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				nics["links"] = flattenAPILink(v.Links)
			}
			if v.ExtId != nil {
				nics["ext_id"] = v.ExtId
			}
			if v.BackingInfo != nil {
				nics["backing_info"] = flattenEmulatedNic(v.BackingInfo)
			}
			if v.NetworkInfo != nil {
				nics["network_info"] = flattenNicNetworkInfo(v.NetworkInfo)
			}
			nicList[k] = nics
		}
		return nicList
	}
	return nil
}

func flattenEmulatedNic(pr *config.EmulatedNic) []map[string]interface{} {
	if pr != nil {
		nicList := make([]map[string]interface{}, 0)
		nic := make(map[string]interface{})

		if pr.Model != nil {
			nic["model"] = flattenEmulatedNicModel(pr.Model)
		}
		if pr.MacAddress != nil {
			nic["mac_address"] = pr.MacAddress
		}
		if pr.IsConnected != nil {
			nic["is_connected"] = pr.IsConnected
		}
		if pr.NumQueues != nil {
			nic["num_queues"] = pr.NumQueues
		}

		nicList = append(nicList, nic)
		return nicList
	}
	return nil
}

func flattenEmulatedNicModel(pr *config.EmulatedNicModel) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == config.EmulatedNicModel(two) {
			return "VIRTIO"
		}
		if *pr == config.EmulatedNicModel(three) {
			return "E1000"
		}
	}
	return "UNKNOWN"
}

func flattenNicNetworkInfo(pr *config.NicNetworkInfo) []map[string]interface{} {
	if pr != nil {
		nicList := make([]map[string]interface{}, 0)
		nic := make(map[string]interface{})

		if pr.NicType != nil {
			nic["nic_type"] = flattenNicType(pr.NicType)
		}
		if pr.NetworkFunctionChain != nil {
			nic["network_function_chain"] = flattenNetworkFunctionChainReference(pr.NetworkFunctionChain)
		}
		if pr.NetworkFunctionNicType != nil {
			nic["network_function_nic_type"] = flattenNetworkFunctionNicType(pr.NetworkFunctionNicType)
		}
		if pr.Subnet != nil {
			nic["subnet"] = flattenSubnetReference(pr.Subnet)
		}
		if pr.VlanMode != nil {
			nic["vlan_mode"] = flattenVlanMode(pr.VlanMode)
		}
		if pr.TrunkedVlans != nil {
			nic["trunked_vlans"] = pr.TrunkedVlans
		}
		if pr.ShouldAllowUnknownMacs != nil {
			nic["should_allow_unknown_macs"] = pr.ShouldAllowUnknownMacs
		}
		if pr.Ipv4Config != nil {
			nic["ipv4_config"] = flattenIpv4Config(pr.Ipv4Config)
		}
		if pr.Ipv4Info != nil {
			nic["ipv4_info"] = flattenIpv4Info(pr.Ipv4Info)
		}

		nicList = append(nicList, nic)
		return nicList
	}
	return nil
}

func flattenNicType(pr *config.NicType) string {
	if pr != nil {
		const two, three, four, five = 2, 3, 4, 5
		if *pr == config.NicType(two) {
			return "NORMAL_NIC"
		}
		if *pr == config.NicType(three) {
			return "DIRECT_NIC"
		}
		if *pr == config.NicType(four) {
			return "NETWORK_FUNCTION_NIC"
		}
		if *pr == config.NicType(five) {
			return "SPAN_DESTINATION_NIC"
		}
	}
	return "UNKNOWN"
}

func flattenNetworkFunctionChainReference(ref *config.NetworkFunctionChainReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenNetworkFunctionNicType(pr *config.NetworkFunctionNicType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == config.NetworkFunctionNicType(two) {
			return "INGRESS"
		}
		if *pr == config.NetworkFunctionNicType(three) {
			return "EGRESS"
		}
		if *pr == config.NetworkFunctionNicType(four) {
			return "TAP"
		}
	}
	return "UNKNOWN"
}

func flattenSubnetReference(ref *config.SubnetReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenVlanMode(pr *config.VlanMode) string {
	if pr != nil {
		const two, three = 2, 3
		if *pr == config.VlanMode(two) {
			return "ACCESS"
		}
		if *pr == config.VlanMode(three) {
			return "TRUNK"
		}
	}
	return "UNKNOWN"
}

func flattenIpv4Config(pr *config.Ipv4Config) []map[string]interface{} {
	if pr != nil {
		cfgList := make([]map[string]interface{}, 0)
		cfg := make(map[string]interface{})

		if pr.ShouldAssignIp != nil {
			cfg["should_assign_ip"] = pr.ShouldAssignIp
		}
		if pr.IpAddress != nil {
			cfg["ip_address"] = flattenIPv4Address(pr.IpAddress)
		}
		if pr.SecondaryIpAddressList != nil {
			cfg["secondary_ip_address_list"] = flattenIPv4AddressList(pr.SecondaryIpAddressList)
		}

		cfgList = append(cfgList, cfg)
		return cfgList
	}
	return nil
}

func flattenIpv4Info(ipv4Info *config.Ipv4Info) []map[string]interface{} {
	if ipv4Info != nil {
		ipv4List := make([]map[string]interface{}, 0)
		ipv4 := make(map[string]interface{})

		if ipv4Info.LearnedIpAddresses != nil {
			ipv4["learned_ip_addresses"] = flattenIPv4AddressList(ipv4Info.LearnedIpAddresses)
		}

		ipv4List = append(ipv4List, ipv4)
		return ipv4List
	}
	return nil
}

func flattenIPv4Address(pr *import4.IPv4Address) []map[string]interface{} {
	if pr != nil {
		ipv4List := make([]map[string]interface{}, 0)
		ipv4s := make(map[string]interface{})

		if pr.PrefixLength != nil {
			ipv4s["prefix_length"] = pr.PrefixLength
		}
		if pr.Value != nil {
			ipv4s["value"] = pr.Value
		}

		ipv4List = append(ipv4List, ipv4s)
		return ipv4List
	}
	return nil
}

func flattenIPv4AddressList(pr []import4.IPv4Address) []interface{} {
	if len(pr) > 0 {
		ipv4List := make([]interface{}, len(pr))

		for k, v := range pr {
			ipv4 := make(map[string]interface{})

			if v.PrefixLength != nil {
				ipv4["prefix_length"] = v.PrefixLength
			}
			if v.Value != nil {
				ipv4["value"] = v.Value
			}
			ipv4List[k] = ipv4
		}
		return ipv4List
	}
	return nil
}

func flattenGpu(pr []config.Gpu) []interface{} {
	if len(pr) > 0 {
		gpus := make([]interface{}, len(pr))

		for k, v := range pr {
			gpu := make(map[string]interface{})

			if v.TenantId != nil {
				gpu["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				gpu["links"] = flattenAPILink(v.Links)
			}
			if v.ExtId != nil {
				gpu["ext_id"] = v.ExtId
			}
			if v.Mode != nil {
				gpu["mode"] = flattenGpuMode(v.Mode)
			}
			if v.DeviceId != nil {
				gpu["device_id"] = v.DeviceId
			}
			if v.Vendor != nil {
				gpu["vendor"] = flattenGpuVendor(v.Vendor)
			}
			if v.PciAddress != nil {
				gpu["pci_address"] = flattenSBDF(v.PciAddress)
			}
			if v.GuestDriverVersion != nil {
				gpu["guest_driver_version"] = v.GuestDriverVersion
			}
			if v.Name != nil {
				gpu["name"] = v.Name
			}
			if v.FrameBufferSizeBytes != nil {
				gpu["frame_buffer_size_bytes"] = v.FrameBufferSizeBytes
			}
			if v.NumVirtualDisplayHeads != nil {
				gpu["num_virtual_display_heads"] = v.NumVirtualDisplayHeads
			}
			if v.Fraction != nil {
				gpu["fraction"] = v.Fraction
			}
			gpus[k] = gpu
		}
		return gpus
	}
	return nil
}

func flattenGpuMode(pr *config.GpuMode) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == config.GpuMode(two) {
			return "PASSTHROUGH_GRAPHICS"
		}
		if *pr == config.GpuMode(three) {
			return "PASSTHROUGH_COMPUTE"
		}
		if *pr == config.GpuMode(four) {
			return "VIRTUAL"
		}
	}
	return "UNKNOWN"
}

func flattenGpuVendor(pr *config.GpuVendor) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == config.GpuVendor(two) {
			return "NVIDIA"
		}
		if *pr == config.GpuVendor(three) {
			return "INTEL"
		}
		if *pr == config.GpuVendor(four) {
			return "AMD"
		}
	}
	return "UNKNOWN"
}

func flattenSBDF(pr *config.SBDF) []map[string]interface{} {
	if pr != nil {
		sbdfList := make([]map[string]interface{}, 0)
		sbdf := make(map[string]interface{})

		if pr.Bus != nil {
			sbdf["bus"] = pr.Bus
		}
		if pr.Device != nil {
			sbdf["device"] = pr.Device
		}
		if pr.Segment != nil {
			sbdf["segment"] = pr.Segment
		}
		if pr.Func != nil {
			sbdf["func"] = pr.Func
		}

		sbdfList = append(sbdfList, sbdf)
		return sbdfList
	}
	return nil
}

func flattenSerialPort(pr []config.SerialPort) []interface{} {
	if len(pr) > 0 {
		portList := make([]interface{}, len(pr))

		for k, v := range pr {
			port := make(map[string]interface{})

			if v.TenantId != nil {
				port["tenant_id"] = v.TenantId
			}
			if v.Links != nil {
				port["links"] = flattenAPILink(v.Links)
			}
			if v.ExtId != nil {
				port["ext_id"] = v.ExtId
			}
			if v.IsConnected != nil {
				port["is_connected"] = v.IsConnected
			}
			if v.Index != nil {
				port["index"] = v.Index
			}
			portList[k] = port
		}
		return portList
	}
	return nil
}

func flattenProtectionType(pr *config.ProtectionType) string {
	if pr != nil {
		const two, three, four = 2, 3, 4
		if *pr == config.ProtectionType(two) {
			return "UNPROTECTED"
		}
		if *pr == config.ProtectionType(three) {
			return "PD_PROTECTED"
		}
		if *pr == config.ProtectionType(four) {
			return "RULE_PROTECTED"
		}
	}
	return "UNKNOWN"
}

func flattenProtectionPolicyState(pr *config.ProtectionPolicyState) []map[string]interface{} {
	if pr != nil {
		stateList := make([]map[string]interface{}, 0)
		state := make(map[string]interface{})

		if pr.Policy != nil {
			state["policy"] = flattenPolicyReference(pr.Policy)
		}

		stateList = append(stateList, state)
		return stateList
	}
	return nil
}

func flattenPolicyReference(ref *config.PolicyReference) []map[string]interface{} {
	if ref != nil {
		refList := make([]map[string]interface{}, 0)

		refs := make(map[string]interface{})

		if ref.ExtId != nil {
			refs["ext_id"] = ref.ExtId
		}
		refList = append(refList, refs)

		return refList
	}
	return nil
}

func flattenAPILink(pr []response.ApiLink) []interface{} {
	if len(pr) > 0 {
		links := make([]interface{}, len(pr))

		for k, v := range pr {
			link := make(map[string]interface{})

			if v.Href != nil {
				link["href"] = v.Href
			}
			if v.Rel != nil {
				link["rel"] = v.Rel
			}
			links[k] = link
		}
		return links
	}
	return nil
}
