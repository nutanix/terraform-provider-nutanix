package monitoringv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	Monitoring struct {
		SdaPolicyExtID string `json:"sda_policy_ext_id"`
		ClusterExtID   string `json:"cluster_ext_id"`
	} `json:"monitoring"`
}

var (
	testVars TestConfig
	path, _  = os.Getwd()
	filepath = path + "/../../../test_config_v2.json"
)

func loadVars(filepath string, varStuct interface{}) {
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
