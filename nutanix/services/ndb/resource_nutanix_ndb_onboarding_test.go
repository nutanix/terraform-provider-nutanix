package ndb

import (
	"context"
	"errors"
	"testing"

	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
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

func TestDiscoverOnboardingOptions(t *testing.T) {
	svc := &fakeOnboardingService{
		storageContainers: map[string]interface{}{
			"entities": []interface{}{
				map[string]interface{}{"name": "container-b"},
				map[string]interface{}{"containerName": "container-a"},
				map[string]interface{}{"vstore_name_list": []interface{}{"container-c", "container-a"}},
			},
		},
		serverConfig: &era.OnboardingEraServerConfig{
			DNSServers: []string{"10.1.1.2", "10.1.1.1", "10.1.1.1"},
			NTPServers: []string{"ntp-b", "ntp-a"},
		},
		networks: &era.ListNetworkResponse{
			&era.NetworkIntentResponse{Name: stringPtrForOnboardingTest("vlan-b")},
			&era.NetworkIntentResponse{Name: stringPtrForOnboardingTest("vlan-a")},
		},
	}

	got, err := discoverOnboardingOptions(context.Background(), svc, "cluster-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertStringSlice(t, got.StorageContainers, []string{"container-a", "container-b", "container-c"})
	assertStringSlice(t, got.DNSServers, []string{"10.1.1.1", "10.1.1.2"})
	assertStringSlice(t, got.NTPServers, []string{"ntp-a", "ntp-b"})
	assertStringSlice(t, got.NetworkNames, []string{"vlan-a", "vlan-b"})
}

func TestDiscoverOnboardingOptionsFailsOnNonRetryableStorageError(t *testing.T) {
	svc := &fakeOnboardingService{
		storageErr: errors.New("401 unauthorized"),
	}

	if _, err := discoverOnboardingOptions(context.Background(), svc, "cluster-1"); err == nil {
		t.Fatalf("expected non-retryable storage error")
	}
}

func TestApplyPrismCentralStepAcceptsValidatedDetails(t *testing.T) {
	svc := &fakeOnboardingService{
		validateDomainResp: map[string]interface{}{
			"name":        "pc-a",
			"external_ip": "10.1.1.10",
		},
	}

	if err := applyPrismCentralStep(context.Background(), svc, prismCentralInfoForTest()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.createClusterRawCalls != 0 {
		t.Fatalf("expected validation success to skip create fallback, got %d calls", svc.createClusterRawCalls)
	}
}

func TestApplyPrismCentralStepFailsOnValidationError(t *testing.T) {
	svc := &fakeOnboardingService{
		validateDomainErr: errors.New("Authentication Error. Please make sure the Prism Central username and password are correct."),
	}

	if err := applyPrismCentralStep(context.Background(), svc, prismCentralInfoForTest()); err == nil {
		t.Fatalf("expected prism central validation error")
	}
	if svc.createClusterRawCalls != 0 {
		t.Fatalf("expected validation error to stop before create fallback, got %d calls", svc.createClusterRawCalls)
	}
}

func TestApplyPrismCentralStepFailsOnValidationErrorPayload(t *testing.T) {
	svc := &fakeOnboardingService{
		validateDomainResp: map[string]interface{}{
			"errorCode": "ERA-PC-0001",
			"message":   "Prism Central and Prism Element versions are not compatible",
		},
	}

	if err := applyPrismCentralStep(context.Background(), svc, prismCentralInfoForTest()); err == nil {
		t.Fatalf("expected prism central validation payload error")
	}
	if svc.createClusterRawCalls != 0 {
		t.Fatalf("expected validation payload error to stop before create fallback, got %d calls", svc.createClusterRawCalls)
	}
}

func TestApplyPrismCentralStepFallsBackWhenValidationEndpointUnsupported(t *testing.T) {
	svc := &fakeOnboardingService{
		validateDomainErr: errors.New("error: Not Found"),
	}

	if err := applyPrismCentralStep(context.Background(), svc, prismCentralInfoForTest()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.createClusterRawCalls != 1 {
		t.Fatalf("expected create fallback call, got %d", svc.createClusterRawCalls)
	}
}

func TestApplyNDBConfigStepAppliesOnlyPresentConfig(t *testing.T) {
	svc := &fakeOnboardingService{}
	cfg := onboardingResolvedConfig{
		DNSServers:         []string{"10.1.1.1"},
		Timezone:           "UTC",
		ApplySMTPEvenEmpty: false,
	}

	if err := applyNDBConfigStep(context.Background(), svc, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(svc.validateConfigParams) != 1 || svc.validateConfigParams[0] != "dns" {
		t.Fatalf("unexpected validation params: %#v", svc.validateConfigParams)
	}
	if len(svc.setConfigParams) != 2 {
		t.Fatalf("expected dns and timezone set calls, got %#v", svc.setConfigParams)
	}
	assertStringSlice(t, svc.setConfigParams[0], []string{"dns"})
	assertStringSlice(t, svc.setConfigParams[1], []string{"timezone"})
}

func TestApplyStorageStepFallsBackToWizardUpload(t *testing.T) {
	svc := &fakeOnboardingService{
		cluster: &era.ListClusterResponse{
			Name:        stringPtrForOnboardingTest("pe-a"),
			Description: stringPtrForOnboardingTest("desc"),
			Cloudtype:   stringPtrForOnboardingTest("NTNX"),
			Version:     stringPtrForOnboardingTest("v2"),
			Ipaddresses: []*string{stringPtrForOnboardingTest("10.1.1.10")},
		},
		replaceErr: errors.New("put rejected"),
	}
	peInfo := []interface{}{
		map[string]interface{}{
			"username": "admin",
			"password": "secret",
		},
	}

	if err := applyStorageStep(context.Background(), svc, "cluster-1", "container-a", peInfo); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.replaceCalls != 1 {
		t.Fatalf("expected one replace call, got %d", svc.replaceCalls)
	}
	if svc.uploadCalls != 1 {
		t.Fatalf("expected one upload fallback call, got %d", svc.uploadCalls)
	}
	if got := svc.lastUploadBody["storageContainer"]; got != "container-a" {
		t.Fatalf("expected storage container fallback payload, got %#v", svc.lastUploadBody)
	}
}

func TestApplyNetworkStepBuildsSelectedNetworkPayload(t *testing.T) {
	svc := &fakeOnboardingService{
		cluster: &era.ListClusterResponse{
			Ipaddresses: []*string{stringPtrForOnboardingTest("10.1.1.10")},
		},
	}
	net := onboardingResolvedNetwork{
		Skip:        false,
		NetworkName: "vlan-a",
		StaticIP:    "10.1.1.20",
		Gateway:     "10.1.1.1",
		SubnetMask:  "255.255.255.0",
	}

	if err := applyNetworkStep(context.Background(), svc, "cluster-1", net); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.uploadCalls != 1 {
		t.Fatalf("expected one upload call, got %d", svc.uploadCalls)
	}
	if svc.createNetworkCalls != 1 {
		t.Fatalf("expected one network create call, got %d", svc.createNetworkCalls)
	}
	if svc.lastNetworkBody == nil ||
		*svc.lastNetworkBody.Name != "vlan-a" ||
		*svc.lastNetworkBody.ClusterID != "cluster-1" ||
		*svc.lastNetworkBody.Type != "Static" {
		t.Fatalf("unexpected network create payload: %#v", svc.lastNetworkBody)
	}
	if svc.lastUploadBody["ip_address"] != "10.1.1.20" ||
		svc.lastUploadBody["vlanName"] != "vlan-a" ||
		svc.lastUploadBody["gateway"] != "10.1.1.1" ||
		svc.lastUploadBody["subnetMask"] != "255.255.255.0" {
		t.Fatalf("unexpected network payload: %#v", svc.lastUploadBody)
	}
}

func TestApplyNetworkStepCreatesDHCPNetworkWithProperties(t *testing.T) {
	svc := &fakeOnboardingService{
		cluster: &era.ListClusterResponse{
			Ipaddresses: []*string{stringPtrForOnboardingTest("10.1.1.10")},
		},
	}
	net := onboardingResolvedNetwork{
		Skip:        false,
		NetworkName: "vlan-a",
	}

	if err := applyNetworkStep(context.Background(), svc, "cluster-1", net); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.createNetworkCalls != 1 {
		t.Fatalf("expected one network create call, got %d", svc.createNetworkCalls)
	}
	if svc.lastNetworkBody == nil ||
		*svc.lastNetworkBody.Type != "DHCP" ||
		len(svc.lastNetworkBody.Properties) != 1 ||
		*svc.lastNetworkBody.Properties[0].Name != "ADVANCED_NETWORKING" ||
		*svc.lastNetworkBody.Properties[0].Value != "FALSE" {
		t.Fatalf("unexpected DHCP network create payload: %#v", svc.lastNetworkBody)
	}
}

func TestApplySetupStepSkipsWhenTriggerFalse(t *testing.T) {
	svc := &fakeOnboardingService{}
	setupRaw := []interface{}{
		map[string]interface{}{
			"trigger":         false,
			"timeout_minutes": 1,
		},
	}

	opID, err := applySetupStep(context.Background(), svc, "cluster-1", setupRaw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opID != "" {
		t.Fatalf("expected empty operation id, got %q", opID)
	}
	if svc.uploadCalls != 0 {
		t.Fatalf("expected setup upload to be skipped, got %d calls", svc.uploadCalls)
	}
}

type fakeOnboardingService struct {
	era.Service

	cluster             *era.ListClusterResponse
	storageContainers   map[string]interface{}
	storageErr          error
	serverConfig        *era.OnboardingEraServerConfig
	networks            *era.ListNetworkResponse
	validateDomainResp  map[string]interface{}
	validateDomainErr   error
	createClusterRawErr error
	replaceErr          error
	uploadErr           error

	validateConfigParams  []string
	setConfigParams       [][]string
	createClusterRawCalls int
	replaceCalls          int
	uploadCalls           int
	createNetworkCalls    int
	lastUploadBody        map[string]interface{}
	lastNetworkBody       *era.NetworkIntentInput
}

func (f *fakeOnboardingService) GetCluster(_ context.Context, _, _ string) (*era.ListClusterResponse, error) {
	if f.cluster != nil {
		return f.cluster, nil
	}
	return &era.ListClusterResponse{}, nil
}

func (f *fakeOnboardingService) GetClusterStorageContainers(_ context.Context, _ string) (map[string]interface{}, error) {
	if f.storageErr != nil {
		return nil, f.storageErr
	}
	return f.storageContainers, nil
}

func (f *fakeOnboardingService) GetEraServerConfig(_ context.Context) (*era.OnboardingEraServerConfig, error) {
	if f.serverConfig != nil {
		return f.serverConfig, nil
	}
	return &era.OnboardingEraServerConfig{}, nil
}

func (f *fakeOnboardingService) ListNetwork(_ context.Context) (*era.ListNetworkResponse, error) {
	if f.networks != nil {
		return f.networks, nil
	}
	return &era.ListNetworkResponse{}, nil
}

func (f *fakeOnboardingService) CreateNetwork(_ context.Context, body *era.NetworkIntentInput) (*era.NetworkIntentResponse, error) {
	f.createNetworkCalls++
	f.lastNetworkBody = body
	id := "network-1"
	return &era.NetworkIntentResponse{ID: &id}, nil
}

func (f *fakeOnboardingService) ValidateDomainClusterDetails(_ context.Context, _ string, _ int, _ string, _ string) (map[string]interface{}, error) {
	if f.validateDomainErr != nil {
		return nil, f.validateDomainErr
	}
	if f.validateDomainResp != nil {
		return f.validateDomainResp, nil
	}
	return map[string]interface{}{"name": "pc-a", "external_ip": "10.1.1.10"}, nil
}

func (f *fakeOnboardingService) CreateClusterRaw(_ context.Context, _ map[string]interface{}) (*era.ProvisionDatabaseResponse, error) {
	f.createClusterRawCalls++
	if f.createClusterRawErr != nil {
		return nil, f.createClusterRawErr
	}
	return &era.ProvisionDatabaseResponse{}, nil
}

func (f *fakeOnboardingService) ValidateEraServerConfig(_ context.Context, _ *era.OnboardingEraServerConfig, configParams []string) (*era.OnboardingEraServerConfig, error) {
	f.validateConfigParams = append([]string{}, configParams...)
	return &era.OnboardingEraServerConfig{}, nil
}

func (f *fakeOnboardingService) SetEraServerConfig(_ context.Context, _ *era.OnboardingEraServerConfig, configParams []string) (*era.OnboardingEraServerConfig, error) {
	f.setConfigParams = append(f.setConfigParams, append([]string{}, configParams...))
	return &era.OnboardingEraServerConfig{}, nil
}

func (f *fakeOnboardingService) ReplaceClusterWizard(_ context.Context, _ string, _ map[string]interface{}) (string, error) {
	f.replaceCalls++
	if f.replaceErr != nil {
		return "", f.replaceErr
	}
	return "", nil
}

func (f *fakeOnboardingService) UploadClusterWizardJSON(_ context.Context, _ string, body map[string]interface{}, _, _, _ bool) (string, error) {
	f.uploadCalls++
	f.lastUploadBody = body
	if f.uploadErr != nil {
		return "", f.uploadErr
	}
	return "", nil
}

func stringPtrForOnboardingTest(v string) *string {
	return &v
}

func prismCentralInfoForTest() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":       "pc-a",
			"ip_address": "10.1.1.10",
			"port":       9440,
			"username":   "admin",
			"password":   "secret",
		},
	}
}

func assertStringSlice(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %#v, got %#v", want, got)
		}
	}
}
