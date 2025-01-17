package prismv2

import (
	"context"
	"encoding/json"
	"log"

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
)

func ResourceNutanixBackupTargetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixBackupTargetV2Create,
		ReadContext:   ResourceNutanixBackupTargetV2Read,
		UpdateContext: ResourceNutanixBackupTargetV2Update,
		DeleteContext: ResourceNutanixBackupTargetV2Delete,
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
													Type:     schema.TypeString,
													Required: true,
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
																Type:     schema.TypeString,
																Required: true,
															},
															"secret_access_key": {
																Type:     schema.TypeString,
																Required: true,
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

	backupTargetExtID := ""
	if location["cluster_location"] != nil && len(location["cluster_location"].([]interface{})) > 0 {
		clusterLocation := location["cluster_location"].([]interface{})[0].(map[string]interface{})
		clusterConfig := clusterLocation["config"].([]interface{})[0].(map[string]interface{})

		clusterConfigBody := management.NewClusterLocation()
		clusterRef := management.NewClusterReference()

		clusterRef.ExtId = utils.StringPtr(clusterConfig["ext_id"].(string))

		clusterConfigBody.Config = clusterRef

		err := OneOfBackupTargetLocation.SetValue(*clusterConfigBody)
		if err != nil {
			return diag.Errorf("error while setting cluster location : %v", err)
		}
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
	}

	body.Location = OneOfBackupTargetLocation

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Backup Target Body: %s", string(aJSON))

	// Get all the backup targets for the domain manager
	// This is to get the backup target ext id
	//by comparing the backup targets before and after creating a new backup target
	listBackupTargets, err := conn.DomainManagerBackupsAPIInstance.ListBackupTargets(utils.StringPtr(domainManagerExtID))
	if err != nil {
		return diag.Errorf("error while Listing Backup Targets for : %s err: %s", domainManagerExtID, err)
	}
	oldBackupTargets := []management.BackupTarget{}
	if listBackupTargets.Data != nil {
		oldBackupTargets = listBackupTargets.Data.GetValue().([]management.BackupTarget)
	}

	resp, err := conn.DomainManagerBackupsAPIInstance.CreateBackupTarget(utils.StringPtr(domainManagerExtID), &body)

	if err != nil {
		return diag.Errorf("error while Creating Backup Target: %s", err)
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
		return diag.Errorf("error waiting for Backup Target to be created: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Backup Target Task Details: %s", err)
	}

	rUUID := resourceUUID.Data.GetValue().(config.Task)
	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Create Backup Target Task Details: %s", string(aJSON))

	//filter := "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
	listBackupTargets, err = conn.DomainManagerBackupsAPIInstance.ListBackupTargets(utils.StringPtr(domainManagerExtID))
	if err != nil {
		return diag.Errorf("error while Listing Backup Targets for : %s err: %s", domainManagerExtID, err)
	}
	newBackupTargets := listBackupTargets.Data.GetValue().([]management.BackupTarget)

	aJSON, _ = json.MarshalIndent(oldBackupTargets, "", "  ")
	log.Printf("[DEBUG] Old Backup Targets: %s", string(aJSON))

	log.Printf("[DEBUG] ###############################")

	aJSON, _ = json.MarshalIndent(newBackupTargets, "", "  ")
	log.Printf("[DEBUG] New Backup Targets: %s", string(aJSON))

	if len(oldBackupTargets) == 0 {
		backupTargetExtID = utils.StringValue(newBackupTargets[0].ExtId)
	} else {
		// Find the new backup target ext id
		for _, newBackupTarget := range newBackupTargets {
			for _, oldBackupTarget := range oldBackupTargets {
				if utils.StringValue(newBackupTarget.ExtId) != utils.StringValue(oldBackupTarget.ExtId) {
					backupTargetExtID = utils.StringValue(newBackupTarget.ExtId)
					break
				}
			}
		}
	}

	if backupTargetExtID == "" {
		return diag.Errorf("error while fetching Created Backup Target Ext ID")
	}

	d.SetId(backupTargetExtID)

	return ResourceNutanixBackupTargetV2Read(ctx, d, meta)
}

func ResourceNutanixBackupTargetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conn := meta.(*conns.Client).PrismAPI

	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), nil)

	if err != nil {
		return diag.Errorf("error while fetching Backup Target: %s", err)
	}

	backupTarget := resp.Data.GetValue().(management.BackupTarget)

	aJSON, _ := json.MarshalIndent(backupTarget, "", "  ")
	log.Printf("[DEBUG] Read Backup Target Details: %s", string(aJSON))

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
	if err := d.Set("location", flattenBackupTargetLocation(backupTarget.Location)); err != nil {
		return diag.Errorf("error setting location: %s", err)
	}

	return nil
}

func ResourceNutanixBackupTargetV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conn := meta.(*conns.Client).PrismAPI
	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	readResp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), nil)
	if err != nil {
		return diag.Errorf("error while fetching Backup Target: %s", err)
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

			clusterConfigBody.Config = clusterRef

			err = oneOfBackupTargetLocation.SetValue(*clusterConfigBody)
			if err != nil {
				return diag.Errorf("error while setting cluster location : %v", err)
			}
		} else if location["object_store_location"] != nil && len(location["object_store_location"].([]interface{})) > 0 {
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
		log.Printf("[DEBUG] No changes in Backup Target Location")
		return nil
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Backup Target Update Body: %s", string(aJSON))

	resp, err := conn.DomainManagerBackupsAPIInstance.UpdateBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), &updateSpec, args)

	if err != nil {
		return diag.Errorf("error while updating Backup Target: %s", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the backup target to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for Backup Target to be updated: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Backup Target Task Details: %s", err)
	}

	rUUID := resourceUUID.Data.GetValue().(config.Task)

	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Update Backup Target Task Details: %s", string(aJSON))

	return ResourceNutanixBackupTargetV2Read(ctx, d, meta)
}

func ResourceNutanixBackupTargetV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	readResp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), nil)
	if err != nil {
		return diag.Errorf("error while fetching Backup Target: %s", err)
	}

	// extract the etag from the read response
	args := make(map[string]interface{})
	eTag := conn.DomainManagerBackupsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(eTag)

	resp, err := conn.DomainManagerBackupsAPIInstance.DeleteBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(d.Id()), args)

	if err != nil {
		return diag.Errorf("error while deleting Backup Target: %s", err)
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
		return diag.Errorf("error waiting for Backup Target to be deleted: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Delete Backup Target Task Details: %s", err)
	}

	rUUID := resourceUUID.Data.GetValue().(config.Task)

	aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Delete Backup Target Task Details: %s", string(aJSON))

	return nil
}

func expandProviderConfig(providerConfig interface{}) *management.AWSS3Config {
	if len(providerConfig.([]interface{})) == 0 {
		return nil
	}
	providerConfigI := providerConfig.([]interface{})

	if len(providerConfigI) == 0 {
		return nil
	}

	providerConfigMap := providerConfigI[0].(map[string]interface{})

	awsS3Config := management.AWSS3Config{
		BucketName:  utils.StringPtr(providerConfigMap["bucket_name"].(string)),
		Region:      utils.StringPtr(providerConfigMap["region"].(string)),
		Credentials: expandAccessKeyCredentials(providerConfigMap["credentials"]),
	}

	return &awsS3Config
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
