package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkingConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
	networkingMappingReq "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/request/vpcvirtualswitchmappings"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVpcVirtualSwitchMappingsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVpcVirtualSwitchMappingsV2Read,
		Schema: map[string]*schema.Schema{
			"vpc_virtual_switch_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
						},
						"virtual_switch_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the virtual switch.",
						},
						"cluster_uuids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "UUID of the cluster.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_all_traffic_permitted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to permit all traffic through virtual switch or only the ICMP and statistics collection requests.",
						},
						"project_ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the project that owns this entity",
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
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceMetadataSchemaV2(),
							},
						},
					},
				},
			},
		},
	}
}

func datasourceNutanixVpcVirtualSwitchMappingsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	listReq := networkingMappingReq.ListVpcVirtualSwitchMappingsRequest{}

	resp, err := conn.VpcVirtualSwitchMappingsAPI.ListVpcVirtualSwitchMappings(ctx, &listReq)
	if err != nil {
		return diag.Errorf("error while fetching VPC Virtual Switch Mappings: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("vpc_virtual_switch_mappings", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of vpc_virtual_switch_mappings.",
		}}
	}

	getResp := resp.Data.GetValue().([]networkingConfig.VpcVirtualSwitchMapping)

	if err := d.Set("vpc_virtual_switch_mappings", flattenVpcVirtualSwitchMappings(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVpcVirtualSwitchMappings(mappings []networkingConfig.VpcVirtualSwitchMapping) []map[string]interface{} {
	if len(mappings) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(mappings))
	for i, m := range mappings {
		result[i] = map[string]interface{}{
			"ext_id":                   utils.StringValue(m.ExtId),
			"virtual_switch_uuid":      utils.StringValue(m.VirtualSwitchUuid),
			"cluster_uuids":            m.ClusterUuids,
			"is_all_traffic_permitted": utils.BoolValue(m.IsAllTrafficPermitted),
			"project_ext_id":           utils.StringValue(m.ProjectExtId),
			"tenant_id":                utils.StringValue(m.TenantId),
			"links":                    flattenLinks(m.Links),
			"metadata":                 flattenMetadata(m.Metadata),
		}
	}
	return result
}
