package monitoringv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitoringModel "github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4/models/monitoring/v4/serviceability"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSdaClusterConfigsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixSdaClusterConfigsV2Read,
		Schema: map[string]*schema.Schema{
			"system_defined_policy_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of the System-Defined Alert Policy.",
			},
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
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixSdaClusterConfigV2(),
			},
		},
	}
}

func datasourceNutanixSdaClusterConfigsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MonitoringAPI

	sdaPolicyExtID := d.Get("system_defined_policy_ext_id").(string)

	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	}

	resp, err := conn.SystemDefinedPoliciesAPI.ListClusterConfigsBySdaId(
		utils.StringPtr(sdaPolicyExtID),
		page, limit, filter, orderBy, selects,
	)
	if err != nil {
		return diag.Errorf("error while fetching SDA cluster configs: %v", err)
	}
	if resp.Data == nil {
		if err := d.Set("cluster_configs", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return nil
	}

	configs := resp.Data.GetValue().([]monitoringModel.ClusterConfig)
	if err := d.Set("cluster_configs", flattenClusterConfigs(configs)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.GenUUID())
	return nil
}
