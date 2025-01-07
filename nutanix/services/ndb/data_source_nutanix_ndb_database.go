package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	Era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixEraDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixEraDatabaseRead,
		Schema: map[string]*schema.Schema{
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
			"clustered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"clone": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"era_created": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbserver_logical_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"info": dataSourceEraDatabaseInfo(),
			"metric": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parent_database_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lcm_config":   dataSourceEraLCMConfig(),
			"time_machine": dataSourceEraTimeMachine(),
			"dbserver_logical_cluster": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"database_nodes":   dataSourceEraDatabaseNodes(),
			"linked_databases": dataSourceEraLinkedDatabases(),
			"databases": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNutanixEraDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	dUUID, ok := d.GetOk("database_id")
	if !ok {
		return diag.Errorf("please provide `database_id`")
	}

	resp, err := conn.Service.GetDatabaseInstance(ctx, dUUID.(string))
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

	if err := d.Set("date_created", resp.Datecreated); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_modified", resp.Datemodified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("properties", flattenDBInstanceProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clone", resp.Clone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clustered", resp.Clustered); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_name", resp.Databasename); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database_cluster_type", resp.Databaseclustertype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dbserver_logical_cluster_id", resp.Dbserverlogicalclusterid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_machine_id", resp.Timemachineid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_zone", resp.Timezone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("info", flattenDBInfo(resp.Info)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("metric", resp.Metric); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("parent_database_id", resp.ParentDatabaseID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lcm_config", flattenDBLcmConfig(resp.Lcmconfig)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_machine", flattenDBTimeMachine(resp.TimeMachine)); err != nil {
		return diag.FromErr(err)
	}

	if resp.Dbserverlogicalcluster != nil {
		if err := d.Set("dbserver_logical_cluster", resp.Dbserverlogicalcluster); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("database_nodes", flattenDBNodes(resp.Databasenodes)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("linked_databases", flattenDBLinkedDbs(resp.Linkeddatabases)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("databases", resp.Databases); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	return nil
}

func flattenDBInstanceProperties(pr []*Era.DBInstanceProperties) []map[string]interface{} {
	if len(pr) > 0 {
		res := []map[string]interface{}{}
		for _, v := range pr {
			prop := map[string]interface{}{}

			prop["description"] = v.Description
			prop["name"] = v.Name
			prop["ref_id"] = v.RefID
			prop["secure"] = v.Secure
			prop["value"] = v.Value

			res = append(res, prop)
		}
		return res
	}
	return nil
}

func flattenDBInstanceMetadata(pr *Era.DBInstanceMetadata) []map[string]interface{} {
	if pr != nil {
		pdbmeta := make([]map[string]interface{}, 0)

		pmeta := make(map[string]interface{})
		pmeta["secure_info"] = pr.Secureinfo
		pmeta["info"] = pr.Info
		pmeta["deregister_info"] = flattenDeRegiserInfo(pr.Deregisterinfo)
		pmeta["tm_activate_operation_id"] = pr.Tmactivateoperationid
		pmeta["created_dbservers"] = pr.Createddbservers
		pmeta["registered_dbservers"] = pr.Registereddbservers
		pmeta["last_refresh_timestamp"] = pr.Lastrefreshtimestamp
		pmeta["last_requested_refresh_timestamp"] = pr.Lastrequestedrefreshtimestamp
		pmeta["capability_reset_time"] = pr.CapabilityResetTime
		pmeta["state_before_refresh"] = pr.Statebeforerefresh
		pmeta["state_before_restore"] = pr.Statebeforerestore
		pmeta["state_before_scaling"] = pr.Statebeforescaling
		pmeta["log_catchup_for_restore_dispatched"] = pr.Logcatchupforrestoredispatched
		pmeta["last_log_catchup_for_restore_operation_id"] = pr.Lastlogcatchupforrestoreoperationid
		pmeta["base_size_computed"] = pr.BaseSizeComputed
		pmeta["original_database_name"] = pr.Originaldatabasename
		pmeta["provision_operation_id"] = pr.ProvisionOperationID
		pmeta["source_snapshot_id"] = pr.SourceSnapshotID
		pmeta["pitr_based"] = pr.PitrBased
		pmeta["refresh_blocker_info"] = pr.RefreshBlockerInfo
		pmeta["deregistered_with_delete_time_machine"] = pr.DeregisteredWithDeleteTimeMachine

		pdbmeta = append(pdbmeta, pmeta)
		return pdbmeta
	}
	return nil
}

func flattenDBNodes(pr []Era.Databasenodes) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			db := map[string]interface{}{}

			db["access_level"] = v.AccessLevel
			db["database_id"] = v.Databaseid
			db["database_status"] = v.Databasestatus
			db["date_created"] = v.Datecreated
			db["date_modified"] = v.Datemodified
			db["dbserver_id"] = v.Dbserverid
			db["description"] = v.Description
			db["id"] = v.ID
			db["name"] = v.Name
			db["primary"] = v.Primary
			db["properties"] = flattenDBInstanceProperties(v.Properties)
			db["protection_domain"] = flattenDBProtectionDomain(v.Protectiondomain)
			db["protection_domain_id"] = v.Protectiondomainid
			db["software_installation_id"] = v.Softwareinstallationid
			db["status"] = v.Status
			db["tags"] = flattenDBTags(v.Tags)

			res[k] = db
		}
		return res
	}
	return nil
}

