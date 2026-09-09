package lifecyclev2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

// TestConfig mirrors the shared test_config_v2.json structure. The Foundation
// Central configuration entity has no infrastructure parameters of its own; the
// struct is kept for parity with the other v2 service packages and future use.
type TestConfig struct {
	Lcm struct {
		EntityModel        string `json:"entity_model"`
		EntityModelVersion string `json:"entity_model_version"`
	} `json:"lcm"`
}

var testVars TestConfig

func loadVars(filepath string, varStruct interface{}) {
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Got this error while reading test_config_v2.json: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(configData, varStruct)
	if err != nil {
		log.Printf("Got this error while unmarshalling test_config_v2.json: %s", err.Error())
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	log.Println("Setup for lifecyclev2 acceptance tests")
	loadVars("../../../test_config_v2.json", &testVars)
	os.Exit(m.Run())
}
