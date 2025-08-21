package clustersv2_test

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clusterPrism "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const timeout = 3 * time.Minute

func associateCategoryToCluster() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Println("Associating category with cluster")
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.ClusterAPI.ClusterEntityAPI

		clusterExtID := ""
		categoryExtID := ""

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_cluster_v2" {
				clusterExtID = rs.Primary.ID
			}
			if rs.Type == "nutanix_category_v2" {
				categoryExtID = rs.Primary.ID
			}
		}

		if clusterExtID == "" || categoryExtID == "" {
			return fmt.Errorf("cluster or category not found in state")
		}

		log.Printf("[DEBUG] Associating category: %s to cluster: %s", categoryExtID, clusterExtID)

		body := config.NewCategoryEntityReferences()

		body.Categories = append(body.Categories, categoryExtID)

		aJSON, _ := json.MarshalIndent(body, "", "  ")
		log.Printf("[DEBUG] Category body: %s", aJSON)

		resp, err := client.AssociateCategoriesToCluster(utils.StringPtr(clusterExtID), body)
		if err != nil {
			return fmt.Errorf("Error associating category to cluster: %v", err)
		}

		TaskRef := resp.Data.GetValue().(clusterPrism.TaskReference)
		taskUUID := TaskRef.ExtId

		taskconn := conn.PrismAPI
		// Wait for the backup target to be deleted
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: taskStateRefreshPrismTaskGroupFunc(utils.StringValue(taskUUID)),
			Timeout: timeout,
		}

		if _, taskErr := stateConf.WaitForState(); taskErr != nil {
			return fmt.Errorf("Error waiting for category association task to complete: %s", taskErr)
		}

		_, err = taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return fmt.Errorf("Error while fetching Category Association Task Details: %s", err)
		}

		return nil
	}
}

func disassociateCategoryFromCluster() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Println("Disassociating category from cluster")
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.ClusterAPI.ClusterEntityAPI

		clusterExtID := ""
		categoryExtID := ""

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_cluster_v2" {
				clusterExtID = rs.Primary.ID
			}
			if rs.Type == "nutanix_category_v2" {
				categoryExtID = rs.Primary.ID
			}
		}

		if clusterExtID == "" || categoryExtID == "" {
			return fmt.Errorf("cluster or category not found in state")
		}

		log.Printf("[DEBUG] Disassociating category: %s from cluster: %s", categoryExtID, clusterExtID)

		body := config.NewCategoryEntityReferences()

		body.Categories = append(body.Categories, categoryExtID)

		aJSON, _ := json.MarshalIndent(body, "", "  ")
		log.Printf("[DEBUG] Category body: %s", aJSON)

		resp, err := client.DisassociateCategoriesFromCluster(utils.StringPtr(clusterExtID), body)
		if err != nil {
			return fmt.Errorf("Error disassociating category from cluster: %v", err)
		}

		TaskRef := resp.Data.GetValue().(clusterPrism.TaskReference)
		taskUUID := TaskRef.ExtId

		taskconn := conn.PrismAPI
		// Wait for the backup target to be deleted
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: taskStateRefreshPrismTaskGroupFunc(utils.StringValue(taskUUID)),
			Timeout: timeout,
		}

		if _, taskErr := stateConf.WaitForState(); taskErr != nil {
			return fmt.Errorf("Error waiting for category disassociation task to complete: %s", taskErr)
		}

		_, err = taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return fmt.Errorf("Error while fetching Category Disassociation Task Details: %s", err)
		}

		return nil
	}
}

// helper function to check the delete task
func taskStateRefreshPrismTaskGroupFunc(taskUUID string) resource.StateRefreshFunc {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := conn.PrismAPI.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)

		if err != nil {
			return "", "", (fmt.Errorf("Error while polling prism task: %v", err))
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

// helper function to flatten the task status to string
func getTaskStatus(pr *prismConfig.TaskStatus) string {
	if pr != nil {
		const QUEUED, RUNNING, SUCCEEDED, FAILED, CANCELED = 2, 3, 5, 6, 7
		if *pr == prismConfig.TaskStatus(FAILED) {
			return "FAILED"
		}
		if *pr == prismConfig.TaskStatus(CANCELED) {
			return "CANCELED"
		}
		if *pr == prismConfig.TaskStatus(QUEUED) {
			return "QUEUED"
		}
		if *pr == prismConfig.TaskStatus(RUNNING) {
			return "RUNNING"
		}
		if *pr == prismConfig.TaskStatus(SUCCEEDED) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}
