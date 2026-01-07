package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixBackupTargetV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixBackupTargetV2Read,
		Schema: map[string]*schema.Schema{
			"domain_manager_ext_id": {
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
			"location": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_location": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
													Type:     schema.TypeString,
													Computed: true,
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
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"credentials": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"access_key_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"secret_access_key": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"backup_policy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rpo_in_minutes": {
													Type:     schema.TypeInt,
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
			},
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

func DatasourceNutanixBackupTargetV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	domainManagerExtID := d.Get("domain_manager_ext_id").(string)
	backupTargetExtID := d.Get("ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(backupTargetExtID), nil)

	if err != nil {
		return diag.Errorf("error while fetching Backup Target: %s", err)
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
	if err := d.Set("location", flattenBackupTargetLocation(backupTarget.Location)); err != nil {
		return diag.Errorf("error setting location: %s", err)
	}

	d.SetId(backupTargetExtID)

	return nil
}

func flattenClusterLocation(location management.ClusterLocation) []map[string]interface{} {
	clusterLocation := make([]map[string]interface{}, 0)
	clusterLocationMap := make(map[string]interface{})
	if location.Config != nil {
		clsRef := location.Config.GetValue().(management.ClusterReference)
		clusterLocationMap["config"] = flattenClusterReference(&clsRef)
		clusterLocation = append(clusterLocation, clusterLocationMap)
	}

	return clusterLocation
}

func flattenClusterReference(clusterReference *management.ClusterReference) []map[string]interface{} {
	if clusterReference == nil {
		return nil
	}

	clusterRef := make([]map[string]interface{}, 0)
	clusterRefMap := make(map[string]interface{})
	clusterRefMap["ext_id"] = clusterReference.ExtId
	clusterRefMap["name"] = clusterReference.Name

	clusterRef = append(clusterRef, clusterRefMap)

	return clusterRef
}

func flattenBackupTargetLocation(location *management.OneOfBackupTargetLocation) []map[string]interface{} {
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
		objectStoreLocationMap["object_store_location"] = flattenObjectStoreLocation(objectStoreLocation)
		backupTargetLocation = append(backupTargetLocation, objectStoreLocationMap)
		return backupTargetLocation
	}

	return backupTargetLocation
}

func flattenObjectStoreLocation(objectStoreLocation management.ObjectStoreLocation) []map[string]interface{} {
	objectStoreLocationMap := make(map[string]interface{})
	objectStoreLocationMap["provider_config"] = flattenProviderConfig(objectStoreLocation.ProviderConfig)
	objectStoreLocationMap["backup_policy"] = flattenBackupPolicy(objectStoreLocation.BackupPolicy)

	objectStoreLocationList := make([]map[string]interface{}, 0)
	objectStoreLocationList = append(objectStoreLocationList, objectStoreLocationMap)

	return objectStoreLocationList
}

func flattenProviderConfig(providerConfig *management.OneOfObjectStoreLocationProviderConfig) []map[string]interface{} {
	if providerConfig == nil {
		return nil
	}

	providerConfigMap := make(map[string]interface{})

	awsConfig := providerConfig.GetValue().(management.AWSS3Config)

	providerConfigMap["bucket_name"] = awsConfig.BucketName
	providerConfigMap["region"] = awsConfig.Region
	providerConfigMap["credentials"] = flattenAccessKeyCredentials(awsConfig.Credentials)

	providerConfigList := make([]map[string]interface{}, 0)
	providerConfigList = append(providerConfigList, providerConfigMap)

	return providerConfigList
}

func flattenAccessKeyCredentials(credentials *management.AccessKeyCredentials) []map[string]interface{} {
	if credentials == nil {
		return nil
	}

	credentialsMap := make(map[string]interface{})
	credentialsMap["access_key_id"] = credentials.AccessKeyId
	credentialsMap["secret_access_key"] = credentials.SecretAccessKey

	credentialsList := make([]map[string]interface{}, 0)
	credentialsList = append(credentialsList, credentialsMap)

	return credentialsList
}

func flattenBackupPolicy(policy *management.BackupPolicy) []map[string]interface{} {
	if policy == nil {
		return nil
	}

	backupPolicyMap := make(map[string]interface{})
	backupPolicyMap["rpo_in_minutes"] = policy.RpoInMinutes

	backupPolicyList := make([]map[string]interface{}, 0)
	backupPolicyList = append(backupPolicyList, backupPolicyMap)

	return backupPolicyList
}
