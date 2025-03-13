package prismv2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// func to check pc task status, and return the task status or error message
func taskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		taskResp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)

		if err != nil {
			return "", "", fmt.Errorf("error while polling prism task: %v", err)
		}

		// get the group results
		v := taskResp.Data.GetValue().(config.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

// func to flatten the task status to string
func getTaskStatus(pr *config.TaskStatus) string {
	if pr != nil {
		const QUEUED, RUNNING, SUCCEEDED, FAILED, CANCELED = 2, 3, 5, 6, 7
		if *pr == config.TaskStatus(FAILED) {
			return "FAILED"
		}
		if *pr == config.TaskStatus(CANCELED) {
			return "CANCELED"
		}
		if *pr == config.TaskStatus(QUEUED) {
			return "QUEUED"
		}
		if *pr == config.TaskStatus(RUNNING) {
			return "RUNNING"
		}
		if *pr == config.TaskStatus(SUCCEEDED) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}

// func to flatten the time to string
func flattenTime(time *time.Time) *string {
	if time == nil {
		return nil
	}
	return utils.StringPtr(time.String())
}

// schemas for links
func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