func flattenDBLinkedDbs(pr []Era.Linkeddatabases) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ld := map[string]interface{}{}

			ld["database_name"] = v.DatabaseName
			ld["database_status"] = v.Databasestatus
			ld["date_created"] = v.Datecreated
			ld["date_modified"] = v.Datemodified
			ld["description"] = v.Description
			ld["id"] = v.ID
			ld["metric"] = v.Metric
			ld["name"] = v.Name
			ld["parent_database_id"] = v.ParentDatabaseID
			ld["parent_linked_database_id"] = v.ParentLinkedDatabaseID
			ld["snapshot_id"] = v.SnapshotID
			ld["status"] = v.Status
			ld["timezone"] = v.TimeZone

			res[k] = ld
		}
		return res
	}
	return nil
}

func flattenDBProtectionDomain(pr *Era.Protectiondomain) []map[string]interface{} {
	pDList := make([]map[string]interface{}, 0)
	if pr != nil {
		pmeta := make(map[string]interface{})

		pmeta["cloud_id"] = pr.Cloudid
		pmeta["date_created"] = pr.Datecreated
		pmeta["date_modified"] = pr.Datemodified
		pmeta["description"] = pr.Description
		pmeta["era_created"] = pr.Eracreated
		pmeta["id"] = pr.ID
		pmeta["name"] = pr.Name
		pmeta["owner_id"] = pr.Ownerid
		pmeta["primary_host"] = pr.PrimaryHost
		pmeta["properties"] = flattenDBInstanceProperties(pr.Properties)
		pmeta["status"] = pr.Status
		if pr.Tags != nil {
			pmeta["tags"] = flattenDBTags(pr.Tags)
		}
		pmeta["type"] = pr.Type

		pDList = append(pDList, pmeta)
		return pDList
	}
	return nil
}

func flattenDBTags(pr []*Era.Tags) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			tag := map[string]interface{}{}

			tag["entity_id"] = v.EntityID
			tag["entity_type"] = v.EntityType
			tag["tag_id"] = v.TagID
			tag["tag_name"] = v.TagName
			tag["value"] = v.Value

			res[k] = tag
		}
		return res
	}
	return nil
}

