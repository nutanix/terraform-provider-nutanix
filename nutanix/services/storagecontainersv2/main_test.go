package storagecontainersv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	StorageContainer struct {
		LogicalAdvertisedCapacityBytes       int `json:"logical_advertised_capacity_bytes"`
		LogicalExplicitReservedCapacityBytes int `json:"logical_explicit_reserved_capacity_bytes"`
		ReplicationFactor                    int `json:"replication_factor"`
		NfsWhitelistAddresses                struct {
			Ipv4 struct {
				Value        string `json:"value"`
				PrefixLength int    `json:"prefix_length"`
			} `json:"ipv4"`
		} `json:"nfs_whitelist_addresses"`
	} `json:"storage_container"`
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
