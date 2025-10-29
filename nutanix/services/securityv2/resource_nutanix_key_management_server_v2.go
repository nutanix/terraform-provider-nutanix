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
							Type:     schema.TypeString,
							Required: true,
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

	TaskRef := resp.Data.GetValue().(securityPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Image to be available
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
		return diag.Errorf("error while fetching Image UUID : %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] create key management server task details: %s", aJSON)

	kmsExtID := taskDetails.EntitiesAffected[0].ExtId
	d.SetId(*kmsExtID)
	return ResourceNutanixKeyManagementServerV2Read(ctx, d, meta)
}

func ResourceNutanixKeyManagementServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.GetKeyManagementServerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching key management server : %v", err)
	}

	getResp := resp.Data.GetValue().(config.KeyManagementServer)

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

	updateSpec := config.KeyManagementServer{}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			updateSpec.Name = utils.StringPtr(v.(string))
		}
	}
	if d.HasChange("access_information") {
		if v, ok := d.GetOk("access_information"); ok {
			accessInfo, expandErr := expandAccessInformation(v.([]interface{}))
			if expandErr != nil {
				return diag.FromErr(expandErr)
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
	TaskRef := updateResp.Data.GetValue().(securityPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the kms to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
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
	TaskRef := resp.Data.GetValue().(securityPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the Delete task to be complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for kms (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func expandAccessInformation(accessInfo []interface{}) (*config.AzureAccessInformation, error) {
	if len(accessInfo) == 0 {
		log.Printf("[DEBUG] access information is nil or empty")
		return nil, fmt.Errorf("access information is nil or empty")
	}

	accessInfoVal := accessInfo[0].(map[string]interface{})

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
