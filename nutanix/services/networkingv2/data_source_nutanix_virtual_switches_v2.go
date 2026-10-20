package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkingConfig "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	networkingVsReq "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/virtualswitches"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVirtualSwitchesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixVirtualSwitchesV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prism Element cluster reference.",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the page number of the result set.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "A URL query parameter that specifies the total number of records returned in the result set.",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to filter a collection of resources.",
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL query parameter that allows clients to specify the sort criteria for the returned list of objects.",
			},
			"virtual_switches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixVirtualSwitchV2(),
			},
		},
	}
}

func datasourceNutanixVirtualSwitchesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	listReq := networkingVsReq.ListVirtualSwitchesRequest{}

	if clusterID, ok := d.GetOk("cluster_id"); ok {
		listReq.XClusterId = utils.StringPtr(clusterID.(string))
	}
	if page, ok := d.GetOk("page"); ok {
		listReq.Page_ = utils.IntPtr(page.(int))
	}
	if limit, ok := d.GetOk("limit"); ok {
		listReq.Limit_ = utils.IntPtr(limit.(int))
	}
	if filter, ok := d.GetOk("filter"); ok {
		listReq.Filter_ = utils.StringPtr(filter.(string))
	}
	if orderBy, ok := d.GetOk("order_by"); ok {
		listReq.Orderby_ = utils.StringPtr(orderBy.(string))
	}

	resp, err := conn.VirtualSwitchAPI.ListVirtualSwitches(ctx, &listReq)
	if err != nil {
		return diag.Errorf("error while fetching Virtual Switches: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("virtual_switches", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of virtual_switches.",
		}}
	}

	getResp := resp.Data.GetValue().([]networkingConfig.VirtualSwitch)

	if err := d.Set("virtual_switches", flattenVirtualSwitchesEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVirtualSwitchesEntities(vsList []networkingConfig.VirtualSwitch) []map[string]interface{} {
	if len(vsList) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(vsList))
	for i, vs := range vsList {
		vsMap := map[string]interface{}{
			"ext_id":                 utils.StringValue(vs.ExtId),
			"name":                   utils.StringValue(vs.Name),
			"description":            utils.StringValue(vs.Description),
			"clusters":               flattenVirtualSwitchClusters(vs.Clusters),
			"mtu":                    vs.Mtu,
			"igmp_spec":              flattenIgmpSpec(vs.IgmpSpec),
			"is_quick_mode":          utils.BoolValue(vs.IsQuickMode),
			"is_default":             utils.BoolValue(vs.IsDefault),
			"has_deployment_error":   utils.BoolValue(vs.HasDeploymentError),
			"has_update_in_progress": utils.BoolValue(vs.HasUpdateInProgress),
			"has_delete_in_progress": utils.BoolValue(vs.HasDeleteInProgress),
			"project_ext_id":         utils.StringValue(vs.ProjectExtId),
			"links":                  flattenLinks(vs.Links),
			"metadata":               flattenMetadata(vs.Metadata),
			"tenant_id":              utils.StringValue(vs.TenantId),
		}
		if vs.BondMode != nil {
			vsMap["bond_mode"] = vs.BondMode.GetName()
		}
		if vs.OwnerType != nil {
			vsMap["owner_type"] = vs.OwnerType.GetName()
		}
		result[i] = vsMap
	}
	return result
}