func flattenDBInfo(pr *Era.Info) []map[string]interface{} {
	infoList := make([]map[string]interface{}, 0)
	if pr != nil {
		info := make(map[string]interface{})

		if pr.Secureinfo != nil {
			info["secure_info"] = pr.Secureinfo
		}
		if pr.Info != nil {
			info["bpg_configs"] = flattenBpgConfig(pr.Info.BpgConfigs)
		}
		infoList = append(infoList, info)
		return infoList
	}
	return nil
}

func flattenBpgConfig(pr *Era.BpgConfigs) []map[string]interface{} {
	bpgList := make([]map[string]interface{}, 0)
	if pr != nil {
		bpg := make(map[string]interface{})

		var bgdbParams []map[string]interface{}
		if pr.BpgDBParam != nil {
			bg := make(map[string]interface{})
			bg["maintenance_work_mem"] = utils.StringValue(&pr.BpgDBParam.MaintenanceWorkMem)
			bg["effective_cache_size"] = utils.StringValue(&pr.BpgDBParam.EffectiveCacheSize)
			bg["max_parallel_workers_per_gather"] = utils.StringValue(&pr.BpgDBParam.MaxParallelWorkersPerGather)
			bg["max_worker_processes"] = utils.StringValue(&pr.BpgDBParam.MaxWorkerProcesses)
			bg["shared_buffers"] = utils.StringValue(&pr.BpgDBParam.SharedBuffers)
			bg["work_mem"] = utils.StringValue(&pr.BpgDBParam.WorkMem)
			bgdbParams = append(bgdbParams, bg)
		}
		bpg["bpg_db_param"] = bgdbParams

		var storg []map[string]interface{}
		if pr.Storage != nil {
			str := make(map[string]interface{})

			var storgArch []map[string]interface{}
			if pr.Storage.ArchiveStorage != nil {
				arc := make(map[string]interface{})

				arc["size"] = pr.Storage.ArchiveStorage.Size
				storgArch = append(storgArch, arc)
			}
			str["archive_storage"] = storgArch

			var stdisk []map[string]interface{}
			if pr.Storage.DataDisks != nil {
				arc := make(map[string]interface{})

				arc["count"] = pr.Storage.DataDisks.Count
				stdisk = append(stdisk, arc)
			}
			str["data_disks"] = stdisk

			var stgLog []map[string]interface{}
			if pr.Storage.LogDisks != nil {
				arc := make(map[string]interface{})

				arc["size"] = pr.Storage.LogDisks.Size
				arc["count"] = pr.Storage.LogDisks.Count
				stgLog = append(stgLog, arc)
			}
			str["log_disks"] = stgLog

			storg = append(storg, str)
		}
		bpg["storage"] = storg

		var vmProp []map[string]interface{}
		if pr.VMProperties != nil {
			vmp := make(map[string]interface{})
			vmp["dirty_background_ratio"] = pr.VMProperties.DirtyBackgroundRatio
			vmp["dirty_expire_centisecs"] = pr.VMProperties.DirtyExpireCentisecs
			vmp["dirty_ratio"] = pr.VMProperties.DirtyRatio
			vmp["dirty_writeback_centisecs"] = pr.VMProperties.DirtyWritebackCentisecs
			vmp["nr_hugepages"] = pr.VMProperties.NrHugepages
			vmp["overcommit_memory"] = pr.VMProperties.OvercommitMemory
			vmp["swappiness"] = pr.VMProperties.Swappiness

			vmProp = append(vmProp, vmp)
		}

		bpg["vm_properties"] = vmProp

		bpgList = append(bpgList, bpg)
		return bpgList
	}
	return nil
}

