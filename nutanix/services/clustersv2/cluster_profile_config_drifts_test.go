package clustersv2

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestClusterProfileConfigDrifts(t *testing.T) {
	tests := []struct {
		name   string
		drifts []config.ConfigType
		want   []interface{}
	}{
		{name: "omitted", want: []interface{}{}},
		{name: "empty", drifts: []config.ConfigType{}, want: []interface{}{}},
		{
			name:   "single drift",
			drifts: []config.ConfigType{config.CONFIGTYPE_NTP_SERVER_CONFIG},
			want:   []interface{}{"NTP_SERVER_CONFIG"},
		},
		{
			name:   "multiple drifts",
			drifts: []config.ConfigType{config.CONFIGTYPE_NTP_SERVER_CONFIG, config.CONFIGTYPE_SNMP_SERVER_CONFIG},
			want:   []interface{}{"NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"},
		},
	}

	for name, resource := range map[string]*schema.Resource{
		"resource":    ResourceNutanixClusterProfileV2(),
		"data source": DatasourceNutanixClusterProfileV2(),
	} {
		t.Run(name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					d := schema.TestResourceDataRaw(t, resource.Schema, nil)
					clusters := []config.ManagedCluster{{
						ExtId:        utils.StringPtr("cluster-id"),
						ConfigDrifts: tt.drifts,
					}}
					if err := d.Set("clusters", flattenClusterProfileClusterList(clusters)); err != nil {
						t.Fatalf("setting clusters: %v", err)
					}
					if got := d.Get("clusters.0.config_drifts"); !reflect.DeepEqual(got, tt.want) {
						t.Fatalf("config_drifts = %#v, want %#v", got, tt.want)
					}
					if got := d.Get("clusters.0.ext_id"); got != "cluster-id" {
						t.Fatalf("cluster ext_id = %v, want cluster-id", got)
					}
				})
			}
		})
	}
}

func TestClusterProfileConfigDriftsStateUpgrade(t *testing.T) {
	tests := []struct {
		name    string
		drifts  interface{}
		want    []string
		flatmap bool
	}{
		{name: "empty string", drifts: "", want: []string{}},
		{name: "nonempty string", drifts: "NTP_SERVER_CONFIG", want: []string{"NTP_SERVER_CONFIG"}},
		{name: "null", want: []string{}},
		{name: "flatmap empty string", drifts: "", want: []string{}, flatmap: true},
		{name: "flatmap nonempty string", drifts: "NTP_SERVER_CONFIG", want: []string{"NTP_SERVER_CONFIG"}, flatmap: true},
	}
	resource := ResourceNutanixClusterProfileV2()
	server := schema.NewGRPCProviderServer(&schema.Provider{
		ResourcesMap: map[string]*schema.Resource{"nutanix_cluster_profile_v2": resource},
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawState, err := json.Marshal(map[string]interface{}{
				"id":   "profile-id",
				"name": "profile",
				"clusters": []interface{}{map[string]interface{}{
					"ext_id":        "cluster-id",
					"config_drifts": tt.drifts,
				}},
			})
			if err != nil {
				t.Fatal(err)
			}
			state := &tfprotov5.RawState{JSON: rawState}
			if tt.flatmap {
				state = &tfprotov5.RawState{Flatmap: map[string]string{
					"id":                       "profile-id",
					"name":                     "profile",
					"clusters.#":               "1",
					"clusters.0.ext_id":        "cluster-id",
					"clusters.0.config_drifts": tt.drifts.(string),
				}}
			}
			resp, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{
				TypeName: "nutanix_cluster_profile_v2",
				Version:  0,
				RawState: state,
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, diagnostic := range resp.Diagnostics {
				if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
					t.Fatalf("upgrading state: %s: %s", diagnostic.Summary, diagnostic.Detail)
				}
			}
			if resp.UpgradedState == nil {
				t.Fatal("missing upgraded state")
			}
			upgradedState, err := msgpack.Unmarshal(resp.UpgradedState.MsgPack, resource.CoreConfigSchema().ImpliedType())
			if err != nil {
				t.Fatal(err)
			}
			clusters := upgradedState.GetAttr("clusters").AsValueSlice()
			if len(clusters) != 1 || clusters[0].GetAttr("ext_id").AsString() != "cluster-id" {
				t.Fatalf("cluster changed during upgrade: %#v", clusters)
			}
			drifts := clusters[0].GetAttr("config_drifts")
			got := []string{}
			if !drifts.IsNull() {
				for _, drift := range drifts.AsValueSlice() {
					got = append(got, drift.AsString())
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("upgraded config_drifts = %#v, want %#v", got, tt.want)
			}
			if upgradedState.GetAttr("name").AsString() != "profile" || upgradedState.GetAttr("id").AsString() != "profile-id" {
				t.Fatal("profile identity changed during upgrade")
			}
		})
	}
}

func TestClusterProfileConfigDriftsRefresh(t *testing.T) {
	for name, resource := range map[string]*schema.Resource{
		"resource":    ResourceNutanixClusterProfileV2(),
		"data source": DatasourceNutanixClusterProfileV2(),
	} {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resource.Schema, nil)
			for _, drifts := range [][]config.ConfigType{
				{config.CONFIGTYPE_NTP_SERVER_CONFIG},
				{},
				{config.CONFIGTYPE_SNMP_SERVER_CONFIG},
				nil,
			} {
				if err := d.Set("clusters", flattenClusterProfileClusterList([]config.ManagedCluster{{ConfigDrifts: drifts}})); err != nil {
					t.Fatalf("setting clusters: %v", err)
				}
				got := d.Get("clusters.0.config_drifts").([]interface{})
				if len(got) != len(drifts) {
					t.Fatalf("config_drifts = %#v, want %d entries", got, len(drifts))
				}
				for i, drift := range drifts {
					if got[i] != drift.GetName() {
						t.Fatalf("config_drifts[%d] = %v, want %s", i, got[i], drift.GetName())
					}
				}
			}
			if err := d.Set("clusters", flattenClusterProfileClusterList(nil)); err != nil {
				t.Fatalf("clearing clusters: %v", err)
			}
			if got := d.Get("clusters.#"); got != 0 {
				t.Fatalf("cluster count = %v after clearing clusters, want 0", got)
			}
		})
	}
}
