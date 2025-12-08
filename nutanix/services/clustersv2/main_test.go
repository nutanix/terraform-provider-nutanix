package clustersv2_test

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
		CvmIP      string `json:"cvm_ip"`
		Nodes      []struct {
			CvmIP    string `json:"cvm_ip"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"nodes"`
		Config struct {
			ClusterFunctions  []string `json:"cluster_functions"`
			AuthPublicKeyList []struct {
				Name string `json:"name"`
				Key  string `json:"key"`
			} `json:"auth_public_keys"`
			RedundancyFactor    int    `json:"redundancy_factor"`
			ClusterArch         string `json:"cluster_arch"`
			FaultToleranceState struct {
				DomainAwarenessLevel string `json:"domain_awareness_level"`
			} `json:"fault_tolerance_state"`
		} `json:"config"`
		Network struct {
			VirtualIP  string   `json:"virtual_ip"`
			IscsiIP    string   `json:"iscsi_ip"`
			IscsiIP1   string   `json:"iscsi_ip1"`
			DNSServers []string `json:"dns_servers"`
			NTPServers []string `json:"ntp_servers"`
			SMTPServer struct {
				IP           string `json:"ip"`
				Port         int    `json:"port"`
				Username     string `json:"username"`
				Password     string `json:"password"`
				Type         string `json:"type"`
				EmailAddress string `json:"email_address"`
			} `json:"smtp_server"`
		} `json:"network"`
		PcExtID       string `json:"pc_ext_id"`
		NodeIP        string `json:"node_ip"`
		RemoteCluster struct {
			ExtID    string `json:"ext_id"`
			IP       string `json:"ip"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"remote_cluster"`
		SSLCertificate struct {
			Passphrase        string `json:"passphrase"`
			PrivateKey        string `json:"private_key"`
			PublicCertificate string `json:"public_certificate"`
			CaChain           string `json:"ca_chain"`
		} `json:"ssl_certificate"`
	} `json:"clusters"`
}

var (
	testVars TestConfig
	path, _  = os.Getwd()
	filepath = path + "/../../../test_config_v2.json"
)

func loadVars(filepath string, varStuct interface{}) {
	// Read config.json from home current path
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
