package prismv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	Prism struct {
		DeployPC struct {
			PeIP           string `json:"pe_ip"`
			Version        string `json:"version"`
			DefaultGateway string `json:"default_gateway"`
			SubnetMask     string `json:"subnet_mask"`
			IPRange        struct {
				Begin string `json:"begin"`
				End   string `json:"end"`
			} `json:"ip_range"`
			NameServers []string `json:"name_servers"`
			NtpServers  []string `json:"ntp_servers"`
		} `json:"deploy_pc"`
		Bucket struct {
			Name      string `json:"name"`
			Region    string `json:"region"`
			AccessKey string `json:"access_key"`
			SecretKey string `json:"secret_key"`
		} `json:"bucket"`
		RestoreSource struct {
			PeIP        string `json:"pe_ip"`
			SSHPassword string `json:"ssh_password"`
			SSHUser     string `json:"ssh_user"`
		} `json:"restore_source"`
		Unregister struct {
			PcExtID string `json:"pc_ext_id"`
		} `json:"unregister"`
		PCRestore struct {
			Username          string `json:"username"`
			Password          string `json:"password"`
			SkipPCRestoreTest bool   `json:"skip_pc_restore_test"`
		} `json:"pc_restore"`
	} `json:"prism"`
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
