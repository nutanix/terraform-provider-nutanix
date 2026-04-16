package selfservice_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	SelfService struct {
		BlueprintName             string `json:"bp_name"`
		BlueprintWithSnapshotName string `json:"bp_name_with_snapshot_config"`
		AppWithSnapshotName       string `json:"app_name_with_snapshot_config"`
	} `json:"self_service"`
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

	os.Exit(m.Run())
}
