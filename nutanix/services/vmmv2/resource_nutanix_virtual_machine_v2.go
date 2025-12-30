package vmmv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	timeout = 3 * time.Minute
	delay   = 3 * time.Second
)

func ResourceNutanixVirtualMachineV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVirtualMachineV2Create,
		ReadContext:   ResourceNutanixVirtualMachineV2Read,
		UpdateContext: ResourceNutanixVirtualMachineV2Update,
		DeleteContext: ResourceNutanixVirtualMachineV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"VM", "VM_RECOVERY_POINT"}, false),
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			"num_numa_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"memory_size_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"is_vcpu_hard_pinning_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_cpu_passthrough_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enabled_cpu_features": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_memory_overcommit_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_gpu_console_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_cpu_hotplug_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_scsi_controller_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"generation_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bios_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"categories": {
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
			"project": {
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
			"ownership_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"owner": {
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
			"host": {
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
			"cluster": {
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
			// not present in API reference
			"availability_zone": {
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
			"guest_customization": schemaForGuestCustomization(),
			"guest_tools": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"capabilities": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"SELF_SERVICE_RESTORE", "VSS_SNAPSHOT"}, false),
							},
						},
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
					},
				},
			},
			"hardware_clock_timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_branding_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"CDROM", "DISK", "NETWORK"}, false),
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
				Optional: true,
				Computed: true,
			},
			"machine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PSERIES", "Q35", "PC"}, false),
			},
			"power_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ON",
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF", "PAUSED", "UNDETERMINED"}, false),
			},
			"vtpm_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_vtpm_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
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
				Optional: true,
				Computed: true,
			},
			"apc_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_apc_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"cpu_model": {
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
									"name": {
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
						"qos_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"throttled_iops": {
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
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
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
						"backing_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_disk": {
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
									"adfs_volume_group_reference": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"volume_group_ext_id": {
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
			"cd_roms": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
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
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"IDE", "SATA"}, false),
									},
									"index": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"backing_info": {
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
						"iso_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION"}, false),
						},
					},
				},
			},
			"nics": {
				Type:     schema.TypeList,
				Optional: true,
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
					},
				},
			},
			"gpus": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": schemaForLinks(),
						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"PASSTHROUGH_GRAPHICS", "PASSTHROUGH_COMPUTE", "VIRTUAL"}, false),
						},
						"device_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"vendor": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"NVIDIA", "AMD", "INTEL"}, false),
						},
						// not present in api reference doc
						"pci_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segment": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"bus": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"device": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"func": {
										Type:     schema.TypeInt,
										Optional: true,
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_connected": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"protection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PD_PROTECTED", "UNPROTECTED", "RULE_PROTECTED"}, false),
			},
			"protection_policy_state": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaForGuestCustomization() *schema.Schema {
	return &schema.Schema{
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
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											ValidateFunc: validation.StringInSlice([]string{"PREPARED", "FRESH"}, false),
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
																// this value is required but not present in API reference
																// the create vm request fails if this is not provided or is empty
																"value": {
																	Type:     schema.TypeString,
																	Required: true,
																},
															},
														},
													},
													"custom_key_values": schemaForCustomKeyValuePairs(),
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
													"custom_key_values": schemaForCustomKeyValuePairs(),
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

func ResourceNutanixVirtualMachineV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	VMConfigMap := resourceDataToMap(d, ResourceNutanixVirtualMachineV2().Schema)
	body := prepareVMConfigFromMap(VMConfigMap)
	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Vm Create Request Payload: %s", string(aJSON))

	resp, err := conn.VMAPIInstance.CreateVm(body)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		return diag.Errorf("error while creating Virtual Machines : %v", e)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for virtual Machine (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		return diag.Errorf("error while fetching vm UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	// read VM

	readResp, errR := conn.VMAPIInstance.GetVmById(utils.StringPtr(*uuid))
	if errR != nil {
		return diag.Errorf("error while reading vm : %v", errR)
	}
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	var PowerTaskRef import1.TaskReference
	if powerState, ok := d.GetOk("power_state"); ok {
		if powerState == "ON" {
			resp, err := conn.VMAPIInstance.PowerOnVm(uuid, args)
			if err != nil {
				return diag.Errorf("error while powering on Virtual Machines : %v", err)
			}
			PowerTaskRef = resp.Data.GetValue().(import1.TaskReference)
		} else if powerState == "OFF" {
			resp, err := conn.VMAPIInstance.PowerOffVm(utils.StringPtr(d.Id()), args)
			if err != nil {
				return diag.Errorf("error while powering OFF : %v", err)
			}
			PowerTaskRef = resp.Data.GetValue().(import1.TaskReference)
		}
	}
	powertaskUUID := PowerTaskRef.ExtId

	// Wait for the VM to be available
	powerStateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(powertaskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := powerStateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vm (%s) to power ON: %s", utils.StringValue(uuid), errWaitTask)
	}

	// If power state is ON and NICs are configured, wait for IP address
	// Skip waiting if no NICs are configured since there won't be any IP address
	nics := d.Get("nics")
	hasNics := nics != nil && len(common.InterfaceToSlice(nics)) > 0

	if d.Get("power_state") == "ON" && hasNics {
		// Wait for the VM to be available
		waitIPConf := &resource.StateChangeConf{
			Pending:    []string{"WAITING"},
			Target:     []string{"AVAILABLE"},
			Refresh:    waitForIPRefreshFunc(conn, utils.StringValue(uuid)),
			Timeout:    timeout,
			Delay:      delay,
			MinTimeout: delay,
		}
		vmIntentResponse, err := waitIPConf.WaitForStateContext(ctx)
		if err != nil {
			log.Printf("[WARN] could not get the IP for VM(%s): %s", utils.StringValue(uuid), err)
		} else {
			vm := vmIntentResponse.(*config.GetVmApiResponse)
			vmResp := vm.Data.GetValue().(config.Vm)

			if len(vmResp.Nics) > 0 && vmResp.Nics[0].NetworkInfo != nil {
				ipAddr := getFirstIPAddress(vmResp.Nics[0])
				if ipAddr != "" {
					d.SetConnInfo(map[string]string{
						"type": "ssh",
						"host": ipAddr,
					})
				}
			}
		}
	}

	return ResourceNutanixVirtualMachineV2Read(ctx, d, meta)
}

func ResourceNutanixVirtualMachineV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching vm : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Vm)
	setVMConfig(d, getResp)

	return nil
}

