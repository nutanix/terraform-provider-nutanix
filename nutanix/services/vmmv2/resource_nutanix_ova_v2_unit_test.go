package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/iam/v4/authn"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/content"
)

func TestUnitExpandOneOfOvaChecksum(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfOvaChecksum([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenOneOfOvaChecksum(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfOvaChecksum(nil))
	})
}

func TestUnitExpandOneOfOvaSource(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfOvaSource([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenOneOfOvaSource(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfOvaSource(nil))
	})
}

func TestUnitExpandOvaURLBasicAuth(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOvaURLBasicAuth([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenOvaDiskFormat(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenOvaDiskFormat(func() *import1.OvaDiskFormat { var v import1.OvaDiskFormat; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenOvaDiskFormat(func() *import1.OvaDiskFormat { v := import1.OvaDiskFormat(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenCreatedBy(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCreatedBy(&import4.User{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCreatedBy(&import4.User{})
		_ = got
	})

}

func TestUnitFlattenCreationType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCreationType(func() *import4.CreationType { var v import4.CreationType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCreationType(func() *import4.CreationType { v := import4.CreationType(1); return &v }())
		_ = got
	})

}
