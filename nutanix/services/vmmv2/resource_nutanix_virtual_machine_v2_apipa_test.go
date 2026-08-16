package vmmv2

import (
	"testing"

	commonconfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestIsAPIPA(t *testing.T) {
	cases := map[string]bool{
		"169.254.0.1":    true,  // APIPA
		"169.254.99.200": true,  // APIPA
		"10.0.0.5":       false, // routable
		"192.168.1.10":   false,
		"172.30.130.229": false, // a real corp desktop
		"127.0.0.1":      false, // loopback, not link-local
		"0.0.0.0":        false, // unspecified
		"":               false,
		"not-an-ip":      false,
		"fe80::1":        false, // IPv6 link-local — we only skip IPv4 APIPA here
	}
	for in, want := range cases {
		if got := isAPIPA(in); got != want {
			t.Errorf("isAPIPA(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestGetFirstIPAddress(t *testing.T) {
	learned := func(ips ...string) config.Nic {
		la := make([]commonconfig.IPv4Address, 0, len(ips))
		for _, ip := range ips {
			la = append(la, commonconfig.IPv4Address{Value: utils.StringPtr(ip)})
		}
		return config.Nic{NetworkInfo: &config.NicNetworkInfo{Ipv4Info: &config.Ipv4Info{LearnedIpAddresses: la}}}
	}
	static := config.Nic{NetworkInfo: &config.NicNetworkInfo{
		Ipv4Config: &config.Ipv4Config{IpAddress: &commonconfig.IPv4Address{Value: utils.StringPtr("192.168.5.5")}},
	}}
	tests := []struct {
		name     string
		nic      config.Nic
		routable bool
		want     string
	}{
		{"routable: apipa then real -> real", learned("169.254.1.5", "10.0.0.5"), true, "10.0.0.5"},
		{"routable: only apipa -> empty (skipped)", learned("169.254.1.5"), true, ""},
		{"routable: real only -> real", learned("10.1.2.3"), true, "10.1.2.3"},
		{"not routable: only apipa -> apipa (accepted)", learned("169.254.1.5"), false, "169.254.1.5"},
		{"no network info -> empty", config.Nic{}, true, ""},
		{"static config fallback", static, true, "192.168.5.5"},
	}
	for _, tc := range tests {
		if got := getFirstIPAddress(tc.nic, tc.routable); got != tc.want {
			t.Errorf("%s: getFirstIPAddress(routable=%v) = %q, want %q", tc.name, tc.routable, got, tc.want)
		}
	}
}
