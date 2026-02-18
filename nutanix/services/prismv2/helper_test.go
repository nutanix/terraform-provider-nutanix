package prismv2_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	vmConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const timeout = 3 * time.Minute
const (
	awsS3ConfigObjectType = "prism.v4.management.AWSS3Config"
)

// checkAttributeLength checks the length of an attribute and make sure it is greater than or equal to minLength
// simply used to check the length of a list returned by List data sources
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

// checkBackupTargetExist checks if the backup target exists
// and deletes it if it does
func checkClusterLocationBackupTargetExistAndDeleteIfExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		outputClusterExtID, ok := s.RootModule().Outputs["clusterExtID"]
		if !ok {
			return fmt.Errorf("output 'clusterExtID' not found")
		}

		clusterExtID := outputClusterExtID.Value.(string)

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_pc_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

				domainManagerExtID := attributes["domain_manager_ext_id"]
				// If backup target exists, delete it
				for i := 0; i < backupTargetsCount; i++ {
					attributes := rs.Primary.Attributes

					backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

					for i := 0; i < backupTargetsCount; i++ {
						clusterLocationCount, _ := strconv.Atoi(attributes["backup_targets."+strconv.Itoa(i)+".location.0.cluster_location.#"])

						if clusterLocationCount > 0 {
							clusterLocationExtID := attributes["backup_targets."+strconv.Itoa(i)+".location.0.cluster_location.0.config.0.ext_id"]

							// delete the backup target with the same cluster location ext_id
							if clusterLocationExtID == clusterExtID {
								log.Printf("[DEBUG] cluster location backup target already exists, ext_id: %s", attributes["backup_targets."+strconv.Itoa(i)+".ext_id"])
								backupTargetExtID := attributes["backup_targets."+strconv.Itoa(i)+".ext_id"]
								readResp, err := client.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(backupTargetExtID), nil)
								if err != nil {
									return fmt.Errorf("error while fetching Backup Target: %s", err)
								}

								// extract the etag from the read response
								args := make(map[string]interface{})
								eTag := client.ApiClient.GetEtag(readResp)
								args["If-Match"] = utils.StringPtr(eTag)

								resp, err := client.DeleteBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(backupTargetExtID), args)

								if err != nil {
									return fmt.Errorf("error while deleting Backup Target: %s", err)
								}
								// wait for the backup target to be deleted
								// if the task is not successful, return the error
								return waitDeleteTask(resp)
							}
						}
					}
				}
				return nil
			}
		}
		return fmt.Errorf("backup target still exists")
	}
}

// checkBackupTargetExist checks if the backup target exists
// and deletes it if it does
func checkObjectStoreLocationBackupTargetExistAndDeleteIfExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_pc_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

				domainManagerExtID := attributes["domain_manager_ext_id"]
				// If backup target exists, delete it
				for i := 0; i < backupTargetsCount; i++ {
					attributes := rs.Primary.Attributes

					backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

					for i := 0; i < backupTargetsCount; i++ {
						objectStoreLocationCount, _ := strconv.Atoi(attributes["backup_targets."+strconv.Itoa(i)+".location.0.object_store_location.#"])

						if objectStoreLocationCount > 0 {
							log.Printf("[DEBUG] object store location backup target already exists, ext_id: %s", attributes["backup_targets."+strconv.Itoa(i)+".ext_id"])
							backupTargetExtID := attributes["backup_targets."+strconv.Itoa(i)+".ext_id"]
							readResp, err := client.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(backupTargetExtID), nil)
							if err != nil {
								return fmt.Errorf("error while fetching Backup Target: %s", err)
							}

							// extract the etag from the read response
							args := make(map[string]interface{})
							eTag := client.ApiClient.GetEtag(readResp)
							args["If-Match"] = utils.StringPtr(eTag)

							resp, err := client.DeleteBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(backupTargetExtID), args)

							if err != nil {
								return fmt.Errorf("error while deleting Backup Target: %s", err)
							}
							// wait for the backup target to be deleted
							// if the task is not successful, return the error
							waitDeleteTask(resp)
						}
					}
				}
				return nil
			}
		}
		return fmt.Errorf("backup target still exists")
	}
}

