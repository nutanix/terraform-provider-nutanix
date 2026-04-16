package foundationcentral

import (
	"context"
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
	c, _ := client.NewClient(&client.Credentials{
		URL:      "https://10.2.242.13:9440",
		Username: "admin",
		Password: "Nutanix.123",
		Port:     "9440",
		Endpoint: "10.2.242.13",
		Insecure: true,
	},
		userAgent,
		absolutePath,
		false)
	c.BaseURL, _ = url.Parse(server.URL)

	return mux, c, server
}

func testHTTPMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func TestOperations_ListImagedNodes(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/imaged_nodes/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"imaged_nodes":[{"node_state": "STATE_AVAILABLE"}]}`)
	})

	list := &ImagedNodesListResponse{}
	list.ImagedNodes = make([]*ImagedNodeDetails, 1)
	list.ImagedNodes[0] = &ImagedNodeDetails{}
	list.ImagedNodes[0].NodeState = utils.StringPtr("STATE_AVAILABLE")

	input := &ImagedNodesListInput{
		Length: utils.IntPtr(1),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *ImagedNodesListInput
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImagedNodesListResponse
		wantErr bool
	}{
		{
			"Test Imaged Nodes",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListImagedNodes(ctx, tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListImagedNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListImagedNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListImagedClusters(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/imaged_clusters/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"imaged_clusters":[{"cluster_name": "Test-Cluster"}]}`)
	})

	list := &ImagedClustersListResponse{}
	list.ImagedClusters = make([]*ImagedClusterDetails, 1)
	list.ImagedClusters[0] = &ImagedClusterDetails{}
	list.ImagedClusters[0].ClusterName = utils.StringPtr("Test-Cluster")

	input := &ImagedClustersListInput{
		Length: utils.IntPtr(1),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *ImagedClustersListInput
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImagedClustersListResponse
		wantErr bool
	}{
		{
			"Test Imaged Clusters",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListImagedClusters(ctx, tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListImagedClusters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListImagedClusters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_ListAPIKeys(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/api_keys/list", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
		fmt.Fprint(w, `{"api_keys":[{"api_key": "00a9-23h9", "alias":"test-1"}]}`)
	})

	list := &ListAPIKeysResponse{}
	list.APIKeys = make([]*CreateAPIKeysResponse, 1)
	list.APIKeys[0] = &CreateAPIKeysResponse{}
	list.APIKeys[0].APIKey = "00a9-23h9"
	list.APIKeys[0].Alias = "test-1"

	input := &ListMetadataInput{
		Length: utils.IntPtr(1),
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		getEntitiesRequest *ListMetadataInput
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListAPIKeysResponse
		wantErr bool
	}{
		{
			"Test List API Keys",
			fields{c},
			args{input},
			list,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.ListAPIKeys(ctx, tt.args.getEntitiesRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.ListAPIKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.ListAPIKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetImagedNode(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/imaged_nodes/0a8x-23d8", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"api_key_uuid": "234d-876f", "available": true, "cvm_ip": "10.0.0.0"}`)
	})

	node := &ImagedNodeDetails{}
	node.APIKeyUUID = utils.StringPtr("234d-876f")
	node.Available = utils.BoolPtr(true)
	node.CvmIP = utils.StringPtr("10.0.0.0")

	type fields struct {
		client *client.Client
	}

	type args struct {
		KeyUUID string
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImagedNodeDetails
		wantErr bool
	}{
		{
			"Get Imaged Node Details",
			fields{c},
			args{"0a8x-23d8"},
			node,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetImagedNode(ctx, tt.args.KeyUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetImagedNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetImagedNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetImagedCluster(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/imaged_clusters/0a8x-23d8", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"cluster_name": "test-cluster", "archived": true, "cluster_external_ip": "10.0.0.0"}`)
	})

	cluster := &ImagedClusterDetails{}
	cluster.ClusterName = utils.StringPtr("test-cluster")
	cluster.Archived = utils.BoolPtr(true)
	cluster.ClusterExternalIP = utils.StringPtr("10.0.0.0")

	type fields struct {
		client *client.Client
	}

	type args struct {
		KeyUUID string
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ImagedClusterDetails
		wantErr bool
	}{
		{
			"Get Imaged Cluster Details",
			fields{c},
			args{"0a8x-23d8"},
			cluster,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetImagedCluster(ctx, tt.args.KeyUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetImagedCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetImagedCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_GetAPIKey(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/api_keys/20ca-4d4c-61fd", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)

		fmt.Fprint(w, `{
			"api_key": "1243-7645", 
			"alias": "test-key",
			"created_timestamp": "2022-04-27T05:15:59.000-07:00",
			"current_time": "2022-04-27T09:33:25.000-07:00",
			"key_uuid": "20ca-4d4c-61fd"
		}`)
	})

	apiKey := &CreateAPIKeysResponse{}
	apiKey.APIKey = "1243-7645"
	apiKey.Alias = "test-key"
	apiKey.CreatedTimestamp = "2022-04-27T05:15:59.000-07:00"
	apiKey.CurrentTime = "2022-04-27T09:33:25.000-07:00"
	apiKey.KeyUUID = "20ca-4d4c-61fd"

	type fields struct {
		client *client.Client
	}

	type args struct {
		KeyUUID string
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateAPIKeysResponse
		wantErr bool
	}{
		{
			"Get API Key Details",
			fields{c},
			args{"20ca-4d4c-61fd"},
			apiKey,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.GetAPIKey(ctx, tt.args.KeyUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.GetAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.GetAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateAPIKey(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/api_keys", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		fmt.Fprint(w, `{
			"api_key": "1243-7645", 
			"alias": "test-key",
			"created_timestamp": "2022-04-27T05:15:59.000-07:00",
			"current_time": "2022-04-27T09:33:25.000-07:00",
			"key_uuid": "20ca-4d4c-61fd"
		}`)
	})

	apiKey := &CreateAPIKeysResponse{}
	apiKey.APIKey = "1243-7645"
	apiKey.Alias = "test-key"
	apiKey.CreatedTimestamp = "2022-04-27T05:15:59.000-07:00"
	apiKey.CurrentTime = "2022-04-27T09:33:25.000-07:00"
	apiKey.KeyUUID = "20ca-4d4c-61fd"

	input := &CreateAPIKeysInput{
		Alias: "test-key",
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		alias *CreateAPIKeysInput
	}

	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateAPIKeysResponse
		wantErr bool
	}{
		{
			"Create API Key",
			fields{c},
			args{input},
			apiKey,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateAPIKey(ctx, tt.args.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_CreateCluster(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/imaged_clusters", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		fmt.Fprint(w, `{
			"imaged_cluster_uuid": "123-654-678"
		}`)
	})

	clusterUUID := &CreateClusterResponse{
		ImagedClusterUUID: utils.StringPtr("123-654-678"),
	}

	input := &CreateClusterInput{
		CommonNetworkSettings: &CommonNetworkSettings{
			CvmDNSServers:        []string{"10.0.0.0"},
			HypervisorDNSServers: []string{"10.0.0.0"},
			CvmNtpServers:        []string{"0.0.0.0"},
			HypervisorNtpServers: []string{"0.0.0.0"},
		},
		RedundancyFactor: utils.IntPtr(2),
		AosPackageURL:    utils.StringPtr("test_aos.tar.gz"),
		ClusterName:      utils.StringPtr("test-cluster"),
		NodesList: []*Node{
			{
				CvmGateway:                 utils.StringPtr("0.0.0.0"),
				IpmiNetmask:                utils.StringPtr("255.255.255.0"),
				ImagedNodeUUID:             utils.StringPtr("12n0-vh87"),
				HypervisorType:             utils.StringPtr("kvm"),
				ImageNow:                   utils.BoolPtr(true),
				HypervisorHostname:         utils.StringPtr("HOST-1"),
				HypervisorNetmask:          utils.StringPtr("255.255.255.0"),
				HypervisorGateway:          utils.StringPtr("0.0.0.0"),
				CvmIP:                      utils.StringPtr("10.0.0.0"),
				CvmNetmask:                 utils.StringPtr("255.255.255.0"),
				IpmiIP:                     utils.StringPtr("10.0.0.0"),
				HypervisorIP:               utils.StringPtr("10.0.0.0"),
				IpmiGateway:                utils.StringPtr("0.0.0.0"),
				UseExistingNetworkSettings: utils.BoolPtr(false),
			},
		},
		HypervisorIsoDetails: &HypervisorIsoDetails{
			URL: utils.StringPtr("hypervisor.iso.tar.gz"),
		},
	}

	type fields struct {
		client *client.Client
	}

	type args struct {
		spec *CreateClusterInput
	}

	ctx := context.TODO()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateClusterResponse
		wantErr bool
	}{
		{
			"Imaged Nodes and create Cluster ",
			fields{c},
			args{input},
			clusterUUID,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			got, err := op.CreateCluster(ctx, tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operations.CreateCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operations.CreateCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperations_DeleteCluster(t *testing.T) {
	mux, c, server := setup()

	defer server.Close()

	mux.HandleFunc("/api/fc/v1/imaged_clusters/4e87-4a75-960f", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	type fields struct {
		client *client.Client
	}

	type args struct {
		UUID string
	}
	ctx := context.TODO()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test Delete Cluster OK",
			fields{c},
			args{"4e87-4a75-960f"},
			false,
		},

		{
			"Test Delete Cluster Errored",
			fields{c},
			args{},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			op := Operations{
				client: tt.fields.client,
			}
			if err := op.DeleteCluster(ctx, tt.args.UUID); (err != nil) != tt.wantErr {
				t.Errorf("Operations.DeleteCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
