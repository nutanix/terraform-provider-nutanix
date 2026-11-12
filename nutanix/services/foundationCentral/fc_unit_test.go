package fc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	fc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFoundationCentralSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixFCAPIKeys().InternalValidate(ResourceNutanixFCAPIKeys().Schema, true))
	require.NoError(t, ResourceNutanixFCOnboardNodes().InternalValidate(ResourceNutanixFCOnboardNodes().Schema, true))
	require.NoError(t, ResourceNutanixFCImageCluster().InternalValidate(ResourceNutanixFCImageCluster().Schema, true))
	require.NoError(t, DataSourceNutanixFCAPIKeys().InternalValidate(DataSourceNutanixFCAPIKeys().Schema, false))
	require.NoError(t, DataSourceNutanixFCListAPIKeys().InternalValidate(DataSourceNutanixFCListAPIKeys().Schema, false))
	require.NoError(t, DataSourceNutanixFCImagedNodesList().InternalValidate(DataSourceNutanixFCImagedNodesList().Schema, false))
	require.NoError(t, DataSourceFCImagedNodeDetails().InternalValidate(DataSourceFCImagedNodeDetails().Schema, false))
	require.NoError(t, DataSourceNutanixFCClusterDetails().InternalValidate(DataSourceNutanixFCClusterDetails().Schema, false))
	require.NoError(t, DataSourceNutanixFCImagedClustersList().InternalValidate(DataSourceNutanixFCImagedClustersList().Schema, false))
}

func TestUnitFlattenAPIKeysList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*fc.CreateAPIKeysResponse
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*fc.CreateAPIKeysResponse{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []*fc.CreateAPIKeysResponse{{
				Alias:            "alias-1",
				APIKey:           "key-1",
				CreatedTimestamp: "ts-1",
				CurrentTime:      "now",
				KeyUUID:          "uuid-1",
			}},
			want: []map[string]interface{}{{
				"alias":             "alias-1",
				"api_key":           "key-1",
				"created_timestamp": "ts-1",
				"current_time":      "now",
				"key_uuid":          "uuid-1",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAPIKeysList(tt.in))
		})
	}
}

func TestUnitFlattenFeature(t *testing.T) {
	t.Parallel()

	feat := "feature-a"
	tests := []struct {
		name string
		in   []*string
		want []string
	}{
		{name: "nil", in: nil, want: []string{}},
		{name: "empty", in: []*string{}, want: []string{}},
		{name: "populated", in: []*string{&feat}, want: []string{"feature-a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFeature(tt.in))
		})
	}
}

func TestUnitFlattenImagedNodes(t *testing.T) {
	t.Parallel()

	aos := "6.5"
	feat := "f1"
	tests := []struct {
		name    string
		in      []*fc.ImagedNodeDetails
		wantLen int
		check   func(t *testing.T, got []map[string]interface{})
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []*fc.ImagedNodeDetails{}, wantLen: 0},
		{
			name:    "populated",
			in:      []*fc.ImagedNodeDetails{{AosVersion: &aos, SupportedFeatures: []*string{&feat}}},
			wantLen: 1,
			check: func(t *testing.T, got []map[string]interface{}) {
				require.Equal(t, &aos, got[0]["aos_version"])
				require.Equal(t, []string{"f1"}, got[0]["supported_features"])
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenImagedNodes(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestUnitFlattenFCCommonNetworkSettings(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenFCCommonNetworkSettings(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenFCCommonNetworkSettings(&fc.CommonNetworkSettings{
			CvmDNSServers:        []string{"8.8.8.8"},
			HypervisorDNSServers: []string{"1.1.1.1"},
			CvmNtpServers:        []string{"ntp1"},
			HypervisorNtpServers: []string{"ntp2"},
		})
		require.Len(t, got, 1)
		m := got[0].(map[string]interface{})
		require.Equal(t, []string{"8.8.8.8"}, utils.StringValueSlice(m["cvm_dns_servers"].([]*string)))
		require.Equal(t, []string{"1.1.1.1"}, utils.StringValueSlice(m["hypervisor_dns_servers"].([]*string)))
		require.Equal(t, []string{"ntp1"}, utils.StringValueSlice(m["cvm_ntp_servers"].([]*string)))
		require.Equal(t, []string{"ntp2"}, utils.StringValueSlice(m["hypervisor_ntp_servers"].([]*string)))
	})
}

func TestUnitFlattenClusterStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *fc.ClusterStatus
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &fc.ClusterStatus{
				IntentPickedUp:           utils.BoolPtr(true),
				ClusterCreationStarted:   utils.BoolPtr(false),
				ImagingStopped:           utils.BoolPtr(true),
				AggregatePercentComplete: utils.Float64Ptr(50),
				CurrentFoundationIP:      utils.StringPtr("10.0.0.2"),
				FoundationSessionID:      utils.StringPtr("sess-1"),
			},
			want: []interface{}{map[string]interface{}{
				"intent_picked_up":           true,
				"cluster_creation_started":   false,
				"imaging_stopped":            true,
				"aggregate_percent_complete": float64(50),
				"current_foundation_ip":      "10.0.0.2",
				"foundation_session_id":      "sess-1",
				"node_progress_details":      []map[string]interface{}{},
				"cluster_progress_details":   []interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenClusterStatus(tt.in))
		})
	}
}

func TestUnitFlattenNodeProgressDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*fc.NodeProgressDetail
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*fc.NodeProgressDetail{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []*fc.NodeProgressDetail{{
				Status:          utils.StringPtr("FAILED"),
				ImagedNodeUUID:  utils.StringPtr("node-1"),
				ImagingStopped:  utils.BoolPtr(true),
				IntentPickedUp:  utils.BoolPtr(true),
				PercentComplete: utils.Float64Ptr(10),
				MessageList:     []*string{utils.StringPtr("err")},
			}},
			want: []map[string]interface{}{{
				"status":           utils.StringPtr("FAILED"),
				"imaged_node_uuid": utils.StringPtr("node-1"),
				"imaging_stopped":  utils.BoolPtr(true),
				"intent_picked_up": utils.BoolPtr(true),
				"percent_complete": utils.Float64Ptr(10),
				"message_list":     []string{"err"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNodeProgressDetails(tt.in))
		})
	}
}

func TestUnitFlattenClusterProgressDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *fc.ClusterProgressDetails
		want []interface{}
	}{
		{name: "nil", in: nil, want: []interface{}{}},
		{
			name: "populated",
			in: &fc.ClusterProgressDetails{
				ClusterName:     utils.StringPtr("cluster-a"),
				Status:          utils.StringPtr("RUNNING"),
				PercentComplete: utils.Float64Ptr(75),
				MessageList:     []*string{utils.StringPtr("ok")},
			},
			want: []interface{}{map[string]interface{}{
				"cluster_name":     "cluster-a",
				"status":           "RUNNING",
				"percent_complete": float64(75),
				"message_list":     []string{"ok"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenClusterProgressDetails(tt.in))
		})
	}
}

func TestUnitFlattenFCFoundationInitConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      *fc.FoundationInitConfig
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{
			name: "populated",
			in: &fc.FoundationInitConfig{
				CvmGateway: "10.0.0.1",
				CvmNetmask: "255.255.255.0",
				// flattenHypervisorIsoURL panics on nil; always set when testing populated config.
				HypervisorIsoURL: &fc.HypervisorIso{
					HypervisorType: utils.StringPtr("AHV"),
					URL:            utils.StringPtr("https://example.test/iso"),
				},
			},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenFCFoundationInitConfig(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				m := got[0].(map[string]interface{})
				require.Equal(t, "10.0.0.1", m["cvm_gateway"])
				require.Equal(t, "255.255.255.0", m["cvm_netmask"])
			}
		})
	}
}

func TestUnitFlattenFCBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []*fc.Blocks
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []*fc.Blocks{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   []*fc.Blocks{{BlockID: "block-1"}},
			want: []map[string]interface{}{{
				"block_id": "block-1",
				"nodes":    []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFCBlock(tt.in))
		})
	}
}

func TestUnitFlattenNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []*fc.Nodes
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []*fc.Nodes{}, wantLen: 0},
		{name: "populated", in: []*fc.Nodes{{CvmIP: "10.0.0.10", ImageNow: true}}, wantLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenNodes(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				require.Equal(t, "10.0.0.10", got[0]["cvm_ip"])
				require.Equal(t, true, got[0]["image_now"])
			}
		})
	}
}

func TestUnitFlattenCluster(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []*fc.Clusters
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []*fc.Clusters{}, wantLen: 0},
		{name: "populated", in: []*fc.Clusters{{ClusterName: "c1", RedundancyFactor: 2}}, wantLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenCluster(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				require.Equal(t, "c1", got[0]["cluster_name"])
				require.Equal(t, 2, got[0]["redundancy_factor"])
			}
		})
	}
}

func TestUnitFlattenHypervisorIsoURL(t *testing.T) {
	t.Parallel()

	ht := "AHV"
	sum := "abc"
	url := "https://example.test/iso"
	got := flattenHypervisorIsoURL(&fc.HypervisorIso{
		HypervisorType: &ht,
		Sha256sum:      &sum,
		URL:            &url,
	})
	require.Equal(t, map[string]interface{}{
		"hypervisor_type": &ht,
		"sha256sum":       &sum,
		"url":             &url,
	}, got)
}

func TestUnitFlattenFCHypervisorIsos(t *testing.T) {
	t.Parallel()

	ht := "ESX"
	tests := []struct {
		name string
		in   []*fc.HypervisorIso
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []*fc.HypervisorIso{}, want: nil},
		{
			name: "populated",
			in:   []*fc.HypervisorIso{{HypervisorType: &ht}},
			want: []map[string]interface{}{{
				"hypervisor_type": &ht,
				"sha256sum":       (*string)(nil),
				"url":             (*string)(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFCHypervisorIsos(tt.in))
		})
	}
}

func TestUnitFlattenNosPackage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *fc.NosPackageURL
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &fc.NosPackageURL{Sha256sum: "sum", URL: "https://example.test/nos"},
			want: []interface{}{map[string]interface{}{
				"sha256sum": "sum",
				"url":       "https://example.test/nos",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenNosPackage(tt.in))
		})
	}
}

func TestUnitFlattenHardwareAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   map[string]interface{}
		want map[string]interface{}
	}{
		{name: "nil", in: nil, want: map[string]interface{}{}},
		{name: "empty", in: map[string]interface{}{}, want: map[string]interface{}{}},
		{name: "populated", in: map[string]interface{}{"cpu": 8}, want: map[string]interface{}{"cpu": "8"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenHardwareAttributes(tt.in))
		})
	}
}

func TestUnitFlattenImagedClusters(t *testing.T) {
	t.Parallel()

	name := "cluster-x"
	tests := []struct {
		name    string
		in      []*fc.ImagedClusterDetails
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []*fc.ImagedClusterDetails{}, wantLen: 0},
		{name: "populated", in: []*fc.ImagedClusterDetails{{ClusterName: &name}}, wantLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenImagedClusters(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				require.Equal(t, &name, got[0]["cluster_name"])
			}
		})
	}
}

func TestUnitExpandNetworklist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []string
	}{
		{name: "empty", in: []interface{}{}, want: []string{}},
		{name: "populated", in: []interface{}{"a", "b"}, want: []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandNetworklist(tt.in))
		})
	}
}

func TestUnitExpandCommonNetworkSettings(t *testing.T) {
	t.Parallel()

	t.Run("missing", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixFCImageCluster().Schema, map[string]interface{}{})
		require.Nil(t, expandCommonNetworkSettings(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixFCImageCluster().Schema, map[string]interface{}{
			"common_network_settings": []interface{}{map[string]interface{}{
				"cvm_dns_servers":        []interface{}{"8.8.8.8"},
				"cvm_ntp_servers":        []interface{}{"ntp1"},
				"hypervisor_dns_servers": []interface{}{"1.1.1.1"},
				"hypervisor_ntp_servers": []interface{}{"ntp2"},
			}},
		})
		got := expandCommonNetworkSettings(d)
		require.NotNil(t, got)
		require.Equal(t, []string{"8.8.8.8"}, got.CvmDNSServers)
		require.Equal(t, []string{"ntp1"}, got.CvmNtpServers)
		require.Equal(t, []string{"1.1.1.1"}, got.HypervisorDNSServers)
		require.Equal(t, []string{"ntp2"}, got.HypervisorNtpServers)
	})
}

