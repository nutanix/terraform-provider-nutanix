package multidomainv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import2 "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/common/v1/response"
	import7 "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/common"
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/management"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []response.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []response.ApiLink{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []response.ApiLink{{
				Rel:  utils.StringPtr("self"),
				Href: utils.StringPtr("https://example.test/projects/1"),
			}},
			want: []map[string]interface{}{{
				"rel":  "self",
				"href": "https://example.test/projects/1",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinks(tt.in))
		})
	}
}

func TestUnitMicrosToRFC3339(t *testing.T) {
	t.Parallel()

	pos := int64(1_704_067_200_000_000) // 2024-01-01T00:00:00Z
	neg := int64(-1)
	zero := int64(0)

	tests := []struct {
		name string
		in   *int64
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "zero", in: &zero, want: ""},
		{name: "negative", in: &neg, want: ""},
		{name: "populated", in: &pos, want: "2024-01-01T00:00:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, microsToRFC3339(tt.in))
		})
	}
}

func TestUnitFlattenGeographicCoordinates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.GeographicCoordinates
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "empty",
			in:   &config.GeographicCoordinates{},
			want: []map[string]interface{}{{
				"latitude":  float64(0),
				"longitude": float64(0),
			}},
		},
		{
			name: "populated",
			in: &config.GeographicCoordinates{
				Latitude:  utils.Float64Ptr(37.77),
				Longitude: utils.Float64Ptr(-122.42),
			},
			want: []map[string]interface{}{{
				"latitude":  37.77,
				"longitude": -122.42,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenGeographicCoordinates(tt.in))
		})
	}
}

func TestUnitFlattenLocations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.Location
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.Location{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.Location{{
				ExtId:    utils.StringPtr("loc-1"),
				Name:     utils.StringPtr("SF"),
				Country:  utils.StringPtr("US"),
				State:    utils.StringPtr("CA"),
				TenantId: utils.StringPtr("tenant-1"),
				Coordinates: &config.GeographicCoordinates{
					Latitude:  utils.Float64Ptr(1.5),
					Longitude: utils.Float64Ptr(2.5),
				},
				Links: []response.ApiLink{{
					Rel:  utils.StringPtr("self"),
					Href: utils.StringPtr("https://example.test/locations/loc-1"),
				}},
			}},
			want: []map[string]interface{}{{
				"ext_id":    "loc-1",
				"name":      "SF",
				"country":   "US",
				"state":     "CA",
				"tenant_id": "tenant-1",
				"coordinates": []map[string]interface{}{{
					"latitude":  1.5,
					"longitude": 2.5,
				}},
				"links": []map[string]interface{}{{
					"rel":  "self",
					"href": "https://example.test/locations/loc-1",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocations(tt.in))
		})
	}
}

func TestUnitFlattenLocationProjections(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.LocationProjection
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.LocationProjection{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.LocationProjection{{
				ExtId: utils.StringPtr("loc-p-1"),
				Name:  utils.StringPtr("NYC"),
			}},
			want: []map[string]interface{}{{
				"ext_id":      "loc-p-1",
				"name":        "NYC",
				"country":     "",
				"state":       "",
				"tenant_id":   "",
				"coordinates": []map[string]interface{}{},
				"links":       []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocationProjections(tt.in))
		})
	}
}

func TestUnitFlattenProjects(t *testing.T) {
	t.Parallel()

	created := int64(1_704_067_200_000_000)
	state := config.PROJECTSTATE_ACTIVE

	tests := []struct {
		name string
		in   []config.Project
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.Project{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.Project{{
				ExtId:             utils.StringPtr("proj-1"),
				Name:              utils.StringPtr("project"),
				Description:       utils.StringPtr("desc"),
				Id:                utils.StringPtr("id-1"),
				TenantId:          utils.StringPtr("tenant-1"),
				IsDefault:         utils.BoolPtr(true),
				IsSystemDefined:   utils.BoolPtr(false),
				CreatedBy:         utils.StringPtr("admin"),
				UpdatedBy:         utils.StringPtr("admin2"),
				CreatedTimestamp:  &created,
				ModifiedTimestamp: &created,
				State:             &state,
				Links: []response.ApiLink{{
					Rel:  utils.StringPtr("self"),
					Href: utils.StringPtr("https://example.test/projects/proj-1"),
				}},
			}},
			want: []map[string]interface{}{{
				"ext_id":             "proj-1",
				"name":               "project",
				"description":        "desc",
				"project_id":         "id-1",
				"tenant_id":          "tenant-1",
				"is_default":         true,
				"is_system_defined":  false,
				"created_by":         "admin",
				"updated_by":         "admin2",
				"created_timestamp":  "2024-01-01T00:00:00Z",
				"modified_timestamp": "2024-01-01T00:00:00Z",
				"state":              "ACTIVE",
				"links": []map[string]interface{}{{
					"rel":  "self",
					"href": "https://example.test/projects/proj-1",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenProjects(tt.in))
		})
	}
}

func TestUnitFlattenProjectProjections(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.ProjectProjection
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.ProjectProjection{}, want: []map[string]interface{}{}},
		{
			name: "populated stub returns empty",
			in:   []config.ProjectProjection{{ExtId: utils.StringPtr("p-1")}},
			want: []map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenProjectProjections(tt.in))
		})
	}
}

