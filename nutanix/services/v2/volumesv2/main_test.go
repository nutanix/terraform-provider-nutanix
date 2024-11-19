package volumesv2_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	// Volumes config
	Volumes struct {
		VolumeGroupExtIdWithCategory string `json:"vg_ext_id_with_category"`
		VmExtId                      string `json:"vm_ext_id"`
		IscsiClient                  struct {
			ExtId              string `json:"ext_id"`
			IscsiInitiatorName string `json:"iscsi_initiator_name"`
		} `json:"iscsi_client"`
		Disk struct {
			DiskDataSourceReference struct {
				ExtId string `json:"ext_id"`
			} `json:"disk_data_source_reference"`
		} `json:"disk"`
	} `json:"volumes"`
}

var testVars TestConfig

func loadVars(filepath string, varStuct interface{}) {
	// Read test_config_v4.json from home current path
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
	loadVars("../../../../test_config_v2.json", &testVars)
	os.Exit(m.Run())
}
