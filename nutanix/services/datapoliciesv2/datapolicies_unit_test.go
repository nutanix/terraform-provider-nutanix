package datapoliciesv2

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/models/common/v1/response"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/models/datapolicies/v4/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/models/dataprotection/v4/common"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitDataPoliciesSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixProtectionPoliciesV2().InternalValidate(ResourceNutanixProtectionPoliciesV2().Schema, true))
	require.NoError(t, ResourceNutanixStoragePoliciesV2().InternalValidate(ResourceNutanixStoragePoliciesV2().Schema, true))
	require.NoError(t, DatasourceNutanixProtectionPolicyV2().InternalValidate(DatasourceNutanixProtectionPolicyV2().Schema, false))
	require.NoError(t, DatasourceNutanixProtectionPoliciesV2().InternalValidate(DatasourceNutanixProtectionPoliciesV2().Schema, false))
	require.NoError(t, DataSourceNutanixStoragePolicyV2().InternalValidate(DataSourceNutanixStoragePolicyV2().Schema, false))
	require.NoError(t, DataSourceNutanixStoragePoliciesV2().InternalValidate(DataSourceNutanixStoragePoliciesV2().Schema, false))
	require.NoError(t, schemaReplicationLocations().InternalValidate(schemaReplicationLocations().Schema, true))
	require.NoError(t, schemaReplicationConfigurations().InternalValidate(schemaReplicationConfigurations().Schema, true))
	require.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, true))
}

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
				Href: utils.StringPtr("https://example.test/policies/1"),
			}},
			want: []map[string]interface{}{{
				"rel":  "self",
				"href": "https://example.test/policies/1",
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

func TestUnitFlattenReplicationLocations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.ReplicationLocation
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.ReplicationLocation{}, want: nil},
		{
			name: "populated",
			in: []config.ReplicationLocation{{
				Label:              utils.StringPtr("primary"),
				DomainManagerExtId: utils.StringPtr("dm-1"),
				IsPrimary:          utils.BoolPtr(true),
			}},
			want: []map[string]interface{}{{
				"label":                 utils.StringPtr("primary"),
				"domain_manager_ext_id": utils.StringPtr("dm-1"),
				"is_primary":            utils.BoolPtr(true),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReplicationLocations(tt.in))
		})
	}
}

func TestUnitFlattenReplicationSubLocation(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenReplicationSubLocation(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		sub := config.NewOneOfReplicationLocationReplicationSubLocation()
		cluster := config.NewNutanixCluster()
		cluster.ClusterExtIds = []string{"c-1", "c-2"}
		require.NoError(t, sub.SetValue(*cluster))

		got := flattenReplicationSubLocation(sub)
		require.Len(t, got, 1)
		require.Contains(t, got[0], "cluster_ext_ids")
	})
}

func TestUnitFlattenReplicationConfigurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []config.ReplicationConfiguration
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []config.ReplicationConfiguration{}, want: nil},
		{
			name: "populated",
			in: []config.ReplicationConfiguration{{
				SourceLocationLabel: utils.StringPtr("primary"),
				RemoteLocationLabel: utils.StringPtr("remote"),
			}},
			want: []map[string]interface{}{{
				"source_location_label": utils.StringPtr("primary"),
				"remote_location_label": utils.StringPtr("remote"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenReplicationConfigurations(tt.in))
		})
	}
}

func TestUnitFlattenSchedule(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSchedule(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		rpt := common.RECOVERYPOINTTYPE_CRASH_CONSISTENT
		got := flattenSchedule(&config.Schedule{
			RecoveryPointType:                 &rpt,
			RecoveryPointObjectiveTimeSeconds: utils.IntPtr(3600),
			StartTime:                         utils.StringPtr("14:00:00"),
			IsReplicationPaused:               utils.BoolPtr(false),
			ScheduleExtId:                     utils.StringPtr("sched-1"),
		})
		require.Equal(t, []map[string]interface{}{{
			"recovery_point_type":                   "CRASH_CONSISTENT",
			"recovery_point_objective_time_seconds": 3600,
			"start_time":                            "14:00:00",
			"is_replication_paused":                 false,
			"schedule_ext_id":                       "sched-1",
		}}, got)
	})
}