func TestUnitFlattenApiCredentialStatus(t *testing.T) {
	t.Parallel()

	valid := config.APICREDENTIALSTATUS_VALID
	tests := []struct {
		name string
		in   *config.ApiCredentialStatus
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "populated", in: &valid, want: "VALID"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenApiCredentialStatus(tt.in))
		})
	}
}

func TestUnitFlattenRegistrationState(t *testing.T) {
	t.Parallel()

	registered := import7.REGISTRATIONSTATE_REGISTERED
	tests := []struct {
		name string
		in   *import7.RegistrationState
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "populated", in: &registered, want: "REGISTERED"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRegistrationState(tt.in))
		})
	}
}

func TestUnitFlattenTunnelProvider(t *testing.T) {
	t.Parallel()

	ikat := import7.TUNNELPROVIDER_IKAT
	tests := []struct {
		name string
		in   *import7.TunnelProvider
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "populated", in: &ikat, want: "IKAT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTunnelProvider(tt.in))
		})
	}
}

func TestUnitFlattenRegisteredDomainPlatformData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.PlatformData
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "empty",
			in:   &config.PlatformData{},
			want: []map[string]interface{}{{
				"cluster_count":       0,
				"connectivity_status": "",
				"pc_env_type":         "",
				"pc_provider_type":    "",
				"pc_version":          "",
			}},
		},
		{
			name: "populated",
			in: &config.PlatformData{
				ClusterCount:       utils.IntPtr(3),
				ConnectivityStatus: utils.StringPtr("CONNECTED"),
				PcEnvType:          utils.StringPtr("ONPREM"),
				PcProviderType:     utils.StringPtr("NTNX"),
				PcVersion:          utils.StringPtr("7.0"),
			},
			want: []map[string]interface{}{{
				"cluster_count":       3,
				"connectivity_status": "CONNECTED",
				"pc_env_type":         "ONPREM",
				"pc_provider_type":    "NTNX",
				"pc_version":          "7.0",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRegisteredDomainPlatformData(tt.in))
		})
	}
}

func TestUnitFlattenRegistrationCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          *import7.RegistrationCredentials
		withSecrets bool
		want        []map[string]interface{}
	}{
		{name: "nil", in: nil, withSecrets: false, want: []map[string]interface{}{}},
		{
			name:        "empty resource",
			in:          &import7.RegistrationCredentials{},
			withSecrets: false,
			want:        []map[string]interface{}{{}},
		},
		{
			name: "populated omits write-only secrets on resource",
			in: &import7.RegistrationCredentials{
				ApiKey: utils.StringPtr("key"),
				KeyId:  utils.StringPtr("kid"),
			},
			withSecrets: false,
			want:        []map[string]interface{}{{}},
		},
		{
			name: "populated includes secrets for datasource",
			in: &import7.RegistrationCredentials{
				ApiKey: utils.StringPtr("key"),
				KeyId:  utils.StringPtr("kid"),
			},
			withSecrets: true,
			want: []map[string]interface{}{{
				"api_key": "key",
				"key_id":  "kid",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRegistrationCredentials(tt.in, tt.withSecrets))
		})
	}
}

