package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
							Elem:     nicsElemSchemaV2(),
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
			nic := config.NewNic()

			val := v.(map[string]interface{})

			if extID, ok := val["ext_id"]; ok && len(extID.(string)) > 0 {
				nic.ExtId = utils.StringPtr(extID.(string))
			}
			// Prefer new nic_backing_info (v2.4.1+). If not present, fall back to legacy backing_info and
			// treat it as nic_backing_info.virtual_ethernet_nic to keep old configs working.
			if nbiRaw, ok := val["nic_backing_info"]; ok && nbiRaw != nil && len(nbiRaw.([]interface{})) > 0 {
				nbi := nbiRaw.([]interface{})[0].(map[string]interface{})
				nicBackingInfo := config.NewOneOfNicNicBackingInfo()

				if venRaw, ok := nbi["virtual_ethernet_nic"]; ok && venRaw != nil && len(venRaw.([]interface{})) > 0 {
					ven := expandVirtualEthernetNic(venRaw)
					if err := nicBackingInfo.SetValue(ven); err != nil {
						log.Printf("[ERROR] Error setting value for nic_backing_info.virtual_ethernet_nic: %v", err)
						diag.Errorf("Error setting value for nic_backing_info.virtual_ethernet_nic: %v", err)
						continue
					}
				} else if sriovRaw, ok := nbi["sriov_nic"]; ok && sriovRaw != nil && len(sriovRaw.([]interface{})) > 0 {
					sriov := expandSriovNic(sriovRaw)
					if err := nicBackingInfo.SetValue(sriov); err != nil {
						log.Printf("[ERROR] Error setting value for nic_backing_info.sriov_nic: %v", err)
						diag.Errorf("Error setting value for nic_backing_info.sriov_nic: %v", err)
						continue
					}
				} else if dpOffloadRaw, ok := nbi["dp_offload_nic"]; ok && dpOffloadRaw != nil && len(dpOffloadRaw.([]interface{})) > 0 {
					dpOffload := expandDpOffloadNic(dpOffloadRaw)
					if err := nicBackingInfo.SetValue(dpOffload); err != nil {
						log.Printf("[ERROR] Error setting value for nic_backing_info.dp_offload_nic: %v", err)
						diag.Errorf("Error setting value for nic_backing_info.dp_offload_nic: %v", err)
						continue
					}
				}
				// The v4 SDK provides SetNicNetworkInfo but not SetNicBackingInfo; set the oneof field directly.
				nic.NicBackingInfo = nicBackingInfo
				if nicBackingInfo != nil && nicBackingInfo.Discriminator != nil {
					if nic.NicBackingInfoItemDiscriminator_ == nil {
						nic.NicBackingInfoItemDiscriminator_ = new(string)
					}
					*nic.NicBackingInfoItemDiscriminator_ = *nicBackingInfo.Discriminator
				}
			} else if backingInfo, ok := val["backing_info"]; ok && backingInfo != nil && len(backingInfo.([]interface{})) > 0 {
				log.Printf("[DEBUG] Expanding legacy backing_info")
				nic.BackingInfo = expandEmulatedNic(backingInfo)
			}
			// Prefer new nic_network_info (v2.4.1+). If not present, fall back to legacy network_info and
			// treat it as nic_network_info.virtual_ethernet_nic_network_info to keep old configs working.
			if nniRaw, ok := val["nic_network_info"]; ok && nniRaw != nil && len(nniRaw.([]interface{})) > 0 {
				nni := nniRaw.([]interface{})[0].(map[string]interface{})
				if venNI, ok := nni["virtual_ethernet_nic_network_info"]; ok && venNI != nil && len(venNI.([]interface{})) > 0 {
					log.Printf("[DEBUG] Expanding new nic_network_info")
					ven := expandVirtualEthernetNic(venNI)
					if err := nic.SetNicNetworkInfo(ven); err != nil {
						log.Printf("[ERROR] Error setting value for nic_network_info.virtual_ethernet_nic_network_info: %v", err)
						diag.Errorf("Error setting value for nic_network_info.virtual_ethernet_nic_network_info: %v", err)
						continue
					}
				} else if sriovNI, ok := nni["sriov_nic_network_info"]; ok && sriovNI != nil && len(sriovNI.([]interface{})) > 0 {
					log.Printf("[DEBUG] Expanding new nic_network_info")
					sriov := expandSriovNic(sriovNI)
					if err := nic.SetNicNetworkInfo(sriov); err != nil {
						log.Printf("[ERROR] Error setting value for nic_network_info.sriov_nic_network_info: %v", err)
						diag.Errorf("Error setting value for nic_network_info.sriov_nic_network_info: %v", err)
						continue
					}
				} else if dpOffloadNI, ok := nni["dp_offload_nic_network_info"]; ok && dpOffloadNI != nil && len(dpOffloadNI.([]interface{})) > 0 {
					log.Printf("[DEBUG] Expanding new nic_network_info")
					dpOffload := expandDpOffloadNic(dpOffloadNI)
					if err := nic.SetNicNetworkInfo(dpOffload); err != nil {
						log.Printf("[ERROR] Error setting value for nic_network_info.dp_offload_nic_network_info: %v", err)
						diag.Errorf("Error setting value for nic_network_info.dp_offload_nic_network_info: %v", err)
						continue
					}
				}
			} else if ntwkInfo, ok := val["network_info"]; ok && ntwkInfo != nil && len(ntwkInfo.([]interface{})) > 0 {
				log.Printf("[DEBUG] Expanding legacy network_info")
				nicNetworkInfo := expandNicNetworkInfo(ntwkInfo)
				if err := nic.SetNicNetworkInfo(nicNetworkInfo); err != nil {
					log.Printf("[ERROR] Error setting value for network_info: %v", err)
					diag.Errorf("Error setting value for network_info: %v", err)
					continue
				}
			}

			nicList[k] = *nic
		}
		return nicList
	}
	return nil
}

