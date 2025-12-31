package vmmv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import4 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	ovaVMDeployTimeout = 30 * time.Minute
	ovaVMDeployDelay   = 10 * time.Second
)

func ResourceNutanixOvaVMDeploymentV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixOvaVMDeploymentCreate,
		ReadContext:   ResourceNutanixOvaVMDeploymentRead,
		UpdateContext: ResourceNutanixOvaVMDeploymentUpdate,
		DeleteContext: ResourceNutanixOvaVMDeploymentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ovaVMDeployTimeout),
			Update: schema.DefaultTimeout(ovaVMDeployTimeout),
			Delete: schema.DefaultTimeout(ovaVMDeployTimeout),
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique identifier of the OVA image to deploy VM from.",
			},
			"override_vm_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Override VM configuration parameters when deploying from OVA.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the VM to be deployed from OVA.",
						},
						"num_sockets": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Number of sockets for the VM CPU configuration.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"num_cores_per_socket": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Number of cores per socket for the VM CPU configuration.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"num_threads_per_core": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Number of threads per core for the VM CPU configuration.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"memory_size_bytes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Memory size of the VM in bytes.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"power_state": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "ON",
							Description:  "Power state of the VM (ON or OFF).",
							ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
						},
						"nics": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Network interface controllers (NICs) for the VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique identifier of the NIC.",
									},
									"backing_info": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Backing information for the NIC including model and MAC address.",
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
						"disks": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Additional disks to attach to the VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
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
			if cats, ok := ovm["categories"]; ok {
				overrideSpec.Categories = expandCategoryReference(cats.([]interface{}))
			}
			if nics, ok := ovm["nics"]; ok {
				overrideSpec.Nics = expandNic(nics.([]interface{}))
			}
			if cdroms, ok := ovm["cd_roms"]; ok {
				overrideSpec.CdRoms = expandCdRom(cdroms.([]interface{}))
			}
			vmDeploymentSpec.OverrideVmConfig = overrideSpec
		}
	}

	log.Printf("[DEBUG] Calling DeployOva API with OVA ext_id: %s", extID)
	resp, err := conn.OvasAPIInstance.DeployOva(&extID, vmDeploymentSpec)
	if err != nil {
		log.Printf("[ERROR] Failed to deploy OVA: %v", err)
		return diag.FromErr(err)
	}

	TaskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId
	log.Printf("[DEBUG] OVA deployment task started with UUID: %s", utils.StringValue(taskUUID))

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"QUEUED", "RUNNING", "PENDING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        ovaVMDeployDelay,
		PollInterval: ovaVMDeployDelay,
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		log.Printf("[ERROR] OVA deployment task failed with UUID %s: %v", utils.StringValue(taskUUID), errWaitTask)
		return diag.Errorf("error in OVA deployment (%s): %s", utils.StringValue(taskUUID), errWaitTask)
	}

	log.Printf("[DEBUG] OVA deployment task completed successfully with UUID: %s", utils.StringValue(taskUUID))

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		return diag.Errorf("error while fetching vm UUID : %v", err)
	}
	taskResult := resourceUUID.Data.GetValue().(import4.Task)

	if len(taskResult.EntitiesAffected) == 0 {
		return diag.Errorf("no entities affected in OVA deployment task")
	}

	var vmUUID *string
	for _, entity := range taskResult.EntitiesAffected {
		if entity.Rel != nil && *entity.Rel == "vmm:ahv:config:vm" {
			vmUUID = entity.ExtId
			log.Printf("[DEBUG] Found VM entity in task result: %s", *vmUUID)
			break
		}
	}

	if vmUUID == nil {
		return diag.Errorf("VM entity (vmm:ahv:vm) not found in task result")
	}

	d.SetId(*vmUUID)
	log.Printf("[DEBUG] OVA VM deployment completed successfully: vm_id=%s", *vmUUID)

	// Handle additional disks after initial VM deployment
	// OVA deployment doesn't support disks in the initial deployment, so add them separately
	if overrideVMConfig, ok := d.GetOk("override_vm_config"); ok {
		overrideVMConfigList := overrideVMConfig.([]interface{})
		if len(overrideVMConfigList) > 0 && overrideVMConfigList[0] != nil {
			overrideConfig := overrideVMConfigList[0].(map[string]interface{})

			if disks, exists := overrideConfig["disks"]; exists && disks != nil {
				disksList := disks.([]interface{})
				if len(disksList) > 0 {
					log.Printf("[DEBUG] Adding %d disks to OVA VM", len(disksList))
					for _, disk := range disksList {
						diskInput := expandDisk([]interface{}{disk})[0]

						readVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
						if err != nil {
							return diag.Errorf("error reading VM for disk creation: %v", err)
						}

						args := make(map[string]interface{})
						args["If-Match"] = getEtagHeader(readVMResp, conn)

						resp, err := conn.VMAPIInstance.CreateDisk(utils.StringPtr(d.Id()), &diskInput, args)
						if err != nil {
							return diag.Errorf("error creating disk: %v", err)
						}
						TaskRef := resp.Data.GetValue().(import3.TaskReference)
						diskTaskUUID := TaskRef.ExtId

						if err := waitForTask(ctx, d, meta, diskTaskUUID, schema.TimeoutCreate, "disk creation"); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	// Handle initial power state if specified as ON after initial VM deployment
	// OVA deployment doesn't support power state in the initial deployment, so add them separately
	if overrideVMConfig, ok := d.GetOk("override_vm_config"); ok {
		overrideVMConfigList := overrideVMConfig.([]interface{})
		if len(overrideVMConfigList) > 0 && overrideVMConfigList[0] != nil {
			overrideConfig := overrideVMConfigList[0].(map[string]interface{})
			if powerState, exists := overrideConfig["power_state"]; exists && powerState.(string) == "ON" {
				log.Printf("[DEBUG] Powering on VM after deployment as requested in configuration")
				if err := callForPowerOnVM(ctx, conn, d, meta); err != nil {
					return err
				}
			}
		}
	}

	// After all disks are created and powrer is handled, read the VM again to get the updated information
	log.Printf("[DEBUG] Reading VM after disk creation to update state")
	return ResourceNutanixOvaVMDeploymentRead(ctx, d, meta)
}

func ResourceNutanixOvaVMDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching vm : %v", err)
	}

	getResp := resp.Data.GetValue().(import2.Vm)
	return setOvaVMConfig(d, getResp)
}

func ResourceNutanixOvaVMDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	log.Printf("[DEBUG] starting OVA VM update for VM ID: %s", d.Id())

	// Check for hot-plug changes that require VM power off
	hotPlugRequired := false
	var oldConfig, newConfig interface{}
	hasOverrideVMConfigChange := d.HasChange("override_vm_config")

	// Consolidated override_vm_config change handling block
	if hasOverrideVMConfigChange {
		oldConfig, newConfig = d.GetChange("override_vm_config")
		if oldList, ok := oldConfig.([]interface{}); ok && len(oldList) > 0 {
			if newList, ok := newConfig.([]interface{}); ok && len(newList) > 0 {
				oldMap := oldList[0].(map[string]interface{})
				newMap := newList[0].(map[string]interface{})

				// Check if hot-plug sensitive fields changed
				hotPlugFields := []string{"num_sockets", "num_cores_per_socket", "num_threads_per_core", "memory_size_bytes"}
				for _, field := range hotPlugFields {
					if oldMap[field] != newMap[field] {
						hotPlugRequired = true
						log.Printf("[DEBUG] hot-plug change detected for field: %s", field)
						break
					}
				}
			}
		}

		// Power off VM if hot-plug changes are required
		if hotPlugRequired && !isVMPowerOff(d, conn) {
			log.Printf("[DEBUG] VM needs to be powered off for hot-plug changes")
			if err := callForPowerOffVM(ctx, conn, d, meta); err != nil {
				return err
			}
		}
		if err := handleDiskChanges(ctx, d, meta, meta, oldConfig, newConfig); err != nil {
			return err
		}

		if err := handlePowerStateChanges(ctx, d, meta, meta, oldConfig, newConfig, hotPlugRequired); err != nil {
			return err
		}

		if err := handleVMConfigurationUpdates(ctx, d, meta, meta, newConfig, hasOverrideVMConfigChange); err != nil {
			return err
		}
	}

	// Power VM back on if it was powered off for hot-plug changes
	if hotPlugRequired {
		// Check if the desired power state is ON
		if overrideVMConfig, ok := d.GetOk("override_vm_config"); ok {
			overrideVMConfigList := overrideVMConfig.([]interface{})
			if len(overrideVMConfigList) > 0 {
				overrideConfig := overrideVMConfigList[0].(map[string]interface{})
				if powerState, exists := overrideConfig["power_state"]; exists && powerState.(string) == "ON" {
					log.Printf("[DEBUG] Powering VM back on after hot-plug changes")
					if err := callForPowerOnVM(ctx, conn, d, meta); err != nil {
						return err
					}
				}
			}
		}
	}

	log.Printf("[DEBUG] OVA VM update completed successfully")
	return ResourceNutanixOvaVMDeploymentRead(ctx, d, meta)
}

func ResourceNutanixOvaVMDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	TaskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	// Wait for VM deletion to complete
	if err := waitForTask(ctx, d, meta, taskUUID, schema.TimeoutDelete, "VM deletion"); err != nil {
		return err
	}
	return nil
}

