package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
)

func TestUnitExpandFlexRuleSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandFlexRuleSpec(nil, "", nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandFlexRuleSpec(nil, "", []interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandFlexRuleSpec(nil, "", []interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIcmpV6TypeCodeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIcmpV6TypeCodeSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIcmpV6TypeCodeSpec([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIcmpV6TypeCodeSpec([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenFlexRuleSpec(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenFlexRuleSpec(import1.FlexRuleSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenFlexRuleSpec(import1.FlexRuleSpec{})
		_ = got
	})

}

func TestUnitFlattenIcmpV6TypeCodeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIcmpV6TypeCodeSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIcmpV6TypeCodeSpec([]import1.IcmpV6TypeCodeSpec{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIcmpV6TypeCodeSpec([]import1.IcmpV6TypeCodeSpec{{}})
		_ = got
	})

}