func expandHostPcieDeviceReference(pr interface{}) *config.HostPcieDeviceReference {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	ref := config.NewHostPcieDeviceReference()
	if extID, ok := val["ext_id"]; ok && extID != nil && extID.(string) != "" {
		ref.ExtId = utils.StringPtr(extID.(string))
	}
	return ref
}

func expandNicProfileReference(pr interface{}) *config.NicProfileReference {
	if pr == nil {
		return nil
	}
	list, ok := pr.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	val, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	ref := config.NewNicProfileReference()
	if extID, ok := val["ext_id"]; ok && extID != nil && extID.(string) != "" {
		ref.ExtId = utils.StringPtr(extID.(string))
	}
	return ref
}

func expandEmulatedNic(pr interface{}) *config.EmulatedNic {
	if len(pr.([]interface{})) > 0 {
		nic := &config.EmulatedNic{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if model, ok := val["model"]; ok && len(model.(string)) > 0 {
			nic.Model = common.ExpandEnum[config.EmulatedNicModel](model.(string))
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
			nic.NicType = common.ExpandEnum[config.NicType](nicType.(string))
		}
		if ntwkFunc, ok := val["network_function_chain"]; ok {
			nic.NetworkFunctionChain = expandNetworkFunctionChainReference(ntwkFunc)
		}
		if ntwkFuncNicType, ok := val["network_function_nic_type"]; ok && len(ntwkFuncNicType.(string)) > 0 {
			nic.NetworkFunctionNicType = common.ExpandEnum[config.NetworkFunctionNicType](ntwkFuncNicType.(string))
		}
		if subnet, ok := val["subnet"]; ok {
			nic.Subnet = expandSubnetReference(subnet)
		}
		if vlanMode, ok := val["vlan_mode"]; ok && len(vlanMode.(string)) > 0 {
			nic.VlanMode = common.ExpandEnum[config.VlanMode](vlanMode.(string))
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

		// If nic_type wasn't provided, infer a sensible default.
		// This prevents perpetual diffs where a NIC in config never gets created because the API expects a nic_type.
		if nic.NicType == nil {
			// If any network-function-specific field is set, assume a network function NIC.
			if nic.NetworkFunctionNicType != nil || nic.NetworkFunctionChain != nil {
				p := config.NICTYPE_NETWORK_FUNCTION_NIC
				nic.NicType = p.Ref()
			} else if nic.Subnet != nil || nic.Ipv4Config != nil || nic.VlanMode != nil || len(nic.TrunkedVlans) > 0 || nic.ShouldAllowUnknownMacs != nil {
				p := config.NICTYPE_NORMAL_NIC
				nic.NicType = p.Ref()
			}
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
