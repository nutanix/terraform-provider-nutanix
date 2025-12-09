package clustersv2

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clsMangPrismConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixClusterUnconfiguredNodeNetworkV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterUnconfiguredNodeNetworkV2Create,
		ReadContext:   ResourceNutanixClusterUnconfiguredNodeNetworkV2Read,
		UpdateContext: ResourceNutanixClusterUnconfiguredNodeNetworkV2Update,
		DeleteContext: ResourceNutanixClusterUnconfiguredNodeNetworkV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_list": nodeListSchema(),
			"request_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nodes_networking_details": nodeListNetworkingDetailsSchema(),
		},
	}
}

func nodeListNetworkingDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"network_info": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"hci": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"hypervisor_type": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"networks": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
							"so": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"hypervisor_type": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"networks": {
											Type:     schema.TypeList,
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
				"uplinks": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.
							Schema{
							"cvm_ip": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     common.SchemaForIPList(false),
							},
							"uplink_list": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"mac": {
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
				"warnings": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func nodeListSchema() *schema.Schema {
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
				"model": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"is_compute_only": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"is_light_compute": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"hypervisor_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice(HypervisorTypeStrings, false),
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
				"current_network_interface": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"is_robo_mixed_hypervisor": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func ResourceNutanixClusterUnconfiguredNodeNetworkV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	body := &config.NodeDetails{}
	clusterExtID := d.Get("ext_id")
	var expand *string

	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}
	if nodeList, ok := d.GetOk("node_list"); ok {
		body.NodeList = expandNodeListNetworkingDetails(nodeList.([]interface{}))
	}
	if requestType, ok := d.GetOk("request_type"); ok {
		body.RequestType = utils.StringPtr(requestType.(string))
	}

	readResp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(clusterExtID.(string)), expand)
	if err != nil {
		return diag.Errorf("error while reading cluster : %v", err)
	}
	// Extract E-Tag Header
	args := getEtagHeader(readResp, conn)

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Fetch Network info Request Body: %s\n", string(aJSON))

	resp, err := conn.ClusterEntityAPI.FetchNodeNetworkingDetails(utils.StringPtr(clusterExtID.(string)), body, args)
	if err != nil {
		return diag.Errorf("error while Fetching Node Networking Details : %v", err)
	}

	TaskRef := resp.Data.GetValue().(clsMangPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the  node to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for  node (%s) to add: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching task : %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	bJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Fetch Network Info Task Details: %s", string(bJSON))

	uuid := strings.Split(utils.StringValue(taskDetails.ExtId), "=:")[1]

	networkingDetails := config.TASKRESPONSETYPE_NETWORKING_DETAILS
	taskResponseType := config.TaskResponseType(networkingDetails)
	networkDetailsResp, taskErr := conn.ClusterEntityAPI.FetchTaskResponse(utils.StringPtr(uuid), &taskResponseType)
	if taskErr != nil {
		return diag.Errorf("error while fetching Task Response for Unconfigured Nodes : %v", taskErr)
	}

	networkDetailsTaskResp := networkDetailsResp.Data.GetValue().(config.TaskResponse)

	if *networkDetailsTaskResp.TaskResponseType != config.TaskResponseType(networkingDetails) {
		return diag.Errorf("error while fetching Task Response for Network Detail Nodes : %v", "task response type mismatch")
	}

	nodeNetworkDetails := networkDetailsTaskResp.Response.GetValue().(config.NodeNetworkingDetails)

	if err := d.Set("nodes_networking_details", flattenNodesNetworkDetails(nodeNetworkDetails)); err != nil {
		return diag.FromErr(err)
	}

	aJSON, _ = json.MarshalIndent(networkDetailsResp, "", " ")
	log.Printf("[DEBUG] fetching Task Response for Unconfigured Nodes Task Details: %s\n", string(aJSON))

	d.SetId(utils.StringValue(taskDetails.ExtId))
	return nil
}

func flattenNodesNetworkDetails(nodeNetworkDetails config.NodeNetworkingDetails) []map[string]interface{} {
	if nodeNetworkDetails.NetworkInfo != nil {
		result := make(map[string]interface{})

		networkInfo := nodeNetworkDetails.NetworkInfo
		uplinks := nodeNetworkDetails.Uplinks
		warnings := nodeNetworkDetails.Warnings

		if networkInfo != nil {
			networkInfoList := make([]map[string]interface{}, 0)
			networkInfoMap := make(map[string]interface{})

			if networkInfo.Hci != nil {
				networkInfoMap["hci"] = flattenNameNetworkRef(networkInfo.Hci)
			}
			if networkInfo.So != nil {
				networkInfoMap["so"] = flattenNameNetworkRef(networkInfo.So)
			}
			networkInfoList = append(networkInfoList, networkInfoMap)
			result["network_info"] = networkInfoList
		}

		if len(uplinks) > 0 {
			uplinksList := make([]map[string]interface{}, 0)
			for _, uplink := range uplinks {
				uplinksMap := make(map[string]interface{})
				uplinksMap["cvm_ip"] = flattenIPAddress(uplink.CvmIp)
				uplinksMap["uplink_list"] = flattenUplinkList(uplink.UplinkList)
				uplinksList = append(uplinksList, uplinksMap)
			}
			result["uplinks"] = uplinksList
		}

		if len(warnings) > 0 {
			warningsList := make([]string, 0)
			warningsList = append(warningsList, warnings...)
			result["warnings"] = warningsList
		}
		return []map[string]interface{}{result}
	}
	return nil
}

func flattenUplinkList(uplinkList []config.NameMacRef) interface{} {
	if len(uplinkList) > 0 {
		result := make([]map[string]interface{}, 0)
		for _, uplink := range uplinkList {
			uplinkMap := make(map[string]interface{})
			uplinkMap["mac"] = uplink.Mac
			uplinkMap["name"] = uplink.Name
			result = append(result, uplinkMap)
		}
		return result
	}
	return nil
}

func flattenNameNetworkRef(nameNetworkRefs []config.NameNetworkRef) []map[string]interface{} {
	if nameNetworkRefs != nil {
		result := make([]map[string]interface{}, 0)
		for _, nameNetworkRef := range nameNetworkRefs {
			networkMap := make(map[string]interface{})
			networkMap["hypervisor_type"] = flattenHypervisorType([]config.HypervisorType{*nameNetworkRef.HypervisorType})[0]
			networkMap["name"] = nameNetworkRef.Name
			networkMap["networks"] = nameNetworkRef.Networks
			result = append(result, networkMap)
		}
		return result
	}
	return nil
}

func ResourceNutanixClusterUnconfiguredNodeNetworkV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterUnconfiguredNodeNetworkV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixClusterV2Create(ctx, d, meta)
}

