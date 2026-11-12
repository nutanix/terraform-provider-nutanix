package networkingv2

import (
	"testing"
)

func TestUnitExpandTCPPortRangeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandTCPPortRangeSpec(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandTCPPortRangeSpec([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandTCPPortRangeSpec([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandUDPPortRangeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandUDPPortRangeSpec(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandUDPPortRangeSpec([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandUDPPortRangeSpec([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIcmpTypeCodeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIcmpTypeCodeSpec(nil, "", nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIcmpTypeCodeSpec(nil, "", []interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIcmpTypeCodeSpec(nil, "icmp_type_code_specs", []interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
