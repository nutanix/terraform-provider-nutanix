package prismv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixListRestorablePcsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixListRestorablePcsV2Read,
		Schema: map[string]*schema.Schema{
			"restore_source_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restorable_pcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixFetchPcV2(),
			},
		},
	}
}

func DatasourceNutanixListRestorablePcsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	var filter *string
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

	restoreSourceExtID := d.Get("restore_source_ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.ListRestorableDomainManagers(utils.StringPtr(restoreSourceExtID), page, limit, filter)
	if err != nil {
		return diag.Errorf("Error while Listing Restorable Domain Managers configurations Details: %v", err)
	}

	aJSON, _ := json.MarshalIndent(resp, "", "  ")
	log.Printf("[DEBUG] ListRestorableDomainManagers Response: %v", string(aJSON))

	if resp.Data == nil {
		if err := d.Set("restorable_pcs", make([]interface{}, 0)); err != nil {
			return diag.Errorf("Error setting Restorable pcs: %v", err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of restorable PCs.",
		}}
	}

	restorablePcs := resp.Data.GetValue().([]management.RestorableDomainManager)
	if err := d.Set("restorable_pcs", flattenRestorablePcs(restorablePcs)); err != nil {
		return diag.Errorf("Error setting pcs: %v", err)
	}

	d.SetId(utils.GenUUID())

	return nil
}

func flattenRestorablePcs(restorableDomainManagers []management.RestorableDomainManager) []map[string]interface{} {
	restorablePcs := make([]map[string]interface{}, 0)

	for _, pc := range restorableDomainManagers {
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
		restorablePcs = append(restorablePcs, pcMap)
	}
	return restorablePcs
}
