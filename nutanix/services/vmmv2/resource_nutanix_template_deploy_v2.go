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
	import4 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	import5 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const defaultValue = 32

func ResourceNutanixTemplateDeployV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixTemplateDeployV2Create,
		ReadContext:   ResourceNutanixTemplateDeployV2Read,
		UpdateContext: ResourceNutanixTemplateDeployV2Update,
		DeleteContext: ResourceNutanixTemplateDeployV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"number_of_vms": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"override_vm_config_map": {
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
																			Default:  defaultValue,
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
						"guest_customization": schemaForGuestCustomization(),
					},
				},
			},
			"cluster_reference": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ResourceNutanixTemplateDeployV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")
	body := &import5.TemplateDeployment{}

	if versionID, ok := d.GetOk("version_id"); ok {
		body.VersionId = utils.StringPtr(versionID.(string))
	}
	if vms, ok := d.GetOk("number_of_vms"); ok {
		body.NumberOfVms = utils.IntPtr(vms.(int))
	}
	if clsRef, ok := d.GetOk("cluster_reference"); ok {
		body.ClusterReference = utils.StringPtr(clsRef.(string))
	}
	if overrideCfg, ok := d.GetOk("override_vm_config_map"); ok {
		body.OverrideVmConfigMap = expandVMConfigOverride(overrideCfg)
	}

	resp, err := conn.TemplatesAPIInstance.DeployTemplate(utils.StringPtr(extID.(string)), body)
	if err != nil {
		return diag.Errorf("error while deploying template : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the template to be deployed
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for template deploy (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching template deploy task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", " ")
	log.Printf("[DEBUG] Template Deploy Task Details: %s", string(aJSON))

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVM, "VM")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))

	return nil
}

func ResourceNutanixTemplateDeployV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixTemplateDeployV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixTemplateDeployV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandVMConfigOverride(pr interface{}) map[string]import5.VmConfigOverride {
	if len(pr.([]interface{})) > 0 {
		// vmcfg := import5.VmConfigOverride{}

		cfg := make(map[string]import5.VmConfigOverride)
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		// using vmConfig as cfg needs map[string] with single object
		vmConfig := cfg["0"]

		if name, ok := val["name"]; ok {
			vmConfig.Name = utils.StringPtr(name.(string))
		}
		if sockets, ok := val["sockets"]; ok {
			vmConfig.NumSockets = utils.IntPtr(sockets.(int))
		}
		if cores, ok := val["num_cores_per_socket"]; ok {
			vmConfig.NumCoresPerSocket = utils.IntPtr(cores.(int))
		}
		if threads, ok := val["num_threads_per_core"]; ok {
			vmConfig.NumThreadsPerCore = utils.IntPtr(threads.(int))
		}
		if mem, ok := val["memory_size_bytes"]; ok {
			vmConfig.MemorySizeBytes = utils.Int64Ptr(int64(mem.(int)))
		}
		if nics, ok := val["nics"]; ok {
			vmConfig.Nics = expandNic(nics.([]interface{}))
		}
		if guest, ok := val["guest_customization"]; ok {
			vmConfig.GuestCustomization = expandTemplateGuestCustomizationParams(guest)
		}

		cfg["0"] = vmConfig
		return cfg
	}
	return nil
}

func expandNic(pr []interface{}) []config.Nic {
	if len(pr) > 0 {
		nicList := make([]config.Nic, len(pr))

		for k, v := range pr {
			nic := config.Nic{}

			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				nic.ExtId = utils.StringPtr(extID.(string))
			}
			if backingInfo, ok := val["backing_info"]; ok && len(backingInfo.([]interface{})) > 0 {
				nic.BackingInfo = expandEmulatedNic(backingInfo)
			}
			if ntwkInfo, ok := val["network_info"]; ok && len(ntwkInfo.([]interface{})) > 0 {
				nic.NetworkInfo = expandNicNetworkInfo(ntwkInfo)
			}

			nicList[k] = nic
		}
		return nicList
	}
	return nil
}

