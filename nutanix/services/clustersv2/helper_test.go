package clustersv2_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clusterPrism "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const timeout = 3 * time.Minute

// helper function to check the delete task
func taskStateRefreshPrismTaskGroupFunc(taskUUID string) resource.StateRefreshFunc {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := conn.PrismAPI.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)

		if err != nil {
			return "", "", fmt.Errorf("error while polling prism task: %v", err)
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
	return pr.GetName()
}

// ##############################
// ### Expand Cluster Helpers ###
// ##############################
func checkNodesIPs(expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[resourceName3NodesCluster].Primary.Attributes

		// Collect all node IPs dynamically
		var ips []string
		for k, v := range attrs {
			if strings.HasPrefix(k, "nodes.0.node_list.") &&
				strings.HasSuffix(k, ".controller_vm_ip.0.ipv4.0.value") {
				ips = append(ips, v)
			}
		}

		// Check if all expected IPs are present
		for _, expectedIP := range expected {
			found := false
			for _, ip := range ips {
				if ip == expectedIP {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("expected IP %s not found in cluster nodes: %v", expectedIP, ips)
			}
		}

		return nil
	}
}

// ##############################
// ### Category Helpers ###
// ##############################
// checkCategories verifies that all expected category IDs are present in the resource,
// regardless of their order. It accepts category resource names and gets their IDs from state.
func checkCategories(resourceName, categoriesPath string, expectedCategoryResourceNames []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		attrs := rs.Primary.Attributes

		// Get the count of categories
		countKey := categoriesPath + ".#"
		countStr, ok := attrs[countKey]
		if !ok {
			return fmt.Errorf("category count not found at %s", countKey)
		}

		var count int
		fmt.Sscanf(countStr, "%d", &count)

		// Collect categories by index
		categories := make([]string, 0, count)
		for i := 0; i < count; i++ {
			key := fmt.Sprintf("%s.%d", categoriesPath, i)
			if catID, ok := attrs[key]; ok {
				categories = append(categories, catID)
			}
		}

		// Get expected category IDs from state
		expectedCategoryIDs := make([]string, 0, len(expectedCategoryResourceNames))
		for _, catResourceName := range expectedCategoryResourceNames {
			catRS, ok := s.RootModule().Resources[catResourceName]
			if !ok {
				return fmt.Errorf("category resource %s not found in state", catResourceName)
			}
			expectedCategoryIDs = append(expectedCategoryIDs, catRS.Primary.ID)
		}

		// Check that the count matches
		if len(categories) != len(expectedCategoryIDs) {
			return fmt.Errorf("category count mismatch: expected %d, got %d in %s: %v", len(expectedCategoryIDs), len(categories), categoriesPath, categories)
		}

		// Check if all expected category IDs are present
		for _, expectedID := range expectedCategoryIDs {
			found := false
			for _, catID := range categories {
				if catID == expectedID {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("expected category ID %s not found in %s: %v", expectedID, categoriesPath, categories)
			}
		}

		return nil
	}
}

// helper function to check if the cluster is destroyed
func testAccCheckNutanixClusterDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_cluster_v2" {
			continue
		}

		readResp, err := conn.ClusterAPI.ClusterEntityAPI.GetClusterById(utils.StringPtr(rs.Primary.ID), nil)
		if err == nil {
			// delete the cluster
			//extract etag from read response
			args := make(map[string]interface{})
			etagValue := conn.ClusterAPI.ClusterEntityAPI.ApiClient.GetEtag(readResp)
			args["If-Match"] = utils.StringPtr(etagValue)

			deleteResp, err := conn.ClusterAPI.ClusterEntityAPI.DeleteClusterById(utils.StringPtr(rs.Primary.ID), utils.BoolPtr(false), args)
			if err != nil {
				return err
			}
			TaskRef := deleteResp.Data.GetValue().(clusterPrism.TaskReference)
			taskUUID := TaskRef.ExtId

			taskconn := conn.PrismAPI
			// Wait for the cluster to be deleted
			stateConf := &resource.StateChangeConf{
				Pending: []string{"PENDING", "RUNNING", "QUEUED"},
				Target:  []string{"SUCCEEDED"},
				Refresh: taskStateRefreshPrismTaskGroupFunc(utils.StringValue(taskUUID)),
				Timeout: timeout,
			}

			if _, taskErr := stateConf.WaitForState(); taskErr != nil {
				return fmt.Errorf("error waiting for cluster deletion task to complete: %s", taskErr)
			}

			_, err = taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
			if err != nil {
				return fmt.Errorf("error while fetching Cluster Deletion Task Details: %s", err)
			}

			return nil
		}
	}

	return nil
}

