package prismv2_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRestorePC = "nutanix_restore_pc_v2.test"

func TestAccV2NutanixRestorePCResource_RestorePC(t *testing.T) {
	var backupTargetExtID,
		pcExtID = new(string), new(string)

	var restorePcConfig, clusterExtID string

	pcDetails := make(map[string]interface{})

	t.Run("create_pre_request", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreventPostDestroyRefresh: true,
			PreCheck:                  func() { acc.TestAccPreCheck(t) },
			Providers:                 acc.TestAccProviders,
			Steps: []resource.TestStep{
				// Step 1: List backup targets and delete if backup target exists
				{
					PreConfig: func() {
						fmt.Printf("Step 1: List backup targets and delete if backup target exists\n")
					},
					Config: testAccListBackupTargetsDatasourceConfig(),
					Check: resource.ComposeTestCheckFunc(
						checkBackupTargetExist(),
					),
				},
				// Step 2: Create backup target, Restore Source
				{
					PreConfig: func() {
						fmt.Printf("Step 2: Create backup target, Restore Source\n")
					},
					Config: testAccPreRequestForRestorePcConfig(),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							pcDetailsOutput, ok := s.RootModule().Outputs["pc_details"]
							if !ok {
								return fmt.Errorf("output 'pc_details' not found")
							}
							clusterExtIDOutput, ok := s.RootModule().Outputs["clusterExtID"]
							if !ok {
								return fmt.Errorf("output 'clusterExtID' not found")
							}
							clusterExtID = clusterExtIDOutput.Value.(string)

							pcDetails = pcDetailsOutput.Value.(map[string]interface{})
							return nil
						},
						createBackupTarget(backupTargetExtID),
						checkLastSyncTimeBackupTargetRestorePC(backupTargetExtID, pcExtID, retries, delay),
					),
				},
				// Step 3: power off PC
				{
					PreConfig: func() {
						fmt.Printf("Step 3: Power off PC\n")
					},
					Config: testAccPowerOffPCConfig(),
					Check: resource.ComposeTestCheckFunc(
						func(s *terraform.State) error {
							// Build the restore PC configuration for the next subtest case
							restorePcConfig = restorePcResourceConfig(pcDetails, clusterExtID)
							fmt.Printf("Restore PC Config: %s\n", restorePcConfig)
							return nil
						},
						powerOffPC(),
					),
				},
			},
		})
	})

	// Restore PC Subtest Case
	t.Run("restore_pc", func(t *testing.T) {
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

						//cleanup
						deleteBackupTarget(backupTargetExtID, pcExtID),
					),
				},
			},
		})
	})
}

func testAccPreRequestForRestorePcConfig() string {
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

data "nutanix_pc_v2" "test" {
  ext_id = local.domainManagerExtId
}

output "pc_details" {
  value = data.nutanix_pc_v2.test
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

func testAccPowerOffPCConfig() string {
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

%s
# Dummy data source to make sure the the second provider is initialized
data "nutanix_subnets_v2" "subnets" {
    provider = nutanix-2
}
`, peHostProviderConfig)
}

func restorePcResourceConfig(pcDetails map[string]interface{}, clusterExtID string) string {
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
		cmd := fmt.Sprintf("/home/nutanix/prism/cli/ncli user reset-password user-name=%s password=%s", "admin", pass)
		remoteCommands += cmd + " ; "
	}

	// Append a fallback command using the previous password.
	fallbackCmd := fmt.Sprintf("/home/nutanix/prism/cli/ncli user reset-password user-name=%s password=%s", "admin", "Nutanix.123")
	remoteCommands += fallbackCmd

	username := "nutanix"
	// Build the two remote password reset commands.
	//Retrieve environment variables.
	pcIP := os.Getenv("NUTANIX_ENDPOINT")

	// Build the full SSH command. Note the single quotes around the remoteCommands.
	resetCommand := fmt.Sprintf("sshpass -p '%s' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null %s@%s '%s'",
		testVars.Prism.RestoreSource.PcPassword, username, pcIP, remoteCommands)

	// pe config
	username = os.Getenv("NUTANIX_USERNAME")
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

resource "nutanix_restore_source_v2" "cluster-location" {
  provider =  nutanix-2
  location {
    cluster_location {
      config {
		# clusterExtID
        ext_id = "%s"
      }
    }
  }
}

data "nutanix_restorable_pcs_v2" "restorable-pcs" {
  provider              = nutanix-2
  restore_source_ext_id = nutanix_restore_source_v2.cluster-location.ext_id
}

locals {
  restorablePc = data.nutanix_restorable_pcs_v2.restorable-pcs.restorable_pcs.0
}

data "nutanix_restore_points_v2" "restore-points" {
  provider                         = nutanix-2
  restorable_domain_manager_ext_id = local.restorablePc.ext_id
  restore_source_ext_id            = nutanix_restore_source_v2.cluster-location.id
}

locals {
  restorePointId = data.nutanix_restore_points_v2.restore-points.restore_points[0].ext_id
}


resource "nutanix_restore_pc_v2" "test" {
	provider                         = nutanix-2
	timeouts {
		create = "120m"
	}
    ext_id                           = local.restorePointId
    restore_source_ext_id            = nutanix_restore_source_v2.cluster-location.id
    restorable_domain_manager_ext_id = local.restorablePc.ext_id
	domain_manager {
		config {
			should_enable_lockdown_mode = false
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
	}
	provisioner "local-exec" {
		command = "%s"		
		on_failure = continue
	}
}
				`, peHostProviderConfig, clusterExtID, version, name, size, strContainerExtIDs, dataDiskSizeBytes, memorySizeBytes, numVcpus,
		externalAddressIPv4Value, nameServersConfig, ntpServersConfig, networkExtID, defaultGatewayIPv4, subnetMaskIPv4,
		ipRangeBeginIPv4, ipRangeEndIPv4, resetCommand)
}
