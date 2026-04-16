package dataprotectionv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	AvailabilityZone struct {
		PcExtID      string `json:"pc_ext_id"`
		ClusterExtID string `json:"cluster_ext_id"`
		RemotePcIP   string `json:"remote_pc_ip"`
	} `json:"availability_zone"`
	DataProtection struct {
		LocalClusterPE   string `json:"local_cluster_pe"`
		LocalClusterVIP  string `json:"local_cluster_vip"`
		RemoteClusterPE  string `json:"remote_cluster_pe"`
		RemoteClusterVIP string `json:"remote_cluster_vip"`
	} `json:"data_protection"`
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
	loadVars("../../../test_config_v2.json", &testVars)
	os.Exit(m.Run())
}