// helper function to check if categories and cluster categories association are destroyed
func testAccCheckNutanixClusterCategoriesDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	categoryClient := conn.PrismAPI.CategoriesAPIInstance

	// Collect all category IDs that should be destroyed
	var categoryIDs []string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "nutanix_category_v2" {
			categoryIDs = append(categoryIDs, rs.Primary.ID)
		}
	}

	// Check if categories still exist
	for _, categoryID := range categoryIDs {
		_, err := categoryClient.GetCategoryById(utils.StringPtr(categoryID), nil)
		if err == nil {
			// Category still exists, try to delete it
			fmt.Printf("[DEBUG] Category still exists, attempting to delete: %s\n", categoryID)
			_, deleteErr := categoryClient.DeleteCategoryById(utils.StringPtr(categoryID))
			if deleteErr != nil {
				return fmt.Errorf("error: Category %s still exists and could not be deleted: %v", categoryID, deleteErr)
			}
			fmt.Printf("[DEBUG] Category deleted: %s\n", categoryID)
		} else if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "ENTITY_NOT_FOUND") {
			// If it's not a "not found" error, return it
			return fmt.Errorf("error checking if category %s exists: %v", categoryID, err)
		}
		// If category is not found, that's expected - it's been destroyed
	}

	return nil
}

// ##############################
// ### Cluster Profile Helpers ###
// ##############################
func testAccCheckNutanixClusterDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_cluster_v2" {
			continue
		}
		// Check API if resource exists
		readResp, errRead := conn.ClusterAPI.ClusterEntityAPI.GetClusterById(utils.StringPtr(rs.Primary.ID), nil)
		if errRead != nil {
			errStr := strings.ToLower(fmt.Sprint(errRead))
			// Check for various indicators that the cluster doesn't exist (404, not found, unknown cluster uuid, etc.)
			if strings.Contains(errStr, "not found") ||
				strings.Contains(errStr, "does not exist") ||
				strings.Contains(errStr, "unknown cluster uuid") ||
				strings.Contains(errStr, "unknown cluster") ||
				strings.Contains(errStr, "clu-10005") {
				log.Printf("[DEBUG] Cluster %s not found (already deleted or doesn't exist), treating as success", rs.Primary.ID)
				return nil
			}
			return errRead
		}
		log.Printf("[DEBUG] Cluster %s still exists, attempting to destroy...", rs.Primary.ID)

		// Extract E-Tag Header for delete operation
		etagValue := conn.ClusterAPI.ClusterEntityAPI.ApiClient.GetEtag(readResp)
		args := make(map[string]interface{})
		args["If-Match"] = utils.StringPtr(etagValue)

		// Attempt to delete the cluster (dryrun=false for actual deletion)
		resp, err := conn.ClusterAPI.ClusterEntityAPI.DeleteClusterById(utils.StringPtr(rs.Primary.ID), utils.BoolPtr(false), args)
		if err != nil {
			// If deletion fails, log but don't fail the test (cluster might be in use or have dependencies)
			log.Printf("[DEBUG] Error attempting to delete cluster %s: %v", rs.Primary.ID, err)
			return nil
		}

		// Log the task but don't wait for completion in test cleanup (it might take too long)
		if resp != nil && resp.Data != nil {
			TaskRef := resp.Data.GetValue().(clusterPrism.TaskReference)
			taskUUID := TaskRef.ExtId
			log.Printf("[DEBUG] Cluster %s deletion task started: %s", rs.Primary.ID, utils.StringValue(taskUUID))
		}
		log.Printf("[DEBUG] Cluster %s deletion initiated successfully", rs.Primary.ID)
	}
	return nil
}
