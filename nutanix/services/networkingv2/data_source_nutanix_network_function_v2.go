package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNetworkFunctionV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixNetworkFunctionV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_handling": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"high_availability_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_forwarding_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
	}
}

func DataSourceNutanixNetworkFunctionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	extID := d.Get("ext_id").(string)
	resp, err := conn.NetworkFunctionAPI.GetNetworkFunctionById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching network function : %v", err)
	}

	raw := resp.Data.GetValue()
	var getResp import1.NetworkFunction
	switch v := raw.(type) {
	case import1.NetworkFunction:
		getResp = v
	case *import1.NetworkFunction:
		if v == nil {
			return diag.Errorf("network function response was nil")
		}
		getResp = *v
	default:
		return diag.Errorf("unexpected network function response type: %T", raw)
	}

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("failure_handling", common.FlattenPtrEnum(getResp.FailureHandling)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("high_availability_mode", common.FlattenPtrEnum(getResp.HighAvailabilityMode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("traffic_forwarding_mode", common.FlattenPtrEnum(getResp.TrafficForwardingMode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_plane_health_check_config", flattenDataPlaneHealthCheckConfig(getResp.DataPlaneHealthCheckConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_pairs", flattenNicPairs(getResp.NicPairs)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