func flattenDBLcmConfig(pr *Era.LcmConfig) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		lcm := map[string]interface{}{}

		lcm["expiry_details"] = flattenEraExpiryDetails(pr.ExpiryDetails)
		lcm["refresh_details"] = flattenEraRefreshDetails(pr.RefreshDetails)

		var preLcmComm []map[string]interface{}
		if pr.PreDeleteCommand != nil {
			pre := map[string]interface{}{}

			pre["command"] = pr.PreDeleteCommand.Command

			preLcmComm = append(preLcmComm, pre)
		}
		lcm["pre_delete_command"] = preLcmComm

		var postLcmComm []map[string]interface{}
		if pr.PreDeleteCommand != nil {
			pre := map[string]interface{}{}

			pre["command"] = pr.PostDeleteCommand.Command

			postLcmComm = append(postLcmComm, pre)
		}
		lcm["post_delete_command"] = postLcmComm

		res = append(res, lcm)
		return res
	}
	return nil
}

func flattenEraExpiryDetails(pr *Era.DBExpiryDetails) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		expiry := map[string]interface{}{}

		expiry["delete_database"] = pr.DeleteDatabase
		expiry["delete_time_machine"] = pr.DeleteTimeMachine
		expiry["delete_vm"] = pr.DeleteVM
		expiry["effective_timestamp"] = pr.EffectiveTimestamp
		expiry["expire_in_days"] = pr.ExpireInDays
		expiry["expiry_date_timezone"] = pr.ExpiryDateTimezone
		expiry["expiry_timestamp"] = pr.ExpiryTimestamp
		expiry["remind_before_in_days"] = pr.RemindBeforeInDays
		expiry["user_created"] = pr.UserCreated

		res = append(res, expiry)
		return res
	}
	return nil
}

func flattenEraRefreshDetails(pr *Era.DBRefreshDetails) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		refresh := map[string]interface{}{}

		refresh["last_refresh_date"] = pr.LastRefreshDate
		refresh["next_refresh_date"] = pr.NextRefreshDate
		refresh["refresh_date_timezone"] = pr.RefreshDateTimezone
		refresh["refresh_in_days"] = pr.RefreshInDays
		refresh["refresh_in_hours"] = pr.RefreshInHours
		refresh["refresh_in_months"] = pr.RefreshInMonths
		refresh["refresh_time"] = pr.RefreshTime

		res = append(res, refresh)
		return res
	}
	return nil
}

func flattenDBTimeMachine(pr *Era.TimeMachine) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		tmac := map[string]interface{}{}

		tmac["id"] = pr.ID
		tmac["name"] = pr.Name
		tmac["description"] = pr.Description
		tmac["date_created"] = pr.DateCreated
		tmac["date_modified"] = pr.DateModified
		tmac["access_level"] = pr.AccessLevel
		tmac["properties"] = flattenDBInstanceProperties(pr.Properties)
		tmac["tags"] = flattenDBTags(pr.Tags)
		tmac["clustered"] = pr.Clustered
		tmac["clone"] = pr.Clone
		tmac["database_id"] = pr.DatabaseID
		tmac["type"] = pr.Type
		tmac["status"] = pr.Status
		tmac["ea_status"] = pr.EaStatus
		tmac["scope"] = pr.Scope
		tmac["sla_id"] = pr.SLAID
		tmac["schedule_id"] = pr.ScheduleID
		tmac["metric"] = pr.Metric
		// tmac["sla_update_metadata"] = pr.SLAUpdateMetadata
		tmac["database"] = pr.Database
		tmac["clones"] = pr.Clones
		tmac["source_nx_clusters"] = pr.SourceNxClusters
		tmac["sla_update_in_progress"] = pr.SLAUpdateInProgress
		tmac["sla"] = flattenDBSLA(pr.SLA)
		tmac["schedule"] = flattenSchedule(pr.Schedule)

		res = append(res, tmac)
		return res
	}
	return nil
}

