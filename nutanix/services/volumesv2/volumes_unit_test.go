package volumesv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/config"
	volumesClientResponse "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/response"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitVolumesSchemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, ResourceNutanixVolumeGroupV2().InternalValidate(ResourceNutanixVolumeGroupV2().Schema, true))
	require.NoError(t, ResourceNutanixVolumeGroupDiskV2().InternalValidate(ResourceNutanixVolumeGroupDiskV2().Schema, true))
	require.NoError(t, ResourceNutanixVolumeGroupIscsiClientV2().InternalValidate(ResourceNutanixVolumeGroupIscsiClientV2().Schema, true))
	require.NoError(t, ResourceNutanixVolumeAttachVMToVolumeGroupV2().InternalValidate(ResourceNutanixVolumeAttachVMToVolumeGroupV2().Schema, true))
	require.NoError(t, ResourceNutanixAssociateCategoryToVolumeGroupV2().InternalValidate(ResourceNutanixAssociateCategoryToVolumeGroupV2().Schema, true))

	require.NoError(t, DatasourceNutanixVolumeGroupV2().InternalValidate(DatasourceNutanixVolumeGroupV2().Schema, false))
	require.NoError(t, DatasourceNutanixVolumeGroupsV2().InternalValidate(DatasourceNutanixVolumeGroupsV2().Schema, false))
	require.NoError(t, DatasourceNutanixVolumeDiskV2().InternalValidate(DatasourceNutanixVolumeDiskV2().Schema, false))
	require.NoError(t, DatasourceNutanixVolumeDisksV2().InternalValidate(DatasourceNutanixVolumeDisksV2().Schema, false))
	require.NoError(t, DatasourceNutanixVolumeIscsiClientV2().InternalValidate(DatasourceNutanixVolumeIscsiClientV2().Schema, false))
	require.NoError(t, DatasourceNutanixVolumeIscsiClientsV2().InternalValidate(DatasourceNutanixVolumeIscsiClientsV2().Schema, false))
}

func TestUnitVolumesNoopCRUD(t *testing.T) {
	t.Parallel()

	t.Run("associate category read", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixAssociateCategoryToVolumeGroupV2().Schema, map[string]interface{}{
			"ext_id": "vg-1",
		})
		require.Nil(t, ResourceNutanixAssociateCategoryToVolumeGroupV2Read(t.Context(), d, nil))
	})
	t.Run("attach vm read update", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVolumeAttachVMToVolumeGroupV2().Schema, map[string]interface{}{
			"volume_group_ext_id": "vg-1",
			"vm_ext_id":           "vm-1",
		})
		require.Nil(t, ResourceNutanixVolumeAttachVMToVolumeGroupV2Read(t.Context(), d, nil))
		require.Nil(t, ResourceNutanixVolumeAttachVMToVolumeGroupV2Update(t.Context(), d, nil))
	})
	t.Run("iscsi client read update", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, ResourceNutanixVolumeGroupIscsiClientV2().Schema, map[string]interface{}{
			"vg_ext_id": "vg-1",
		})
		require.Nil(t, ResourceNutanixVolumeGroupIscsiClientV2Read(t.Context(), d, nil))
		require.Nil(t, ResourceNutanixVolumeGroupIscsiClientV2Update(t.Context(), d, nil))
	})
}

