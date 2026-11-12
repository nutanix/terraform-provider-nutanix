package lcmv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lcmConfigPkg "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/common"
	lcmEntityPkg "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/resources"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitLcmSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixLcmPerformInventoryV2().InternalValidate(ResourceNutanixLcmPerformInventoryV2().Schema, true))
	require.NoError(t, ResourceNutanixPreChecksV2().InternalValidate(ResourceNutanixPreChecksV2().Schema, true))
	require.NoError(t, ResourceLcmUpgradeV2().InternalValidate(ResourceLcmUpgradeV2().Schema, true))
	require.NoError(t, ResourceNutanixLcmConfigV2().InternalValidate(ResourceNutanixLcmConfigV2().Schema, true))
	require.NoError(t, DatasourceNutanixLcmEntitiesV2().InternalValidate(DatasourceNutanixLcmEntitiesV2().Schema, false))
	require.NoError(t, DatasourceNutanixLcmEntityV2().InternalValidate(DatasourceNutanixLcmEntityV2().Schema, false))
	require.NoError(t, DatasourceNutanixLcmStatusV2().InternalValidate(DatasourceNutanixLcmStatusV2().Schema, false))
	require.NoError(t, DatasourceNutanixLcmConfigV2().InternalValidate(DatasourceNutanixLcmConfigV2().Schema, false))
	require.NoError(t, schemaForKeyValuePairs().InternalValidate(schemaForKeyValuePairs().Schema, true))
	require.NoError(t, schemaForAvailableVersions().InternalValidate(schemaForAvailableVersions().Schema, true))
	require.NoError(t, schemaForSubEntities().InternalValidate(schemaForSubEntities().Schema, true))
	require.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, schemaForValue().Elem.(*schema.Resource).InternalValidate(schemaForValue().Elem.(*schema.Resource).Schema, true))
}

func TestUnitFlattenLocationInfo(t *testing.T) {
	t.Parallel()

	locType := common.LOCATIONTYPE_NODE
	uuid := "loc-1"

	tests := []struct {
		name string
		in   *common.LocationInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &common.LocationInfo{
				Uuid:         &uuid,
				LocationType: &locType,
			},
			want: []map[string]interface{}{{
				"uuid":          &uuid,
				"location_type": "NODE",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocationInfo(tt.in))
		})
	}
}

func TestUnitFlattenAvailableVersions(t *testing.T) {
	t.Parallel()

	status := common.AVAILABLEVERSIONSTATUS_RECOMMENDED
	version := "1.2.3"
	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name string
		in   []lcmEntityPkg.AvailableVersion
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []lcmEntityPkg.AvailableVersion{}, want: nil},
		{
			name: "populated",
			in: []lcmEntityPkg.AvailableVersion{{
				Version:     &version,
				Status:      &status,
				ReleaseDate: &when,
			}},
			want: []map[string]interface{}{{
				"version":                &version,
				"status":                 "RECOMMENDED",
				"is_enabled":             (*bool)(nil),
				"available_version_uuid": (*string)(nil),
				"order":                  (*int64)(nil),
				"disablement_reason":     (*string)(nil),
				"release_notes":          (*string)(nil),
				"release_date":           "2024-01-02T03:04:05Z",
				"custom_message":         (*string)(nil),
				"child_entities":         []string(nil),
				"group_uuid":             (*string)(nil),
				"dependencies":           []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAvailableVersions(tt.in))
		})
	}
}

func TestUnitFlattenDependencies(t *testing.T) {
	t.Parallel()

	extID := "dep-1"
	tests := []struct {
		name string
		in   []lcmEntityPkg.DependentEntity
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []lcmEntityPkg.DependentEntity{}, want: nil},
		{
			name: "populated",
			in: []lcmEntityPkg.DependentEntity{{
				ExtId: &extID,
			}},
			want: []map[string]interface{}{{
				"tenant_id":          (*string)(nil),
				"ext_id":             &extID,
				"links":              []map[string]interface{}(nil),
				"entity_class":       (*string)(nil),
				"entity_model":       (*string)(nil),
				"entity_type":        (*common.EntityType)(nil),
				"entity_version":     (*string)(nil),
				"hardware_family":    (*string)(nil),
				"dependent_versions": []map[string]interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDependencies(tt.in))
		})
	}
}

func TestUnitFlattenKeyValuePairs(t *testing.T) {
	t.Parallel()

	name := "key1"
	tests := []struct {
		name string
		in   []lcmConfigPkg.KVPair
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []lcmConfigPkg.KVPair{}, want: nil},
		{
			name: "populated",
			in: []lcmConfigPkg.KVPair{{
				Name: &name,
			}},
			want: []map[string]interface{}{{
				"name":  &name,
				"value": []interface{}(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenKeyValuePairs(tt.in))
		})
	}
}

func TestUnitFlattenKVValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "string", in: "abc", want: []interface{}{map[string]interface{}{"string": "abc"}}},
		{name: "int", in: 7, want: []interface{}{map[string]interface{}{"integer": 7}}},
		{name: "bool", in: true, want: []interface{}{map[string]interface{}{"boolean": true}}},
		{name: "string_list", in: []string{"a"}, want: []interface{}{map[string]interface{}{"string_list": []string{"a"}}}},
		{name: "integer_list", in: []int{1, 2}, want: []interface{}{map[string]interface{}{"integer_list": []int{1, 2}}}},
		{name: "object", in: map[string]string{"k": "v"}, want: []interface{}{map[string]interface{}{"object": map[string]string{"k": "v"}}}},
		{name: "unknown", in: 1.5, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenKVValue(tt.in))
		})
	}
}

func TestUnitFlattenSubEntities(t *testing.T) {
	t.Parallel()

	entityType := common.ENTITYTYPE_SOFTWARE
	extID := "sub-1"
	tests := []struct {
		name string
		in   []common.EntityBaseModel
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []common.EntityBaseModel{}, want: nil},
		{
			name: "populated",
			in: []common.EntityBaseModel{{
				ExtId:      &extID,
				EntityType: &entityType,
			}},
			want: []map[string]interface{}{{
				"tenant_id":       (*string)(nil),
				"ext_id":          &extID,
				"links":           []map[string]interface{}(nil),
				"entity_class":    (*string)(nil),
				"entity_model":    (*string)(nil),
				"entity_type":     "SOFTWARE",
				"entity_version":  (*string)(nil),
				"hardware_family": (*string)(nil),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSubEntities(tt.in))
		})
	}
}

func TestUnitFlattenEntityTypes(t *testing.T) {
	t.Parallel()

	software := common.ENTITYTYPE_SOFTWARE
	firmware := common.ENTITYTYPE_FIRMWARE
	tests := []struct {
		name string
		in   *common.EntityType
		want string
	}{
		{name: "software", in: &software, want: "SOFTWARE"},
		{name: "firmware", in: &firmware, want: "FIRMWARE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEntityTypes(tt.in))
		})
	}
}

func TestUnitFlattenTime(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		in   *time.Time
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "populated", in: &when, want: "2024-06-01T12:00:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTime(tt.in))
		})
	}
}

func TestUnitFlattenLcmEntities(t *testing.T) {
	t.Parallel()

	entityType := common.ENTITYTYPE_FIRMWARE
	extID := "ent-1"
	tests := []struct {
		name string
		in   []lcmEntityPkg.Entity
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []lcmEntityPkg.Entity{}, want: nil},
		{
			name: "populated",
			in: []lcmEntityPkg.Entity{{
				ExtId:      &extID,
				EntityType: &entityType,
			}},
			want: []map[string]interface{}{{
				"tenant_id":                     (*string)(nil),
				"ext_id":                        &extID,
				"links":                         []map[string]interface{}(nil),
				"entity_class":                  (*string)(nil),
				"entity_model":                  (*string)(nil),
				"entity_type":                   "FIRMWARE",
				"entity_version":                (*string)(nil),
				"hardware_family":               (*string)(nil),
				"entity_description":            (*string)(nil),
				"location_info":                 []map[string]interface{}(nil),
				"target_version":                (*string)(nil),
				"last_updated_time":             "",
				"device_id":                     (*string)(nil),
				"group_uuid":                    (*string)(nil),
				"entity_details":                []map[string]interface{}(nil),
				"child_entities":                []string(nil),
				"available_versions":            []map[string]interface{}(nil),
				"sub_entities":                  []map[string]interface{}(nil),
				"cluster_ext_id":                (*string)(nil),
				"hardware_vendor":               (*string)(nil),
				"has_previous_inventory_failed": false,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLcmEntities(tt.in))
		})
	}
}

func TestUnitFlattenFrameworkVersion(t *testing.T) {
	t.Parallel()

	cur := "3.0"
	avail := "3.1"
	needed := true
	tests := []struct {
		name string
		in   *lcmEntityPkg.FrameworkVersionInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &lcmEntityPkg.FrameworkVersionInfo{
				CurrentVersion:   &cur,
				AvailableVersion: &avail,
				IsUpdateNeeded:   &needed,
			},
			want: []map[string]interface{}{{
				"current_version":   &cur,
				"available_version": &avail,
				"is_update_needed":  &needed,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFrameworkVersion(tt.in))
		})
	}
}