func flattenDBSLA(pr *Era.ListSLAResponse) []map[string]interface{} {
	res := []map[string]interface{}{}
	if pr != nil {
		sla := map[string]interface{}{}

		sla["id"] = pr.ID
		sla["name"] = pr.Name
		sla["continuous_retention"] = pr.Continuousretention
		sla["daily_retention"] = pr.Dailyretention
		sla["date_modified"] = pr.Datemodified
		sla["date_created"] = pr.Datecreated
		sla["description"] = pr.Description
		sla["monthly_retention"] = pr.Monthlyretention
		sla["owner_id"] = pr.Ownerid
		sla["quarterly_retention"] = pr.Quarterlyretention
		sla["reference_count"] = pr.Referencecount
		sla["system_sla"] = pr.Systemsla
		sla["unique_name"] = pr.Uniquename
		sla["weekly_retention"] = pr.Weeklyretention
		sla["yearly_retention"] = pr.Yearlyretention

		res = append(res, sla)
		return res
	}
	return nil
}

func flattenSchedule(pr *Era.Schedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		sch := map[string]interface{}{}

		sch["continuous_schedule"] = flattenContinousSch(pr.Continuousschedule)
		sch["date_created"] = pr.Datecreated
		sch["date_modified"] = pr.Datemodified
		sch["description"] = pr.Description
		sch["global_policy"] = pr.GlobalPolicy
		sch["id"] = pr.ID
		sch["monthly_schedule"] = flattenMonthlySchedule(pr.Monthlyschedule)
		sch["name"] = pr.Name
		sch["owner_id"] = pr.OwnerID
		sch["quartely_schedule"] = flattenQuartelySchedule(pr.Quartelyschedule)
		sch["reference_count"] = pr.ReferenceCount
		sch["snapshot_time_of_day"] = flattenSnapshotTimeOfDay(pr.Snapshottimeofday)
		sch["start_time"] = pr.StartTime
		sch["system_policy"] = pr.SystemPolicy
		sch["time_zone"] = pr.TimeZone
		sch["unique_name"] = pr.UniqueName
		sch["weekly_schedule"] = flattenWeeklySchedule(pr.Weeklyschedule)
		sch["yearly_schedule"] = flattenYearlylySchedule(pr.Yearlyschedule)
		sch["daily_schedule"] = flattenDailySchedule(pr.Dailyschedule)

		res = append(res, sch)
		return res
	}
	return nil
}

