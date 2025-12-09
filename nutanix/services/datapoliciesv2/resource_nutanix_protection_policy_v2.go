package datapoliciesv2

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/dataprotection/v4/common"
	prism "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixProtectionPoliciesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixProtectionPoliciesV2Create,
		ReadContext:   ResourceNutanixProtectionPoliciesV2Read,
		UpdateContext: ResourceNutanixProtectionPoliciesV2Update,
		DeleteContext: ResourceNutanixProtectionPoliciesV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"replication_locations": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 3, //nolint:gomnd
				Elem:     schemaReplicationLocations(),
			},
			"replication_configurations": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 9, //nolint:gomnd
				Elem:     schemaReplicationConfigurations(),
			},
			"category_ids": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 10, //nolint:gomnd
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: categoryIdsDiffSuppressFunc,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func ResourceNutanixProtectionPoliciesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	bodySpec := config.NewProtectionPolicy()

	if name, ok := d.GetOk("name"); ok {
		bodySpec.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		bodySpec.Description = utils.StringPtr(description.(string))
	}
	if replicationLocations, ok := d.GetOk("replication_locations"); ok {
		bodySpec.ReplicationLocations = expandReplicationLocations(replicationLocations.([]interface{}))
	}
	if replicationConfigurations, ok := d.GetOk("replication_configurations"); ok {
		bodySpec.ReplicationConfigurations = expandReplicationConfigurations(replicationConfigurations.([]interface{}))
	}
	if categoryIds, ok := d.GetOk("category_ids"); ok {
		bodySpec.CategoryIds = commonUtils.ExpandListOfString(categoryIds.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(bodySpec, "", "  ")
	log.Printf("[DEBUG] Create Protection Policy Body Spec: %s", string(aJSON))

	resp, err := conn.ProtectionPolicies.CreateProtectionPolicy(bodySpec)
	if err != nil {
		return diag.Errorf("error while creating Protection Policy: %v", err)
	}

	TaskRef := resp.Data.GetValue().(prism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the protection policy to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for protection policy (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching protection policy task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Protection Policy Task Details: %s", string(aJSON))

	// Extract UUID from completion details
	uuid, err := commonUtils.ExtractCompletionDetailFromTask(taskDetails, utils.CompletionDetailsNameProtectionPolicy, "Protection Policy")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid)

	return ResourceNutanixProtectionPoliciesV2Read(ctx, d, meta)
}

func ResourceNutanixProtectionPoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	extID := d.Id()

	resp, err := conn.ProtectionPolicies.GetProtectionPolicyById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching Protection Policy: %s", err)
	}

	getResp := resp.Data.GetValue().(config.ProtectionPolicy)

	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] Read Protection Policy Response Details: %s", string(aJSON))

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
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

	return nil
}

func ResourceNutanixProtectionPoliciesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	readResp, err := conn.ProtectionPolicies.GetProtectionPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Protection Policy: %v", err)
	}
	// extract e-tag
	args := make(map[string]interface{})
	etag := conn.ProtectionPolicies.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etag)

	updateSpec := config.NewProtectionPolicy()

	if name, ok := d.GetOk("name"); ok {
		updateSpec.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		updateSpec.Description = utils.StringPtr(description.(string))
	}
	if replicationLocations, ok := d.GetOk("replication_locations"); ok {
		updateSpec.ReplicationLocations = expandReplicationLocations(replicationLocations.([]interface{}))
	}
	if replicationConfigurations, ok := d.GetOk("replication_configurations"); ok {
		updateSpec.ReplicationConfigurations = expandReplicationConfigurations(replicationConfigurations.([]interface{}))
	}
	if categoryIds, ok := d.GetOk("category_ids"); ok {
		updateSpec.CategoryIds = commonUtils.ExpandListOfString(categoryIds.([]interface{}))
	}

	resp, err := conn.ProtectionPolicies.UpdateProtectionPolicyById(utils.StringPtr(d.Id()), updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating Protection Policy: %v", err)
	}

	TaskRef := resp.Data.GetValue().(prism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the protection policy to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for protection policy (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching protection policy task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update Protection Policy Task Details: %s", string(aJSON))

	return ResourceNutanixProtectionPoliciesV2Read(ctx, d, meta)
}

func ResourceNutanixProtectionPoliciesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

	resp, err := conn.ProtectionPolicies.DeleteProtectionPolicyById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting Protection Policy: %v", err)
	}
	TaskRef := resp.Data.GetValue().(prism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the protection policy to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for protection policy (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching protection policy delete task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Protection Policy Task Details: %s", string(aJSON))

	return nil
}

// schemas funcs
func schemaReplicationLocations() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_manager_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replication_sub_location": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_ext_ids": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1, //nolint:gomnd
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_ext_ids": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,   //nolint:gomnd
										MaxItems: 200, //nolint:gomnd
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"is_primary": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func schemaReplicationConfigurations() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_location_label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_location_label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recovery_point_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"CRASH_CONSISTENT", "APPLICATION_CONSISTENT"}, false),
						},
						"recovery_point_objective_time_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"retention": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1, //nolint:gomnd
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"linear_retention": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1, //nolint:gomnd
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"remote": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"auto_rollup_retention": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1, //nolint:gomnd
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1, //nolint:gomnd
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"snapshot_interval_type": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice([]string{"YEARLY", "WEEKLY", "DAILY", "MONTHLY", "HOURLY"}, false),
															},
															"frequency": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 24),
															},
														},
													},
												},
												"remote": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1, //nolint:gomnd
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"snapshot_interval_type": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice([]string{"YEARLY", "WEEKLY", "DAILY", "MONTHLY", "HOURLY"}, false),
															},
															"frequency": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 24),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sync_replication_auto_suspend_timeout_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtMost(300), //nolint:gomnd
						},
					},
				},
			},
		},
	}
}