func expandEmulatedNic(pr interface{}) *config.EmulatedNic {
	if len(pr.([]interface{})) > 0 {
		nic := &config.EmulatedNic{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if model, ok := val["model"]; ok && len(model.(string)) > 0 {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"VIRTIO": two,
				"E1000":  three,
			}
			pVal := subMap[model.(string)]
			p := config.EmulatedNicModel(pVal.(int))
			nic.Model = &p
		}
		if macAdd, ok := val["mac_address"]; ok && len(macAdd.(string)) > 0 {
			nic.MacAddress = utils.StringPtr(macAdd.(string))
		}
		if isConn, ok := val["is_connected"]; ok {
			nic.IsConnected = utils.BoolPtr(isConn.(bool))
		}
		if numQ, ok := val["num_queues"]; ok {
			nic.NumQueues = utils.IntPtr(numQ.(int))
		}
		return nic
	}
	return nil
}

func expandNicNetworkInfo(pr interface{}) *config.NicNetworkInfo {
	if len(pr.([]interface{})) > 0 {
		nic := &config.NicNetworkInfo{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if nicType, ok := val["nic_type"]; ok && len(nicType.(string)) > 0 {
			const two, three, four, five = 2, 3, 4, 5
			subMap := map[string]interface{}{
				"NORMAL_NIC":           two,
				"DIRECT_NIC":           three,
				"NETWORK_FUNCTION_NIC": four,
				"SPAN_DESTINATION_NIC": five,
			}
			pVal := subMap[nicType.(string)]
			p := config.NicType(pVal.(int))
			nic.NicType = &p
		}
		if ntwkFunc, ok := val["network_function_chain"]; ok {
			nic.NetworkFunctionChain = expandNetworkFunctionChainReference(ntwkFunc)
		}
		if ntwkFuncNicType, ok := val["network_function_nic_type"]; ok && len(ntwkFuncNicType.(string)) > 0 {
			const two, three, four = 2, 3, 4
			subMap := map[string]interface{}{
				"INGRESS": two,
				"EGRESS":  three,
				"TAP":     four,
			}
			pVal := subMap[ntwkFuncNicType.(string)]
			p := config.NetworkFunctionNicType(pVal.(int))
			nic.NetworkFunctionNicType = &p
		}
		if subnet, ok := val["subnet"]; ok {
			nic.Subnet = expandSubnetReference(subnet)
		}
		if vlanMode, ok := val["vlan_mode"]; ok && len(vlanMode.(string)) > 0 {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"ACCESS": two,
				"TRUNK":  three,
			}
			pVal := subMap[vlanMode.(string)]
			p := config.VlanMode(pVal.(int))
			nic.VlanMode = &p
		}
		if trunkedVlans, ok := val["trunked_vlans"]; ok {
			vlansList := trunkedVlans.([]interface{})
			trunkInts := make([]int, len(vlansList))

			for k, v := range vlansList {
				trunkInts[k] = v.(int)
			}
			nic.TrunkedVlans = trunkInts
		}
		if unknownMacs, ok := val["should_allow_unknown_macs"]; ok && unknownMacs.(bool) {
			nic.ShouldAllowUnknownMacs = utils.BoolPtr(unknownMacs.(bool))
		} else if unknownMacs, ok := val["should_allow_unknown_macs"]; ok && !unknownMacs.(bool) {
			nic.ShouldAllowUnknownMacs = nil
		}
		if ipv4, ok := val["ipv4_config"]; ok {
			nic.Ipv4Config = expandIpv4Config(ipv4)
		}
		if ipv4Info, ok := val["ipv4_info"]; ok {
			nic.Ipv4Info = expandIPv4Info(ipv4Info)
		}
		return nic
	}
	return nil
}

