package prismv2

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRestoreSourceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRestoreSourceV2Create,
		ReadContext:   ResourceNutanixRestoreSourceV2Read,
		UpdateContext: ResourceNutanixRestoreSourceV2Update,
		DeleteContext: ResourceNutanixRestoreSourceV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
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

func ResourceNutanixRestoreSourceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Restore Source Create. ID: %s", d.Id())
	conn := meta.(*conns.Client).PrismAPI

	body := management.RestoreSource{}

	oneOfRestoreSourceLocation := management.NewOneOfRestoreSourceLocation()
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

		err := oneOfRestoreSourceLocation.SetValue(*clusterConfigBody)
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

		err := oneOfRestoreSourceLocation.SetValue(*objectStoreLocationBody)
		if err != nil {
			return diag.Errorf("error while setting object store location : %v", err)
		}
	}

	body.Location = oneOfRestoreSourceLocation

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Restore Source Create Payload: %s", string(aJSON))

	resp, err := conn.DomainManagerBackupsAPIInstance.CreateRestoreSource(&body)

	if err != nil {
		return diag.Errorf("error while creating restore source: %s", err)
	}

	restoreSource := resp.Data.GetValue().(management.RestoreSource)

	d.SetId(utils.StringValue(restoreSource.ExtId))

	aJSON, _ = json.MarshalIndent(resp, "", "  ")
	log.Printf("[DEBUG] Restore Source create response: %s", string(aJSON))

	return ResourceNutanixRestoreSourceV2Read(ctx, d, meta)
}

func ResourceNutanixRestoreSourceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	resp, err := conn.DomainManagerBackupsAPIInstance.GetRestoreSourceById(utils.StringPtr(d.Id()))

	if err != nil {
		log.Printf("[DEBUG] Restore Source read error: %s", err)
		errMessage := utils.ExtractErrorFromV4APIResponse(err)
		if strings.Contains(errMessage, "not found") {
			// If the resource is not found, its Auto-Deleted create a new one
			log.Printf("[DEBUG] Restore Source automatically deleted, recreating it")
			return ResourceNutanixRestoreSourceV2Create(ctx, d, meta)
		}
		return diag.Errorf("error while fetching restore source: %s", err)
	}

	restoreSource := resp.Data.GetValue().(management.RestoreSource)

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

	return nil
}

func ResourceNutanixRestoreSourceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Since the resource is auto-deleted, we will recreate it
	return ResourceNutanixRestoreSourceV2Create(ctx, d, meta)
}

func ResourceNutanixRestoreSourceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	readResp, err := conn.DomainManagerBackupsAPIInstance.GetRestoreSourceById(utils.StringPtr(d.Id()))
	if err != nil {
		log.Printf("[DEBUG] Restore Source Read Error: %s", err)
		errMessage := utils.ExtractErrorFromV4APIResponse(err)
		if strings.Contains(errMessage, "not found") {
			// If the resource is not found, its Auto-Deleted create a new one
			log.Printf("[DEBUG] Restore Source Automatically Deleted")
			return nil
		}
		return diag.Errorf("error while fetching Restore Source: %s", err)
	}

	// extract the etag from the read response
	args := make(map[string]interface{})
	eTag := conn.DomainManagerBackupsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(eTag)

	resp, err := conn.DomainManagerBackupsAPIInstance.DeleteRestoreSourceById(utils.StringPtr(d.Id()), args)

	if err != nil {
		return diag.Errorf("error while deleting restore source: %s", err)
	}

	aJSON, _ := json.MarshalIndent(resp, "", "  ")
	log.Printf("[DEBUG] Restore Source delete response: %s", string(aJSON))

	return nil
}