func TestUnitFlattenRetention(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenRetention(nil))
	})

	t.Run("linear", func(t *testing.T) {
		t.Parallel()
		ret := config.NewOneOfScheduleRetention()
		linear := config.NewLinearRetention()
		linear.Local = utils.IntPtr(5)
		linear.Remote = utils.IntPtr(3)
		require.NoError(t, ret.SetValue(*linear))

		got := flattenRetention(ret)
		require.Equal(t, []map[string]interface{}{{
			"linear_retention": []map[string]interface{}{{
				"local":  5,
				"remote": 3,
			}},
		}}, got)
	})
}

func TestUnitFlattenAutoRollupRetentionDetails(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenAutoRollupRetentionDetails(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		interval := config.SNAPSHOTINTERVALTYPE_DAILY
		got := flattenAutoRollupRetentionDetails(&config.AutoRollupRetentionDetails{
			SnapshotIntervalType: &interval,
			Frequency:            utils.IntPtr(7),
		})
		require.Equal(t, []map[string]interface{}{{
			"snapshot_interval_type": "DAILY",
			"frequency":              7,
		}}, got)
	})
}

func TestUnitFlattenProtectionPolicies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      []config.ProtectionPolicy
		wantLen int
	}{
		{name: "nil", in: nil, wantLen: 0},
		{name: "empty", in: []config.ProtectionPolicy{}, wantLen: 0},
		{
			name: "populated",
			in: []config.ProtectionPolicy{{
				ExtId:       utils.StringPtr("pp-1"),
				Name:        utils.StringPtr("policy"),
				TenantId:    utils.StringPtr("t-1"),
				Description: utils.StringPtr("desc"),
			}},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := flattenProtectionPolicies(tt.in)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				require.Equal(t, "policy", got[0]["name"])
				require.Equal(t, utils.StringPtr("pp-1"), got[0]["ext_id"])
			}
		})
	}
}

func TestUnitFlattenCompressionSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenCompressionSpec(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		state := config.COMPRESSIONSTATE_INLINE
		got := flattenCompressionSpec(&config.CompressionSpec{CompressionState: &state})
		require.Equal(t, []interface{}{map[string]interface{}{"compression_state": "INLINE"}}, got)
	})
}

func TestUnitFlattenEncryptionSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenEncryptionSpec(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		state := config.ENCRYPTIONSTATE_ENABLED
		got := flattenEncryptionSpec(&config.EncryptionSpec{EncryptionState: &state})
		require.Equal(t, []interface{}{map[string]interface{}{"encryption_state": "ENABLED"}}, got)
	})
}

func TestUnitFlattenQosSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenQosSpec(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		iops := 1000
		got := flattenQosSpec(&config.QosSpec{ThrottledIops: &iops})
		require.Equal(t, []interface{}{map[string]interface{}{"throttled_iops": &iops}}, got)
	})
}

func TestUnitFlattenFaultToleranceSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenFaultToleranceSpec(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		rf := config.REPLICATIONFACTOR_TWO
		got := flattenFaultToleranceSpec(&config.FaultToleranceSpec{ReplicationFactor: &rf})
		require.Equal(t, []interface{}{map[string]interface{}{"replication_factor": "TWO"}}, got)
	})
}

