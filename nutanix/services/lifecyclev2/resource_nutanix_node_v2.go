package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import3 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	prismConfigLc "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixNodeV2 registers and onboards a node for management by Foundation Central.
func ResourceNutanixNodeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixNodeV2Create,
		ReadContext:   ResourceNutanixNodeV2Read,
		UpdateContext: ResourceNutanixNodeV2Update,
		DeleteContext: ResourceNutanixNodeV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identifiers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"manufacturer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aos_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_serial_number": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"memory_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"custom_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_connection_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cpu_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity_g_hz": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
						"logical_core_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"manufacturer": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"socket_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"network_details": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bmc":  nodeNetworkDetailsSchema(),
						"cvm":  nodeComponentNetworkDetailsSchema(),
						"host": nodeComponentNetworkDetailsSchema(),
					},
				},
			},
			// Read-only projection data.
			"cvm_connectivity_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_connectivity_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_data": nodeProviderDataSchema(),
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
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixNodeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	spec := expandNode(d)

	aJSON, _ := json.MarshalIndent(spec, "", " ")
	log.Printf("[DEBUG] Node create payload : %s", string(aJSON))

	resp, err := conn.NodesAPIInstance.CreateNode(spec)
	if err != nil {
		return diag.Errorf("error while creating node : %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfigLc.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for node (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching node creation task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeNode, "Node")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixNodeV2Read(ctx, d, meta)
}

func ResourceNutanixNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.NodesAPIInstance.GetNodeById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching node : %v", err)
	}

	return flattenNodeToState(d, resp.Data.GetValue())
}

