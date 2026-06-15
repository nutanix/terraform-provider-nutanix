package ndb

import (
	"errors"
	"testing"
)

func TestExpandStringList(t *testing.T) {
	in := []interface{}{"a", "", "b"}
	out := expandStringList(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 values, got %d", len(out))
	}
	if out[0] != "a" || out[1] != "b" {
		t.Fatalf("unexpected output: %#v", out)
	}
}

func TestResolveStorageContainer(t *testing.T) {
	tests := []struct {
		name          string
		selectionMode string
		raw           interface{}
		available     []string
		want          string
		wantErr       bool
	}{
		{
			name:          "uses user value when present",
			selectionMode: "auto",
			raw: []interface{}{
				map[string]interface{}{"container_name": "user-container"},
			},
			available: []string{"a", "b"},
			want:      "user-container",
		},
		{
			name:          "strict mode rejects unknown user value",
			selectionMode: "strict",
			raw: []interface{}{
				map[string]interface{}{"container_name": "missing"},
			},
			available: []string{"a", "b"},
			wantErr:   true,
		},
		{
			name:          "auto mode picks first available",
			selectionMode: "auto",
			raw:           []interface{}{},
			available:     []string{"first", "second"},
			want:          "first",
		},
		{
			name:          "strict mode fails when no value and no discovered",
			selectionMode: "strict",
			raw:           []interface{}{},
			available:     []string{},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveStorageContainer(tt.selectionMode, tt.raw, tt.available)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestResolveNetworkDetails(t *testing.T) {
	discovery := &onboardingDiscovery{NetworkNames: []string{"discovered-vlan"}}

	tests := []struct {
		name string
		raw  interface{}
		want onboardingResolvedNetwork
	}{
		{
			name: "defaults to skip true when missing block",
			raw:  nil,
			want: onboardingResolvedNetwork{Skip: true},
		},
		{
			name: "uses existing network name",
			raw: []interface{}{
				map[string]interface{}{
					"skip":                  false,
					"existing_network_name": "custom-vlan",
				},
			},
			want: onboardingResolvedNetwork{
				Skip:        false,
				NetworkName: "custom-vlan",
			},
		},
		{
			name: "fallback does not force network when skip true",
			raw: []interface{}{
				map[string]interface{}{
					"skip": true,
				},
			},
			want: onboardingResolvedNetwork{
				Skip: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveNetworkDetails(tt.raw, discovery)
			if got.Skip != tt.want.Skip || got.NetworkName != tt.want.NetworkName {
				t.Fatalf("unexpected network resolution: %#v", got)
			}
		})
	}
}

func TestResolveNDBConfig(t *testing.T) {
	tests := []struct {
		name      string
		raw       interface{}
		discovery *onboardingDiscovery
		wantDNS   []string
		wantNTP   []string
		wantTZ    string
	}{
		{
			name:      "uses discovered values by default",
			raw:       nil,
			discovery: &onboardingDiscovery{DNSServers: []string{"10.1.1.1"}, NTPServers: []string{"10.2.2.2"}},
			wantDNS:   []string{"10.1.1.1"},
			wantNTP:   []string{"10.2.2.2"},
			wantTZ:    "UTC",
		},
		{
			name:      "uses explicit config overrides",
			discovery: &onboardingDiscovery{DNSServers: []string{"10.1.1.1"}, NTPServers: []string{"10.2.2.2"}},
			raw: []interface{}{
				map[string]interface{}{
					"dns_servers": []interface{}{"8.8.8.8"},
					"ntp_servers": []interface{}{"1.1.1.1"},
					"timezone":    "Asia/Kolkata",
				},
			},
			wantDNS: []string{"8.8.8.8"},
			wantNTP: []string{"1.1.1.1"},
			wantTZ:  "Asia/Kolkata",
		},
		{
			name:      "keeps empty when nothing discovered and no hardcoded fallback",
			raw:       nil,
			discovery: &onboardingDiscovery{},
			wantDNS:   []string{},
			wantNTP:   []string{},
			wantTZ:    "UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveNDBConfig(tt.raw, tt.discovery)
			if got.Timezone != tt.wantTZ {
				t.Fatalf("expected timezone %q, got %q", tt.wantTZ, got.Timezone)
			}
			if len(got.DNSServers) != len(tt.wantDNS) {
				t.Fatalf("expected dns %#v, got %#v", tt.wantDNS, got.DNSServers)
			}
			for i := range tt.wantDNS {
				if got.DNSServers[i] != tt.wantDNS[i] {
					t.Fatalf("expected dns %#v, got %#v", tt.wantDNS, got.DNSServers)
				}
			}
			if len(got.NTPServers) != len(tt.wantNTP) {
				t.Fatalf("expected ntp %#v, got %#v", tt.wantNTP, got.NTPServers)
			}
			for i := range tt.wantNTP {
				if got.NTPServers[i] != tt.wantNTP[i] {
					t.Fatalf("expected ntp %#v, got %#v", tt.wantNTP, got.NTPServers)
				}
			}
		})
	}
}

func TestIsRetryableStep4Error(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil is not retryable", err: nil, want: false},
		{name: "prism rest caller is retryable", err: errors.New("Could not get prism rest caller for cloudId"), want: true},
		{name: "era sql internal is retryable", err: errors.New("ERA-SQL-0000001"), want: true},
		{name: "generic internal error is not retryable", err: errors.New("An internal error has occurred"), want: false},
		{name: "generic era code is not retryable", err: errors.New("ERA-0000000"), want: false},
		{name: "unrelated error is not retryable", err: errors.New("401 unauthorized"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableStep4Error(tt.err)
			if got != tt.want {
				t.Fatalf("expected %v, got %v for err %v", tt.want, got, tt.err)
			}
		})
	}
}
