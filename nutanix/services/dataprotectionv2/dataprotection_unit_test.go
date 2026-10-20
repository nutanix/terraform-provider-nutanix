package dataprotectionv2

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/common"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitDataProtectionSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixRecoveryPointsV2().InternalValidate(ResourceNutanixRecoveryPointsV2().Schema, true))
	require.NoError(t, ResourceNutanixRecoveryPointRestoreV2().InternalValidate(ResourceNutanixRecoveryPointRestoreV2().Schema, true))
	require.NoError(t, ResourceNutanixRecoveryPointReplicateV2().InternalValidate(ResourceNutanixRecoveryPointReplicateV2().Schema, true))
	require.NoError(t, ResourceNutanixPromoteProtectedResourceV2().InternalValidate(ResourceNutanixPromoteProtectedResourceV2().Schema, true))
	require.NoError(t, ResourceNutanixRestoreProtectedResourceV2().InternalValidate(ResourceNutanixRestoreProtectedResourceV2().Schema, true))
	require.NoError(t, DatasourceNutanixRecoveryPointV2().InternalValidate(DatasourceNutanixRecoveryPointV2().Schema, false))
	require.NoError(t, DatasourceNutanixRecoveryPointsV2().InternalValidate(DatasourceNutanixRecoveryPointsV2().Schema, false))
	require.NoError(t, DatasourceNutanixGetProtectedResourceV2().InternalValidate(DatasourceNutanixGetProtectedResourceV2().Schema, false))
	require.NoError(t, DatasourceNutanixVMRecoveryPointInfoV2().InternalValidate(DatasourceNutanixVMRecoveryPointInfoV2().Schema, false))
	require.NoError(t, SchemaForSourceSiteReference().InternalValidate(SchemaForSourceSiteReference().Schema, true))
	require.NoError(t, SchemaForSiteProtectionInfo().InternalValidate(SchemaForSiteProtectionInfo().Schema, true))
	require.NoError(t, SchemaForReplicationStates().InternalValidate(SchemaForReplicationStates().Schema, true))
	require.NoError(t, SchemaForLinks().Elem.(*schema.Resource).InternalValidate(SchemaForLinks().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, SchemaForDiskRecoveryPoints().Elem.(*schema.Resource).InternalValidate(SchemaForDiskRecoveryPoints().Elem.(*schema.Resource).Schema, true))
}

func TestUnitFlattenTime(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		name string
		in   *time.Time
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "populated", in: &when, want: "2024-01-02T03:04:05Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenTime(tt.in))
		})
	}
}

func TestUnitFlattenStatus(t *testing.T) {
	t.Parallel()

	complete := common.RecoveryPointStatus(2)
	other := common.RecoveryPointStatus(99)
	tests := []struct {
		name string
		in   *common.RecoveryPointStatus
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "complete", in: &complete, want: "COMPLETE"},
		{name: "other", in: &other, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStatus(tt.in))
		})
	}
}

func TestUnitFlattenRecoveryPointType(t *testing.T) {
	t.Parallel()

	crash := common.RecoveryPointType(2)
	app := common.RecoveryPointType(3)
	tests := []struct {
		name string
		in   *common.RecoveryPointType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "crash", in: &crash, want: "CRASH_CONSISTENT"},
		{name: "app", in: &app, want: "APPLICATION_CONSISTENT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenRecoveryPointType(tt.in))
		})
	}
}

func TestUnitFlattenBackupType(t *testing.T) {
	t.Parallel()

	full := common.BackupType(2)
	copyBackup := common.BackupType(3)
	tests := []struct {
		name string
		in   *common.BackupType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "full", in: &full, want: "FULL_BACKUP"},
		{name: "copy", in: &copyBackup, want: "COPY_BACKUP"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenBackupType(tt.in))
		})
	}
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	rel := utils.StringPtr("self")
	href := utils.StringPtr("https://example.test/rp/1")
	tests := []struct {
		name string
		in   []response.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []response.ApiLink{}, want: nil},
		{
			name: "populated",
			in:   []response.ApiLink{{Rel: rel, Href: href}},
			want: []map[string]interface{}{{"rel": rel, "href": href}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLinks(tt.in))
		})
	}
}