func TestUnitFlattenVolumesEntities(t *testing.T) {
	t.Parallel()

	sharing := volumesClient.SHARINGSTATUS_SHARED
	auth := volumesClient.AUTHENTICATIONTYPE_CHAP
	usage := volumesClient.USAGETYPE_USER
	tests := []struct {
		name string
		in   []volumesClient.VolumeGroup
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []volumesClient.VolumeGroup{}, want: nil},
		{
			name: "populated",
			in: []volumesClient.VolumeGroup{{
				Name:                           utils.StringPtr("vg-1"),
				ExtId:                          utils.StringPtr("ext-1"),
				TenantId:                       utils.StringPtr("t-1"),
				Description:                    utils.StringPtr("desc"),
				ShouldLoadBalanceVmAttachments: utils.BoolPtr(true),
				SharingStatus:                  &sharing,
				TargetName:                     utils.StringPtr("target-1"),
				EnabledAuthentications:         &auth,
				CreatedBy:                      utils.StringPtr("admin"),
				ClusterReference:               utils.StringPtr("c-1"),
				UsageType:                      &usage,
				IsHidden:                       utils.BoolPtr(false),
				ProjectExtId:                   utils.StringPtr("p-1"),
				IscsiFeatures: &volumesClient.IscsiFeatures{
					EnabledAuthentications: &auth,
				},
				StorageFeatures: &volumesClient.StorageFeatures{
					FlashMode: &volumesClient.FlashMode{IsEnabled: utils.BoolPtr(true)},
				},
			}},
			want: []interface{}{map[string]interface{}{
				"name":                               utils.StringPtr("vg-1"),
				"ext_id":                             utils.StringPtr("ext-1"),
				"tenant_id":                          utils.StringPtr("t-1"),
				"description":                        utils.StringPtr("desc"),
				"should_load_balance_vm_attachments": utils.BoolPtr(true),
				"sharing_status":                     "SHARED",
				"target_name":                        utils.StringPtr("target-1"),
				"enabled_authentications":            "CHAP",
				"created_by":                         utils.StringPtr("admin"),
				"cluster_reference":                  utils.StringPtr("c-1"),
				"usage_type":                         "USER",
				"is_hidden":                          utils.BoolPtr(false),
				"project_ext_id":                     utils.StringPtr("p-1"),
				"iscsi_features": []map[string]interface{}{{
					"enabled_authentications": "CHAP",
				}},
				"storage_features": []map[string]interface{}{{
					"flash_mode": []map[string]interface{}{{
						"is_enabled": utils.BoolPtr(true),
					}},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenVolumesEntities(tt.in))
		})
	}
}

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []volumesClientResponse.ApiLink
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []volumesClientResponse.ApiLink{}, want: nil},
		{
			name: "populated",
			in: []volumesClientResponse.ApiLink{{
				Href: utils.StringPtr("https://example.test/vg/1"),
				Rel:  utils.StringPtr("self"),
			}},
			want: []map[string]interface{}{{
				"href": utils.StringPtr("https://example.test/vg/1"),
				"rel":  utils.StringPtr("self"),
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

func TestUnitFlattenSharingStatus(t *testing.T) {
	t.Parallel()

	shared := volumesClient.SHARINGSTATUS_SHARED
	notShared := volumesClient.SHARINGSTATUS_NOT_SHARED
	tests := []struct {
		name string
		in   *volumesClient.SharingStatus
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "shared", in: &shared, want: "SHARED"},
		{name: "not shared", in: &notShared, want: "NOT_SHARED"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenSharingStatus(tt.in))
		})
	}
}

func TestUnitFlattenEnabledAuthentications(t *testing.T) {
	t.Parallel()

	chap := volumesClient.AUTHENTICATIONTYPE_CHAP
	none := volumesClient.AUTHENTICATIONTYPE_NONE
	tests := []struct {
		name string
		in   *volumesClient.AuthenticationType
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "chap", in: &chap, want: "CHAP"},
		{name: "none", in: &none, want: "NONE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenEnabledAuthentications(tt.in))
		})
	}
}

func TestUnitFlattenIscsiFeatures(t *testing.T) {
	t.Parallel()

	chap := volumesClient.AUTHENTICATIONTYPE_CHAP
	tests := []struct {
		name string
		in   *volumesClient.IscsiFeatures
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &volumesClient.IscsiFeatures{EnabledAuthentications: &chap},
			want: []map[string]interface{}{{"enabled_authentications": "CHAP"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIscsiFeatures(tt.in))
		})
	}
}

func TestUnitFlattenFlashMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *volumesClient.FlashMode
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &volumesClient.FlashMode{IsEnabled: utils.BoolPtr(true)},
			want: []map[string]interface{}{{"is_enabled": utils.BoolPtr(true)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFlashMode(tt.in))
		})
	}
}