func TestUnitFlattenStoragePolicies(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenStoragePolicies(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenStoragePolicies([]config.StoragePolicy{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		pt := config.POLICYTYPE_USER
		got := flattenStoragePolicies([]config.StoragePolicy{{
			ExtId:      utils.StringPtr("sp-1"),
			Name:       utils.StringPtr("storage"),
			TenantId:   utils.StringPtr("t-1"),
			PolicyType: &pt,
		}})
		require.Len(t, got, 1)
		require.Equal(t, "storage", got[0]["name"])
		require.Equal(t, "USER", got[0]["policy_type"])
	})
}

func TestUnitExpandReplicationLocations(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandReplicationLocations(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandReplicationLocations([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandReplicationLocations([]interface{}{
			map[string]interface{}{
				"label":                    "primary",
				"domain_manager_ext_id":    "dm-1",
				"is_primary":               true,
				"replication_sub_location": []interface{}{},
			},
		})
		require.Len(t, got, 1)
		require.Equal(t, "primary", *got[0].Label)
		require.Equal(t, "dm-1", *got[0].DomainManagerExtId)
		require.True(t, *got[0].IsPrimary)
	})
}

func TestUnitExpandOneOfReplicationLocationReplicationSubLocation(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandOneOfReplicationLocationReplicationSubLocation(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandOneOfReplicationLocationReplicationSubLocation([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfReplicationLocationReplicationSubLocation([]interface{}{
			map[string]interface{}{
				"cluster_ext_ids": []interface{}{
					map[string]interface{}{
						"cluster_ext_ids": []interface{}{"c-1"},
					},
				},
			},
		})
		require.NotNil(t, got)
	})
}

func TestUnitExpandReplicationConfigurations(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandReplicationConfigurations(nil, false, cty.NullVal(cty.DynamicPseudoType)))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandReplicationConfigurations([]interface{}{}, false, cty.NullVal(cty.DynamicPseudoType)))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandReplicationConfigurations([]interface{}{
			map[string]interface{}{
				"source_location_label": "primary",
				"remote_location_label": "remote",
				"schedule":              []interface{}{},
			},
		}, false, cty.NullVal(cty.DynamicPseudoType))
		require.Len(t, got, 1)
		require.Equal(t, "primary", *got[0].SourceLocationLabel)
		require.Equal(t, "remote", *got[0].RemoteLocationLabel)
	})
}

func TestUnitExpandSchedule(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandSchedule(nil, false, false))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandSchedule([]interface{}{}, false, false))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandSchedule([]interface{}{
			map[string]interface{}{
				"recovery_point_type":                   "CRASH_CONSISTENT",
				"recovery_point_objective_time_seconds": 3600,
				"start_time":                            "14:00:00",
				"is_replication_paused":                 false,
				"retention":                             []interface{}{},
			},
		}, false, false)
		require.NotNil(t, got)
		require.Equal(t, common.RECOVERYPOINTTYPE_CRASH_CONSISTENT, *got.RecoveryPointType)
		require.Equal(t, 3600, *got.RecoveryPointObjectiveTimeSeconds)
		require.Equal(t, "14:00:00", *got.StartTime)
	})
}

func TestUnitExpandRecoveryPointType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want *common.RecoveryPointType
	}{
		{name: "empty", in: "", want: nil},
		{name: "crash", in: "CRASH_CONSISTENT", want: recoveryPointTypePtr(common.RECOVERYPOINTTYPE_CRASH_CONSISTENT)},
		{name: "app", in: "APPLICATION_CONSISTENT", want: recoveryPointTypePtr(common.RECOVERYPOINTTYPE_APPLICATION_CONSISTENT)},
		{name: "unknown", in: "OTHER", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandRecoveryPointType(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
			} else {
				require.Equal(t, *tt.want, *got)
			}
		})
	}
}

func TestUnitExpandRetention(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandRetention(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandRetention([]interface{}{}))
	})

	t.Run("linear", func(t *testing.T) {
		t.Parallel()
		got := expandRetention([]interface{}{
			map[string]interface{}{
				"linear_retention": []interface{}{
					map[string]interface{}{"local": 5, "remote": 3},
				},
			},
		})
		require.NotNil(t, got)
	})
}

