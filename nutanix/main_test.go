package nutanix

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	SubnetName                     string `json:"subnet_name"`
	DefaultContainerName           string `json:"default_container_name"`
	UserGroupWithDistinguishedName struct {
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
	NodeOsVersion string `json:"node_os_version"`
	AdRuleTarget  struct {
		Name   string `json:"name"`
		Values string `json:"values"`
	} `json:"ad_rule_target"`
}

var testVars TestConfig

func TestMain(m *testing.M) {
	log.Println("Do some crazy stuff before tests!")

	// Read config.json from home current path
	configData, err := os.ReadFile("../test_config.json")
	if err != nil {
		log.Printf("Got this error while reading config.json: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(configData, &testVars)
	if err != nil {
		log.Printf("Got this error while unmarshalling config.json: %s", err.Error())
		os.Exit(1)
	}

	log.Println(testVars)

	os.Exit(m.Run())
}
