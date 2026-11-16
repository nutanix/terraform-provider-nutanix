package clustersv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clusterapi "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/api"
	clusterclient "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/client"
	prismapi "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/api"
	prismclient "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/client"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
)

// These tests run Terraform against a local API stub; they never contact a cluster.
func TestSSLCertificateWriteOnly(t *testing.T) {
	const name = "nutanix_ssl_certificate_v2.test"
	var mu sync.Mutex
	var uploaded []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/ssl-certificate") && r.Method == http.MethodGet:
			w.Header().Set("ETag", "test-etag")
			fmt.Fprint(w, `{"data":{"$objectType":"clustermgmt.v4.config.SSLCertificate","publicCertificate":"Y2VydGlmaWNhdGU=","privateKeyAlgorithm":"RSA_2048"}}`)
		case strings.HasSuffix(r.URL.Path, "/ssl-certificate") && r.Method == http.MethodPut:
			var body struct {
				PrivateKey string `json:"privateKey"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
				w.WriteHeader(400)
				return
			}
			mu.Lock()
			uploaded = append(uploaded, body.PrivateKey)
			mu.Unlock()
			fmt.Fprint(w, `{"data":{"$objectType":"prism.v4.config.TaskReference","extId":"test-task"}}`)
		case strings.HasSuffix(r.URL.Path, "/tasks/test-task") && r.Method == http.MethodGet:
			fmt.Fprint(w, `{"data":{"$objectType":"prism.v4.config.Task","extId":"test-task","status":"SUCCEEDED"}}`)
		default:
			t.Errorf("unexpected API request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)
	port, _ := strconv.Atoi(u.Port())
	factory := func() *schema.Provider {
		c := clusterclient.NewApiClient()
		c.Host = u.Hostname()
		c.Port = port
		c.Scheme = "http"
		c.AllowVersionNegotiation = false
		p := prismclient.NewApiClient()
		p.Host = u.Hostname()
		p.Port = port
		p.Scheme = "http"
		p.AllowVersionNegotiation = false
		return &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{"nutanix_ssl_certificate_v2": ResourceNutanixSSLCertificateV2()},
			ConfigureContextFunc: func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
				return &conns.Client{ClusterAPI: &clusters.Client{SSLCertificateAPI: clusterapi.NewSSLCertificateServiceApi(c)}, PrismAPI: &prism.Client{TaskRefAPI: prismapi.NewTasksServiceApi(p)}}, nil
			},
		}
	}
	config := func(keyArguments string) string {
		return fmt.Sprintf(`resource "nutanix_ssl_certificate_v2" "test" {
   cluster_ext_id = "test-cluster"
   public_certificate = "Y2VydGlmaWNhdGU="
   ca_chain = "Y2E="
   private_key_algorithm = "RSA_2048"
   %s
  }`, keyArguments)
	}
	ephemeralConfig := func(key string, version int) string {
		return fmt.Sprintf(`variable "private_key" {
          type = string
          sensitive = true
          ephemeral = true
          default = %q
        }
        `, key) + config(fmt.Sprintf("private_key = var.private_key\nprivate_key_version = %d", version))
	}

	check := func(expected []string, writeOnly bool) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			mu.Lock()
			defer mu.Unlock()
			if fmt.Sprint(uploaded) != fmt.Sprint(expected) {
				return fmt.Errorf("unexpected uploads: got %d, want %d", len(uploaded), len(expected))
			}
			if writeOnly {
				attrs := s.RootModule().Resources[name].Primary.Attributes
				if attrs["private_key"] != "" {
					return fmt.Errorf("private key persisted in state")
				}
				serialized, err := json.Marshal(s)
				if err != nil {
					return err
				}
				for _, secret := range expected {
					if strings.Contains(string(serialized), secret) {
						return fmt.Errorf("secret present in serialized state")
					}
				}
			}
			return nil
		}
	}
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){"nutanix": func() (*schema.Provider, error) { return factory(), nil }},
		Steps: []resource.TestStep{
			{Config: config(`private_key = "legacy-key"`), Check: check([]string{"legacy-key"}, true)},
			// Existing configurations remain valid without a version; ephemeral inputs also work.
			{Config: ephemeralConfig("write-only-key", 1), Check: check([]string{"legacy-key", "write-only-key"}, true)},
			// The certificate is unchanged: the version alone must cause an upload.
			{Config: ephemeralConfig("rotated-key", 2), Check: check([]string{"legacy-key", "write-only-key", "rotated-key"}, true)},
			{Config: config(`private_key_version = 1`), PlanOnly: true, ExpectError: regexp.MustCompile("private_key must be configured")},
			{Config: config("private_key = \"key\"\nprivate_key_version = 0"), PlanOnly: true, ExpectError: regexp.MustCompile("expected private_key_version to be at least")},
			// Leave a valid configuration for the harness's final destroy.
			{Config: ephemeralConfig("rotated-key", 2), Check: check([]string{"legacy-key", "write-only-key", "rotated-key"}, true)},
		},
	})
}