func ResourceNutanixNodeV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.NodesAPIInstance.GetNodeById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching node : %v", err)
	}

	updateSpec := getNodeFromResponse(resp.Data.GetValue())
	if updateSpec == nil {
		return diag.Errorf("error while reading node for update: unexpected response type")
	}

	if d.HasChange("manufacturer") {
		updateSpec.Manufacturer = utils.StringPtr(d.Get("manufacturer").(string))
	}
	if d.HasChange("model") {
		updateSpec.Model = utils.StringPtr(d.Get("model").(string))
	}
	if d.HasChange("hostname") {
		updateSpec.Hostname = utils.StringPtr(d.Get("hostname").(string))
	}
	if d.HasChange("host_version") {
		updateSpec.HostVersion = utils.StringPtr(d.Get("host_version").(string))
	}
	if d.HasChange("aos_version") {
		updateSpec.AosVersion = utils.StringPtr(d.Get("aos_version").(string))
	}
	if d.HasChange("block_serial_number") {
		updateSpec.BlockSerialNumber = utils.StringPtr(d.Get("block_serial_number").(string))
	}
	if d.HasChange("memory_gb") {
		updateSpec.MemoryGB = utils.IntPtr(d.Get("memory_gb").(int))
	}
	if d.HasChange("host_type") {
		updateSpec.HostType = common.ExpandEnum[import3.HostType](d.Get("host_type").(string))
	}
	if d.HasChange("custom_attributes") {
		updateSpec.CustomAttributes = common.ExpandListOfString(d.Get("custom_attributes").([]interface{}))
	}
	if d.HasChange("owner_ext_id") {
		updateSpec.OwnerExtId = utils.StringPtr(d.Get("owner_ext_id").(string))
	}
	if d.HasChange("provider_ext_id") {
		updateSpec.ProviderExtId = utils.StringPtr(d.Get("provider_ext_id").(string))
	}
	if d.HasChange("provider_connection_ext_id") {
		updateSpec.ProviderConnectionExtId = utils.StringPtr(d.Get("provider_connection_ext_id").(string))
	}
	if d.HasChange("identifiers") {
		updateSpec.Identifiers = expandNodeIdentifiers(d.Get("identifiers").([]interface{}))
	}
	if d.HasChange("cpu_info") {
		updateSpec.CpuInfo = expandCPUInfo(d.Get("cpu_info").([]interface{}))
	}
	if d.HasChange("network_details") {
		updateSpec.NetworkDetails = expandNodeNetworkDetails(d.Get("network_details").([]interface{}))
	}

	updateResp, err := conn.NodesAPIInstance.UpdateNodeById(utils.StringPtr(d.Id()), updateSpec)
	if err != nil {
		return diag.Errorf("error while updating node : %v", err)
	}

	taskRef := updateResp.Data.GetValue().(prismConfigLc.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for node (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return ResourceNutanixNodeV2Read(ctx, d, meta)
}

func ResourceNutanixNodeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	resp, err := conn.NodesAPIInstance.DeleteNodeById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting node : %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfigLc.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for node (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	d.SetId("")
	return nil
}

// expandNode builds the Node request body from the resource data.
func expandNode(d *schema.ResourceData) *import3.Node {
	spec := &import3.Node{}

	if v, ok := d.GetOk("identifiers"); ok {
		spec.Identifiers = expandNodeIdentifiers(v.([]interface{}))
	}
	if v, ok := d.GetOk("manufacturer"); ok {
		spec.Manufacturer = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("model"); ok {
		spec.Model = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("hostname"); ok {
		spec.Hostname = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("host_type"); ok {
		spec.HostType = common.ExpandEnum[import3.HostType](v.(string))
	}
	if v, ok := d.GetOk("host_version"); ok {
		spec.HostVersion = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("aos_version"); ok {
		spec.AosVersion = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("block_serial_number"); ok {
		spec.BlockSerialNumber = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("memory_gb"); ok {
		spec.MemoryGB = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("custom_attributes"); ok {
		spec.CustomAttributes = common.ExpandListOfString(v.([]interface{}))
	}
	if v, ok := d.GetOk("owner_ext_id"); ok {
		spec.OwnerExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("provider_ext_id"); ok {
		spec.ProviderExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("provider_connection_ext_id"); ok {
		spec.ProviderConnectionExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("cpu_info"); ok {
		spec.CpuInfo = expandCPUInfo(v.([]interface{}))
	}
	if v, ok := d.GetOk("network_details"); ok {
		spec.NetworkDetails = expandNodeNetworkDetails(v.([]interface{}))
	}
	return spec
}

func expandNodeIdentifiers(pr []interface{}) []import3.NodeIdentifier {
	if len(pr) == 0 {
		return nil
	}
	out := make([]import3.NodeIdentifier, 0, len(pr))
	for _, item := range pr {
		if item == nil {
			continue
		}
		val := item.(map[string]interface{})
		identifier := import3.NodeIdentifier{}
		if t, ok := val["type"]; ok && t.(string) != "" {
			identifier.Type = common.ExpandEnum[import3.NodeIdentifierType](t.(string))
		}
		if v, ok := val["value"]; ok && v.(string) != "" {
			identifier.Value = utils.StringPtr(v.(string))
		}
		out = append(out, identifier)
	}
	return out
}

func expandCPUInfo(pr []interface{}) *import3.CpuInfo {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	cpu := &import3.CpuInfo{}
	if v, ok := val["capacity_g_hz"]; ok && v.(float64) != 0 {
		f := float32(v.(float64))
		cpu.CapacityGHz = &f
	}
	if v, ok := val["logical_core_count"]; ok && v.(int) != 0 {
		cpu.LogicalCoreCount = utils.IntPtr(v.(int))
	}
	if v, ok := val["manufacturer"]; ok && v.(string) != "" {
		cpu.Manufacturer = utils.StringPtr(v.(string))
	}
	if v, ok := val["model"]; ok && v.(string) != "" {
		cpu.Model = utils.StringPtr(v.(string))
	}
	if v, ok := val["socket_count"]; ok && v.(int) != 0 {
		cpu.SocketCount = utils.IntPtr(v.(int))
	}
	return cpu
}

func expandNodeNetworkDetails(pr []interface{}) *import3.NodeNetworkDetails {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	nd := &import3.NodeNetworkDetails{}
	if v, ok := val["bmc"]; ok {
		nd.Bmc = expandNetworkDetails(v.([]interface{}))
	}
	if v, ok := val["cvm"]; ok {
		nd.Cvm = expandComponentNetworkDetails(v.([]interface{}))
	}
	if v, ok := val["host"]; ok {
		nd.Host = expandComponentNetworkDetails(v.([]interface{}))
	}
	return nd
}

func expandNetworkDetails(pr []interface{}) *import3.NetworkDetails {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	nd := &import3.NetworkDetails{}
	if v, ok := val["gateway"]; ok {
		nd.Gateway = expandIPAddress(v.([]interface{}))
	}
	if v, ok := val["ip"]; ok {
		nd.Ip = expandIPAddress(v.([]interface{}))
	}
	if v, ok := val["vlan_id"]; ok && v.(int) != 0 {
		nd.VlanId = utils.IntPtr(v.(int))
	}
	return nd
}

func expandComponentNetworkDetails(pr []interface{}) *import3.ComponentNetworkDetails {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	cnd := &import3.ComponentNetworkDetails{}
	if v, ok := val["management_network"]; ok {
		cnd.ManagementNetwork = expandNetworkDetails(v.([]interface{}))
	}
	return cnd
}

// getNodeFromResponse returns a *Node from either a Node or NodeProjection GetById OneOf response.
func getNodeFromResponse(value interface{}) *import3.Node {
	switch v := value.(type) {
	case import3.Node:
		return &v
	case *import3.Node:
		return v
	case import3.NodeProjection:
		return nodeFromProjection(&v)
	case *import3.NodeProjection:
		return nodeFromProjection(v)
	default:
		return nil
	}
}

// nodeFromProjection maps the shared fields of a NodeProjection onto a Node, so the
// read-modify-write Update flow works regardless of which OneOf branch the API returns.
func nodeFromProjection(p *import3.NodeProjection) *import3.Node {
	if p == nil {
		return nil
	}
	return &import3.Node{
		AosVersion:              p.AosVersion,
		BlockSerialNumber:       p.BlockSerialNumber,
		CpuInfo:                 p.CpuInfo,
		CreatedTime:             p.CreatedTime,
		CustomAttributes:        p.CustomAttributes,
		CvmConnectivityStatus:   p.CvmConnectivityStatus,
		ExtId:                   p.ExtId,
		HostConnectivityStatus:  p.HostConnectivityStatus,
		HostType:                p.HostType,
		HostVersion:             p.HostVersion,
		Hostname:                p.Hostname,
		Identifiers:             p.Identifiers,
		Links:                   p.Links,
		Manufacturer:            p.Manufacturer,
		MemoryGB:                p.MemoryGB,
		Model:                   p.Model,
		NetworkDetails:          p.NetworkDetails,
		OwnerExtId:              p.OwnerExtId,
		ProviderConnectionExtId: p.ProviderConnectionExtId,
		ProviderData:            p.ProviderData,
		ProviderExtId:           p.ProviderExtId,
		State:                   p.State,
		TenantId:                p.TenantId,
	}
}

// flattenNodeToState writes a Node/NodeProjection GetById response onto the resource state.
func flattenNodeToState(d *schema.ResourceData, value interface{}) diag.Diagnostics {
	node := getNodeFromResponse(value)
	if node == nil {
		return diag.Errorf("error while reading node: unexpected response type")
	}

	if err := d.Set("ext_id", utils.StringValue(node.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("identifiers", flattenNodeIdentifiers(node.Identifiers)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("manufacturer", utils.StringValue(node.Manufacturer)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("model", utils.StringValue(node.Model)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hostname", utils.StringValue(node.Hostname)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_type", common.FlattenPtrEnum(node.HostType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_version", utils.StringValue(node.HostVersion)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("aos_version", utils.StringValue(node.AosVersion)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("block_serial_number", utils.StringValue(node.BlockSerialNumber)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_gb", utils.IntValue(node.MemoryGB)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_attributes", node.CustomAttributes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", utils.StringValue(node.OwnerExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("provider_ext_id", utils.StringValue(node.ProviderExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("provider_connection_ext_id", utils.StringValue(node.ProviderConnectionExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_info", flattenCPUInfo(node.CpuInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_details", flattenNodeNetworkDetails(node.NetworkDetails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cvm_connectivity_status", common.FlattenPtrEnum(node.CvmConnectivityStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_connectivity_status", common.FlattenPtrEnum(node.HostConnectivityStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", nodeStateName(node.State)); err != nil {
		return diag.FromErr(err)
	}
	if node.CreatedTime != nil {
		if err := d.Set("created_time", node.CreatedTime.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("provider_data", flattenExtendedNodeProviderData(node.ProviderData)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenNodeLinks(node.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(node.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenNodeIdentifiers(pr []import3.NodeIdentifier) []map[string]interface{} {
	if len(pr) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(pr))
	for i, id := range pr {
		out[i] = map[string]interface{}{
			"type":  common.FlattenPtrEnum(id.Type),
			"value": utils.StringValue(id.Value),
		}
	}
	return out
}

func flattenCPUInfo(cpu *import3.CpuInfo) []map[string]interface{} {
	if cpu == nil {
		return nil
	}
	out := map[string]interface{}{
		"logical_core_count": utils.IntValue(cpu.LogicalCoreCount),
		"manufacturer":       utils.StringValue(cpu.Manufacturer),
		"model":              utils.StringValue(cpu.Model),
		"socket_count":       utils.IntValue(cpu.SocketCount),
	}
	if cpu.CapacityGHz != nil {
		out["capacity_g_hz"] = float64(*cpu.CapacityGHz)
	}
	return []map[string]interface{}{out}
}

func flattenNodeNetworkDetails(nd *import3.NodeNetworkDetails) []map[string]interface{} {
	if nd == nil {
		return nil
	}
	out := map[string]interface{}{
		"bmc":  flattenNetworkDetails(nd.Bmc),
		"cvm":  flattenComponentNetworkDetails(nd.Cvm),
		"host": flattenComponentNetworkDetails(nd.Host),
	}
	return []map[string]interface{}{out}
}

func flattenNetworkDetails(nd *import3.NetworkDetails) []map[string]interface{} {
	if nd == nil {
		return nil
	}
	out := map[string]interface{}{
		"gateway": flattenIPAddress(nd.Gateway),
		"ip":      flattenIPAddress(nd.Ip),
		"vlan_id": utils.IntValue(nd.VlanId),
	}
	return []map[string]interface{}{out}
}

func flattenComponentNetworkDetails(cnd *import3.ComponentNetworkDetails) []map[string]interface{} {
	if cnd == nil {
		return nil
	}
	out := map[string]interface{}{
		"management_network": flattenNetworkDetails(cnd.ManagementNetwork),
	}
	return []map[string]interface{}{out}
}

func flattenExtendedNodeProviderData(pd *import3.ExtendedNodeProviderData) []map[string]interface{} {
	if pd == nil {
		return nil
	}
	out := map[string]interface{}{
		"domain":        utils.StringValue(pd.Domain),
		"is_available":  utils.BoolValue(pd.IsAvailable),
		"is_configured": utils.BoolValue(pd.IsConfigured),
		"mode":          utils.StringValue(pd.Mode),
		"name":          utils.StringValue(pd.Name),
		"node_position": utils.StringValue(pd.NodePosition),
		"tags":          flattenKVStringPairs(pd.Tags),
	}
	if pd.DisplayNameMapping != nil {
		out["display_name_mapping"] = []map[string]interface{}{{
			"domain": utils.StringValue(pd.DisplayNameMapping.Domain),
			"groups": utils.StringValue(pd.DisplayNameMapping.Groups),
			"mode":   utils.StringValue(pd.DisplayNameMapping.Mode),
			"tags":   utils.StringValue(pd.DisplayNameMapping.Tags),
		}}
	}
	if len(pd.Groups) > 0 {
		groups := make([]map[string]interface{}, len(pd.Groups))
		for i, g := range pd.Groups {
			groups[i] = map[string]interface{}{
				"id":   utils.StringValue(g.Id),
				"name": utils.StringValue(g.Name),
			}
		}
		out["groups"] = groups
	}
	return []map[string]interface{}{out}
}

// ipAddressSchema returns the shared ipv4/ipv6 value + prefix_length schema block.
// When computed is true the block is fully Computed (no MaxItems, since it is not configurable).
func ipAddressSchema(computed bool) *schema.Schema {
	s := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": ipVersionAddressSchema(computed),
				"ipv6": ipVersionAddressSchema(computed),
			},
		},
	}
	if computed {
		s.Computed = true
	} else {
		s.Optional = true
		s.Computed = true
		s.MaxItems = 1
	}
	return s
}

func ipVersionAddressSchema(computed bool) *schema.Schema {
	s := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Optional: !computed,
					Computed: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Optional: !computed,
					Computed: true,
				},
			},
		},
	}
	if computed {
		s.Computed = true
	} else {
		s.Optional = true
		s.Computed = true
		s.MaxItems = 1
	}
	return s
}

// nodeNetworkDetailsSchema returns the resource (input) schema for a NetworkDetails block.
func nodeNetworkDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"gateway": ipAddressSchema(false),
				"ip":      ipAddressSchema(false),
				"vlan_id": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func nodeComponentNetworkDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"management_network": nodeNetworkDetailsSchema(),
			},
		},
	}
}

func nodeProviderDataSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"domain":        {Type: schema.TypeString, Computed: true},
				"is_available":  {Type: schema.TypeBool, Computed: true},
				"is_configured": {Type: schema.TypeBool, Computed: true},
				"mode":          {Type: schema.TypeString, Computed: true},
				"name":          {Type: schema.TypeString, Computed: true},
				"node_position": {Type: schema.TypeString, Computed: true},
				"display_name_mapping": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"domain": {Type: schema.TypeString, Computed: true},
							"groups": {Type: schema.TypeString, Computed: true},
							"mode":   {Type: schema.TypeString, Computed: true},
							"tags":   {Type: schema.TypeString, Computed: true},
						},
					},
				},
				"groups": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id":   {Type: schema.TypeString, Computed: true},
							"name": {Type: schema.TypeString, Computed: true},
						},
					},
				},
				"tags": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name":  {Type: schema.TypeString, Computed: true},
							"value": {Type: schema.TypeString, Computed: true},
						},
					},
				},
			},
		},
	}
}
