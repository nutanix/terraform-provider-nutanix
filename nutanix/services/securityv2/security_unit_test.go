package securityv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonCfg "github.com/nutanix-core/ntnx-api-golang-sdk-internal/security-go-client/v17/models/common/v1/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/security-go-client/v17/models/common/v1/response"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/security-go-client/v17/models/security/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/security-go-client/v17/models/security/v4/report"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitSecuritySchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixKeyManagementServerV2().InternalValidate(ResourceNutanixKeyManagementServerV2().Schema, true))
	require.NoError(t, DatasourceNutanixKeyManagementServerV2().InternalValidate(DatasourceNutanixKeyManagementServerV2().Schema, false))
	require.NoError(t, DatasourceNutanixKeyManagementServersV2().InternalValidate(DatasourceNutanixKeyManagementServersV2().Schema, false))
	require.NoError(t, DatasourceNutanixStigsControlsV2().InternalValidate(DatasourceNutanixStigsControlsV2().Schema, false))
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	href := utils.StringPtr("https://example.test/kms/1")
	rel := utils.StringPtr("self")

	tests := []struct {
		name string
		in   []response.ApiLink
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []response.ApiLink{}, want: nil},
		{
			name: "populated",
			in:   []response.ApiLink{{Href: href, Rel: rel}},
			want: []interface{}{map[string]interface{}{"href": href, "rel": rel}},
		},
		{
			name: "missing fields",
			in:   []response.ApiLink{{}},
			want: []interface{}{map[string]interface{}{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinks(tt.in))
		})
	}
}

func TestUnitFlattenIPv4Address(t *testing.T) {
	t.Parallel()

	value := utils.StringPtr("10.0.0.1")
	prefix := utils.IntPtr(24)

	tests := []struct {
		name string
		in   *commonCfg.IPv4Address
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &commonCfg.IPv4Address{},
			want: []map[string]interface{}{{"value": "", "prefix_length": 0}},
		},
		{
			name: "populated",
			in:   &commonCfg.IPv4Address{Value: value, PrefixLength: prefix},
			want: []map[string]interface{}{{"value": "10.0.0.1", "prefix_length": 24}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPv4Address(tt.in))
		})
	}
}

func TestUnitFlattenIPv6Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *commonCfg.IPv6Address
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &commonCfg.IPv6Address{},
			want: []map[string]interface{}{{"value": "", "prefix_length": 0}},
		},
		{
			name: "populated",
			in: &commonCfg.IPv6Address{
				Value:        utils.StringPtr("2001:db8::1"),
				PrefixLength: utils.IntPtr(64),
			},
			want: []map[string]interface{}{{"value": "2001:db8::1", "prefix_length": 64}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPv6Address(tt.in))
		})
	}
}

func TestUnitFlattenFQDN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *commonCfg.FQDN
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &commonCfg.FQDN{},
			want: []map[string]interface{}{{"value": ""}},
		},
		{
			name: "populated",
			in:   &commonCfg.FQDN{Value: utils.StringPtr("kms.example.test")},
			want: []map[string]interface{}{{"value": "kms.example.test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFQDN(tt.in))
		})
	}
}

func TestUnitFlattenIPAddressOrFQDN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *commonCfg.IPAddressOrFQDN
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "empty",
			in:   &commonCfg.IPAddressOrFQDN{},
			want: []map[string]interface{}{{
				"ipv4": []map[string]interface{}(nil),
				"ipv6": []map[string]interface{}(nil),
				"fqdn": []map[string]interface{}(nil),
			}},
		},
		{
			name: "populated ipv4",
			in: &commonCfg.IPAddressOrFQDN{
				Ipv4: &commonCfg.IPv4Address{
					Value:        utils.StringPtr("10.0.0.8"),
					PrefixLength: utils.IntPtr(24),
				},
			},
			want: []map[string]interface{}{{
				"ipv4": []map[string]interface{}{{"value": "10.0.0.8", "prefix_length": 24}},
				"ipv6": []map[string]interface{}(nil),
				"fqdn": []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIPAddressOrFQDN(tt.in))
		})
	}
}

