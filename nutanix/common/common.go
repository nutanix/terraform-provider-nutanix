package common

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ExpandListOfString to expand a list of interface{}
// which is defined in the schema
// its return a list of string
func ExpandListOfString(list []interface{}) []string {
	stringListStr := make([]string, 0)
	for i, v := range list {
		if v == nil || v == "" {
			log.Printf("[DEBUG] Skipping nil or empty value at index %d", i)
			continue // Skip nil or empty values
		}
		stringListStr = append(stringListStr, v.(string))
	}
	return stringListStr
}

// IsExplicitlySet defined to determine whether a particular key (or configuration attribute) within a Terraform resource configuration has been explicitly set by the user.
// Returns a Boolean (true or false). true indicates that the key was explicitly set with a non-null value; false implies it was either not set, is unknown, or explicitly set to null.
func IsExplicitlySet(d *schema.ResourceData, key string) bool {
	rawConfig := d.GetRawConfig() // Get raw Terraform config as cty.Value
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return false // If rawConfig is null/unknown, key wasn't explicitly set
	}

	// Convert rawConfig to map and check if key exists
	configMap := rawConfig.AsValueMap()
	if val, exists := configMap[key]; exists {
		log.Printf("[DEBUG] Key: %s, Value: %s", key, val)
		return !val.IsNull() // Ensure key exists and isn't explicitly null
	}
	return false
}

func TaskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
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
	return pr.GetName()
}