func expandIPv4Info(ipv4Info interface{}) *config.Ipv4Info {
	if len(ipv4Info.([]interface{})) > 0 {
		ipv4InfoObj := &config.Ipv4Info{}
		ipv4InfoData := ipv4Info.([]interface{})[0].(map[string]interface{})

		if learnedIPAddresses, ok := ipv4InfoData["learned_ip_addresses"]; ok {
			ipAddressesList := make([]import4.IPv4Address, len(learnedIPAddresses.([]interface{})))
			for i, learnedIP := range learnedIPAddresses.([]interface{}) {
				learnedIPData := learnedIP.(map[string]interface{})
				ipAddressesList[i] = import4.IPv4Address{
					Value:        utils.StringPtr(learnedIPData["value"].(string)),
					PrefixLength: utils.IntPtr(learnedIPData["prefix_length"].(int)),
				}
			}
			ipv4InfoObj.LearnedIpAddresses = ipAddressesList
		}
	}
	return nil
}

func expandGuestCustomizationParams(pr interface{}) *config.GuestCustomizationParams {
	if len(pr.([]interface{})) > 0 {
		guest := &config.GuestCustomizationParams{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if config, ok := val["config"]; ok {
			guest.Config = expandOneOfGuestCustomizationParamsConfig(config)
		}

		return guest
	}
	return nil
}

func expandNetworkFunctionChainReference(pr interface{}) *config.NetworkFunctionChainReference {
	if len(pr.([]interface{})) > 0 {
		ntwk := &config.NetworkFunctionChainReference{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if extID, ok := val["ext_id"]; ok {
			ntwk.ExtId = utils.StringPtr(extID.(string))
		}
		return ntwk
	}
	return nil
}

func expandSubnetReference(pr interface{}) *config.SubnetReference {
	if len(pr.([]interface{})) > 0 {
		ntwk := &config.SubnetReference{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if extID, ok := val["ext_id"]; ok {
			ntwk.ExtId = utils.StringPtr(extID.(string))
		}
		return ntwk
	}
	return nil
}

func expandIpv4Config(pr interface{}) *config.Ipv4Config {
	if len(pr.([]interface{})) > 0 {
		ipv4 := &config.Ipv4Config{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if assignIP, ok := val["should_assign_ip"]; ok {
			ipv4.ShouldAssignIp = utils.BoolPtr(assignIP.(bool))
		}
		if ipAdd, ok := val["ip_address"]; ok {
			ipv4.IpAddress = expandIPv4Address(ipAdd)
		}
		if secondaryIP, ok := val["secondary_ip_address_list"]; ok {
			ipv4.SecondaryIpAddressList = expandIPv4AddressList(secondaryIP.([]interface{}))
		}
		return ipv4
	}
	return nil
}

func expandIPv4Address(pr interface{}) *import4.IPv4Address {
	if len(pr.([]interface{})) > 0 {
		ipv4 := &import4.IPv4Address{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv4.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4
	}
	return nil
}

func expandIPv4AddressList(pr []interface{}) []import4.IPv4Address {
	if len(pr) > 0 {
		ipv4List := make([]import4.IPv4Address, len(pr))

		for k, v := range pr {
			ipv4 := import4.IPv4Address{}
			val := v.(map[string]interface{})

			if value, ok := val["value"]; ok {
				ipv4.Value = utils.StringPtr(value.(string))
			}
			if prefix, ok := val["prefix_length"]; ok {
				ipv4.PrefixLength = utils.IntPtr(prefix.(int))
			}

			ipv4List[k] = ipv4
		}
		return ipv4List
	}
	return nil
}

func expandOneOfGuestCustomizationParamsConfig(pr interface{}) *config.OneOfGuestCustomizationParamsConfig {
	if len(pr.([]interface{})) > 0 {
		guestCfgs := &config.OneOfGuestCustomizationParamsConfig{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if sysprep, ok := val["sysprep"]; ok && len(sysprep.([]interface{})) > 0 {
			sysPrepInput := config.NewSysprep()
			prI := sysprep.([]interface{})
			val := prI[0].(map[string]interface{})

			if installType, ok := val["install_type"]; ok {
				const two, three = 2, 3
				subMap := map[string]interface{}{
					"FRESH":    two,
					"PREPARED": three,
				}
				pVal := subMap[installType.(string)]
				p := config.InstallType(pVal.(int))
				sysPrepInput.InstallType = &p
			}
			if sysScript, ok := val["sysprep_script"]; ok {
				sysPrepInput.SysprepScript = expandOneOfSysprepSysprepScript(sysScript)
			}

			guestCfgs.SetValue(*sysPrepInput)
		}
		if cloudInit, ok := val["cloud_init"]; ok && len(cloudInit.([]interface{})) > 0 {
			cloud := config.NewCloudInit()
			prI := cloudInit.([]interface{})
			val := prI[0].(map[string]interface{})

			if ds, ok := val["datasource_type"]; ok && len(ds.(string)) > 0 {
				const two = 2
				subMap := map[string]interface{}{
					"CONFIG_DRIVE_V2": two,
				}
				pVal := subMap[ds.(string)]
				p := config.CloudInitDataSourceType(pVal.(int))
				cloud.DatasourceType = &p
			}
			if meta, ok := val["metadata"]; ok && len(meta.(string)) > 0 {
				cloud.Metadata = utils.StringPtr(meta.(string))
			}
			if cloudScript, ok := val["cloud_init_script"]; ok {
				cloud.CloudInitScript = expandOneOfCloudInitCloudInitScript(cloudScript)
			}
			guestCfgs.SetValue(*cloud)
		}
		return guestCfgs
	}
	return nil
}

func expandOneOfSysprepSysprepScript(pr interface{}) *config.OneOfSysprepSysprepScript {
	if len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		scripts := &config.OneOfSysprepSysprepScript{}

		if unXML, ok := val["unattend_xml"]; ok && len(unXML.([]interface{})) > 0 {
			xml := config.NewUnattendxml()
			xI := unXML.([]interface{})
			xmlVal := xI[0].(map[string]interface{})

			if vall, ok := xmlVal["value"]; ok {
				xml.Value = utils.StringPtr(vall.(string))
			}
			scripts.SetValue(*xml)
		}
		if customKeyVal, ok := val["custom_key_values"]; ok && len(customKeyVal.([]interface{})) > 0 {
			ckey := config.NewCustomKeyValues()
			cI := customKeyVal.([]interface{})
			cVal := cI[0].(map[string]interface{})

			if keyval, ok := cVal["key_value_pairs"]; ok {
				ckey.KeyValuePairs = expandTemplateKVPairs(keyval.([]interface{}))
			}
			scripts.SetValue(*ckey)
		}
		return scripts
	}
	return nil
}

func expandOneOfCloudInitCloudInitScript(pr interface{}) *config.OneOfCloudInitCloudInitScript {
	if len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		cloudInit := &config.OneOfCloudInitCloudInitScript{}

		if userdata, ok := val["user_data"]; ok && len(userdata.([]interface{})) > 0 {
			user := config.NewUserdata()
			cI := userdata.([]interface{})
			cVal := cI[0].(map[string]interface{})

			if value, ok := cVal["value"]; ok {
				user.Value = utils.StringPtr(value.(string))
			}
			cloudInit.SetValue(*user)
		}
		if customKeyVal, ok := val["custom_key_values"]; ok && len(customKeyVal.([]interface{})) > 0 {
			ckey := config.NewCustomKeyValues()
			cI := customKeyVal.([]interface{})
			cVal := cI[0].(map[string]interface{})

			if keyval, ok := cVal["key_value_pairs"]; ok {
				ckey.KeyValuePairs = expandTemplateKVPairs(keyval.([]interface{}))
			}
			cloudInit.SetValue(*ckey)
		}
		return cloudInit
	}
	return nil
}
