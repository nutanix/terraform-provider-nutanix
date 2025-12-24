package prismv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var exactlyOneOfLocation = []string{
	"location.0.cluster_location",
	"location.0.object_store_location",
}

const (
	clustersLocationObjectType    = "prism.v4.management.ClusterLocation"
	objectStoreLocationObjectType = "prism.v4.management.ObjectStoreLocation"
	awsS3ConfigObjectType         = "prism.v4.management.AWSS3Config"
)

func ResourceNutanixBackupTargetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixBackupTargetV2Create,
		ReadContext:   ResourceNutanixBackupTargetV2Read,
		UpdateContext: ResourceNutanixBackupTargetV2Update,
		DeleteContext: ResourceNutanixBackupTargetV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				const expectedPartsCount = 2
				parts := strings.Split(d.Id(), "/")
				if len(parts) != expectedPartsCount {
					return nil, fmt.Errorf("invalid import uuid (%q), expected domain_manager_ext_id/backup_target_ext_id", d.Id())
				}
				d.Set("domain_manager_ext_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"domain_manager_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_location": {
							Type:         schema.TypeList,
							Optional:     true,
							ExactlyOneOf: exactlyOneOfLocation,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"object_store_location": {
							Type:         schema.TypeList,
							Optional:     true,
							ExactlyOneOf: exactlyOneOfLocation,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider_config": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket_name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"region": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "us-east-1",
												},
												"credentials": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"access_key_id": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
															"secret_access_key": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
											},
										},
									},
									"backup_policy": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rpo_in_minutes": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(60, 1440), //nolint:gomnd
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
			// computed attributes for read
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"last_sync_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_backup_paused": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backup_pause_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixBackupTargetV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	body := management.BackupTarget{}

	OneOfBackupTargetLocation := management.NewOneOfBackupTargetLocation()
	locationI := d.Get("location").([]interface{})
	location := locationI[0].(map[string]interface{})

	clusterExtID := ""
	isClusterLocation := false
	bucketName := ""
	isObjectStoreLocation := false
	if location["cluster_location"] != nil && len(location["cluster_location"].([]interface{})) > 0 {
		clusterLocation := location["cluster_location"].([]interface{})[0].(map[string]interface{})
		clusterConfig := clusterLocation["config"].([]interface{})[0].(map[string]interface{})

		clusterConfigBody := management.NewClusterLocation()
		clusterRef := management.NewClusterReference()
		clusterExtID = clusterConfig["ext_id"].(string)
		clusterRef.ExtId = utils.StringPtr(clusterExtID)
		// From IRIS SDK, the cluster location config is a OneOfClusterLocationConfig
		// so we need to set the value of the OneOfClusterLocationConfig
		oneOfClusterLocationConfig := management.NewOneOfClusterLocationConfig()
		oneOfClusterLocationConfig.SetValue(*clusterRef)
		clusterConfigBody.Config = oneOfClusterLocationConfig

		err := OneOfBackupTargetLocation.SetValue(*clusterConfigBody)
		if err != nil {
			return diag.Errorf("error while setting cluster location : %v", err)
		}
		isClusterLocation = true
	} else if location["object_store_location"] != nil && len(location["object_store_location"].([]interface{})) > 0 {
		objectStoreLocation := location["object_store_location"].([]interface{})[0].(map[string]interface{})
		providerConfig := objectStoreLocation["provider_config"]
		backupPolicy := objectStoreLocation["backup_policy"]

		objectStoreLocationBody := management.NewObjectStoreLocation()

		objectStoreLocationBody.ProviderConfig = expandProviderConfig(providerConfig)
		objectStoreLocationBody.BackupPolicy = expandBackupPolicy(backupPolicy)

		err := OneOfBackupTargetLocation.SetValue(*objectStoreLocationBody)
		if err != nil {
			return diag.Errorf("error while setting object store location : %v", err)
		}
		if utils.StringValue(objectStoreLocationBody.ProviderConfig.ObjectType_) == awsS3ConfigObjectType {
			// Since the backup target ext ID is not returned in the task details response
			// we need to find the backup target by bucket name
			bucketName = utils.StringValue(objectStoreLocationBody.ProviderConfig.GetValue().(management.AWSS3Config).BucketName)
		} else {
			return diag.Errorf("unsupported object store provider config type: %s", utils.StringValue(objectStoreLocationBody.ProviderConfig.ObjectType_))
		}

		isObjectStoreLocation = true
	}

	body.Location = OneOfBackupTargetLocation

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Payload backup target Body: %s", string(aJSON))

	resp, err := conn.DomainManagerBackupsAPIInstance.CreateBackupTarget(utils.StringPtr(domainManagerExtID), &body)

	if err != nil {
		return diag.Errorf("error while creating backup target: %s", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for backup target to be created: %s", err)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching backup target task details: %s", err)
	}

	taskDetails := taskResp.Data.GetValue().(config.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create backup target task details: %s", string(aJSON))

	listBackupTargets, err := conn.DomainManagerBackupsAPIInstance.ListBackupTargets(utils.StringPtr(domainManagerExtID))
	if err != nil {
		return diag.Errorf("error while Listing Backup Targets for : %s err: %s", domainManagerExtID, err)
	}
	backupTargets := listBackupTargets.Data.GetValue().([]management.BackupTarget)

	// Find the new backup target ext id since the response does not contain the ext id
	for _, backupTarget := range backupTargets {
		backupTargetLocation := backupTarget.Location
		if isClusterLocation && utils.StringValue(backupTargetLocation.ObjectType_) == clustersLocationObjectType {
			log.Printf("[DEBUG] Cluster Backup Target with Ext ID: %s", utils.StringValue(backupTarget.ExtId))
			clusterLocation := backupTarget.Location.GetValue().(management.ClusterLocation)
			// From IRIS SDK, the cluster location config is a OneOfClusterLocationConfig
			// so we need to get the value of the OneOfClusterLocationConfig
			clusterConfig := clusterLocation.Config.GetValue().(management.ClusterReference)
			if utils.StringValue(clusterConfig.ExtId) == clusterExtID {
				d.SetId(utils.StringValue(backupTarget.ExtId))
				break
			}
		} else if isObjectStoreLocation && utils.StringValue(backupTargetLocation.ObjectType_) == objectStoreLocationObjectType {
			objectStoreLocation := backupTarget.Location.GetValue().(management.ObjectStoreLocation)
			// Since the backup target ext ID is not returned in the task details response
			// we need to find the backup target by bucket name
			if *objectStoreLocation.ProviderConfig.ObjectType_ == awsS3ConfigObjectType {
				awsS3Config := objectStoreLocation.ProviderConfig.GetValue().(management.AWSS3Config)
				if utils.StringValue(awsS3Config.BucketName) == bucketName {
					d.SetId(utils.StringValue(backupTarget.ExtId))
					break
				}
			}
		}
	}

	if d.Id() == "" {
		return diag.Errorf("error while setting backup target ID")
	}

	return ResourceNutanixBackupTargetV2Read(ctx, d, meta)
}

func ResourceNutanixBackupTargetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), nil)

	if err != nil {
		return diag.Errorf("error while fetching backup target: %s", err)
	}

	backupTarget := resp.Data.GetValue().(management.BackupTarget)

	if err := d.Set("tenant_id", backupTarget.TenantId); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("ext_id", backupTarget.ExtId); err != nil {
		return diag.Errorf("error setting ext_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(backupTarget.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("last_sync_time", flattenTime(backupTarget.LastSyncTime)); err != nil {
		return diag.Errorf("error setting last_sync_time: %s", err)
	}
	if err := d.Set("is_backup_paused", backupTarget.IsBackupPaused); err != nil {
		return diag.Errorf("error setting is_backup_paused: %s", err)
	}
	if err := d.Set("backup_pause_reason", backupTarget.BackupPauseReason); err != nil {
		return diag.Errorf("error setting backup_pause_reason: %s", err)
	}
	if err := d.Set("location", flattenResourceBackupTargetLocation(backupTarget.Location, d.Get("location"))); err != nil {
		return diag.Errorf("error setting location: %s", err)
	}

	return nil
}

func flattenResourceBackupTargetLocation(location *management.OneOfBackupTargetLocation, d interface{}) []map[string]interface{} {
	if location == nil {
		return nil
	}

	backupTargetLocation := make([]map[string]interface{}, 0)

	if utils.StringValue(location.ObjectType_) == clustersLocationObjectType {
		clusterLocation := location.GetValue().(management.ClusterLocation)

		clusterLocationMap := make(map[string]interface{})
		clusterLocationMap["cluster_location"] = flattenClusterLocation(clusterLocation)
		backupTargetLocation = append(backupTargetLocation, clusterLocationMap)
		return backupTargetLocation
	}

	if utils.StringValue(location.ObjectType_) == objectStoreLocationObjectType {
		objectStoreLocation := location.GetValue().(management.ObjectStoreLocation)

		objectStoreLocationMap := make(map[string]interface{})
		objectStoreLocationMap["object_store_location"] = flattenResourceObjectStoreLocation(&objectStoreLocation, d)
		backupTargetLocation = append(backupTargetLocation, objectStoreLocationMap)
		return backupTargetLocation
	}

	return backupTargetLocation
}

