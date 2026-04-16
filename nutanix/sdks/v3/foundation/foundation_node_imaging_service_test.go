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

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func setup() (*http.ServeMux, *client.Client, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c, _ := client.NewBaseClient(&client.Credentials{
		URL:      "",
		Username: "username",
		Password: "password",
		Port:     "",
		Endpoint: "0.0.0.0",
		Insecure: true},
		absolutePath,
		true)
	c.UserAgent = userAgent
	c.BaseURL, _ = url.Parse(server.URL)

	return mux, c, server
}

func testHTTPMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func TestNodeImagingOperations_ImageNodes(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()

	mux.HandleFunc("/foundation/image_nodes", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

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
					"redundancy_factor":   float64(1),
					"cluster_init_now":    true,
					"cluster_external_ip": nil,
					"cluster_name":        "test_cluster",
					"cluster_members":     []interface{}{"0.0.0.0"},
				},
			},
		}

		// checks
		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
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
		t.Fatalf("NodeImagingOperations.ImageNodes() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("NodeImagingOperations.ImageNodes() got = %#v, want = %#v", got, out)
	}
}

func TestNodeImagingOperations_ImageNodesProgress(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	sessionID := "123456-1234-123456"
	mux.HandleFunc("/foundation/progress", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)

		// mock response
		fmt.Fprintf(w, `{
			"session_id": "%v",
			"imaging_stopped": true,
			"aggregate_percent_complete": 100.00,
			"clusters": [{
				"cluster_name": "test_cluster",
				"time_elapsed": 102.33,
				"cluster_members": [
					"0.0.0.0"
				],
				"percent_complete": 100.00
			}],
			"nodes": [{
				"cvm_ip": "0.0.0.0",
				"hypervisor_ip": "0.0.0.0",
				"time_elapsed": 102.33,
				"percent_complete": 100.00
			}]
		}`, sessionID)
	})
	ctx := context.TODO()

	out := &ImageNodesProgressResponse{
		SessionID:                "123456-1234-123456",
		ImagingStopped:           utils.BoolPtr(true),
		AggregatePercentComplete: utils.Float64Ptr(100.00),
		Clusters: []*ClusterProgress{
			{
				ClusterName:     "test_cluster",
				TimeElapsed:     utils.Float64Ptr(102.33),
				ClusterMembers:  []string{"0.0.0.0"},
				PercentComplete: utils.Float64Ptr(100.00),
			},
		},
		Nodes: []*NodeProgress{
			{
				CvmIP:           "0.0.0.0",
				HypervisorIP:    "0.0.0.0",
				TimeElapsed:     utils.Float64Ptr(102.33),
				PercentComplete: utils.Float64Ptr(100.00),
			},
		},
	}

	op := NodeImagingOperations{
		client: c,
	}

	// checks
	got, err := op.ImageNodesProgress(ctx, sessionID)
	if err != nil {
		t.Fatalf("NodeImagingOperations.ImageNodesProgress() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("NodeImagingOperations.ImageNodesProgress() got = %#v, want = %#v", got, out)
	}
}