func TestUnitFlattenRegisteredDomainRegistrationConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		in          *config.RegistrationConfig
		withSecrets bool
		want        []map[string]interface{}
	}{
		{name: "nil", in: nil, withSecrets: false, want: []map[string]interface{}{}},
		{
			name:        "empty",
			in:          &config.RegistrationConfig{},
			withSecrets: false,
			want: []map[string]interface{}{{
				"domain_ext_id": "",
				"target_url":    "",
				"tenant_ext_id": "",
				"credentials":   []map[string]interface{}{},
			}},
		},
		{
			name: "populated resource omits secrets",
			in: &config.RegistrationConfig{
				DomainExtId: utils.StringPtr("dom-1"),
				TargetUrl:   utils.StringPtr("https://example.test"),
				TenantExtId: utils.StringPtr("ten-1"),
				Credentials: &import7.RegistrationCredentials{
					ApiKey: utils.StringPtr("key"),
					KeyId:  utils.StringPtr("kid"),
				},
			},
			withSecrets: false,
			want: []map[string]interface{}{{
				"domain_ext_id": "dom-1",
				"target_url":    "https://example.test",
				"tenant_ext_id": "ten-1",
				"credentials":   []map[string]interface{}{{}},
			}},
		},
		{
			name: "populated datasource includes secrets",
			in: &config.RegistrationConfig{
				DomainExtId: utils.StringPtr("dom-1"),
				TargetUrl:   utils.StringPtr("https://example.test"),
				TenantExtId: utils.StringPtr("ten-1"),
				Credentials: &import7.RegistrationCredentials{
					ApiKey: utils.StringPtr("key"),
					KeyId:  utils.StringPtr("kid"),
				},
			},
			withSecrets: true,
			want: []map[string]interface{}{{
				"domain_ext_id": "dom-1",
				"target_url":    "https://example.test",
				"tenant_ext_id": "ten-1",
				"credentials": []map[string]interface{}{{
					"api_key": "key",
					"key_id":  "kid",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRegisteredDomainRegistrationConfig(tt.in, tt.withSecrets))
		})
	}
}

func TestUnitFlattenRegisteredDomainTunnel(t *testing.T) {
	t.Parallel()

	provider := import7.TUNNELPROVIDER_IKAT
	tests := []struct {
		name string
		in   *config.DomainTunnel
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "empty",
			in:   &config.DomainTunnel{},
			want: []map[string]interface{}{{
				"endpoint": "",
				"provider": "",
			}},
		},
		{
			name: "populated",
			in: &config.DomainTunnel{
				Endpoint: utils.StringPtr("tunnel.example.test"),
				Provider: &provider,
			},
			want: []map[string]interface{}{{
				"endpoint": "tunnel.example.test",
				"provider": "IKAT",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRegisteredDomainTunnel(tt.in))
		})
	}
}

