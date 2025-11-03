package prismv2_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameRestorePC = "nutanix_pc_restore_v2.test"

func TestAccV2NutanixRestorePCResource_ClusterLocationRestorePC(t *testing.T) {
	if testVars.Prism.PCRestore.SkipPCRestoreTest {
		// We are skipping the PC restore tests because they require powering off the PC VM,
		// which could affect the execution of other test cases running in parallel.
		// The PC restore test cases can be run separately.
		t.Skip("Skipping PC restore test: We are skipping the PC restore tests because they require powering off the PC VM, which could affect the execution of other test cases running in parallel. The PC restore test cases can be run separately.")
	}
	var backupTargetExtID, domainManagerExtID, restoreSourceExtID = new(string), new(string), new(string)
	var restorePcConfig string

	pcDetails := make(map[string]interface{})

	t.Run("pre request: backup target and restore source ", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreventPostDestroyRefresh: true,
			PreCheck:                  func() { acc.TestAccPreCheck(t) },
			Providers:                 acc.TestAccProviders,
			Steps: []resource.TestStep{
				// Step 1: List backup targets and delete if backup target exists
				{
					PreConfig: func() {
						fmt.Printf("Step 1: List backup targets and create if backup target does not exist\n")
					},
					Config: testAccPreRequestForRestoreSourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkClusterLocationBackupTargetExistAndCreateIfNot(backupTargetExtID, domainManagerExtID),
					),
				},
				// Step 2: Check last sync time for backup target
				{
					PreConfig: func() {
						fmt.Printf("Step 2: Create Restore Source\n")
					},
					Config: testAccPreRequestForRestoreSourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkLastSyncTimeBackupTargetRestorePC(backupTargetExtID, retries, delay),
						createClusterLocationRestoreSource(restoreSourceExtID),
					),
				},
			},
		})
	})

	// fetch the restore point and extract the pc details
	t.Run("power of PC", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreventPostDestroyRefresh: true,
			PreCheck:                  func() { acc.TestAccPreCheck(t) },
			Providers:                 acc.TestAccProviders,
			Steps: []resource.TestStep{
				// Step 3: power off PC
				{
					PreConfig: func() {
						fmt.Printf("Step 3: build PC restore Block and Power off PC\n")
					},
					Config: testAccPowerOffPCConfig(utils.StringValue(restoreSourceExtID), utils.StringValue(domainManagerExtID)),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							// Build the restore PC configuration for the next sub-test case
							pcDetailsOutput, ok := s.RootModule().Outputs["pc_details"]
							if !ok {
								return fmt.Errorf("output 'pc_details' not found")
							}
							pcDetails = pcDetailsOutput.Value.(map[string]interface{})

							restoreSourceExtIDOutput, ok := s.RootModule().Outputs["restoreSourceExtID"]
							if !ok {
								return fmt.Errorf("output 'restoreSourceExtID' not found")
							}
							restoreSourceExtID := restoreSourceExtIDOutput.Value.(string)

							restorePointExtIDOutput, ok := s.RootModule().Outputs["restorePointExtID"]
							if !ok {
								return fmt.Errorf("output 'restorePointExtID' not found")
							}
							restorePointExtID := restorePointExtIDOutput.Value.(string)

							restorablePcExtIDOutput, ok := s.RootModule().Outputs["restorablePcExtID"]
							if !ok {
								return fmt.Errorf("output 'restorablePcExtID' not found")
							}
							restorablePcExtID := restorablePcExtIDOutput.Value.(string)

							restorePcConfig = restorePcResourceConfig(pcDetails, restoreSourceExtID, restorePointExtID, restorablePcExtID)
							log.Printf("[DEBUG] Restore PC Config: %s\n", restorePcConfig)
							return nil
						},
						powerOffPC(),
					),
				},
			},
		})
	})

	// Restore PC Sub-test Case
	t.Run("PC restore test: ", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { acc.TestAccPreCheck(t) },
			Providers: acc.TestAccProviders,
			Steps: []resource.TestStep{
				// Step 5: Restore PC
				{
					PreConfig: func() {
						fmt.Printf("Step 4: Restore PC\n")
					},
					Config: restorePcConfig,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "id"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.config.0.build_info.0.version"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.config.0.name"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.config.0.size"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.external_address.0.ipv4.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.name_servers.0.ipv4.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.name_servers.1.ipv4.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.ntp_servers.0.fqdn.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.ntp_servers.1.fqdn.0.value"),
					),
				},
			},
		})
	})
}

