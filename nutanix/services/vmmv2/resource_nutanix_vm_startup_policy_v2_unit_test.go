package vmmv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/common/v1/response"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
)

func TestUnitExpandVmStartupPolicyGroups(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVmStartupPolicyGroups(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVmStartupPolicyGroups([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVmStartupPolicyGroups([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandVmStartupPolicyStartConditions(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVmStartupPolicyStartConditions(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVmStartupPolicyStartConditions([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVmStartupPolicyStartConditions([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenVmStartupPolicy(t *testing.T) {
	t.Parallel()
	schemaMap := ResourceNutanixVmStartupPolicyV2().Schema
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenVmStartupPolicy(d, &import1.VmStartupPolicy{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenVmStartupPolicy(d, &import1.VmStartupPolicy{ExtId: utils.StringPtr("ext-1")}))
		require.Equal(t, "ext-1", d.Get("ext_id"))
	})
}

func TestUnitFlattenUserReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenUserReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenUserReference(&import1.UserReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenUserReference(&import1.UserReference{})
		_ = got
	})

}

func TestUnitFlattenVmStartupPolicyGroups(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmStartupPolicyGroups(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmStartupPolicyGroups([]import1.DependencyGroup{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmStartupPolicyGroups([]import1.DependencyGroup{import1.DependencyGroup{}})
		_ = got
	})

}

func TestUnitFlattenVmStartupPolicyStartConditions(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmStartupPolicyStartConditions(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmStartupPolicyStartConditions([]import1.StartCondition{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmStartupPolicyStartConditions([]import1.StartCondition{import1.StartCondition{}})
		_ = got
	})

}

func TestUnitFlattenVmStartupPolicyLinks(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmStartupPolicyLinks(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmStartupPolicyLinks([]response.ApiLink{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmStartupPolicyLinks([]response.ApiLink{response.ApiLink{}})
		_ = got
	})

}
