package objectstoresv2_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	import1 "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/config"
)

type TestConfig struct {
	ObjectStore struct {
		ImgURL              string   `json:"img_url"`
		SubnetName          string   `json:"subnet_name"`
		BucketName          string   `json:"bucket_name"`
		Domain              string   `json:"domain"`
		PublicNetworkIPs    []string `json:"public_network_ips"`
		StorageNetworkDNSIP string   `json:"storage_network_dns_ip"`
		StorageNetworkVip   string   `json:"storage_network_vip"`
	} `json:"object_store"`
}

var testVars TestConfig

var (
	path, _             = os.Getwd()
	filepath            = path + "/../../../test_config_v2.json"
	certificateJSONFile = path + "/../../../object_store_cert.json"
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
	createCertificateJSONFile()
	os.Exit(m.Run())
}

func createCertificateJSONFile() error {
	alternateIps := testVars.ObjectStore.PublicNetworkIPs

	certificate := config.NewCertificate()
	if len(alternateIps) > 0 {
		certificate.AlternateIps = make([]import1.IPAddress, 1)
		certificate.AlternateIps[0] = import1.IPAddress{
			Ipv4: &import1.IPv4Address{
				Value: utils.StringPtr(alternateIps[0]),
			},
		}
	}

	// Marshal the certificate data to JSON
	certificateJSON, err := json.MarshalIndent(certificate, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal certificate data: %w", err)
	}
	// Write the JSON data to a file
	err = os.WriteFile(certificateJSONFile, certificateJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}
	log.Printf("Certificate JSON file created at: %s", certificateJSONFile)
	return nil
}