// checkBackupTargetExistAndCreateIfNotExists checks if the cluster location backup target exists
// and creates a new one if it does not
func checkClusterLocationBackupTargetExistAndCreateIfNotExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		outputClusterExtID, ok := s.RootModule().Outputs["clusterExtID"]
		if !ok {
			return fmt.Errorf("output 'clusterExtID' not found")
		}

		clusterExtID := outputClusterExtID.Value.(string)

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_pc_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])
				for i := 0; i < backupTargetsCount; i++ {
					clusterLocationCount, _ := strconv.Atoi(attributes["backup_targets."+strconv.Itoa(i)+".location.0.cluster_location.#"])
					if clusterLocationCount > 0 {
						clusterLocationExtID := attributes["backup_targets."+strconv.Itoa(i)+".location.0.cluster_location.0.config.0.ext_id"]

						if clusterLocationExtID == clusterExtID {
							log.Printf("[DEBUG] Backup Target already exists, ext_id: %s", attributes["backup_targets.0.ext_id"])
							return nil
						}
					}
				}
				log.Printf("[DEBUG] Backup Target not found, creating new Backup Target")
				break
			}
		}

		// Extract the output value for use in later steps
		outputDomainManagerExtID, ok := s.RootModule().Outputs["domainManagerExtID"]
		if !ok {
			return fmt.Errorf("output 'domainManagerExtID' not found")
		}

		domainManagerExtID := outputDomainManagerExtID.Value.(string)

		// Create Backup Target
		body := management.BackupTarget{}

		OneOfBackupTargetLocation := management.NewOneOfBackupTargetLocation()

		clusterConfigBody := management.NewClusterLocation()
		clusterRef := management.NewClusterReference()

		clusterRef.ExtId = utils.StringPtr(clusterExtID)

		oneOfClusterLocationConfig := management.NewOneOfClusterLocationConfig()
		oneOfClusterLocationConfig.SetValue(*clusterRef)
		clusterConfigBody.Config = oneOfClusterLocationConfig

		err := OneOfBackupTargetLocation.SetValue(*clusterConfigBody)
		if err != nil {
			return fmt.Errorf("error while setting cluster location : %v", err)
		}

		body.Location = OneOfBackupTargetLocation

		resp, err := client.CreateBackupTarget(utils.StringPtr(domainManagerExtID), &body)

		if err != nil {
			return fmt.Errorf("error while Creating Backup Target: %s", err)
		}

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
		//nolint:staticcheck
		if _, taskErr := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error waiting for Backup Target to be deleted: %s", taskErr)
		}

		_, err = taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target Task Details: %s", err)
		}

		listResp, err := client.ListBackupTargets(utils.StringPtr(domainManagerExtID), nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target: %s", err)
		}
		backupTargets := listResp.Data.GetValue().([]management.BackupTarget)

		// Find the new backup target ext id
		for _, backupTarget := range backupTargets {
			backupTargetLocation := backupTarget.Location
			if utils.StringValue(backupTargetLocation.ObjectType_) == "prism.v4.management.ClusterLocation" {
				clusterLocation := backupTarget.Location.GetValue().(management.ClusterLocation)
				clusterConfig := clusterLocation.Config.GetValue().(management.ClusterReference)
				if utils.StringValue(clusterConfig.ExtId) == clusterExtID {
					break
				}
			}
		}

		return nil
	}
}

