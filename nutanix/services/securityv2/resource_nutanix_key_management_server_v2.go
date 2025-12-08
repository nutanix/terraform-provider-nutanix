// Package securityv2 provides resources for managing security-related configurations in Nutanix.
package securityv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	securityPrism "github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/security/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixKeyManagementServerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixKeyManagementServerV2Create,
		ReadContext:   ResourceNutanixKeyManagementServerV2Read,
		UpdateContext: ResourceNutanixKeyManagementServerV2Update,
		DeleteContext: ResourceNutanixKeyManagementServerV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Update: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Delete: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_information": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_secret": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(8, 256),
						},
						"credential_expiry_date": {
							Type:     schema.TypeString,
							Required: true,
						},
						"truncated_client_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// Computed Attributes
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

func ResourceNutanixKeyManagementServerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	body := config.NewKeyManagementServer()

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if accessInfo, ok := d.GetOk("access_information"); ok {
		accessInfoVal, expandErr := expandAccessInformation(accessInfo.([]interface{}))
		if expandErr != nil {
			return diag.FromErr(expandErr)
		}
		body.AccessInformation = accessInfoVal
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] key management server payload: %s", aJSON)

	resp, err := conn.KeyManagementServersAPIInstance.CreateKeyManagementServer(body)
	if err != nil {
		return diag.Errorf("error while creating Key Management Server : %v", err)
	}

	taskRefValue, ok := resp.Data.GetValue().(securityPrism.TaskReference)
	if !ok {
		return diag.Errorf("error: unexpected response type from create API, expected TaskReference")
	}
	TaskRef := taskRefValue
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Key Management Server to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for kms to be created: %s", errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Key Management Server Task UUID : %v", err)
	}
	taskDetailsValue, ok := taskResp.Data.GetValue().(prismConfig.Task)
	if !ok {
		return diag.Errorf("error: unexpected response type from task API, expected Task")
	}
	taskDetails := taskDetailsValue

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] create key management server task details: %s", aJSON)

	if len(taskDetails.EntitiesAffected) == 0 {
		return diag.Errorf("error: task completed but no entities affected found in task response")
	}
	kmsExtID := taskDetails.EntitiesAffected[0].ExtId
	if kmsExtID == nil {
		return diag.Errorf("error: task completed but entity ext_id is nil")
	}
	d.SetId(*kmsExtID)
	return ResourceNutanixKeyManagementServerV2Read(ctx, d, meta)
}

func ResourceNutanixKeyManagementServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.GetKeyManagementServerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching key management server : %v", err)
	}

	getRespValue, ok := resp.Data.GetValue().(config.KeyManagementServer)
	if !ok {
		return diag.Errorf("error: unexpected response type from get API, expected KeyManagementServer")
	}
	getResp := getRespValue

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	accessInfo, flattenErr := flattenAccessInformation(getResp.AccessInformation)
	if flattenErr != nil {
		return diag.FromErr(flattenErr)
	}
	if err := d.Set("access_information", accessInfo); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixKeyManagementServerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.GetKeyManagementServerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching key management server : %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.KeyManagementServersAPIInstance.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	// Get current state to preserve values for fields that haven't changed
	updateSpec, ok := resp.Data.GetValue().(config.KeyManagementServer)
	if !ok {
		return diag.Errorf("error: unexpected response type from get API, expected KeyManagementServer")
	}

	// Update name if it has changed
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			updateSpec.Name = utils.StringPtr(v.(string))
		}
	}

	// Update access_information if it has changed
	if d.HasChange("access_information") {
		if v, ok := d.GetOk("access_information"); ok {
			accessInfo, expandErr := expandAccessInformation(v.([]interface{}))
			if expandErr != nil {
				return diag.FromErr(expandErr)
			}
			// The API requires clientId and clientSecret to be present in update requests
			// Get values from raw config as they are not stored in state (sensitive fields)
			rawConfig := d.GetRawConfig()
			configMap := rawConfig.AsValueMap()

			var configClientSecret, configClientId string
			if accessInfoVal, exists := configMap["access_information"]; exists && !accessInfoVal.IsNull() && accessInfoVal.IsKnown() {
				if accessInfoVal.Type().IsListType() || accessInfoVal.Type().IsTupleType() {
					accessInfoList := accessInfoVal.AsValueSlice()
					if len(accessInfoList) > 0 {
						firstItem := accessInfoList[0]
						if firstItem.Type().IsObjectType() {
							itemMap := firstItem.AsValueMap()
							if clientSecretVal, exists := itemMap["client_secret"]; exists && !clientSecretVal.IsNull() && clientSecretVal.IsKnown() {
								configClientSecret = clientSecretVal.AsString()
							}
							if clientIdVal, exists := itemMap["client_id"]; exists && !clientIdVal.IsNull() && clientIdVal.IsKnown() {
								configClientId = clientIdVal.AsString()
							}
						}
					}
				}
			}

			// If client_secret hasn't changed, use the value from raw config
			if !d.HasChange("access_information.0.client_secret") {
				if configClientSecret != "" {
					accessInfo.ClientSecret = utils.StringPtr(configClientSecret)
				}
			} else if accessInfo.ClientSecret != nil && utils.StringValue(accessInfo.ClientSecret) == "" {
				// If it has changed but is empty, use the value from raw config
				if configClientSecret != "" {
					accessInfo.ClientSecret = utils.StringPtr(configClientSecret)
				}
			}
			// client_id
			if !d.HasChange("access_information.0.client_id") && configClientId != "" {
				accessInfo.ClientId = utils.StringPtr(configClientId)
			} else if accessInfo.ClientId != nil && utils.StringValue(accessInfo.ClientId) == "" && configClientId != "" {
				accessInfo.ClientId = utils.StringPtr(configClientId)
			}
			updateSpec.AccessInformation = accessInfo
		}
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] update key management server payload: %s", aJSON)

	updateResp, err := conn.KeyManagementServersAPIInstance.UpdateKeyManagementServerById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating key management server : %v", err)
	}
	taskRefValue, ok := updateResp.Data.GetValue().(securityPrism.TaskReference)
	if !ok {
		return diag.Errorf("error: unexpected response type from update API, expected TaskReference")
	}
	TaskRef := taskRefValue
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the kms to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for kms (%s) to updated: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return ResourceNutanixKeyManagementServerV2Read(ctx, d, meta)
}

func ResourceNutanixKeyManagementServerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.DeleteKeyManagementServerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting key management server : %v", err)
	}
	taskRefValue, ok := resp.Data.GetValue().(securityPrism.TaskReference)
	if !ok {
		return diag.Errorf("error: unexpected response type from delete API, expected TaskReference")
	}
	TaskRef := taskRefValue
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Delete task to be complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for kms (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandAccessInformation(accessInfo []interface{}) (*config.AzureAccessInformation, error) {
	if len(accessInfo) == 0 {
		return nil, fmt.Errorf("access information is required")
	}

	accessInfoVal, ok := accessInfo[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("access information must be a map")
	}

	expiryStr := accessInfoVal["credential_expiry_date"].(string)
	expiryTime, err := time.Parse("2006-01-02", expiryStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credential_expiry_date %q: %w", expiryStr, err)
	}

	return &config.AzureAccessInformation{
		EndpointUrl:          utils.StringPtr(accessInfoVal["endpoint_url"].(string)),
		KeyId:                utils.StringPtr(accessInfoVal["key_id"].(string)),
		TenantId:             utils.StringPtr(accessInfoVal["tenant_id"].(string)),
		ClientId:             utils.StringPtr(accessInfoVal["client_id"].(string)),
		ClientSecret:         utils.StringPtr(accessInfoVal["client_secret"].(string)),
		CredentialExpiryDate: utils.Time(expiryTime),
	}, nil
}
