package sdkconfig

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

// fakeAPIClient records the default headers a v4 SDK client would be configured with.
type fakeAPIClient struct {
	headers map[string]string
}

func newFakeAPIClient() *fakeAPIClient {
	return &fakeAPIClient{headers: map[string]string{}}
}

func (f *fakeAPIClient) AddDefaultHeader(name, value string) {
	f.headers[name] = value
}

func baseCreds() client.Credentials {
	return client.Credentials{
		Endpoint: "pc.example.com",
		Username: "admin",
		Password: "secret",
	}
}

// TestConfigureV4Client_VerifySSLFollowsInsecure is the regression test for the
// hardcoded VerifySSL:false. Every v4 SDK client used to skip certificate validation
// regardless of the provider's `insecure` setting, exposing credentials to MITM.
func TestConfigureV4Client_VerifySSLFollowsInsecure(t *testing.T) {
	for _, tc := range []struct {
		name          string
		insecure      bool
		wantVerifySSL bool
	}{
		{"secure by default", false, true},
		{"explicit opt-out", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			creds := baseCreds()
			creds.Insecure = tc.insecure

			cfg, err := ConfigureV4Client(creds, newFakeAPIClient())
			if err != nil {
				t.Fatalf("ConfigureV4Client returned error: %v", err)
			}
			if cfg == nil {
				t.Fatal("expected a config for complete credentials")
			}
			if cfg.VerifySSL != tc.wantVerifySSL {
				t.Errorf("VerifySSL = %v, want %v (insecure=%v)", cfg.VerifySSL, tc.wantVerifySSL, tc.insecure)
			}
		})
	}
}

func TestConfigureV4Client_Port(t *testing.T) {
	for _, tc := range []struct {
		name     string
		port     string
		wantPort int
		wantErr  bool
	}{
		{"default when unset", "", DefaultPort, false},
		{"honors configured port", "9441", 9441, false},
		{"rejects non-numeric", "94o1", 0, true},
		{"rejects out of range", "70000", 0, true},
		{"rejects zero", "0", 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			creds := baseCreds()
			creds.Port = tc.port

			cfg, err := ConfigureV4Client(creds, newFakeAPIClient())
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error for port %q, got nil", tc.port)
				}
				return
			}
			if err != nil {
				t.Fatalf("ConfigureV4Client returned error: %v", err)
			}
			if cfg.Port != tc.wantPort {
				t.Errorf("Port = %d, want %d", cfg.Port, tc.wantPort)
			}
		})
	}
}

func TestConfigureV4Client_InsufficientCredentials(t *testing.T) {
	for _, tc := range []struct {
		name  string
		creds client.Credentials
	}{
		{"no endpoint", client.Credentials{Username: "admin", Password: "secret"}},
		{"no auth", client.Credentials{Endpoint: "pc.example.com"}},
		{"username without password", client.Credentials{Endpoint: "pc.example.com", Username: "admin"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := ConfigureV4Client(tc.creds, newFakeAPIClient())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg != nil {
				t.Errorf("expected nil config for incomplete credentials, got %+v", cfg)
			}
		})
	}
}

func TestConfigureV4Client_APIKeyTakesPrecedence(t *testing.T) {
	creds := baseCreds()
	creds.APIKey = "test-key"

	api := newFakeAPIClient()
	cfg, err := ConfigureV4Client(creds, api)
	if err != nil {
		t.Fatalf("ConfigureV4Client returned error: %v", err)
	}
	if api.headers["X-Ntnx-Api-Key"] != "test-key" {
		t.Errorf("API key header = %q, want %q", api.headers["X-Ntnx-Api-Key"], "test-key")
	}
	if cfg.Username != "" || cfg.Password != "" {
		t.Error("basic auth must not be set when an API key is supplied")
	}
}

func TestConfigureV4Client_CustomHeaders(t *testing.T) {
	creds := baseCreds()
	creds.CustomHeaders = map[string]string{"CF-Access-Client-Id": "abc"}

	api := newFakeAPIClient()
	if _, err := ConfigureV4Client(creds, api); err != nil {
		t.Fatalf("ConfigureV4Client returned error: %v", err)
	}
	if api.headers["CF-Access-Client-Id"] != "abc" {
		t.Errorf("custom header not applied: %v", api.headers)
	}
}
