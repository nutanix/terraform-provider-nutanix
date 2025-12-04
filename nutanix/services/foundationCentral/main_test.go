package fc_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type IPMIConfig struct {
	IpmiGateway  string `json:"ipmi_gateway"`
	IpmiNetmask  string `json:"ipmi_netmask"`
	IpmiUser     string `json:"ipmi_user"`
	IpmiPassword string `json:"ipmi_password"`
	IpmiIP       string `json:"ipmi_ip"`
	IpmiMac      string `json:"ipmi_mac"`
}

type FoundationVars struct {
	IPv6Addresses []string   `json:"ipv6_addresses"`
	IpmiConfig    IPMIConfig `json:"ipmi_config"`
	Blocks        []struct {
		Nodes []struct {
			IpmiIP                  string `json:"ipmi_ip"`
			IpmiPassword            string `json:"ipmi_password"`
			IpmiUser                string `json:"ipmi_user"`
			IpmiNetmask             string `json:"ipmi_netmask"`
			IpmiGateway             string `json:"ipmi_gateway"`
			CvmIP                   string `json:"cvm_ip"`
			HypervisorIP            string `json:"hypervisor_ip"`
			Hypervisor              string `json:"hypervisor"`
			HypervisorHostname      string `json:"hypervisor_hostname"`
			NodePosition            string `json:"node_position"`
			IPv6Address             string `json:"ipv6_address"`
			CurrentNetworkInterface string `json:"current_network_interface"`
			ImagedNodeUUID          string `json:"imaged_node_uuid"`
			HypervisorType          string `json:"hypervisor_type"`
		} `json:"nodes"`
		BlockID                    string `json:"block_id"`
		CvmGateway                 string `json:"cvm_gateway"`
		HypervisorGateway          string `json:"hypervisor_gateway"`
		CvmNetmask                 string `json:"cvm_netmask"`
		HypervisorNetmask          string `json:"hypervisor_netmask"`
		IpmiUser                   string `json:"ipmi_user"`
		AosPackageURL              string `json:"aos_package_url"`
		UseExistingNetworkSettings bool   `json:"use_existing_network_settings"`
		ImageNow                   bool   `json:"image_now"`
		CommonNetworkSettings      struct {
			CvmDNSServers        []string `json:"cvm_dns_servers"`
			HypervisorDNSServers []string `json:"hypervisor_dns_servers"`
			CvmNtpServers        []string `json:"cvm_ntp_servers"`
			HypervisorNtpServers []string `json:"hypervisor_ntp_servers"`
		} `json:"common_network_settings"`
	} `json:"blocks"`
	OnboardNodes []struct {
		NodeSerial string `json:"node_serial"`
	} `json:"onboard_nodes"`
}

var foundationVars FoundationVars

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
	loadVars("../../../test_foundation_config.json", &foundationVars)

	os.Exit(m.Run())
}
