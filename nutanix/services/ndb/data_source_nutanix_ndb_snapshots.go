package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
)

func DataSourceNutanixNDBSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_machine_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceNutanixNDBSnapshotsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	tmsID := ""
	if filter, ok := d.GetOk("filters"); ok {
		filterList := filter.([]interface{})

		for _, v := range filterList {
			val := v.(map[string]interface{})

			if tms, ok := val["time_machine_id"]; ok {
				tmsID = tms.(string)
			}
		}
	}

	resp, err := conn.Service.ListSnapshots(ctx, tmsID)
	if err != nil {
		return diag.FromErr(err)
	}

	if e := d.Set("snapshots", flattenSnapshotsList(resp)); e != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()
	if er != nil {
		return diag.Errorf("Error generating UUID for era snapshots: %+v", er)
	}
	d.SetId(uuid)
	return nil
}

func flattenSnapshotsList(sn *era.ListSnapshots) []map[string]interface{} {
	if sn != nil {
		snpList := []map[string]interface{}{}
		for _, val := range *sn {
			snap := map[string]interface{}{}

			snap["id"] = val.ID
			snap["name"] = val.Name
			snap["description"] = val.Description
			snap["date_created"] = val.DateCreated
			snap["date_modified"] = val.DateModified
			snap["properties"] = flattenDBInstanceProperties(val.Properties)
			snap["tags"] = flattenDBTags(val.Tags)
			snap["snapshot_uuid"] = val.SnapshotUUID
			snap["nx_cluster_id"] = val.NxClusterID
			snap["protection_domain_id"] = val.ProtectionDomainID
			snap["parent_snapshot_id"] = val.ParentSnapshotID
			snap["time_machine_id"] = val.TimeMachineID
			snap["database_node_id"] = val.DatabaseNodeID
			snap["app_info_version"] = val.AppInfoVersion
			snap["status"] = val.Status
			snap["type"] = val.Type
			snap["applicable_types"] = val.ApplicableTypes
			snap["snapshot_timestamp"] = val.SnapshotTimeStamp
			snap["software_snapshot_id"] = val.SoftwareSnapshotID
			snap["software_database_snapshot"] = val.SoftwareDatabaseSnapshot
			snap["dbserver_storage_metadata_version"] = val.DBServerStorageMetadataVersion
			snap["santized_from_snapshot_id"] = val.SanitizedFromSnapshotID
			snap["santized"] = val.Sanitized
			snap["timezone"] = val.TimeZone
			snap["processed"] = val.Processed
			snap["database_snapshot"] = val.DatabaseSnapshot
			snap["from_timestamp"] = val.FromTimeStamp
			snap["to_timestamp"] = val.ToTimeStamp
			snap["dbserver_id"] = val.DbserverID
			snap["dbserver_name"] = val.DbserverName
			snap["dbserver_ip"] = val.DbserverIP
			snap["replicated_snapshots"] = val.ReplicatedSnapshots
			snap["software_snapshot"] = val.SoftwareSnapshot
			snap["santized_snapshots"] = val.SanitizedSnapshots
			snap["snapshot_family"] = val.SnapshotFamily
			snap["snapshot_timestamp_date"] = val.SnapshotTimeStampDate
			snap["lcm_config"] = flattenDBLcmConfig(val.LcmConfig)
			snap["parent_snapshot"] = val.ParentSnapshot
			snap["snapshot_size"] = val.SnapshotSize

			snpList = append(snpList, snap)
		}
		return snpList
	}
	return nil
}