// expander funcs
func expandReplicationLocations(replicationLocations []interface{}) []config.ReplicationLocation {
	if len(replicationLocations) == 0 {
		return nil
	}

	replicationLocationsSpec := make([]config.ReplicationLocation, 0)

	for _, replicationLocation := range replicationLocations {
		replicationLocationVal := replicationLocation.(map[string]interface{})

		replicationLocationSpec := config.NewReplicationLocation()
		if sourceLocationLabel, ok := replicationLocationVal["label"]; ok {
			replicationLocationSpec.Label = utils.StringPtr(sourceLocationLabel.(string))
		}
		if domainManagerExtID, ok := replicationLocationVal["domain_manager_ext_id"]; ok {
			replicationLocationSpec.DomainManagerExtId = utils.StringPtr(domainManagerExtID.(string))
		}
		if replicationSubLocation, ok := replicationLocationVal["replication_sub_location"]; ok {
			replicationLocationSpec.ReplicationSubLocation = expandOneOfReplicationLocationReplicationSubLocation(replicationSubLocation.([]interface{}))
		}
		if isPrimary, ok := replicationLocationVal["is_primary"]; ok {
			replicationLocationSpec.IsPrimary = utils.BoolPtr(isPrimary.(bool))
		}
		replicationLocationsSpec = append(replicationLocationsSpec, *replicationLocationSpec)
	}

	return replicationLocationsSpec
}

func expandOneOfReplicationLocationReplicationSubLocation(oneOfReplicationLocationReplicationSubLocations []interface{}) *config.OneOfReplicationLocationReplicationSubLocation {
	if len(oneOfReplicationLocationReplicationSubLocations) == 0 {
		return nil
	}

	oneOfReplicationLocationReplicationSubLocationI := oneOfReplicationLocationReplicationSubLocations[0]
	oneOfReplicationLocationReplicationSubLocationVal := oneOfReplicationLocationReplicationSubLocationI.(map[string]interface{})

	oneOfReplicationLocationReplicationSubLocationSpec := config.NewOneOfReplicationLocationReplicationSubLocation()

	if clusterExtIdsMap, ok := oneOfReplicationLocationReplicationSubLocationVal["cluster_ext_ids"]; ok && len(clusterExtIdsMap.([]interface{})) > 0 {
		nutanixCluster := config.NewNutanixCluster()
		clusterExtIdsI := clusterExtIdsMap.([]interface{})[0]
		clusterExtIdsVal := clusterExtIdsI.(map[string]interface{})

		if clusterExtIds, ok := clusterExtIdsVal["cluster_ext_ids"]; ok {
			nutanixCluster.ClusterExtIds = commonUtils.ExpandListOfString(clusterExtIds.([]interface{}))
		}

		err := oneOfReplicationLocationReplicationSubLocationSpec.SetValue(*nutanixCluster)
		if err != nil {
			log.Printf("[ERROR] Error while setting value for OneOfReplicationLocationReplicationSubLocation: %v", err)
			return nil
		}
	}

	return oneOfReplicationLocationReplicationSubLocationSpec
}

