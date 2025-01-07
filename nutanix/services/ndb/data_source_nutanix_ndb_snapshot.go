package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBSnapshotRead,
		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "UTC",
						},
						"load_replicated_child_snapshots": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
					},
				},
			},

			// computed args
			"id": {
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
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"properties": dataSourceEraDatabaseProperties(),
			"tags":       dataSourceEraDBInstanceTags(),
			"snapshot_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_info_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"applicable_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"snapshot_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"software_snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"software_database_snapshot": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dbserver_storage_metadata_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"santized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"santized_from_snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"processed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"database_snapshot": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"from_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"to_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbserver_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbserver_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbserver_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replicated_snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"software_snapshot": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"santized_snapshots": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_timestamp_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lcm_config": dataSourceEraLCMConfig(),
			"parent_snapshot": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"snapshot_size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixNDBSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	snapID := ""
	if snapshotID, ok := d.GetOk("snapshot_id"); ok {
		snapID = snapshotID.(string)
	}

	filterParams := &era.FilterParams{}
	if filter, ok := d.GetOk("filters"); ok {
		filterList := filter.([]interface{})

		for _, v := range filterList {
			val := v.(map[string]interface{})

			if timezone, tok := val["timezone"]; tok {
				filterParams.TimeZone = timezone.(string)
			}

			if loadRep, lok := val["load_replicated_child_snapshots"]; lok {
				filterParams.LoadReplicatedChildSnapshots = loadRep.(string)
			}
		}
	} else {
		filterParams.TimeZone = "UTC"
		filterParams.LoadReplicatedChildSnapshots = "false"
	}

	resp, err := conn.Service.GetSnapshot(ctx, snapID, filterParams)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_created", resp.DateCreated); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_modified", resp.DateModified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("properties", flattenDBInstanceProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_uuid", resp.SnapshotUUID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nx_cluster_id", resp.NxClusterID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("protection_domain_id", resp.ProtectionDomainID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_snapshot_id", resp.ParentSnapshotID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_machine_id", resp.TimeMachineID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_node_id", resp.DatabaseNodeID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("app_info_version", resp.AppInfoVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("applicable_types", resp.ApplicableTypes); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_timestamp", resp.SnapshotTimeStamp); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("software_snapshot_id", resp.SoftwareSnapshotID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("software_database_snapshot", resp.SoftwareDatabaseSnapshot); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_storage_metadata_version", resp.DBServerStorageMetadataVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("santized", resp.Sanitized); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("santized_from_snapshot_id", resp.SanitizedFromSnapshotID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("timezone", resp.TimeZone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("processed", resp.Processed); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_snapshot", resp.DatabaseSnapshot); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("from_timestamp", resp.FromTimeStamp); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("to_timestamp", resp.ToTimeStamp); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_id", resp.DbserverID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_name", resp.DbserverName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_ip", resp.DbserverIP); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("replicated_snapshots", resp.ReplicatedSnapshots); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("software_snapshot", resp.SoftwareSnapshot); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("santized_snapshots", resp.SanitizedSnapshots); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_family", resp.SnapshotFamily); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_timestamp_date", resp.SnapshotTimeStampDate); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_snapshot", resp.ParentSnapshot); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_size", resp.SnapshotSize); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lcm_config", flattenDBLcmConfig(resp.LcmConfig)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snapID)
	return nil
}

func flattenClonedMetadata(pr *era.ClonedMetadata) []interface{} {
	if pr != nil {
		cloneMetadata := make([]interface{}, 0)
		meta := make(map[string]interface{})

		meta["secure_info"] = pr.SecureInfo
		meta["info"] = pr.Info
		meta["deregister_info"] = pr.DeregisterInfo
		meta["from_timestamp"] = pr.FromTimeStamp
		meta["to_timestamp"] = pr.ToTimeStamp
		meta["replication_retry_count"] = pr.ReplicationRetryCount
		meta["last_replication_retyr_source_snapshot_id"] = pr.LastReplicationRetrySourceSnapshotID
		meta["async"] = pr.Async
		meta["stand_by"] = pr.Standby
		meta["curation_retry_count"] = pr.CurationRetryCount
		meta["operations_using_snapshot"] = pr.OperationsUsingSnapshot

		cloneMetadata = append(cloneMetadata, meta)

		return cloneMetadata
	}
	return nil
}