func TestUnitFlattenDiskRecoveryPoints(t *testing.T) {
	t.Parallel()

	diskRP := utils.StringPtr("drp-1")
	diskExt := utils.StringPtr("disk-1")
	tests := []struct {
		name string
		in   []common.DiskRecoveryPoint
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []common.DiskRecoveryPoint{}, want: nil},
		{
			name: "populated",
			in: []common.DiskRecoveryPoint{{
				DiskRecoveryPointExtId: diskRP,
				DiskExtId:              diskExt,
			}},
			want: []map[string]interface{}{{
				"disk_recovery_point_ext_id": diskRP,
				"disk_ext_id":                diskExt,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDiskRecoveryPoints(tt.in))
		})
	}
}

func TestUnitFlattenLocationReferences(t *testing.T) {
	t.Parallel()

	loc := utils.StringPtr("loc-1")
	tests := []struct {
		name string
		in   []config.LocationReference
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.LocationReference{}, want: nil},
		{
			name: "populated",
			in:   []config.LocationReference{{LocationExtId: loc}},
			want: []map[string]interface{}{{"location_ext_id": loc}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenLocationReferences(tt.in))
		})
	}
}

func TestUnitFlattenVolumeGroupRecoveryPoints(t *testing.T) {
	t.Parallel()

	extID := utils.StringPtr("vg-rp-1")
	tests := []struct {
		name string
		in   []config.VolumeGroupRecoveryPoint
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.VolumeGroupRecoveryPoint{}, want: nil},
		{
			name: "populated",
			in:   []config.VolumeGroupRecoveryPoint{{ExtId: extID}},
			want: []map[string]interface{}{{"ext_id": extID}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenVolumeGroupRecoveryPoints(tt.in))
		})
	}
}

func TestUnitFlattenVMRecoveryPoints(t *testing.T) {
	t.Parallel()

	extID := utils.StringPtr("vm-rp-1")
	name := utils.StringPtr("vm-rp")
	tests := []struct {
		name string
		in   []config.VmRecoveryPoint
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.VmRecoveryPoint{}, want: nil},
		{
			name: "populated",
			in:   []config.VmRecoveryPoint{{ExtId: extID, Name: name}},
			want: []map[string]interface{}{{"ext_id": extID, "name": name}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenVMRecoveryPoints(tt.in))
		})
	}
}

func TestUnitFlattenApplicationConsistentProperties(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenApplicationConsistentProperties(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		props := config.NewOneOfVmRecoveryPointApplicationConsistentProperties()
		vss := common.NewVssProperties()
		bt := common.BackupType(2)
		vss.BackupType = &bt
		vss.ShouldIncludeWriters = utils.BoolPtr(true)
		vss.Writers = []string{"writer-1"}
		vss.ShouldStoreVssMetadata = utils.BoolPtr(false)
		vss.ObjectType_ = utils.StringPtr(ApplicationConsistentPropertiesVss1)
		require.NoError(t, props.SetValue(*vss))

		got := flattenApplicationConsistentProperties(props)
		require.Len(t, got, 1)
		require.Equal(t, "FULL_BACKUP", got[0]["backup_type"])
	})
}

func TestUnitFlattenRecoveryPoints(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenRecoveryPoints(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Len(t, flattenRecoveryPoints([]config.RecoveryPoint{}), 0)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		extID := utils.StringPtr("rp-1")
		got := flattenRecoveryPoints([]config.RecoveryPoint{{ExtId: extID, Name: utils.StringPtr("rp")}})
		require.Len(t, got, 1)
		m := got[0].(map[string]interface{})
		require.Equal(t, extID, m["ext_id"])
	})
}

func TestUnitFlattenEntityType(t *testing.T) {
	t.Parallel()

	vm := config.PROTECTEDENTITYTYPE_VM
	tests := []struct {
		name string
		in   *config.ProtectedEntityType
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "vm", in: &vm, want: "VM"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEntityType(tt.in))
		})
	}
}

func TestUnitFlattenSourceSiteReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSourceSiteReference(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSourceSiteReference(&config.DataProtectionSiteReference{
			DomainManagerExtId: utils.StringPtr("dm-1"),
			ClusterExtId:       utils.StringPtr("c-1"),
		})
		require.Equal(t, []map[string]interface{}{{
			"mgmt_cluster_ext_id": "dm-1",
			"cluster_ext_id":      "c-1",
		}}, got)
	})
}

func TestUnitFlattenSynchronousReplicationRole(t *testing.T) {
	t.Parallel()

	primary := config.SYNCHRONOUSREPLICATIONROLE_PRIMARY
	tests := []struct {
		name string
		in   *config.SynchronousReplicationRole
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "primary", in: &primary, want: "PRIMARY"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSynchronousReplicationRole(tt.in))
		})
	}
}