func TestUnitFlattenAccessInformation(t *testing.T) {
	t.Parallel()

	expiry := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got, err := flattenAccessInformation(nil)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("unsupported type", func(t *testing.T) {
		t.Parallel()
		got, err := flattenAccessInformation("bad")
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("azure populated", func(t *testing.T) {
		t.Parallel()
		azure := config.AzureAccessInformation{
			EndpointUrl:           utils.StringPtr("https://vault.azure.net"),
			KeyId:                 utils.StringPtr("key-1"),
			TenantId:              utils.StringPtr("tenant-1"),
			ClientId:              utils.StringPtr("client-1"),
			TruncatedClientSecret: utils.StringPtr("****cret"),
			CredentialExpiryDate:  &expiry,
		}
		got, err := flattenAccessInformation(azure)
		require.NoError(t, err)
		require.Equal(t, []map[string]interface{}{{
			"azure_key_vault": []map[string]interface{}{{
				"endpoint_url":            "https://vault.azure.net",
				"key_id":                  "key-1",
				"tenant_id":               "tenant-1",
				"client_id":               "client-1",
				"truncated_client_secret": "****cret",
				"credential_expiry_date":  "2026-06-15",
			}},
		}}, got)
	})

	t.Run("kmip populated", func(t *testing.T) {
		t.Parallel()
		kmip := config.KmipAccessInformation{
			CaName:     utils.StringPtr("ca"),
			CaPem:      utils.StringPtr("ca-pem"),
			CertPem:    utils.StringPtr("cert-pem"),
			PrivateKey: utils.StringPtr("private-key"),
			Endpoints: []config.EndpointInfo{{
				IpAddress: &commonCfg.IPAddressOrFQDN{
					Ipv4: &commonCfg.IPv4Address{
						Value:        utils.StringPtr("10.1.1.1"),
						PrefixLength: utils.IntPtr(32),
					},
				},
				Port: utils.IntPtr(5696),
			}},
		}
		got, err := flattenAccessInformation(kmip)
		require.NoError(t, err)
		require.Equal(t, []map[string]interface{}{{
			"kmip_key_vault": []map[string]interface{}{{
				"ca_name": "ca",
				"endpoint_url": []interface{}{
					map[string]interface{}{
						"ip_address": []map[string]interface{}{{
							"ipv4": []map[string]interface{}{{"value": "10.1.1.1", "prefix_length": 32}},
							"ipv6": []map[string]interface{}(nil),
							"fqdn": []map[string]interface{}(nil),
						}},
						"port": 5696,
					},
				},
			}},
		}}, got)
	})

	t.Run("azure pointer", func(t *testing.T) {
		t.Parallel()
		azure := &config.AzureAccessInformation{
			EndpointUrl:           utils.StringPtr("https://vault.azure.net"),
			KeyId:                 utils.StringPtr("key-1"),
			TenantId:              utils.StringPtr("tenant-1"),
			ClientId:              utils.StringPtr("client-1"),
			TruncatedClientSecret: utils.StringPtr("****cret"),
			CredentialExpiryDate:  &expiry,
		}
		got, err := flattenAccessInformation(azure)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Contains(t, got[0], "azure_key_vault")
	})
}

func TestUnitFlattenKeyManagementServer(t *testing.T) {
	t.Parallel()

	expiry := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	created := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	href := utils.StringPtr("https://example.test/kms/ext-1")
	rel := utils.StringPtr("self")

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got, err := flattenKeyManagementServer(nil)
		require.NoError(t, err)
		require.Equal(t, []map[string]interface{}{}, got)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got, err := flattenKeyManagementServer([]config.KeyManagementServer{})
		require.NoError(t, err)
		require.Equal(t, []map[string]interface{}{}, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		kms := config.NewKeyManagementServer()
		kms.Name = utils.StringPtr("kms-1")
		kms.ExtId = utils.StringPtr("ext-1")
		kms.TenantId = utils.StringPtr("tenant-1")
		kms.CreationTimestamp = &created
		kms.Links = []response.ApiLink{{Href: href, Rel: rel}}

		azure := config.NewAzureAccessInformation()
		azure.EndpointUrl = utils.StringPtr("https://vault.azure.net")
		azure.KeyId = utils.StringPtr("key-1")
		azure.TenantId = utils.StringPtr("tenant-1")
		azure.ClientId = utils.StringPtr("client-1")
		azure.TruncatedClientSecret = utils.StringPtr("****cret")
		azure.CredentialExpiryDate = &expiry
		require.NoError(t, kms.SetAccessInformation(*azure))

		got, err := flattenKeyManagementServer([]config.KeyManagementServer{*kms})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "kms-1", got[0]["name"])
		require.Equal(t, "ext-1", got[0]["ext_id"])
		require.Equal(t, "tenant-1", got[0]["tenant_id"])
		require.Equal(t, created.Format(time.RFC3339), got[0]["creation_timestamp"])
		require.Equal(t, []interface{}{map[string]interface{}{"href": href, "rel": rel}}, got[0]["links"])
		require.NotNil(t, got[0]["access_information"])
	})

	t.Run("missing access information", func(t *testing.T) {
		t.Parallel()
		got, err := flattenKeyManagementServer([]config.KeyManagementServer{{}})
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestUnitFlattenSTIGs(t *testing.T) {
	t.Parallel()

	severity := report.SEVERITY_HIGH
	status := report.STIGSTATUS_APPLICABLE
	extID := utils.StringPtr("stig-1")
	tenantID := utils.StringPtr("tenant-1")

	tests := []struct {
		name string
		in   []report.Stig
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []report.Stig{}, want: []map[string]interface{}{}},
		{
			name: "empty entity",
			in:   []report.Stig{{}},
			want: []map[string]interface{}{{
				"ext_id":            (*string)(nil),
				"tenant_id":         (*string)(nil),
				"links":             []interface{}(nil),
				"title":             "",
				"rule_id":           "",
				"stig_version":      "",
				"identifiers":       []string(nil),
				"affected_clusters": []string(nil),
				"severity":          "",
				"status":            "",
				"comments":          "",
				"fix_text":          "",
				"benchmark_id":      "",
			}},
		},
		{
			name: "populated",
			in: []report.Stig{{
				ExtId:            extID,
				TenantId:         tenantID,
				Title:            utils.StringPtr("Ensure password complexity"),
				RuleId:           utils.StringPtr("V-123"),
				StigVersion:      utils.StringPtr("1.0"),
				Identifiers:      []string{"CCI-0001"},
				AffectedClusters: []string{"cluster-1"},
				Severity:         &severity,
				Status:           &status,
				Comments:         utils.StringPtr("ok"),
				FixText:          utils.StringPtr("fix it"),
				BenchmarkId:      utils.StringPtr("bench-1"),
			}},
			want: []map[string]interface{}{{
				"ext_id":            extID,
				"tenant_id":         tenantID,
				"links":             []interface{}(nil),
				"title":             "Ensure password complexity",
				"rule_id":           "V-123",
				"stig_version":      "1.0",
				"identifiers":       []string{"CCI-0001"},
				"affected_clusters": []string{"cluster-1"},
				"severity":          "HIGH",
				"status":            "APPLICABLE",
				"comments":          "ok",
				"fix_text":          "fix it",
				"benchmark_id":      "bench-1",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSTIGs(tt.in))
		})
	}
}

func TestUnitExpandIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4Address(nil))
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4Address([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Address([]interface{}{map[string]interface{}{
			"value":         "10.0.0.1",
			"prefix_length": 24,
		}})
		require.NotNil(t, got)
		require.Equal(t, "10.0.0.1", utils.StringValue(got.Value))
		require.Equal(t, 24, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandIPv6Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv6Address(nil))
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv6Address([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Address([]interface{}{map[string]interface{}{
			"value":         "2001:db8::1",
			"prefix_length": 64,
		}})
		require.NotNil(t, got)
		require.Equal(t, "2001:db8::1", utils.StringValue(got.Value))
		require.Equal(t, 64, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandFQDN(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandFQDN(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandFQDN([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandFQDN([]interface{}{map[string]interface{}{"value": "kms.example.test"}})
		require.NotNil(t, got)
		require.Equal(t, "kms.example.test", utils.StringValue(got.Value))
	})
}

func TestUnitExpandIPAddressOrFQDN(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPAddressOrFQDN(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPAddressOrFQDN([]interface{}{}))
	})

	t.Run("populated ipv4", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddressOrFQDN([]interface{}{map[string]interface{}{
			"ipv4": []interface{}{map[string]interface{}{
				"value":         "10.0.0.1",
				"prefix_length": 24,
			}},
		}})
		require.NotNil(t, got)
		require.NotNil(t, got.Ipv4)
		require.Equal(t, "10.0.0.1", utils.StringValue(got.Ipv4.Value))
		require.Equal(t, 24, utils.IntValue(got.Ipv4.PrefixLength))
	})
}

func TestUnitExpandKMIPEndpoints(t *testing.T) {
	t.Parallel()

	t.Run("invalid type", func(t *testing.T) {
		t.Parallel()
		got, err := expandKMIPEndpoints("bad")
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got, err := expandKMIPEndpoints([]interface{}{})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got, err := expandKMIPEndpoints([]interface{}{map[string]interface{}{
			"ip_address": []interface{}{map[string]interface{}{
				"ipv4": []interface{}{map[string]interface{}{
					"value":         "10.1.1.1",
					"prefix_length": 32,
				}},
			}},
			"port": 5696,
		}})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, 5696, utils.IntValue(got[0].Port))
		require.NotNil(t, got[0].IpAddress)
		require.True(t, got[0].IpAddress.IsValid())
	})
}

func TestUnitExpandAccessInformation(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got, err := expandAccessInformation(nil)
		require.Error(t, err)
		require.Nil(t, got)

		got, err = expandAccessInformation([]interface{}{})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("neither vault", func(t *testing.T) {
		t.Parallel()
		got, err := expandAccessInformation([]interface{}{map[string]interface{}{}})
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("azure from schema", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixKeyManagementServerV2().Schema, map[string]interface{}{
			"name": "kms-azure",
			"access_information": []interface{}{map[string]interface{}{
				"azure_key_vault": []interface{}{map[string]interface{}{
					"endpoint_url":           "https://vault.azure.net",
					"key_id":                 "key-1",
					"tenant_id":              "tenant-1",
					"client_id":              "client-1",
					"client_secret":          "secret12",
					"credential_expiry_date": "2026-06-15",
				}},
			}},
		})
		got, err := expandAccessInformation(d.Get("access_information").([]interface{}))
		require.NoError(t, err)
		azure, ok := got.(config.AzureAccessInformation)
		require.True(t, ok)
		require.Equal(t, "https://vault.azure.net", utils.StringValue(azure.EndpointUrl))
		require.Equal(t, "key-1", utils.StringValue(azure.KeyId))
		require.Equal(t, "tenant-1", utils.StringValue(azure.TenantId))
		require.Equal(t, "client-1", utils.StringValue(azure.ClientId))
		require.Equal(t, "secret12", utils.StringValue(azure.ClientSecret))
		require.Equal(t, "2026-06-15", azure.CredentialExpiryDate.Format("2006-01-02"))
	})

	t.Run("kmip from schema", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixKeyManagementServerV2().Schema, map[string]interface{}{
			"name": "kms-kmip",
			"access_information": []interface{}{map[string]interface{}{
				"kmip_key_vault": []interface{}{map[string]interface{}{
					"cert_pem":    "cert",
					"private_key": "key",
					"ca_name":     "ca",
					"ca_pem":      "ca-pem",
					"endpoint_url": []interface{}{map[string]interface{}{
						"ip_address": []interface{}{map[string]interface{}{
							"ipv4": []interface{}{map[string]interface{}{
								"value":         "10.1.1.1",
								"prefix_length": 32,
							}},
						}},
						"port": 5696,
					}},
				}},
			}},
		})
		got, err := expandAccessInformation(d.Get("access_information").([]interface{}))
		require.NoError(t, err)
		kmip, ok := got.(config.KmipAccessInformation)
		require.True(t, ok)
		require.Equal(t, "ca", utils.StringValue(kmip.CaName))
		require.Equal(t, "ca-pem", utils.StringValue(kmip.CaPem))
		require.Equal(t, "cert", utils.StringValue(kmip.CertPem))
		require.Equal(t, "key", utils.StringValue(kmip.PrivateKey))
		require.Len(t, kmip.Endpoints, 1)
	})
}
