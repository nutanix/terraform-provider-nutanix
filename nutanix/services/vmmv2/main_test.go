package vmmv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	// Volumes config
	VMM struct {
		ImageName  string `json:"image_name"`
		AssignedIP string `json:"assigned_ip"`
		Subnet     struct {
			NetworkID int    `json:"network_id"`
			IP        string `json:"ip"`
			Prefix    int    `json:"prefix"`
			GatewayIP string `json:"gateway_ip"`
			StartIP   string `json:"start_ip"`
			EndIP     string `json:"end_ip"`
		}
		StorageContainer string `json:"storage_container"`
		NGT              struct {
			Credential struct {
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"credential"`
		} `json:"ngt"`
		GPUS []struct {
			DeviceID int    `json:"device_id"`
			Mode     string `json:"mode"`
			Vendor   string `json:"vendor"`
		} `json:"gpus"`
	} `json:"vmm"`
}

var testVars TestConfig

func loadVars(filepath string, varStuct interface{}) {
	// Read test_config_v2.json from home current path
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
