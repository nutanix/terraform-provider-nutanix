package ndb_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	SubnetName                     string `json:"subnet_name"`
	DefaultContainerName           string `json:"default_container_name"`
	UserGroupWithDistinguishedName []struct {
		DistinguishedName string `json:"distinguished_name"`
		DisplayName       string `json:"display_name"`
		UUID              string `json:"uuid"`
	} `json:"user_group_with_distinguished_name"`
	Permissions []struct {
		Name string `json:"name"`
		UUID string `json:"uuid"`
	} `json:"permissions"`
	Users []struct {
		PrincipalName        string `json:"principal_name"`
		ExpectedDisplayName  string `json:"expected_display_name"`
		DirectoryServiceUUID string `json:"directory_service_uuid"`
	} `json:"users"`
	KubernetesVersion string `json:"kubernetes_version"`
	NodeOsVersion     string `json:"node_os_version"`
	AdRuleTarget      struct {
		Name   string `json:"name"`
		Values string `json:"values"`
	} `json:"ad_rule_target"`
	// here UUID = availability_zone_url
	ProtectionPolicy struct {
		LocalAz struct {
			UUID        string `json:"uuid"`
			ClusterUUID string `json:"cluster_uuid"`
		} `json:"local_az"`
		DestinationAz struct {
			UUID        string `json:"uuid"`
			ClusterUUID string `json:"cluster_uuid"`
		} `json:"destination_az"`
	} `json:"protection_policy"`
	// sshKey required for ndb database provision test
	SSHKey string `json:"ssh_key"`
	// NDB config
	NDB struct {
		RegisterClusterInfo struct {
			ClusterIP        string `json:"cluster_ip"`
			Username         string `json:"username"`
			Password         string `json:"password"`
			DNS              string `json:"dns"`
			NTP              string `json:"ntp"`
			StaticIP         string `json:"static_ip"`
			Gateway          string `json:"gateway"`
			SubnetMask       string `json:"subnet_mask"`
			StorageContainer string `json:"strorage_container"`
		} `json:"register_cluster_info"`
		TestStaticNetwork string `json:"test_static_network"`
	} `json:"ndb"`
}

var testVars TestConfig

func loadVars(filepath string, varStuct interface{}) {
	// Read config.json from home current path
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Got this error while reading config.json: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(configData, varStuct)
	if err != nil {
		log.Printf("Got this error while unmarshalling config.json: %s", err.Error())
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	log.Println("Do some crazy stuff before tests!")
	loadVars("../../../test_config.json", &testVars)

	// Fallback for SSH Key from Environment Variables if JSON is empty
	if testVars.SSHKey == "" {
		testVars.SSHKey = os.Getenv("NDB_SSH_KEY")
	}
	if testVars.SSHKey == "" {
		testVars.SSHKey = os.Getenv("SSH_KEY")
	}

	// Final fallback: Use a dummy key if still empty (ensures CI passes provider validation)
	if testVars.SSHKey == "" {
		log.Println("[INFO] No SSH key found in config or env. Using dummy key for CI validation.")
		testVars.SSHKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8Z2... dummy-ci-key"
	}

	os.Exit(m.Run())
}
