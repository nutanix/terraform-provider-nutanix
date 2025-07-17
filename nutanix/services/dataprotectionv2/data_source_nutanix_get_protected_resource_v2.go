package dataprotectionv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixGetProtectedResourceV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixGetProtectedResourceV2Create,

		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": SchemaForLinks(),
			"entity_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_site_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     SchemaForSourceSiteReference(),
			},
			"site_protection_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     SchemaForSiteProtectionInfo(),
			},
			"replication_states": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     SchemaForReplicationStates(),
			},
			"consistency_group_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_fq_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// DatasourceNutanixGetProtectedResourceV2Create to Get Protected Resource
func DatasourceNutanixGetProtectedResourceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataProtectionAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.ProtectedResource.GetProtectedResourceById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("Error while fetching protected resource: %s", err)
	}

	protectedResource := resp.Data.GetValue().(config.ProtectedResource)

	if err := d.Set("tenant_id", utils.StringValue(protectedResource.TenantId)); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(protectedResource.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("entity_ext_id", utils.StringValue(protectedResource.EntityExtId)); err != nil {
		return diag.Errorf("error setting entity_ext_id: %s", err)
	}
	if err := d.Set("entity_type", flattenEntityType(protectedResource.EntityType)); err != nil {
		return diag.Errorf("error setting entity_type: %s", err)
	}
	if err := d.Set("source_site_reference", flattenSourceSiteReference(protectedResource.SourceSiteReference)); err != nil {
		return diag.Errorf("error setting source_site_reference: %s", err)
	}
	if err := d.Set("site_protection_info", flattenSiteProtectionInfo(protectedResource.SiteProtectionInfo)); err != nil {
		return diag.Errorf("error setting site_protection_info: %s", err)
	}
	if err := d.Set("replication_states", flattenReplicationStates(protectedResource.ReplicationStates)); err != nil {
		return diag.Errorf("error setting replication_states: %s", err)
	}
	if err := d.Set("consistency_group_ext_id", utils.StringValue(protectedResource.ConsistencyGroupExtId)); err != nil {
		return diag.Errorf("error setting consistency_group_ext_id: %s", err)
	}
	if err := d.Set("category_fq_names", protectedResource.CategoryFqNames); err != nil {
		return diag.Errorf("error setting category_fq_names: %s", err)
	}

	d.SetId(utils.StringValue(protectedResource.EntityExtId))

	return nil
}

func flattenReplicationStates(replicationStates []config.ReplicationState) []map[string]interface{} {
	if replicationStates == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0)

	for _, replicationState := range replicationStates {
		result = append(result, map[string]interface{}{
			"protection_policy_ext_id":         utils.StringValue(replicationState.ProtectionPolicyExtId),
			"recovery_point_objective_seconds": utils.Int64Value(replicationState.RecoveryPointObjectiveSeconds),
			"replication_status":               flattenReplicationStatus(replicationState.ReplicationStatus),
			"target_site_reference":            flattenSourceSiteReference(replicationState.TargetSiteReference),
		})
	}

	return result
}

// schema func

func SchemaForSourceSiteReference() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mgmt_cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func SchemaForSiteProtectionInfo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recovery_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restorable_time_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"location_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     SchemaForSourceSiteReference(),
			},
			"synchronous_replication_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func SchemaForReplicationStates() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protection_policy_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recovery_point_objective_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replication_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_site_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     SchemaForSourceSiteReference(),
			},
		},
	}
}

// flatten func

func flattenEntityType(entityType *config.ProtectedEntityType) string {
	if entityType == nil {
		return "UNKNOWN"
	}

	return config.ProtectedEntityType.GetName(*entityType)
}

func flattenSourceSiteReference(sourceSiteReference *config.DataProtectionSiteReference) []map[string]interface{} {
	if sourceSiteReference == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mgmt_cluster_ext_id": utils.StringValue(sourceSiteReference.DomainManagerExtId),
			"cluster_ext_id":      utils.StringValue(sourceSiteReference.ClusterExtId),
		},
	}
}

func flattenSiteProtectionInfo(siteProtectionInfo []config.SiteProtectionInfo) []map[string]interface{} {
	if siteProtectionInfo == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0)

	for _, siteProtectionInfoEntity := range siteProtectionInfo {
		result = append(result, map[string]interface{}{
			"recovery_info":                flattenRecoveryInfo(siteProtectionInfoEntity.RecoveryInfo),
			"location_reference":           flattenSourceSiteReference(siteProtectionInfoEntity.LocationReference),
			"synchronous_replication_role": flattenSynchronousReplicationRole(siteProtectionInfoEntity.SynchronousReplicationRole),
		})
	}

	return result
}

func flattenSynchronousReplicationRole(synchronousReplicationRole *config.SynchronousReplicationRole) string {
	if synchronousReplicationRole == nil {
		return "UNKNOWN"
	}

	return config.SynchronousReplicationRole.GetName(*synchronousReplicationRole)
}

func flattenRecoveryInfo(recoveryInfo *config.RecoveryInfo) []map[string]interface{} {
	if recoveryInfo == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0)

	result = append(result, map[string]interface{}{
		"restorable_time_ranges": flattenRestorableTimeRanges(recoveryInfo.RestorableTimeRanges),
	})

	return result
}

func flattenRestorableTimeRanges(restorableTimeRanges []config.RestorableTimeRange) []map[string]interface{} {
	if restorableTimeRanges == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0)

	for _, restorableTimeRange := range restorableTimeRanges {
		result = append(result, map[string]interface{}{
			"start_time": restorableTimeRange.StartTime.String(),
			"end_time":   restorableTimeRange.EndTime.String(),
		})
	}

	return result
}

func flattenReplicationStatus(protectedResourceReplicationStatus *config.ProtectedResourceReplicationStatus) string {
	if protectedResourceReplicationStatus == nil {
		return "UNKNOWN"
	}

	return config.ProtectedResourceReplicationStatus.GetName(*protectedResourceReplicationStatus)
}