func setOvaVMConfig(d *schema.ResourceData, vm import2.Vm) diag.Diagnostics {
	log.Printf("[DEBUG] Setting OVA VM state from API response")

	// Get existing override_vm_config to preserve all user configuration
	if v, ok := d.GetOk("override_vm_config"); ok {
		existingList := v.([]interface{})
		if len(existingList) > 0 {
			if existingConfig, ok := existingList[0].(map[string]interface{}); ok {
				overrideConfig := make(map[string]interface{})
				for k, v := range existingConfig {
					overrideConfig[k] = v
				}

				if vm.Name != nil {
					overrideConfig["name"] = utils.StringValue(vm.Name)
				}

				if vm.NumSockets != nil && utils.IntValue(vm.NumSockets) > 0 {
					overrideConfig["num_sockets"] = utils.IntValue(vm.NumSockets)
					log.Printf("[DEBUG] Preserved num_sockets from API: %d", utils.IntValue(vm.NumSockets))
				} else {
					log.Printf("[DEBUG] API did not return num_sockets or returned 0, preserving user config")
				}

				if vm.NumCoresPerSocket != nil && utils.IntValue(vm.NumCoresPerSocket) > 0 {
					overrideConfig["num_cores_per_socket"] = utils.IntValue(vm.NumCoresPerSocket)
					log.Printf("[DEBUG] Preserved num_cores_per_socket from API: %d", utils.IntValue(vm.NumCoresPerSocket))
				} else {
					log.Printf("[DEBUG] API did not return num_cores_per_socket or returned 0, preserving user config")
				}

				if vm.NumThreadsPerCore != nil && utils.IntValue(vm.NumThreadsPerCore) > 0 {
					overrideConfig["num_threads_per_core"] = utils.IntValue(vm.NumThreadsPerCore)
					log.Printf("[DEBUG] Preserved num_threads_per_core from API: %d", utils.IntValue(vm.NumThreadsPerCore))
				} else {
					log.Printf("[DEBUG] API did not return num_threads_per_core or returned 0, preserving user config")
				}

				if vm.MemorySizeBytes != nil && utils.Int64Value(vm.MemorySizeBytes) > 0 {
					overrideConfig["memory_size_bytes"] = int(utils.Int64Value(vm.MemorySizeBytes))
				}
				if vm.PowerState != nil {
					overrideConfig["power_state"] = vm.PowerState.GetName()
					log.Printf("[DEBUG] Set power_state: %s", vm.PowerState.GetName())
				}

				// Update disks to capture ext_id from API by matching bus_type and index
				if len(vm.Disks) > 0 && overrideConfig["disks"] != nil {
					if existingDisks, ok := overrideConfig["disks"].([]interface{}); ok {
						for _, apiDisk := range vm.Disks {
							if apiDisk.ExtId == nil || apiDisk.DiskAddress == nil {
								continue
							}

							var apiBusType string
							var apiIndex *int
							if apiDisk.DiskAddress.BusType != nil {
								switch *apiDisk.DiskAddress.BusType {
								case import2.DISKBUSTYPE_SCSI:
									apiBusType = "SCSI"
								case import2.DISKBUSTYPE_IDE:
									apiBusType = "IDE"
								case import2.DISKBUSTYPE_PCI:
									apiBusType = "PCI"
								case import2.DISKBUSTYPE_SATA:
									apiBusType = "SATA"
								case import2.DISKBUSTYPE_SPAPR:
									apiBusType = "SPAPR"
								}
							}
							if apiDisk.DiskAddress.Index != nil {
								apiIndex = apiDisk.DiskAddress.Index
							}

							for i, existingDiskInterface := range existingDisks {
								if existingDisk, ok := existingDiskInterface.(map[string]interface{}); ok {
									var configBusType string
									var configIndex *int

									if diskAddress, exists := existingDisk["disk_address"]; exists {
										if diskAddressList, ok := diskAddress.([]interface{}); ok && len(diskAddressList) > 0 {
											if diskAddressMap, ok := diskAddressList[0].(map[string]interface{}); ok {
												if busType, exists := diskAddressMap["bus_type"]; exists {
													configBusType = busType.(string)
												}
												if index, exists := diskAddressMap["index"]; exists {
													if indexInt, ok := index.(int); ok {
														configIndex = &indexInt
													}
												}
											}
										}
									}

									busTypeMatches := apiBusType == configBusType
									indexMatches := (apiIndex == nil && configIndex == nil) ||
										(apiIndex != nil && configIndex != nil && *apiIndex == *configIndex)

									if busTypeMatches && indexMatches {
										existingDisk["ext_id"] = utils.StringValue(apiDisk.ExtId)
										log.Printf("[DEBUG] Matched and updated disk %d with bus_type=%s, index=%v, ext_id=%s",
											i, apiBusType, apiIndex, utils.StringValue(apiDisk.ExtId))
										break
									}
								}
							}
						}
					}
				}

				if len(vm.Nics) > 0 {
					nicsList := make([]interface{}, 0)

					for _, nic := range vm.Nics {
						nicMap := make(map[string]interface{})

						if nic.ExtId != nil {
							nicMap["ext_id"] = utils.StringValue(nic.ExtId)
						}

						if existingNics, ok := existingConfig["nics"].([]interface{}); ok && len(existingNics) > len(nicsList) {
							if existingNic, ok := existingNics[len(nicsList)].(map[string]interface{}); ok {
								if existingBackingInfo, ok := existingNic["backing_info"]; ok {
									nicMap["backing_info"] = existingBackingInfo
								}
							}
						}

						if nic.NetworkInfo != nil {
							networkInfoList := make([]map[string]interface{}, 0)
							networkInfo := make(map[string]interface{})

							if nic.NetworkInfo.NicType != nil {
								networkInfo["nic_type"] = nic.NetworkInfo.NicType.GetName()
							}

							if nic.NetworkInfo.Subnet != nil && nic.NetworkInfo.Subnet.ExtId != nil {
								subnetList := make([]map[string]interface{}, 0)
								subnet := make(map[string]interface{})
								subnet["ext_id"] = utils.StringValue(nic.NetworkInfo.Subnet.ExtId)
								subnetList = append(subnetList, subnet)
								networkInfo["subnet"] = subnetList
							}

							if nic.NetworkInfo.VlanMode != nil {
								networkInfo["vlan_mode"] = nic.NetworkInfo.VlanMode.GetName()
							}

							if len(nic.NetworkInfo.TrunkedVlans) > 0 {
								networkInfo["trunked_vlans"] = nic.NetworkInfo.TrunkedVlans
							} else {
								networkInfo["trunked_vlans"] = []int{}
							}

							if nic.NetworkInfo.ShouldAllowUnknownMacs != nil {
								networkInfo["should_allow_unknown_macs"] = utils.BoolValue(nic.NetworkInfo.ShouldAllowUnknownMacs)
							}

							networkInfoList = append(networkInfoList, networkInfo)
							nicMap["network_info"] = networkInfoList
						}

						nicsList = append(nicsList, nicMap)
					}

					overrideConfig["nics"] = nicsList
					log.Printf("[DEBUG] Updated NICs configuration with %d NICs", len(nicsList))
				}

				// Set the complete override_vm_config with preserved user settings
				overrideConfigList := []map[string]interface{}{overrideConfig}
				if err := d.Set("override_vm_config", overrideConfigList); err != nil {
					return diag.FromErr(fmt.Errorf("failed setting override_vm_config: %w", err))
				}
			}
		}
	}

	log.Printf("[DEBUG] OVA VM state set successfully")
	return nil
}

