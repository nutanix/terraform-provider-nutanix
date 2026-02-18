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
	clustermgmtConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixClusterDiscoverUnconfiguredNodesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Create,
		ReadContext:   DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Read,
		UpdateContext: DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Update,
		DeleteContext: DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6"}, false),
			},
			"ip_filter_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     common.SchemaForIPList(false),
			},
			"uuid_filter_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"interface_filter_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_manual_discovery": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"unconfigured_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     unconfiguredNodeSchemaV2(),
			},
		},
	}
}

func unconfiguredNodeSchemaV2() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_workload": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_model_supported": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_robo_mixed_hypervisor": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"lcm_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"should_work_with_1g_nic": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"current_cvm_vlan_tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_network_interface": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cvm_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     common.SchemaForIPList(false),
			},
			"foundation_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     common.SchemaForIPList(false),
			},
			"hypervisor_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipmi_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     common.SchemaForIPList(false),
			},
			"is_secure_booted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_position": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nos_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rackable_unit_max_nodes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rackable_unit_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rackable_unit_serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	body := &config.NodeDiscoveryParams{}

	// initialize query params
	var extID *string

	if extIDf, ok := d.GetOk("ext_id"); ok {
		extID = utils.StringPtr(extIDf.(string))
	}
	if addressType, ok := d.GetOk("address_type"); ok {
		if addressType == nil || addressType == "" {
			body.AddressType = nil
		} else {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"IPV4": two,
				"IPV6": three,
			}
			pVal := subMap[addressType.(string)]
			p := config.AddressType(pVal.(int))
			body.AddressType = &p
		}
	}
	if ipFilterList, ok := d.GetOk("ip_filter_list"); ok {
		body.IpFilterList = expandIPFilterList(ipFilterList)
	}
	if uuidFilterList, ok := d.GetOk("uuid_filter_list"); ok {
		filteredUUIDList := uuidFilterList.([]interface{})
		filteredUUIDListStr := make([]string, len(filteredUUIDList))
		for i, v := range filteredUUIDList {
			filteredUUIDListStr[i] = v.(string)
		}
		body.UuidFilterList = filteredUUIDListStr
	}
	if timeout, ok := d.GetOk("timeout"); ok {
		body.Timeout = utils.Int64Ptr(int64(timeout.(int)))
	}
	if interfaceFilters, ok := d.GetOk("interface_filter_list"); ok {
		interfaceFilterList := interfaceFilters.([]interface{})
		interfaceFilterListStr := make([]string, len(interfaceFilterList))
		for i, v := range interfaceFilterList {
			interfaceFilterListStr[i] = v.(string)
		}
		body.InterfaceFilterList = interfaceFilterListStr
	}
	if isManualDiscovery, ok := d.GetOk("is_manual_discovery"); ok {
		body.IsManualDiscovery = utils.BoolPtr(isManualDiscovery.(bool))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Discover Unconfigured Nodes body : %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.DiscoverUnconfiguredNodes(extID, body)
	if err != nil {
		return diag.Errorf("error while Discover Unconfigured Nodes : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the discover unconfigured nodes operation to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for unconfigured nodes (%s) to discover: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching discover unconfigured nodes task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Discover Unconfigured Nodes Task Details: %s", string(aJSON))

	uuid := strings.Split(utils.StringValue(taskDetails.ExtId), "=:")[1]

	const unconfiguredNodes = config.TASKRESPONSETYPE_UNCONFIGURED_NODES
	taskResponseType := config.TaskResponseType(unconfiguredNodes)
	unconfiguredNodesResp, taskErr := conn.ClusterEntityAPI.FetchTaskResponse(utils.StringPtr(uuid), &taskResponseType)
	if taskErr != nil {
		return diag.Errorf("error while fetching Task Response for Unconfigured Nodes : %v", taskErr)
	}

	unconfiguredNodesTaskResp := unconfiguredNodesResp.Data.GetValue().(config.TaskResponse)

	if *unconfiguredNodesTaskResp.TaskResponseType != config.TaskResponseType(unconfiguredNodes) {
		return diag.Errorf("error while fetching Task Response for Unconfigured Nodes : %v", "task response type mismatch")
	}

	unconfiguredNodeDetails := unconfiguredNodesTaskResp.Response.GetValue().(config.UnconfigureNodeDetails)

	if err := d.Set("unconfigured_nodes", flattenUnconfiguredNodes(unconfiguredNodeDetails.NodeList)); err != nil {
		return diag.FromErr(err)
	}

	bJSON, _ := json.MarshalIndent(unconfiguredNodesResp, "", "  ")
	log.Printf("[DEBUG] Fetch Task Response for Unconfigured Nodes: %s", string(bJSON))

	// Set the ID
	d.SetId(utils.StringValue(taskDetails.ExtId))

	return nil
}

func DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixClusterV2Create(ctx, d, meta)
}

