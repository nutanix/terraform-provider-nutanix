package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixListPcsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixListPcsV2Read,
		Schema: map[string]*schema.Schema{
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixFetchPcV2(),
			},
		},
	}
}

func DatasourceNutanixListPcsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	var selects *string

	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	resp, err := conn.DomainManagerAPIInstance.ListDomainManagers(selects)
	if err != nil {
		return diag.Errorf("Error while Listing Domain Managers configurations Details: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("pcs", []map[string]interface{}{}); err != nil {
			return diag.Errorf("Error setting pcs: %v", err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of PCs.",
		}}
	}

	pcs := resp.Data.GetValue().([]config.DomainManager)

	if err := d.Set("pcs", flattenPcs(pcs)); err != nil {
		return diag.Errorf("Error setting pcs: %v", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenPcs(domainManagers []config.DomainManager) []map[string]interface{} {
	pcs := make([]map[string]interface{}, 0)
	for _, pc := range domainManagers {
		pcMap := map[string]interface{}{
			"ext_id":                             utils.StringValue(pc.ExtId),
			"tenant_id":                          utils.StringValue(pc.TenantId),
			"links":                              flattenLinks(pc.Links),
			"config":                             flattenPCConfig(pc.Config),
			"is_registered_with_hosting_cluster": utils.BoolValue(pc.IsRegisteredWithHostingCluster),
			"network":                            flattenPCNetwork(pc.Network),
			"hosting_cluster_ext_id":             utils.StringValue(pc.HostingClusterExtId),
			"should_enable_high_availability":    utils.BoolValue(pc.ShouldEnableHighAvailability),
			"node_ext_ids":                       pc.NodeExtIds,
		}
		pcs = append(pcs, pcMap)
	}
	return pcs
}