func TestUnitFlattenRegisteredDomains(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	status := config.APICREDENTIALSTATUS_VALID
	state := import7.REGISTRATIONSTATE_REGISTERED

	tests := []struct {
		name string
		in   []config.RegisteredDomain
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.RegisteredDomain{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.RegisteredDomain{{
				ExtId:               utils.StringPtr("rd-1"),
				Name:                utils.StringPtr("domain"),
				DnsName:             utils.StringPtr("dns"),
				Fqdn:                utils.StringPtr("fqdn.example"),
				LocationExtId:       utils.StringPtr("loc-1"),
				CredentialIssuer:    utils.StringPtr("issuer"),
				CreatedBy:           utils.StringPtr("admin"),
				ApiCredentialStatus: &status,
				RegistrationState:   &state,
				RegisteredTime:      &when,
				TenantId:            utils.StringPtr("tenant-1"),
			}},
			want: []map[string]interface{}{{
				"ext_id":                "rd-1",
				"name":                  "domain",
				"dns_name":              "dns",
				"fqdn":                  "fqdn.example",
				"domain_ext_id":         "rd-1",
				"location_ext_id":       "loc-1",
				"credential_issuer":     "issuer",
				"created_by":            "admin",
				"api_credential_status": "VALID",
				"registration_state":    "REGISTERED",
				"registered_time":       "2024-01-02T03:04:05Z",
				"platform_data":         []map[string]interface{}{},
				"registration_config":   []map[string]interface{}{},
				"tunnel":                []map[string]interface{}{},
				"tenant_id":             "tenant-1",
				"links":                 []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRegisteredDomains(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroupStorageContainers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.StorageContainerDetails
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.StorageContainerDetails{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in:   []config.StorageContainerDetails{{ExtId: utils.StringPtr("sc-1")}},
			want: []map[string]interface{}{{"ext_id": "sc-1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroupStorageContainers(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroupPlacementTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.TargetDetails
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.TargetDetails{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.TargetDetails{{
				ClusterExtId: utils.StringPtr("cluster-1"),
				StorageContainers: []config.StorageContainerDetails{
					{ExtId: utils.StringPtr("sc-1")},
				},
			}},
			want: []map[string]interface{}{{
				"cluster_ext_id": "cluster-1",
				"storage_containers": []map[string]interface{}{
					{"ext_id": "sc-1"},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroupPlacementTargets(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroupCapabilities(t *testing.T) {
	t.Parallel()

	val := import2.NewOneOfKVPairValue()
	require.NoError(t, val.SetValue("cap-val"))

	tests := []struct {
		name string
		in   []import2.KVPair
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []import2.KVPair{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []import2.KVPair{{
				Name:  utils.StringPtr("cap"),
				Value: val,
			}},
			want: []map[string]interface{}{{
				"name":  "cap",
				"value": "cap-val",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroupCapabilities(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroupStorageContainersIncludingCapabilities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.StorageContainerDetails
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.StorageContainerDetails{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.StorageContainerDetails{{
				ExtId: utils.StringPtr("sc-2"),
			}},
			want: []map[string]interface{}{{
				"ext_id":       "sc-2",
				"capabilities": []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroupStorageContainersIncludingCapabilities(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroupPlacementTargetsIncludingCapabilities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.TargetDetails
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.TargetDetails{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.TargetDetails{{
				ClusterExtId: utils.StringPtr("cluster-2"),
			}},
			want: []map[string]interface{}{{
				"cluster_ext_id":     "cluster-2",
				"storage_containers": []map[string]interface{}{},
				"capabilities":       []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroupPlacementTargetsIncludingCapabilities(tt.in))
		})
	}
}

func TestUnitFlattenResourceGroups(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	tests := []struct {
		name string
		in   []config.ResourceGroup
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []config.ResourceGroup{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []config.ResourceGroup{{
				ExtId:          utils.StringPtr("rg-1"),
				Name:           utils.StringPtr("rg"),
				ProjectExtId:   utils.StringPtr("proj-1"),
				TenantId:       utils.StringPtr("tenant-1"),
				CreatedBy:      utils.StringPtr("admin"),
				LastUpdatedBy:  utils.StringPtr("admin2"),
				CreateTime:     &when,
				LastUpdateTime: &when,
			}},
			want: []map[string]interface{}{{
				"ext_id":            "rg-1",
				"name":              "rg",
				"project_ext_id":    "proj-1",
				"tenant_id":         "tenant-1",
				"created_by":        "admin",
				"last_updated_by":   "admin2",
				"placement_targets": []map[string]interface{}{},
				"links":             []map[string]interface{}{},
				"capabilities":      []map[string]interface{}{},
				"create_time":       "2024-02-03T04:05:06Z",
				"last_update_time":  "2024-02-03T04:05:06Z",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenResourceGroups(tt.in))
		})
	}
}

func TestUnitFlattenLocalDomainGeographicCoordinates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.GeographicCoordinates
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "empty",
			in:   &config.GeographicCoordinates{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &config.GeographicCoordinates{
				Latitude:  utils.Float64Ptr(10.1),
				Longitude: utils.Float64Ptr(20.2),
			},
			want: []map[string]interface{}{{
				"latitude":  10.1,
				"longitude": 20.2,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocalDomainGeographicCoordinates(tt.in))
		})
	}
}

func TestUnitFlattenLocalDomainLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.Location
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "empty",
			in:   &config.Location{},
			want: []map[string]interface{}{{
				"ext_id":      "",
				"name":        "",
				"country":     "",
				"state":       "",
				"coordinates": []map[string]interface{}{},
			}},
		},
		{
			name: "populated",
			in: &config.Location{
				ExtId:   utils.StringPtr("loc-ld"),
				Name:    utils.StringPtr("Local"),
				Country: utils.StringPtr("US"),
				State:   utils.StringPtr("WA"),
			},
			want: []map[string]interface{}{{
				"ext_id":      "loc-ld",
				"name":        "Local",
				"country":     "US",
				"state":       "WA",
				"coordinates": []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocalDomainLocation(tt.in))
		})
	}
}

func TestUnitFlattenLocalDomainServiceStatuses(t *testing.T) {
	t.Parallel()

	svc := management.SERVICENAME_IDENTITY_MANAGEMENT
	st := management.CONNECTIONSTATUS_ACTIVE
	tests := []struct {
		name string
		in   []management.ServiceStatus
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: []management.ServiceStatus{}, want: []map[string]interface{}{}},
		{
			name: "populated",
			in: []management.ServiceStatus{{
				ServiceName: &svc,
				Status:      &st,
			}},
			want: []map[string]interface{}{{
				"service_name": "IDENTITY_MANAGEMENT",
				"status":       "ACTIVE",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocalDomainServiceStatuses(tt.in))
		})
	}
}

func TestUnitFlattenLocalDomainConnectionStatus(t *testing.T) {
	t.Parallel()

	st := management.CONNECTIONSTATUS_ACTIVE
	tests := []struct {
		name string
		in   *management.LocalDomainConnectionStatus
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "empty",
			in:   &management.LocalDomainConnectionStatus{},
			want: []map[string]interface{}{{
				"service_statuses": []map[string]interface{}{},
			}},
		},
		{
			name: "populated",
			in: &management.LocalDomainConnectionStatus{
				Status: &st,
			},
			want: []map[string]interface{}{{
				"status":           "ACTIVE",
				"service_statuses": []map[string]interface{}{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocalDomainConnectionStatus(tt.in))
		})
	}
}

func TestUnitFlattenLocalDomainRegistrationStatus(t *testing.T) {
	t.Parallel()

	state := import7.REGISTRATIONSTATE_REGISTERED
	tests := []struct {
		name string
		in   *management.LocalDomainRegistrationStatus
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "empty",
			in:   &management.LocalDomainRegistrationStatus{},
			want: []map[string]interface{}{{}},
		},
		{
			name: "populated",
			in: &management.LocalDomainRegistrationStatus{
				DaysToExpiry:          utils.IntPtr(30),
				IdentityProviderExtId: utils.StringPtr("idp-1"),
				State:                 &state,
			},
			want: []map[string]interface{}{{
				"days_to_expiry":           30,
				"identity_provider_ext_id": "idp-1",
				"state":                    "REGISTERED",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocalDomainRegistrationStatus(tt.in))
		})
	}
}

func TestUnitExpandRegisteredDomainTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		wantNil bool
		want    time.Time
	}{
		{name: "empty", in: "", wantNil: true},
		{name: "invalid", in: "not-a-time", wantNil: true},
		{
			name:    "populated",
			in:      "2024-01-02T03:04:05Z",
			wantNil: false,
			want:    time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandRegisteredDomainTime(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.True(t, got.Equal(tt.want))
		})
	}
}

func TestUnitExpandRegisteredDomainPlatformData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantNil bool
		check   func(t *testing.T, got *config.PlatformData)
	}{
		{name: "nil", in: nil, wantNil: true},
		{name: "empty", in: []interface{}{}, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"cluster_count":       2,
				"connectivity_status": "UP",
				"pc_env_type":         "CLOUD",
				"pc_provider_type":    "AWS",
				"pc_version":          "6.8",
			}},
			check: func(t *testing.T, got *config.PlatformData) {
				t.Helper()
				require.Equal(t, 2, utils.IntValue(got.ClusterCount))
				require.Equal(t, "UP", utils.StringValue(got.ConnectivityStatus))
				require.Equal(t, "CLOUD", utils.StringValue(got.PcEnvType))
				require.Equal(t, "AWS", utils.StringValue(got.PcProviderType))
				require.Equal(t, "6.8", utils.StringValue(got.PcVersion))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandRegisteredDomainPlatformData(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			tt.check(t, got)
		})
	}
}

func TestUnitExpandRegistrationCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantNil bool
		apiKey  string
		keyID   string
	}{
		{name: "nil", in: nil, wantNil: true},
		{name: "empty", in: []interface{}{}, wantNil: true},
		{
			name:   "populated",
			in:     []interface{}{map[string]interface{}{"api_key": "k", "key_id": "id"}},
			apiKey: "k",
			keyID:  "id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandRegistrationCredentials(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.apiKey, utils.StringValue(got.ApiKey))
			require.Equal(t, tt.keyID, utils.StringValue(got.KeyId))
		})
	}
}

func TestUnitExpandRegisteredDomainRegistrationConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantNil bool
	}{
		{name: "nil", in: nil, wantNil: true},
		{name: "empty", in: []interface{}{}, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"domain_ext_id": "dom",
				"target_url":    "https://t",
				"tenant_ext_id": "ten",
				"credentials": []interface{}{
					map[string]interface{}{"api_key": "k", "key_id": "id"},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandRegisteredDomainRegistrationConfig(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "dom", utils.StringValue(got.DomainExtId))
			require.Equal(t, "https://t", utils.StringValue(got.TargetUrl))
			require.Equal(t, "ten", utils.StringValue(got.TenantExtId))
			require.NotNil(t, got.Credentials)
			require.Equal(t, "k", utils.StringValue(got.Credentials.ApiKey))
		})
	}
}

func TestUnitExpandRegisteredDomainTunnel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantNil bool
	}{
		{name: "nil", in: nil, wantNil: true},
		{name: "empty", in: []interface{}{}, wantNil: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"endpoint": "ep",
				"provider": "IKAT",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandRegisteredDomainTunnel(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, "ep", utils.StringValue(got.Endpoint))
			require.NotNil(t, got.Provider)
			require.Equal(t, "IKAT", flattenTunnelProvider(got.Provider))
		})
	}
}

func TestUnitExpandRegisteredDomain(t *testing.T) {
	t.Parallel()

	dEmpty := schema.TestResourceDataRaw(t, ResourceNutanixRegisteredDomainV2().Schema, map[string]interface{}{})
	emptyBody := expandRegisteredDomain(dEmpty)
	require.NotNil(t, emptyBody)
	require.Nil(t, emptyBody.Name)

	d := schema.TestResourceDataRaw(t, ResourceNutanixRegisteredDomainV2().Schema, map[string]interface{}{
		"name":                  "dom",
		"dns_name":              "dns",
		"fqdn":                  "fqdn",
		"domain_ext_id":         "dom-ext",
		"location_ext_id":       "loc",
		"credential_issuer":     "issuer",
		"created_by":            "admin",
		"api_credential_status": "VALID",
		"registration_state":    "REGISTERED",
		"registered_time":       "2024-01-02T03:04:05Z",
		"platform_data": []interface{}{map[string]interface{}{
			"cluster_count": 1,
			"pc_version":    "7.0",
		}},
		"registration_config": []interface{}{map[string]interface{}{
			"domain_ext_id": "dom-cfg",
			"target_url":    "https://example.test",
		}},
		"tunnel": []interface{}{map[string]interface{}{
			"endpoint": "tunnel",
			"provider": "IKAT",
		}},
	})
	got := expandRegisteredDomain(d)
	require.Equal(t, "dom", utils.StringValue(got.Name))
	require.Equal(t, "dns", utils.StringValue(got.DnsName))
	require.Equal(t, "fqdn", utils.StringValue(got.Fqdn))
	require.Equal(t, "dom-ext", utils.StringValue(got.DomainExtId))
	require.Equal(t, "loc", utils.StringValue(got.LocationExtId))
	require.NotNil(t, got.PlatformData)
	require.NotNil(t, got.RegistrationConfig)
	require.NotNil(t, got.Tunnel)
}

func TestUnitExpandResourceGroupStorageContainers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []config.StorageContainerDetails
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in:   []interface{}{map[string]interface{}{"ext_id": "sc-9"}},
			want: []config.StorageContainerDetails{{ExtId: utils.StringPtr("sc-9")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandResourceGroupStorageContainers(tt.in)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnitExpandResourceGroupPlacementTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []config.TargetDetails
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"cluster_ext_id": "c-1",
				"storage_containers": []interface{}{
					map[string]interface{}{"ext_id": "sc-1"},
				},
			}},
			want: []config.TargetDetails{{
				ClusterExtId: utils.StringPtr("c-1"),
				StorageContainers: []config.StorageContainerDetails{
					{ExtId: utils.StringPtr("sc-1")},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandResourceGroupPlacementTargets(tt.in))
		})
	}
}

func TestUnitExpandLocalDomainRegistrationCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []interface{}
		wantNil bool
		apiKey  string
		keyID   string
	}{
		{name: "nil", in: nil, wantNil: true},
		{name: "empty", in: []interface{}{}, wantNil: true},
		{
			name:   "populated",
			in:     []interface{}{map[string]interface{}{"api_key": "api", "key_id": "kid"}},
			apiKey: "api",
			keyID:  "kid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandLocalDomainRegistrationCredentials(tt.in)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.apiKey, utils.StringValue(got.ApiKey))
			require.Equal(t, tt.keyID, utils.StringValue(got.KeyId))
		})
	}
}

func TestUnitRegisteredDomainRefreshReadAndDelete(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixRegisteredDomainRefreshApiCredentialsV2().Schema, map[string]interface{}{
		"ext_id": "rd-1",
	})
	require.Nil(t, ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Read(t.Context(), d, nil))
	require.Nil(t, ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Delete(t.Context(), d, nil))
}