func TestUnitExpandLinearRetention(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandLinearRetention(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandLinearRetention([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandLinearRetention([]interface{}{
			map[string]interface{}{"local": 5, "remote": 3},
		})
		require.Equal(t, 5, *got.Local)
		require.Equal(t, 3, *got.Remote)
	})
}

func TestUnitExpandAutoRollupRetention(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAutoRollupRetention(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAutoRollupRetention([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandAutoRollupRetention([]interface{}{
			map[string]interface{}{
				"local": []interface{}{
					map[string]interface{}{
						"snapshot_interval_type": "DAILY",
						"frequency":              7,
					},
				},
				"remote": []interface{}{},
			},
		})
		require.NotNil(t, got)
		require.NotNil(t, got.Local)
		require.Equal(t, config.SNAPSHOTINTERVALTYPE_DAILY, *got.Local.SnapshotIntervalType)
		require.Equal(t, 7, *got.Local.Frequency)
	})
}

func TestUnitExpandAutoRollupRetentionDetails(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAutoRollupRetentionDetails(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandAutoRollupRetentionDetails([]interface{}{}))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandAutoRollupRetentionDetails([]interface{}{
			map[string]interface{}{
				"snapshot_interval_type": "HOURLY",
				"frequency":              2,
			},
		})
		require.Equal(t, config.SNAPSHOTINTERVALTYPE_HOURLY, *got.SnapshotIntervalType)
		require.Equal(t, 2, *got.Frequency)
	})
}

func TestUnitExpandSnapshotIntervalType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want *config.SnapshotIntervalType
	}{
		{name: "empty", in: "", want: nil},
		{name: "hourly", in: "HOURLY", want: snapshotIntervalPtr(config.SNAPSHOTINTERVALTYPE_HOURLY)},
		{name: "daily", in: "DAILY", want: snapshotIntervalPtr(config.SNAPSHOTINTERVALTYPE_DAILY)},
		{name: "weekly", in: "WEEKLY", want: snapshotIntervalPtr(config.SNAPSHOTINTERVALTYPE_WEEKLY)},
		{name: "monthly", in: "MONTHLY", want: snapshotIntervalPtr(config.SNAPSHOTINTERVALTYPE_MONTHLY)},
		{name: "yearly", in: "YEARLY", want: snapshotIntervalPtr(config.SNAPSHOTINTERVALTYPE_YEARLY)},
		{name: "unknown", in: "OTHER", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandSnapshotIntervalType(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
			} else {
				require.Equal(t, *tt.want, *got)
			}
		})
	}
}

func TestUnitBuildSpecs(t *testing.T) {
	t.Parallel()

	t.Run("compression", func(t *testing.T) {
		t.Parallel()
		got := buildCompressionSpec(map[string]interface{}{"compression_state": "INLINE"})
		require.Equal(t, "INLINE", got.CompressionState.GetName())
	})

	t.Run("encryption", func(t *testing.T) {
		t.Parallel()
		got := buildEncryptionSpec(map[string]interface{}{"encryption_state": "ENABLED"})
		require.Equal(t, "ENABLED", got.EncryptionState.GetName())
	})

	t.Run("qos", func(t *testing.T) {
		t.Parallel()
		got := buildQosSpec(map[string]interface{}{"throttled_iops": 500})
		require.Equal(t, 500, *got.ThrottledIops)
	})

	t.Run("fault tolerance", func(t *testing.T) {
		t.Parallel()
		got := buildFaultToleranceSpec(map[string]interface{}{"replication_factor": "TWO"})
		require.Equal(t, "TWO", got.ReplicationFactor.GetName())
	})
}

func TestUnitIsLatestRetentionSetInConfig(t *testing.T) {
	t.Parallel()

	t.Run("null", func(t *testing.T) {
		t.Parallel()
		require.False(t, isLatestRetentionSetInConfig(cty.NullVal(cty.DynamicPseudoType), 0))
	})
}

func recoveryPointTypePtr(v common.RecoveryPointType) *common.RecoveryPointType { return &v }

func snapshotIntervalPtr(v config.SnapshotIntervalType) *config.SnapshotIntervalType { return &v }
