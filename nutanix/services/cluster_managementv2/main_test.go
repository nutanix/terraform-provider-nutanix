package cluster_managementv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	Clusters struct {
		PEUsername string `json:"pe_username"`
		PEPassword string `json:"pe_password"`
		PcExtID    string `json:"pc_ext_id"`
		Network    struct {
			SMTPServer struct {
				IP   string `json:"ip"`
				Port int    `json:"port"`
			} `json:"smtp_server"`
		} `json:"network"`
	} `json:"clusters"`
}

var (
	testVars TestConfig
	path, _  = os.Getwd()
	filepath = path + "/../../../test_config_v2.json"
)

func loadVars(filepath string, varStuct interface{}) {
	configData, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Got this error while reading test_config_v2.json: %s", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(configData, varStuct)
	if err != nil {
		log.Printf("Got this error while unmarshalling test_config_v2.json: %s", err.Error())
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	log.Println("Do some crazy stuff before tests!")
	loadVars(filepath, &testVars)
	os.Exit(m.Run())
}
