package prismv2_test

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const timeout = 3 * time.Minute

func checkAttributeLength(resourceName, attribute string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}

		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}

		if count < minLength {
			return fmt.Errorf("expected %s to be >= %d, got %d", attrKey, minLength, count)
		}

		return nil
	}
}

func checkBackupTargetExist() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

				domainManagerExtID := attributes["domain_manager_ext_id"]
				for i := 0; i < backupTargetsCount; i++ {
					extID := attributes["backup_targets."+strconv.Itoa(i)+".ext_id"]

					readResp, err := client.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(extID), nil)
					if err != nil {
						return fmt.Errorf("error while fetching Backup Target: %s", err)
					}

					// extract the etag from the read response
					args := make(map[string]interface{})
					eTag := client.ApiClient.GetEtag(readResp)
					args["If-Match"] = utils.StringPtr(eTag)

					resp, err := client.DeleteBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(extID), args)

					if err != nil {
						return fmt.Errorf("error while deleting Backup Target: %s", err)
					}
					return waitDeleteTask(resp)
				}

				return nil
			}
		}
		return fmt.Errorf("backup target still exists")
	}
}

func checkLastSyncTimeBackupTarget(retries int, delay time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_backup_target_v2" {
				attributes := rs.Primary.Attributes

				domainManagerExtID := attributes["domain_manager_ext_id"]
				backupTargetExtID := attributes["ext_id"]

				for i := 0; i < retries; i++ {
					readResp, err := client.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(backupTargetExtID), nil)
					if err != nil {
						return fmt.Errorf("error while fetching Backup Target: %s", err)
					}

					backupTarget := readResp.Data.GetValue().(management.BackupTarget)

					fmt.Printf("LastSyncTime: %v\n", backupTarget.LastSyncTime)
					if backupTarget.LastSyncTime != nil {
						fmt.Printf(" Restore Point Created after %d minutes\n", i*30/60)
						return nil
					}
					fmt.Printf("Waiting for 30 seconds to Fetch backup target\n")
					time.Sleep(delay)
				}
			}
		}
		return fmt.Errorf("backup Target restore point not created")
	}
}

func waitDeleteTask(resp *management.DeleteBackupTargetApiResponse) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := conn.PrismAPI
	// Wait for the backup target to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(utils.StringValue(taskUUID)),
		Timeout: timeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for Backup Target to be deleted: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return fmt.Errorf("error while fetching Backup Target Task Details: %s", err)
	}

	rUUID := resourceUUID.Data.GetValue().(config.Task)

	aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Delete Backup Target Task Details: %s", string(aJSON))
	return nil
}

func taskStateRefreshPrismTaskGroupFunc(taskUUID string) resource.StateRefreshFunc {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := conn.PrismAPI.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)

		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(config.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

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
