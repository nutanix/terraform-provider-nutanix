package vmmv2_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

type TestConfig struct {
	// Volumes config
	VMM struct {
		ImageName   string `json:"image_name"`
		AssignedIP  string `json:"assigned_ip"`
		UnattendXML string `json:"unattend_xml"`
		Subnet      struct {
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
		OvaURL string `json:"ova_url"`
	} `json:"vmm"`
}

var testVars TestConfig
var (
	path, _             = os.Getwd()
	filepath            = path + "/../../../test_config_v2.json"
	untendedXMLFilePath = path + "/../../../unattendxml.txt"
)

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
	downloadFile(testVars.VMM.UnattendXML, untendedXMLFilePath)
	os.Exit(m.Run())
}

// downloadFile downloads a file from a given URL and saves it to the specified path.
func downloadFile(url, destinationFilePath string) error {
	// Send HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the destination file
	out, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy data from response to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}
