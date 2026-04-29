package monitoringv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	AlertExtID string `json:"alert_ext_id"`
}

var testVars TestConfig

func loadVars(filepath string, varStruct interface{}) {
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Got this error while reading config.json: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(configData, varStruct)
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
