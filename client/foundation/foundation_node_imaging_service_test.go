package foundation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func setup() (*http.ServeMux, *client.Client, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c, _ := client.NewClient(&client.Credentials{
		URL:      "",
		Username: "username",
		Password: "password",
		Port:     "",
		Endpoint: "0.0.0.0",
		Insecure: true},
		userAgent,
		absolutePath,
		true)
	c.UserAgent = userAgent
	c.BaseURL, _ = url.Parse(server.URL)

	return mux, c, server
}

func TestOperation_ImageNodes(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()

	mux.HandleFunc("/foundation/image_nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Request method = %v, expected %v", r.Method, http.MethodPost)
		}

		expected := map[string]interface{}{
			"ipmi_password":      "test_password",
			"ipmi_user":          "test_user",
			"cvm_gateway":        "0.0.0.0",
			"cvm_netmask":        "255.255.255.0",
			"hypervisor_gateway": "0.0.0.0",
			"hypervisor_netmask": "255.255.255.0",
			"nos_package":        "test_nos.tar.gz",
			"hypervisor_iso":     map[string]interface{}{},
			"blocks": []interface{}{
				map[string]interface{}{
					"block_id": "N123",
					"nodes": []interface{}{
						map[string]interface{}{
							"ipmi_configure_now":  true,
							"ipmi_ip":             "0.0.0.0",
							"cvm_ip":              "0.0.0.0",
							"hypervisor_ip":       "0.0.0.0",
							"image_now":           true,
							"ipmi_password":       "test_password",
							"ipmi_user":           "test_user",
							"hypervisor_hostname": "test_hostname",
							"hypervisor":          "kvm",
							"node_position":       "A",
						},
					},
				},
			},
			"clusters": []interface{}{
				map[string]interface{}{
					"redundancy_factor":   1,
					"cluster_init_now":    true,
					"cluster_external_ip": interface{}(nil),
					"cluster_name":        "test_cluster",
					"cluster_members":     []interface{}{"0.0.0.0"},
				},
			},
		}

		// checks
		var b map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}
		if !reflect.DeepEqual(b, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", b, expected)
		}

		// mock response
		fmt.Fprintf(w, `{
			"session_id" : "123456-1234-123456"
		}`)
	})
	ctx := context.TODO()
	inp := &ImageNodesInput{
		IpmiPassword:      "test_password",
		IpmiUser:          "test_user",
		CvmGateway:        "0.0.0.0",
		CvmNetmask:        "255.255.255.0",
		HypervisorGateway: "0.0.0.0",
		HypervisorNetmask: "255.255.255.0",
		NosPackage:        "test_nos.tar.gz",
		Blocks: []*Block{
			{
				BlockID: "N123",
				Nodes: []*Node{
					{
						IpmiConfigureNow:   utils.BoolPtr(true),
						IpmiIP:             "0.0.0.0",
						IpmiUser:           "test_user",
						IpmiPassword:       "test_password",
						CvmIP:              "0.0.0.0",
						ImageNow:           utils.BoolPtr(true),
						HypervisorIP:       "0.0.0.0",
						HypervisorHostname: "test_hostname",
						Hypervisor:         "kvm",
						NodePosition:       "A",
					},
				},
			},
		},
		Clusters: []*Clusters{
			{
				RedundancyFactor: utils.Int64Ptr(1),
				ClusterInitNow:   utils.BoolPtr(true),
				ClusterName:      "test_cluster",
				ClusterMembers:   []string{"0.0.0.0"},
			},
		},
	}

	out := &ImageNodesAPIResponse{
		SessionID: "123456-1234-123456",
	}

	op := NodeImagingOperations{
		client: c,
	}

	// checks
	got, err := op.ImageNodes(ctx, inp)
	if err != nil {
		t.Errorf("NodeImagingOperations.ImageNodes() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("NodeImagingOperations.ImageNodes() = %+v, want %+v", got, out)
	}

}
