package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clustermgmtPrism "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixClusterAddNodeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterAddNodeV2Create,
		ReadContext:   ResourceNutanixClusterAddNodeV2Read,
		UpdateContext: ResourceNutanixClusterAddNodeV2Update,
		DeleteContext: ResourceNutanixClusterAddNodeV2Delete,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_params": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"rack_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"node_list":         AddNodeListSchema(),
						"compute_node_list": computedNodeListSchema(),
						"hypervisor_isos": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"XEN", "HYPERV", "NATIVEHOST", "ESX", "AHV"}, false),
									},
									"md5_sum": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"hyperv_sku": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"bundle_info": {
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
								},
							},
						},
						"should_skip_host_networking": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"config_params": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"should_skip_discovery": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"should_skip_imaging": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"should_validate_rack_awareness": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_nos_compatible": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_compute_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_never_schedulable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"target_hypervisor": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"hiperv": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_details":           userInfoSchema(),
									"failover_cluster_details": userInfoSchema(),
								},
							},
						},
					},
				},
			},
			"should_skip_add_node": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"should_skip_pre_expand_checks": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// remove node params
			"remove_node_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extra_params": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"should_skip_upgrade_check": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"skip_space_check": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"should_skip_add_check": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},

						"should_skip_remove": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"should_skip_prechecks": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func AddNodeListSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"node_uuid": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"block_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"node_position": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"hypervisor_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice([]string{"XEN", "HYPERV", "NATIVEHOST", "ESX", "AHV"}, false),
				},
				"is_robo_mixed_hypervisor": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"hypervisor_hostname": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"hypervisor_version": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"nos_version": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"is_light_compute": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"ipmi_ip": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem:     common.SchemaForIPList(false),
				},
				"digital_certificate_map_list": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"key": {
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
				"cvm_ip": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem:     common.SchemaForIPList(false),
				},
				"hypervisor_ip": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem:     common.SchemaForIPList(false),
				},
				"model": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"current_network_interface": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"networks": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"networks": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"uplinks": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"active":  uplinkFiledSchema(),
										"standby": uplinkFiledSchema(),
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

func computedNodeListSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"node_uuid": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"block_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"node_position": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"hypervisor_ip": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem:     common.SchemaForIPList(false),
				},
				"ipmi_ip": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem:     common.SchemaForIPList(false),
				},
				"digital_certificate_map_list": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"key": {
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
				"hypervisor_hostname": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"model": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func uplinkFiledSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mac": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func userInfoSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"username": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"password": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"cluster_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func ResourceNutanixClusterAddNodeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id")

	body := config.ExpandClusterParams{}

	if nodeParams, ok := d.GetOk("node_params"); ok {
		body.NodeParams = expandClusterNodeParams(nodeParams)
	}
	if configParams, ok := d.GetOk("config_params"); ok {
		body.ConfigParams = expandClusterConfigParams(configParams)
	}
	if skipAddNode, ok := d.GetOk("should_skip_add_node"); ok {
		body.ShouldSkipAddNode = utils.BoolPtr(skipAddNode.(bool))
	}
	if skipPreExpandChecks, ok := d.GetOk("should_skip_pre_expand_checks"); ok {
		body.ShouldSkipPreExpandChecks = utils.BoolPtr(skipPreExpandChecks.(bool))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Add Node Request Body: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.ExpandCluster(utils.StringPtr(clusterExtID.(string)), &body)
	if err != nil {
		return diag.Errorf("error while adding node : %v", err)
	}

	TaskRef := resp.Data.GetValue().(clustermgmtPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the  node to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for  node (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching  node UUID : %v", err)
	}

	aJSON, _ = json.Marshal(resourceUUID)
	log.Printf("[DEBUG] Add Node Response: %s", string(aJSON))

	uuid := clusterExtID.(string)
	d.SetId(uuid)
	return nil
}

func ResourceNutanixClusterAddNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterAddNodeV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterAddNodeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	body := &config.NodeRemovalParams{}
	clusterExtID := d.Get("cluster_ext_id")

	nodeUUIDList := make([]string, 0)

	// set node UUID
	NodeParams := d.Get("node_params").([]interface{})[0].(map[string]interface{})
	nodeListMap := NodeParams["node_list"].([]interface{})[0].(map[string]interface{})
	nodeUUID := nodeListMap["node_uuid"]
	nodeUUIDList = append(nodeUUIDList, nodeUUID.(string))

	if len(nodeUUIDList) > 0 {
		body.NodeUuids = nodeUUIDList
	} else {
		return diag.Errorf("error while removing node : Node UUID is required for remove node")
	}

	if removeNodeParams, ok := d.GetOk("remove_node_params"); ok {
		removeNodeParamsList := removeNodeParams.([]interface{})
		if len(removeNodeParamsList) > 0 {
			removeNodeParamsMap := removeNodeParamsList[0].(map[string]interface{})
			if extraParams, ok := removeNodeParamsMap["extra_params"]; ok {
				body.ExtraParams = expandExtraParams(extraParams)
			}
			if skipRemove, ok := removeNodeParamsMap["should_skip_remove"]; ok {
				body.ShouldSkipRemove = utils.BoolPtr(skipRemove.(bool))
			}
			if skipPrechecks, ok := removeNodeParamsMap["should_skip_prechecks"]; ok {
				body.ShouldSkipPrechecks = utils.BoolPtr(skipPrechecks.(bool))
			}
		}
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Remove Node Request Body: %s", string(aJSON))
	resp, err := conn.ClusterEntityAPI.RemoveNode(utils.StringPtr(clusterExtID.(string)), body)
	if err != nil {
		return diag.Errorf("error while Removing node : %v", err)
	}

	TaskRef := resp.Data.GetValue().(clustermgmtPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the node to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		resourceUUID, _ := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		rUUID := resourceUUID.Data.GetValue().(import2.Task)
		aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
		log.Printf("Error Remove Node Task Details : %s", string(aJSON))
		return diag.Errorf("error waiting for  node (%s) to Remove: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching  node UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)

	bJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("Remove Node Task Details : %s", string(bJSON))
	return nil
}

func expandClusterNodeParams(pr interface{}) *config.NodeParam {
	if pr != nil {
		nConf := config.NodeParam{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if blockList, ok := val["block_list"]; ok {
			nConf.BlockList = expandBlockList(blockList.([]interface{}))
		}
		if nodeList, ok := val["node_list"]; ok {
			nConf.NodeList = expandNodeList(nodeList.([]interface{}))
		}
		if computeNodeList, ok := val["compute_node_list"]; ok {
			nConf.ComputeNodeList = expandComputeNodeList(computeNodeList.([]interface{}))
		}
		if hypervisorIsos, ok := val["hypervisor_isos"]; ok {
			nConf.HypervisorIsos = expandHypervisorIsos(hypervisorIsos.([]interface{}))
		}
		if hypervSku, ok := val["hyperv_sku"]; ok && hypervSku != "" {
			nConf.HypervSku = utils.StringPtr(hypervSku.(string))
		}
		if bundleInfo, ok := val["bundle_info"]; ok {
			nConf.BundleInfo = expandBundleInfo(bundleInfo)
		}
		if skipHostNetworking, ok := val["should_skip_host_networking"]; ok {
			nConf.ShouldSkipHostNetworking = utils.BoolPtr(skipHostNetworking.(bool))
		}

		return &nConf
	}
	return nil
}

func expandBlockList(pr []interface{}) []config.BlockItem {
	if len(pr) > 0 {
		blockList := make([]config.BlockItem, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			block := config.BlockItem{}

			if blockID, ok := val["block_id"]; ok && blockID != "" {
				block.BlockId = utils.StringPtr(blockID.(string))
			}
			if rackName, ok := val["rack_name"]; ok && rackName != "" {
				block.RackName = utils.StringPtr(rackName.(string))
			}
			blockList[k] = block
		}
		return blockList
	}
	return nil
}

func expandNodeList(pr []interface{}) []config.NodeItem {
	if len(pr) > 0 {
		nodeList := make([]config.NodeItem, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			node := config.NodeItem{}

			if nodeUUID, ok := val["node_uuid"]; ok && nodeUUID != "" {
				node.NodeUuid = utils.StringPtr(nodeUUID.(string))
			}
			if blockID, ok := val["block_id"]; ok && blockID != "" {
				node.BlockId = utils.StringPtr(blockID.(string))
			}
			if nodePosition, ok := val["node_position"]; ok && nodePosition != "" {
				node.NodePosition = utils.StringPtr(nodePosition.(string))
			}
			if hypervisorType, ok := val["hypervisor_type"]; ok {
				node.HypervisorType = expandHypervisorType(hypervisorType)
			}
			if roboMixedHypervisor, ok := val["is_robo_mixed_hypervisor"]; ok {
				node.IsRoboMixedHypervisor = utils.BoolPtr(roboMixedHypervisor.(bool))
			}
			if hypervisorHostname, ok := val["hypervisor_hostname"]; ok && hypervisorHostname != "" {
				node.HypervisorHostname = utils.StringPtr(hypervisorHostname.(string))
			}
			if hypervisorVersion, ok := val["hypervisor_version"]; ok && hypervisorVersion != "" {
				node.HypervisorVersion = utils.StringPtr(hypervisorVersion.(string))
			}
			if nosVersion, ok := val["nos_version"]; ok && nosVersion != "" {
				node.NosVersion = utils.StringPtr(nosVersion.(string))
			}
			if isLightCompute, ok := val["is_light_compute"]; ok {
				node.IsLightCompute = utils.BoolPtr(isLightCompute.(bool))
			}
			if ipmiIP, ok := val["ipmi_ip"]; ok {
				node.IpmiIp = expandIPAddress(ipmiIP)
			}
			if digitalCertificateMapList, ok := val["digital_certificate_map_list"]; ok {
				node.DigitalCertificateMapList = expandKeyValueMap(digitalCertificateMapList.([]interface{}))
			}
			if cvmIP, ok := val["cvm_ip"]; ok {
				node.CvmIp = expandIPAddress(cvmIP)
			}
			if hypervisorIP, ok := val["hypervisor_ip"]; ok {
				node.HypervisorIp = expandIPAddress(hypervisorIP)
			}
			if model, ok := val["model"]; ok {
				node.Model = utils.StringPtr(model.(string))
			}
			if currentNetworkInterface, ok := val["current_network_interface"]; ok {
				node.CurrentNetworkInterface = utils.StringPtr(currentNetworkInterface.(string))
			}
			if networks, ok := val["networks"]; ok {
				node.Networks = expandNetworks(networks.([]interface{}))
			}
			nodeList[k] = node
		}
		return nodeList
	}
	return nil
}

func expandComputeNodeList(pr []interface{}) []config.ComputeNodeItem {
	if len(pr) > 0 {
		nodeList := make([]config.ComputeNodeItem, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			node := config.NewComputeNodeItem()

			if nodeUUID, ok := val["node_uuid"]; ok && nodeUUID != "" {
				node.NodeUuid = utils.StringPtr(nodeUUID.(string))
			}
			if blockID, ok := val["block_id"]; ok && blockID != "" {
				node.BlockId = utils.StringPtr(blockID.(string))
			}
			if nodePosition, ok := val["node_position"]; ok && nodePosition != "" {
				node.NodePosition = utils.StringPtr(nodePosition.(string))
			}
			if hypervisorHostname, ok := val["hypervisor_hostname"]; ok && hypervisorHostname != "" {
				node.HypervisorHostname = utils.StringPtr(hypervisorHostname.(string))
			}
			if ipmiIP, ok := val["ipmi_ip"]; ok {
				node.IpmiIp = expandIPAddress(ipmiIP)
			}
			if digitalCertificateMapList, ok := val["digital_certificate_map_list"]; ok {
				node.DigitalCertificateMapList = expandKeyValueMap(digitalCertificateMapList.([]interface{}))
			}
			if hypervisorIP, ok := val["hypervisor_ip"]; ok {
				node.HypervisorIp = expandIPAddress(hypervisorIP)
			}
			if model, ok := val["model"]; ok && model != "" {
				node.Model = utils.StringPtr(model.(string))
			}
			nodeList[k] = *node
		}
		return nodeList
	}
	return nil
}

func expandKeyValueMap(pr []interface{}) []config.DigitalCertificateMapReference {
	if len(pr) > 0 {
		dcmList := make([]config.DigitalCertificateMapReference, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			dcm := config.DigitalCertificateMapReference{}

			if key, ok := val["key"]; ok {
				dcm.Key = utils.StringPtr(key.(string))
			}
			if value, ok := val["value"]; ok {
				dcm.Value = utils.StringPtr(value.(string))
			}
			dcmList[k] = dcm
		}
		return dcmList
	}
	return nil
}

func expandNetworks(pr []interface{}) []config.UplinkNetworkItem {
	if len(pr) > 0 {
		networkList := make([]config.UplinkNetworkItem, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			network := config.UplinkNetworkItem{}

			if name, ok := val["name"]; ok {
				network.Name = utils.StringPtr(name.(string))
			}
			if networkTypes, ok := val["networks"]; ok {
				networksList := networkTypes.([]interface{})
				networksListStr := make([]string, len(networksList))
				for i, v := range networksList {
					networksListStr[i] = v.(string)
				}
				network.Networks = networksListStr
			}
			if uplinks, ok := val["uplinks"]; ok {
				network.Uplinks = expandUplink(uplinks)
			}
			networkList[k] = network
		}
		return networkList
	}
	return nil
}

func expandUplink(pr interface{}) *config.Uplinks {
	if pr != nil {
		nConf := config.Uplinks{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if active, ok := val["active"]; ok {
			nConf.Active = expandUplinkParams(active.([]interface{}))
		}
		if standby, ok := val["standby"]; ok {
			nConf.Standby = expandUplinkParams(standby.([]interface{}))
		}

		return &nConf
	}
	return nil
}

func expandUplinkParams(pr []interface{}) []config.UplinksField {
	if len(pr) > 0 {
		networkList := make([]config.UplinksField, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			network := config.UplinksField{}

			if mac, ok := val["mac"]; ok {
				network.Mac = utils.StringPtr(mac.(string))
			}
			if name, ok := val["name"]; ok {
				network.Name = utils.StringPtr(name.(string))
			}
			if value, ok := val["value"]; ok {
				network.Value = utils.StringPtr(value.(string))
			}
			networkList[k] = network
		}
		return networkList
	}
	return nil
}

func expandHypervisorIsos(pr []interface{}) []config.HypervisorIsoMap {
	aJSON, _ := json.MarshalIndent(pr, "", " ")
	log.Printf("[DEBUG] expandHypervisorIsos pr: %s", string(aJSON))
	if len(pr) > 0 {
		itemList := make([]config.HypervisorIsoMap, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			item := config.NewHypervisorIsoMap()

			if md5Sum, ok := val["md5_sum"]; ok && md5Sum != "" {
				item.Md5Sum = utils.StringPtr(md5Sum.(string))
			}
			if hypervisorType, ok := val["type"]; ok {
				item.Type = expandHypervisorType(hypervisorType)
			}

			itemList[k] = *item
		}
		return itemList
	}
	return nil
}

func expandBundleInfo(pr interface{}) *config.BundleInfo {
	if pr != nil {
		if len(pr.([]interface{})) == 0 {
			return nil
		}

		nConf := config.BundleInfo{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if name, ok := val["name"]; ok {
			nConf.Name = utils.StringPtr(name.(string))
		}
		return &nConf
	}
	return nil
}

func expandClusterConfigParams(pr interface{}) *config.ConfigParams {
	if pr != nil {
		if len(pr.([]interface{})) == 0 {
			return nil
		}

		conf := config.ConfigParams{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if skipDiscovery, ok := val["should_skip_discovery"]; ok {
			conf.ShouldSkipDiscovery = utils.BoolPtr(skipDiscovery.(bool))
		}
		if skipImaging, ok := val["should_skip_imaging"]; ok {
			conf.ShouldSkipImaging = utils.BoolPtr(skipImaging.(bool))
		}
		if validateRackAwareness, ok := val["should_validate_rack_awareness"]; ok {
			conf.ShouldValidateRackAwareness = utils.BoolPtr(validateRackAwareness.(bool))
		}
		if isNosCompatible, ok := val["is_nos_compatible"]; ok {
			conf.IsNosCompatible = utils.BoolPtr(isNosCompatible.(bool))
		}
		if isComputeOnly, ok := val["is_compute_only"]; ok {
			conf.IsComputeOnly = utils.BoolPtr(isComputeOnly.(bool))
		}
		if neverSchedulable, ok := val["is_never_schedulable"]; ok {
			conf.IsNeverScheduleable = utils.BoolPtr(neverSchedulable.(bool))
		}
		if targetHypervisor, ok := val["target_hypervisor"]; ok {
			conf.TargetHypervisor = utils.StringPtr(targetHypervisor.(string))
		}
		if hyperv, ok := val["hyperv"]; ok {
			conf.Hyperv = expandHyperv(hyperv)
		}

		return &conf
	}

	return nil
}

func expandHyperv(pr interface{}) *config.HypervCredentials {
	if pr != nil {
		if len(pr.([]interface{})) == 0 {
			return nil
		}

		hyperv := config.HypervCredentials{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if domainDetails, ok := val["domain_details"]; ok {
			hyperv.DomainDetails = expandDetails(domainDetails)
		}
		if failoverClusterDetails, ok := val["failover_cluster_details"]; ok {
			hyperv.FailoverClusterDetails = expandDetails(failoverClusterDetails)
		}
		return &hyperv
	}
	return nil
}

func expandDetails(pr interface{}) *config.UserInfo {
	if pr != nil {
		if len(pr.([]interface{})) == 0 {
			return nil
		}

		userInfo := config.UserInfo{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if username, ok := val["username"]; ok {
			userInfo.UserName = utils.StringPtr(username.(string))
		}
		if password, ok := val["password"]; ok {
			userInfo.Password = utils.StringPtr(password.(string))
		}
		if clusterName, ok := val["cluster_name"]; ok {
			userInfo.ClusterName = utils.StringPtr(clusterName.(string))
		}
		return &userInfo
	}
	return nil
}

func expandHypervisorType(hypervisorType interface{}) *config.HypervisorType {
	return common.ExpandEnum(hypervisorType, HypervisorTypeMap, "hypervisor_type")
}

func expandExtraParams(pr interface{}) *config.NodeRemovalExtraParam {
	if pr != nil {
		extraParams := config.NodeRemovalExtraParam{}
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if skipUpgradeCheck, ok := val["should_skip_upgrade_check"]; ok {
			extraParams.ShouldSkipUpgradeCheck = utils.BoolPtr(skipUpgradeCheck.(bool))
		}
		if skipSpaceCheck, ok := val["skip_space_check"]; ok {
			extraParams.ShouldSkipSpaceCheck = utils.BoolPtr(skipSpaceCheck.(bool))
		}
		if skipAddCheck, ok := val["should_skip_add_check"]; ok {
			extraParams.ShouldSkipAddCheck = utils.BoolPtr(skipAddCheck.(bool))
		}
		return &extraParams
	}
	return nil
}
