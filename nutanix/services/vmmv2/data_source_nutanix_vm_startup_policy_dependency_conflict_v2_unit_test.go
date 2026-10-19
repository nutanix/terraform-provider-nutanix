package vmmv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
)

func TestUnitFlattenDependencyConflict(t *testing.T) {
	t.Parallel()
	schemaMap := DatasourceNutanixVmStartupPolicyDependencyConflictV2().Schema
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenDependencyConflict(d, &import1.DependencyConflict{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenDependencyConflict(d, &import1.DependencyConflict{ExtId: utils.StringPtr("ext-1")}))
		require.Equal(t, "ext-1", d.Get("ext_id"))
	})
}

func TestUnitFlattenCategoryRef(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenCategoryRef(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryRef(&import1.CategoryReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryRef(&import1.CategoryReference{})
		_ = got
	})

}

func TestUnitFlattenCategoryRefList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenCategoryRefList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryRefList([]import1.CategoryReference{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryRefList([]import1.CategoryReference{import1.CategoryReference{}})
		_ = got
	})

}

func TestUnitFlattenCategoryDependencyChain(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenCategoryDependencyChain(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryDependencyChain([]import1.CategoryDependency{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryDependencyChain([]import1.CategoryDependency{import1.CategoryDependency{}})
		_ = got
	})

}

func TestUnitFlattenDependencyConflictToMap(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDependencyConflictToMap(&import1.DependencyConflict{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDependencyConflictToMap(&import1.DependencyConflict{})
		_ = got
	})

}