func flattenResourceObjectStoreLocation(objectStoreLocation *management.ObjectStoreLocation, d interface{}) []map[string]interface{} {
	if objectStoreLocation == nil {
		return nil
	}

	objectStoreLocationMap := make(map[string]interface{})
	objectStoreLocationMap["provider_config"] = flattenProviderConfig(objectStoreLocation.ProviderConfig)
	objectStoreLocationMap["backup_policy"] = flattenBackupPolicy(objectStoreLocation.BackupPolicy)

	// extract the credentials from the state file since they are not returned in the response
	locationI, ok := d.([]interface{})

	if !ok || len(locationI) == 0 {
		// no previous state, just return flattened response
		return []map[string]interface{}{objectStoreLocationMap}
	}

	location, ok := locationI[0].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{objectStoreLocationMap}
	}

	objectStoreLocationI := location["object_store_location"].([]interface{})[0].(map[string]interface{})
	providerConfig := objectStoreLocationI["provider_config"].([]interface{})[0].(map[string]interface{})
	credentials := providerConfig["credentials"].([]interface{})
	objectStoreLocationMap["provider_config"].([]map[string]interface{})[0]["credentials"] = credentials

	objectStoreLocationList := make([]map[string]interface{}, 0)
	objectStoreLocationList = append(objectStoreLocationList, objectStoreLocationMap)

	return objectStoreLocationList
}

func ResourceNutanixBackupTargetV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	readResp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), nil)
	if err != nil {
		return diag.Errorf("error while fetching backup target: %s", err)
	}

	// extract the etag from the read response
	args := make(map[string]interface{})
	eTag := conn.DomainManagerBackupsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(eTag)

	updateSpec := management.BackupTarget{}

	if d.HasChange("location") {
		oneOfBackupTargetLocation := management.NewOneOfBackupTargetLocation()
		locationI := d.Get("location").([]interface{})
		location := locationI[0].(map[string]interface{})

		if location["cluster_location"] != nil && len(location["cluster_location"].([]interface{})) > 0 {
			clusterLocation := location["cluster_location"].([]interface{})[0].(map[string]interface{})
			clusterConfig := clusterLocation["config"].([]interface{})[0].(map[string]interface{})

			clusterConfigBody := management.NewClusterLocation()
			clusterRef := management.NewClusterReference()

			clusterRef.ExtId = utils.StringPtr(clusterConfig["ext_id"].(string))
			// From IRIS SDK, the cluster location config is a OneOfClusterLocationConfig
			// so we need to set the value of the OneOfClusterLocationConfig
			oneOfClusterLocationConfig := management.NewOneOfClusterLocationConfig()
			oneOfClusterLocationConfig.SetValue(*clusterRef)
			clusterConfigBody.Config = oneOfClusterLocationConfig

			err = oneOfBackupTargetLocation.SetValue(*clusterConfigBody)
			if err != nil {
				return diag.Errorf("error while setting cluster location : %v", err)
			}
		} else if location["object_store_location"] != nil && len(location["object_store_location"].([]interface{})) > 0 {
			aJSON, _ := json.MarshalIndent(location, "", "  ")
			log.Printf("[DEBUG] Object Store Location: %s", string(aJSON))

			objectStoreLocation := location["object_store_location"].([]interface{})[0].(map[string]interface{})
			providerConfig := objectStoreLocation["provider_config"]
			backupPolicy := objectStoreLocation["backup_policy"]

			objectStoreLocationBody := management.NewObjectStoreLocation()

			objectStoreLocationBody.ProviderConfig = expandProviderConfig(providerConfig)
			objectStoreLocationBody.BackupPolicy = expandBackupPolicy(backupPolicy)

			err = oneOfBackupTargetLocation.SetValue(*objectStoreLocationBody)
			if err != nil {
				return diag.Errorf("error while setting object store location : %v", err)
			}
		}

		updateSpec.Location = oneOfBackupTargetLocation
	} else {
		log.Printf("[DEBUG] No changes in backup target Location")
		return nil
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Payload to update backup target: %s", string(aJSON))

	resp, err := conn.DomainManagerBackupsAPIInstance.UpdateBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), &updateSpec, args)

	if err != nil {
		return diag.Errorf("error while updating backup target: %s", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the backup target to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for backup target to be updated: %s", err)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching backup target Task Details: %s", err)
	}

	taskDetails := taskResp.Data.GetValue().(config.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update backup target task details: %s", string(aJSON))

	return ResourceNutanixBackupTargetV2Read(ctx, d, meta)
}