func TestAccV2NutanixRestorePCResource_ObjectRestoreSourceRestorePC(t *testing.T) {
	if testVars.Prism.PCRestore.SkipPCRestoreTest {
		// We are skipping the PC restore tests because they require powering off the PC VM,
		// which could affect the execution of other test cases running in parallel.
		// The PC restore test cases can be run separately.
		t.Skip("Skipping PC restore test: We are skipping the PC restore tests because they require powering off the PC VM, which could affect the execution of other test cases running in parallel. The PC restore test cases can be run separately.")
	}
	var backupTargetExtID, domainManagerExtID, restoreSourceExtID = new(string), new(string), new(string)
	var restorePcConfig string

	pcDetails := make(map[string]interface{})

	t.Run("pre request: backup target and restore source ", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreventPostDestroyRefresh: true,
			PreCheck:                  func() { acc.TestAccPreCheck(t) },
			Providers:                 acc.TestAccProviders,
			Steps: []resource.TestStep{
				// Step 1: List backup targets and create if backup target does not exist
				{
					PreConfig: func() {
						fmt.Printf("Step 1: List backup targets and create if backup target does not exist\n")
					},
					Config: testAccPreRequestForRestoreSourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkObjectRestoreLocationBackupTargetExistAndCreateIfNot(backupTargetExtID, domainManagerExtID),
					),
				},
				// Step 2: Check last sync time for backup target
				{
					PreConfig: func() {
						fmt.Printf("Step 2: Create Restore Source\n")
					},
					Config: testAccPreRequestForRestoreSourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkLastSyncTimeBackupTargetRestorePC(backupTargetExtID, retries, delay),
						createObjectStoreLocationLocationRestoreSource(restoreSourceExtID),
					),
				},
			},
		})
	})

	// fetch the restore point and extract the pc details
	t.Run("power off PC: ", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreventPostDestroyRefresh: true,
			PreCheck:                  func() { acc.TestAccPreCheck(t) },
			Providers:                 acc.TestAccProviders,
			Steps: []resource.TestStep{
				// Step 3: power off PC
				{
					PreConfig: func() {
						fmt.Printf("Step 3: build PC restore config block and power off PC\n")
					},
					Config: testAccPowerOffPCConfig(utils.StringValue(restoreSourceExtID), utils.StringValue(domainManagerExtID)),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							// Build the restore PC configuration for the next sub-test case
							pcDetailsOutput, ok := s.RootModule().Outputs["pc_details"]
							if !ok {
								return fmt.Errorf("output 'pc_details' not found")
							}
							pcDetails = pcDetailsOutput.Value.(map[string]interface{})

							restoreSourceExtIDOutput, ok := s.RootModule().Outputs["restoreSourceExtID"]
							if !ok {
								return fmt.Errorf("output 'restoreSourceExtID' not found")
							}
							restoreSourceExtID := restoreSourceExtIDOutput.Value.(string)

							restorePointExtIDOutput, ok := s.RootModule().Outputs["restorePointExtID"]
							if !ok {
								return fmt.Errorf("output 'restorePointExtID' not found")
							}
							restorePointExtID := restorePointExtIDOutput.Value.(string)

							restorablePcExtIDOutput, ok := s.RootModule().Outputs["restorablePcExtID"]
							if !ok {
								return fmt.Errorf("output 'restorablePcExtID' not found")
							}
							restorablePcExtID := restorablePcExtIDOutput.Value.(string)

							restorePcConfig = restorePcResourceConfig(pcDetails, restoreSourceExtID, restorePointExtID, restorablePcExtID)
							log.Printf("[DEBUG] Restore PC Config: %s\n", restorePcConfig)
							return nil
						},
						powerOffPC(),
					),
				},
			},
		})
	})

	// Restore PC Sub-test Case
	t.Run("PC restore test: ", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { acc.TestAccPreCheck(t) },
			Providers: acc.TestAccProviders,
			Steps: []resource.TestStep{
				// Step 5: Restore PC
				{
					PreConfig: func() {
						fmt.Printf("Step 4: Restore PC\n")
					},
					Config: restorePcConfig,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "id"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.config.0.build_info.0.version"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.config.0.name"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.config.0.size"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.external_address.0.ipv4.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.name_servers.0.ipv4.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.name_servers.1.ipv4.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.ntp_servers.0.fqdn.0.value"),
						resource.TestCheckResourceAttrSet(resourceNameRestorePC, "domain_manager.0.network.0.ntp_servers.1.fqdn.0.value"),
					),
				},
			},
		})
	})
}

func testAccPreRequestForRestoreSourceConfig() string {
	// pe config
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	peHostProviderConfig := fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}
`, username, password, testVars.Prism.RestoreSource.PeIP, insecure, port)

	return fmt.Sprintf(`
# peHostProviderConfig
%s
data "nutanix_clusters_v2" "cls" {
	provider = nutanix
	filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}
