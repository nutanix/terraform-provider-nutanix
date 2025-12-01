package foundation

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

func TestNewFoundationAPIClient(t *testing.T) {
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
	foundationClient, err := NewFoundationAPIClient(cred)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	outURL := fmt.Sprintf("http://%s:%s/", cred.FoundationEndpoint, cred.FoundationPort)
	if foundationClient.client.BaseURL.String() != outURL {
		t.Errorf("NewFoundationAPIClient(%v) BaseUrl in base client of foundation client = %v, expected %v", cred, foundationClient.client.BaseURL.String(), outURL)
	}

	// verify missing client scenario
	cred2 := client.Credentials{
		URL:      "foo.com",
		Username: "username",
		Password: "password",
		Port:     "",
		Endpoint: "0.0.0.0",
		Insecure: true,
		RequiredFields: map[string][]string{
			"foundation": {"foundation_endpoint"},
		},
	}
	foundationClient2, err2 := NewFoundationAPIClient(cred2)
	if err2 != nil {
		t.Errorf("%s", err2.Error())
	}

	if foundationClient2.client.ErrorMsg == "" {
		t.Errorf("NewFoundationAPIClient(%v) expected the base client in foundation client to have some error message", cred2)
	}
}