// checkBackupTargetExistAndCreateIfNotExists checks if the cluster location backup target exists
// and creates a new one if it does not
func checkObjectStoreLocationBackupTargetExistAndCreateIfNotExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_pc_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

				for i := 0; i < backupTargetsCount; i++ {
					objectStoreLocationCount, _ := strconv.Atoi(attributes["backup_targets."+strconv.Itoa(i)+".location.0.object_store_location.#"])

					if objectStoreLocationCount > 0 {
						log.Printf("[DEBUG] object store location backup target already exists, ext_id: %s", attributes["backup_targets."+strconv.Itoa(i)+".ext_id"])
						return nil
					}
				}

				log.Printf("[DEBUG] Backup Target not found, creating new Backup Target")
				break
			}
		}

		// Extract the output value for use in later steps
		outputDomainManagerExtID, ok := s.RootModule().Outputs["domainManagerExtID"]
		if !ok {
			return fmt.Errorf("output 'domainManagerExtID' not found")
		}

		domainManagerExtID := outputDomainManagerExtID.Value.(string)

		// Create Backup Target
		body := management.BackupTarget{}

		bucket := testVars.Prism.Bucket

		OneOfBackupTargetLocation := management.NewOneOfBackupTargetLocation()

		objectStoreLocationBody := management.NewObjectStoreLocation()

		// Set the provider config for AWS S3
		providerConfig := management.NewOneOfObjectStoreLocationProviderConfig()

		awsS3Config := management.NewAWSS3Config()
		awsS3Config.BucketName = utils.StringPtr(bucket.Name)
		awsS3Config.Region = utils.StringPtr(bucket.Region)
		awsS3Config.Credentials = &management.AccessKeyCredentials{
			AccessKeyId:     utils.StringPtr(bucket.AccessKey),
			SecretAccessKey: utils.StringPtr(bucket.SecretKey),
		}

		if err := providerConfig.SetValue(*awsS3Config); err != nil {
			return fmt.Errorf("error while setting provider config for AWS S3: %v", err)
		}

		objectStoreLocationBody.ProviderConfig = providerConfig

		objectStoreLocationBody.BackupPolicy = &management.BackupPolicy{
			RpoInMinutes: utils.IntPtr(60),
		}

		err := OneOfBackupTargetLocation.SetValue(*objectStoreLocationBody)
		if err != nil {
			return fmt.Errorf("error while setting object store location : %v", err)
		}

		body.Location = OneOfBackupTargetLocation

		log.Printf("[DEBUG] Creating Backup Target")
		resp, err := client.CreateBackupTarget(utils.StringPtr(domainManagerExtID), &body)

		if err != nil {
			return fmt.Errorf("error while Creating Backup Target: %s", err)
		}

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
		//nolint:staticcheck
		if _, taskErr := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error waiting for Backup Target to be deleted: %s", taskErr)
		}

		_, err = taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target Task Details: %s", err)
		}

		return nil
	}
}

// checkLastSyncTimeBackupTarget checks the last sync time of the backup target to know if the restore point is created
func checkLastSyncTimeBackupTarget(domainManagerExtID, backupTargetExtID *string, retries int, delay time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance
		if *backupTargetExtID == "" && *domainManagerExtID == "" {
			return fmt.Errorf("backup target ext_id and domain manager ext_id not set")
		}
		for i := 0; i < retries; i++ {
			readResp, err := client.GetBackupTargetById(domainManagerExtID, backupTargetExtID, nil)
			if err != nil {
				return fmt.Errorf("error while fetching Backup Target: %s", err)
			}

			backupTarget := readResp.Data.GetValue().(management.BackupTarget)

			log.Printf("[DEBUG] LastSyncTime: %v\n", backupTarget.LastSyncTime)
			if backupTarget.LastSyncTime != nil {
				log.Printf("[DEBUG]  Restore Point Created after %d minutes\n", i*30/60)
				return nil
			}
			log.Printf("[DEBUG] Waiting for 30 seconds to Fetch backup target\n")
			time.Sleep(delay)
		}

		return fmt.Errorf("backup Target restore point not created")
	}
}

// helper function to check the delete task
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
	//nolint:staticcheck
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

// helper function to check the delete task
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

// helper function to flatten the task status to string
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

