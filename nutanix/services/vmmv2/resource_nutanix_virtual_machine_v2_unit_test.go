// Unit tests for the routable-IP wait added for issue #871.
//
// These live in package vmmv2 (not vmmv2_test, like the acceptance tests in
// resource_nutanix_virtual_machine_v2_test.go) because they exercise the unexported
// isAPIPA/getFirstIPAddress helpers directly, which an external test package cannot reach.
// The acceptance coverage for the same feature is
// TestAccV2NutanixVmsResource_WaitForRoutableIP in resource_nutanix_virtual_machine_v2_test.go.
package vmmv2

import (
	"context"
	"testing"
	"time"

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

// TestFirstIPAddress checks that the IP wait and the connection info pick the first NIC
// with a usable address, so a VM whose first NIC has only APIPA still gets connection info
// from a later NIC.
func TestFirstIPAddress(t *testing.T) {
	learned := func(ips ...string) config.Nic {
		la := make([]commonconfig.IPv4Address, 0, len(ips))
		for _, ip := range ips {
			la = append(la, commonconfig.IPv4Address{Value: utils.StringPtr(ip)})
		}
		return config.Nic{NetworkInfo: &config.NicNetworkInfo{Ipv4Info: &config.Ipv4Info{LearnedIpAddresses: la}}}
	}
	tests := []struct {
		name     string
		nics     []config.Nic
		routable bool
		want     string
	}{
		{"first NIC routable", []config.Nic{learned("10.0.0.5"), learned("10.0.1.5")}, true, "10.0.0.5"},
		{"first NIC only apipa, second routable", []config.Nic{learned("169.254.1.5"), learned("10.0.1.5")}, true, "10.0.1.5"},
		{"first NIC only apipa, not routable", []config.Nic{learned("169.254.1.5"), learned("10.0.1.5")}, false, "169.254.1.5"},
		{"only apipa anywhere", []config.Nic{learned("169.254.1.5"), learned("169.254.2.5")}, true, ""},
		{"no NICs", nil, true, ""},
	}
	for _, tc := range tests {
		if got := firstIPAddress(tc.nics, tc.routable); got != tc.want {
			t.Errorf("%s: firstIPAddress(routable=%v) = %q, want %q", tc.name, tc.routable, got, tc.want)
		}
	}
}

// TestWaitForIPDuration checks the IP wait is capped to the time left on the create
// deadline, so a large wait_for_ip_timeout cannot make create fail after the VM exists.
func TestWaitForIPDuration(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	withDeadline := func(left time.Duration) context.Context {
		ctx, cancel := context.WithDeadline(context.Background(), now.Add(left))
		t.Cleanup(cancel)
		return ctx
	}
	tests := []struct {
		name    string
		ctx     context.Context
		minutes int
		want    time.Duration
	}{
		{"disabled", context.Background(), 0, 0},
		{"no deadline", context.Background(), 30, 30 * time.Minute},
		{"fits the deadline", withDeadline(20 * time.Minute), 5, 5 * time.Minute},
		{"capped to the deadline", withDeadline(20 * time.Minute), 30, 20*time.Minute - waitForIPReadMargin},
		{"no time left", withDeadline(time.Minute), 5, 0},
	}
	for _, tc := range tests {
		if got := waitForIPDuration(tc.ctx, tc.minutes, now); got != tc.want {
			t.Errorf("%s: waitForIPDuration(%d) = %s, want %s", tc.name, tc.minutes, got, tc.want)
		}
	}
}

// TestWaitForIPSchemaDefaults pins the schema defaults of the two new arguments to the
// package constants, so a future change to either has to be deliberate.
func TestWaitForIPSchemaDefaults(t *testing.T) {
	s := ResourceNutanixVirtualMachineV2().Schema

	if got := s["wait_for_ip_timeout"].Default; got != defaultWaitForIPTimeoutMinutes {
		t.Errorf("wait_for_ip_timeout default = %v, want %v", got, defaultWaitForIPTimeoutMinutes)
	}
	if got := s["wait_for_ip_routable"].Default; got != defaultWaitForIPRoutable {
		t.Errorf("wait_for_ip_routable default = %v, want %v", got, defaultWaitForIPRoutable)
	}
}