func TestUnitFlattenInProgressOperation(t *testing.T) {
	t.Parallel()

	opType := common.OPERATIONTYPE_INVENTORY
	opID := "op-1"
	tests := []struct {
		name string
		in   *lcmEntityPkg.InProgressOpInfo
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &lcmEntityPkg.InProgressOpInfo{
				OperationType: &opType,
				OperationId:   &opID,
			},
			want: []map[string]interface{}{{
				"operation_type": "INVENTORY",
				"operation_id":   &opID,
			}},
		},
		{
			name: "missing operation type",
			in:   &lcmEntityPkg.InProgressOpInfo{OperationId: &opID},
			want: []map[string]interface{}{{
				"operation_type": "",
				"operation_id":   &opID,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenInProgressOperation(tt.in))
		})
	}
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []import2.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []import2.ApiLink{}, want: nil},
		{
			name: "populated",
			in: []import2.ApiLink{{
				Rel:  utils.StringPtr("self"),
				Href: utils.StringPtr("https://example.test/lcm"),
			}},
			want: []map[string]interface{}{{
				"rel":  utils.StringPtr("self"),
				"href": utils.StringPtr("https://example.test/lcm"),
			}},
		},
		{
			name: "missing fields",
			in:   []import2.ApiLink{{}},
			want: []map[string]interface{}{{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinks(tt.in))
		})
	}
}

func TestUnitExpandHypervisorType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want *common.HypervisorType
	}{
		{name: "empty", in: "", want: nil},
		{name: "unknown", in: "OTHER", want: nil},
		{name: "ahv", in: "AHV", want: hypervisorTypePtr(common.HYPERVISORTYPE_AHV)},
		{name: "esx", in: "ESX", want: hypervisorTypePtr(common.HYPERVISORTYPE_ESX)},
		{name: "hyperv", in: "HYPERV", want: hypervisorTypePtr(common.HYPERVISORTYPE_HYPERV)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandHypervisorType(tt.in))
		})
	}
}

func TestUnitExpandManagementServer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want *common.ManagementServer
	}{
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"hypervisor_type": "AHV",
				"ip":              "10.0.0.1",
				"username":        "admin",
				"password":        "secret",
			}},
			want: &common.ManagementServer{
				HypervisorType: hypervisorTypePtr(common.HYPERVISORTYPE_AHV),
				Ip:             utils.StringPtr("10.0.0.1"),
				Username:       utils.StringPtr("admin"),
				Password:       utils.StringPtr("secret"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandManagementServer(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			require.Equal(t, tt.want.HypervisorType, got.HypervisorType)
			require.Equal(t, tt.want.Ip, got.Ip)
			require.Equal(t, tt.want.Username, got.Username)
			require.Equal(t, tt.want.Password, got.Password)
		})
	}
}

func TestUnitExpandEntityUpdateSpecs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []common.EntityUpdateSpec
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"entity_uuid": "e-1",
				"to_version":  "2.0",
			}},
			want: []common.EntityUpdateSpec{{
				EntityUuid: utils.StringPtr("e-1"),
				ToVersion:  utils.StringPtr("2.0"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandEntityUpdateSpecs(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, tt.want[0].EntityUuid, got[0].EntityUuid)
			require.Equal(t, tt.want[0].ToVersion, got[0].ToVersion)
		})
	}
}

func TestUnitExpandSystemAutoMgmtFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []common.SystemAutoMgmtFlag
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{
			name: "populated",
			in:   []interface{}{"POWER_OFF_UVMS", "OTHER"},
			want: []common.SystemAutoMgmtFlag{common.SYSTEMAUTOMGMTFLAG_POWER_OFF_UVMS},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandSystemAutoMgmtFlag(tt.in))
		})
	}
}

func TestUnitIsExplicitlySet(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixLcmConfigV2().Schema, map[string]interface{}{
		"url": "https://example.test/repo",
	})
	require.False(t, IsExplicitlySet(d, "url"))
	require.False(t, IsExplicitlySet(d, "missing"))
}

func TestUnitLcmResourceReadAndDeleteNoOps(t *testing.T) {
	t.Parallel()

	inv := schema.TestResourceDataRaw(t, ResourceNutanixLcmPerformInventoryV2().Schema, map[string]interface{}{})
	require.Nil(t, ResourceNutanixLcmPerformInventoryV2Read(t.Context(), inv, nil))
	require.Nil(t, ResourceNutanixLcmPerformInventoryV2Delete(t.Context(), inv, nil))

	pre := schema.TestResourceDataRaw(t, ResourceNutanixPreChecksV2().Schema, map[string]interface{}{})
	require.Nil(t, ResourceNutanixLcmPreChecksV2Read(t.Context(), pre, nil))
	require.Nil(t, ResourceNutanixLcmPreChecksV2Delete(t.Context(), pre, nil))

	up := schema.TestResourceDataRaw(t, ResourceLcmUpgradeV2().Schema, map[string]interface{}{})
	require.Nil(t, ResourceLcmUpgradeV2Read(t.Context(), up, nil))
	require.Nil(t, ResourceLcmUpgradeV2Delete(t.Context(), up, nil))

	cfg := schema.TestResourceDataRaw(t, ResourceNutanixLcmConfigV2().Schema, map[string]interface{}{})
	require.Nil(t, ResourceNutanixLcmConfigV2Read(t.Context(), cfg, nil))
	require.Nil(t, ResourceNutanixLcmConfigV2Delete(t.Context(), cfg, nil))
}

func hypervisorTypePtr(v common.HypervisorType) *common.HypervisorType { return &v }