func TestUnitFlattenReplicationStatus(t *testing.T) {
	t.Parallel()

	inSync := config.PROTECTEDRESOURCEREPLICATIONSTATUS_IN_SYNC
	tests := []struct {
		name string
		in   *config.ProtectedResourceReplicationStatus
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "in sync", in: &inSync, want: "IN_SYNC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReplicationStatus(tt.in))
		})
	}
}

func TestUnitFlattenReplicationStates(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenReplicationStates(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenReplicationStates([]config.ReplicationState{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		status := config.PROTECTEDRESOURCEREPLICATIONSTATUS_IN_SYNC
		got := flattenReplicationStates([]config.ReplicationState{{
			ProtectionPolicyExtId:         utils.StringPtr("pp-1"),
			RecoveryPointObjectiveSeconds: utils.Int64Ptr(3600),
			ReplicationStatus:             &status,
		}})
		require.Len(t, got, 1)
		require.Equal(t, "pp-1", got[0]["protection_policy_ext_id"])
		require.Equal(t, "IN_SYNC", got[0]["replication_status"])
	})
}

func TestUnitFlattenSiteProtectionInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSiteProtectionInfo(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenSiteProtectionInfo([]config.SiteProtectionInfo{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		role := config.SYNCHRONOUSREPLICATIONROLE_PRIMARY
		got := flattenSiteProtectionInfo([]config.SiteProtectionInfo{{
			SynchronousReplicationRole: &role,
		}})
		require.Len(t, got, 1)
		require.Equal(t, "PRIMARY", got[0]["synchronous_replication_role"])
	})
}

func TestUnitFlattenRecoveryInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenRecoveryInfo(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRecoveryInfo(&config.RecoveryInfo{})
		require.Equal(t, []map[string]interface{}{{
			"restorable_time_ranges": []map[string]interface{}(nil),
		}}, got)
	})
}

func TestUnitFlattenRestorableTimeRanges(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenRestorableTimeRanges(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenRestorableTimeRanges([]config.RestorableTimeRange{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
		got := flattenRestorableTimeRanges([]config.RestorableTimeRange{{
			StartTime: &start,
			EndTime:   &end,
		}})
		require.Len(t, got, 1)
		require.Equal(t, start.String(), got[0]["start_time"])
		require.Equal(t, end.String(), got[0]["end_time"])
	})
}

func TestUnitExpandVolumeGroupRecoveryPoints(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVolumeGroupRecoveryPoints(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVolumeGroupRecoveryPoints([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVolumeGroupRecoveryPoints([]interface{}{
			map[string]interface{}{"volume_group_ext_id": "vg-1"},
		})
		require.Len(t, got, 1)
		require.Equal(t, "vg-1", *got[0].VolumeGroupExtId)
	})
}

func TestUnitExpandVMRecoveryPoints(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got, err := expandVMRecoveryPoints(nil)
		require.NoError(t, err)
		require.Nil(t, got)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got, err := expandVMRecoveryPoints([]interface{}{})
		require.NoError(t, err)
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got, err := expandVMRecoveryPoints([]interface{}{
			map[string]interface{}{
				"vm_ext_id":                         "vm-1",
				"name":                              "vm-rp",
				"status":                            "COMPLETE",
				"recovery_point_type":               "CRASH_CONSISTENT",
				"expiration_time":                   "2024-01-02T03:04:05Z",
				"application_consistent_properties": []interface{}{},
			},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, "vm-1", *got[0].VmExtId)
		require.Equal(t, "vm-rp", *got[0].Name)
		require.Equal(t, common.RecoveryPointStatus(2), *got[0].Status)
		require.Equal(t, common.RecoveryPointType(2), *got[0].RecoveryPointType)
	})
}

func TestUnitExpandApplicationConsistentProperties(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got, err := expandApplicationConsistentProperties(nil)
		require.NoError(t, err)
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got, err := expandApplicationConsistentProperties([]interface{}{
			map[string]interface{}{
				"backup_type":               "FULL_BACKUP",
				"should_include_writers":    true,
				"writers":                   []interface{}{"w1"},
				"should_store_vss_metadata": false,
				"object_type":               ApplicationConsistentPropertiesVss1,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, got)
	})
}

func TestUnitExpandWritersList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []string
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []interface{}{}, want: nil},
		{name: "populated", in: []interface{}{"a", "b"}, want: []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandWritersList(tt.in))
		})
	}
}