func TestUnitExpandHyperVisorIsoDetails(t *testing.T) {
	t.Parallel()

	t.Run("nil slice", func(t *testing.T) {
		t.Parallel()
		var in []interface{}
		require.Nil(t, expandHyperVisorIsoDetails(in))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandHyperVisorIsoDetails([]interface{}{map[string]interface{}{
			"hyperv_sku":         "sku",
			"url":                "https://example.test/iso",
			"hyperv_product_key": "key",
			"sha256sum":          "sum",
			"hypervisor_type":    "AHV",
		}})
		require.NotNil(t, got)
		require.Len(t, *got, 1)
		require.Equal(t, utils.StringPtr("AHV"), (*got)[0].HypervisorType)
		require.Equal(t, utils.StringPtr("https://example.test/iso"), (*got)[0].URL)
	})
}

func TestUnitExpandNodesList(t *testing.T) {
	t.Parallel()

	t.Run("missing", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixFCImageCluster().Schema, map[string]interface{}{})
		require.Nil(t, expandNodesList(d))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixFCImageCluster().Schema, map[string]interface{}{
			"node_list": []interface{}{map[string]interface{}{
				"cvm_ip":           "10.0.0.10",
				"imaged_node_uuid": "node-uuid",
				"hypervisor_type":  "AHV",
				"image_now":        true,
				"cvm_vlan_id":      10,
			}},
		})
		got := expandNodesList(d)
		require.Len(t, got, 1)
		require.Equal(t, utils.StringPtr("10.0.0.10"), got[0].CvmIP)
		require.Equal(t, utils.StringPtr("node-uuid"), got[0].ImagedNodeUUID)
		require.Equal(t, utils.StringPtr("AHV"), got[0].HypervisorType)
		require.Equal(t, utils.BoolPtr(true), got[0].ImageNow)
		require.Equal(t, utils.IntPtr(10), got[0].CvmVlanID)
	})
}

func TestUnitCollectIndividualErrorDiagnosticsFC(t *testing.T) {
	t.Parallel()

	t.Run("no failures", func(t *testing.T) {
		t.Parallel()
		progress := &fc.ImagedClusterDetails{
			ClusterStatus: &fc.ClusterStatus{
				NodeProgressDetails: []*fc.NodeProgressDetail{{
					PercentComplete: utils.Float64Ptr(100),
				}},
				ClusterProgressDetails: &fc.ClusterProgressDetails{
					PercentComplete: utils.Float64Ptr(100),
				},
			},
		}
		require.Empty(t, collectIndividualErrorDiagnosticsFC(progress))
	})

	t.Run("node failure", func(t *testing.T) {
		t.Parallel()
		progress := &fc.ImagedClusterDetails{
			ClusterStatus: &fc.ClusterStatus{
				NodeProgressDetails: []*fc.NodeProgressDetail{{
					PercentComplete: utils.Float64Ptr(10),
					ImagedNodeUUID:  utils.StringPtr("node-1"),
					Status:          utils.StringPtr("FAILED"),
					MessageList:     []*string{utils.StringPtr("boom")},
				}},
			},
		}
		diags := collectIndividualErrorDiagnosticsFC(progress)
		require.Len(t, diags, 1)
		require.Equal(t, diag.Error, diags[0].Severity)
		require.Contains(t, diags[0].Summary, "node-1")
		require.Equal(t, "boom", diags[0].Detail)
	})
}

func TestUnitImageClusterUpdateNoOp(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixFCImageCluster().Schema, map[string]interface{}{})
	require.Nil(t, resourceNutanixFCImageClusterUpdate(t.Context(), d, nil))
}