// waitForTask waits for a Nutanix task to complete
func waitForTask(ctx context.Context, d *schema.ResourceData, meta interface{}, taskUUID *string, timeoutType string, operation string) diag.Diagnostics {
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"QUEUED", "RUNNING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout:      d.Timeout(timeoutType),
		Delay:        ovaVMDeployDelay,
		PollInterval: ovaVMDeployDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for %s (%s): %s", operation, utils.StringValue(taskUUID), errWaitTask)
	}

	log.Printf("[DEBUG] %s completed successfully", operation)
	return nil
}

func handleDiskDeletions(ctx context.Context, d *schema.ResourceData, meta interface{}, vmmConn interface{}, deletedDisks []interface{}) diag.Diagnostics {
	if len(deletedDisks) == 0 {
		return nil
	}

	conn := vmmConn.(*conns.Client).VmmAPI
	log.Printf("[DEBUG] Handling deletion of %d disks", len(deletedDisks))
	for _, disk := range deletedDisks {
		diskInput := expandDisk([]interface{}{disk})[0]
		diskExtID := diskInput.ExtId

		readVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
		if err != nil {
			return diag.Errorf("error reading VM for disk deletion: %v", err)
		}

		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readVMResp, conn)

		resp, err := conn.VMAPIInstance.DeleteDiskById(utils.StringPtr(d.Id()), diskExtID, args)
		if err != nil {
			return diag.Errorf("error deleting disk: %v", err)
		}
		TaskRef := resp.Data.GetValue().(import3.TaskReference)
		taskUUID := TaskRef.ExtId

		if err := waitForTask(ctx, d, meta, taskUUID, schema.TimeoutUpdate, "disk deletion"); err != nil {
			return err
		}
	}
	return nil
}

