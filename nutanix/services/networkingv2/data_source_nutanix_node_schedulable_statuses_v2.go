package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkingConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
	networkingNodesReq "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/request/virtualswitchnodesinfo"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixNodeSchedulableStatusesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixNodeSchedulableStatusesV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Prism Element cluster reference.",
			},
			"node_schedulable_statuses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
						},
						"is_never_schedulable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The boolean value to indicate whether or not node is a storage only node",
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
					},
				},
			},
		},
	}
}

func datasourceNutanixNodeSchedulableStatusesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	listReq := networkingNodesReq.ListNodeSchedulableStatusRequest{}

	if clusterID, ok := d.GetOk("cluster_id"); ok {
		listReq.XClusterId = utils.StringPtr(clusterID.(string))
	}

	resp, err := conn.VirtualSwitchNodesInfoAPI.ListNodeSchedulableStatus(ctx, &listReq)
	if err != nil {
		return diag.Errorf("error while fetching node schedulable statuses: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("node_schedulable_statuses", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of node_schedulable_statuses.",
		}}
	}

	getResp := resp.Data.GetValue().([]networkingConfig.NodeSchedulableStatus)

	if err := d.Set("node_schedulable_statuses", flattenNodeSchedulableStatuses(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenNodeSchedulableStatuses(statuses []networkingConfig.NodeSchedulableStatus) []map[string]interface{} {
	if len(statuses) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(statuses))
	for i, s := range statuses {
		result[i] = map[string]interface{}{
			"ext_id":               utils.StringValue(s.ExtId),
			"is_never_schedulable": utils.BoolValue(s.IsNeverSchedulable),
			"tenant_id":            utils.StringValue(s.TenantId),
			"links":                flattenLinks(s.Links),
		}
	}
	return result
}
