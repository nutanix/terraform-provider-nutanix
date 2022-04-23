package foundation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestNtwOperations_DiscoverNodes(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	mux.HandleFunc("/foundation/discover_nodes", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)

		// mock response
		fmt.Fprintf(w, `[{
			"model": "XCV10",
			"nodes": [{
				"node_position": "A",
				"hypervisor": "kvm",
				"svm_ip": "0.0.0.0",
				"configured": true
			}],
			"block_id": "GMD1"
		}, {
			"model": "XCV10",
			"nodes": [{
					"node_position": "A",
					"hypervisor": "kvm",
					"svm_ip": "0.0.0.0",
					"configured": false
				}
			],
			"block_id": "GMD2"
		}]`)
	})
	ctx := context.TODO()

	out := &DiscoverNodesAPIResponse{
		{
			Model: "XCV10",
			Nodes: []DiscoveredNode{
				{
					NodePosition: "A",
					Hypervisor:   "kvm",
					SvmIP:        "0.0.0.0",
					Configured:   utils.BoolPtr(true),
				},
			},
			BlockID: "GMD1",
		},
		{
			Model: "XCV10",
			Nodes: []DiscoveredNode{
				{
					NodePosition: "A",
					Hypervisor:   "kvm",
					SvmIP:        "0.0.0.0",
					Configured:   utils.BoolPtr(false),
				},
			},
			BlockID: "GMD2",
		},
	}

	op := NetworkingOperations{
		client: c,
	}

	// checks
	got, err := op.DiscoverNodes(ctx)
	if err != nil {
		t.Fatalf("NetworkingOperations.DiscoverNodes() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("NetworkingOperations.DiscoverNodes() got = %#v, want = %#v", got, out)
	}
}

func TestNtwOperations_NodeNetworkDetails(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	mux.HandleFunc("/foundation/node_network_details", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"nodes": []interface{}{
				map[string]interface{}{
					"ipv6_address": "ffff::ffff:fffff:ffff",
				},
				map[string]interface{}{
					"ipv6_address": "ec12::ec12:ec12:ec12",
				},
			},
			"timeout": "30",
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
			"nodes" : [
				{	
					"cvm_ip" : "0.0.0.0",
					"node_serial" : "NX1234",
					"ipmi_ip" : "0.0.0.0"
				},
				{	
					"cvm_ip" : "0.0.0.0",
					"node_serial" : "NX1235",
					"ipmi_ip" : "0.0.0.0"
				}
			]
		}`)
	})
	ctx := context.TODO()
	inp := &NodeNetworkDetailsInput{
		Nodes: []NodeIpv6Input{
			{
				Ipv6Address: "ffff::ffff:fffff:ffff",
			},
			{
				Ipv6Address: "ec12::ec12:ec12:ec12",
			},
		},
		Timeout: "30",
	}
	out := &NodeNetworkDetailsResponse{
		Nodes: []NodeNetworkDetail{
			{
				CvmIP:      "0.0.0.0",
				NodeSerial: "NX1234",
				IpmiIP:     "0.0.0.0",
			},
			{
				CvmIP:      "0.0.0.0",
				NodeSerial: "NX1235",
				IpmiIP:     "0.0.0.0",
			},
		},
	}

	op := NetworkingOperations{
		client: c,
	}

	// checks
	got, err := op.NodeNetworkDetails(ctx, inp)
	if err != nil {
		t.Fatalf("NetworkingOperations.NodeNetworkDetails() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("NetworkingOperations.NodeNetworkDetails() got = %#v, want = %#v", got, out)
	}
}

func TestNtwOperations_ConfigureIPMI(t *testing.T) {
	mux, c, server := setup()
	defer server.Close()
	mux.HandleFunc("/foundation/ipmi_config", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]interface{}{
			"blocks": []interface{}{
				map[string]interface{}{
					"nodes": []interface{}{
						map[string]interface{}{
							"ipmi_ip":            "0.0.0.0",
							"ipmi_mac":           "ac:da:af:fa:af:fa",
							"ipmi_configure_now": true,
						},
					},
					"block_id": "GMD10",
				},
			},
			"ipmi_netmask":  "255.255.255.0",
			"ipmi_gateway":  "0.0.0.0",
			"ipmi_user":     "username",
			"ipmi_password": "password",
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
			"blocks": [
				{
					"nodes":[ 
						{
							"ipmi_ip":            "0.0.0.0",
							"ipmi_mac":           "ac:da:af:fa:af:fa",
							"ipmi_configure_now": true,
							"ipmi_configure_successful": true,
							"ipmi_message" : "success"
						}
					],
					"block_id": "GMD10"
				}
			],
			"ipmi_netmask":  "255.255.255.0",
			"ipmi_gateway":  "0.0.0.0",
			"ipmi_user":     "username",
			"ipmi_password": "password"
		}`)
	})
	ctx := context.TODO()
	inp := &IPMIConfigAPIInput{
		IpmiUser:     "username",
		IpmiPassword: "password",
		IpmiNetmask:  "255.255.255.0",
		IpmiGateway:  "0.0.0.0",
		Blocks: []IPMIConfigBlockInput{
			{
				Nodes: []IPMIConfigNodeInput{
					{
						IpmiIP:           "0.0.0.0",
						IpmiMac:          "ac:da:af:fa:af:fa",
						IpmiConfigureNow: true,
					},
				},
				BlockID: "GMD10",
			},
		},
	}
	out := &IPMIConfigAPIResponse{
		IpmiUser:     "username",
		IpmiPassword: "password",
		IpmiNetmask:  "255.255.255.0",
		IpmiGateway:  "0.0.0.0",
		Blocks: []IPMIConfigBlockResponse{
			{
				Nodes: []IPMIConfigNodeResponse{
					{
						IpmiIP:                  "0.0.0.0",
						IpmiMac:                 "ac:da:af:fa:af:fa",
						IpmiConfigureNow:        true,
						IpmiConfigureSuccessful: true,
						IpmiMessage:             "success",
					},
				},
				BlockID: "GMD10",
			},
		},
	}

	op := NetworkingOperations{
		client: c,
	}

	// checks
	got, err := op.ConfigureIPMI(ctx, inp)
	if err != nil {
		t.Fatalf("NetworkingOperations.ConfigureIPMI() error = %v", err)
	}
	if !reflect.DeepEqual(got, out) {
		t.Errorf("NetworkingOperations.ConfigureIPMI() got = %#v, want = %#v", got, out)
	}
}
