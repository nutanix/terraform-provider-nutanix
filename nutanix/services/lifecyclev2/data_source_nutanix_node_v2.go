package lifecyclev2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixNodeV2 returns the hardware, network, and state details of a node by its external ID.
func DatasourceNutanixNodeV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixNodeV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identifiers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":  {Type: schema.TypeString, Computed: true},
						"value": {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"manufacturer":               {Type: schema.TypeString, Computed: true},
			"model":                      {Type: schema.TypeString, Computed: true},
			"hostname":                   {Type: schema.TypeString, Computed: true},
			"host_type":                  {Type: schema.TypeString, Computed: true},
			"host_version":               {Type: schema.TypeString, Computed: true},
			"aos_version":                {Type: schema.TypeString, Computed: true},
			"block_serial_number":        {Type: schema.TypeString, Computed: true},
			"memory_gb":                  {Type: schema.TypeInt, Computed: true},
			"owner_ext_id":               {Type: schema.TypeString, Computed: true},
			"provider_ext_id":            {Type: schema.TypeString, Computed: true},
			"provider_connection_ext_id": {Type: schema.TypeString, Computed: true},
			"cvm_connectivity_status":    {Type: schema.TypeString, Computed: true},
			"host_connectivity_status":   {Type: schema.TypeString, Computed: true},
			"state":                      {Type: schema.TypeString, Computed: true},
			"created_time":               {Type: schema.TypeString, Computed: true},
			"tenant_id":                  {Type: schema.TypeString, Computed: true},
			"custom_attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cpu_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity_g_hz":      {Type: schema.TypeFloat, Computed: true},
						"logical_core_count": {Type: schema.TypeInt, Computed: true},
						"manufacturer":       {Type: schema.TypeString, Computed: true},
						"model":              {Type: schema.TypeString, Computed: true},
						"socket_count":       {Type: schema.TypeInt, Computed: true},
					},
				},
			},
			"network_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bmc":  datasourceNetworkDetailsSchema(),
						"cvm":  datasourceComponentNetworkDetailsSchema(),
						"host": datasourceComponentNetworkDetailsSchema(),
					},
				},
			},
			"provider_data": nodeProviderDataSchema(),
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
		},
	}
}

func datasourceNutanixNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id").(string)
	resp, err := conn.NodesAPIInstance.GetNodeById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching node : %v", err)
	}

	if diags := flattenNodeToState(d, resp.Data.GetValue()); diags != nil {
		return diags
	}
	d.SetId(extID)
	return nil
}

// datasourceNetworkDetailsSchema returns the fully-Computed NetworkDetails datasource block.
func datasourceNetworkDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"gateway": ipAddressSchema(true),
				"ip":      ipAddressSchema(true),
				"vlan_id": {Type: schema.TypeInt, Computed: true},
			},
		},
	}
}

func datasourceComponentNetworkDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"management_network": datasourceNetworkDetailsSchema(),
			},
		},
	}
}
