package vmmv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/policies"
)

func TestUnitFlattenStartConditionConflict(t *testing.T) {
	t.Parallel()
	schemaMap := DatasourceNutanixVmStartupPolicyStartConditionConflictV2().Schema
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenStartConditionConflict(d, &import1.StartConditionConflict{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.NoError(t, flattenStartConditionConflict(d, &import1.StartConditionConflict{ExtId: utils.StringPtr("ext-1")}))
		require.Equal(t, "ext-1", d.Get("ext_id"))
	})
}

func TestUnitFlattenStartConditionSingle(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenStartConditionSingle(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenStartConditionSingle(&import1.StartCondition{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenStartConditionSingle(&import1.StartCondition{})
		_ = got
	})

}

func TestUnitFlattenStartConditionConflictToMap(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenStartConditionConflictToMap(&import1.StartConditionConflict{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenStartConditionConflictToMap(&import1.StartConditionConflict{})
		_ = got
	})

}
