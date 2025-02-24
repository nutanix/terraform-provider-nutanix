package prismv2_test

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

					log.Printf("[DEBUG] LastSyncTime: %v\n", backupTarget.LastSyncTime)
					if backupTarget.LastSyncTime != nil {
						log.Printf("[DEBUG]  Restore Point Created after %d minutes\n", i*30/60)
						return nil
					}
					log.Printf("[DEBUG] Waiting for 30 seconds to Fetch backup target\n")
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

func checkBackupTargetExistAndCreateIfNot(backupTargetExtID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		for _, rs := range s.RootModule().Resources {
			if rs.Type == "nutanix_backup_targets_v2" {
				attributes := rs.Primary.Attributes

				backupTargetsCount, _ := strconv.Atoi(attributes["backup_targets.#"])

				if backupTargetsCount > 0 {
					log.Printf("[DEBUG] Backup Target already exists, ext_id: %s", attributes["backup_targets.0.ext_id"])
					*backupTargetExtID = attributes["backup_targets.0.ext_id"]
					return nil
				} else {
					log.Printf("[DEBUG] Backup Target not found, creating new Backup Target")
					break
				}
			}
		}

		// Extract the output value for use in later steps
		outputDomainManagerExtID, ok := s.RootModule().Outputs["domainManagerExtID"]
		if !ok {
			return fmt.Errorf("output 'domainManagerExtID' not found")
		}

		domainManagerExtID := outputDomainManagerExtID.Value.(string)

		outputClusterExtID, ok := s.RootModule().Outputs["clusterExtID"]
		if !ok {
			return fmt.Errorf("output 'clusterExtID' not found")
		}

		clusterExtID := outputClusterExtID.Value.(string)

		// Create Backup Target
		body := management.BackupTarget{}

		OneOfBackupTargetLocation := management.NewOneOfBackupTargetLocation()

		clusterConfigBody := management.NewClusterLocation()
		clusterRef := management.NewClusterReference()

		clusterRef.ExtId = utils.StringPtr(clusterExtID)

		clusterConfigBody.Config = clusterRef

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
				if utils.StringValue(clusterLocation.Config.ExtId) == clusterExtID {
					*backupTargetExtID = utils.StringValue(backupTarget.ExtId)
					break
				}
			}
		}

		return nil
	}
}

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

func createRestoreSource(restoreSourceExtID *string) resource.TestCheckFunc {
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

		clusterConfigBody.Config = clusterRef

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

func deleteBackupTarget(backupTargetExtID, pcExtID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.PrismAPI.DomainManagerBackupsAPIInstance

		readResp, err := client.GetBackupTargetById(pcExtID, backupTargetExtID)
		if err != nil {
			return fmt.Errorf("error while fetching Backup Target: %s", err)
		}

		// extract the etag from the read response
		args := make(map[string]interface{})
		eTag := client.ApiClient.GetEtag(readResp)
		args["If-Match"] = utils.StringPtr(eTag)

		_, err = client.DeleteBackupTargetById(pcExtID, backupTargetExtID, args)

		if err != nil {
			return fmt.Errorf("error while deleting Backup Target: %s", err)
		}

		return nil
	}
}

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
			if vm.MachineType.GetName() == "PC" && utils.StringValue(vm.Description) == "NutanixPrismCentral" &&
				strings.Contains(utils.StringValue(vm.Name), "auto_pc_") {
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
					//return fmt.Errorf("error while powering off PC: %s", err)
					return nil
				}

				return nil
			}
		}
		return fmt.Errorf("PC not found")
	}
}

func expandDomainManagerConfigBlock(pcDetails map[string]interface{}) string {
	// Extract nested values from the map
	config, ok := pcDetails["config"].([]interface{})
	if !ok || len(config) == 0 {
		panic("config is not a slice or is empty")
	}
	configMap, ok := config[0].(map[string]interface{})
	if !ok {
		panic("config[0] is not a map")
	}

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

	network, ok := pcDetails["network"].([]interface{})
	if !ok || len(network) == 0 {
		panic("network is not a slice or is empty")
	}
	networkMap, ok := network[0].(map[string]interface{})
	if !ok {
		panic("network[0] is not a map")
	}

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
		  }

`, nameServerIPv4Value)
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
		  }

`, ntpServerFQDN)
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
	domain_manager {
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
		}
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

	}`, showEnableLockdownMode, version, name, size, strContainerExtIDs, dataDiskSizeBytes, memorySizeBytes, numVcpus,
		externalAddressIPv4Value, nameServersConfig, ntpServersConfig, networkExtID, defaultGatewayIPv4, subnetMaskIPv4,
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