func expandReplicationConfigurations(replicationConfigurationsData []interface{}) []config.ReplicationConfiguration {
	if len(replicationConfigurationsData) == 0 {
		return nil
	}

	replicationConfigurations := make([]config.ReplicationConfiguration, 0)

	for _, replicationConfigurationData := range replicationConfigurationsData {
		replicationConfigurationDataMap := replicationConfigurationData.(map[string]interface{})

		replicationConfiguration := config.ReplicationConfiguration{}
		if sourceLocationLabel, ok := replicationConfigurationDataMap["source_location_label"]; ok {
			replicationConfiguration.SourceLocationLabel = utils.StringPtr(sourceLocationLabel.(string))
		}
		if remoteLocationLabel, ok := replicationConfigurationDataMap["remote_location_label"]; ok && remoteLocationLabel != "" {
			replicationConfiguration.RemoteLocationLabel = utils.StringPtr(remoteLocationLabel.(string))
		}
		if schedule, ok := replicationConfigurationDataMap["schedule"]; ok {
			replicationConfiguration.Schedule = expandSchedule(schedule.([]interface{}))
		}
		replicationConfigurations = append(replicationConfigurations, replicationConfiguration)
	}

	return replicationConfigurations
}

func expandSchedule(scheduleData []interface{}) *config.Schedule {
	if len(scheduleData) == 0 {
		return nil
	}

	scheduleDataMap := scheduleData[0].(map[string]interface{})

	schedule := config.NewSchedule()
	if recoveryPointType, ok := scheduleDataMap["recovery_point_type"]; ok && recoveryPointType != "" {
		schedule.RecoveryPointType = expandRecoveryPointType(recoveryPointType.(string))
	}
	if recoveryPointObjectiveTimeSeconds, ok := scheduleDataMap["recovery_point_objective_time_seconds"]; ok {
		schedule.RecoveryPointObjectiveTimeSeconds = utils.IntPtr(recoveryPointObjectiveTimeSeconds.(int))
	}
	if retention, ok := scheduleDataMap["retention"]; ok {
		schedule.Retention = expandRetention(retention.([]interface{}))
	}
	if startTime, ok := scheduleDataMap["start_time"]; ok && startTime != "" {
		schedule.StartTime = utils.StringPtr(startTime.(string))
	}
	if syncReplicationAutoSuspendTimeoutSeconds, ok := scheduleDataMap["sync_replication_auto_suspend_timeout_seconds"]; ok {
		schedule.SyncReplicationAutoSuspendTimeoutSeconds = utils.IntPtr(syncReplicationAutoSuspendTimeoutSeconds.(int))
	}

	return schedule
}

func expandRecoveryPointType(recoveryPointType string) *common.RecoveryPointType {
	if recoveryPointType == "" {
		return nil
	}

	const CrashConsistent, ApplicationConsistent = 2, 3
	switch recoveryPointType {
	case "CRASH_CONSISTENT":
		p := common.RecoveryPointType(CrashConsistent)
		return &p
	case "APPLICATION_CONSISTENT":
		p := common.RecoveryPointType(ApplicationConsistent)
		return &p
	}
	return nil
}

func expandRetention(retention []interface{}) *config.OneOfScheduleRetention {
	if len(retention) == 0 {
		return nil
	}

	retentionData := retention[0].(map[string]interface{})

	retentionSpec := config.NewOneOfScheduleRetention()

	if linearRetention, ok := retentionData["linear_retention"]; ok && len(linearRetention.([]interface{})) > 0 {
		linearRetentionSpec := expandLinearRetention(linearRetention.([]interface{}))
		err := retentionSpec.SetValue(*linearRetentionSpec)
		if err != nil {
			log.Printf("[ERROR] Error while setting value for LinearRetention: %v", err)
			return nil
		}
	} else if autoRollupRetention, ok := retentionData["auto_rollup_retention"]; ok && len(autoRollupRetention.([]interface{})) > 0 {
		autoRollupRetentionSpec := expandAutoRollupRetention(autoRollupRetention.([]interface{}))
		err := retentionSpec.SetValue(*autoRollupRetentionSpec)
		if err != nil {
			log.Printf("[ERROR] Error while setting value for AutoRollupRetention: %v", err)
			return nil
		}
	}

	return retentionSpec
}