func handleDiskUpdates(ctx context.Context, d *schema.ResourceData, meta interface{}, vmmConn interface{}, updatedDisks []interface{}) diag.Diagnostics {
	if len(updatedDisks) == 0 {
		return nil
	}

	conn := vmmConn.(*conns.Client).VmmAPI
	log.Printf("[DEBUG] Handling updates of %d disks", len(updatedDisks))
	for _, disk := range updatedDisks {
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

		readVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
		if err != nil {
			return diag.Errorf("error reading VM for disk update: %v", err)
		}

		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readVMResp, conn)

		resp, err := conn.VMAPIInstance.UpdateDiskById(utils.StringPtr(d.Id()), diskExtID, &diskInput, args)
		if err != nil {
			return diag.Errorf("error updating disk: %v", err)
		}
		TaskRef := resp.Data.GetValue().(import3.TaskReference)
		taskUUID := TaskRef.ExtId

		if err := waitForTask(ctx, d, meta, taskUUID, schema.TimeoutUpdate, "disk update"); err != nil {
			return err
		}
	}
	return nil
}

func handleDiskAdditions(ctx context.Context, d *schema.ResourceData, meta interface{}, vmmConn interface{}, newDisks []interface{}) diag.Diagnostics {
	if len(newDisks) == 0 {
		return nil
	}

	conn := vmmConn.(*conns.Client).VmmAPI
	log.Printf("[DEBUG] Handling addition of %d disks", len(newDisks))
	for _, disk := range newDisks {
		diskInput := expandDisk([]interface{}{disk})[0]

		readVMResp, err := conn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
		if err != nil {
			return diag.Errorf("error reading VM for disk creation: %v", err)
		}

		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readVMResp, conn)

		resp, err := conn.VMAPIInstance.CreateDisk(utils.StringPtr(d.Id()), &diskInput, args)
		if err != nil {
			return diag.Errorf("error creating disk: %v", err)
		}
		TaskRef := resp.Data.GetValue().(import3.TaskReference)
		taskUUID := TaskRef.ExtId

		if err := waitForTask(ctx, d, meta, taskUUID, schema.TimeoutUpdate, "disk creation"); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Reading VM after disk addition to update state")
	return ResourceNutanixOvaVMDeploymentRead(ctx, d, meta)
}

