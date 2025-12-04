package foundationcentral

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

func TestNewFoundationCentralClient(t *testing.T) {
	// verifies positive client creation
	cred := client.Credentials{
		URL:                "foo.com",
		Username:           "username",
		Password:           "password",
		Port:               "",
		Endpoint:           "0.0.0.0",
		Insecure:           true,
		FoundationEndpoint: "10.0.0.0",
		FoundationPort:     "8000",
		RequiredFields:     nil,
	}
	_, err := NewFoundationCentralClient(cred)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	// verify missing client scenario
	cred2 := client.Credentials{
		URL:      "foo.com",
		Insecure: true,
		RequiredFields: map[string][]string{
			"prism_central":      {"username", "password", "endpoint"},
			"foundation_central": {"username", "password", "endpoint"},
		},
	}
	FcClient2, err2 := NewFoundationCentralClient(cred2)
	if err2 != nil {
		t.Errorf("%s", err2.Error())
	}

	if FcClient2.client.ErrorMsg == "" {
		t.Errorf("NewFoundationCentralClient(%v) expected the base client in v3 client to have some error message", cred2)
	}
}