// checkClusterLocationBackupTargetExistAndCreateIfNot checks if the cluster location backup target exists
// and creates a new one if it does not
// its set the backupTargetExtID to the ext_id of the created backup target
// this method is used to check the backup target for cluster location restore PC test
func checkClusterLocationBackupTargetExistAndCreateIfNot(backupTargetExtID, domainManagerExtID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		// Extract the output value for use in later steps
		outputDomainManagerExtID, ok := s.RootModule().Outputs["domainManagerExtID"]
		if !ok {
			return fmt.Errorf("output 'domainManagerExtID' not found")
		}

		*domainManagerExtID = outputDomainManagerExtID.Value.(string)

		outputClusterExtID, ok := s.RootModule().Outputs["clusterExtID"]
		if !ok {
			return fmt.Errorf("output 'clusterExtID' not found")
		}

		clusterExtID := outputClusterExtID.Value.(string)

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_pc_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

				for i := 0; i < backupTargetsCount; i++ {
					clusterLocationCount, _ := strconv.Atoi(attributes["backup_targets."+strconv.Itoa(i)+".location.0.cluster_location.#"])

					if clusterLocationCount > 0 {
						if attributes["backup_targets."+strconv.Itoa(i)+".location.0.cluster_location.0.config.0.ext_id"] == clusterExtID {
							log.Printf("[DEBUG] cluster location backup target already exists, ext_id: %s", attributes["backup_targets."+strconv.Itoa(i)+".ext_id"])
							*backupTargetExtID = attributes["backup_targets."+strconv.Itoa(i)+".ext_id"]
							return nil
						}
					}
				}
				log.Printf("[DEBUG] cluster location backup target target not found, creating new cluster location backup target")
				break
			}
		}

		// Create Backup Target
		body := management.BackupTarget{}

		OneOfBackupTargetLocation := management.NewOneOfBackupTargetLocation()

		clusterConfigBody := management.NewClusterLocation()
		clusterRef := management.NewClusterReference()

		clusterRef.ExtId = utils.StringPtr(clusterExtID)
		oneOfClusterLocationConfig := management.NewOneOfClusterLocationConfig()
		oneOfClusterLocationConfig.SetValue(*clusterRef)
		clusterConfigBody.Config = oneOfClusterLocationConfig

		err := OneOfBackupTargetLocation.SetValue(*clusterConfigBody)
		if err != nil {
			return fmt.Errorf("error while setting cluster location : %v", err)
		}

		body.Location = OneOfBackupTargetLocation

		resp, err := client.CreateBackupTarget(domainManagerExtID, &body)

		if err != nil {
			return fmt.Errorf("error while Creating Backup Target: %s", err)
		}

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
		//nolint:staticcheck
		if _, taskErr := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error waiting for Backup Target to be deleted: %s", taskErr)
		}

		_, err = taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target Task Details: %s", err)
		}

		listResp, err := client.ListBackupTargets(domainManagerExtID, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target: %s", err)
		}
		backupTargets := listResp.Data.GetValue().([]management.BackupTarget)

		// Find the new backup target ext id
		for _, backupTarget := range backupTargets {
			backupTargetLocation := backupTarget.Location
			if utils.StringValue(backupTargetLocation.ObjectType_) == "prism.v4.management.ClusterLocation" {
				clusterLocation := backupTarget.Location.GetValue().(management.ClusterLocation)
				clusterConfig := clusterLocation.Config.GetValue().(management.ClusterReference)
				if utils.StringValue(clusterConfig.ExtId) == clusterExtID {
					*backupTargetExtID = utils.StringValue(backupTarget.ExtId)
					break
				}
			}
		}

		return nil
	}
}

// checkObjectRestoreLocationBackupTargetExistAndCreateIfNot checks if the object restore location backup target exists
// and creates a new one if it does not
// its set the backupTargetExtID to the ext_id of the created backup target
// this method is used to check the backup target for object restore location restore PC test
func checkObjectRestoreLocationBackupTargetExistAndCreateIfNot(backupTargetExtID, domainManagerExtID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		// Extract the output value for use in later steps
		outputDomainManagerExtID, ok := s.RootModule().Outputs["domainManagerExtID"]
		if !ok {
			return fmt.Errorf("output 'domainManagerExtID' not found")
		}

		*domainManagerExtID = outputDomainManagerExtID.Value.(string)

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_pc_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

				for i := 0; i < backupTargetsCount; i++ {
					objectStoreLocationCount, _ := strconv.Atoi(attributes["backup_targets."+strconv.Itoa(i)+".location.0.object_store_location.#"])

					if objectStoreLocationCount > 0 {
						log.Printf("[DEBUG] Object store location backup target already exists, ext_id: %s", attributes["backup_targets."+strconv.Itoa(i)+".ext_id"])
						*backupTargetExtID = attributes["backup_targets."+strconv.Itoa(i)+".ext_id"]
						return nil
					}
				}
				log.Printf("[DEBUG] Object store location backup target not found, creating new Object store location backup target")
				break
			}
		}

		// Create Backup Target Aws S3
		body := management.BackupTarget{}

		bucket := testVars.Prism.Bucket

		OneOfBackupTargetLocation := management.NewOneOfBackupTargetLocation()

		objectStoreLocationBody := management.NewObjectStoreLocation()

		// Set the provider config for AWS S3
		providerConfig := management.NewOneOfObjectStoreLocationProviderConfig()

		awsS3Config := management.NewAWSS3Config()
		awsS3Config.BucketName = utils.StringPtr(bucket.Name)
		awsS3Config.Region = utils.StringPtr(bucket.Region)
		awsS3Config.Credentials = &management.AccessKeyCredentials{
			AccessKeyId:     utils.StringPtr(bucket.AccessKey),
			SecretAccessKey: utils.StringPtr(bucket.SecretKey),
		}

		if err := providerConfig.SetValue(*awsS3Config); err != nil {
			return fmt.Errorf("error while setting provider config for AWS S3: %v", err)
		}

		objectStoreLocationBody.ProviderConfig = providerConfig

		objectStoreLocationBody.BackupPolicy = &management.BackupPolicy{
			RpoInMinutes: utils.IntPtr(60),
		}

		err := OneOfBackupTargetLocation.SetValue(*objectStoreLocationBody)
		if err != nil {
			return fmt.Errorf("error while setting object store location : %v", err)
		}

		body.Location = OneOfBackupTargetLocation

		log.Printf("[DEBUG] Creating Backup Target")
		resp, err := client.CreateBackupTarget(domainManagerExtID, &body)

		if err != nil {
			return fmt.Errorf("error while Creating Backup Target: %s", err)
		}

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
		//nolint:staticcheck
		if _, taskErr := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error waiting for Backup Target to be deleted: %s", taskErr)
		}

		_, err = taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target Task Details: %s", err)
		}

		listResp, err := client.ListBackupTargets(domainManagerExtID, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target: %s", err)
		}
		backupTargets := listResp.Data.GetValue().([]management.BackupTarget)

		// Find the new backup target ext id
		for _, backupTarget := range backupTargets {
			backupTargetLocation := backupTarget.Location
			if utils.StringValue(backupTargetLocation.ObjectType_) == "prism.v4.management.ObjectStoreLocation" {
				objectStoreLocation := backupTarget.Location.GetValue().(management.ObjectStoreLocation)
				if *objectStoreLocation.ProviderConfig.ObjectType_ == awsS3ConfigObjectType {
					awsS3Config := objectStoreLocation.ProviderConfig.GetValue().(management.AWSS3Config)
					if utils.StringValue(awsS3Config.BucketName) == bucket.Name {
						*backupTargetExtID = utils.StringValue(backupTarget.ExtId)
						log.Printf("[DEBUG] AWS S3Object store location backup target Ext ID: %s", *backupTargetExtID)
						break
					}
				}
			}
		}

		return nil
	}
}