func handleDiskChanges(ctx context.Context, d *schema.ResourceData, meta interface{}, vmmConn interface{}, oldConfig, newConfig interface{}) diag.Diagnostics {
	var oldDisks, newDisks []interface{}

	if oldList, ok := oldConfig.([]interface{}); ok && len(oldList) > 0 {
		if oldMap, ok := oldList[0].(map[string]interface{}); ok {
			if disks, exists := oldMap["disks"]; exists && disks != nil {
				oldDisks = disks.([]interface{})
			}
		}
	}

	if newList, ok := newConfig.([]interface{}); ok && len(newList) > 0 {
		if newMap, ok := newList[0].(map[string]interface{}); ok {
			if disks, exists := newMap["disks"]; exists && disks != nil {
				newDisks = disks.([]interface{})
			}
		}
	}

	newAddedDisk, oldDeletedDisk, updatedDisk := diffConfig(oldDisks, newDisks)

	if err := handleDiskDeletions(ctx, d, meta, vmmConn, oldDeletedDisk); err != nil {
		return err
	}

	if err := handleDiskUpdates(ctx, d, meta, vmmConn, updatedDisk); err != nil {
		return err
	}

	if err := handleDiskAdditions(ctx, d, meta, vmmConn, newAddedDisk); err != nil {
		return err
	}

	return nil
}

