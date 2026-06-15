package ndb

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestShouldRetryClusterCreateWithRawFallback(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil error", err: nil, want: false},
		{name: "clusterIP mismatch", err: errors.New("Unrecognized field 'clusterIP'"), want: true},
		{name: "clusterType mismatch", err: errors.New("Unrecognized field 'clusterType'"), want: true},
		{name: "clusterDescription mismatch", err: errors.New("Unrecognized field 'clusterDescription'"), want: true},
		{name: "other error", err: errors.New("timeout"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldRetryClusterCreateWithRawFallback(tt.err)
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestBuildClusterCreateRawFallbackRequest(t *testing.T) {
	baseSchema := ResourceNutanixNDBCluster().Schema

	t.Run("without prism central info", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, baseSchema, map[string]interface{}{
			"name":         "pe-a",
			"cluster_ip":   "10.1.1.10",
			"cluster_type": "NTNX",
			"version":      "v2",
			"username":     "admin",
			"password":     "secret",
		})

		got := buildClusterCreateRawFallbackRequest(d)
		if got["name"] != "pe-a" {
			t.Fatalf("unexpected name: %#v", got["name"])
		}
		if got["cloudType"] != "NTNX" {
			t.Fatalf("unexpected cloudType: %#v", got["cloudType"])
		}
		if _, ok := got["managementServerInfo"]; ok {
			t.Fatalf("managementServerInfo should be absent when prism_central_info is not set")
		}
	})

	t.Run("with prism central info", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, baseSchema, map[string]interface{}{
			"name":         "pe-a",
			"cluster_ip":   "10.1.1.10",
			"cluster_type": "NTNX",
			"version":      "v2",
			"username":     "admin",
			"password":     "secret",
			"prism_central_info": []interface{}{
				map[string]interface{}{
					"ip_address": "10.2.2.20",
					"port":       9440,
					"username":   "pc-admin",
					"password":   "pc-secret",
				},
			},
		})

		got := buildClusterCreateRawFallbackRequest(d)
		mgmt, ok := got["managementServerInfo"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected managementServerInfo map, got %#v", got["managementServerInfo"])
		}
		if mgmt["ipAddress"] != "10.2.2.20" {
			t.Fatalf("unexpected ipAddress: %#v", mgmt["ipAddress"])
		}
		if mgmt["username"] != "pc-admin" {
			t.Fatalf("unexpected username: %#v", mgmt["username"])
		}
	})
}
