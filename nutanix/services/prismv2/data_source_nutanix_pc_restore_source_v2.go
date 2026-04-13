package prismv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRestoreSourceV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRestoreSourceV2Read,
		Schema: map[string]*schema.Schema{
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
												//"name": {
												//	Type:     schema.TypeString,
												//	Computed: true,
												//},
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
		},
	}
}

func DatasourceNutanixRestoreSourceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	restoreSourceExtID := d.Get("ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.GetRestoreSourceById(utils.StringPtr(restoreSourceExtID), nil)

	if err != nil {
		return diag.Errorf("error while fetching Restore Source: %s", err)
	}

	restoreSource := resp.Data.GetValue().(management.RestoreSource)

	aJSON, _ := json.MarshalIndent(restoreSource, "", "  ")
	log.Printf("[DEBUG] Restore Source Read Response: %s", string(aJSON))

	if err := d.Set("tenant_id", restoreSource.TenantId); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("ext_id", restoreSource.ExtId); err != nil {
		return diag.Errorf("error setting ext_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(restoreSource.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("location", flattenRestoreSourceLocation(restoreSource.Location)); err != nil {
		return diag.Errorf("error setting location: %s", err)
	}

	d.SetId(utils.StringValue(restoreSource.ExtId))
	return nil
}

func flattenRestoreSourceLocation(location *management.OneOfRestoreSourceLocation) []map[string]interface{} {
	if location == nil {
		return nil
	}

	restoreSourceLocation := make([]map[string]interface{}, 0)

	if utils.StringValue(location.ObjectType_) == clustersLocationObjectType {
		clusterLocation := location.GetValue().(management.ClusterLocation)

		clusterLocationMap := make(map[string]interface{})
		clusterLocationMap["cluster_location"] = flattenRestoreSourceClusterLocation(clusterLocation)
		restoreSourceLocation = append(restoreSourceLocation, clusterLocationMap)
		return restoreSourceLocation
	}

	if utils.StringValue(location.ObjectType_) == objectStoreLocationObjectType {
		objectStoreLocation := location.GetValue().(management.ObjectStoreLocation)

		objectStoreLocationMap := make(map[string]interface{})
		objectStoreLocationMap["object_store_location"] = flattenRestoreSourceObjectStoreLocation(objectStoreLocation)
		restoreSourceLocation = append(restoreSourceLocation, objectStoreLocationMap)
		return restoreSourceLocation
	}

	return restoreSourceLocation
}

func flattenRestoreSourceClusterLocation(location management.ClusterLocation) []map[string]interface{} {
	clusterLocation := make([]map[string]interface{}, 0)
	clusterLocationMap := make(map[string]interface{})
	// From IRIS SDK, the cluster location config is a OneOfClusterLocationConfig
	// so we need to get the value of the OneOfClusterLocationConfig
	clusterConfig := location.Config.GetValue().(management.OneOfClusterLocationConfig)
	clusterConfigValue := clusterConfig.GetValue().(management.ClusterReference)
	clusterLocationMap["config"] = flattenRestoreSourceClusterReference(&clusterConfigValue)

	clusterLocation = append(clusterLocation, clusterLocationMap)

	return clusterLocation
}

func flattenRestoreSourceClusterReference(clusterReference *management.ClusterReference) []map[string]interface{} {
	if clusterReference == nil {
		return nil
	}

	clusterRef := make([]map[string]interface{}, 0)
	clusterRefMap := make(map[string]interface{})
	clusterRefMap["ext_id"] = clusterReference.ExtId
	//clusterRefMap["name"] = clusterReference.Name

	clusterRef = append(clusterRef, clusterRefMap)

	return clusterRef
}

func flattenRestoreSourceObjectStoreLocation(objectStoreLocation management.ObjectStoreLocation) []map[string]interface{} {
	objectStoreLocationMap := make(map[string]interface{})
	objectStoreLocationMap["provider_config"] = flattenProviderConfig(objectStoreLocation.ProviderConfig)
	//objectStoreLocationMap["backup_policy"] = flattenBackupPolicy(objectStoreLocation.BackupPolicy)

	objectStoreLocationList := make([]map[string]interface{}, 0)
	objectStoreLocationList = append(objectStoreLocationList, objectStoreLocationMap)

	return objectStoreLocationList
}