func TestUnitFlattenStorageFeatures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *volumesClient.StorageFeatures
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &volumesClient.StorageFeatures{
				FlashMode: &volumesClient.FlashMode{IsEnabled: utils.BoolPtr(false)},
			},
			want: []map[string]interface{}{{
				"flash_mode": []map[string]interface{}{{"is_enabled": utils.BoolPtr(false)}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenStorageFeatures(tt.in))
		})
	}
}

func TestUnitFlattenUsageType(t *testing.T) {
	t.Parallel()

	user := volumesClient.USAGETYPE_USER
	internal := volumesClient.USAGETYPE_INTERNAL
	temporary := volumesClient.USAGETYPE_TEMPORARY
	backup := volumesClient.USAGETYPE_BACKUP_TARGET
	tests := []struct {
		name string
		in   *volumesClient.UsageType
		want string
	}{
		{name: "nil", in: nil, want: ""},
		{name: "user", in: &user, want: "USER"},
		{name: "internal", in: &internal, want: "INTERNAL"},
		{name: "temporary", in: &temporary, want: "TEMPORARY"},
		{name: "backup", in: &backup, want: "BACKUP_TARGET"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenUsageType(tt.in))
		})
	}
}

func TestUnitFlattenIscsiClientsEntities(t *testing.T) {
	t.Parallel()

	site := volumesClient.VOLUMEGROUPATTACHMENTSITE_PRIMARY
	auth := volumesClient.AUTHENTICATIONTYPE_NONE
	tests := []struct {
		name string
		in   []volumesClient.IscsiClient
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []volumesClient.IscsiClient{}, want: nil},
		{
			name: "populated",
			in: []volumesClient.IscsiClient{{
				ExtId:                  utils.StringPtr("ic-1"),
				TenantId:               utils.StringPtr("t-1"),
				IscsiInitiatorName:     utils.StringPtr("iqn.2024.test"),
				EnabledAuthentications: &auth,
				AttachmentSite:         &site,
				ClusterReference:       utils.StringPtr("c-1"),
				AttachedTargets: []volumesClient.TargetParam{{
					NumVirtualTargets: utils.IntPtr(1),
					IscsiTargetName:   utils.StringPtr("tgt-1"),
				}},
				IscsiInitiatorNetworkId: &config.IPAddressOrFQDN{
					Ipv4: &config.IPv4Address{Value: utils.StringPtr("10.0.0.2"), PrefixLength: utils.IntPtr(32)},
				},
			}},
			want: []interface{}{map[string]interface{}{
				"ext_id":                  utils.StringPtr("ic-1"),
				"tenant_id":               utils.StringPtr("t-1"),
				"iscsi_initiator_name":    utils.StringPtr("iqn.2024.test"),
				"enabled_authentications": "NONE",
				"attachment_site":         "PRIMARY",
				"cluster_reference":       utils.StringPtr("c-1"),
				"attached_targets": []interface{}{map[string]interface{}{
					"num_virtual_targets": utils.IntPtr(1),
					"iscsi_target_name":   utils.StringPtr("tgt-1"),
				}},
				"iscsi_initiator_network_id": []interface{}{map[string]interface{}{
					"ipv4": []interface{}{map[string]interface{}{
						"value":         utils.StringPtr("10.0.0.2"),
						"prefix_length": utils.IntPtr(32),
					}},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIscsiClientsEntities(tt.in))
		})
	}
}

func TestUnitFlattenAttachmentSite(t *testing.T) {
	t.Parallel()

	primary := volumesClient.VOLUMEGROUPATTACHMENTSITE_PRIMARY
	secondary := volumesClient.VOLUMEGROUPATTACHMENTSITE_SECONDARY
	tests := []struct {
		name string
		in   *volumesClient.VolumeGroupAttachmentSite
		want string
	}{
		{name: "nil", in: nil, want: "UNKNOWN"},
		{name: "primary", in: &primary, want: "PRIMARY"},
		{name: "secondary", in: &secondary, want: "SECONDARY"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAttachmentSite(tt.in))
		})
	}
}

func TestUnitFlattenAttachedTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []volumesClient.TargetParam
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []volumesClient.TargetParam{}, want: nil},
		{
			name: "populated",
			in: []volumesClient.TargetParam{{
				NumVirtualTargets: utils.IntPtr(2),
				IscsiTargetName:   utils.StringPtr("tgt"),
			}},
			want: []interface{}{map[string]interface{}{
				"num_virtual_targets": utils.IntPtr(2),
				"iscsi_target_name":   utils.StringPtr("tgt"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenAttachedTargets(tt.in))
		})
	}
}

func TestUnitFlattenIscsiInitiatorNetworkID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.IPAddressOrFQDN
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &config.IPAddressOrFQDN{
				Fqdn: &config.FQDN{Value: utils.StringPtr("host.example.test")},
			},
			want: []interface{}{map[string]interface{}{
				"fqdn": []interface{}{map[string]interface{}{"value": utils.StringPtr("host.example.test")}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIscsiInitiatorNetworkID(tt.in))
		})
	}
}

func TestUnitFlattenIP4Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.IPv4Address
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &config.IPv4Address{Value: utils.StringPtr("10.0.0.1"), PrefixLength: utils.IntPtr(24)},
			want: []interface{}{map[string]interface{}{
				"value":         utils.StringPtr("10.0.0.1"),
				"prefix_length": utils.IntPtr(24),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIP4Address(tt.in))
		})
	}
}

func TestUnitFlattenIP6Address(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.IPv6Address
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &config.IPv6Address{Value: utils.StringPtr("::1"), PrefixLength: utils.IntPtr(128)},
			want: []interface{}{map[string]interface{}{
				"value":         utils.StringPtr("::1"),
				"prefix_length": utils.IntPtr(128),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenIP6Address(tt.in))
		})
	}
}

func TestUnitFlattenFQDN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *config.FQDN
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in:   &config.FQDN{Value: utils.StringPtr("example.test")},
			want: []interface{}{map[string]interface{}{"value": utils.StringPtr("example.test")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenFQDN(tt.in))
		})
	}
}

func TestUnitFlattenDisksEntities(t *testing.T) {
	t.Parallel()

	entityType := config.ENTITYTYPE_VM_DISK
	tests := []struct {
		name string
		in   []volumesClient.VolumeDisk
		want []interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []volumesClient.VolumeDisk{}, want: nil},
		{
			name: "populated",
			in: []volumesClient.VolumeDisk{{
				ExtId:              utils.StringPtr("vd-1"),
				TenantId:           utils.StringPtr("t-1"),
				Index:              utils.IntPtr(0),
				DiskSizeBytes:      utils.Int64Ptr(1024),
				StorageContainerId: utils.StringPtr("sc-1"),
				Description:        utils.StringPtr("disk"),
				DiskDataSourceReference: &config.EntityReference{
					ExtId:      utils.StringPtr("src-1"),
					Name:       utils.StringPtr("src"),
					Uris:       []string{"uri-1"},
					EntityType: &entityType,
				},
				DiskStorageFeatures: &volumesClient.DiskStorageFeatures{
					FlashMode: &volumesClient.FlashMode{IsEnabled: utils.BoolPtr(true)},
				},
			}},
			want: []interface{}{map[string]interface{}{
				"ext_id":               utils.StringPtr("vd-1"),
				"tenant_id":            utils.StringPtr("t-1"),
				"index":                utils.IntPtr(0),
				"disk_size_bytes":      utils.Int64Ptr(1024),
				"storage_container_id": utils.StringPtr("sc-1"),
				"description":          utils.StringPtr("disk"),
				"disk_data_source_reference": []map[string]interface{}{{
					"ext_id":      utils.StringPtr("src-1"),
					"name":        utils.StringPtr("src"),
					"uris":        []string{"uri-1"},
					"entity_type": &entityType,
				}},
				"disk_storage_features": []map[string]interface{}{{
					"flash_mode": []map[string]interface{}{{"is_enabled": utils.BoolPtr(true)}},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDisksEntities(tt.in))
		})
	}
}

func TestUnitFlattenDiskDataSourceReference(t *testing.T) {
	t.Parallel()

	entityType := config.ENTITYTYPE_VOLUME_DISK
	tests := []struct {
		name string
		in   *config.EntityReference
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &config.EntityReference{
				ExtId:      utils.StringPtr("e-1"),
				Name:       utils.StringPtr("n-1"),
				Uris:       []string{"u-1"},
				EntityType: &entityType,
			},
			want: []map[string]interface{}{{
				"ext_id":      utils.StringPtr("e-1"),
				"name":        utils.StringPtr("n-1"),
				"uris":        []string{"u-1"},
				"entity_type": &entityType,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDiskDataSourceReference(tt.in))
		})
	}
}

