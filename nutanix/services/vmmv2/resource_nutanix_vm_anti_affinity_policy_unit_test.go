package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
)

func TestUnitExpandPolicyCategoryReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPolicyCategoryReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandPolicyCategoryReference([]string{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandPolicyCategoryReference([]string{"a"})
		_ = got
	})

}

func TestUnitFlattenPolicyCategoryReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenPolicyCategoryReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicyCategoryReference([]import1.CategoryReference{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicyCategoryReference([]import1.CategoryReference{import1.CategoryReference{}})
		_ = got
	})

}