func flattenContinousSch(pr *Era.Continuousschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["log_backup_interval"] = pr.Logbackupinterval
		cr["snapshots_per_day"] = pr.Snapshotsperday

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenMonthlySchedule(pr *Era.Monthlyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_month"] = pr.Dayofmonth

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenQuartelySchedule(pr *Era.Quartelyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_month"] = pr.Dayofmonth
		cr["start_month"] = pr.Startmonth

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenSnapshotTimeOfDay(pr *Era.Snapshottimeofday) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["hours"] = pr.Hours
		cr["minutes"] = pr.Minutes
		cr["seconds"] = pr.Seconds

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenWeeklySchedule(pr *Era.Weeklyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_week"] = pr.Dayofweek

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenYearlylySchedule(pr *Era.Yearlyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		cr["day_of_month"] = pr.Dayofmonth
		cr["month"] = pr.Month

		res = append(res, cr)
		return res
	}
	return nil
}

func flattenDailySchedule(pr *Era.Dailyschedule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	if pr != nil {
		cr := map[string]interface{}{}

		cr["enabled"] = pr.Enabled
		res = append(res, cr)
		return res
	}
	return nil
}

func flattenTimeMachineMetadata(pr *Era.TimeMachineMetadata) []map[string]interface{} {
	if pr != nil {
		tmMeta := make([]map[string]interface{}, 0)
		tm := make(map[string]interface{})

		tm["secure_info"] = pr.SecureInfo
		tm["info"] = pr.Info
		if pr.DeregisterInfo != nil {
			tm["deregister_info"] = flattenDeRegiserInfo(pr.DeregisterInfo)
		}
		tm["capability_reset_time"] = pr.CapabilityResetTime
		tm["auto_heal"] = pr.AutoHeal
		tm["auto_heal_snapshot_count"] = pr.AutoHealSnapshotCount
		tm["auto_heal_log_catchup_count"] = pr.AutoHealLogCatchupCount
		tm["first_snapshot_captured"] = pr.FirstSnapshotCaptured
		tm["first_snapshot_dispatched"] = pr.FirstSnapshotDispatched
		tm["last_snapshot_time"] = pr.LastSnapshotTime
		tm["last_auto_snapshot_time"] = pr.LastAutoSnapshotTime
		tm["last_snapshot_operation_id"] = pr.LastSnapshotOperationID
		tm["last_auto_snapshot_operation_id"] = pr.LastAutoSnapshotOperationID
		tm["last_successful_snapshot_operation_id"] = pr.LastSuccessfulSnapshotOperationID
		tm["snapshot_successive_failure_count"] = pr.SnapshotSuccessiveFailureCount
		tm["last_heal_snapshot_operation"] = pr.LastHealSnapshotOperation
		tm["last_log_catchup_time"] = pr.LastLogCatchupTime
		tm["last_successful_log_catchup_operation_id"] = pr.LastSuccessfulLogCatchupOperationID
		tm["last_log_catchup_operation_id"] = pr.LastLogCatchupOperationID
		tm["log_catchup_successive_failure_count"] = pr.LogCatchupSuccessiveFailureCount
		tm["last_pause_time"] = pr.LastPauseTime
		tm["last_pause_by_force"] = pr.LastPauseByForce
		tm["last_resume_time"] = pr.LastResumeTime
		tm["last_pause_reason"] = pr.LastPauseReason
		tm["state_before_restore"] = pr.StateBeforeRestore
		tm["last_health_alerted_time"] = pr.LastHealthAlertedTime
		tm["last_ea_breakdown_time"] = pr.LastEaBreakdownTime
		tm["authorized_dbservers"] = pr.AuthorizedDbservers
		tm["last_heal_time"] = pr.LastHealTime
		tm["last_heal_system_triggered"] = pr.LastHealSystemTriggered

		tmMeta = append(tmMeta, tm)
		return tmMeta
	}
	return nil
}

func flattenDeRegiserInfo(pr *Era.DeregisterInfo) []map[string]interface{} {
	if pr != nil {
		Deregis := make([]map[string]interface{}, 0)
		regis := map[string]interface{}{}

		regis["message"] = pr.Message
		regis["operations"] = utils.StringValueSlice(pr.Operations)

		Deregis = append(Deregis, regis)
		return Deregis
	}
	return nil
}

func dataSourceEraDatabaseProperties() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of all the properties",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"ref_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"secure": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"description": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceEraDatabaseInfo() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"secure_info": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"bpg_configs": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"storage": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"data_disks": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"count": {
														Type:     schema.TypeFloat,
														Computed: true,
													},
												},
											},
										},
										"log_disks": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"count": {
														Type:     schema.TypeFloat,
														Computed: true,
													},
													"size": {
														Type:     schema.TypeFloat,
														Computed: true,
													},
												},
											},
										},
										"archive_storage": {
											Type:     schema.TypeList,
											Computed: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"size": {
														Type:     schema.TypeFloat,
														Computed: true,
													},
												},
											},
										},
									},
								},
							},
							"vm_properties": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"nr_hugepages": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
										"overcommit_memory": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
										"dirty_ratio": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
										"dirty_background_ratio": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
										"dirty_expire_centisecs": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
										"dirty_writeback_centisecs": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
										"swappiness": {
											Type:     schema.TypeFloat,
											Computed: true,
										},
									},
								},
							},
							"bpg_db_param": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"shared_buffers": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"maintenance_work_mem": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"work_mem": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"effective_cache_size": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"max_worker_processes": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"max_parallel_workers_per_gather": {
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
		},
	}
}

func dataSourceEraLCMConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expiry_details": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"remind_before_in_days": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"effective_timestamp": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"expiry_timestamp": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"expiry_date_timezone": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"user_created": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"expire_in_days": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"delete_database": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"delete_time_machine": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"delete_vm": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
				"refresh_details": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"refresh_in_days": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"refresh_in_hours": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"refresh_in_months": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"last_refresh_date": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"next_refresh_date": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"refresh_time": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"refresh_date_timezone": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"pre_delete_command": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"command": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"post_delete_command": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"command": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEraTimeMachine() *schema.Schema {
	return &schema.Schema{
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
				"access_level": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"properties": dataSourceEraDatabaseProperties(),
				"tags":       dataSourceEraDBInstanceTags(),
				"clustered": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"clone": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"database_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"ea_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"scope": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sla_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"schedule_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"database": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"clones": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"source_nx_clusters": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"sla_update_in_progress": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"metric": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sla_update_metadata": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sla": {
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
							"unique_name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"description": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"owner_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"system_sla": {
								Type:     schema.TypeBool,
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

							"continuous_retention": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"daily_retention": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"weekly_retention": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"monthly_retention": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"quarterly_retention": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"yearly_retention": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"reference_count": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"pitr_enabled": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"current_active_frequency": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"schedule": {
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
							"unique_name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"description": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"owner_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"system_policy": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"global_policy": {
								Type:     schema.TypeBool,
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
							"snapshot_time_of_day": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"hours": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"minutes": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"seconds": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"extra": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"continuous_schedule": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"log_backup_interval": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"snapshots_per_day": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"enabled": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"weekly_schedule": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"day_of_week": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"day_of_week_value": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"enabled": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"monthly_schedule": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"day_of_month": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"enabled": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"yearly_schedule": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"month": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"month_value": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"day_of_month": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"enabled": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"quartely_schedule": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"start_month": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"start_month_value": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"day_of_month": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"enabled": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"daily_schedule": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"enabled": {
											Type:     schema.TypeBool,
											Computed: true,
										},
									},
								},
							},
							"reference_count": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"start_time": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"time_zone": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEraDatabaseNodes() *schema.Schema {
	return &schema.Schema{
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
				"access_level": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"properties": dataSourceEraDatabaseProperties(),
				"tags":       dataSourceEraDBInstanceTags(),
				"database_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"database_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"primary": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"dbserver_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"software_installation_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"protection_domain_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"info": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"secure_info": {
								Type:     schema.TypeMap,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"info": {
								Type:     schema.TypeMap,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"dbserver": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"protection_domain": {
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
							"type": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"cloud_id": {
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
							"owner_id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"status": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"primary_host": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"properties": {
								Type:        schema.TypeList,
								Description: "List of all the properties",
								Computed:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"value": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"ref_id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"secure": {
											Type:     schema.TypeBool,
											Computed: true,
										},
										"description": {
											Type:     schema.TypeString,
											Computed: true,
										},
									},
								},
							},
							"era_created": {
								Type:     schema.TypeBool,
								Computed: true,
							},
							"assoc_entities": {
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
		},
	}
}

func dataSourceEraLinkedDatabases() *schema.Schema {
	return &schema.Schema{
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
				"database_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"database_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"parent_database_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"parent_linked_database_id": {
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
				"timezone": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"info": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"secure_info": {
								Type:     schema.TypeMap,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"info": {
								Type:     schema.TypeMap,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"metric": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"snapshot_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceEraDBInstanceMetadata() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"secure_info": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"info": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"deregister_info": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"message": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"operations": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"tm_activate_operation_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"created_dbservers": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"registered_dbservers": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"last_refresh_timestamp": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"last_requested_refresh_timestamp": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"capability_reset_time": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"state_before_refresh": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"state_before_restore": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"state_before_scaling": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"log_catchup_for_restore_dispatched": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"last_log_catchup_for_restore_operation_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"base_size_computed": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"original_database_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"provision_operation_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"source_snapshot_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"pitr_based": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"refresh_blocker_info": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"deregistered_with_delete_time_machine": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceEraDBInstanceTags() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tag_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"entity_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"entity_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"tag_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}
