package vmmv2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitIsVmmEtagMismatchErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "other", err: errors.New("connection reset"), want: false},
		{name: "if-match", err: errors.New("If-Match header value passed"), want: true},
		{name: "etag mismatch", err: errors.New("VM_ETAG_MISMATCH"), want: true},
		{name: "vmm code", err: errors.New("VMM-30303"), want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, isVmmEtagMismatchErr(tt.err))
		})
	}
}

func TestUnitStripDataSourceFromDiskBackingInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.NotPanics(t, func() { StripDataSourceFromDiskBackingInfo(nil) })
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.NotPanics(t, func() { StripDataSourceFromDiskBackingInfo(map[string]interface{}{}) })
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		disk := map[string]interface{}{
			"backing_info": []interface{}{
				map[string]interface{}{
					"vm_disk": []interface{}{
						map[string]interface{}{
							"data_source": []interface{}{map[string]interface{}{"ext_id": "ds-1"}},
							"ext_id":      "disk-1",
						},
					},
				},
			},
		}
		StripDataSourceFromDiskBackingInfo(disk)
		vmDisk := disk["backing_info"].([]interface{})[0].(map[string]interface{})["vm_disk"].([]interface{})[0].(map[string]interface{})
		_, ok := vmDisk["data_source"]
		require.False(t, ok)
		require.Equal(t, "disk-1", vmDisk["ext_id"])
	})
}

func TestUnitHasNonEmptyBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		m    map[string]interface{}
		key  string
		want bool
	}{
		{name: "nil map", m: nil, key: "x", want: false},
		{name: "empty", m: map[string]interface{}{}, key: "x", want: false},
		{name: "populated", m: map[string]interface{}{"x": []interface{}{map[string]interface{}{"a": 1}}}, key: "x", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, hasNonEmptyBlock(tt.m, tt.key))
		})
	}
}

func TestUnitVolumeGroupExtIDFromBackingInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", volumeGroupExtIDFromBackingInfo(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", volumeGroupExtIDFromBackingInfo([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := volumeGroupExtIDFromBackingInfo([]interface{}{
			map[string]interface{}{
				"adfs_volume_group_reference": []interface{}{
					map[string]interface{}{"volume_group_ext_id": "vg-1"},
				},
			},
		})
		require.Equal(t, "vg-1", got)
	})
}
