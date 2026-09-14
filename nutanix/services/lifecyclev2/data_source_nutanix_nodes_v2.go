package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import3 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixNodesV2 returns a paginated list of all nodes managed by Foundation Central.
func DatasourceNutanixNodesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixNodesV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixNodeV2(),
			},
		},
	}
}

func datasourceNutanixNodesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	var page, limit *int
	var filter, orderBy, expand, selects *string

	if v, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(v.(string))
	}

	resp, err := conn.NodesAPIInstance.ListNodes(page, limit, filter, orderBy, expand, selects)
	if err != nil {
		return diag.Errorf("error while fetching nodes : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("nodes", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of nodes.",
		}}
	}

	if err := d.Set("nodes", flattenNodesEntities(resp.Data.GetValue())); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// flattenNodesEntities flattens a ListNodes OneOf ([]Node or []NodeProjection) into state maps.
func flattenNodesEntities(value interface{}) []map[string]interface{} {
	var nodes []*import3.Node
	switch v := value.(type) {
	case []import3.Node:
		for i := range v {
			nodes = append(nodes, &v[i])
		}
	case []import3.NodeProjection:
		for i := range v {
			nodes = append(nodes, nodeFromProjection(&v[i]))
		}
	default:
		return nil
	}

	out := make([]map[string]interface{}, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, flattenNodeEntity(n))
	}
	return out
}

func flattenNodeEntity(node *import3.Node) map[string]interface{} {
	if node == nil {
		return nil
	}
	entity := map[string]interface{}{
		"ext_id":                     utils.StringValue(node.ExtId),
		"identifiers":                flattenNodeIdentifiers(node.Identifiers),
		"manufacturer":               utils.StringValue(node.Manufacturer),
		"model":                      utils.StringValue(node.Model),
		"hostname":                   utils.StringValue(node.Hostname),
		"host_type":                  common.FlattenPtrEnum(node.HostType),
		"host_version":               utils.StringValue(node.HostVersion),
		"aos_version":                utils.StringValue(node.AosVersion),
		"block_serial_number":        utils.StringValue(node.BlockSerialNumber),
		"memory_gb":                  utils.IntValue(node.MemoryGB),
		"custom_attributes":          node.CustomAttributes,
		"owner_ext_id":               utils.StringValue(node.OwnerExtId),
		"provider_ext_id":            utils.StringValue(node.ProviderExtId),
		"provider_connection_ext_id": utils.StringValue(node.ProviderConnectionExtId),
		"cpu_info":                   flattenCPUInfo(node.CpuInfo),
		"network_details":            flattenNodeNetworkDetails(node.NetworkDetails),
		"cvm_connectivity_status":    common.FlattenPtrEnum(node.CvmConnectivityStatus),
		"host_connectivity_status":   common.FlattenPtrEnum(node.HostConnectivityStatus),
		"state":                      nodeStateName(node.State),
		"provider_data":              flattenExtendedNodeProviderData(node.ProviderData),
		"links":                      flattenNodeLinks(node.Links),
		"tenant_id":                  utils.StringValue(node.TenantId),
	}
	if node.CreatedTime != nil {
		entity["created_time"] = node.CreatedTime.String()
	}
	return entity
}