func ResourceNutanixClusterUnconfiguredNodeNetworkV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandNodeListNetworkingDetails(pr []interface{}) []config.NodeListNetworkingDetails {
	if len(pr) > 0 {
		nodeList := make([]config.NodeListNetworkingDetails, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			node := config.NodeListNetworkingDetails{}

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
				node.HypervisorType = common.ExpandEnum(hypervisorType, HypervisorTypeMap, "hypervisor_type")
			}
			if roboMixedHypervisor, ok := val["is_robo_mixed_hypervisor"]; ok {
				node.IsRoboMixedHypervisor = utils.BoolPtr(roboMixedHypervisor.(bool))
			}
			if hypervisorVersion, ok := val["hypervisor_version"]; ok && hypervisorVersion != "" {
				node.HypervisorVersion = utils.StringPtr(hypervisorVersion.(string))
			}
			if nosVersion, ok := val["nos_version"]; ok && nosVersion != "" {
				node.NosVersion = utils.StringPtr(nosVersion.(string))
			}
			if isLightCompute, ok := val["is_light_compute"]; ok && isLightCompute != "" {
				node.IsLightCompute = utils.BoolPtr(isLightCompute.(bool))
			}
			if isComputeOnly, ok := val["is_compute_only"]; ok {
				node.IsComputeOnly = utils.BoolPtr(isComputeOnly.(bool))
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
			if model, ok := val["model"]; ok && model != "" {
				node.Model = utils.StringPtr(model.(string))
			}
			if currentNetworkInterface, ok := val["current_network_interface"]; ok && currentNetworkInterface != "" {
				node.CurrentNetworkInterface = utils.StringPtr(currentNetworkInterface.(string))
			}
			nodeList[k] = node
		}
		return nodeList
	}
	return nil
}
