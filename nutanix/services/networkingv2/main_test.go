package networkingv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	Networking struct {
		FloatingIP struct {
			VMNicReference string `json:"vm_nic_reference"`
		} `json:"floating_ip"`
		Subnets struct {
			ProjectID     string `json:"project_id"`
			VlanID        int    `json:"vlan_id"`
			NetworkIP     string `json:"network_ip"`
			NetworkPrefix int    `json:"network_prefix"`
			GatewayIP     string `json:"gateway_ip"`
			DHCP          struct {
				StartIP string `json:"start_ip"`
				EndIP   string `json:"end_ip"`
			}
		}
	} `json:"networking"`
}

var testVars TestConfig

var (
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
