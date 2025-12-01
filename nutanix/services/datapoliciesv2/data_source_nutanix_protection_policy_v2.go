package datapoliciesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/dataprotection/v4/common"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixProtectionPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixProtectionPolicyV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"links": schemaForLinks(),
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaReplicationLocations(),
			},
			"replication_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaReplicationConfigurations(),
			},
			"category_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_approval_policy_needed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixProtectionPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	extID := d.Get("ext_id")

	resp, err := conn.ProtectionPolicies.GetProtectionPolicyById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching Protection Policy: %s", err)
	}

	getResp := resp.Data.GetValue().(config.ProtectionPolicy)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("replication_locations", flattenReplicationLocations(getResp.ReplicationLocations)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("replication_configurations", flattenReplicationConfigurations(getResp.ReplicationConfigurations)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("category_ids", getResp.CategoryIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_approval_policy_needed", getResp.IsApprovalPolicyNeeded); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}

// schema
func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// flatten funcs
func flattenLinks(links []response.ApiLink) []map[string]interface{} {
	if len(links) > 0 {
		linkList := make([]map[string]interface{}, 0)
		for _, link := range links {
			linkMap := make(map[string]interface{})
			if link.Rel != nil {
				linkMap["rel"] = utils.StringValue(link.Rel)
			}
			if link.Href != nil {
				linkMap["href"] = utils.StringValue(link.Href)
			}

			linkList = append(linkList, linkMap)
		}
		return linkList
	}
	return []map[string]interface{}{}
}

func flattenReplicationLocations(replicationLocations []config.ReplicationLocation) []map[string]interface{} {
	if len(replicationLocations) > 0 {
		replicationLocationList := make([]map[string]interface{}, 0)

		for _, location := range replicationLocations {
			replicationLocation := make(map[string]interface{})
			if location.Label != nil {
				replicationLocation["label"] = location.Label
			}
			if location.DomainManagerExtId != nil {
				replicationLocation["domain_manager_ext_id"] = location.DomainManagerExtId
			}
			if location.ReplicationSubLocation != nil {
				replicationLocation["replication_sub_location"] = flattenReplicationSubLocation(location.ReplicationSubLocation)
			}
			if location.IsPrimary != nil {
				replicationLocation["is_primary"] = location.IsPrimary
			}

			replicationLocationList = append(replicationLocationList, replicationLocation)
		}
		return replicationLocationList
	}
	return nil
}

func flattenReplicationSubLocation(replicationSubLocation *config.OneOfReplicationLocationReplicationSubLocation) []map[string]interface{} {
	if replicationSubLocation != nil {
		replicationSubLocationList := make([]map[string]interface{}, 0)

		objectType := utils.StringValue(replicationSubLocation.ObjectType_)

		if objectType == "datapolicies.v4.config.NutanixCluster" {
			nutanixCluster := make([]map[string]interface{}, 0)

			nutanixClusterObj := replicationSubLocation.GetValue().(config.NutanixCluster)

			clusterExtIds := make([]string, 0)

			clusterExtIds = append(clusterExtIds, nutanixClusterObj.ClusterExtIds...)

			nutanixCluster = append(nutanixCluster, map[string]interface{}{
				"cluster_ext_ids": clusterExtIds,
			})

			replicationSubLocationList = append(replicationSubLocationList, map[string]interface{}{
				"cluster_ext_ids": nutanixCluster,
			})
		}

		return replicationSubLocationList
	}
	return nil
}

func flattenReplicationConfigurations(replicationConfigurations []config.ReplicationConfiguration) []map[string]interface{} {
	if len(replicationConfigurations) > 0 {
		replicationConfigurationsList := make([]map[string]interface{}, 0)

		for _, configuration := range replicationConfigurations {
			replicationConfiguration := make(map[string]interface{})
			if configuration.SourceLocationLabel != nil {
				replicationConfiguration["source_location_label"] = configuration.SourceLocationLabel
			}
			if configuration.RemoteLocationLabel != nil {
				replicationConfiguration["remote_location_label"] = configuration.RemoteLocationLabel
			}
			if configuration.Schedule != nil {
				replicationConfiguration["schedule"] = flattenSchedule(configuration.Schedule)
			}

			replicationConfigurationsList = append(replicationConfigurationsList, replicationConfiguration)
		}
		return replicationConfigurationsList
	}
	return nil
}

func flattenSchedule(schedule *config.Schedule) []map[string]interface{} {
	if schedule != nil {
		scheduleList := make([]map[string]interface{}, 0)

		scheduleMap := make(map[string]interface{})

		if schedule.RecoveryPointType != nil {
			switch *schedule.RecoveryPointType {
			case common.RECOVERYPOINTTYPE_CRASH_CONSISTENT:
				scheduleMap["recovery_point_type"] = "CRASH_CONSISTENT"
			case common.RECOVERYPOINTTYPE_APPLICATION_CONSISTENT:
				scheduleMap["recovery_point_type"] = "APPLICATION_CONSISTENT"
			default:
				scheduleMap["recovery_point_type"] = "UNKNOWN"
			}
		}
		if schedule.RecoveryPointObjectiveTimeSeconds != nil {
			scheduleMap["recovery_point_objective_time_seconds"] = utils.IntValue(schedule.RecoveryPointObjectiveTimeSeconds)
		}
		if schedule.Retention != nil {
			scheduleMap["retention"] = flattenRetention(schedule.Retention)
		}
		if schedule.StartTime != nil {
			scheduleMap["start_time"] = utils.StringValue(schedule.StartTime)
		}
		if schedule.SyncReplicationAutoSuspendTimeoutSeconds != nil {
			scheduleMap["sync_replication_auto_suspend_timeout_seconds"] = utils.IntValue(schedule.SyncReplicationAutoSuspendTimeoutSeconds)
		}

		scheduleList = append(scheduleList, scheduleMap)

		return scheduleList
	}
	return nil
}

func flattenRetention(retention *config.OneOfScheduleRetention) []map[string]interface{} {
	if retention != nil {
		retentionList := make([]map[string]interface{}, 0)

		if *retention.ObjectType_ == "datapolicies.v4.config.LinearRetention" {
			linearRetention := retention.GetValue().(config.LinearRetention)

			linearRetentionList := make([]map[string]interface{}, 0)
			linearRetentionMap := make(map[string]interface{})

			if linearRetention.Local != nil {
				linearRetentionMap["local"] = utils.IntValue(linearRetention.Local)
			}
			if linearRetention.Remote != nil {
				linearRetentionMap["remote"] = utils.IntValue(linearRetention.Remote)
			}

			linearRetentionList = append(linearRetentionList, linearRetentionMap)

			retentionList = append(retentionList, map[string]interface{}{
				"linear_retention": linearRetentionList},
			)
		} else if *retention.ObjectType_ == "datapolicies.v4.config.AutoRollupRetention" {
			autoRollupRetention := retention.GetValue().(config.AutoRollupRetention)

			autoRollupRetentionList := make([]map[string]interface{}, 0)
			autoRollupRetentionMap := make(map[string]interface{})

			if autoRollupRetention.Local != nil {
				autoRollupRetentionMap["local"] = flattenAutoRollupRetentionDetails(autoRollupRetention.Local)
			}
			if autoRollupRetention.Remote != nil {
				autoRollupRetentionMap["remote"] = flattenAutoRollupRetentionDetails(autoRollupRetention.Remote)
			}

			autoRollupRetentionList = append(autoRollupRetentionList, autoRollupRetentionMap)

			retentionList = append(retentionList, map[string]interface{}{
				"auto_rollup_retention": autoRollupRetentionList},
			)
		}

		return retentionList
	}
	return nil
}

func flattenAutoRollupRetentionDetails(retentionDetails *config.AutoRollupRetentionDetails) []map[string]interface{} {
	if retentionDetails != nil {
		retentionDetailsList := make([]map[string]interface{}, 0)

		retentionDetailsMap := make(map[string]interface{})

		if retentionDetails.SnapshotIntervalType != nil {
			switch *retentionDetails.SnapshotIntervalType {
			case config.SNAPSHOTINTERVALTYPE_HOURLY:
				retentionDetailsMap["snapshot_interval_type"] = "HOURLY"
			case config.SNAPSHOTINTERVALTYPE_DAILY:
				retentionDetailsMap["snapshot_interval_type"] = "DAILY"
			case config.SNAPSHOTINTERVALTYPE_WEEKLY:
				retentionDetailsMap["snapshot_interval_type"] = "WEEKLY"
			case config.SNAPSHOTINTERVALTYPE_MONTHLY:
				retentionDetailsMap["snapshot_interval_type"] = "MONTHLY"
			case config.SNAPSHOTINTERVALTYPE_YEARLY:
				retentionDetailsMap["snapshot_interval_type"] = "YEARLY"
			default:
				retentionDetailsMap["retention_type"] = "UNKNOWN"
			}
		}
		if retentionDetails.Frequency != nil {
			retentionDetailsMap["frequency"] = utils.IntValue(retentionDetails.Frequency)
		}

		retentionDetailsList = append(retentionDetailsList, retentionDetailsMap)

		return retentionDetailsList
	}
	return nil
}