func ResourceNutanixBackupTargetV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	readResp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), nil)
	if err != nil {
		return diag.Errorf("error while fetching backup target: %s", err)
	}

	// extract the etag from the read response
	args := make(map[string]interface{})
	eTag := conn.DomainManagerBackupsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(eTag)

	resp, err := conn.DomainManagerBackupsAPIInstance.DeleteBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), args)

	if err != nil {
		return diag.Errorf("error while deleting backup target: %s", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the backup target to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for backup target to be deleted: %s", err)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching delete backup target task details: %s", err)
	}

	taskDetails := taskResp.Data.GetValue().(config.Task)

	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete backup target task details: %s", string(aJSON))

	return nil
}

func expandProviderConfig(providerConfig interface{}) *management.OneOfObjectStoreLocationProviderConfig {
	if len(providerConfig.([]interface{})) == 0 {
		return nil
	}
	providerConfigI := providerConfig.([]interface{})

	if len(providerConfigI) == 0 {
		return nil
	}

	providerConfigMap := providerConfigI[0].(map[string]interface{})

	providerConfigObj := management.NewOneOfObjectStoreLocationProviderConfig()

	awsS3Config := management.NewAWSS3Config()
	awsS3Config.BucketName = utils.StringPtr(providerConfigMap["bucket_name"].(string))
	awsS3Config.Region = utils.StringPtr(providerConfigMap["region"].(string))
	awsS3Config.Credentials = expandAccessKeyCredentials(providerConfigMap["credentials"])

	if err := providerConfigObj.SetValue(*awsS3Config); err != nil {
		log.Printf("[ERROR] Error while setting AWS S3 config: %v", err)
	}

	return providerConfigObj
}

func expandAccessKeyCredentials(credentials interface{}) *management.AccessKeyCredentials {
	if len(credentials.([]interface{})) == 0 {
		return nil
	}

	credentialsMap := credentials.([]interface{})[0].(map[string]interface{})
	accessKeyCredentials := management.AccessKeyCredentials{
		AccessKeyId:     utils.StringPtr(credentialsMap["access_key_id"].(string)),
		SecretAccessKey: utils.StringPtr(credentialsMap["secret_access_key"].(string)),
	}
	return &accessKeyCredentials
}

func expandBackupPolicy(policy interface{}) *management.BackupPolicy {
	if len(policy.([]interface{})) == 0 {
		return nil
	}

	policyI := policy.([]interface{})[0].(map[string]interface{})

	backupPolicy := management.BackupPolicy{
		RpoInMinutes: utils.IntPtr(policyI["rpo_in_minutes"].(int)),
	}

	return &backupPolicy
}