func TestUnitFlattenDiskStorageFeatures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *volumesClient.DiskStorageFeatures
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: nil},
		{
			name: "populated",
			in: &volumesClient.DiskStorageFeatures{
				FlashMode: &volumesClient.FlashMode{IsEnabled: utils.BoolPtr(true)},
			},
			want: []map[string]interface{}{{
				"flash_mode": []map[string]interface{}{{"is_enabled": utils.BoolPtr(true)}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, flattenDiskStorageFeatures(tt.in))
		})
	}
}

func TestUnitExpandIscsiFeatures(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIscsiFeatures([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIscsiFeatures([]interface{}{map[string]interface{}{
			"target_secret":           "secret",
			"enabled_authentications": "CHAP",
		}})
		require.NotNil(t, got)
		require.Equal(t, "secret", utils.StringValue(got.TargetSecret))
		require.Equal(t, "CHAP", got.EnabledAuthentications.GetName())
	})
}

func TestUnitExpandStorageFeatures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"flash_mode": []interface{}{map[string]interface{}{"is_enabled": true}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandStorageFeatures(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got.FlashMode)
			require.Equal(t, true, utils.BoolValue(got.FlashMode.IsEnabled))
		})
	}
}

func TestUnitExpandFlashMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{name: "populated", in: []interface{}{map[string]interface{}{"is_enabled": false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandFlashMode(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.Equal(t, false, utils.BoolValue(got.IsEnabled))
		})
	}
}

func TestUnitExpandDisks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want int
	}{
		{name: "nil", in: nil, want: 0},
		{name: "empty", in: []interface{}{}, want: 0},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"index":           1,
				"disk_size_bytes": 2048,
				"description":     "d1",
			}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandDisks(tt.in)
			if tt.want == 0 {
				require.Nil(t, got)
				return
			}
			require.Len(t, got, 1)
			require.Equal(t, 1, utils.IntValue(got[0].Index))
			require.Equal(t, int64(2048), utils.Int64Value(got[0].DiskSizeBytes))
			require.Equal(t, "d1", utils.StringValue(got[0].Description))
		})
	}
}

func TestUnitExpandDiskStorageFeatures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []interface{}
		nilOut bool
	}{
		{name: "nil", in: nil, nilOut: true},
		{name: "empty", in: []interface{}{}, nilOut: true},
		{
			name: "populated",
			in: []interface{}{map[string]interface{}{
				"flash_mode": []interface{}{map[string]interface{}{"is_enabled": true}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandDiskStorageFeatures(tt.in)
			if tt.nilOut {
				require.Nil(t, got)
				return
			}
			require.True(t, utils.BoolValue(got.FlashMode.IsEnabled))
		})
	}
}

func TestUnitExpandDiskDataSourceReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandDiskDataSourceReference(nil))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandDiskDataSourceReference([]interface{}{map[string]interface{}{
			"ext_id":      "e-1",
			"name":        "n-1",
			"uris":        []interface{}{"u-1"},
			"entity_type": "VM_DISK",
		}})
		require.NotNil(t, got)
		require.Equal(t, "e-1", utils.StringValue(got.ExtId))
		require.Equal(t, "n-1", utils.StringValue(got.Name))
		require.Equal(t, []string{"u-1"}, got.Uris)
		require.Equal(t, config.ENTITYTYPE_VM_DISK, *got.EntityType)
	})
}

func TestUnitExpandCategoryEntityReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandCategoryEntityReference(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, expandCategoryEntityReference([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandCategoryEntityReference([]interface{}{map[string]interface{}{
			"ext_id":      "cat-1",
			"name":        "AppType",
			"entity_type": "CATEGORY",
			"uris":        []interface{}{"uri-1"},
		}})
		require.Len(t, got, 1)
		require.Equal(t, "cat-1", utils.StringValue(got[0].ExtId))
		require.Equal(t, "AppType", utils.StringValue(got[0].Name))
		require.Equal(t, config.ENTITYTYPE_CATEGORY, *got[0].EntityType)
		require.Equal(t, []string{"uri-1"}, got[0].Uris)
	})
}

func TestUnitExpandCategoryEntityType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want *config.EntityType
	}{
		{name: "empty", in: "", want: nil},
		{name: "unknown", in: "VM", want: nil},
		{name: "category", in: "CATEGORY", want: config.ENTITYTYPE_CATEGORY.Ref()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := expandCategoryEntityType(tt.in)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, *tt.want, *got)
		})
	}
}

func TestUnitExpandListOfString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []interface{}
		want []string
	}{
		{name: "nil", in: nil, want: []string{}},
		{name: "empty", in: []interface{}{}, want: []string{}},
		{name: "populated", in: []interface{}{"a", "b"}, want: []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandListOfString(tt.in))
		})
	}
}

func TestUnitExpandiscsiInitiatorNetworkID(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandiscsiInitiatorNetworkID(nil))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandiscsiInitiatorNetworkID([]interface{}{map[string]interface{}{
			"ipv4": []interface{}{map[string]interface{}{"value": "10.0.0.5", "prefix_length": 32}},
			"ipv6": []interface{}{map[string]interface{}{"value": "::1", "prefix_length": 128}},
			"fqdn": []interface{}{map[string]interface{}{"value": "host.example.test"}},
		}})
		require.NotNil(t, got)
		require.Equal(t, "10.0.0.5", utils.StringValue(got.Ipv4.Value))
		require.Equal(t, "::1", utils.StringValue(got.Ipv6.Value))
		require.Equal(t, "host.example.test", utils.StringValue(got.Fqdn.Value))
	})
}

func TestUnitExpandFloatingIPV4Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandFloatingIPV4Address(nil))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandFloatingIPV4Address([]interface{}{map[string]interface{}{
			"value":         "10.0.0.9",
			"prefix_length": 24,
		}})
		require.Equal(t, "10.0.0.9", utils.StringValue(got.Value))
		require.Equal(t, 24, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandFloatingIPV6Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandFloatingIPV6Address(nil))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandFloatingIPV6Address([]interface{}{map[string]interface{}{
			"value":         "fe80::1",
			"prefix_length": 64,
		}})
		require.Equal(t, "fe80::1", utils.StringValue(got.Value))
		require.Equal(t, 64, utils.IntValue(got.PrefixLength))
	})
}

func TestUnitExpandFQDN(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandFQDN(nil))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandFQDN([]interface{}{map[string]interface{}{"value": "fqdn.example.test"}})
		require.Equal(t, "fqdn.example.test", utils.StringValue(got.Value))
	})
}

func TestUnitDiskIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want int
	}{
		{name: "nil", in: nil, want: -1},
		{name: "empty", in: map[string]interface{}{}, want: -1},
		{name: "populated", in: map[string]interface{}{"index": 3}, want: 3},
		{name: "wrong type", in: "disk", want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, diskIndex(tt.in))
		})
	}
}

func TestUnitDiffVolumeGroupDisks(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		added, removed, updated := diffVolumeGroupDisks(nil, nil)
		require.Empty(t, added)
		require.Empty(t, removed)
		require.Empty(t, updated)
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		oldDisks := []interface{}{
			map[string]interface{}{"index": 1, "description": "old"},
			map[string]interface{}{"index": 2, "description": "keep"},
		}
		newDisks := []interface{}{
			map[string]interface{}{"index": 2, "description": "keep"},
			map[string]interface{}{"index": 3, "description": "new"},
			map[string]interface{}{"index": 1, "description": "changed"},
			map[string]interface{}{"description": "no-index"},
		}
		added, removed, updated := diffVolumeGroupDisks(oldDisks, newDisks)
		require.Len(t, added, 2)
		require.Len(t, removed, 0)
		require.Len(t, updated, 1)
		require.Equal(t, 1, diskIndex(updated[0]))
	})
}