func TestUnitExpandVolumeGroupRecoveryPointRestoreOverrides(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVolumeGroupRecoveryPointRestoreOverrides([]interface{}{})
		require.Equal(t, []config.VolumeGroupRecoveryPointRestoreOverride{}, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVolumeGroupRecoveryPointRestoreOverrides([]interface{}{
			map[string]interface{}{
				"volume_group_recovery_point_ext_id": "vg-rp-1",
				"volume_group_override_spec": []interface{}{
					map[string]interface{}{"name": "new-vg"},
				},
			},
		})
		require.Len(t, got, 1)
		require.Equal(t, "vg-rp-1", *got[0].VolumeGroupRecoveryPointExtId)
		require.Equal(t, "new-vg", *got[0].VolumeGroupOverrideSpec.Name)
	})
}

func TestUnitExpandVolumeGroupOverrideSpec(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVolumeGroupOverrideSpec([]interface{}{})
		require.NotNil(t, got)
		require.Nil(t, got.Name)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVolumeGroupOverrideSpec([]interface{}{
			map[string]interface{}{"name": "vg-name"},
		})
		require.Equal(t, "vg-name", *got.Name)
	})
}

func TestUnitExpandVMRecoveryPointRestoreOverrides(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVMRecoveryPointRestoreOverrides([]interface{}{})
		require.Equal(t, []config.VmRecoveryPointRestoreOverride{}, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVMRecoveryPointRestoreOverrides([]interface{}{
			map[string]interface{}{"vm_recovery_point_ext_id": "vm-rp-1"},
		})
		require.Len(t, got, 1)
		require.Equal(t, "vm-rp-1", *got[0].VmRecoveryPointExtId)
	})
}

func TestUnitRemoveVMRecoveryPointByExtID(t *testing.T) {
	t.Parallel()

	keep := config.VmRecoveryPoint{ExtId: utils.StringPtr("keep")}
	remove := config.VmRecoveryPoint{ExtId: utils.StringPtr("remove")}
	got := removeVMRecoveryPointByExtID([]config.VmRecoveryPoint{keep, remove}, remove)
	require.Len(t, got, 1)
	require.Equal(t, "keep", *got[0].ExtId)
}

func TestUnitNoOpCRUD(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixPromoteProtectedResourceV2().Schema, map[string]interface{}{
		"ext_id": "ext-1",
	})
	require.Nil(t, ResourceNutanixPromoteProtectedResourceV2Read(t.Context(), d, nil))
	require.Nil(t, ResourceNutanixPromoteProtectedResourceV2Update(t.Context(), d, nil))
	require.Nil(t, ResourceNutanixPromoteProtectedResourceV2Delete(t.Context(), d, nil))

	restoreD := schema.TestResourceDataRaw(t, ResourceNutanixRestoreProtectedResourceV2().Schema, map[string]interface{}{
		"ext_id": "ext-1",
	})
	require.Nil(t, ResourceNutanixRestoreProtectedResourceV2Read(t.Context(), restoreD, nil))
	require.Nil(t, ResourceNutanixRestoreProtectedResourceV2Update(t.Context(), restoreD, nil))
	require.Nil(t, ResourceNutanixRestoreProtectedResourceV2Delete(t.Context(), restoreD, nil))

	rpRestore := schema.TestResourceDataRaw(t, ResourceNutanixRecoveryPointRestoreV2().Schema, map[string]interface{}{
		"ext_id":                              "rp-1",
		"cluster_ext_id":                      "c-1",
		"vm_recovery_point_restore_overrides": []interface{}{},
	})
	require.Nil(t, ResourceNutanixRecoveryPointRestoreV2Read(t.Context(), rpRestore, nil))
	require.Nil(t, ResourceNutanixRecoveryPointRestoreV2Update(t.Context(), rpRestore, nil))
	require.Nil(t, ResourceNutanixRecoveryPointRestoreV2Delete(t.Context(), rpRestore, nil))

	replicate := schema.TestResourceDataRaw(t, ResourceNutanixRecoveryPointReplicateV2().Schema, map[string]interface{}{
		"ext_id":         "rp-1",
		"pc_ext_id":      "pc-1",
		"cluster_ext_id": "c-1",
	})
	require.Nil(t, ResourceNutanixRecoveryPointReplicateV2Read(t.Context(), replicate, nil))
	require.Nil(t, ResourceNutanixRecoveryPointReplicateV2Update(t.Context(), replicate, nil))
	require.Nil(t, ResourceNutanixRecoveryPointReplicateV2Delete(t.Context(), replicate, nil))
}