func handlePowerStateChanges(ctx context.Context, d *schema.ResourceData, meta interface{}, conn interface{}, oldConfig, newConfig interface{}, hotPlugRequired bool) diag.Diagnostics {
	var oldPowerState, newPowerState string

	if oldList, ok := oldConfig.([]interface{}); ok && len(oldList) > 0 {
		if oldMap, ok := oldList[0].(map[string]interface{}); ok {
			if ps, exists := oldMap["power_state"]; exists && ps != nil {
				oldPowerState = ps.(string)
			}
		}
	}

	if newList, ok := newConfig.([]interface{}); ok && len(newList) > 0 {
		if newMap, ok := newList[0].(map[string]interface{}); ok {
			if ps, exists := newMap["power_state"]; exists && ps != nil {
				newPowerState = ps.(string)
			}
		}
	}

	if oldPowerState != newPowerState && newPowerState != "" && !hotPlugRequired {
		log.Printf("[DEBUG] Handling power state change from '%s' to '%s'", oldPowerState, newPowerState)

		vmmConn := conn.(*conns.Client).VmmAPI
		readResp, err := vmmConn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
		if err != nil {
			return diag.Errorf("error reading VM for power state change: %v", err)
		}

		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(readResp, vmmConn)

		var powerResp interface{}
		var taskUUID *string

		switch newPowerState {
		case "ON":
			powerResp, err = vmmConn.VMAPIInstance.PowerOnVm(utils.StringPtr(d.Id()), args)
			if err != nil {
				return diag.Errorf("error powering on VM: %v", err)
			}
			TaskRef := powerResp.(*import3.TaskReference)
			taskUUID = TaskRef.ExtId
		case "OFF":
			powerResp, err = vmmConn.VMAPIInstance.PowerOffVm(utils.StringPtr(d.Id()), args)
			if err != nil {
				return diag.Errorf("error powering off VM: %v", err)
			}
			TaskRef := powerResp.(*import3.TaskReference)
			taskUUID = TaskRef.ExtId
		default:
			return diag.Errorf("unsupported power state: %s", newPowerState)
		}

		if taskUUID != nil {
			if err := waitForTask(ctx, d, meta, taskUUID, schema.TimeoutUpdate, fmt.Sprintf("power state change to %s", newPowerState)); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleVMConfigurationUpdates(ctx context.Context, d *schema.ResourceData, meta interface{}, conn interface{}, newConfig interface{}, hasChanges bool) diag.Diagnostics {
	if !hasChanges {
		return nil
	}

	vmmConn := conn.(*conns.Client).VmmAPI
	updatedVMResp, err := vmmConn.VMAPIInstance.GetVmById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error reading VM for update: %v", err)
	}

	respVM := updatedVMResp.Data.GetValue().(import2.Vm)
	updateSpec := respVM
	checkForUpdateParams := false

	if newList, ok := newConfig.([]interface{}); ok && len(newList) > 0 {
		if overrideConfig, ok := newList[0].(map[string]interface{}); ok {
			if name, exists := overrideConfig["name"]; exists && name != nil && name.(string) != "" {
				currentName := ""
				if respVM.Name != nil {
					currentName = *respVM.Name
				}
				if name.(string) != currentName {
					updateSpec.Name = utils.StringPtr(name.(string))
					checkForUpdateParams = true
					log.Printf("[DEBUG] updating VM name from '%s' to '%s'", currentName, name.(string))
				}
			}

			if numSockets, exists := overrideConfig["num_sockets"]; exists && numSockets != nil && numSockets.(int) > 0 {
				currentSockets := 0
				if respVM.NumSockets != nil {
					currentSockets = *respVM.NumSockets
				}
				if numSockets.(int) != currentSockets {
					updateSpec.NumSockets = utils.IntPtr(numSockets.(int))
					checkForUpdateParams = true
					log.Printf("[DEBUG] updating VM sockets from %d to %d", currentSockets, numSockets.(int))
				}
			}

			if numCoresPerSocket, exists := overrideConfig["num_cores_per_socket"]; exists && numCoresPerSocket != nil && numCoresPerSocket.(int) > 0 {
				currentCores := 0
				if respVM.NumCoresPerSocket != nil {
					currentCores = *respVM.NumCoresPerSocket
				}
				if numCoresPerSocket.(int) != currentCores {
					updateSpec.NumCoresPerSocket = utils.IntPtr(numCoresPerSocket.(int))
					checkForUpdateParams = true
					log.Printf("[DEBUG] updating VM cores per socket from %d to %d", currentCores, numCoresPerSocket.(int))
				}
			}

			if numThreadsPerCore, exists := overrideConfig["num_threads_per_core"]; exists && numThreadsPerCore != nil && numThreadsPerCore.(int) > 0 {
				currentThreads := 0
				if respVM.NumThreadsPerCore != nil {
					currentThreads = *respVM.NumThreadsPerCore
				}
				if numThreadsPerCore.(int) != currentThreads {
					updateSpec.NumThreadsPerCore = utils.IntPtr(numThreadsPerCore.(int))
					checkForUpdateParams = true
					log.Printf("[DEBUG] updating VM threads per core from %d to %d", currentThreads, numThreadsPerCore.(int))
				}
			}

			if memorySizeBytes, exists := overrideConfig["memory_size_bytes"]; exists && memorySizeBytes != nil && memorySizeBytes.(int) > 0 {
				currentMemory := int64(0)
				if respVM.MemorySizeBytes != nil {
					currentMemory = *respVM.MemorySizeBytes
				}
				if int64(memorySizeBytes.(int)) != currentMemory {
					updateSpec.MemorySizeBytes = utils.Int64Ptr(int64(memorySizeBytes.(int)))
					checkForUpdateParams = true
					log.Printf("[DEBUG] updating VM memory from %d to %d", currentMemory, memorySizeBytes.(int))
				}
			}
		}
	}

	if checkForUpdateParams {
		log.Printf("[DEBUG] Applying VM configuration updates")
		args := make(map[string]interface{})
		args["If-Match"] = getEtagHeader(updatedVMResp, vmmConn)

		updateResp, err := vmmConn.VMAPIInstance.UpdateVmById(utils.StringPtr(d.Id()), &updateSpec, args)
		if err != nil {
			return diag.Errorf("error updating VM: %v", err)
		}

		TaskRef := updateResp.Data.GetValue().(import3.TaskReference)
		taskUUID := TaskRef.ExtId

		if err := waitForTask(ctx, d, meta, taskUUID, schema.TimeoutUpdate, "VM configuration update"); err != nil {
			return err
		}
	}
	return nil
}
