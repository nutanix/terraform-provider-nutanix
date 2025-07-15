package passwordmanagerv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
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
				Required: true,
			},
			"new_password": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNutanixPasswordManagerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmt.Printf("[DEBUG] Creating Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	conn := meta.(*conns.Client).ClusterAPI
	extID := utils.StringPtr(d.Get("ext_id").(string))
	body := &clusterConfig.ChangePasswordSpec{}
	if currPassword, ok := d.GetOk("current_password"); ok {
		body.CurrentPassword = utils.StringPtr(currPassword.(string))
	}
	if newPassword, ok := d.GetOk("new_password"); ok {
		body.NewPassword = utils.StringPtr(newPassword.(string))
	}
	resp, err := conn.PasswordManagerAPI.ChangeSystemUserPasswordById(extID, body)
	if err != nil {
		return diag.Errorf("error while performing password change: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI

	// Wait for the PreChecks to be successful
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroup(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Change Password Request failed for ext_id %s with error %s", utils.StringValue(extID), errWaitTask)
	}

	// set the resource id to random uuid
	d.SetId(utils.GenUUID())
	return resourceNutanixPasswordManagerV2Read(ctx, d, meta)
}

func resourceNutanixPasswordManagerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmt.Printf("[DEBUG] Reading Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	return nil
}

func resourceNutanixPasswordManagerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmt.Printf("[DEBUG] Updating Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	return resourceNutanixPasswordManagerV2Create(ctx, d, meta)
	// Note: The update operation is the same as create in this case.
}

func resourceNutanixPasswordManagerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmt.Printf("[DEBUG] Deleting Password Manager V2 resource with ext_id: %s", d.Get("ext_id").(string))
	return nil
}

func taskStateRefreshPrismTaskGroup(client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)
		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(prismConfig.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(pr *prismConfig.TaskStatus) string {
	const two, three, five, six, seven = 2, 3, 5, 6, 7
	if pr != nil {
		if *pr == prismConfig.TaskStatus(six) {
			return "FAILED"
		}
		if *pr == prismConfig.TaskStatus(seven) {
			return "CANCELED"
		}
		if *pr == prismConfig.TaskStatus(two) {
			return "QUEUED"
		}
		if *pr == prismConfig.TaskStatus(three) {
			return "RUNNING"
		}
		if *pr == prismConfig.TaskStatus(five) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}