func ResourceNutanixVirtualMachineV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI
	// respImages := resp.Data.GetValue().(config.Vm)
	// updateSpec := respImages

	if checkForHotPlugChanges(d) && !isVMPowerOff(d, conn) {
		log.Printf("[DEBUG] callingForPowerOffVM func")
		callForPowerOffVM(ctx, conn, d, meta)
	}

	updatedVMResp, _ := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))

	respVM := updatedVMResp.Data.GetValue().(config.Vm)

	updateSpec := respVM

	checkForUpdateParams := false

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
		checkForUpdateParams = true
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
		checkForUpdateParams = true
	}
	if d.HasChange("source") {
		updateSpec.Source = expandVMSourceReference(d.Get("source"))
		checkForUpdateParams = true
	}
	if d.HasChange("num_sockets") {
		updateSpec.NumSockets = utils.IntPtr(d.Get("num_sockets").(int))
		checkForUpdateParams = true
	}
	if d.HasChange("num_cores_per_socket") {
		updateSpec.NumCoresPerSocket = utils.IntPtr(d.Get("num_cores_per_socket").(int))
		checkForUpdateParams = true
	}
	if d.HasChange("num_threads_per_core") {
		updateSpec.NumThreadsPerCore = utils.IntPtr(d.Get("num_threads_per_core").(int))
		checkForUpdateParams = true
	}
	if d.HasChange("num_numa_nodes") {
		updateSpec.NumNumaNodes = utils.IntPtr(d.Get("num_numa_nodes").(int))
		checkForUpdateParams = true
	}
	if d.HasChange("memory_size_bytes") {
		updateSpec.MemorySizeBytes = utils.Int64Ptr(int64(d.Get("memory_size_bytes").(int)))
		checkForUpdateParams = true
	}
	if d.HasChange("is_vcpu_hard_pinning_enabled") {
		updateSpec.IsVcpuHardPinningEnabled = utils.BoolPtr(d.Get("is_vcpu_hard_pinning_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("is_cpu_passthrough_enabled") {
		updateSpec.IsCpuPassthroughEnabled = utils.BoolPtr(d.Get("is_cpu_passthrough_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("enabled_cpu_features") {
		updateSpec.EnabledCpuFeatures = expandCPUFeature(d.Get("enabled_cpu_features").([]interface{}))
		checkForUpdateParams = true
	}
	if d.HasChange("is_memory_overcommit_enabled") {
		updateSpec.IsMemoryOvercommitEnabled = utils.BoolPtr(d.Get("is_memory_overcommit_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("is_gpu_console_enabled") {
		updateSpec.IsGpuConsoleEnabled = utils.BoolPtr(d.Get("is_gpu_console_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("is_cpu_hotplug_enabled") {
		updateSpec.IsCpuHotplugEnabled = utils.BoolPtr(d.Get("is_cpu_hotplug_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("is_scsi_controller_enabled") {
		updateSpec.IsScsiControllerEnabled = utils.BoolPtr(d.Get("is_scsi_controller_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("generation_uuid") {
		updateSpec.GenerationUuid = utils.StringPtr(d.Get("generation_uuid").(string))
		checkForUpdateParams = true
	}
	if d.HasChange("bios_uuid") {
		updateSpec.BiosUuid = utils.StringPtr(d.Get("bios_uuid").(string))
		checkForUpdateParams = true
	}
	if d.HasChange("ownership_info") {
		updateSpec.OwnershipInfo = expandOwnershipInfo(d.Get("ownership_info"))
		checkForUpdateParams = true
	}
	if d.HasChange("host") {
		updateSpec.Host = expandHostReference(d.Get("host"))
		checkForUpdateParams = true
	}
	if d.HasChange("cluster") {
		updateSpec.Cluster = expandClusterReference(d.Get("cluster"))
		checkForUpdateParams = true
	}
	if d.HasChange("guest_customization") {
		updateSpec.GuestCustomization = expandTemplateGuestCustomizationParams(d.Get("guest_customization"))
		checkForUpdateParams = true
	}
	if d.HasChange("hardware_clock_timezone") {
		updateSpec.HardwareClockTimezone = utils.StringPtr(d.Get("hardware_clock_timezone").(string))
		checkForUpdateParams = true
	}
	if d.HasChange("is_branding_enabled") {
		updateSpec.IsBrandingEnabled = utils.BoolPtr(d.Get("is_branding_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("boot_config") {
		updateSpec.BootConfig = expandOneOfVMBootConfig(d.Get("boot_config"))
		checkForUpdateParams = true
	}
	if d.HasChange("is_vga_console_enabled") {
		updateSpec.IsVgaConsoleEnabled = utils.BoolPtr(d.Get("is_vga_console_enabled").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("machine_type") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"PC":      two,
			"PSERIES": three,
			"Q35":     four,
		}
		pVal := subMap[d.Get("machine_type").(string)]
		p := config.MachineType(pVal.(int))
		updateSpec.MachineType = &p
		checkForUpdateParams = true
	}
	if d.HasChange("vtpm_config") {
		updateSpec.VtpmConfig = expandVtpmConfig(d.Get("vtpm_config"))
		checkForUpdateParams = true
	}
	if d.HasChange("is_agent_vm") {
		updateSpec.IsAgentVm = utils.BoolPtr(d.Get("is_agent_vm").(bool))
		checkForUpdateParams = true
	}
	if d.HasChange("apc_config") {
		updateSpec.ApcConfig = expandApcConfig(d.Get("apc_config"))
		checkForUpdateParams = true
	}
	if d.HasChange("storage_config") {
		updateSpec.StorageConfig = expandADSFVmStorageConfig(d.Get("storage_config"))
		checkForUpdateParams = true
	}
	if d.HasChange("protection_type") {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"UNPROTECTED":    two,
			"PD_PROTECTED":   three,
			"RULE_PROTECTED": four,
		}
		pVal := subMap[d.Get("protection_type").(string)]
		p := config.ProtectionType(pVal.(int))
		updateSpec.ProtectionType = &p
		checkForUpdateParams = true
	}
	if d.HasChange("protection_policy_state") {
		updateSpec.ProtectionPolicyState = expandProtectionPolicyState(d.Get("protection_policy_state"))
		checkForUpdateParams = true
	}
	if d.HasChange("guest_tools") {
		updateSpec.GuestTools = expandGuestTools(d.Get("guest_tools"))
		checkForUpdateParams = true
	}

	if checkForUpdateParams {
		updateResp, err := conn.VMAPIInstance.UpdateVmById(utils.StringPtr(d.Id()), &updateSpec)
		if err != nil {
			return diag.Errorf("error while updating Virtual Machines : %v", err)
		}

		TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
		taskUUID := TaskRef.ExtId

		taskconn := meta.(*conns.Client).PrismAPI
		// Wait for the VM to be available
		stateConf := &resource.StateChangeConf{
			Pending: []string{"QUEUED", "RUNNING"},
			Target:  []string{"SUCCEEDED"},
			Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for virtual machine (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
		}
	}

	// now calling different APIs to update the VM

	if d.HasChange("disks") {
		oldDisk, newDisk := d.GetChange("disks")
		newAddedDisk, oldDeletedDisk, updatedDisk := diffConfig(oldDisk.([]interface{}), newDisk.([]interface{}))

		if len(oldDeletedDisk) > 0 {
			for _, disk := range oldDeletedDisk {

				diskInput := expandDisk([]interface{}{disk})[0]

				diskExtID := diskInput.ExtId

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.DeleteDiskById(utils.StringPtr(d.Id()), diskExtID, args)
				if err != nil {
					return diag.Errorf("error while deleting Disk : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for disk (%s) to be deleted: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}

		if len(updatedDisk) > 0 {
			for _, disk := range updatedDisk {
				if diskMap, ok := disk.(map[string]interface{}); ok {
					if backingInfoRaw, ok := diskMap["backing_info"]; ok {
						if backingInfoSlice, ok := backingInfoRaw.([]interface{}); ok {
							if backingInfoMap, ok := backingInfoSlice[0].(map[string]interface{}); ok {
								if vmDiskArray, ok := backingInfoMap["vm_disk"].([]interface{}); ok {
									if vmDiskMap, ok := vmDiskArray[0].(map[string]interface{}); ok {
										if vmDiskMap["data_source"] != nil {
											delete(vmDiskMap, "data_source")
										}
									}
								}
							}
						}
					}
				}

				diskInput := expandDisk([]interface{}{disk})[0]

				diskExtID := diskInput.ExtId

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.UpdateDiskById(utils.StringPtr(d.Id()), diskExtID, &diskInput, args)
				if err != nil {
					return diag.Errorf("error while updating Disk : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for disk (%s) to be updated: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}

		if len(newAddedDisk) > 0 {
			for _, disk := range newAddedDisk {
				diskInput := expandDisk([]interface{}{disk})[0]

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.CreateDisk(utils.StringPtr(d.Id()), &diskInput, args)
				if err != nil {
					return diag.Errorf("error while creating Disk : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for disk (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
	}

	if d.HasChange("nics") {
		oldNic, newNic := d.GetChange("nics")
		aJSON, _ := json.MarshalIndent(common.InterfaceToSlice(oldNic), "", "  ")
		log.Println("[DEBUG] oldNic raw config:", string(aJSON))
		aJSON, _ = json.MarshalIndent(common.InterfaceToSlice(newNic), "", "  ")
		log.Println("[DEBUG] newNic raw config:", string(aJSON))
		newAddedNic, oldDeletedNic, updatedNic := diffConfig(common.InterfaceToSlice(oldNic), common.InterfaceToSlice(newNic))

		aJSON, _ = json.MarshalIndent(newAddedNic, "", "  ")
		log.Println("[DEBUG] newAddedNic diff config:", string(aJSON))
		aJSON, _ = json.MarshalIndent(oldDeletedNic, "", "  ")
		log.Println("[DEBUG] oldDeletedNic diff config:", string(aJSON))
		aJSON, _ = json.MarshalIndent(updatedNic, "", "  ")
		log.Println("[DEBUG] updatedNic diff config:", string(aJSON))

		if len(oldDeletedNic) > 0 {
			for _, nic := range oldDeletedNic {
				nicInput := expandNic([]interface{}{nic})[0]

				nicExtID := nicInput.ExtId

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.DeleteNicById(utils.StringPtr(d.Id()), nicExtID, args)
				if err != nil {
					return diag.Errorf("error while deleting nic : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for nic (%s) to be deleted: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
		if len(updatedNic) > 0 {
			for _, nic := range updatedNic {
				nicInput := expandNic([]interface{}{nic})[0]

				nicExtID := nicInput.ExtId

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				log.Printf("[DEBUG] updating nic: %+v", *nicExtID)
				aJSON, _ := json.MarshalIndent(nicInput, "", "  ")
				log.Printf("[DEBUG] update nic payload: %s", string(aJSON))

				resp, err := conn.VMAPIInstance.UpdateNicById(utils.StringPtr(d.Id()), nicExtID, &nicInput, args)
				if err != nil {
					return diag.Errorf("error while updating Nic : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for nic (%s) to be updated: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
		if len(newAddedNic) > 0 {
			for _, nic := range newAddedNic {
				nicInput := expandNic([]interface{}{nic})[0]

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.CreateNic(utils.StringPtr(d.Id()), &nicInput, args)
				if err != nil {
					return diag.Errorf("error while creating NIC : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for NIC (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
	}

	if d.HasChange("cd_roms") {
		oldCdRom, newCdRom := d.GetChange("cd_roms")
		newAddedCdRom, oldDeletedCdRom, _ := diffConfig(oldCdRom.([]interface{}), newCdRom.([]interface{}))
		if len(newAddedCdRom) > 0 {
			for _, cdrom := range newAddedCdRom {
				cdromInput := expandCdRom([]interface{}{cdrom})[0]

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.CreateCdRom(utils.StringPtr(d.Id()), &cdromInput, args)
				if err != nil {
					return diag.Errorf("error while creating CdRom : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for CdRom (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}

		if len(oldDeletedCdRom) > 0 {
			for _, cdrom := range oldDeletedCdRom {
				cdromInput := expandCdRom([]interface{}{cdrom})[0]

				cdromExtID := cdromInput.ExtId

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.DeleteCdRomById(utils.StringPtr(d.Id()), cdromExtID, args)
				if err != nil {
					return diag.Errorf("error while deleting cdrom : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for cdrom (%s) to be deleted: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
	}

	if d.HasChange("serial_ports") {
		oldSerialPorts, newSerialPorts := d.GetChange("serial_ports")
		newAddedSerialPorts, oldDeletedSerialPorts, updatedSerialPorts := diffConfig(oldSerialPorts.([]interface{}), newSerialPorts.([]interface{}))

		if len(oldDeletedSerialPorts) > 0 {
			for _, serialPort := range oldDeletedSerialPorts {
				serialPortInput := expandSerialPort([]interface{}{serialPort})[0]

				serialPortExtID := serialPortInput.ExtId

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.DeleteSerialPortById(utils.StringPtr(d.Id()), serialPortExtID, args)
				if err != nil {
					return diag.Errorf("error while deleting serial port : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for serial port (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
		if len(updatedSerialPorts) > 0 {
			for _, serialPort := range updatedSerialPorts {
				serialPortInput := expandSerialPort([]interface{}{serialPort})[0]

				portExtTD := serialPortInput.ExtId

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.UpdateSerialPortById(utils.StringPtr(d.Id()), portExtTD, &serialPortInput, args)
				if err != nil {
					return diag.Errorf("error while updating serial port : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for seial port (%s) to be updated: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
		if len(newAddedSerialPorts) > 0 {
			for _, serialPort := range newAddedSerialPorts {
				serialPortInput := expandSerialPort([]interface{}{serialPort})[0]

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.CreateSerialPort(utils.StringPtr(d.Id()), &serialPortInput, args)
				if err != nil {
					return diag.Errorf("error while creating SerialPort : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for SerialPort (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
	}

	if d.HasChange("gpus") {
		oldGpus, newGpus := d.GetChange("gpus")
		newAddedGpus, oldDeletedGpus, _ := diffConfig(oldGpus.([]interface{}), newGpus.([]interface{}))

		if len(newAddedGpus) > 0 {
			for _, gpu := range newAddedGpus {
				gpuInput := expandGpu([]interface{}{gpu})[0]

				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.CreateGpu(utils.StringPtr(d.Id()), &gpuInput, args)
				if err != nil {
					return diag.Errorf("error while creating Gpu : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for Gpu (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}

		if len(oldDeletedGpus) > 0 {
			for _, gpu := range oldDeletedGpus {
				gpuInput := expandGpu([]interface{}{gpu})[0]

				gpuExtID := gpuInput.ExtId
				ReadVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
				if err != nil {
					return diag.Errorf("error while fetching vm : %v", err)
				}

				// // Extract E-Tag Header
				args := make(map[string]interface{})
				args["If-Match"] = getEtagHeader(ReadVMResp, conn)

				resp, err := conn.VMAPIInstance.DeleteGpuById(utils.StringPtr(d.Id()), gpuExtID, args)
				if err != nil {
					return diag.Errorf("error while deleting gpu : %v", err)
				}
				TaskRef := resp.Data.GetValue().(import1.TaskReference)
				taskUUID := TaskRef.ExtId

				taskconn := meta.(*conns.Client).PrismAPI
				// Wait for the VM to be available
				stateConf := &resource.StateChangeConf{
					Pending: []string{"QUEUED", "RUNNING"},
					Target:  []string{"SUCCEEDED"},
					Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for gpu (%s) to be deleted: %s", utils.StringValue(taskUUID), errWaitTask)
				}
			}
		}
	}

	if d.HasChange("categories") {
		oldCategories, newCategories := d.GetChange("categories")
		newAddedCategories, oldDeletedCategories, _ := diffConfig(oldCategories.([]interface{}), newCategories.([]interface{}))

		if len(oldDeletedCategories) > 0 {
			body := config.DisassociateVmCategoriesParams{}

			body.Categories = expandCategoryReference(oldDeletedCategories)

			readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
			if err != nil {
				return diag.Errorf("error while reading vm : %v", err)
			}
			// Extract E-Tag Header
			args := make(map[string]interface{})
			args["If-Match"] = getEtagHeader(readResp, conn)

			resp, err := conn.VMAPIInstance.DisassociateCategories(utils.StringPtr(d.Id()), &body, args)
			if err != nil {
				return diag.Errorf("error while diassociate categories : %v", err)
			}

			TaskRef := resp.Data.GetValue().(import1.TaskReference)
			taskUUID := TaskRef.ExtId

			taskconn := meta.(*conns.Client).PrismAPI
			// Wait for the VM to be available
			stateConf := &resource.StateChangeConf{
				Pending: []string{"QUEUED", "RUNNING"},
				Target:  []string{"SUCCEEDED"},
				Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
				Timeout: d.Timeout(schema.TimeoutCreate),
			}

			if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
				return diag.Errorf("error waiting for categories (%s) to diassociate: %s", utils.StringValue(taskUUID), errWaitTask)
			}
		}

		if len(newAddedCategories) > 0 {
			body := config.AssociateVmCategoriesParams{}

			body.Categories = expandCategoryReference(newAddedCategories)

			readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
			if err != nil {
				return diag.Errorf("error while reading vm : %v", err)
			}
			// Extract E-Tag Header
			args := make(map[string]interface{})
			args["If-Match"] = getEtagHeader(readResp, conn)

			resp, err := conn.VMAPIInstance.AssociateCategories(utils.StringPtr(d.Id()), &body, args)
			if err != nil {
				return diag.Errorf("error while associating categories : %v", err)
			}

			TaskRef := resp.Data.GetValue().(import1.TaskReference)
			taskUUID := TaskRef.ExtId

			taskconn := meta.(*conns.Client).PrismAPI
			// Wait for the VM to be available
			stateConf := &resource.StateChangeConf{
				Pending: []string{"QUEUED", "RUNNING"},
				Target:  []string{"SUCCEEDED"},
				Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
				Timeout: d.Timeout(schema.TimeoutCreate),
			}

			if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
				return diag.Errorf("error waiting for categories (%s) to attach: %s", utils.StringValue(taskUUID), errWaitTask)
			}
		}
	}

	// call for power on VM after updating
	if checkForHotPlugChanges(d) {
		if power, ok := d.GetOk("power_state"); ok {
			if power == "ON" {
				callForPowerOnVM(ctx, conn, d, meta)
			}
		}
	}

	if d.HasChange("power_state") {
		if power, ok := d.GetOk("power_state"); ok {
			log.Printf("[DEBUG] Power state change detected: %s", power)
			if power == "ON" {
				log.Printf("[DEBUG] Powering on the VM")
				callForPowerOnVM(ctx, conn, d, meta)
			} else {
				log.Printf("[DEBUG] Powering off the VM")
				callForPowerOffVM(ctx, conn, d, meta)
			}
		}
	}

	return ResourceNutanixVirtualMachineV2Read(ctx, d, meta)
}

func ResourceNutanixVirtualMachineV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while reading vm : %v", err)
	}
	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	resp, err := conn.VMAPIInstance.DeleteVmById(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting vm : %v", err)
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
		return diag.Errorf("error waiting for vm (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandAvailabilityZoneReference(zone interface{}) *config.AvailabilityZoneReference {
	zoneList, ok := zone.([]interface{})
	if !ok || len(zoneList) == 0 {
		return nil // or handle error accordingly
	}

	zoneRef, ok := zoneList[0].(map[string]interface{})
	if !ok {
		return nil // or handle error accordingly
	}

	return &config.AvailabilityZoneReference{
		ExtId: utils.StringPtr(zoneRef["ext_id"].(string)),
	}
}

func expandVMSourceReference(pr interface{}) *config.VmSourceReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		srcRef := &config.VmSourceReference{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if entity, ok := val["entity_type"]; ok {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"VM":                two,
				"VM_RECOVERY_POINT": three,
			}
			pVal := subMap[entity.(string)]
			p := config.VmSourceReferenceEntityType(pVal.(int))
			srcRef.EntityType = &p
		}
		return srcRef
	}
	return nil
}

func expandCPUFeature(pr []interface{}) []config.CpuFeature {
	if len(pr) > 0 {
		feats := make([]config.CpuFeature, len(pr))
		const two = 2
		for k, v := range pr {
			subMap := map[string]interface{}{
				"HARDWARE_VIRTUALIZATION": two,
			}
			pVal := subMap[v.(string)]
			p := config.CpuFeature(pVal.(int))
			feats[k] = p
		}
		return feats
	}
	return nil
}

func expandCategoryReference(pr []interface{}) []config.CategoryReference {
	if len(pr) > 0 {
		catsRef := make([]config.CategoryReference, len(pr))

		for k, v := range pr {
			cats := config.CategoryReference{}
			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				cats.ExtId = utils.StringPtr(extID.(string))
			}
			catsRef[k] = cats
		}
		return catsRef
	}
	return nil
}

func expandOwnershipInfo(pr interface{}) *config.OwnershipInfo {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		ownerInfo := &config.OwnershipInfo{}

		if owner, ok := val["owner"]; ok {
			ownerInfo.Owner = expandOwnerReference(owner)
		}
		return ownerInfo
	}
	return nil
}

func expandOwnerReference(pr interface{}) *config.OwnerReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		ownerRef := &config.OwnerReference{}

		if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
			ownerRef.ExtId = utils.StringPtr(extID.(string))
		}
		return ownerRef
	}
	return nil
}

func expandHostReference(pr interface{}) *config.HostReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		ownerRef := &config.HostReference{}

		if extID, ok := val["ext_id"]; ok {
			ownerRef.ExtId = utils.StringPtr(extID.(string))
		}
		return ownerRef
	}
	return nil
}

func expandClusterReference(pr interface{}) *config.ClusterReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		ownerRef := &config.ClusterReference{}

		if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
			ownerRef.ExtId = utils.StringPtr(extID.(string))
		}
		return ownerRef
	}
	return nil
}

func expandGuestTools(pr interface{}) *config.GuestTools {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		tools := &config.GuestTools{}

		if isEnabled, ok := val["is_enabled"]; ok {
			tools.IsEnabled = utils.BoolPtr(isEnabled.(bool))
		}
		if capabilities, ok := val["capabilities"]; ok && len(capabilities.([]interface{})) > 0 {
			feats := make([]config.NgtCapability, len(capabilities.([]interface{})))

			for k, v := range capabilities.([]interface{}) {
				const two, three = 2, 3
				subMap := map[string]interface{}{
					"SELF_SERVICE_RESTORE": two,
					"VSS_SNAPSHOT":         three,
				}
				if subMap[v.(string)] == nil {
					return nil
				}
				pVal := subMap[v.(string)]
				p := config.NgtCapability(pVal.(int))
				feats[k] = p
			}
			tools.Capabilities = feats
		}
		return tools
	}
	return nil
}

func expandOneOfVMBootConfig(pr interface{}) *config.OneOfVmBootConfig {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		vmBootConfig := &config.OneOfVmBootConfig{}

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

func expandOneOfUefiBootBootDevice(bootDevice interface{}) *config.OneOfUefiBootBootDevice {
	if bootDevice != nil && len(bootDevice.([]interface{})) > 0 {
		prI := bootDevice.([]interface{})
		val := prI[0].(map[string]interface{})
		BootDevice := &config.OneOfUefiBootBootDevice{}
		if bootDisk, ok := val["boot_device_disk"]; ok && len(bootDisk.([]interface{})) > 0 {
			brI := bootDisk.([]interface{})
			bootVal := brI[0].(map[string]interface{})

			diskOut := config.NewBootDeviceDisk()
			if diskAddress, ok := bootVal["disk_address"]; ok && len(diskAddress.([]interface{})) > 0 {
				daI := diskAddress.([]interface{})
				diskVal := daI[0].(map[string]interface{})
				diskAddOut := config.NewDiskAddress()
				if busType, ok := diskVal["bus_type"]; ok {
					const two, three, four, five, six = 2, 3, 4, 5, 6
					subMap := map[string]interface{}{
						"SCSI":  two,
						"IDE":   three,
						"PCI":   four,
						"SATA":  five,
						"SPAPR": six,
					}
					pVal := subMap[busType.(string)]
					p := config.DiskBusType(pVal.(int))
					diskAddOut.BusType = &p
				}
				if index, ok := diskVal["index"]; ok {
					diskAddOut.Index = utils.IntPtr(index.(int))
				}
				diskOut.DiskAddress = diskAddOut
			}
			BootDevice.SetValue(*diskOut)
		}
		if bootNic, ok := val["boot_device_nic"]; ok && len(bootNic.([]interface{})) > 0 {
			brI := bootNic.([]interface{})
			bootVal := brI[0].(map[string]interface{})

			diskNicOut := config.NewBootDeviceNic()

			if mac, ok := bootVal["mac_address"]; ok {
				diskNicOut.MacAddress = utils.StringPtr(mac.(string))
			}
			BootDevice.SetValue(*diskNicOut)
		}
		return BootDevice
	}
	return nil
}

func expandOneOfLegacyBootBootDevice(pr interface{}) *config.OneOfLegacyBootBootDevice {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		BootDevice := &config.OneOfLegacyBootBootDevice{}
		if bootDisk, ok := val["boot_device_disk"]; ok && len(bootDisk.([]interface{})) > 0 {
			brI := bootDisk.([]interface{})
			bootVal := brI[0].(map[string]interface{})

			diskOut := config.NewBootDeviceDisk()
			if diskAddress, ok := bootVal["disk_address"]; ok && len(diskAddress.([]interface{})) > 0 {
				daI := diskAddress.([]interface{})
				diskVal := daI[0].(map[string]interface{})
				diskAddOut := config.NewDiskAddress()
				if busType, ok := diskVal["bus_type"]; ok {
					const two, three, four, five, six = 2, 3, 4, 5, 6
					subMap := map[string]interface{}{
						"SCSI":  two,
						"IDE":   three,
						"PCI":   four,
						"SATA":  five,
						"SPAPR": six,
					}
					pVal := subMap[busType.(string)]
					p := config.DiskBusType(pVal.(int))
					diskAddOut.BusType = &p
				}
				if index, ok := diskVal["index"]; ok {
					diskAddOut.Index = utils.IntPtr(index.(int))
				}
				diskOut.DiskAddress = diskAddOut
			}
			BootDevice.SetValue(*diskOut)
		}
		if bootNic, ok := val["boot_device_nic"]; ok && len(bootNic.([]interface{})) > 0 {
			brI := bootNic.([]interface{})
			bootVal := brI[0].(map[string]interface{})

			diskNicOut := config.NewBootDeviceNic()

			if mac, ok := bootVal["mac_address"]; ok {
				diskNicOut.MacAddress = utils.StringPtr(mac.(string))
			}
			BootDevice.SetValue(*diskNicOut)
		}
		return BootDevice
	}
	return nil
}

func expandNvramDevice(pr interface{}) *config.NvramDevice {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		nvramOut := &config.NvramDevice{}

		if backingStorageInfo, ok := val["backing_storage_info"]; ok {
			nvramOut.BackingStorageInfo = expandVMDisk(backingStorageInfo)
		}
		return nvramOut
	}
	return nil
}

func expandVMDisk(pr interface{}) *config.VmDisk {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		vmDiskOutput := &config.VmDisk{}

		if diskSizeBytes, ok := val["disk_size_bytes"]; ok && diskSizeBytes.(int) > 0 {
			vmDiskOutput.DiskSizeBytes = utils.Int64Ptr(int64(diskSizeBytes.(int)))
		}
		if storageContainer, ok := val["storage_container"]; ok && len(storageContainer.([]interface{})) > 0 {
			vmDiskOutput.StorageContainer = expandVMDiskContainerReference(storageContainer)
		}
		if storageConfig, ok := val["storage_config"]; ok && len(storageConfig.([]interface{})) > 0 {
			vmDiskOutput.StorageConfig = expandVMDiskStorageConfig(storageConfig)
		}
		if datasource, ok := val["data_source"]; ok && len(datasource.([]interface{})) > 0 {
			vmDiskOutput.DataSource = expandDataSource(datasource)
		}
		return vmDiskOutput
	}
	return nil
}

func expandVMDiskStorageConfig(pr interface{}) *config.VmDiskStorageConfig {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		storageConfig := &config.VmDiskStorageConfig{}

		if isFlashMode, ok := val["is_flash_mode_enabled"]; ok {
			storageConfig.IsFlashModeEnabled = utils.BoolPtr(isFlashMode.(bool))
		}
		return storageConfig
	}
	return nil
}

func expandDataSource(pr interface{}) *config.DataSource {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		dsOut := &config.DataSource{}

		if ref, ok := val["reference"]; ok {
			dsOut.Reference = expandOneOfDataSourceReference(ref)
		}
		return dsOut
	}
	return nil
}

func expandOneOfDataSourceReference(pr interface{}) *config.OneOfDataSourceReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		dataOut := &config.OneOfDataSourceReference{}

		if imgRef, ok := val["image_reference"]; ok && len(imgRef.([]interface{})) > 0 {
			imgI := imgRef.([]interface{})
			imgVal := imgI[0].(map[string]interface{})

			imgOut := config.NewImageReference()

			if imgExtID, ok := imgVal["image_ext_id"]; ok && len(imgExtID.(string)) > 0 {
				imgOut.ImageExtId = utils.StringPtr(imgExtID.(string))
			}
			dataOut.SetValue(*imgOut)
		}
		if vmdiskRef, ok := val["vm_disk_reference"]; ok && len(vmdiskRef.([]interface{})) > 0 {
			vmI := vmdiskRef.([]interface{})
			vmVal := vmI[0].(map[string]interface{})

			vmDiskOut := config.NewVmDiskReference()

			if diskExtID, ok := vmVal["disk_ext_id"]; ok && len(diskExtID.(string)) > 0 {
				vmDiskOut.DiskExtId = utils.StringPtr(diskExtID.(string))
			}
			if diskAdd, ok := vmVal["disk_address"]; ok && len(diskAdd.([]interface{})) > 0 {
				daI := diskAdd.([]interface{})
				diskVal := daI[0].(map[string]interface{})
				diskAddOut := config.DiskAddress{}
				if busType, ok := diskVal["bus_type"]; ok {
					const two, three, four, five, six = 2, 3, 4, 5, 6
					subMap := map[string]interface{}{
						"SCSI":  two,
						"IDE":   three,
						"PCI":   four,
						"SATA":  five,
						"SPAPR": six,
					}
					pVal := subMap[busType.(string)]
					p := config.DiskBusType(pVal.(int))
					diskAddOut.BusType = &p
				}
				if index, ok := diskVal["index"]; ok {
					diskAddOut.Index = utils.IntPtr(index.(int))
				}
				vmDiskOut.DiskAddress = &diskAddOut
			}
			if vmRef, ok := vmVal["vm_reference"]; ok && len(vmRef.([]interface{})) > 0 {
				vmDiskOut.VmReference = expandVMReference(vmRef)
			}
			dataOut.SetValue(*vmDiskOut)
		}
		return dataOut
	}
	return nil
}

func expandVMDiskContainerReference(pr interface{}) *config.VmDiskContainerReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		ownerRef := &config.VmDiskContainerReference{}

		if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
			ownerRef.ExtId = utils.StringPtr(extID.(string))
		}
		return ownerRef
	}
	return nil
}

func expandVMReference(pr interface{}) *config.VmReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		vmRef := &config.VmReference{}

		if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
			vmRef.ExtId = utils.StringPtr(extID.(string))
		}
		return vmRef
	}
	return nil
}

func expandVtpmConfig(pr interface{}) *config.VtpmConfig {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		vtpmConfig := &config.VtpmConfig{}

		if isEnabled, ok := val["is_vtpm_enabled"]; ok {
			vtpmConfig.IsVtpmEnabled = utils.BoolPtr(isEnabled.(bool))
		}
		return vtpmConfig
	}
	return nil
}

func expandApcConfig(pr interface{}) *config.ApcConfig {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		apcConfig := &config.ApcConfig{}

		if isEnabled, ok := val["is_apc_enabled"]; ok {
			apcConfig.IsApcEnabled = utils.BoolPtr(isEnabled.(bool))
		}
		if cpuModel, ok := val["cpu_model"]; ok {
			apcConfig.CpuModel = expandCPUModelReference(cpuModel)
		}
		return apcConfig
	}
	return nil
}

func expandCPUModelReference(pr interface{}) *config.CpuModelReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		cpuModel := &config.CpuModelReference{}

		if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
			cpuModel.ExtId = utils.StringPtr(extID.(string))
		}
		if name, ok := val["name"]; ok {
			cpuModel.Name = utils.StringPtr(name.(string))
		}

		return cpuModel
	}
	return nil
}

func expandADSFVmStorageConfig(pr interface{}) *config.ADSFVmStorageConfig {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		storageConfig := &config.ADSFVmStorageConfig{}

		if isFlash, ok := val["is_flash_mode_enabled"]; ok {
			storageConfig.IsFlashModeEnabled = utils.BoolPtr(isFlash.(bool))
		}
		if qosConfig, ok := val["qos_config"]; ok {
			storageConfig.QosConfig = expandQosConfig(qosConfig)
		}

		return storageConfig
	}
	return nil
}

func expandQosConfig(pr interface{}) *config.QosConfig {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		qos := &config.QosConfig{}

		if throttled, ok := val["throttled_iops"]; ok {
			qos.ThrottledIops = utils.IntPtr(throttled.(int))
		}
		return qos
	}
	return nil
}

func expandDisk(disk []interface{}) []config.Disk {
	if len(disk) > 0 {
		diskList := make([]config.Disk, len(disk))

		for k, v := range disk {
			val := v.(map[string]interface{})
			disk := config.Disk{}

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				disk.ExtId = utils.StringPtr(extID.(string))
			}
			if diskAdd, ok := val["disk_address"]; ok {
				disk.DiskAddress = expandDiskAddress(diskAdd)
			}
			if backingInfo, ok := val["backing_info"]; ok {
				disk.BackingInfo = expandOneOfDiskBackingInfo(backingInfo)
			}

			diskList[k] = disk
		}
		return diskList
	}
	return nil
}

func expandDiskAddress(disk interface{}) *config.DiskAddress {
	if disk != nil && len(disk.([]interface{})) > 0 {
		daI := disk.([]interface{})
		diskVal := daI[0].(map[string]interface{})
		diskAddOut := config.DiskAddress{}
		if busType, ok := diskVal["bus_type"]; ok {
			const two, three, four, five, six = 2, 3, 4, 5, 6
			subMap := map[string]interface{}{
				"SCSI":  two,
				"IDE":   three,
				"PCI":   four,
				"SATA":  five,
				"SPAPR": six,
			}
			pVal := subMap[busType.(string)]
			p := config.DiskBusType(pVal.(int))
			diskAddOut.BusType = &p
		}
		if index, ok := diskVal["index"]; ok {
			diskAddOut.Index = utils.IntPtr(index.(int))
		}
		return &diskAddOut
	}
	return nil
}

func expandOneOfDiskBackingInfo(pr interface{}) *config.OneOfDiskBackingInfo {
	if pr != nil && len(pr.([]interface{})) > 0 {
		backInfoOut := &config.OneOfDiskBackingInfo{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if vmDisk, ok := val["vm_disk"]; ok && len(vmDisk.([]interface{})) > 0 {
			vmDiskOut := config.NewVmDisk()
			vmI := vmDisk.([]interface{})
			vmVal := vmI[0].(map[string]interface{})

			if diskBytes, ok := vmVal["disk_size_bytes"]; ok && diskBytes.(int) > 0 {
				vmDiskOut.DiskSizeBytes = utils.Int64Ptr(int64(diskBytes.(int)))
			}
			if storageContainer, ok := vmVal["storage_container"]; ok && len(storageContainer.([]interface{})) > 0 {
				vmDiskOut.StorageContainer = expandVMDiskContainerReference(storageContainer)
			}
			if storageCfg, ok := vmVal["storage_config"]; ok && len(storageCfg.([]interface{})) > 0 {
				vmDiskOut.StorageConfig = expandVMDiskStorageConfig(storageCfg)
			}
			if ds, ok := vmVal["data_source"]; ok && len(ds.([]interface{})) > 0 {
				vmDiskOut.DataSource = expandDataSource(ds)
			}

			backInfoOut.SetValue(*vmDiskOut)
		}
		if adfs, ok := val["adfs_volume_group_reference"]; ok && len(adfs.([]interface{})) > 0 {
			adfsOut := config.NewADSFVolumeGroupReference()
			adI := adfs.([]interface{})
			adVal := adI[0].(map[string]interface{})

			if vgExtID, ok := adVal["volume_group_ext_id"]; ok {
				adfsOut.VolumeGroupExtId = utils.StringPtr(vgExtID.(string))
			}

			backInfoOut.SetValue(*adfsOut)
		}
		return backInfoOut
	}
	return nil
}

func expandCdRom(cd []interface{}) []config.CdRom {
	if len(cd) > 0 {
		cdList := make([]config.CdRom, len(cd))

		for k, v := range cd {
			val := v.(map[string]interface{})
			cds := config.CdRom{}

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				cds.ExtId = utils.StringPtr(extID.(string))
			}
			if diskAdd, ok := val["disk_address"]; ok && len(diskAdd.([]interface{})) > 0 {
				cds.DiskAddress = expandCdRomAddress(diskAdd)
			}
			if backingInfo, ok := val["backing_info"]; ok && len(backingInfo.([]interface{})) > 0 {
				cds.BackingInfo = expandVMDisk(backingInfo)
			}
			if isoType, ok := val["iso_type"]; ok && len(isoType.(string)) > 0 {
				const two, three, four = 2, 3, 4
				subMap := map[string]interface{}{
					"OTHER":               two,
					"GUEST_TOOLS":         three,
					"GUEST_CUSTOMIZATION": four,
				}
				pVal := subMap[isoType.(string)]
				p := config.IsoType(pVal.(int))
				cds.IsoType = &p
			}

			cdList[k] = cds
		}
		return cdList
	}
	return nil
}

func expandCdRomAddress(disk interface{}) *config.CdRomAddress {
	if disk != nil && len(disk.([]interface{})) > 0 {
		cdRomAdd := &config.CdRomAddress{}
		adI := disk.([]interface{})
		adVal := adI[0].(map[string]interface{})

		if busType, ok := adVal["bus_type"]; ok {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"IDE":  two,
				"SATA": three,
			}
			pVal := subMap[busType.(string)]
			p := config.CdRomBusType(pVal.(int))
			cdRomAdd.BusType = &p
		}
		if index, ok := adVal["index"]; ok {
			cdRomAdd.Index = utils.IntPtr(index.(int))
		}
		return cdRomAdd
	}
	return nil
}

func expandGpu(gpu []interface{}) []config.Gpu {
	if len(gpu) > 0 {
		gpuList := make([]config.Gpu, len(gpu))

		for k, v := range gpu {
			gpus := config.Gpu{}
			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				gpus.ExtId = utils.StringPtr(extID.(string))
			}
			if mode, ok := val["mode"]; ok {
				const two, three, four = 2, 3, 4
				subMap := map[string]interface{}{
					"PASSTHROUGH_GRAPHICS": two,
					"PASSTHROUGH_COMPUTE":  three,
					"VIRTUAL":              four,
				}
				pVal := subMap[mode.(string)]
				p := config.GpuMode(pVal.(int))
				gpus.Mode = &p
			}
			if deviceID, ok := val["device_id"]; ok {
				gpus.DeviceId = utils.IntPtr(deviceID.(int))
			}
			if vendor, ok := val["vendor"]; ok {
				const two, three, four = 2, 3, 4
				subMap := map[string]interface{}{
					"NVIDIA": two,
					"INTEL":  three,
					"AMD":    four,
				}
				pVal := subMap[vendor.(string)]
				p := config.GpuVendor(pVal.(int))
				gpus.Vendor = &p
			}
			if pciAddress, ok := val["pci_address"]; ok && len(pciAddress.([]interface{})) > 0 {
				pciObj := config.SBDF{}
				pciI := pciAddress.([]interface{})
				pciVal := pciI[0].(map[string]interface{})

				if pciVal["segment"] != nil {
					pciObj.Segment = utils.IntPtr(pciVal["segment"].(int))
				}
				if pciVal["bus"] != nil {
					pciObj.Bus = utils.IntPtr(pciVal["bus"].(int))
				}
				if pciVal["device"] != nil {
					pciObj.Device = utils.IntPtr(pciVal["device"].(int))
				}
				if pciVal["func"] != nil {
					pciObj.Func = utils.IntPtr(pciVal["func"].(int))
				}
				gpus.PciAddress = &pciObj
			}
			gpuList[k] = gpus
		}
		return gpuList
	}
	return nil
}

func expandSerialPort(serial []interface{}) []config.SerialPort {
	if len(serial) > 0 {
		serialPortList := make([]config.SerialPort, len(serial))

		for k, v := range serial {
			val := v.(map[string]interface{})
			serials := config.SerialPort{}

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				serials.ExtId = utils.StringPtr(extID.(string))
			}
			if isConn, ok := val["is_connected"]; ok {
				serials.IsConnected = utils.BoolPtr(isConn.(bool))
			}
			if index, ok := val["index"]; ok {
				serials.Index = utils.IntPtr(index.(int))
			}

			serialPortList[k] = serials
		}

		return serialPortList
	}
	return nil
}

func expandProtectionPolicyState(pr interface{}) *config.ProtectionPolicyState {
	if pr != nil && len(pr.([]interface{})) > 0 {
		policyState := &config.ProtectionPolicyState{}

		prI := pr.([]interface{})
		prVal := prI[0].(map[string]interface{})

		if policy, ok := prVal["policy"]; ok {
			policyState.Policy = expandPolicyReference(policy)
		}
		return policyState
	}
	return nil
}

func expandPolicyReference(pr interface{}) *config.PolicyReference {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		ownerRef := &config.PolicyReference{}

		if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
			ownerRef.ExtId = utils.StringPtr(extID.(string))
		}
		return ownerRef
	}
	return nil
}

func flattenPowerState(pr *config.PowerState) string {
	if pr != nil {
		const two, three, four, five = 2, 3, 4, 5
		if *pr == config.PowerState(two) {
			return "ON"
		}
		if *pr == config.PowerState(three) {
			return "OFF"
		}
		if *pr == config.PowerState(four) {
			return "PAUSED"
		}
		if *pr == config.PowerState(five) {
			return "UNDETERMINED"
		}
	}
	return "UNKNOWN"
}

func callForPowerOffVM(ctx context.Context, conn *vmm.Client, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	readResp, errR := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
	if errR != nil {
		return diag.Errorf("error while reading vm : %v", errR)
	}
	// checking current state of VM
	vmResp := readResp.Data.GetValue().(config.Vm)

	if vmResp.PowerState.GetName() == "OFF" {
		log.Printf("[DEBUG] VM is already in %s state", d.Get("power_state").(string))
		return nil
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)

	// Power off the VM
	powerOffResp, err := conn.VMAPIInstance.PowerOffVm(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while powering off Virtual Machine : %v", err)
	}

	TaskRef := powerOffResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	prismConn := meta.(*conns.Client).PrismAPI

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, prismConn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for virtual machine (%s) to power off: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func callForPowerOnVM(ctx context.Context, conn *vmm.Client, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	readResp, errR := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
	if errR != nil {
		return diag.Errorf("error while reading vm : %v", errR)
	}

	vmResp := readResp.Data.GetValue().(config.Vm)

	if vmResp.PowerState.GetName() == "ON" {
		log.Printf("[DEBUG] VM is already in %s state", d.Get("power_state").(string))
		return nil
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = getEtagHeader(readResp, conn)
	// Power on the VM
	powerOnResp, err := conn.VMAPIInstance.PowerOnVm(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while powering on Virtual Machine : %v", err)
	}

	aJSON, _ := json.Marshal(powerOnResp)
	log.Printf("[DEBUG] PowerOn Response: %s", string(aJSON))

	TaskRef := powerOnResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	prismConn := meta.(*conns.Client).PrismAPI

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, prismConn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for virtual machine (%s) to power on: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func diffConfig(oldValue []interface{}, newValue []interface{}) ([]interface{}, []interface{}, []interface{}) {
	newlyAdded := make([]interface{}, 0)
	removed := make([]interface{}, 0)
	updated := make([]interface{}, 0)

	// check for newly added
	for _, newItem := range newValue {
		found := false
		for _, oldItem := range oldValue {
			if reflect.DeepEqual(newItem, oldItem) {
				found = true
				break
			}
		}
		if !found {
			newlyAdded = append(newlyAdded, newItem)
		}
	}

	// check for removed
	for _, oldItem := range oldValue {
		found := false
		for _, newItem := range newValue {
			if reflect.DeepEqual(oldItem, newItem) {
				found = true
				break
			}
		}
		if !found {
			removed = append(removed, oldItem)
		}
	}

	for _, oldItem := range oldValue {
		oldMap := oldItem.(map[string]interface{})

		if oldMap["ext_id"] != "" {
			for _, newItem := range newValue {
				if oldMap["ext_id"] == newItem.(map[string]interface{})["ext_id"] {
					// Only add to updated if the items are actually different
					if !reflect.DeepEqual(oldItem, newItem) {
						updated = append(updated, newItem)
					}
					break
				}
			}
		}
	}

	// deleting updated list from newly added and removed
	for _, item := range updated {
		itemMap := item.(map[string]interface{})

		for i, newItem := range newlyAdded {
			newMap := newItem.(map[string]interface{})

			if newMap["ext_id"] == itemMap["ext_id"] {
				newlyAdded = append(newlyAdded[:i], newlyAdded[i+1:]...)
			}
		}

		for i, oldItem := range removed {
			oldMap := oldItem.(map[string]interface{})

			if oldMap["ext_id"] == itemMap["ext_id"] {
				removed = append(removed[:i], removed[i+1:]...)
			}
		}
	}

	return newlyAdded, removed, updated
}

// Check if VM is in power off state to perform update operations
func checkForHotPlugChanges(d *schema.ResourceData) bool {
	if d.HasChange(("num_sockets")) || d.HasChange(("num_cores_per_socket")) || d.HasChange(("memory_size_bytes")) ||
		d.HasChange(("num_threads_per_core")) || d.HasChange(("cd_rom")) || d.HasChange(("num_numa_nodes")) ||
		d.HasChange("cluster") || d.HasChange("is_cpu_passthrough_enabled") || d.HasChange("enabled_cpu_features") ||
		d.HasChange("is_vcpu_hard_pinning_enabled") || d.HasChange("guest_customization") || d.HasChange("guest_tools") ||
		d.HasChange("serial_ports") || d.HasChange("gpus") || d.HasChange("boot_config") {
		return true
	}
	return false
}

func isVMPowerOff(d *schema.ResourceData, conn *vmm.Client) bool {
	readResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
	if err != nil {
		return false
	}
	vmResp := readResp.Data.GetValue().(config.Vm)

	return vmResp.PowerState.GetName() == "OFF"
}

// getFirstIPAddress returns the first available IP address from a NIC.
// It checks both DHCP learned IPs and statically configured IPs.
func getFirstIPAddress(nic config.Nic) string {
	if nic.NetworkInfo == nil {
		return ""
	}
	// Check for DHCP learned IPs first
	if nic.NetworkInfo.Ipv4Info != nil && len(nic.NetworkInfo.Ipv4Info.LearnedIpAddresses) > 0 {
		if nic.NetworkInfo.Ipv4Info.LearnedIpAddresses[0].Value != nil {
			return *nic.NetworkInfo.Ipv4Info.LearnedIpAddresses[0].Value
		}
	}
	// Check for statically configured IP
	if nic.NetworkInfo.Ipv4Config != nil && nic.NetworkInfo.Ipv4Config.IpAddress != nil && nic.NetworkInfo.Ipv4Config.IpAddress.Value != nil {
		return *nic.NetworkInfo.Ipv4Config.IpAddress.Value
	}
	return ""
}

func waitForIPRefreshFunc(client *vmm.Client, vmUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.VMAPIInstance.GetVmById(utils.StringPtr(vmUUID))
		if err != nil {
			return nil, "", err
		}

		getResp := resp.Data.GetValue().(config.Vm)

		if len(getResp.Nics) > 0 {
			for _, nic := range getResp.Nics {
				if nic.NetworkInfo != nil {
					// Check for DHCP learned IPs
					if nic.NetworkInfo.Ipv4Info != nil {
						for _, ip := range nic.NetworkInfo.Ipv4Info.LearnedIpAddresses {
							if ip.Value != nil {
								return resp, "AVAILABLE", nil
							}
						}
					}
					// Check for statically configured IPs
					if nic.NetworkInfo.Ipv4Config != nil && nic.NetworkInfo.Ipv4Config.IpAddress != nil && nic.NetworkInfo.Ipv4Config.IpAddress.Value != nil {
						return resp, "AVAILABLE", nil
					}
				}
			}
		}
		return resp, "WAITING", nil
	}
}

func expandProjectReference(pr []interface{}) *config.ProjectReference {
	if len(pr) > 0 {
		val := pr[0].(map[string]interface{})
		project := config.ProjectReference{}
		if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
			log.Printf("[DEBUG] Project ExtID: %s", extID.(string))
			project.ExtId = utils.StringPtr(extID.(string))
		}
		return &project
	}
	return nil
}

func prepareVMConfigFromMap(m map[string]interface{}) *config.Vm {
	body := &config.Vm{}
	if extID, ok := m["ext_id"]; ok {
		body.ExtId = utils.StringPtr(extID.(string))
	}
	if name, ok := m["name"]; ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := m["description"]; ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if source, ok := m["source"]; ok {
		body.Source = expandVMSourceReference(source)
	}
	if numSock, ok := m["num_sockets"]; ok {
		body.NumSockets = utils.IntPtr(numSock.(int))
	}
	if numCoresPerSock, ok := m["num_cores_per_socket"]; ok {
		body.NumCoresPerSocket = utils.IntPtr(numCoresPerSock.(int))
	}
	if numThreadsPerCore, ok := m["num_threads_per_core"]; ok {
		body.NumThreadsPerCore = utils.IntPtr(numThreadsPerCore.(int))
	}
	if numNumaNodes, ok := m["num_numa_nodes"]; ok {
		body.NumNumaNodes = utils.IntPtr(numNumaNodes.(int))
	}
	if memorySize, ok := m["memory_size_bytes"]; ok {
		body.MemorySizeBytes = utils.Int64Ptr(int64(memorySize.(int)))
	}
	if isVcpuEnabled, ok := m["is_vcpu_hard_pinning_enabled"]; ok {
		body.IsVcpuHardPinningEnabled = utils.BoolPtr(isVcpuEnabled.(bool))
	}
	if isCPUPassThrough, ok := m["is_cpu_passthrough_enabled"]; ok {
		body.IsCpuPassthroughEnabled = utils.BoolPtr(isCPUPassThrough.(bool))
	}
	if cpuFeatures, ok := m["enabled_cpu_features"]; ok {
		body.EnabledCpuFeatures = expandCPUFeature(cpuFeatures.([]interface{}))
	}
	if memoryOverCommit, ok := m["is_memory_overcommit_enabled"]; ok {
		body.IsMemoryOvercommitEnabled = utils.BoolPtr(memoryOverCommit.(bool))
	}
	if isGpuConsole, ok := m["is_gpu_console_enabled"]; ok {
		body.IsGpuConsoleEnabled = utils.BoolPtr(isGpuConsole.(bool))
	}
	if isCPUHotplugEnabled, ok := m["is_cpu_hotplug_enabled"]; ok {
		body.IsCpuHotplugEnabled = utils.BoolPtr(isCPUHotplugEnabled.(bool))
	}
	if isScsiControllerEnabled, ok := m["is_scsi_controller_enabled"]; ok {
		body.IsScsiControllerEnabled = utils.BoolPtr(isScsiControllerEnabled.(bool))
	}
	if genUUID, ok := m["generation_uuid"]; ok {
		body.GenerationUuid = utils.StringPtr(genUUID.(string))
	}
	if bios, ok := m["bios_uuid"]; ok {
		body.BiosUuid = utils.StringPtr(bios.(string))
	}
	if categories, ok := m["categories"]; ok {
		body.Categories = expandCategoryReference(categories.([]interface{}))
	}
	if project, ok := m["project"]; ok {
		body.Project = expandProjectReference(project.([]interface{}))
	}
	if ownerRef, ok := m["ownership_info"]; ok {
		body.OwnershipInfo = expandOwnershipInfo(ownerRef)
	}
	if host, ok := m["host"]; ok {
		body.Host = expandHostReference(host)
	}
	if cls, ok := m["cluster"]; ok {
		body.Cluster = expandClusterReference(cls)
	}
	if availabilityZone, ok := m["availability_zone"]; ok {
		body.AvailabilityZone = expandAvailabilityZoneReference(availabilityZone)
	}
	if guestCstm, ok := m["guest_customization"]; ok {
		body.GuestCustomization = expandTemplateGuestCustomizationParams(guestCstm)
	}
	if guestTools, ok := m["guest_tools"]; ok {
		body.GuestTools = expandGuestTools(guestTools)
	}
	if hardwareClock, ok := m["hardware_clock_timezone"]; ok {
		body.HardwareClockTimezone = utils.StringPtr(hardwareClock.(string))
	}
	if isBranding, ok := m["is_branding_enabled"]; ok {
		body.IsBrandingEnabled = utils.BoolPtr(isBranding.(bool))
	}
	if bootConfig, ok := m["boot_config"]; ok {
		body.BootConfig = expandOneOfVMBootConfig(bootConfig)
	}
	if vgaConsole, ok := m["is_vga_console_enabled"]; ok {
		body.IsVgaConsoleEnabled = utils.BoolPtr(vgaConsole.(bool))
	}
	if machineType, ok := m["machine_type"]; ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"PC":      two,
			"PSERIES": three,
			"Q35":     four,
		}
		if val, ok := machineType.(string); ok {
			if pVal, exists := subMap[val]; exists {
				p := config.MachineType(pVal.(int))
				body.MachineType = &p
			}
		}
	}
	if vtpmConfig, ok := m["vtpm_config"]; ok {
		body.VtpmConfig = expandVtpmConfig(vtpmConfig)
	}
	if isAgentVM, ok := m["is_agent_vm"]; ok {
		body.IsAgentVm = utils.BoolPtr(isAgentVM.(bool))
	}
	if apcConfig, ok := m["apc_config"]; ok {
		body.ApcConfig = expandApcConfig(apcConfig)
	}
	if storageConfig, ok := m["storage_config"]; ok {
		body.StorageConfig = expandADSFVmStorageConfig(storageConfig)
	}
	if disks, ok := m["disks"]; ok {
		body.Disks = expandDisk(disks.([]interface{}))
	}
	if cdroms, ok := m["cd_roms"]; ok {
		body.CdRoms = expandCdRom(cdroms.([]interface{}))
	}
	if nics, ok := m["nics"]; ok {
		body.Nics = expandNic(nics.([]interface{}))
	}
	if gpus, ok := m["gpus"]; ok {
		body.Gpus = expandGpu(gpus.([]interface{}))
	}
	if serialPorts, ok := m["serial_ports"]; ok {
		body.SerialPorts = expandSerialPort(serialPorts.([]interface{}))
	}
	if protectionType, ok := m["protection_type"]; ok {
		const two, three, four = 2, 3, 4
		subMap := map[string]interface{}{
			"UNPROTECTED":    two,
			"PD_PROTECTED":   three,
			"RULE_PROTECTED": four,
		}
		if val, ok := protectionType.(string); ok {
			if pVal, exists := subMap[val]; exists {
				p := config.ProtectionType(pVal.(int))
				body.ProtectionType = &p
			}
		}
	}
	if protectionPolicyState, ok := m["protection_policy_state"]; ok {
		body.ProtectionPolicyState = expandProtectionPolicyState(protectionPolicyState)
	}
	return body
}

func resourceDataToMap(d *schema.ResourceData, schemaMap map[string]*schema.Schema) map[string]interface{} {
	result := make(map[string]interface{})
	for key := range schemaMap {
		if v, ok := d.GetOk(key); ok {
			result[key] = v
		}
	}
	return result
}

func extractVMConfigFields(getResp config.Vm) (map[string]interface{}, diag.Diagnostics) {
	fields := make(map[string]interface{})
	var diags diag.Diagnostics

	fields["name"] = getResp.Name
	fields["ext_id"] = getResp.ExtId
	fields["description"] = getResp.Description

	if getResp.CreateTime != nil {
		fields["create_time"] = getResp.CreateTime.String()
	}
	if getResp.UpdateTime != nil {
		fields["update_time"] = getResp.UpdateTime.String()
	}
	fields["source"] = flattenVMSourceReference(getResp.Source)
	fields["num_sockets"] = getResp.NumSockets
	fields["num_cores_per_socket"] = getResp.NumCoresPerSocket
	fields["num_threads_per_core"] = getResp.NumThreadsPerCore
	fields["num_numa_nodes"] = getResp.NumNumaNodes
	fields["memory_size_bytes"] = getResp.MemorySizeBytes
	fields["is_vcpu_hard_pinning_enabled"] = getResp.IsVcpuHardPinningEnabled
	fields["is_cpu_passthrough_enabled"] = getResp.IsCpuPassthroughEnabled
	fields["enabled_cpu_features"] = flattenCPUFeature(getResp.EnabledCpuFeatures)
	fields["is_memory_overcommit_enabled"] = getResp.IsMemoryOvercommitEnabled
	fields["is_gpu_console_enabled"] = getResp.IsGpuConsoleEnabled
	fields["is_cpu_hotplug_enabled"] = getResp.IsCpuHotplugEnabled
	fields["is_scsi_controller_enabled"] = getResp.IsScsiControllerEnabled
	fields["generation_uuid"] = getResp.GenerationUuid
	fields["bios_uuid"] = getResp.BiosUuid
	fields["categories"] = flattenCategoryReference(getResp.Categories)
	fields["project"] = flattenProjectReference(getResp.Project)
	fields["ownership_info"] = flattenOwnershipInfo(getResp.OwnershipInfo)
	fields["host"] = flattenHostReference(getResp.Host)
	fields["cluster"] = flattenClusterReference(getResp.Cluster)
	fields["availability_zone"] = flattenAvailabilityZoneReference(getResp.AvailabilityZone)
	fields["guest_customization"] = flattenGuestCustomizationParams(getResp.GuestCustomization)
	fields["guest_tools"] = flattenGuestTools(getResp.GuestTools)
	fields["hardware_clock_timezone"] = getResp.HardwareClockTimezone
	fields["is_branding_enabled"] = getResp.IsBrandingEnabled
	fields["boot_config"] = flattenOneOfVMBootConfig(getResp.BootConfig)
	fields["is_vga_console_enabled"] = getResp.IsVgaConsoleEnabled
	fields["machine_type"] = flattenMachineType(getResp.MachineType)
	fields["power_state"] = flattenPowerState(getResp.PowerState)
	fields["vtpm_config"] = flattenVtpmConfig(getResp.VtpmConfig)
	fields["is_agent_vm"] = getResp.IsAgentVm
	fields["apc_config"] = flattenApcConfig(getResp.ApcConfig)
	fields["storage_config"] = flattenADSFVmStorageConfig(getResp.StorageConfig)
	fields["disks"] = flattenDisk(getResp.Disks)
	fields["cd_roms"] = flattenCdRom(getResp.CdRoms)
	fields["nics"] = flattenNic(getResp.Nics)
	fields["gpus"] = flattenGpu(getResp.Gpus)
	fields["serial_ports"] = flattenSerialPort(getResp.SerialPorts)
	fields["protection_type"] = flattenProtectionType(getResp.ProtectionType)
	fields["protection_policy_state"] = flattenProtectionPolicyState(getResp.ProtectionPolicyState)

	return fields, diags
}

func setVMConfig(d *schema.ResourceData, getResp config.Vm) diag.Diagnostics {
	fields, diags := extractVMConfigFields(getResp)
	if diags.HasError() {
		return diags
	}
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("failed setting %q: %w", k, err))
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
	if err := d.Set("memory_size_bytes", getResp.MemorySizeBytes); err != nil {
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
	if err := d.Set("project", flattenProjectReference(getResp.Project)); err != nil {
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
	if err := d.Set("availability_zone", flattenAvailabilityZoneReference(getResp.AvailabilityZone)); err != nil {
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
		return diag.FromErr(err)
	}
	if err := d.Set("cd_roms", flattenCdRom(getResp.CdRoms)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nics", flattenNic(getResp.Nics)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpus", flattenGpu(getResp.Gpus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("serial_ports", flattenSerialPort(getResp.SerialPorts)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("protection_type", flattenProtectionType(getResp.ProtectionType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("protection_policy_state", flattenProtectionPolicyState(getResp.ProtectionPolicyState)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