func DatasourceNutanixClusterDiscoverUnconfiguredNodesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func flattenUnconfiguredNodes(nodeListItems []config.UnconfiguredNodeListItem) []interface{} {
	if len(nodeListItems) > 0 {
		nodeList := make([]interface{}, len(nodeListItems))

		for k, v := range nodeListItems {
			node := make(map[string]interface{})
			if v.Arch != nil {
				node["arch"] = *v.Arch
			}
			if v.Attributes != nil {
				node["attributes"] = flattenAttributes(v.Attributes)
			}
			if v.ClusterId != nil {
				node["cluster_id"] = *v.ClusterId
			}
			if v.CpuType != nil {
				node["cpu_type"] = v.CpuType
			}
			if v.CurrentCvmVlanTag != nil {
				node["current_cvm_vlan_tag"] = *v.CurrentCvmVlanTag
			}
			if v.CurrentNetworkInterface != nil {
				node["current_network_interface"] = *v.CurrentNetworkInterface
			}
			if v.CvmIp != nil {
				node["cvm_ip"] = flattenIPAddress(v.CvmIp)
			}
			if v.FoundationVersion != nil {
				node["foundation_version"] = *v.FoundationVersion
			}
			if v.HostType != nil {
				node["host_type"] = flattenHostTypeEnum(v.HostType)
			}
			if v.HypervisorIp != nil {
				node["hypervisor_ip"] = flattenIPAddress(v.HypervisorIp)
			}
			if v.HypervisorType != nil {
				hypervisorTypeList := []config.HypervisorType{*v.HypervisorType}
				aJSON, _ := json.MarshalIndent(hypervisorTypeList, "", " ")
				log.Printf("[DEBUG] Hypervisor Type : %v", string(aJSON))
				aJSON, _ = json.MarshalIndent(v.HypervisorType, "", " ")
				log.Printf("[DEBUG] Hypervisor Type : %v", string(aJSON))
				node["hypervisor_type"] = flattenHypervisorType([]config.HypervisorType{*v.HypervisorType})[0]
			}
			if v.HypervisorVersion != nil {
				node["hypervisor_version"] = *v.HypervisorVersion
			}
			if v.InterfaceIpv6 != nil {
				node["interface_ipv6"] = *v.InterfaceIpv6
			}
			if v.IpmiIp != nil {
				node["ipmi_ip"] = flattenIPAddress(v.IpmiIp)
			}
			if v.NodeUuid != nil {
				node["node_uuid"] = *v.NodeUuid
			}
			if v.NosVersion != nil {
				node["nos_version"] = *v.NosVersion
			}
			if v.RackableUnitMaxNodes != nil {
				node["rackable_unit_max_nodes"] = *v.RackableUnitMaxNodes
			}
			if v.RackableUnitModel != nil {
				node["rackable_unit_model"] = *v.RackableUnitModel
			}
			if v.RackableUnitSerial != nil {
				node["rackable_unit_serial"] = *v.RackableUnitSerial
			}
			if v.IsSecureBooted != nil {
				node["is_secure_booted"] = *v.IsSecureBooted
			}
			if v.NodePosition != nil {
				node["node_position"] = *v.NodePosition
			}

			nodeList[k] = node
		}
		return nodeList
	}
	return nil
}

func flattenAttributes(attributes *config.UnconfiguredNodeAttributeMap) []interface{} {
	if attributes != nil {
		attributeMap := make([]interface{}, 1)
		attribute := make(map[string]interface{})

		if attributes.DefaultWorkload != nil {
			attribute["default_workload"] = *attributes.DefaultWorkload
		}
		if attributes.IsModelSupported != nil {
			attribute["is_model_supported"] = *attributes.IsModelSupported
		}
		if attributes.IsRoboMixedHypervisor != nil {
			attribute["is_robo_mixed_hypervisor"] = *attributes.IsRoboMixedHypervisor
		}
		if attributes.LcmFamily != nil {
			attribute["lcm_family"] = *attributes.LcmFamily
		}
		if attributes.ShouldWorkWith1GNic != nil {
			attribute["should_work_with_1g_nic"] = *attributes.ShouldWorkWith1GNic
		}

		attributeMap[0] = attribute
		return attributeMap
	}
	return nil
}

func expandIPFilterList(pr interface{}) []clustermgmtConfig.IPAddress {
	if len(pr.([]interface{})) > 0 {
		ipFilterList := make([]clustermgmtConfig.IPAddress, len(pr.([]interface{})))

		for i, v := range pr.([]interface{}) {
			ipFilter := clustermgmtConfig.IPAddress{}

			if v.(map[string]interface{})["ipv4"] != nil {
				ipFilter.Ipv4 = expandIPv4Address(v.(map[string]interface{})["ipv4"])
			}
			if v.(map[string]interface{})["ipv6"] != nil {
				ipFilter.Ipv6 = expandIPv6Address(v.(map[string]interface{})["ipv6"])
			}

			ipFilterList[i] = ipFilter
		}
		return ipFilterList
	}
	return nil
}
