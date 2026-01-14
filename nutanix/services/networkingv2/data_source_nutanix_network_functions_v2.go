package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNetworkFunctionsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNetworkFunctionsV2Read,
		Schema: map[string]*schema.Schema{
			"page":     {Type: schema.TypeInt, Optional: true},
			"limit":    {Type: schema.TypeInt, Optional: true},
			"filter":   {Type: schema.TypeString, Optional: true},
			"order_by": {Type: schema.TypeString, Optional: true},
			"network_functions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id":    {Type: schema.TypeString, Computed: true},
						"tenant_id": {Type: schema.TypeString, Computed: true},
						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {Type: schema.TypeString, Computed: true},
									"rel":  {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceMetadataSchemaV2(),
							},
						},
						"name":                    {Type: schema.TypeString, Computed: true},
						"description":             {Type: schema.TypeString, Computed: true},
						"failure_handling":        {Type: schema.TypeString, Computed: true},
						"high_availability_mode":  {Type: schema.TypeString, Computed: true},
						"traffic_forwarding_mode": {Type: schema.TypeString, Computed: true},
						"data_plane_health_check_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"failure_threshold": {Type: schema.TypeInt, Computed: true},
									"interval_secs":     {Type: schema.TypeInt, Computed: true},
									"success_threshold": {Type: schema.TypeInt, Computed: true},
									"timeout_secs":      {Type: schema.TypeInt, Computed: true},
								},
							},
						},
						"nic_pairs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ingress_nic_reference":    {Type: schema.TypeString, Computed: true},
									"egress_nic_reference":     {Type: schema.TypeString, Computed: true},
									"is_enabled":               {Type: schema.TypeBool, Computed: true},
									"vm_reference":             {Type: schema.TypeString, Computed: true},
									"data_plane_health_status": {Type: schema.TypeString, Computed: true},
									"high_availability_state":  {Type: schema.TypeString, Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceNutanixNetworkFunctionsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	var filter, orderBy *string
	var page, limit *int

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

	resp, err := conn.NetworkFunctionAPI.ListNetworkFunctions(page, limit, filter, orderBy)
	if err != nil {
		return diag.Errorf("error while fetching network functions : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("network_functions", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of network functions.",
		}}
	}

	raw := resp.Data.GetValue()
	var items []import1.NetworkFunction
	switch v := raw.(type) {
	case []import1.NetworkFunction:
		items = v
	case []*import1.NetworkFunction:
		for _, it := range v {
			if it != nil {
				items = append(items, *it)
			}
		}
	default:
		return diag.Errorf("unexpected network functions response type: %T", raw)
	}

	if err := d.Set("network_functions", flattenNetworkFunctionEntities(items)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenNetworkFunctionEntities(items []import1.NetworkFunction) []interface{} {
	if len(items) == 0 {
		return nil
	}
	out := make([]interface{}, len(items))
	for i, nf := range items {
		m := make(map[string]interface{})
		m["ext_id"] = nf.ExtId
		m["tenant_id"] = nf.TenantId
		m["links"] = flattenLinks(nf.Links)
		m["metadata"] = flattenMetadata(nf.Metadata)
		m["name"] = nf.Name
		m["description"] = nf.Description
		m["failure_handling"] = common.FlattenPtrEnum(nf.FailureHandling)
		m["high_availability_mode"] = common.FlattenPtrEnum(nf.HighAvailabilityMode)
		m["traffic_forwarding_mode"] = common.FlattenPtrEnum(nf.TrafficForwardingMode)
		m["data_plane_health_check_config"] = flattenDataPlaneHealthCheckConfig(nf.DataPlaneHealthCheckConfig)
		m["nic_pairs"] = flattenNicPairs(nf.NicPairs)
		out[i] = m
	}
	return out
}
