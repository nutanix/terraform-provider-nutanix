package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixOvaVMDeploymentV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixOvaVMDeploymentCreate,
		ReadContext:   ResourceNutanixOvaVMDeploymentRead,
		UpdateContext: ResourceNutanixOvaVMDeploymentUpdate,
		DeleteContext: ResourceNutanixOvaVMDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"override_vm_config": {
				Type:     schema.TypeList,
				Required: true,
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
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"backing_info": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"model": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mac_address": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"is_connected": {
													Type:     schema.TypeBool,
													Optional: true,
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
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"nic_type": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"SPAN_DESTINATION_NIC",
														"NORMAL_NIC", "DIRECT_NIC", "NETWORK_FUNCTION_NIC",
													}, false),
												},
												"network_function_chain": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ext_id": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"network_function_nic_type": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"TAP", "EGRESS",
														"INGRESS",
													}, false),
												},
												"subnet": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ext_id": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"vlan_mode": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"TRUNK", "ACCESS"}, false),
												},
												"trunked_vlans": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"should_allow_unknown_macs": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ipv4_config": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"should_assign_ip": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"ip_address": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Optional: true,
																		},
																	},
																},
															},
															"secondary_ip_address_list": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"ipv4_info": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"learned_ip_addresses": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"prefix_length": {
																			Type:     schema.TypeInt,
																			Optional: true,
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
						"cd_roms": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_address": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bus_type": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"IDE", "SATA"}, false),
												},
												"index": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"backing_info": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disk_size_bytes": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"storage_container": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ext_id": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"storage_config": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_flash_mode_enabled": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
												},
												"data_source": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"reference": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_reference": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_ext_id": {
																						Type:     schema.TypeString,
																						Optional: true,
																					},
																				},
																			},
																		},
																		"vm_disk_reference": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"disk_ext_id": {
																						Type:     schema.TypeString,
																						Optional: true,
																					},
																					"disk_address": {
																						Type:     schema.TypeList,
																						Optional: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bus_type": {
																									Type:     schema.TypeString,
																									Optional: true,
																									ValidateFunc: validation.StringInSlice([]string{
																										"SCSI", "SPAPR", "PCI",
																										"IDE", "SATA",
																									}, false),
																								},
																								"index": {
																									Type:     schema.TypeInt,
																									Optional: true,
																								},
																							},
																						},
																					},
																					"vm_reference": {
																						Type:     schema.TypeList,
																						Optional: true,
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"ext_id": {
																									Type:     schema.TypeString,
																									Optional: true,
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
									"iso_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"OTHER", "GUEST_TOOLS", "GUEST_CUSTOMIZATION"}, false),
									},
								},
							},
						},
						"categories": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"cluster_location_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ResourceNutanixOvaVMDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)
	vmDeploymentSpec := &import1.OvaDeploymentSpec{}
	if clusterLocationExtID := d.Get("cluster_location_ext_id").(string); clusterLocationExtID != "" {
		vmDeploymentSpec.ClusterLocationExtId = &clusterLocationExtID
	}
	if overrideVMConfig, ok := d.GetOk("override_vm_config"); ok {
		overrideVMConfigList := overrideVMConfig.([]interface{})
		if len(overrideVMConfigList) > 0 && overrideVMConfigList[0] != nil {
			ovm := overrideVMConfigList[0].(map[string]interface{})
			overrideSpec := &import2.OvaVmConfigOverrideSpec{}

			if v, ok := ovm["name"].(string); ok && v != "" {
				overrideSpec.Name = &v
			}
			if v, ok := ovm["num_sockets"].(int); ok && v != 0 {
				overrideSpec.NumSockets = &v
			}
			if v, ok := ovm["num_cores_per_socket"].(int); ok && v != 0 {
				overrideSpec.NumCoresPerSocket = &v
			}
			if v, ok := ovm["num_threads_per_core"].(int); ok && v != 0 {
				overrideSpec.NumThreadsPerCore = &v
			}
			if v, ok := ovm["memory_size_bytes"].(int); ok && v != 0 {
				mem := int64(v)
				overrideSpec.MemorySizeBytes = &mem
			}
			// Handle categories
			if cats, ok := ovm["categories"]; ok {
				overrideSpec.Categories = expandCategoryReference(cats.([]interface{}))
			}
			// Handle nics
			if nics, ok := ovm["nics"]; ok {
				overrideSpec.Nics = expandNic(nics.([]interface{}))
			}
			// Handle cd_roms
			if cdroms, ok := ovm["cd_roms"]; ok {
				overrideSpec.CdRoms = expandCdRom(cdroms.([]interface{}))
			}
			vmDeploymentSpec.OverrideVmConfig = overrideSpec
		}
	}

	resp, err := conn.OvasAPIInstance.DeployOva(&extID, vmDeploymentSpec)
	if err != nil {
		return diag.FromErr(err)
	}

	TaskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the OVA VM to be deployed
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for OVA VM deployment (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(extID)
	return nil
}

func ResourceNutanixOvaVMDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixOvaVMDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixOvaVMDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