// checkLastSyncTimeBackupTarget checks the last sync time of the backup target to know if the restore point is created
// this method is used to check the last sync time of the backup target for restore PC test
func checkLastSyncTimeBackupTargetRestorePC(backupTargetExtID *string, retries int, delay time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] Checking Last Sync Time\n")

		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		// Extract the output value for use in later steps
		outputDomainManagerExtID, ok := s.RootModule().Outputs["domainManagerExtID"]
		if !ok {
			return fmt.Errorf("output 'domainManagerExtID' not found")
		}

		pcExtID := outputDomainManagerExtID.Value.(string)

		for i := 0; i < retries; i++ {
			readResp, err := client.GetBackupTargetById(utils.StringPtr(pcExtID), backupTargetExtID, nil)
			if err != nil {
				return fmt.Errorf("error while fetching Backup Target: %s", err)
			}

			backupTarget := readResp.Data.GetValue().(management.BackupTarget)

			log.Printf("[DEBUG] LastSyncTime: %v\n", backupTarget.LastSyncTime)
			if backupTarget.LastSyncTime != nil {
				log.Printf("[DEBUG]  Restore Point Created after %d minutes\n", i*30/60)
				return nil
			}
			log.Printf("[DEBUG] Waiting for 30 seconds to Fetch backup target\n")
			time.Sleep(delay)
		}

		return fmt.Errorf("backup Target restore point not created")
	}
}

// createRestoreSource to create the object store location restore source location restore source
// this method is used to create the restore source for restore PC test
func createClusterLocationRestoreSource(restoreSourceExtID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] Create Restore Source\n")
		conn := acc.TestAccProvider2.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		// Extract the output value for use in later steps
		outputClusterExtID, ok := s.RootModule().Outputs["clusterExtID"]
		if !ok {
			return fmt.Errorf("output 'clusterExtID' not found")
		}

		clusterExtID := outputClusterExtID.Value.(string)

		// Create Backup Target
		body := management.RestoreSource{}

		oneOfRestoreSourceLocation := management.NewOneOfRestoreSourceLocation()

		clusterConfigBody := management.NewClusterLocation()
		clusterRef := management.NewClusterReference()

		clusterRef.ExtId = utils.StringPtr(clusterExtID)

		oneOfClusterLocationConfig := management.NewOneOfClusterLocationConfig()
		oneOfClusterLocationConfig.SetValue(*clusterRef)
		clusterConfigBody.Config = oneOfClusterLocationConfig

		err := oneOfRestoreSourceLocation.SetValue(*clusterConfigBody)
		if err != nil {
			return fmt.Errorf("error while setting cluster location : %v", err)
		}

		body.Location = oneOfRestoreSourceLocation

		resp, err := client.CreateRestoreSource(&body)

		if err != nil {
			return fmt.Errorf("error while Creating Restore Source: %s", err)
		}

		restoreSource := resp.Data.GetValue().(management.RestoreSource)
		*restoreSourceExtID = utils.StringValue(restoreSource.ExtId)

		return nil
	}
}

