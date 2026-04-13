package securityv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	Security struct {
		KMS struct {
			EndpointURL  string `json:"endpoint_url"`
			TenantID     string `json:"tenant_id"`
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			KeyID        string `json:"key_id"`
		} `json:"kms"`
	} `json:"security"`
}

var (
	testVars TestConfig
	path, _  = os.Getwd()
	filepath = path + "/../../../test_config_v2.json"
)

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
	loadVars("../../../test_config_v2.json", &testVars)
	os.Exit(m.Run())
}