data "nutanix_clusters_v2" "clusters" {
    provider = nutanix
}
locals {
  domainManagerExtId = data.nutanix_clusters_v2.cls.cluster_entities.0.ext_id
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

data "nutanix_pc_backup_targets_v2" "test" {
  domain_manager_ext_id = local.domainManagerExtId
}

output "domainManagerExtID" {
	  value = local.domainManagerExtId
}

output "clusterExtID" {
    value = local.clusterExtId
}

# Dummy data source to make sure the the second provider is initialized
data "nutanix_subnets_v2" "subnets" {
    provider = nutanix-2
}
`, peHostProviderConfig)
}

func testAccPowerOffPCConfig(restoreSourceExtID, domainManagerExtID string) string {
	// pe config
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	peHostProviderConfig := fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}
`, username, password, testVars.Prism.RestoreSource.PeIP, insecure, port)

	return fmt.Sprintf(`
%[1]s


data "nutanix_restorable_pcs_v2" "restorable-pcs" {
  provider              = nutanix-2
  restore_source_ext_id = "%[2]s"
  filter = "extId eq %[3]s"
}

locals {
  restorablePcExtId = data.nutanix_restorable_pcs_v2.restorable-pcs.restorable_pcs.0.ext_id
}

data "nutanix_pc_restore_points_v2" "restore-points" {
  provider                         = nutanix-2
  restorable_domain_manager_ext_id = local.restorablePcExtId
  order_by						   = "creationTime desc"
  restore_source_ext_id            = "%[2]s"
}

data "nutanix_pc_restore_point_v2" "restore-point" {
  provider = nutanix-2
  restore_source_ext_id = "%[2]s"
  restorable_domain_manager_ext_id = local.restorablePcExtId
  ext_id   = data.nutanix_pc_restore_points_v2.restore-points.restore_points[0].ext_id
}

locals {
  restorePoint = data.nutanix_pc_restore_point_v2.restore-point
}

output "pc_details" {
  value = local.restorePoint.domain_manager[0]
}

output "restoreSourceExtID" {
  value = "%[2]s"
}

output "restorePointExtID" {
  value = local.restorePoint.ext_id
}

output "restorablePcExtID" {
  value = local.restorablePcExtId
}

`, peHostProviderConfig, restoreSourceExtID, domainManagerExtID)
}

func restorePcResourceConfig(pcDetails map[string]interface{}, restoreSourceExtID, restorePointExtID, restorablePcExtID string) string {
	// Extract Pc details from the output
	// Extract Config values from the map
	configBlock, ok := pcDetails["config"].([]interface{})
	if !ok || len(configBlock) == 0 {
		panic("config is not a slice or is empty")
	}
	configMap, ok := configBlock[0].(map[string]interface{})
	if !ok {
		panic("config[0] is not a map")
	}

	// Extract Network values from the map
	network, ok := pcDetails["network"].([]interface{})
	if !ok || len(network) == 0 {
		panic("network is not a slice or is empty")
	}
	networkMap, ok := network[0].(map[string]interface{})
	if !ok {
		panic("network[0] is not a map")
	}

	configString := expandPCConfigBlock(configMap)
	networkString := expandPCNetworkBlock(networkMap)

	// Generate 9 unique passwords.
	const numPasswords = 9
	uniquePasswords := make(map[string]struct{})

	for len(uniquePasswords) < numPasswords {
		pass, err := generatePassword()
		if err != nil {
			log.Fatalf("Error generating password: %v", err)
		}
		uniquePasswords[pass] = struct{}{}
	}

	// Build remote commands to reset the admin password.
	remoteCommands := ""
	for pass := range uniquePasswords {
		cmd := fmt.Sprintf("/home/nutanix/prism/cli/ncli user reset-password user-name=%s password=%s", testVars.Prism.PCRestore.Username, pass)
		remoteCommands += cmd + " ; "
	}

	// Append a fallback command using the previous password.
	fallbackCmd := fmt.Sprintf("/home/nutanix/prism/cli/ncli user reset-password user-name=%s password=%s", testVars.Prism.PCRestore.Username, testVars.Prism.PCRestore.Password)
	remoteCommands += fallbackCmd

	// Build the two remote password reset commands.
	//Retrieve environment variables.
	pcIP := os.Getenv("NUTANIX_ENDPOINT")

	// Build the full SSH command. Note the single quotes around the remoteCommands.
	resetCommand := fmt.Sprintf("ssh-keygen -f ~/.ssh/known_hosts -R %[3]s; sshpass -p '%[1]s' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null %[2]s@%[3]s '%[4]s'",
		testVars.Prism.RestoreSource.SSHPassword, testVars.Prism.RestoreSource.SSHUser, pcIP, remoteCommands)

	// pe config
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	peHostProviderConfig := fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}
`, username, password, testVars.Prism.RestoreSource.PeIP, insecure, port)

	return fmt.Sprintf(`
# peHostProviderConfig
%[1]s

resource "nutanix_pc_restore_v2" "test" {
	provider                         = nutanix-2
	timeouts {
		create = "140m"
	}
	ext_id                           = "%[2]s"
	restore_source_ext_id            = "%[3]s"
	restorable_domain_manager_ext_id = "%[4]s"
	domain_manager {
		# Config
		%[5]s

		# Network
		%[6]s
	}
	provisioner "local-exec" {
		command = "%[7]s"
		on_failure = continue
	}
}
`, peHostProviderConfig, restorePointExtID, restoreSourceExtID, restorablePcExtID, configString, networkString, resetCommand)
}