// createRestoreSource to create the cluster object store location restore source
// this method is used to create the restore source for restore PC test
func createObjectStoreLocationLocationRestoreSource(restoreSourceExtID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] Create object store location restore source\n")
		conn := acc.TestAccProvider2.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		// Create Backup Target
		body := management.RestoreSource{}
		bucket := testVars.Prism.Bucket
		oneOfRestoreSourceLocation := management.NewOneOfRestoreSourceLocation()

		objectStoreLocationBody := management.NewObjectStoreLocation()

		// Set the provider config for AWS S3
		providerConfig := management.NewOneOfObjectStoreLocationProviderConfig()

		awsS3Config := management.NewAWSS3Config()
		awsS3Config.BucketName = utils.StringPtr(bucket.Name)
		awsS3Config.Region = utils.StringPtr(bucket.Region)
		awsS3Config.Credentials = &management.AccessKeyCredentials{
			AccessKeyId:     utils.StringPtr(bucket.AccessKey),
			SecretAccessKey: utils.StringPtr(bucket.SecretKey),
		}

		if err := providerConfig.SetValue(*awsS3Config); err != nil {
			return fmt.Errorf("error while setting provider config for AWS S3: %v", err)
		}

		objectStoreLocationBody.ProviderConfig = providerConfig

		err := oneOfRestoreSourceLocation.SetValue(*objectStoreLocationBody)
		if err != nil {
			return fmt.Errorf("error while setting cluster location : %v", err)
		}

		body.Location = oneOfRestoreSourceLocation

		resp, err := client.CreateRestoreSource(&body)

		if err != nil {
			return fmt.Errorf("error while creating object store location restore source: %s", err)
		}

		restoreSource := resp.Data.GetValue().(management.RestoreSource)
		*restoreSourceExtID = utils.StringValue(restoreSource.ExtId)
		aJSON, _ := json.MarshalIndent(restoreSource, "", "  ")
		log.Printf("[DEBUG] Restore Source Create Response: %s", string(aJSON))

		return nil
	}
}

// powerOffPC to power off the PC
// this method is used to power off the PC for restore PC test
func powerOffPC() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		vmClient := conn.VmmAPI.VMAPIInstance

		// Cluster filter
		vmsResp, err := vmClient.ListVms(nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error while fetching VMs: %s", err)
		}

		vms := vmsResp.Data.GetValue().([]vmConfig.Vm)

		for _, vm := range vms {
			if utils.StringValue(vm.Description) == "NutanixPrismCentral" {
				// if the PC VM does not have any NICs, continue to the next VM
				if len(vm.Nics) == 0 {
					continue
				}
				for _, nic := range vm.Nics {
					// if the PC VM does not have any network info, continue to the next VM
					if nic.NetworkInfo == nil || nic.NetworkInfo.Ipv4Info == nil {
						continue
					}

					// loop through the learned IP addresses to find the correct PC VM
					for _, learnedIPAddress := range nic.NetworkInfo.Ipv4Info.LearnedIpAddresses {
						if utils.StringValue(learnedIPAddress.Value) == os.Getenv("NUTANIX_ENDPOINT") {
							// get etag
							readResp, err := vmClient.GetVmById(vm.ExtId, nil)
							if err != nil {
								return fmt.Errorf("error while fetching PC: %s", err)
							}
							args := make(map[string]interface{})
							eTag := vmClient.ApiClient.GetEtag(readResp)
							args["If-Match"] = utils.StringPtr(eTag)

							// Power off the PC
							_, err = vmClient.PowerOffVm(vm.ExtId, args)
							if err != nil {
								log.Printf("[DEBUG] error while powering off PC: %s", err)
								// after the PC is powered off, the API returns timeout error
								return nil
							}

							return nil
						}
					}
				}
			}
		}
		return fmt.Errorf("PC VM not found")
	}
}