func expandLinearRetention(linearRetentionData []interface{}) *config.LinearRetention {
	if len(linearRetentionData) == 0 {
		return nil
	}

	linearRetentionDataMap := linearRetentionData[0].(map[string]interface{})

	linearRetention := config.NewLinearRetention()
	if local, ok := linearRetentionDataMap["local"]; ok && local != "" {
		linearRetention.Local = utils.IntPtr(local.(int))
	}
	if remote, ok := linearRetentionDataMap["remote"]; ok && remote.(int) > 0 {
		linearRetention.Remote = utils.IntPtr(remote.(int))
	}

	return linearRetention
}

func expandAutoRollupRetention(autoRollupRetentionData []interface{}) *config.AutoRollupRetention {
	if len(autoRollupRetentionData) == 0 {
		return nil
	}

	autoRollupRetentionDataMap := autoRollupRetentionData[0].(map[string]interface{})

	autoRollupRetention := config.NewAutoRollupRetention()
	if local, ok := autoRollupRetentionDataMap["local"]; ok && len(local.([]interface{})) > 0 {
		autoRollupRetention.Local = expandAutoRollupRetentionDetails(local.([]interface{}))
	}
	if remote, ok := autoRollupRetentionDataMap["remote"]; ok && len(remote.([]interface{})) > 0 {
		autoRollupRetention.Remote = expandAutoRollupRetentionDetails(remote.([]interface{}))
	}

	return autoRollupRetention
}

func expandAutoRollupRetentionDetails(autoRollupRetentionLocal []interface{}) *config.AutoRollupRetentionDetails {
	if len(autoRollupRetentionLocal) == 0 {
		return nil
	}

	autoRollupRetentionLocalDataMap := autoRollupRetentionLocal[0].(map[string]interface{})

	autoRollupRetentionLocalSpec := config.AutoRollupRetentionDetails{}
	if snapshotIntervalType, ok := autoRollupRetentionLocalDataMap["snapshot_interval_type"]; ok && snapshotIntervalType != "" {
		autoRollupRetentionLocalSpec.SnapshotIntervalType = expandSnapshotIntervalType(snapshotIntervalType.(string))
	}
	if frequency, ok := autoRollupRetentionLocalDataMap["frequency"]; ok && frequency != "" {
		autoRollupRetentionLocalSpec.Frequency = utils.IntPtr(frequency.(int))
	}

	return &autoRollupRetentionLocalSpec
}

func expandSnapshotIntervalType(snapshotIntervalType string) *config.SnapshotIntervalType {
	if snapshotIntervalType == "" {
		return nil
	}

	const HOURLY, DAILY, WEEKLY, MONTHLY, YEARLY = 2, 3, 4, 5, 6
	switch snapshotIntervalType {
	case "YEARLY":
		p := config.SnapshotIntervalType(YEARLY)
		return &p
	case "WEEKLY":
		p := config.SnapshotIntervalType(WEEKLY)
		return &p
	case "DAILY":
		p := config.SnapshotIntervalType(DAILY)
		return &p
	case "MONTHLY":
		p := config.SnapshotIntervalType(MONTHLY)
		return &p
	case "HOURLY":
		p := config.SnapshotIntervalType(HOURLY)
		return &p
	}
	return nil
}

func categoryIdsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.HasChange("category_ids") {
		oldCap, newCap := d.GetChange("category_ids")
		log.Printf("[DEBUG] oldCap : %v", oldCap)
		log.Printf("[DEBUG] newCap : %v", newCap)

		oldList := oldCap.([]interface{})
		newList := newCap.([]interface{})

		if len(oldList) != len(newList) {
			log.Printf("[DEBUG] category_ids are different")
			return false
		}

		sort.SliceStable(oldList, func(i, j int) bool {
			return oldList[i].(string) < oldList[j].(string)
		})
		sort.SliceStable(newList, func(i, j int) bool {
			return newList[i].(string) < newList[j].(string)
		})

		aJSON, _ := json.Marshal(oldList)
		log.Printf("[DEBUG] oldList : %s", aJSON)
		aJSON, _ = json.Marshal(newList)
		log.Printf("[DEBUG] newList : %s", aJSON)

		if reflect.DeepEqual(oldList, newList) {
			log.Printf("[DEBUG] category_ids are same")
			return true
		}
		log.Printf("[DEBUG] category_ids are different")
		return false
	}
	return false
}
