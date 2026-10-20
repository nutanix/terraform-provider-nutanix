package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	catalogConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/catalogCommon/v1/config"
)

func TestUnitExpandCategoriesFilter(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandCategoriesFilter(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandCategoriesFilter([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenCategoriesFilter(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenCategoriesFilter(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoriesFilter(&catalogConfig.CategoriesFilter{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoriesFilter(&catalogConfig.CategoriesFilter{})
		_ = got
	})

}
