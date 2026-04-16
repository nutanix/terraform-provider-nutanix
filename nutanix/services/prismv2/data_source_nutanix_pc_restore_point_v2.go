package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixFetchRestorePointV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRestorePointV2Read,
		Schema: map[string]*schema.Schema{
			"restore_source_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restorable_domain_manager_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_manager": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForPcConfig(),
						},
						"is_registered_with_hosting_cluster": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForPcNetwork(),
						},
						"hosting_cluster_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"should_enable_high_availability": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"node_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixRestorePointV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	restoreSourceExtID := utils.StringPtr(d.Get("restore_source_ext_id").(string))
	restorableDomainManagerExtID := utils.StringPtr(d.Get("restorable_domain_manager_ext_id").(string))
	extID := utils.StringPtr(d.Get("ext_id").(string))

	resp, err := conn.DomainManagerBackupsAPIInstance.GetRestorePointById(restoreSourceExtID, restorableDomainManagerExtID, extID)

	if err != nil {
		return diag.Errorf("error while fetching Domain Manager Restore Point Detail: %s", err)
	}

	restorePoint := resp.Data.GetValue().(management.RestorePoint)

	if err := d.Set("tenant_id", utils.StringValue(restorePoint.TenantId)); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(restorePoint.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("creation_time", flattenTime(restorePoint.CreationTime)); err != nil {
		return diag.Errorf("error setting creation_time: %s", err)
	}
	if err := d.Set("domain_manager", flattenDomainManager(restorePoint.DomainManager)); err != nil {
		return diag.Errorf("error setting domain_manager: %s", err)
	}

	d.SetId(utils.StringValue(extID))
	return nil
}

func flattenDomainManager(domainManager *config.DomainManager) []map[string]interface{} {
	if domainManager == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"tenant_id":                          utils.StringValue(domainManager.TenantId),
			"ext_id":                             utils.StringValue(domainManager.ExtId),
			"config":                             flattenPCConfig(domainManager.Config),
			"is_registered_with_hosting_cluster": utils.BoolValue(domainManager.IsRegisteredWithHostingCluster),
			"network":                            flattenPCNetwork(domainManager.Network),
			"hosting_cluster_ext_id":             utils.StringValue(domainManager.HostingClusterExtId),
			"should_enable_high_availability":    utils.BoolValue(domainManager.ShouldEnableHighAvailability),
			"node_ext_ids":                       domainManager.NodeExtIds,
		},
	}
}