// method to expand the PC config block
// to be used in the restore PC resource configuration
func expandPCConfigBlock(configMap map[string]interface{}) string {
	// Extract nested values from the map
	buildInfo, ok := configMap["build_info"].([]interface{})
	if !ok || len(buildInfo) == 0 {
		panic("build_info is not a slice or is empty")
	}
	buildInfoMap, ok := buildInfo[0].(map[string]interface{})
	if !ok {
		panic("build_info[0] is not a map")
	}
	version := buildInfoMap["version"].(string)
	showEnableLockdownMode := configMap["should_enable_lockdown_mode"].(bool)
	name := configMap["name"].(string)
	size := configMap["size"].(string)

	resourceConfig, ok := configMap["resource_config"].([]interface{})
	if !ok || len(resourceConfig) == 0 {
		panic("resource_config is not a slice or is empty")
	}
	resourceConfigMap, ok := resourceConfig[0].(map[string]interface{})
	if !ok {
		panic("resource_config[0] is not a map")
	}
	containerExtIDs := resourceConfigMap["container_ext_ids"].([]interface{})
	// Convert all elements to strings (add quotes implicitly in Go)
	strContainerExtIDs := make([]string, len(containerExtIDs))
	for i, extID := range containerExtIDs {
		strContainerExtIDs[i] = fmt.Sprintf("\"%v\"", extID) // Convert to string
	}

	dataDiskSizeBytesStr := resourceConfigMap["data_disk_size_bytes"].(json.Number).String()
	dataDiskSizeBytes, err := strconv.Atoi(dataDiskSizeBytesStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert data_disk_size_bytes to int: %v", err))
	}
	memorySizeBytesStr := resourceConfigMap["memory_size_bytes"].(json.Number).String()
	memorySizeBytes, err := strconv.Atoi(memorySizeBytesStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert memory_size_bytes to int: %v", err))
	}

	numVcpusStr := resourceConfigMap["num_vcpus"].(json.Number).String()
	numVcpus, err := strconv.Atoi(numVcpusStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert num_vcpus to int: %v", err))
	}
	return fmt.Sprintf(`
		config {
			should_enable_lockdown_mode = %t
			build_info {
				version = "%s"
			}
			name = "%s"
			size = "%s"
			resource_config {
				container_ext_ids    = %v
				data_disk_size_bytes = %d
				memory_size_bytes    = %d
				num_vcpus            = %d
			}
		}`, showEnableLockdownMode, version, name, size, strContainerExtIDs, dataDiskSizeBytes, memorySizeBytes, numVcpus)
}

