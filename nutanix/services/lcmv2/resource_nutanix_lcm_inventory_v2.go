package lcmv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixLcmPerformInventoryV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixLcmPerformInventoryV2Create,
		ReadContext:   ResourceNutanixLcmPerformInventoryV2Read,
		UpdateContext: ResourceNutanixLcmPerformInventoryV2Update,
		DeleteContext: ResourceNutanixLcmPerformInventoryV2Delete,
		Schema: map[string]*schema.Schema{
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ResourceNutanixLcmPerformInventoryV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterExtID := d.Get("x_cluster_id").(string)
	var clusterID *string
	if clusterExtID != "" {
		clusterID = utils.StringPtr(clusterExtID)
	} else {
		clusterID = nil
	}
	// pass nil for the body as it is not required and its implemented in hercules Sdk
	// it will be implemented in the future releases of terraform
	resp, err := conn.LcmInventoryAPIInstance.PerformInventory(nil, clusterID, nil)
	if err != nil {
		return diag.Errorf("error while performing the inventory: %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the inventory to be successful
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroup(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Perform inventory task failed: %s", errWaitTask)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching the Lcm inventory task : %v", err)
	}

	task := resourceUUID.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(task, "", "  ")
	log.Printf("[DEBUG] Perform Inventory Task Response: %s", string(aJSON))

	// randomly generating the id
	d.SetId(utils.GenUUID())
	return nil
}

func ResourceNutanixLcmPerformInventoryV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixLcmPerformInventoryV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixLcmPerformInventoryV2Create(ctx, d, meta)
}

func ResourceNutanixLcmPerformInventoryV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func taskStateRefreshPrismTaskGroup(client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
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
