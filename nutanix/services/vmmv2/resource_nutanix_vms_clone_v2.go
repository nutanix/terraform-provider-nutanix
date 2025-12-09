package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVMCloneV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVMCloneV2Create,
		ReadContext:   ResourceNutanixVMCloneV2Read,
		UpdateContext: ResourceNutanixVMCloneV2Update,
		DeleteContext: ResourceNutanixVMCloneV2Delete,
		Schema: map[string]*schema.Schema{
			"vm_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"num_sockets": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"num_cores_per_socket": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"num_threads_per_core": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"memory_size_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
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
			"boot_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"legacy_boot": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"boot_device": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"boot_device_disk": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_address": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bus_type": {
																			Type:         schema.TypeString,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR", "PCI", "IDE", "SATA"}, false),
																		},
																		"index": {
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
												"boot_device_nic": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
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
									"boot_order": {
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
						"uefi_boot": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_secure_boot_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"nvram_device": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"backing_storage_info": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_ext_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"disk_size_bytes": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
															"storage_container": {
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
															"storage_config": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_flash_mode_enabled": {
																			Type:     schema.TypeBool,
																			Optional: true,
																			Computed: true,
																		},
																	},
																},
															},
															"data_source": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"reference": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Computed: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_reference": {
																						Type:     schema.TypeList,
																						Optional: true,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_ext_id": {
																									Type:     schema.TypeString,
																									Optional: true,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"vm_disk_reference": {
																						Type:     schema.TypeList,
																						Optional: true,
																						Computed: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"disk_ext_id": {
																									Type:     schema.TypeString,
																									Optional: true,
																									Computed: true,
																								},
																								"disk_address": {
																									Type:     schema.TypeList,
																									Optional: true,
																									Computed: true,
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bus_type": {
																												Type:     schema.TypeString,
																												Optional: true,
																												Computed: true,
																												ValidateFunc: validation.StringInSlice([]string{
																													"SCSI", "SPAPR", "PCI",
																													"IDE", "SATA",
																												}, false),
																											},
																											"index": {
																												Type:     schema.TypeInt,
																												Optional: true,
																												Computed: true,
																											},
																										},
																									},
																								},
																								"vm_reference": {
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
									"boot_device": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"boot_device_disk": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_address": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bus_type": {
																			Type:         schema.TypeString,
																			Optional:     true,
																			Computed:     true,
																			ValidateFunc: validation.StringInSlice([]string{"SCSI", "SPAPR", "PCI", "IDE", "SATA"}, false),
																		},
																		"index": {
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
												"boot_device_nic": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
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
									"boot_order": {
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

			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			"num_numa_nodes": {
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

func ResourceNutanixVMCloneV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	vmExtID := d.Get("vm_ext_id")

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(vmExtID.(string)))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	body := &config.CloneOverrideParams{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if numSock, ok := d.GetOk("num_sockets"); ok {
		body.NumSockets = utils.IntPtr(numSock.(int))
	}
	if numCoresPerSock, ok := d.GetOk("num_cores_per_socket"); ok {
		body.NumCoresPerSocket = utils.IntPtr(numCoresPerSock.(int))
	}
	if numThreadsPerCore, ok := d.GetOk("num_threads_per_core"); ok {
		body.NumThreadsPerCore = utils.IntPtr(numThreadsPerCore.(int))
	}
	if memorySize, ok := d.GetOk("memory_size_bytes"); ok {
		body.MemorySizeBytes = utils.Int64Ptr(int64(memorySize.(int)))
	}
	if guestCstm, ok := d.GetOk("guest_customization"); ok {
		body.GuestCustomization = expandGuestCustomizationParams(guestCstm)
	}
	if bootConfig, ok := d.GetOk("boot_config"); ok {
		body.BootConfig = expandOneOfCloneVMBootConfig(bootConfig)
	}

	resp, err := conn.VMAPIInstance.CloneVm(utils.StringPtr(vmExtID.(string)), body, args)
	if err != nil {
		return diag.Errorf("error while Cloning Vm : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be cloned
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM clone (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching VM clone task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", " ")
	log.Printf("[DEBUG] Clone VM Task Details: %s", string(aJSON))

	// The cloned VM is the second entity (index 1) in EntitiesAffected
	uuid := taskDetails.EntitiesAffected[1].ExtId
	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixVMCloneV2Read(ctx, d, meta)
}

func ResourceNutanixVMCloneV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
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

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
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
	if err := d.Set("memory_size_bytes", getResp.MemorySizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_customization", flattenGuestCustomizationParams(getResp.GuestCustomization)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("boot_config", flattenOneOfVMBootConfig(getResp.BootConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenAPILink(getResp.Links)); err != nil {
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
	if err := d.Set("num_numa_nodes", getResp.NumNumaNodes); err != nil {
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

	log.Printf("[INFO] Successfully read cloned vm with ext_id %s", d.Id())
	d.SetId(utils.StringValue(getResp.ExtId))

	return nil
}

func ResourceNutanixVMCloneV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixVMCloneV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixVirtualMachineV2Delete(ctx, d, meta)
}

func expandOneOfCloneVMBootConfig(pr interface{}) *config.OneOfCloneOverrideParamsBootConfig {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		vmBootConfig := &config.OneOfCloneOverrideParamsBootConfig{}

		if legacy, ok := val["legacy_boot"]; ok && len(legacy.([]interface{})) > 0 {
			legacyBootInput := config.NewLegacyBoot()

			prI := legacy.([]interface{})
			val := prI[0].(map[string]interface{})

			if bootDevice, ok := val["boot_device"]; ok && len(bootDevice.([]interface{})) > 0 {
				legacyBootInput.BootDevice = expandOneOfLegacyBootBootDevice(bootDevice)
			}
			if bootOrder, ok := val["boot_order"]; ok && len(bootOrder.([]interface{})) > 0 {
				bootOrderLen := len(bootOrder.([]interface{}))
				orders := make([]config.BootDeviceType, bootOrderLen)

				for k, v := range bootOrder.([]interface{}) {
					const two, three, four = 2, 3, 4
					subMap := map[string]interface{}{
						"CDROM":   two,
						"DISK":    three,
						"NETWORK": four,
					}
					pVal := subMap[v.(string)]
					p := config.BootDeviceType(pVal.(int))
					orders[k] = p
				}
				legacyBootInput.BootOrder = orders
			}
			vmBootConfig.SetValue(*legacyBootInput)
		}
		if uefi, ok := val["uefi_boot"]; ok && len(uefi.([]interface{})) > 0 {
			uefiBootInput := config.NewUefiBoot()

			prI := uefi.([]interface{})
			val := prI[0].(map[string]interface{})

			if secureBootEnabled, ok := val["is_secure_boot_enabled"]; ok {
				uefiBootInput.IsSecureBootEnabled = utils.BoolPtr(secureBootEnabled.(bool))
			}
			if nvram, ok := val["nvram_device"]; ok && len(nvram.([]interface{})) > 0 {
				uefiBootInput.NvramDevice = expandNvramDevice(nvram)
			}
			if bootDevice, ok := val["boot_device"]; ok && len(bootDevice.([]interface{})) > 0 {
				uefiBootInput.BootDevice = expandOneOfUefiBootBootDevice(bootDevice)
			}
			if bootOrder, ok := val["boot_order"]; ok && len(bootOrder.([]interface{})) > 0 {
				bootOrderLen := len(bootOrder.([]interface{}))
				orders := make([]config.BootDeviceType, bootOrderLen)

				for k, v := range bootOrder.([]interface{}) {
					const two, three, four = 2, 3, 4
					subMap := map[string]interface{}{
						"CDROM":   two,
						"DISK":    three,
						"NETWORK": four,
					}
					pVal := subMap[v.(string)]
					p := config.BootDeviceType(pVal.(int))
					orders[k] = p
				}
				uefiBootInput.BootOrder = orders
			}
			vmBootConfig.SetValue(*uefiBootInput)
		}
		return vmBootConfig
	}
	return nil
}
