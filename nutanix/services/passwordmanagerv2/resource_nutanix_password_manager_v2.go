package passwordmanagerv2

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixPasswordManagerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixPasswordManagerV2Create,
		ReadContext:   resourceNutanixPasswordManagerV2Read,
		UpdateContext: resourceNutanixPasswordManagerV2Update,
		DeleteContext: resourceNutanixPasswordManagerV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"current_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"new_password": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNutanixPasswordManagerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	conn := meta.(*conns.Client).ClusterAPI
	extID := utils.StringPtr(d.Get("ext_id").(string))
	body := &clusterConfig.ChangePasswordSpec{}
	if currPassword, ok := d.GetOk("current_password"); ok {
		body.CurrentPassword = utils.StringPtr(currPassword.(string))
	}
	if newPassword, ok := d.GetOk("new_password"); ok {
		body.NewPassword = utils.StringPtr(newPassword.(string))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Change Password Request body: %s", aJSON)

	resp, err := conn.PasswordManagerAPI.ChangeSystemUserPasswordById(extID, body)
	if err != nil {
		return diag.Errorf("error while performing password change: %v", err)
	}

	aJSON, _ = json.MarshalIndent(resp, "", "  ")
	log.Printf("[DEBUG] Change Password Response: %s", aJSON)

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	_, taskErr := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if taskErr != nil {
		log.Printf("[DEBUG] Checking if the error is due to password change for the user configured in the provider configuration")
		// if the error is 401, it means the password has been changed for the
		// user configured in the provider configuration, so we need to set the new password
		// and retry to fetch the task
		// Convert the struct to JSON bytes
		rawJSON, _ := json.Marshal(taskErr)

		var taskErrMap map[string]interface{}

		if unmarshalErr := json.Unmarshal(rawJSON, &taskErrMap); unmarshalErr == nil {
			log.Printf("[DEBUG]  Error message: %s", taskErrMap["Message"])
			if status, ok := taskErrMap["Status"].(string); ok && strings.Contains(status, "401") {
				log.Printf("[DEBUG]  Status code is %s, indicating a 401 Unauthorized error", taskErrMap["Status"])
				// password changed for the user configured in the provider configuration
				// set the task conn client password to the new password and retry to fetch the task

				newCredentials := client.Credentials{
					Username: taskconn.TaskRefAPI.ApiClient.Username,
					Password: utils.StringValue(body.NewPassword),
					Endpoint: taskconn.TaskRefAPI.ApiClient.Host,
					Port:     strconv.Itoa(taskconn.TaskRefAPI.ApiClient.Port),
				}

				newPrismClient, prismErr := prism.NewPrismClient(newCredentials)
				if prismErr != nil {
					return diag.Errorf("error while creating new prism client: %v", prismErr)
				}

				// retry to fetch the task
				_, taskErr = newPrismClient.TaskRefAPI.GetTaskById(taskUUID, nil)
				if taskErr != nil {
					return diag.Errorf("error while fetching task by ID %s: %v", utils.StringValue(taskUUID), taskErr)
				}

				// Wait for the password change to complete
				stateConf := &resource.StateChangeConf{
					Pending: []string{"PENDING", "RUNNING", "QUEUED"},
					Target:  []string{"SUCCEEDED"},
					Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, newPrismClient, utils.StringValue(taskUUID)),
					Timeout: d.Timeout(schema.TimeoutCreate),
				}

				if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
					return diag.Errorf("error waiting for password change (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
				}

				// Get task details for logging
				taskResp, err := newPrismClient.TaskRefAPI.GetTaskById(taskUUID, nil)
				if err != nil {
					return diag.Errorf("error while fetching password change task: %v", err)
				}
				taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
				aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
				log.Printf("[DEBUG] Create Password Manager Task Details: %s", string(aJSON))

				// This is an action resource that does not maintain state.
				// The resource ID is set to the task ExtId for traceability.
				d.SetId(utils.StringValue(taskDetails.ExtId))
				return resourceNutanixPasswordManagerV2Read(ctx, d, meta)
			}
			return diag.Errorf("error while fetching task by ID %s: %v", utils.StringValue(taskUUID), taskErr)
		}
	}

	// The password change is not for the user configured in the provider configuration
	// Wait for the password change to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for password change (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details for logging
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching password change task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Password Manager Task Details: %s", string(aJSON))

	// This is an action resource that does not maintain state.
	// The resource ID is set to the task ExtId for traceability.
	d.SetId(utils.StringValue(taskDetails.ExtId))
	return resourceNutanixPasswordManagerV2Read(ctx, d, meta)
}

func resourceNutanixPasswordManagerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	return nil
}

func resourceNutanixPasswordManagerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	return resourceNutanixPasswordManagerV2Create(ctx, d, meta)
	// Note: The update operation is the same as create in this case.
}

func resourceNutanixPasswordManagerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	return nil
}