// method to expand the PC network block
// to be used in the restore PC resource configuration
func expandPCNetworkBlock(networkMap map[string]interface{}) string {
	// Extract nested values from the map
	externalAddress, ok := networkMap["external_address"].([]interface{})
	if !ok || len(externalAddress) == 0 {
		panic("external_address is not a slice or is empty")
	}
	externalAddressMap, ok := externalAddress[0].(map[string]interface{})
	if !ok {
		panic("external_address[0] is not a map")
	}
	externalAddressIPv4, ok := externalAddressMap["ipv4"].([]interface{})
	if !ok || len(externalAddressIPv4) == 0 {
		panic("external_address.ipv4 is not a slice or is empty")
	}
	externalAddressIPv4Map, ok := externalAddressIPv4[0].(map[string]interface{})
	if !ok {
		panic("external_address.ipv4[0] is not a map")
	}
	externalAddressIPv4Value := externalAddressIPv4Map["value"].(string)

	nameServers := networkMap["name_servers"].([]interface{})
	nameServersConfig := ""
	for _, nameServer := range nameServers {
		nameServerMap, okNameServerMap := nameServer.(map[string]interface{})
		if !okNameServerMap {
			panic("name_server is not a map")
		}
		ipv4, okIpv4 := nameServerMap["ipv4"].([]interface{})
		if !okIpv4 || len(ipv4) == 0 {
			panic("ipv4 is not a slice or is empty")
		}
		ipv4Map, okIpv4Map := ipv4[0].(map[string]interface{})
		if !okIpv4Map {
			panic("ipv4[0] is not a map")
		}
		nameServerIPv4Value := ipv4Map["value"].(string)
		nameServersConfig += fmt.Sprintf(`
			name_servers {
				ipv4 {
					value = "%s"
				}
			}`, nameServerIPv4Value)
	}

	ntpServers := networkMap["ntp_servers"].([]interface{})

	ntpServersConfig := ""
	for _, ntpServer := range ntpServers {
		ntpServerMap, okNtpServerMap := ntpServer.(map[string]interface{})
		if !okNtpServerMap {
			panic("ntp_server is not a map")
		}
		fqdn, okFqdn := ntpServerMap["fqdn"].([]interface{})
		if !okFqdn || len(fqdn) == 0 {
			panic("fqdn is not a slice or is empty")
		}
		fqdnMap, okFqdnMap := fqdn[0].(map[string]interface{})
		if !okFqdnMap {
			panic("fqdn[0] is not a map")
		}
		ntpServerFQDN := fqdnMap["value"].(string)
		ntpServersConfig += fmt.Sprintf(`
			ntp_servers {
				fqdn {
					value = "%s"
				}
			}`, ntpServerFQDN)
	}

	externalNetworks, okExternalNetworks := networkMap["external_networks"].([]interface{})
	if !okExternalNetworks || len(externalNetworks) == 0 {
		panic("external_networks is not a slice or is empty")
	}
	externalNetworksMap, ok := externalNetworks[0].(map[string]interface{})
	if !ok {
		panic("external_networks[0] is not a map")
	}
	networkExtID := externalNetworksMap["network_ext_id"].(string)
	defaultGatewayIPv4 := externalNetworksMap["default_gateway"].([]interface{})[0].(map[string]interface{})["ipv4"].([]interface{})[0].(map[string]interface{})["value"].(string)
	subnetMaskIPv4 := externalNetworksMap["subnet_mask"].([]interface{})[0].(map[string]interface{})["ipv4"].([]interface{})[0].(map[string]interface{})["value"].(string)
	ipRanges := externalNetworksMap["ip_ranges"].([]interface{})[0].(map[string]interface{})
	ipRangeBeginIPv4 := ipRanges["begin"].([]interface{})[0].(map[string]interface{})["ipv4"].([]interface{})[0].(map[string]interface{})["value"].(string)
	ipRangeEndIPv4 := ipRanges["end"].([]interface{})[0].(map[string]interface{})["ipv4"].([]interface{})[0].(map[string]interface{})["value"].(string)

	return fmt.Sprintf(`
		network {
			external_address {
				ipv4 {
					value = "%s"
				}
			}
			# name servers
			%s
			# ntp servers
			%s

			external_networks {
				network_ext_id = "%s"
				default_gateway {
					ipv4 {
						value = "%s"
					}
				}
				subnet_mask {
					ipv4 {
						value = "%s"
					}
				}
				ip_ranges {
					begin {
						ipv4 {
							value = "%s"
						}
					}
					end {
						ipv4 {
							value = "%s"
						}
					}
				}
			}
		}
`, externalAddressIPv4Value, nameServersConfig, ntpServersConfig, networkExtID, defaultGatewayIPv4, subnetMaskIPv4,
		ipRangeBeginIPv4, ipRangeEndIPv4)
}

// generate Random Passwords
var (
	lowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")
	upperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits       = []rune("0123456789")
	specials     = []rune("@#$")
)

// allChars is the union of all allowed characters.
var allChars = append(append(append(lowerLetters, upperLetters...), digits...), specials...)

// getRandomRune returns a random rune from a given set.
func getRandomRune(set []rune) rune {
	return set[rand.Intn(len(set))]
}

// hasConsecutiveDuplicates returns true if there are three identical runes in a row.
func hasConsecutiveDuplicates(p []rune) bool {
	for i := 2; i < len(p); i++ {
		if p[i] == p[i-1] && p[i] == p[i-2] {
			return true
		}
	}
	return false
}

// meetsRequirements checks that p contains at least one lowercase letter,
// one uppercase letter, one digit, and one special character.
func meetsRequirements(p []rune) bool {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, c := range p {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case strings.ContainsRune(string(specials), c):
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasDigit && hasSpecial
}

// generatePassword builds a password that meets the requirements.
func generatePassword() (string, error) {
	// Choose a random length between 9 and 15 characters.
	length := rand.Intn(7) + 9 // 9-15 characters

	// Try up to 100 times to generate a valid password.
	for attempt := 0; attempt < 100; attempt++ {
		password := make([]rune, 0, length)

		// Guarantee one character from each required set.
		password = append(password, getRandomRune(lowerLetters))
		password = append(password, '.')
		password = append(password, getRandomRune(upperLetters))
		password = append(password, '.')
		password = append(password, getRandomRune(digits))
		password = append(password, '.')
		password = append(password, getRandomRune(specials))

		// Fill remaining characters.
		for len(password) < length {
			password = append(password, '.')
			password = append(password, getRandomRune(allChars))
		}

		// Validate constraints.
		if hasConsecutiveDuplicates(password) {
			continue
		}
		if !meetsRequirements(password) {
			continue
		}

		// Password meets all requirements.
		return string(password), nil
	}

	return "", fmt.Errorf("failed to generate valid password after 100 attempts")
}
